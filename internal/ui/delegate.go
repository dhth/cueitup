package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newAppItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(cueitupColor)).
		BorderLeftForeground(lipgloss.Color(cueitupColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
