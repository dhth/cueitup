package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/cueitup/internal/utils"
)

var listWidth = 54

func (m Model) View() string {
	var content string
	var footer string
	var mode string
	var statusBar string
	var debugMsg string

	if m.message != "" {
		statusBar = m.message
	}

	m.msgsList.Styles.Title = m.msgsList.Styles.Title.Background(lipgloss.Color(inactivePaneColor))
	msgValTitleStyleToUse := msgValueTitleStyle

	switch m.activeView {
	case msgsListView:
		m.msgsList.Styles.Title = m.msgsList.Styles.Title.Background(lipgloss.Color(cueitupColor))
	case msgValueView:
		msgValTitleStyleToUse = msgValTitleStyleToUse.Background(lipgloss.Color(cueitupColor))
	}

	if !m.behaviours.DeleteMessages {
		mode += " " + deletingMsgsStyle.Render("not deleting msgs!")
	} else {
		mode += " " + deletingMsgsStyle.Render("deleting msgs!")
	}

	if !m.pollForQueueMsgCount {
		mode += " " + pollingMsgStyle.Render("not polling for msg count!")
	}

	if m.behaviours.PersistMessages {
		mode += " " + persistingStyle.Render("persisting msgs!")
	}

	if m.behaviours.SkipMessages {
		mode += " " + skippingStyle.Render("skipping msgs!")
	}

	var errorMsg string
	if m.errorMsg != "" {
		errorMsg = " error: " + utils.Trim(m.errorMsg, 120)
	}

	var msgValueVP string
	if !m.msgValueVPReady {
		msgValueVP = "\n  Initializing..."
	} else {
		msgValueVP = msgValueVPStyle.Render(fmt.Sprintf("%s\n\n%s\n", msgValTitleStyleToUse.Render("Message Value"), m.msgValueVP.View()))
	}
	var helpVP string
	if !m.helpVPReady {
		helpVP = "\n  Initializing..."
	} else {
		helpVP = helpVPStyle.Render(fmt.Sprintf("  %s\n\n%s\n", helpVPTitleStyle.Render("Help"), m.helpVP.View()))
	}

	switch m.activeView {
	case msgsListView:
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			msgListStyle.Render(m.msgsList.View()),
			msgValueVP,
		)
	case msgValueView:
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			msgListStyle.Render(m.msgsList.View()),
			msgValueVP,
		)
	case helpView:
		content = helpVP
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282828")).
		Background(lipgloss.Color("#7c6f64"))

	var helpMsg string
	if m.showHelpIndicator {
		helpMsg = " " + helpMsgStyle.Render("Press ? for help")
	}

	if m.debugMode {
		debugMsg += fmt.Sprintf(" %v", m.activeView)
	}

	footerStr := fmt.Sprintf("%s%s%s%s%s",
		modeStyle.Render("cueitup"),
		debugMsg,
		helpMsg,
		mode,
		errorMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
