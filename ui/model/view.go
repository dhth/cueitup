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
	var headerViewPtr string
	var valueViewPtr string
	var mode string
	var statusBar string
	var debugMsg string
	var msgValVPTitleStyle lipgloss.Style

	if m.msg != "" {
		statusBar = Trim(m.msg, 120)
	}
	m.msgsList.Styles.Title.Background(lipgloss.Color(inactivePaneColor))
	msgValVPTitleStyle = msgValueTitleStyle.Copy()

	switch m.activeView {
	case kMsgsListView:
		m.msgsList.Styles.Title.Background(lipgloss.Color(activeHeaderColor))
	case kMsgValueView:
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

	var msgMetadataVP string
	if !m.msgValueVPReady {
		msgMetadataVP = "\n  Initializing..."
	} else {
		msgMetadataVP = msgValueVPStyle.Render(fmt.Sprintf("%s%s\n\n%s\n", msgDetailsTitleStyle.Render("Message Metadata"), headerViewPtr, m.msgMetadataVP.View()))
	}

	var msgValueVP string
	if !m.msgValueVPReady {
		msgValueVP = "\n  Initializing..."
	} else {
		msgValueVP = msgValueVPStyle.Render(fmt.Sprintf("%s%s\n\n%s\n", msgValVPTitleStyle.Render("Message Value"), valueViewPtr, m.msgValueVP.View()))
	}
	var helpVP string
	if !m.helpVPReady {
		helpVP = "\n  Initializing..."
	} else {
		helpVP = helpVPStyle.Render(fmt.Sprintf("  %s\n\n%s\n", helpVPTitleStyle.Render("Help"), m.helpVP.View()))
	}

	switch m.vpFullScreen {
	case false:
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			msgListStyle.Render(m.msgsList.View()),
			msgValueVP,
		)
	case true:
		switch m.activeView {
		case kMsgMetadataView:
			content = msgMetadataVP
		case kMsgValueView:
			content = msgValueVP
		case helpView:
			content = helpVP
		}
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
