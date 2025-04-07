package ui

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type (
	MsgCountTickMsg struct{}
	HideHelpMsg     struct{}
)

type SQSMsgFetchedMsg struct {
	messages      []types.Message
	messageValues []string
	keyValues     []string
	err           error
}

type QueueMsgCountFetchedMsg struct {
	approxMsgCount int
	err            error
}

type SQSMsgsDeletedMsg struct {
	err error
}

type MsgChosenMsg struct {
	key string
}

type ContextSearchValuesSetMsg struct {
	values []string
}

type RecordSavedToDiskMsg struct {
	path string
	err  error
}

type KMsgValueReadyMsg struct {
	storeKey string
	msgValue string
	err      error
}
