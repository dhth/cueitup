package model

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tidwall/pretty"
)

func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func getRecordValueJSONFull(message *types.Message) (string, error) {
	if message.Body == nil {
		return "", nil
	}

	prettyJSON := pretty.Pretty([]byte(*message.Body))

	return string(prettyJSON), nil
}

func getRecordValueJSONNested(message *types.Message, extractKey string, contextKey string) (string, string, error) {
	if message.Body == nil {
		return "", "", errors.New("body is nil")
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(*message.Body), &data)
	if err != nil {
		return "", "", err
	}

	subsetKey, ok := data[extractKey]
	if !ok {
		return "", "", errors.New("nested object couln't be accessed")
	}

	var nestedData map[string]interface{}
	switch n := subsetKey.(type) {
	case map[string]interface{}:
		// If it's a JSON object, directly access the nested key
		nestedData = n
	case string:
		// May be stringified JSON; attempt to convert it to JSON
		if err := json.Unmarshal([]byte(n), &nestedData); err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.New("Unexpected type")
	}

	nestedBytes, err := json.MarshalIndent(nestedData, "", "    ")
	if err != nil {
		return "", "", err
	}

	contextualValue, ok := nestedData[contextKey]
	if !ok {
		return string(nestedBytes), "", nil
	}

	return string(nestedBytes), contextualValue.(string), nil
}

func getMessageData(message *types.Message, msgConsumptionConf MsgConsumptionConf) (string, string, error) {
	var msgValue, keyValue string
	var err error

	switch msgConsumptionConf.Format {
	case JsonFmt:
		if msgConsumptionConf.SubsetKey != "" {
			msgValue,
				keyValue,
				err = getRecordValueJSONNested(message,
				msgConsumptionConf.SubsetKey,
				msgConsumptionConf.ContextKey,
			)
		} else {
			msgValue, err = getRecordValueJSONFull(message)
		}
	case PlainTxtFmt:
		msgValue = *message.Body
	}
	if err != nil {
		return "", "", err
	} else {
		return msgValue, keyValue, nil
	}
}
