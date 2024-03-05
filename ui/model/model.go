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

type DeserializationFmt uint

const (
	JsonFmt DeserializationFmt = iota
	ProtobufFmt
)

const msgCountTickInterval = time.Second * 20

type model struct {
	deserializationFmt   DeserializationFmt
	sqsClient            *sqs.Client
	queueUrl             string
	extractJSONObject    string
	keyProperty          string
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
