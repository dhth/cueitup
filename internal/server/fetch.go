package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	t "github.com/dhth/cueitup/internal/types"
)

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json; charset=utf-8"
	unexpected      = "something unexpected happened (let @dhth know about this via https://github.com/dhth/cueitup/issues)"
)

type MessageCount struct {
	Count int `json:"count"`
}

func getMessages(client *sqs.Client, config t.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		numMessagesStr := queryParams.Get("num")

		numMessages := 1
		if numMessagesStr != "" {
			num, err := strconv.Atoi(numMessagesStr)
			if err != nil || num < 1 {
				http.Error(w, fmt.Sprintf("incorrect value provided for query param \"num\": %s", err.Error()), http.StatusBadRequest)
				return
			}
			numMessages = num
		}
		if numMessages > 10 {
			numMessages = 10
		}

		deleteStr := queryParams.Get("delete")
		var deleteMessages bool
		if deleteStr != "" {
			parsed, err := strconv.ParseBool(deleteStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("incorrect value provided for query param \"delete\": %s", err.Error()), http.StatusBadRequest)
				return
			}
			deleteMessages = parsed
		}

		result, err := client.ReceiveMessage(context.TODO(),
			&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(config.QueueURL),
				MaxNumberOfMessages: int32(numMessages),
				WaitTimeSeconds:     0,
				VisibilityTimeout:   30,
			})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch messages: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		messages := make([]t.SerializableMessage, len(result.Messages))
		for i, message := range result.Messages {
			messages[i] = t.GetMessageData(&message, config).ToSerializable()
		}

		if deleteMessages && len(messages) > 0 {
			deleteEntries := make([]sqstypes.DeleteMessageBatchRequestEntry, len(messages))
			for i := range result.Messages {
				deleteEntries[i].Id = aws.String(fmt.Sprintf("%v", i))
				deleteEntries[i].ReceiptHandle = result.Messages[i].ReceiptHandle
			}
			_, err = client.DeleteMessageBatch(context.TODO(),
				&sqs.DeleteMessageBatchInput{
					Entries:  deleteEntries,
					QueueUrl: aws.String(config.QueueURL),
				})
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to delete messages on SQS: %s", err.Error()), http.StatusInternalServerError)
				return
			}
		}

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		if _, err := w.Write(jsonBytes); err != nil {
			log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
		}
	}
}

func getMessageCount(client *sqs.Client, config t.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		approxMsgCountType := sqstypes.QueueAttributeNameApproximateNumberOfMessages
		attribute, err := client.GetQueueAttributes(context.TODO(),
			&sqs.GetQueueAttributesInput{
				QueueUrl:       aws.String(config.QueueURL),
				AttributeNames: []sqstypes.QueueAttributeName{approxMsgCountType},
			})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get message message count: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		countStr := attribute.Attributes[string(approxMsgCountType)]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to convert message count to an int: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(MessageCount{count})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		if _, err := w.Write(jsonBytes); err != nil {
			log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
		}
	}
}

func getConfig(config t.Config) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		jsonBytes, err := json.Marshal(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		if _, err := w.Write(jsonBytes); err != nil {
			log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
		}
	}
}

func getBehaviours(behaviours t.WebBehaviours) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		jsonBytes, err := json.Marshal(behaviours)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, applicationJSON)
		if _, err := w.Write(jsonBytes); err != nil {
			log.Printf("failed to write bytes to HTTP connection: %s", err.Error())
		}
	}
}
