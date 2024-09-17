package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/dhth/cueitup/ui"
	"github.com/dhth/cueitup/ui/model"
)

var (
	errQueueURLEmpty        = errors.New("queue URL is empty")
	errAWSProfileEmpty      = errors.New("AWS profile is empty")
	errAWSRegionEmpty       = errors.New("AWS region is empty")
	errQueueURLIncorrect    = errors.New("queue URL is incorrect")
	errMsgFormatInvalid     = errors.New("message format is invalid")
	errInvalidFlagUsage     = errors.New("invalid flag usage")
	errCouldntLoadAWSConfig = errors.New("couldn't load AWS config")
)

var (
	queueURL   = flag.String("queue-url", "", "url of the queue to consume from")
	awsProfile = flag.String("aws-profile", "", "aws profile to use")
	awsRegion  = flag.String("aws-region", "", "aws region to use")
	msgFormat  = flag.String("msg-format", "json", "message format")
	subsetKey  = flag.String("subset-key", "", "extract a nested object inside the JSON body")
	contextKey = flag.String("context-key", "", "the key to use as for context in the list")
)

func Execute() error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Inspect messages in an AWS SQS queue in a simple and deliberate manner.\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n------\n%s", model.HelpText)
	}
	flag.Parse()

	if *queueURL == "" {
		return errQueueURLEmpty
	}

	if !strings.HasPrefix(*queueURL, "https://") {
		return fmt.Errorf("%w: must begin with https", errQueueURLIncorrect)
	}

	if *awsProfile == "" {
		return errAWSProfileEmpty
	}

	if *awsRegion == "" {
		return errAWSRegionEmpty
	}

	var msgFmt model.MsgFmt
	switch *msgFormat {
	case "json":
		msgFmt = model.JSONFmt
	case "plaintext":
		msgFmt = model.PlainTxtFmt
	default:
		return fmt.Errorf("%w: supported values: json, plaintext", errMsgFormatInvalid)
	}

	if *subsetKey != "" && msgFmt != model.JSONFmt {
		return fmt.Errorf("%w: subset-key can only be used when msg-format=json", errInvalidFlagUsage)
	}
	if *contextKey != "" && msgFmt != model.JSONFmt {
		return fmt.Errorf("%w: context-key can only be used when msg-format=json", errInvalidFlagUsage)
	}

	msgConsumptionConf := model.MsgConsumptionConf{
		Format:     msgFmt,
		SubsetKey:  *subsetKey,
		ContextKey: *contextKey,
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(*awsProfile),
		config.WithRegion(*awsRegion),
	)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntLoadAWSConfig, err.Error())
	}

	sqsClient := sqs.NewFromConfig(sdkConfig)

	return ui.RenderUI(sqsClient, *queueURL, msgConsumptionConf)
}
