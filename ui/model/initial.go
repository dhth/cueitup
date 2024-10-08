package model

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(sqsClient *sqs.Client, queueURL string, msgConsumptionConf MsgConsumptionConf) Model {
	appDelegate := newAppItemDelegate()
	jobItems := make([]list.Item, 0)

	queueParts := strings.Split(queueURL, "/")
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

	m := Model{
		sqsClient:            sqsClient,
		queueURL:             queueURL,
		msgConsumptionConf:   msgConsumptionConf,
		pollForQueueMsgCount: true,
		msgsList:             list.New(jobItems, appDelegate, listWidth+10, 0),
		recordValueStore:     make(map[string]string),
		persistDir:           persistDir,
		contextSearchInput:   ti,
		showHelpIndicator:    true,
		debugMode:            dbg,
		firstFetch:           true,
	}
	m.msgsList.Title = "Messages"
	m.msgsList.SetStatusBarItemName("message", "messages")
	m.msgsList.SetFilteringEnabled(false)
	m.msgsList.DisableQuitKeybindings()
	m.msgsList.SetShowHelp(false)
	m.msgsList.Styles.Title = m.msgsList.Styles.Title.
		Background(lipgloss.Color(listColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)

	return m
}
