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
	errMessageIDNil                = errors.New("message ID is null")
	errMessageBodyNil              = errors.New("message body is null")
	errCouldntUnmarshalBytes       = errors.New("couldn't unmarshal message body bytes as JSON")
	errSubsetKeyNotFound           = errors.New("subset key not found in message body")
	errCouldntUnmarshalSubsetValue = errors.New("couldn't unmarshal subset value")
	errSubsetTypeIsUnsupported     = errors.New("subset type is unsupported")
	errCouldntMarshalBytes         = errors.New("couldn't marshall nested object")
	errContextKeyNotFound          = errors.New("context key not found in message body")
	errContextValueTypeUnsupported = errors.New("context value type is unsupported")
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
	DeleteMessages   bool
	PersistMessages  bool
	ShowMessageCount bool
	SkipMessages     bool
}

func (b TUIBehaviours) Display() string {
	return fmt.Sprintf(`
- delete messages         %v
- persist messages        %v
- show message count      %v
- skip messages           %v
`,
		b.DeleteMessages,
		b.PersistMessages,
		b.ShowMessageCount,
		b.SkipMessages,
	)
}

type WebBehaviours struct {
	DeleteMessages   bool `json:"delete_messages"`
	SelectOnHover    bool `json:"select_on_hover"`
	ShowMessageCount bool `json:"show_message_count"`
}

func (b WebBehaviours) Display() string {
	return fmt.Sprintf(`
- delete messages         %v
- select on hover         %v
- show message count      %v
`,
		b.DeleteMessages,
		b.SelectOnHover,
		b.ShowMessageCount,
	)
}

type Message struct {
	ID           string  `json:"id"`
	Body         string  `json:"body"`
	ContextKey   *string `json:"context_key"`
	ContextValue *string `json:"context_value"`
	Err          error   `json:"-"`
}

type SerializableMessage struct {
	Message
	Err *string `json:"error"`
}

func (m Message) ToSerializable() SerializableMessage {
	var err *string
	if m.Err != nil {
		errStr := m.Err.Error()
		err = &errStr
	}

	return SerializableMessage{
		Message: m,
		Err:     err,
	}
}

func (m Message) Title() string {
	if m.Err != nil {
		return "error"
	}

	return fmt.Sprintf("%s: %s", utils.RightPadTrim("message ID", 12), m.ID)
}

func (m Message) Description() string {
	if m.Err != nil {
		return ""
	}

	if m.ContextKey != nil && m.ContextValue != nil {
		return fmt.Sprintf("%s: %s", utils.RightPadTrim(*m.ContextKey, 12), *m.ContextValue)
	}

	return ""
}

func (m Message) FilterValue() string {
	return m.ID
}

func GetMessageData(message *sqstypes.Message, config Config) Message {
	switch config.Format {
	case JSON:
		return getJSONMessage(message, config.SubsetKey, config.ContextKey)
	default:
		return getPlainMessage(message)
	}
}

func getJSONMessage(message *sqstypes.Message, subsetKey *string, contextKey *string) Message {
	if message.MessageId == nil {
		return Message{
			Err: errMessageIDNil,
		}
	}
	messageID := *message.MessageId
	if message.Body == nil {
		return Message{
			Err: errMessageBodyNil,
		}
	}
	messageBody := *message.Body
	bodyBytes := []byte(messageBody)

	var data map[string]any
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(fmt.Errorf("%w: %s", errCouldntUnmarshalBytes, err.Error()), messageID, bodyBytes),
		}
	}

	if subsetKey == nil && contextKey == nil {
		return getJSONMessageWithNoSubsetAndContext(messageID, bodyBytes, data)
	}

	if subsetKey == nil && contextKey != nil {
		return getJSONMessageWithContextButNoSubset(messageID, bodyBytes, data, *contextKey)
	}

	if subsetKey != nil && contextKey == nil {
		return getJSONMessageWithSubsetButNoContext(messageID, bodyBytes, *subsetKey)
	}

	return getJSONMessageWithSubsetAndContext(messageID, bodyBytes, *subsetKey, *contextKey)
}

