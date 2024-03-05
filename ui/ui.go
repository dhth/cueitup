package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/cueitup/ui/model"
)

func RenderUI(sqsClient *sqs.Client, queueUrl string, extractJSONObject string, keyProperty string) {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(model.InitialModel(sqsClient, queueUrl, extractJSONObject, keyProperty), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Something went wrong %s", err)
	}
}
