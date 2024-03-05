package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	tea "github.com/charmbracelet/bubbletea"
)

func FetchMessages(client *sqs.Client, queueUrl string, maxMessages int32, waitTime int32) tea.Cmd {
	return func() tea.Msg {

		var messages []types.Message
		result, err := client.ReceiveMessage(context.TODO(),
			&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueUrl),
				MaxNumberOfMessages: maxMessages,
				WaitTimeSeconds:     waitTime,
				VisibilityTimeout:   30,
			})
		if err != nil {
			return KMsgFetchedMsg{
				messages: nil,
				err:      err,
			}
		} else {
			messages = result.Messages
		}

		return KMsgFetchedMsg{
			messages: messages,
			err:      nil,
		}
	}
}

func saveRecordValue(message *types.Message, extractJSONObject string) tea.Cmd {
	return func() tea.Msg {
		var msgValue string
		var err error
		if extractJSONObject != "" {
			msgValue, err = getRecordValueJSON(message, extractJSONObject)
		} else {
			msgValue, err = getRecordValueJSONFull(message)
		}
		if err != nil {
			return KMsgValueReadyMsg{err: err}
		} else {
			return KMsgValueReadyMsg{storeKey: *message.MessageId, record: message, msgValue: msgValue}
		}
	}
}

func showItemDetails(key string) tea.Cmd {
	return func() tea.Msg {
		return KMsgChosenMsg{key}
	}
}
