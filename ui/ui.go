package ui

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/cueitup/ui/model"
)

func RenderUI(sqsClient *sqs.Client, queueUrl string, msgConsumptionConf model.MsgConsumptionConf) {
	p := tea.NewProgram(model.InitialModel(sqsClient, queueUrl, msgConsumptionConf), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Something went wrong %s", err)
	}
}
