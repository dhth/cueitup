package model

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type KMsgFetchedMsg struct {
	messages []types.Message
	err      error
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
