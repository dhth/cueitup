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
	var msgsViewPtr string
	var headerViewPtr string
	var valueViewPtr string
	var mode string
	var statusBar string
	var debugMsg string

	if m.msg != "" {
		statusBar = Trim(m.msg, 120)
	}

	switch m.activeView {
	case kMsgsListView:
		msgsViewPtr = " 👇"
	case kMsgMetadataView:
		headerViewPtr = " 👇"
	case kMsgValueView:
		valueViewPtr = " 👇"
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

	m.kMsgsList.Title += msgsViewPtr

	var errorMsg string
	if m.errorMsg != "" {
		errorMsg = " error: " + Trim(m.errorMsg, 120)
	}

	var msgMetadataVP string
	if !m.msgValueVPReady {
		msgMetadataVP = "\n  Initializing..."
	} else {
		msgMetadataVP = viewPortStyle.Render(fmt.Sprintf("%s%s\n\n%s\n", kMsgMetadataTitleStyle.Render("Message Metadata"), headerViewPtr, m.msgMetadataVP.View()))
	}

	var msgValueVP string
	if !m.msgValueVPReady {
		msgValueVP = "\n  Initializing..."
	} else {
		msgValueVP = viewPortStyle.Render(fmt.Sprintf("%s%s\n\n%s\n", kMsgValueTitleStyle.Render("Message Value"), valueViewPtr, m.msgValueVP.View()))
	}
	var helpVP string
	if !m.helpVPReady {
		helpVP = "\n  Initializing..."
	} else {
		helpVP = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n", kMsgValueTitleStyle.Render("Help"), m.helpVP.View()))
	}

	switch m.vpFullScreen {
	case false:
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			stackListStyle.Render(m.kMsgsList.View()),
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
