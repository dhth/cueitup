package utils

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	t "github.com/dhth/cueitup/internal/types"
	"github.com/tidwall/pretty"
)

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

	var data map[string]any
	err := json.Unmarshal([]byte(*message.Body), &data)
	if err != nil {
		return "", "", err
	}

	subsetKey, ok := data[extractKey]
	if !ok {
		return "", "", errors.New("nested object couln't be accessed")
	}

	var nestedData map[string]any
	switch n := subsetKey.(type) {
	case map[string]any:
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

func GetMessageData(message *types.Message, profile t.Profile) (string, string, error) {
	var msgValue, keyValue string
	var err error

	switch profile.Format {
	case t.JSON:
		if profile.SubsetKey != "" {
			msgValue,
				keyValue,
				err = getRecordValueJSONNested(message,
				profile.SubsetKey,
				profile.ContextKey,
			)
		} else {
			msgValue, err = getRecordValueJSONFull(message)
		}
	case t.None:
		msgValue = *message.Body
	}
	if err != nil {
		return "", "", err
	}
	return msgValue, keyValue, nil
}
