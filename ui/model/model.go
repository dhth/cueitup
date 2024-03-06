package model

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type stateView uint

const (
	kMsgsListView stateView = iota
	kMsgMetadataView
	kMsgValueView
	helpView
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
	kMsgsList            list.Model
	helpVP               viewport.Model
	helpSeen             uint
	msgMetadataVP        viewport.Model
	msgValueVP           viewport.Model
	recordMetadataStore  map[string]string
	recordValueStore     map[string]string
	deleteMsgs           bool
	skipRecords          bool
	persistRecords       bool
	persistDir           string
	filteredKeys         []string
	msgMetadataVPReady   bool
	msgValueVPReady      bool
	helpVPReady          bool
	vpFullScreen         bool
	terminalWidth        int
	terminalHeight       int
	msg                  string
	errorMsg             string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		GetQueueMsgCount(m.sqsClient, m.queueUrl),
		tickEvery(msgCountTickInterval),
	)
}
