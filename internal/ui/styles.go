package ui

import "github.com/charmbracelet/lipgloss"

const (
	defaultForegroundColor = "#282828"
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
	errorColor             = "#fb4934"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(defaultForegroundColor))

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingBottom(1)

	msgListStyle = baseListStyle.
			Width(listWidth+5).
			Border(lipgloss.ThickBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color(listPaneBorderColor))

	msgValueVPStyle = baseListStyle.PaddingLeft(4)

	helpVPStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
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

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(errorColor))
)
