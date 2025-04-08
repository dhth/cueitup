package ui

import "github.com/charmbracelet/lipgloss"

const (
	defaultForegroundColor = "#282828"
	activeHeaderColor      = "#fe8019"
	inactivePaneColor      = "#928374"
	listPaneBorderColor    = "#3c3836"
	helpMsgColor           = "#83a598"
	helpViewTitleColor     = "#83a598"
	helpHeaderColor        = "#83a598"
	helpSectionColor       = "#fabd2f"
	cueitupColor           = "#d3869b"
	persistingColor        = "#fb4934"
	deletingMsgsColor      = "#d3869b"
	skippingColor          = "#fabd2f"
	pollingColor           = "#928374"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(defaultForegroundColor))

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingLeft(1).PaddingBottom(1)

	msgListStyle = baseListStyle.
			Width(listWidth+5).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(listPaneBorderColor))

	msgValueVPStyle   = baseListStyle.Width(150).PaddingLeft(3)
	msgValueVPFSStyle = lipgloss.NewStyle().
				PaddingTop(1).
				PaddingRight(2).
				PaddingLeft(1).
				PaddingBottom(1)

	helpVPStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	modeStyle = baseStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(cueitupColor))

	msgValueTitleStyle = baseStyle.
				Bold(true).
				Background(lipgloss.Color(inactivePaneColor)).
				Align(lipgloss.Left)

	persistingStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color(persistingColor))

	deletingMsgsStyle = baseStyle.
				Bold(true).
				Foreground(lipgloss.Color(deletingMsgsColor))

	skippingStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color(skippingColor))

	pollingMsgStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color(pollingColor))

	helpMsgStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color(helpMsgColor))

	helpVPTitleStyle = baseStyle.
				Bold(true).
				Background(lipgloss.Color(helpViewTitleColor)).
				Align(lipgloss.Left)

	helpHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(helpHeaderColor))

	helpSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(helpSectionColor))
)
