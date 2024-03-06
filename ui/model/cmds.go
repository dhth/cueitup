package model

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) FetchMessages(maxMessages int32, waitTime int32) tea.Cmd {
	return func() tea.Msg {

		var messages []types.Message
		var messagesValues []string
		var keyValues []string
		result, err := m.sqsClient.ReceiveMessage(context.TODO(),
			// WaitTimeSeconds > 0 enables long polling
			// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-short-and-long-polling.html#sqs-long-polling
			&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(m.queueUrl),
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
			for _, message := range messages {
				msgValue, keyValue, _ := getMessageData(&message, m.msgConsumptionConf)
				messagesValues = append(messagesValues, msgValue)
				keyValues = append(keyValues, keyValue)
			}
		}

		return KMsgFetchedMsg{
			messages:      messages,
			messageValues: messagesValues,
			keyValues:     keyValues,
			err:           nil,
		}
	}
}

func DeleteMessages(client *sqs.Client, queueUrl string, messages []types.Message) tea.Cmd {
	return func() tea.Msg {

		entries := make([]types.DeleteMessageBatchRequestEntry, len(messages))
		for msgIndex := range messages {
			entries[msgIndex].Id = aws.String(fmt.Sprintf("%v", msgIndex))
			entries[msgIndex].ReceiptHandle = messages[msgIndex].ReceiptHandle
		}
		_, err := client.DeleteMessageBatch(context.TODO(),
			&sqs.DeleteMessageBatchInput{
				Entries:  entries,
				QueueUrl: aws.String(queueUrl),
			})
		if err != nil {
			return SQSMsgsDeletedMsg{
				err: err,
			}
		}

		return SQSMsgsDeletedMsg{}
	}
}

func GetQueueMsgCount(client *sqs.Client, queueUrl string) tea.Cmd {
	return func() tea.Msg {

		approxMsgCountType := types.QueueAttributeNameApproximateNumberOfMessages
		attribute, err := client.GetQueueAttributes(context.TODO(),
			&sqs.GetQueueAttributesInput{
				QueueUrl:       aws.String(queueUrl),
				AttributeNames: []types.QueueAttributeName{approxMsgCountType},
			})

		if err != nil {
			return QueueMsgCountFetchedMsg{
				approxMsgCount: -1,
				err:            err,
			}
		}

		countStr := attribute.Attributes[string(approxMsgCountType)]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return QueueMsgCountFetchedMsg{
				approxMsgCount: -1,
				err:            err,
			}
		}
		return QueueMsgCountFetchedMsg{
			approxMsgCount: count,
		}
	}
}

func saveRecordValueToDisk(filePath string, msgValue string, msgFmt MsgFmt) tea.Cmd {
	return func() tea.Msg {
		dir := filepath.Dir(filePath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return RecordSavedToDiskMsg{err: err}
		}
		var data string
		switch msgFmt {
		case JsonFmt:
			data = fmt.Sprintf("```json\n%s\n```", msgValue)
		case PlainTxtFmt:
			data = msgValue
		}
		err = os.WriteFile(filePath, []byte(data), 0644)
		if err != nil {
			return RecordSavedToDiskMsg{err: err}
		}
		return RecordSavedToDiskMsg{path: filePath}
	}
}

func setContextSearchValues(userInput string) tea.Cmd {
	return func() tea.Msg {
		valuesEls := strings.Split(userInput, ",")
		var values []string
		for _, v := range valuesEls {
			values = append(values, strings.TrimSpace(v))
		}
		return ContextSearchValuesSetMsg{values}
	}
}

func showItemDetails(key string) tea.Cmd {
	return func() tea.Msg {
		return KMsgChosenMsg{key}
	}
}

func tickEvery(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return MsgCountTickMsg{}
	})
}
