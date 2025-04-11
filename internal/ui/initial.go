package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	t "github.com/dhth/cueitup/internal/types"
)

func InitialModel(
	sqsClient *sqs.Client,
	queueURL string,
	config t.Config,
	behaviours t.TUIBehaviours,
) Model {
	appDelegate := newAppItemDelegate()
	jobItems := make([]list.Item, 0)

	queueParts := strings.Split(queueURL, "/")
	queueName := queueParts[len(queueParts)-1]
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02-15-04-05")
	persistDir := fmt.Sprintf("messages/%s/%s", queueName, timeString)

	var dbg bool
	if len(os.Getenv("DEBUG")) > 0 {
		dbg = true
	}

	m := Model{
		sqsClient:           sqsClient,
		queueURL:            queueURL,
		config:              config,
		behaviours:          behaviours,
		msgsList:            list.New(jobItems, appDelegate, listWidth, 0),
		msgListCurrentIndex: -1,
		persistDir:          persistDir,
		showHelpIndicator:   true,
		debugMode:           dbg,
		firstFetch:          true,
	}
	m.msgsList.Title = "Messages"
	m.msgsList.SetStatusBarItemName("message", "messages")
	m.msgsList.SetFilteringEnabled(false)
	m.msgsList.DisableQuitKeybindings()
	m.msgsList.SetShowHelp(false)
	m.msgsList.Styles.Title = m.msgsList.Styles.Title.
		Background(lipgloss.Color(cueitupColor)).
		Foreground(lipgloss.Color(defaultForegroundColor)).
		Bold(true)

	return m
}
