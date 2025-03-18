package ui

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	tea "github.com/charmbracelet/bubbletea"
	t "github.com/dhth/cueitup/internal/types"
)

var errFailedToConfigureDebugging = errors.New("failed to configure debugging")

func RenderUI(
	sqsClient *sqs.Client,
	queueURL string,
	profile t.Profile,
	behaviours t.Behaviours,
) error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("%w: %s", errFailedToConfigureDebugging, err.Error())
		}
		defer f.Close()
	}
	p := tea.NewProgram(InitialModel(sqsClient, queueURL, profile, behaviours), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
