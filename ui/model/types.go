package model

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/charmbracelet/bubbles/key"
)

type delegateKeyMap struct {
	choose key.Binding
}

type KMsgItem struct {
	message      types.Message
	messageValue string
	msgMetadata  string
	msgValue     string
}

func (item KMsgItem) Title() string {
	return string(*item.message.MessageId)
}

func (item KMsgItem) Description() string {
	return ""
}

func (item KMsgItem) FilterValue() string {
	return string(*item.message.MessageId)
}
