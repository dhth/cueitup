package ui

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/cueitup/ui/model"
)

var errFailedToConfigureDebugging = errors.New("failed to configure debugging")

func RenderUI(sqsClient *sqs.Client, queueURL string, msgConsumptionConf model.MsgConsumptionConf) error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("%w: %s", errFailedToConfigureDebugging, err.Error())
		}
		defer f.Close()
	}
	p := tea.NewProgram(model.InitialModel(sqsClient, queueURL, msgConsumptionConf), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
