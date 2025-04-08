package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		keyMsg, keyMsgOK := msg.(tea.KeyMsg)
		if !keyMsgOK {
			return nil
		}
		if key.Matches(keyMsg,
			list.DefaultKeyMap().CursorUp,
			list.DefaultKeyMap().CursorDown,
			list.DefaultKeyMap().GoToStart,
			list.DefaultKeyMap().GoToEnd,
			list.DefaultKeyMap().NextPage,
			list.DefaultKeyMap().PrevPage,
		) {
			selected := m.SelectedItem()
			if selected == nil {
				return nil
			}
			key := selected.FilterValue()
			return showItemDetails(key)
		}
		return nil
	}

	return d
}
