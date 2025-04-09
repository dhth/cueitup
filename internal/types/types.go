package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/dhth/cueitup/internal/utils"
)

const (
	typeJSON = "json"
	typeNone = "none"
)

var (
	errIncorrectMessageFmtProvided      = errors.New("encoding format is incorrect")
	errIncorrectQueueURLProvided        = errors.New("queue URL is incorrect")
	errConfigSourceEmpty                = errors.New("config source is empty")
	errContextKeyCannotBeUsed           = errors.New("context key can only be used when message format is JSON")
	errSubsetKeyCannotBeUsed            = errors.New("subset key can only be used when message format is JSON")
	errContextKeyEmpty                  = errors.New("context key is empty")
	errSubsetKeyEmpty                   = errors.New("subset key is empty")
	errUnexpectedTypeAfterUnmarshalling = errors.New("unexpected type found after unmarshalling JSON")
	errMessageIDNil                     = errors.New("message ID is null")
	errMessageBodyNil                   = errors.New("message body is null")
	errCouldntUnmarshalBytes            = errors.New("couldn't unmarshal message body bytes as JSON")
	errNestedObjectCouldntBeAccessed    = errors.New("nested object couln't be accessed")
	errCouldntUnmarshalNestedObject     = errors.New("couldn't unmarshall nested object")
	errCouldntMarshalNestedObject       = errors.New("couldn't marshall nested object")
	errContextKeyNotFound               = errors.New("context key not found")
	errContextValueIsNotAString         = errors.New("context value is not a string")
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
		value = typeJSON
	case None:
		value = typeNone
	}

	return value
}

func (f MessageFormat) Extension() string {
	var value string
	switch f {
	case JSON:
		value = typeJSON
	case None:
		value = "txt"
	}

	return value
}

type Config struct {
	ProfileName     string        `json:"profile_name"`
	QueueURL        string        `json:"queue_url"`
	AWSConfigSource string        `json:"aws_config_source"`
	Format          MessageFormat `json:"-"`
	ContextKey      *string       `json:"context_key"`
	SubsetKey       *string       `json:"subset_key"`
}

func (p Config) Display() string {
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
			p.ProfileName,
			p.QueueURL,
			p.AWSConfigSource,
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
			p.ProfileName,
			p.QueueURL,
			p.AWSConfigSource,
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

	return fmt.Errorf("%w: %q", errIncorrectQueueURLProvided, pc.QueueURL)
}

