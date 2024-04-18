package model

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tidwall/pretty"
)

const useHighPerformanceRenderer = false

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.msg = ""
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.activeView == contextualSearchView {
				m.activeView = m.lastView
				return m, tea.Batch(cmds...)
			}
		case "enter":
			if m.activeView == contextualSearchView {
				m.activeView = m.lastView
				if len(m.contextSearchInput.Value()) > 0 {
					return m, setContextSearchValues(m.contextSearchInput.Value())
				}
			}
		}
	}

	switch m.activeView {
	case contextualSearchView:
		m.contextSearchInput, cmd = m.contextSearchInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.vpFullScreen == false {
				return m, tea.Quit
			}
			m.msgMetadataVP.Height = m.terminalHeight/2 - 8
			m.msgValueVP.Height = m.terminalHeight - 8
			m.vpFullScreen = false
			m.activeView = kMsgsListView
			return m, nil
		case "n", " ":
			m.msg = " ..."
			return m,
				m.FetchMessages(1, 0)
		case "N":
			m.msg = " ..."
			for i := 0; i < 10; i++ {
				cmds = append(cmds,
					m.FetchMessages(1, 0),
				)
			}
			return m, tea.Batch(cmds...)
		case "}":
			m.msg = " ..."
			for i := 0; i < 20; i++ {
				cmds = append(cmds,
					m.FetchMessages(5, 0),
				)
			}
			return m, tea.Batch(cmds...)
		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
			m.vpFullScreen = true
			return m, nil
		case "d":
			if m.activeView == kMsgsListView {
				m.deleteMsgs = !m.deleteMsgs
			}
			return m, nil
		case "p":
			if m.persistRecords == false {
				m.skipRecords = false
			}
			m.persistRecords = !m.persistRecords
			return m, nil
		case "s":
			if m.skipRecords == false {
				m.persistRecords = false
			}
			m.skipRecords = !m.skipRecords
			return m, nil
		case "[":
			m.msgsList.CursorUp()
			m.msgValueVP.SetContent(m.recordValueStore[m.msgsList.SelectedItem().FilterValue()])
		case "]":
			m.msgsList.CursorDown()
			m.msgValueVP.SetContent(m.recordValueStore[m.msgsList.SelectedItem().FilterValue()])

		case "ctrl+p":
			m.pollForQueueMsgCount = !m.pollForQueueMsgCount
			if m.pollForQueueMsgCount {
				return m,
					tea.Batch(GetQueueMsgCount(m.sqsClient,
						m.queueUrl),
						tickEvery(msgCountTickInterval),
					)
			}
		case "ctrl+s":
			if m.activeView != contextualSearchView {
				m.lastView = m.activeView
				m.activeView = contextualSearchView
			}
			return m, tea.Batch(cmds...)
		case "ctrl+f":
			if len(m.contextSearchValues) > 0 {
				m.filterMessages = !m.filterMessages
			}
		case "ctrl+r":
			m.msgsList.SetItems(make([]list.Item, 0))
			m.msgValueVP.SetContent("")
		case "1":
			m.msgValueVP.Height = m.terminalHeight - 7
			m.vpFullScreen = true
			m.lastView = kMsgsListView
			m.activeView = kMsgValueView
		case "f":
			switch m.activeView {
			case kMsgMetadataView:
				switch m.vpFullScreen {
				case false:
					m.msgMetadataVP.Height = m.terminalHeight - 7
					m.lastView = kMsgMetadataView
					m.vpFullScreen = true
				case true:
					m.msgMetadataVP.Height = m.terminalHeight/2 - 8
					m.msgValueVP.Height = m.terminalHeight - 8
					m.vpFullScreen = false
					m.activeView = m.lastView
				}
			case kMsgValueView:
				switch m.vpFullScreen {
				case false:
					m.msgValueVP.Height = m.terminalHeight - 7
					m.lastView = kMsgValueView
					m.vpFullScreen = true
				case true:
					m.msgValueVP.Height = m.terminalHeight - 8
					m.msgMetadataVP.Height = m.terminalHeight/2 - 8
					m.vpFullScreen = false
					m.activeView = m.lastView
				}
			}
		case "tab":
			if m.vpFullScreen {
				return m, nil
			}
			if m.activeView == kMsgsListView {
				m.activeView = kMsgValueView
			} else if m.activeView == kMsgValueView {
				m.activeView = kMsgsListView
			}
		case "shift+tab":
			if m.vpFullScreen {
				return m, nil
			}
			if m.activeView == kMsgsListView {
				m.activeView = kMsgsListView
			} else if m.activeView == kMsgValueView {
				m.activeView = kMsgsListView
			}
		}

	case tea.WindowSizeMsg:
		_, h := msgListStyle.GetFrameSize()
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.msgsList.SetHeight(msg.Height - h - 2)

		if !m.msgMetadataVPReady {
			m.msgMetadataVP = viewport.New(120, m.terminalHeight/2-8)
			m.msgMetadataVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.msgMetadataVPReady = true
		} else {
			m.msgMetadataVP.Width = 120
			m.msgMetadataVP.Height = 12
		}

		if !m.msgValueVPReady {
			m.msgValueVP = viewport.New(120, m.terminalHeight-8)
			m.msgValueVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.msgValueVPReady = true
		} else {
			m.msgValueVP.Width = 120
			m.msgValueVP.Height = 12
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(120, m.terminalHeight-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(HelpText)
			m.helpVPReady = true
		}
	case KMsgValueReadyMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.recordValueStore[msg.storeKey] = msg.msgValue
		}
		return m, tea.Batch(cmds...)

	case ContextSearchValuesSetMsg:
		m.contextSearchValues = msg.values
		m.contextSearchInput.SetValue("")
		m.filterMessages = true

	case HideHelpMsg:
		m.showHelpIndicator = false

	case KMsgFetchedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			switch m.skipRecords {
			case false:
				vPresenceMap := make(map[string]bool)
				if m.filterMessages && len(m.contextSearchValues) > 0 {
					for _, p := range m.contextSearchValues {
						vPresenceMap[p] = true
					}
				}
				for i, message := range msg.messages {

					// only save/persist values that are requested to be filtered
					if m.filterMessages && !(msg.keyValues[i] != "" && vPresenceMap[msg.keyValues[i]]) {
						continue
					}

					m.msgsList.InsertItem(len(m.msgsList.Items()),
						KMsgItem{message: message,
							messageValue:    msg.messageValues[i],
							contextKeyName:  m.msgConsumptionConf.ContextKey,
							contextKeyValue: msg.keyValues[i],
						},
					)
					m.recordValueStore[*message.MessageId] = msg.messageValues[i]
					if m.persistRecords {
						prefix := time.Now().Unix()
						filePath := fmt.Sprintf("%s/%d-%s.md", m.persistDir, prefix, *message.MessageId)
						cmds = append(cmds,
							saveRecordValueToDisk(
								filePath,
								msg.messageValues[i],
								m.msgConsumptionConf.Format,
							),
						)
					}
				}
				if m.deleteMsgs {
					cmds = append(cmds,
						DeleteMessages(m.sqsClient,
							m.queueUrl,
							msg.messages),
					)
				}
			}
		}
	case KMsgChosenMsg:
		switch m.deserializationFmt {
		case JsonFmt:
			result := string(pretty.Color([]byte(m.recordValueStore[msg.key]), nil))
			m.msgValueVP.SetContent(result)
		default:
			m.msgValueVP.SetContent(m.recordValueStore[msg.key])
		}
	case MsgCountTickMsg:
		cmds = append(cmds, GetQueueMsgCount(m.sqsClient, m.queueUrl))
		if m.pollForQueueMsgCount {
			cmds = append(cmds, tickEvery(msgCountTickInterval))
		}
		return m, tea.Batch(cmds...)
	case QueueMsgCountFetchedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.msgsList.Title = fmt.Sprintf("Messages (%d in queue)", msg.approxMsgCount)
		}
	}

	switch m.activeView {
	case kMsgsListView:
		m.msgsList, cmd = m.msgsList.Update(msg)
		cmds = append(cmds, cmd)
	case kMsgMetadataView:
		m.msgMetadataVP, cmd = m.msgMetadataVP.Update(msg)
		cmds = append(cmds, cmd)
	case kMsgValueView:
		m.msgValueVP, cmd = m.msgValueVP.Update(msg)
		cmds = append(cmds, cmd)
	case helpView:
		m.helpVP, cmd = m.helpVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
