package model

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
)

func InitialModel(sqsClient *sqs.Client, queueUrl string, extractJSONObject string, keyProperty string) model {

	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)
	jobItems := make([]list.Item, 0)

	m := model{
		sqsClient:            sqsClient,
		queueUrl:             queueUrl,
		extractJSONObject:    extractJSONObject,
		keyProperty:          keyProperty,
		pollForQueueMsgCount: true,
		kMsgsList:            list.New(jobItems, appDelegate, 60, 0),
		deleteMsgs:           true,
		persistRecords:       false,
		recordMetadataStore:  make(map[string]string),
		recordValueStore:     make(map[string]string),
	}
	m.kMsgsList.Title = "Messages"
	m.kMsgsList.SetFilteringEnabled(false)
	m.kMsgsList.SetShowHelp(false)

	return m
}
