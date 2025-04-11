package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	tea "github.com/charmbracelet/bubbletea"
	t "github.com/dhth/cueitup/internal/types"
)

func (m Model) FetchMessages(maxMessages int32, waitTime int32) tea.Cmd {
	return func() tea.Msg {
		result, err := m.sqsClient.ReceiveMessage(context.TODO(),
			// WaitTimeSeconds > 0 enables long polling
			// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-short-and-long-polling.html#sqs-long-polling
			&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(m.queueURL),
				MaxNumberOfMessages: maxMessages,
				WaitTimeSeconds:     waitTime,
				VisibilityTimeout:   30,
			})
		if err != nil {
			return SQSMsgsFetchedMsg{
				err: err,
			}
		}
		messages := make([]t.Message, len(result.Messages))
		for i, message := range result.Messages {
			messages[i] = t.GetMessageData(&message, m.config)
		}

		return SQSMsgsFetchedMsg{
			messages:    messages,
			sqsMessages: result.Messages,
		}
	}
}

func DeleteMessages(client *sqs.Client, queueURL string, messages []sqstypes.Message) tea.Cmd {
	return func() tea.Msg {
		entries := make([]sqstypes.DeleteMessageBatchRequestEntry, len(messages))
		for i := range messages {
			entries[i].Id = aws.String(fmt.Sprintf("%v", i))
			entries[i].ReceiptHandle = messages[i].ReceiptHandle
		}
		_, err := client.DeleteMessageBatch(context.TODO(),
			&sqs.DeleteMessageBatchInput{
				Entries:  entries,
				QueueUrl: aws.String(queueURL),
			})
		if err != nil {
			return SQSMsgsDeletedMsg{
				err: err,
			}
		}

		return SQSMsgsDeletedMsg{}
	}
}

func GetQueueMsgCount(client *sqs.Client, queueURL string) tea.Cmd {
	return func() tea.Msg {
		approxMsgCountType := sqstypes.QueueAttributeNameApproximateNumberOfMessages
		attribute, err := client.GetQueueAttributes(context.TODO(),
			&sqs.GetQueueAttributesInput{
				QueueUrl:       aws.String(queueURL),
				AttributeNames: []sqstypes.QueueAttributeName{approxMsgCountType},
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

func saveMessageToDisk(id string, value string, format t.MessageFormat, dir string) tea.Cmd {
	return func() tea.Msg {
		now := time.Now().Unix()
		fileName := fmt.Sprintf("%d-%s.%s", now, id, format.Extension())
		fp := filepath.Join(dir, fileName)
		dir := filepath.Dir(fp)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return RecordSavedToDiskMsg{err: err}
		}

		err = os.WriteFile(fp, []byte(value), 0o644)
		if err != nil {
			return RecordSavedToDiskMsg{err: err}
		}

		return RecordSavedToDiskMsg{path: fp}
	}
}

func tickEvery(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return MsgCountTickMsg{}
	})
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return HideHelpMsg{}
	})
}
