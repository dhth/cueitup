package model

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type stateView uint

const (
	msgsListView stateView = iota
	msgValueView
	helpView
	contextualSearchView
)

type MsgFmt uint

const (
	JsonFmt MsgFmt = iota
	PlainTxtFmt
)

const msgCountTickInterval = time.Second * 3

type MsgConsumptionConf struct {
	Format     MsgFmt
	SubsetKey  string
	ContextKey string
}

type model struct {
	deserializationFmt   MsgFmt
	sqsClient            *sqs.Client
	queueUrl             string
	msgConsumptionConf   MsgConsumptionConf
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
	deleteMsgs           bool
	skipRecords          bool
	persistRecords       bool
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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		GetQueueMsgCount(m.sqsClient, m.queueUrl),
		tickEvery(msgCountTickInterval),
		hideHelp(time.Minute*1),
	)
}
