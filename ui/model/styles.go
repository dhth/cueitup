package model

import "github.com/charmbracelet/lipgloss"

const (
	defaultBackgroundColor = "#282828"
	listColor              = "#fe8019"
	activeHeaderColor      = "#fe8019"
	inactivePaneColor      = "#928374"
	listPaneBorderColor    = "#3c3836"
	helpMsgColor           = "#83a598"
	helpViewTitleColor     = "#83a598"
	helpHeaderColor        = "#83a598"
	helpSectionColor       = "#fabd2f"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#282828"))

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
			Background(lipgloss.Color("#b8bb26"))

	msgValueTitleStyle = baseStyle.
				Bold(true).
				Background(lipgloss.Color(inactivePaneColor)).
				Align(lipgloss.Left)

	persistingStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color("#fb4934"))

	deletingMsgsStyle = baseStyle.
				Bold(true).
				Foreground(lipgloss.Color("#d3869b"))

	skippingStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color("#fabd2f"))

	pollingMsgStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color("#928374"))

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
