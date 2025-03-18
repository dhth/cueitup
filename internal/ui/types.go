package ui

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/dhth/cueitup/internal/utils"
)

type msgItem struct {
	message         types.Message
	messageValue    string
	contextKeyName  string
	contextKeyValue string
}

func (item msgItem) Title() string {
	return utils.RightPadTrim(fmt.Sprintf("%s: %s", utils.RightPadTrim("msgId", 10), *item.message.MessageId), listWidth)
}

func (item msgItem) Description() string {
	if item.contextKeyValue != "" {
		return utils.RightPadTrim(fmt.Sprintf("%s: %s", utils.RightPadTrim(item.contextKeyName, 10), item.contextKeyValue), listWidth)
	}
	return ""
}

func (item msgItem) FilterValue() string {
	return *item.message.MessageId
}
