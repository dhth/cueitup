package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	t "github.com/dhth/cueitup/internal/types"
	"github.com/tidwall/pretty"
)

const (
	useHighPerformanceRenderer = false
	fetchingIndicator          = " ..."
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	m.message = ""
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			switch m.activeView {
			case msgsListView:
				return m, tea.Quit
			case msgValueView:
				if m.vpFullScreen {
					m.vpFullScreen = false
					m.msgValueVP.Width = 120
				} else {
					m.activeView = msgsListView
				}
			case helpView:
				m.activeView = m.lastView
			}
		case "esc":
			if m.activeView == contextualSearchView {
				m.activeView = m.lastView
			}
		case "enter":
			if m.activeView == contextualSearchView {
				m.activeView = m.lastView
				if len(m.contextSearchInput.Value()) > 0 {
					cmds = append(cmds, setContextSearchValues(m.contextSearchInput.Value()))
				} else {
					m.filterMessages = false
				}
			}
		case "n", " ":
			m.message = fetchingIndicator
			cmds = append(cmds, m.FetchMessages(1, 0))
		case "N":
			m.message = fetchingIndicator
			for range 10 {
				cmds = append(cmds,
					m.FetchMessages(1, 0),
				)
			}
		case "}":
			m.message = fetchingIndicator
			for range 20 {
				cmds = append(cmds,
					m.FetchMessages(5, 0),
				)
			}
		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
		case "d":
			if m.activeView == msgsListView {
				m.behaviours.DeleteMessages = !m.behaviours.DeleteMessages
			}
		case "p":
			if m.activeView == msgsListView {
				m.behaviours.PersistMessages = !m.behaviours.PersistMessages
			}
		case "s":
			if m.activeView == msgsListView {
				m.behaviours.SkipMessages = !m.behaviours.SkipMessages
			}
		case "[", "h":
			if m.activeView == msgValueView {
				m.msgsList.CursorUp()
				selected := m.msgsList.SelectedItem()
				if selected != nil {
					result := string(pretty.Color([]byte(m.recordValueStore[selected.FilterValue()]), nil))
					m.msgValueVP.SetContent(result)
				}
			}
		case "]", "l":
			if m.activeView == msgValueView {
				m.msgsList.CursorDown()
				selected := m.msgsList.SelectedItem()
				if selected != nil {
					result := string(pretty.Color([]byte(m.recordValueStore[selected.FilterValue()]), nil))
					m.msgValueVP.SetContent(result)
				}
			}
		case "ctrl+p":
			m.pollForQueueMsgCount = !m.pollForQueueMsgCount
			if m.pollForQueueMsgCount {
				cmds = append(cmds,
					tea.Batch(GetQueueMsgCount(m.sqsClient, m.queueURL),
						tickEvery(msgCountTickInterval),
					),
				)
			}
		case "ctrl+s":
			if m.activeView == msgsListView {
				m.lastView = m.activeView
				m.activeView = contextualSearchView
			}
		case "ctrl+f":
			if len(m.contextSearchValues) > 0 {
				m.filterMessages = !m.filterMessages
			}
		case "ctrl+r":
			if m.activeView == msgsListView {
				m.msgsList.SetItems(make([]list.Item, 0))
				m.msgValueVP.SetContent("")
				m.firstFetch = true
				m.filterMessages = false
			}
		case "1":
			m.msgValueVP.Width = m.terminalWidth - 1
			m.msgValueVP.Height = m.terminalHeight - 7
			m.vpFullScreen = true
			m.lastView = msgsListView
			m.activeView = msgValueView
		case "tab":
			if !m.vpFullScreen {
				switch m.activeView {
				case msgsListView:
					m.activeView = msgValueView
				case msgValueView:
					m.activeView = msgsListView
				}
			}
		case "shift+tab":
			if !m.vpFullScreen {
				switch m.activeView {
				case msgsListView:
					m.activeView = msgsListView
				case msgValueView:
					m.activeView = msgsListView
				}
			}
		}

	case tea.WindowSizeMsg:
		_, h := msgListStyle.GetFrameSize()
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width - 1
		m.msgsList.SetHeight(msg.Height - h - 2)

		if !m.msgValueVPReady {
			m.msgValueVP = viewport.New(120, m.terminalHeight-7)
			m.msgValueVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.msgValueVPReady = true
		} else {
			m.msgValueVP.Width = 120
			m.msgValueVP.Height = 12
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(msg.Width, msg.Height-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(HelpText)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = msg.Width - 1
			m.helpVP.Height = msg.Height - 7
		}
	case KMsgValueReadyMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.recordValueStore[msg.storeKey] = msg.msgValue
		}

	case ContextSearchValuesSetMsg:
		m.contextSearchValues = msg.values
		m.contextSearchInput.SetValue("")
		m.filterMessages = true

	case HideHelpMsg:
		m.showHelpIndicator = false

	case SQSMsgFetchedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			if !m.behaviours.SkipMessages {
				vPresenceMap := make(map[string]bool)
				if m.filterMessages && len(m.contextSearchValues) > 0 {
					for _, p := range m.contextSearchValues {
						vPresenceMap[p] = true
					}
				}
				for i, message := range msg.messages {
					// only save/persist values that are requested to be filtered
					// nolint: staticcheck
					if m.filterMessages && !(msg.keyValues[i] != "" && vPresenceMap[msg.keyValues[i]]) {
						continue
					}

					m.msgsList.InsertItem(len(m.msgsList.Items()),
						msgItem{
							message:         message,
							messageValue:    msg.messageValues[i],
							contextKeyName:  m.profile.ContextKey,
							contextKeyValue: msg.keyValues[i],
						},
					)
					m.recordValueStore[*message.MessageId] = msg.messageValues[i]

					if m.behaviours.PersistMessages {
						cmds = append(cmds,
							saveMessageToDisk(
								*message.MessageId,
								msg.messageValues[i],
								m.profile.Format,
								m.persistDir,
							),
						)
					}
				}

				if m.behaviours.DeleteMessages {
					cmds = append(cmds,
						DeleteMessages(m.sqsClient,
							m.queueURL,
							msg.messages),
					)
				}

				if m.firstFetch {
					selected := m.msgsList.SelectedItem()
					if selected != nil {
						recordValue, ok := m.recordValueStore[selected.FilterValue()]
						if ok {
							switch m.profile.Format {
							case t.JSON:
								result := string(pretty.Color([]byte(recordValue), nil))
								m.msgValueVP.SetContent(result)
							case t.None:
								m.msgValueVP.SetContent(recordValue)
							}
							m.firstFetch = false
						}
					}
				}
			}
		}
	case MsgChosenMsg:
		switch m.profile.Format {
		case t.JSON:
			result := string(pretty.Color([]byte(m.recordValueStore[msg.key]), nil))
			m.msgValueVP.SetContent(result)
		default:
			m.msgValueVP.SetContent(m.recordValueStore[msg.key])
		}
	case MsgCountTickMsg:
		cmds = append(cmds, GetQueueMsgCount(m.sqsClient, m.queueURL))
		if m.pollForQueueMsgCount {
			cmds = append(cmds, tickEvery(msgCountTickInterval))
		}
	case QueueMsgCountFetchedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.msgsList.Title = fmt.Sprintf("Messages (%d in queue)", msg.approxMsgCount)
		}
	}

	var updateCmd tea.Cmd
	switch m.activeView {
	case msgsListView:
		m.msgsList, updateCmd = m.msgsList.Update(msg)
		cmds = append(cmds, updateCmd)
	case msgValueView:
		m.msgValueVP, updateCmd = m.msgValueVP.Update(msg)
		cmds = append(cmds, updateCmd)
	case helpView:
		m.helpVP, updateCmd = m.helpVP.Update(msg)
		cmds = append(cmds, updateCmd)
	case contextualSearchView:
		m.contextSearchInput, updateCmd = m.contextSearchInput.Update(msg)
		cmds = append(cmds, updateCmd)
	}

	return m, tea.Batch(cmds...)
}
