package model

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MsgCountTickMsg struct{}

type KMsgFetchedMsg struct {
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

type KMsgChosenMsg struct {
	key string
}

type RecordSavedToDiskMsg struct {
	path string
	err  error
}

type KMsgValueReadyMsg struct {
	storeKey string
	record   *types.Message
	msgValue string
	err      error
}
