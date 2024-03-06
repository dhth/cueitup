package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
)

func InitialModel(sqsClient *sqs.Client, queueUrl string, msgConsumptionConf MsgConsumptionConf) model {

	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)
	jobItems := make([]list.Item, 0)

	queueParts := strings.Split(queueUrl, "/")
	queueName := queueParts[len(queueParts)-1]
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02-15-04-05")
	persistDir := fmt.Sprintf("messages/%s/%s", queueName, timeString)

	m := model{
		sqsClient:            sqsClient,
		queueUrl:             queueUrl,
		msgConsumptionConf:   msgConsumptionConf,
		pollForQueueMsgCount: true,
		kMsgsList:            list.New(jobItems, appDelegate, 60, 0),
		recordMetadataStore:  make(map[string]string),
		recordValueStore:     make(map[string]string),
		persistDir:           persistDir,
	}
	m.kMsgsList.Title = "Messages"
	m.kMsgsList.SetFilteringEnabled(false)
	m.kMsgsList.SetShowHelp(false)

	return m
}
