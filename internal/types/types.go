package types

import (
	"encoding/json"
	"errors"
	"fmt"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/dhth/cueitup/internal/utils"
)

const (
	typeJSON = "json"
	typeNone = "none"
)

var (
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
