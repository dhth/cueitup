package ui

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	t "github.com/dhth/cueitup/internal/types"
)

type stateView uint

const (
	msgsListView stateView = iota
	msgValueView
	helpView
	contextualSearchView
)

const msgCountTickInterval = time.Second * 3

type Model struct {
	sqsClient            *sqs.Client
	queueURL             string
	config               t.Config
	behaviours           t.Behaviours
	activeView           stateView
	lastView             stateView
	pollForQueueMsgCount bool
	msgsList             list.Model
	helpVP               viewport.Model
	showHelpIndicator    bool
	msgValueVP           viewport.Model
	recordValueStore     map[string]string
	contextSearchInput   textinput.Model
	contextSearchValues  []string
	filterMessages       bool
	persistDir           string
	msgValueVPReady      bool
	helpVPReady          bool
	vpFullScreen         bool
	terminalWidth        int
	terminalHeight       int
	message              string
	errorMsg             string
	debugMode            bool
	firstFetch           bool
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		GetQueueMsgCount(m.sqsClient, m.queueURL),
		tickEvery(msgCountTickInterval),
		hideHelp(time.Minute*1),
	)
}
