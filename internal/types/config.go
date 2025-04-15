package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	cfgSrcSharedProfilePrefix = "profile:"
	notProvided               = "<NOT PROVIDED>"
)

type ConfigSourceKind uint

const (
	Env ConfigSourceKind = iota
	SharedProfile
)

type ConfigSource struct {
	Kind  ConfigSourceKind
	Value string
}

func (cs ConfigSource) Display() string {
	return cs.Value
}

func (cs ConfigSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs.Value)
}

var (
	errProfileNameEmpty            = errors.New("profile name is empty")
	errIncorrectMessageFmtProvided = errors.New("encoding format is incorrect")
	errIncorrectQueueURLProvided   = errors.New("queue URL is incorrect")
	errConfigSourceEmpty           = errors.New("config source is empty")
	errIncorrectConfigSource       = errors.New("incorrect config source provided")
	errContextKeyCannotBeUsed      = errors.New("context key can only be used when message format is JSON")
	errSubsetKeyCannotBeUsed       = errors.New("subset key can only be used when message format is JSON")
	errContextKeyEmpty             = errors.New("context key is empty")
	errSubsetKeyEmpty              = errors.New("subset key is empty")
)

type Config struct {
	ProfileName     string        `json:"profile_name"`
	QueueURL        string        `json:"queue_url"`
	AWSConfigSource ConfigSource  `json:"aws_config_source"`
	Format          MessageFormat `json:"-"`
	ContextKey      *string       `json:"context_key"`
	SubsetKey       *string       `json:"subset_key"`
}

func (p Config) Display() string {
	var value string

	switch p.Format {
	case JSON:
		var contextKey string
		if p.ContextKey != nil {
			contextKey = *p.ContextKey
		} else {
			contextKey = notProvided
		}
		var subsetKey string
		if p.SubsetKey != nil {
			subsetKey = *p.SubsetKey
		} else {
			subsetKey = notProvided
		}
		value = fmt.Sprintf(`
- name                    %s
- queue URL               %s
- AWS config source       %s
- format                  %v
- context key             %s
- subset key              %s
        `,
			p.ProfileName,
			p.QueueURL,
			p.AWSConfigSource.Display(),
			p.Format.Display(),
			contextKey,
			subsetKey,
		)
	case None:
		value = fmt.Sprintf(`
- name                    %s
- queue URL               %s
- AWS config source       %s
- format                  %v
        `,
			p.ProfileName,
			p.QueueURL,
			p.AWSConfigSource.Display(),
			p.Format.Display(),
		)
	}

	return value
}

type CueitupConfig struct {
	Profiles []ProfileConfig
}

type ProfileConfig struct {
	Name            string  `yaml:"name"`
	QueueURL        string  `yaml:"queue_url"`
	AWSConfigSource string  `yaml:"aws_config_source"`
	Format          string  `yaml:"format"`
	ContextKey      *string `yaml:"context_key"`
	SubsetKey       *string `yaml:"subset_key"`
}

func (pc *ProfileConfig) validateProfileName() (string, error) {
	var zero string
	if len(strings.TrimSpace(pc.Name)) == 0 {
		return zero, errProfileNameEmpty
	}

	return strings.TrimSpace(pc.Name), nil
}

func (pc *ProfileConfig) validateMessageFormat() (MessageFormat, error) {
	switch pc.Format {
	case typeJSON:
		return JSON, nil
	case typeNone:
		return None, nil
	default:
		return JSON, fmt.Errorf("%w: %q; possible values: [%s, %s]", errIncorrectMessageFmtProvided, pc.Format, typeJSON, typeNone)
	}
}

func (pc *ProfileConfig) validateQueueURL() error {
	if strings.HasPrefix(pc.QueueURL, "https://") {
		return nil
	}

	return fmt.Errorf("%w (%q): needs to be a proper URL", errIncorrectQueueURLProvided, pc.QueueURL)
}

func (pc *ProfileConfig) validateContextKey(format MessageFormat) error {
	if format != JSON && pc.ContextKey != nil {
		return errContextKeyCannotBeUsed
	}

	if pc.ContextKey != nil && strings.TrimSpace(*pc.ContextKey) == "" {
		return errContextKeyEmpty
	}

	return nil
}

func (pc *ProfileConfig) validateSubsetKey(format MessageFormat) error {
	if format != JSON && pc.SubsetKey != nil {
		return errSubsetKeyCannotBeUsed
	}

	if pc.SubsetKey != nil && strings.TrimSpace(*pc.SubsetKey) == "" {
		return errSubsetKeyEmpty
	}

	return nil
}

func ParseProfileConfig(config ProfileConfig) (Config, []error) {
	var errors []error

	profileName, err := config.validateProfileName()
	if err != nil {
		errors = append(errors, err)
	}

	msgFmt, err := config.validateMessageFormat()
	if err != nil {
		errors = append(errors, err)
	}

	err = config.validateQueueURL()
	if err != nil {
		errors = append(errors, err)
	}

	cfgSrc, err := parseConfigSource(config.AWSConfigSource)
	if err != nil {
		errors = append(errors, err)
	}

	err = config.validateContextKey(msgFmt)
	if err != nil {
		errors = append(errors, err)
	}

	err = config.validateSubsetKey(msgFmt)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return Config{}, errors
	}

	return Config{
		ProfileName:     profileName,
		QueueURL:        config.QueueURL,
		AWSConfigSource: cfgSrc,
		Format:          msgFmt,
		ContextKey:      config.ContextKey,
		SubsetKey:       config.SubsetKey,
	}, nil
}

func parseConfigSource(value string) (ConfigSource, error) {
	var zero ConfigSource
	if strings.TrimSpace(value) == "" {
		return zero, errConfigSourceEmpty
	}

	if value == "env" {
		return ConfigSource{Env, "env"}, nil
	}

	if strings.HasPrefix(value, cfgSrcSharedProfilePrefix) {
		value := strings.TrimPrefix(value, cfgSrcSharedProfilePrefix)
		if strings.TrimSpace(value) == "" {
			return zero, errConfigSourceEmpty
		}
		return ConfigSource{
			SharedProfile,
			value,
		}, nil
	}

	return zero, fmt.Errorf(`%w; possible values: "env", "profile:<aws-shared-config-profile-name>", "assume:<arn-of-role-to-assume>"`, errIncorrectConfigSource)
}
