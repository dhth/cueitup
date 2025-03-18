package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	t "github.com/dhth/cueitup/internal/types"
	"github.com/dhth/cueitup/internal/ui"
	"github.com/dhth/cueitup/internal/utils"
	"github.com/spf13/cobra"
)

const (
	configFileName = "cueitup/cueitup.yml"
)

var (
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
		configPath      string
		configPathFull  string
		homeDir         string
		profile         t.Profile
		deleteMessages  bool
		persistMessages bool
		skipMessages    bool
		debug           bool
	)

	rootCmd := &cobra.Command{
		Use:   "cueitup <PROFILE>",
		Short: "cueitup lets you inspect messages in an AWS SQS queue in a simple and deliberate manner",
		Long: `cueitup lets you inspect messages in an AWS SQS queue in a simple and deliberate manner.

cueitup relies on a configuration file that contains profiles for various SQS topics, each with its
own details related to authentication, deserialization, etc.
`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		PersistentPreRunE: func(_ *cobra.Command, args []string) error {
			configPathFull = utils.ExpandTilde(configPath, homeDir)
			configBytes, err := os.ReadFile(configPathFull)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
			}

			profile, err = getProfile(configBytes, args[0])
			if errors.Is(err, errProfileNotFound) {
				return err
			} else if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			behaviours := t.Behaviours{
				DeleteMessages:  deleteMessages,
				PersistMessages: persistMessages,
				SkipMessages:    skipMessages,
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
					profile.Display(),
					behaviours.Display(),
				)
				return nil
			}

			sdkConfig, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedConfigProfile(profile.ConfigSource),
			)
			if err != nil {
				return fmt.Errorf("%w: %s", errCouldntLoadAWSConfig, err.Error())
			}

			sqsClient := sqs.NewFromConfig(sdkConfig)

			return ui.RenderUI(sqsClient, profile.QueueURL, profile, behaviours)
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

	rootCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of cueitup's config file")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "whether to only display config picked up by cueitup")
	rootCmd.Flags().BoolVarP(&deleteMessages, "delete-messages", "D", false, "whether to start the TUI with the setting \"delete messages\" ON")
	rootCmd.Flags().BoolVarP(&persistMessages, "persist-messages", "P", false, "whether to start the TUI with the setting \"persist messages\" ON")
	rootCmd.Flags().BoolVarP(&skipMessages, "skip-messages", "S", false, "whether to start the TUI with the setting \"skip messages\" ON")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}
