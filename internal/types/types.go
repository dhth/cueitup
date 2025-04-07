package config

import (
	"errors"
	"fmt"
	"strings"
)

const (
	json = "json"
	none = "none"
)

var (
	errIncorrectMessageFmtProvided = errors.New("encoding format is incorrect")
	errIncorrectQueueURLProvided   = errors.New("queue URL is incorrect")
	errConfigSourceEmpty           = errors.New("config source is empty")
	errContextKeyCannotBeUsed      = errors.New("context key can only be used when message format is JSON")
	errSubsetKeyCannotBeUsed       = errors.New("subset key can only be used when message format is JSON")
	errContextKeyEmpty             = errors.New("context key is empty")
	errSubsetKeyEmpty              = errors.New("subset key is empty")
)

type MessageFormat uint

const (
	JSON MessageFormat = iota
	None
)

func (f MessageFormat) Display() string {
	var value string
	switch f {
	case JSON:
		value = json
	case None:
		value = none
	}

	return value
}

func (f MessageFormat) Extension() string {
	var value string
	switch f {
	case JSON:
		value = json
	case None:
		value = "txt"
	}

	return value
}

type Profile struct {
	Name         string
	QueueURL     string
	ConfigSource string
	Format       MessageFormat
	ContextKey   string
	SubsetKey    string
}

func (p Profile) Display() string {
	var value string
	switch p.Format {
	case JSON:
		value = fmt.Sprintf(`
- name                    %s
- queue URL               %s
- config source           %s
- format                  %v
- context key             %v
- subset key              %v
        `,
			p.Name,
			p.QueueURL,
			p.ConfigSource,
			p.Format.Display(),
			p.ContextKey,
			p.SubsetKey,
		)
	case None:
		value = fmt.Sprintf(`
- name                    %s
- queue URL               %s
- config source           %s
- format                  %v
        `,
			p.Name,
			p.QueueURL,
			p.ConfigSource,
			p.Format.Display(),
		)
	}

	return value
}

type Config struct {
	Profiles []ProfileConfig
}

type ProfileConfig struct {
	Name         string  `yaml:"name"`
	QueueURL     string  `yaml:"queue_url"`
	ConfigSource string  `yaml:"config_source"`
	Format       string  `yaml:"format"`
	ContextKey   *string `yaml:"context_key"`
	SubsetKey    *string `yaml:"subset_key"`
}

func (pc *ProfileConfig) validateMessageFormat() (MessageFormat, error) {
	switch pc.Format {
	case json:
		return JSON, nil
	case none:
		return None, nil
	default:
		return JSON, fmt.Errorf("%w: %q; possible values: [%s, %s]", errIncorrectMessageFmtProvided, pc.Format, json, none)
	}
}

func (pc *ProfileConfig) validateQueueURL() error {
	if strings.HasPrefix(pc.QueueURL, "https://") {
		return nil
	}

	return fmt.Errorf("%w: %q", errIncorrectQueueURLProvided, pc.QueueURL)
}

func (pc *ProfileConfig) validateConfigSource() error {
	if strings.TrimSpace(pc.ConfigSource) == "" {
		return errConfigSourceEmpty
	}

	return nil
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

func ParseProfileConfig(config ProfileConfig) (Profile, []error) {
	var errors []error

	msgFmt, err := config.validateMessageFormat()
	if err != nil {
		errors = append(errors, err)
	}

	err = config.validateQueueURL()
	if err != nil {
		errors = append(errors, err)
	}

	err = config.validateConfigSource()
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
		return Profile{}, errors
	}

	var ck string
	if config.ContextKey != nil {
		ck = *config.ContextKey
	}

	var sk string
	if config.SubsetKey != nil {
		sk = *config.SubsetKey
	}

	return Profile{
		Name:         config.Name,
		QueueURL:     config.QueueURL,
		ConfigSource: config.ConfigSource,
		Format:       msgFmt,
		ContextKey:   ck,
		SubsetKey:    sk,
	}, nil
}

type Behaviours struct {
	DeleteMessages  bool
	PersistMessages bool
	SkipMessages    bool
}

func (b Behaviours) Display() string {
	return fmt.Sprintf(`
- delete messages         %v
- persist messages        %v
- skip messages           %v
`,
		b.DeleteMessages,
		b.PersistMessages,
		b.SkipMessages,
	)
}
