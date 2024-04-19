package model

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	listWidth = 50
)

func (m model) View() string {
	var content string
	var footer string
	var mode string
	var statusBar string
	var debugMsg string
	var msgValVPTitleStyle lipgloss.Style

	if m.message != "" {
		statusBar = Trim(m.message, 120)
	}
	m.msgsList.Styles.Title.Background(lipgloss.Color(inactivePaneColor))
	msgValVPTitleStyle = msgValueTitleStyle.Copy()

	switch m.activeView {
	case msgsListView:
		m.msgsList.Styles.Title.Background(lipgloss.Color(activeHeaderColor))
	case msgValueView:
		msgValVPTitleStyle.Background(lipgloss.Color(activeHeaderColor))
	case contextualSearchView:
		statusBar = m.contextSearchInput.View()
	}

	if !m.deleteMsgs {
		mode += " " + deletingMsgsStyle.Render("not deleting msgs!")
	}

	if !m.pollForQueueMsgCount {
		mode += " " + pollingMsgStyle.Render("not polling for msg count!")
	}

	if m.persistRecords {
		mode += " " + persistingStyle.Render("persisting msgs!")
	}

	if m.skipRecords {
		mode += " " + skippingStyle.Render("skipping msgs!")
	}

	if m.filterMessages && len(m.contextSearchValues) > 0 {
		mode += " " + skippingStyle.Render(fmt.Sprintf("filtering where %s in : %v", m.msgConsumptionConf.ContextKey, m.contextSearchValues))
	}

	var errorMsg string
	if m.errorMsg != "" {
		errorMsg = " error: " + Trim(m.errorMsg, 120)
	}

	var msgValueVP string
	if !m.msgValueVPReady {
		msgValueVP = "\n  Initializing..."
	} else {
		switch m.vpFullScreen {
		case true:
			msgValueVP = msgValueVPFSStyle.Render(fmt.Sprintf("  %s\n\n%s\n", msgValVPTitleStyle.Render("Message Value"), m.msgValueVP.View()))
		case false:
			msgValueVP = msgValueVPStyle.Render(fmt.Sprintf("  %s\n\n%s\n", msgValVPTitleStyle.Render("Message Value"), m.msgValueVP.View()))
		}
	}
	var helpVP string
	if !m.helpVPReady {
		helpVP = "\n  Initializing..."
	} else {
		helpVP = helpVPStyle.Render(fmt.Sprintf("  %s\n\n%s\n", helpVPTitleStyle.Render("Help"), m.helpVP.View()))
	}

	switch m.activeView {
	case msgsListView, contextualSearchView:
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			msgListStyle.Render(m.msgsList.View()),
			msgValueVP,
		)
	case msgValueView:
		switch m.vpFullScreen {
		case true:
			content = msgValueVP
		case false:
			content = lipgloss.JoinHorizontal(
				lipgloss.Top,
				msgListStyle.Render(m.msgsList.View()),
				msgValueVP,
			)
		}
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
