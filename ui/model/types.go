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
	message          types.Message
	messageValue     string
	msgMetadata      string
	msgValue         string
	keyPropertyName  string
	keyPropertyValue string
}

func (item KMsgItem) Title() string {
	return fmt.Sprintf("%s: %s", RightPadTrim("msgId", 10), RightPadTrim(*item.message.MessageId, 40))
}

func (item KMsgItem) Description() string {
	if item.keyPropertyValue != "" {
		return fmt.Sprintf("%s: %s", RightPadTrim(item.keyPropertyName, 10), RightPadTrim(item.keyPropertyValue, 40))
	}
	return ""
}

func (item KMsgItem) FilterValue() string {
	return string(*item.message.MessageId)
}
