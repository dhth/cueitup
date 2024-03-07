package model

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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

	ti := textinput.New()
	ti.Prompt = fmt.Sprintf("Filter messages where %s in > ", msgConsumptionConf.ContextKey)
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 100

	var dbg bool
	if len(os.Getenv("DEBUG")) > 0 {
		dbg = true
	}

	m := model{
		sqsClient:            sqsClient,
		queueUrl:             queueUrl,
		msgConsumptionConf:   msgConsumptionConf,
		pollForQueueMsgCount: true,
		kMsgsList:            list.New(jobItems, appDelegate, listWidth+10, 0),
		recordMetadataStore:  make(map[string]string),
		recordValueStore:     make(map[string]string),
		persistDir:           persistDir,
		contextSearchInput:   ti,
		showHelpIndicator:    true,
		debugMode:            dbg,
	}
	m.kMsgsList.Title = "Messages"
	m.kMsgsList.SetStatusBarItemName("message", "messages")
	m.kMsgsList.SetFilteringEnabled(false)
	m.kMsgsList.SetShowHelp(false)

	return m
}
