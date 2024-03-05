package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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
			return m, FetchMessages(m.sqsClient, m.queueUrl, 1, 0, m.extractJSONObject, m.keyProperty)
		case "N":
			m.msg = " ..."
			for i := 0; i < 10; i++ {
				cmds = append(cmds, FetchMessages(m.sqsClient, m.queueUrl, 1, 0, m.extractJSONObject, m.keyProperty))
			}
			return m, tea.Batch(cmds...)
		case "}":
			m.msg = " ..."
			for i := 0; i < 100; i++ {
				cmds = append(cmds, FetchMessages(m.sqsClient, m.queueUrl, 1, 0, m.extractJSONObject, m.keyProperty))
			}
			return m, tea.Batch(cmds...)
		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
			m.vpFullScreen = true
			if m.helpSeen < 2 {
				m.helpSeen += 1
			}
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
			m.kMsgsList.CursorUp()
			m.msgValueVP.SetContent(m.recordValueStore[m.kMsgsList.SelectedItem().FilterValue()])
		case "]":
			m.kMsgsList.CursorDown()
			m.msgValueVP.SetContent(m.recordValueStore[m.kMsgsList.SelectedItem().FilterValue()])

		case "ctrl+p":
			m.pollForQueueMsgCount = !m.pollForQueueMsgCount
			if m.pollForQueueMsgCount {
				return m, tea.Batch(GetQueueMsgCount(m.sqsClient, m.queueUrl), tickEvery(msgCountTickInterval))
			}
			return m, nil
		case "ctrl+r":
			deleteMsgsFlag := m.deleteMsgs
			persistMsgsFlag := m.persistRecords
			m = InitialModel(m.sqsClient, m.queueUrl, m.extractJSONObject, m.keyProperty)
			m.deleteMsgs = deleteMsgsFlag
			m.persistRecords = persistMsgsFlag
			m.msgValueVP.SetContent("")
			return m, nil
		case "1":
			m.msgValueVP.Height = m.terminalHeight - 7
			m.vpFullScreen = true
			m.lastView = kMsgsListView
			m.activeView = kMsgValueView
			return m, nil
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
				return m, nil
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
		_, h := stackListStyle.GetFrameSize()
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.kMsgsList.SetHeight(msg.Height - h - 2)

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
			m.helpVP.SetContent(helpText)
			m.helpVPReady = true
		}
	case KMsgValueReadyMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.recordValueStore[msg.storeKey] = msg.msgValue
		}
		return m, tea.Batch(cmds...)

	case KMsgFetchedMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			switch m.skipRecords {
			case false:
				for i, message := range msg.messages {
					m.kMsgsList.InsertItem(len(m.kMsgsList.Items()),
						KMsgItem{message: message,
							messageValue:     msg.messageValues[i],
							keyPropertyName:  m.keyProperty,
							keyPropertyValue: msg.keyValues[i],
						},
					)
					if m.deleteMsgs {
						cmds = append(cmds, DeleteMessages(m.sqsClient, m.queueUrl, msg.messages))
					}
					m.recordValueStore[*message.MessageId] = msg.messageValues[i]
				}
			}
		}
	case KMsgChosenMsg:
		m.msgValueVP.SetContent(m.recordValueStore[msg.key])
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
			m.kMsgsList.Title = fmt.Sprintf("Messages (%d in queue)", msg.approxMsgCount)
		}
	}

	switch m.activeView {
	case kMsgsListView:
		m.kMsgsList, cmd = m.kMsgsList.Update(msg)
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
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
