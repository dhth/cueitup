package model

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type msgItem struct {
	message         types.Message
	messageValue    string
	msgMetadata     string
	msgValue        string
	contextKeyName  string
	contextKeyValue string
}

func (item msgItem) Title() string {
	return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim("msgId", 10), *item.message.MessageId), listWidth)
}

func (item msgItem) Description() string {
	if item.contextKeyValue != "" {
		return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim(item.contextKeyName, 10), item.contextKeyValue), listWidth)
	}
	return ""
}

func (item msgItem) FilterValue() string {
	return string(*item.message.MessageId)
}
