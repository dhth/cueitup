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
				m.activeView = msgsListView
			case helpView:
				m.activeView = m.lastView
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
			}
		case "]", "l":
			if m.activeView == msgValueView {
				m.msgsList.CursorDown()
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
		case "ctrl+r":
			if m.activeView == msgsListView {
				m.msgsList.SetItems(make([]list.Item, 0))
				m.msgValueVP.SetContent("")
				m.firstFetch = true
			}
		case "tab":
			switch m.activeView {
			case msgsListView:
				m.activeView = msgValueView
			case msgValueView:
				m.activeView = msgsListView
			}
		case "shift+tab":
			switch m.activeView {
			case msgsListView:
				m.activeView = msgsListView
			case msgValueView:
				m.activeView = msgsListView
			}
		}

	case tea.WindowSizeMsg:
		w, h := msgListStyle.GetFrameSize()
		w2, _ := msgValueVPStyle.GetFrameSize()
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width - 1
		m.msgsList.SetHeight(msg.Height - h - 2)

		if !m.msgValueVPReady {
			m.msgValueVP = viewport.New(msg.Width-2-w-w2-listWidth, m.terminalHeight-12)
			m.msgValueVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.msgValueVPReady = true
		} else {
			m.msgValueVP.Width = msg.Width - 2 - w - w2 - listWidth
			m.msgValueVP.Height = msg.Height - 12
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(msg.Width-1, msg.Height-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(HelpText)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = msg.Width - 1
			m.helpVP.Height = msg.Height - 7
		}

	case HideHelpMsg:
		m.showHelpIndicator = false

	case SQSMsgsFetchedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			if !m.behaviours.SkipMessages {
				for _, message := range msg.messages {
					m.msgsList.InsertItem(len(m.msgsList.Items()), message)

					if m.behaviours.PersistMessages {
						cmds = append(cmds,
							saveMessageToDisk(
								message.ID,
								message.Body,
								m.config.Format,
								m.persistDir,
							),
						)
					}
				}

				if m.behaviours.DeleteMessages {
					cmds = append(cmds,
						DeleteMessages(m.sqsClient,
							m.queueURL,
							msg.sqsMessages),
					)
				}
			}
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
	}

	if m.activeView == msgsListView || m.activeView == msgValueView {
		if len(m.msgsList.Items()) > 0 && m.msgsList.Index() != m.msgListCurrentIndex {
			m.msgListCurrentIndex = m.msgsList.Index()
			message, ok := m.msgsList.SelectedItem().(t.Message)

			if ok {
				var vpContent string
				if message.Err != nil {
					vpContent = errorStyle.Render(fmt.Sprintf("error: %s", message.Err.Error()))
				} else {
					switch m.config.Format {
					case t.JSON:
						vpContent = string(pretty.Color([]byte(message.Body), nil))
					default:
						vpContent = message.Body
					}
				}
				m.msgValueVP.SetContent(vpContent)
			}

		}
	}

	return m, tea.Batch(cmds...)
}
