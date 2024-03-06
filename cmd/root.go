package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/dhth/cueitup/ui"
	"github.com/dhth/cueitup/ui/model"
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var (
	queueUrl   = flag.String("queue-url", "", "url of the queue to consume from")
	awsProfile = flag.String("aws-profile", "", "aws profile to use")
	awsRegion  = flag.String("aws-region", "", "aws region to use")
	msgFormat  = flag.String("msg-format", "json", "message format")
	subsetKey  = flag.String("subset-key", "", "extract a nested object inside the JSON body")
	contextKey = flag.String("context-key", "", "the key to use as for context in the list")
)

func Execute() {
	flag.Parse()

	if *queueUrl == "" {
		die("queue-url cannot be empty")
	} else if !strings.HasPrefix(*queueUrl, "https://") {
		die("queue-url must begin with https")
	}

	if *awsProfile == "" {
		die("aws-profile cannot be empty")
	}

	if *awsRegion == "" {
		die("aws-region cannot be empty")
	}

	var msgFmt model.MsgFmt
	switch *msgFormat {
	case "json":
		msgFmt = model.JsonFmt
	case "plaintext":
		msgFmt = model.PlainTxtFmt
	default:
		die("cueitup only supports the following msg-format values: json, plaintext")
	}

	if *subsetKey != "" && msgFmt != model.JsonFmt {
		die("subset-key can only be used when msg-format=json")
	}
	if *contextKey != "" && msgFmt != model.JsonFmt {
		die("context-key can only be used when msg-format=json")
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
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	sqsClient := sqs.NewFromConfig(sdkConfig)

	ui.RenderUI(sqsClient, *queueUrl, msgConsumptionConf)

}
