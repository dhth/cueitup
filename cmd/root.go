package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/dhth/cueitup/internal/aws"
	"github.com/dhth/cueitup/internal/server"
	t "github.com/dhth/cueitup/internal/types"
	"github.com/dhth/cueitup/internal/ui"
	"github.com/dhth/cueitup/internal/utils"
	"github.com/spf13/cobra"
)

const (
	configFileName = "cueitup/cueitup.yml"
)

var (
	errConfigFileNotYAML       = errors.New("config needs to be a YAML file")
	errCouldntLoadAWSConfig    = errors.New("couldn't load AWS config")
	errCouldntGetUserHomeDir   = errors.New("couldn't get your home directory")
	errCouldntGetUserConfigDir = errors.New("couldn't get your config directory")
	ErrCouldntReadConfigFile   = errors.New("couldn't read config file")
)

func Execute() error {
	rootCmd, err := NewRootCommand()
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func NewRootCommand() (*cobra.Command, error) {
	var (
		configPath       string
		configPathFull   string
		configBytes      []byte
		homeDir          string
		deleteMessages   bool
		persistMessages  bool
		skipMessages     bool
		selectOnHover    bool
		showMessageCount bool
		webOpen          bool
		debug            bool
		listConfig       bool
	)

	rootCmd := &cobra.Command{
		Use:   "cueitup",
		Short: "cueitup lets you inspect messages in an AWS SQS queue in a simple and deliberate manner",
		Long: `cueitup lets you inspect messages in an AWS SQS queue in a simple and deliberate manner.

cueitup relies on a configuration file that contains profiles for various SQS topics, each with its
own details related to authentication, deserialization, etc.
`,
		SilenceUsage: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if !strings.HasSuffix(configPath, ".yml") && !strings.HasSuffix(configPath, ".yaml") {
				return errConfigFileNotYAML
			}

			var err error
			configPathFull = utils.ExpandTilde(configPath, homeDir)
			configBytes, err = os.ReadFile(configPathFull)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
			}

			return nil
		},
	}

	configCmd := &cobra.Command{
		Use:          "config",
		Short:        "interact with cueitup's config",
		SilenceUsage: true,
	}

	validateConfigCmd := &cobra.Command{
		Use:          "validate",
		Short:        "validate cueitup's config",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			if listConfig {
				fmt.Printf("%s\n---\n\n", configBytes)
			}
			errors := validateConfig(configBytes)
			if len(errors) > 0 {
				fmt.Println("config has some errors:")
				for _, err := range errors {
					fmt.Printf("- %s\n", err.Error())
				}
				// nolint:revive
				os.Exit(1)
			}

			fmt.Println("config looks good ✅")
			return nil
		},
	}

	tuiCmd := &cobra.Command{
		Use:          "tui <PROFILE>",
		Short:        "open cueitup's TUI",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := getConfig(configBytes, args[0])
			if err != nil {
				return err
			}

			behaviours := t.TUIBehaviours{
				DeleteMessages:   deleteMessages,
				PersistMessages:  persistMessages,
				SkipMessages:     skipMessages,
				ShowMessageCount: showMessageCount,
			}

			if debug {
				fmt.Printf(`Debug info:
===

Profile
---
%s
Behaviours 
---
%s`,
					cfg.Display(),
					behaviours.Display(),
				)
				return nil
			}

			sdkConfig, err := aws.GetAWSConfig(cfg.AWSConfigSource)
			if err != nil {
				return fmt.Errorf("%w: %s", errCouldntLoadAWSConfig, err.Error())
			}

			sqsClient := sqs.NewFromConfig(sdkConfig)

			return ui.RenderUI(sqsClient, cfg.QueueURL, cfg, behaviours)
		},
	}

	serveCmd := &cobra.Command{
		Use:          "serve <PROFILE>",
		Short:        "open cueitup's web interface",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := getConfig(configBytes, args[0])
			if err != nil {
				return err
			}

			behaviours := t.WebBehaviours{
				DeleteMessages:   deleteMessages,
				SelectOnHover:    selectOnHover,
				ShowMessageCount: showMessageCount,
			}

			if debug {
				fmt.Printf(`Debug info:
===

Profile
---
%s
Behaviours 
---
%s`,
					cfg.Display(),
					behaviours.Display(),
				)
				return nil
			}

			sdkConfig, err := aws.GetAWSConfig(cfg.AWSConfigSource)
			if err != nil {
				return fmt.Errorf("%w: %s", errCouldntLoadAWSConfig, err.Error())
			}

			sqsClient := sqs.NewFromConfig(sdkConfig)

			return server.Serve(sqsClient, cfg, behaviours, webOpen)
		},
	}

	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserHomeDir, err.Error())
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserConfigDir, err.Error())
	}

	defaultConfigPath := filepath.Join(configDir, configFileName)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of cueitup's config file")

	tuiCmd.Flags().BoolVarP(&debug, "debug", "d", false, "whether to only display config picked up by cueitup")
	tuiCmd.Flags().BoolVarP(&deleteMessages, "delete-messages", "D", true, "whether to start the TUI with the setting \"delete messages\" ON")
	tuiCmd.Flags().BoolVarP(&persistMessages, "persist-messages", "P", false, "whether to start the TUI with the setting \"persist messages\" ON")
	tuiCmd.Flags().BoolVarP(&skipMessages, "skip-messages", "S", false, "whether to start the TUI with the setting \"skip messages\" ON")
	tuiCmd.Flags().BoolVarP(&showMessageCount, "show-message-count", "M", true, "whether to start the TUI with the setting \"show message count\" ON")

	serveCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of cueitup's config file")
	serveCmd.Flags().BoolVarP(&deleteMessages, "delete-messages", "D", true, "whether to start the web interface with the setting \"delete messages\" ON")
	serveCmd.Flags().BoolVarP(&selectOnHover, "select-on-hover", "S", false, "whether to start the web interface with the setting \"select on hover\" ON")
	serveCmd.Flags().BoolVarP(&showMessageCount, "show-message-count", "M", true, "whether to start the web interface with the setting \"show message count\" ON")
	serveCmd.Flags().BoolVarP(&webOpen, "open", "o", false, "whether to open web interface in browser automatically")
	serveCmd.Flags().BoolVarP(&debug, "debug", "d", false, "whether to only display config picked up by cueitup")

	validateConfigCmd.Flags().BoolVarP(&listConfig, "list", "l", false, "whether to list the config as well")
	configCmd.AddCommand(validateConfigCmd)

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(serveCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}