func (pc *ProfileConfig) validateConfigSource() error {
	if strings.TrimSpace(pc.AWSConfigSource) == "" {
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

func ParseProfileConfig(config ProfileConfig) (Config, []error) {
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
		return Config{}, errors
	}

	return Config{
		ProfileName:     config.Name,
		QueueURL:        config.QueueURL,
		AWSConfigSource: config.AWSConfigSource,
		Format:          msgFmt,
		ContextKey:      config.ContextKey,
		SubsetKey:       config.SubsetKey,
	}, nil
}

type TUIBehaviours struct {
	DeleteMessages  bool
	PersistMessages bool
	SkipMessages    bool
}

func (b TUIBehaviours) Display() string {
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

type WebBehaviours struct {
	DeleteMessages bool `json:"delete_messages"`
	SelectOnHover  bool `json:"select_on_hover"`
	ShowLiveCount  bool `json:"show_live_count"`
}

func (b WebBehaviours) Display() string {
	return fmt.Sprintf(`
- delete messages         %v
- select on hover         %v
- show live count         %v
`,
		b.DeleteMessages,
		b.SelectOnHover,
		b.ShowLiveCount,
	)
}

type Message struct {
	ID           string  `json:"id"`
	Body         string  `json:"body"`
	ContextKey   *string `json:"context_key"`
	ContextValue *string `json:"context_value"`
	Err          error   `json:"error"`
}

func (item Message) Title() string {
	return fmt.Sprintf("%s: %s", utils.RightPadTrim("msgId", 12), item.ID)
}

func (item Message) Description() string {
	if item.ContextKey != nil && item.ContextValue != nil {
		return fmt.Sprintf("%s: %s", utils.RightPadTrim(*item.ContextKey, 12), *item.ContextValue)
	}
	return ""
}

func (item Message) FilterValue() string {
	return item.ID
}

func GetMessageData(message *sqstypes.Message, config Config) Message {
	switch config.Format {
	case JSON:
		return getJSONMessage(message, config.SubsetKey, config.ContextKey)
	default:
		return getPlainMessage(message)
	}
}

// TODO: improve this
func getJSONMessage(message *sqstypes.Message, subsetKey *string, contextKey *string) Message {
	if message.MessageId == nil {
		return Message{
			Err: errMessageIDNil,
		}
	}
	if message.Body == nil {
		return Message{
			Err: errMessageBodyNil,
		}
	}

	var data map[string]any
	messageID := *message.MessageId
	bodyBytes := []byte(*message.Body)
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return Message{
			Err: fmt.Errorf("%w: %s", errCouldntUnmarshalBytes, err.Error()),
		}
	}

	if subsetKey == nil && contextKey == nil {
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return Message{
				Err: fmt.Errorf("%w: %s", errCouldntMarshalNestedObject, err.Error()),
			}
		}

		return Message{
			ID:   messageID,
			Body: string(bytes),
		}
	}

	if subsetKey == nil && contextKey != nil {
		context, ok := data[*contextKey]
		if !ok {
			return Message{
				ContextKey: contextKey,
				Err:        errContextKeyNotFound,
			}
		}

		var contextValue string
		switch c := context.(type) {
		case string:
			contextValue = c
		default:
			return Message{
				ContextKey: contextKey,
				Err:        fmt.Errorf("%w (key: %s): %v", errContextValueIsNotAString, *contextKey, context),
			}
		}

		return Message{
			ID:           messageID,
			Body:         string(bodyBytes),
			ContextKey:   contextKey,
			ContextValue: &contextValue,
		}
	}

	if subsetKey != nil && contextKey == nil {
		sKey := *subsetKey
		subset, ok := data[sKey]
		if !ok {
			return Message{
				Err: fmt.Errorf("%w (key: %s)", errNestedObjectCouldntBeAccessed, sKey),
			}
		}

		var subsetData map[string]any
		switch n := subset.(type) {
		case map[string]any:
			// It's a JSON object, directly access the nested key
			subsetData = n
		case string:
			// May be stringified JSON; attempt to convert it to JSON
			if err := json.Unmarshal([]byte(n), &subsetData); err != nil {
				return Message{
					Err: fmt.Errorf("%w (key: %s): %s", errCouldntUnmarshalNestedObject, sKey, err.Error()),
				}
			}
		default:
			return Message{
				Err: fmt.Errorf("%w; data: %v", errUnexpectedTypeAfterUnmarshalling, subsetData),
			}
		}

		subsetBytes, err := json.MarshalIndent(subsetData, "", "  ")
		if err != nil {
			return Message{
				Err: fmt.Errorf("%w: %s", errCouldntMarshalNestedObject, err.Error()),
			}
		}

		return Message{
			ID:   messageID,
			Body: string(subsetBytes),
		}
	}

	sKey := *subsetKey
	subset, ok := data[sKey]
	if !ok {
		return Message{
			ContextKey: contextKey,
			Err:        fmt.Errorf("%w (key: %s)", errNestedObjectCouldntBeAccessed, sKey),
		}
	}

	var subsetData map[string]any
	switch n := subset.(type) {
	case map[string]any:
		// It's a JSON object, directly access the nested key
		subsetData = n
	case string:
		// May be stringified JSON; attempt to convert it to JSON
		if err := json.Unmarshal([]byte(n), &subsetData); err != nil {
			return Message{
				ContextKey: contextKey,
				Err:        fmt.Errorf("%w (key: %s): %s", errCouldntUnmarshalNestedObject, sKey, err.Error()),
			}
		}
	default:
		return Message{
			ContextKey: contextKey,
			Err:        fmt.Errorf("%w; data: %v", errUnexpectedTypeAfterUnmarshalling, subsetData),
		}
	}

	subsetBytes, err := json.MarshalIndent(subsetData, "", "  ")
	if err != nil {
		return Message{
			ContextKey: contextKey,
			Err:        fmt.Errorf("%w: %s", errCouldntMarshalNestedObject, err.Error()),
		}
	}
	context, ok := subsetData[*contextKey]
	if !ok {
		return Message{
			ContextKey: contextKey,
			Err:        errContextKeyNotFound,
		}
	}

	var contextValue *string
	switch c := context.(type) {
	case string:
		contextValue = &c
	default:
		return Message{
			ContextKey: contextKey,
			Err:        fmt.Errorf("%w (key: %s): %v", errContextValueIsNotAString, *contextKey, context),
		}
	}

	return Message{
		ID:           messageID,
		Body:         string(subsetBytes),
		ContextKey:   contextKey,
		ContextValue: contextValue,
	}
}

func getPlainMessage(message *sqstypes.Message) Message {
	if message.MessageId == nil {
		return Message{
			Err: errMessageIDNil,
		}
	}
	if message.Body == nil {
		return Message{
			Err: errMessageBodyNil,
		}
	}

	return Message{
		ID:   *message.MessageId,
		Body: *message.Body,
	}
}
