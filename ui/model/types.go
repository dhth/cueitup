package model

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/charmbracelet/bubbles/key"
)

type delegateKeyMap struct {
	choose key.Binding
}

type KMsgItem struct {
	message         types.Message
	messageValue    string
	msgMetadata     string
	msgValue        string
	contextKeyName  string
	contextKeyValue string
}

func (item KMsgItem) Title() string {
	return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim("msgId", 10), *item.message.MessageId), listWidth)
}

func (item KMsgItem) Description() string {
	if item.contextKeyValue != "" {
		return RightPadTrim(fmt.Sprintf("%s: %s", RightPadTrim(item.contextKeyName, 10), item.contextKeyValue), listWidth)
	}
	return ""
}

func (item KMsgItem) FilterValue() string {
	return string(*item.message.MessageId)
}
