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
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

var (
	queueUrl          = flag.String("queue-url", "", "url of the queue to consume from")
	awsProfile        = flag.String("aws-profile", "", "aws profile to use")
	awsRegion         = flag.String("aws-region", "", "aws region to use")
	extractJSONObject = flag.String("json-extract", "", "extract a nested object inside the JSON body")
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

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(*awsProfile),
		config.WithRegion(*awsRegion),
	)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	sqsClient := sqs.NewFromConfig(sdkConfig)

	ui.RenderUI(sqsClient, *queueUrl, *extractJSONObject)

}