func getJSONMessageWithNoSubsetAndContext(messageID string, bodyBytes []byte, data map[string]any) Message {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(fmt.Errorf("%w: %s", errCouldntMarshalBytes, err.Error()), messageID, bodyBytes),
		}
	}

	return Message{
		ID:   messageID,
		Body: string(bytes),
	}
}

func getJSONMessageWithContextButNoSubset(messageID string, bodyBytes []byte, data map[string]any, contextKey string) Message {
	contextValue, err := getContextValue(data, contextKey)
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(err, messageID, bodyBytes),
		}
	}

	return Message{
		ID:           messageID,
		Body:         string(bodyBytes),
		ContextKey:   &contextKey,
		ContextValue: contextValue,
	}
}

func getJSONMessageWithSubsetButNoContext(messageID string, bodyBytes []byte, subsetKey string) Message {
	subset, err := getSubset(bodyBytes, subsetKey)
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(err, messageID, bodyBytes),
		}
	}

	subsetBytes, err := json.MarshalIndent(subset, "", "  ")
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(fmt.Errorf("%w: %s", errCouldntMarshalBytes, err.Error()), messageID, bodyBytes),
		}
	}

	return Message{
		ID:   messageID,
		Body: string(subsetBytes),
	}
}

func getJSONMessageWithSubsetAndContext(messageID string, bodyBytes []byte, subsetKey, contextKey string) Message {
	subset, err := getSubset(bodyBytes, subsetKey)
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(err, messageID, bodyBytes),
		}
	}
	subsetBytes, err := json.MarshalIndent(subset, "", "  ")
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(fmt.Errorf("%w: %s", errCouldntMarshalBytes, err.Error()), messageID, bodyBytes),
		}
	}

	contextValue, err := getContextValue(subset, contextKey)
	if err != nil {
		return Message{
			Err: wrapErrWithDetails(err, messageID, bodyBytes),
		}
	}

	return Message{
		ID:           messageID,
		Body:         string(subsetBytes),
		ContextKey:   &contextKey,
		ContextValue: contextValue,
	}
}

func getSubset(bodyBytes []byte, key string) (map[string]any, error) {
	var data map[string]any
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return data, fmt.Errorf("%w: %s", errCouldntUnmarshalBytes, err.Error())
	}

	subset, ok := data[key]
	if !ok {
		return data, fmt.Errorf("%w (key: %q)", errSubsetKeyNotFound, key)
	}

	var subsetData map[string]any
	switch s := subset.(type) {
	case map[string]any:
		// It's a JSON object, directly access the nested key
		subsetData = s
	case string:
		// May be stringified JSON; attempt to convert it to JSON
		if err := json.Unmarshal([]byte(s), &subsetData); err != nil {
			return data, fmt.Errorf("%w (key: %q): %s", errCouldntUnmarshalSubsetValue, key, err.Error())
		}
	default:
		return data, fmt.Errorf("%w (key: %q); subset needs to be an object or stringified JSON; determined type: %T", errSubsetTypeIsUnsupported, key, s)
	}

	return subsetData, nil
}

func getContextValue(data map[string]any, key string) (*string, error) {
	var contextValue *string
	context, ok := data[key]
	if !ok {
		return contextValue, errContextKeyNotFound
	}

	switch c := context.(type) {
	case string:
		contextValue = &c
	case map[string]any:
		return contextValue, fmt.Errorf("%w (key: %q); determined type: object; context value needs to be a string", errContextValueTypeUnsupported, key)
	case []any:
		return contextValue, fmt.Errorf("%w (key: %q); determined type: array; context value needs to be a string", errContextValueTypeUnsupported, key)
	default:
		return contextValue, fmt.Errorf("%w (key: %q); determined type: %T; context value needs to be a string", errContextValueTypeUnsupported, key, c)
	}

	return contextValue, nil
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

func wrapErrWithDetails(err error, messageID string, bodyBytes []byte) error {
	return fmt.Errorf("%w\n\n- message id: %s\n- message body:\n>>>\n%s\n<<<", err, messageID, string(bodyBytes))
}
