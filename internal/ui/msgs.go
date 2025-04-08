package ui

import (
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	t "github.com/dhth/cueitup/internal/types"
)

type (
	MsgCountTickMsg struct{}
	HideHelpMsg     struct{}
)

type SQSMsgsFetchedMsg struct {
	messages    []t.Message
	sqsMessages []sqsTypes.Message
	err         error
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
