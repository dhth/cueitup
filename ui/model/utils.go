package model

import (
	"encoding/json"
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

	var result string
	prettyJSON := pretty.Pretty([]byte(*message.Body))
	result = string(pretty.Color(prettyJSON, nil))

	return result, nil
}

func getRecordValueJSON(message *types.Message, extractKey string, keyProperty string) (string, string, error) {
	if message.Body == nil {
		return "", "", nil
	}

	var result string
	var data map[string]interface{}
	var keyValue string
	err := json.Unmarshal([]byte(*message.Body), &data)
	if err != nil {
		return "", "", err
	}
	if data[extractKey] != nil {
		prettyJSON := pretty.Pretty([]byte(data[extractKey].(string)))
		result = string(pretty.Color(prettyJSON, nil))

		if keyProperty != "" {
			var nested map[string]interface{}
			_ = json.Unmarshal([]byte(data[extractKey].(string)), &nested)
			if nested[keyProperty] != nil {
				keyValue = nested[keyProperty].(string)
			}
		}
	}

	return result, keyValue, nil
}

func getMessageData(message *types.Message, extractJSONObject string, keyProperty string) (string, string, error) {
	var msgValue string
	var err error
	var keyValue string
	if extractJSONObject != "" {
		msgValue, keyValue, err = getRecordValueJSON(message, extractJSONObject, keyProperty)
	} else {
		msgValue, err = getRecordValueJSONFull(message)
	}
	if err != nil {
		return "", "", err
	} else {
		return msgValue, keyValue, nil
	}
}
