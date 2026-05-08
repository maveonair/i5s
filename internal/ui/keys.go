package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	if m.filtering {
		switch key {
		case "esc":
			m.filtering = false
			m.filter = ""
			m.selectedIndex = 0
			return m, nil
		case "enter":
			m.filtering = false
			return m, nil
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.selectedIndex = 0
			}
			return m, nil
		default:
			if len(key) == 1 {
				m.filter += key
				m.selectedIndex = 0
			}
			return m, nil
		}
	}

	switch m.view {
	case ViewInstances:
		return m.handleInstanceKey(key)
	case ViewLogs, ViewConsole, ViewHelp:
		return m.handleViewportKey(key)
	case ViewRemotePicker, ViewProjectPicker, ViewLogPicker:
		return m.handlePickerKey(key)
	case ViewConfirmStop, ViewConfirmDelete, ViewConfirmClearConsole:
		return m.handleConfirmKey(key)
	}

	return m, nil
}

func (m Model) handleInstanceKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "?":
		m.view = ViewHelp
		m.viewport.SetContent(helpText())
		m.viewport.GotoTop()
		return m, nil
	case "r":
		if m.loading {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		m.loading = true
		m.operation = "Refreshing..."
		return m, m.refreshCmd()
	case "/":
		m.filtering = true
		return m, nil
	case "down", "j":
		m.moveSelection(1)
		return m, m.loadSelectedStateCmd()
	case "up", "k":
		m.moveSelection(-1)
		return m, m.loadSelectedStateCmd()
	case "enter":
		if m.busy {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		selected := m.selected()
		if selected.Name == "" {
			return m, nil
		}
		if !selected.IsRunning() {
			m.setError("Instance must be running to open shell")
			return m, nil
		}
		return m, tea.Exec(shellCommand{service: m.service, name: selected.Name}, func(err error) tea.Msg { return shellDoneMsg{err: err} })
	case "e":
		if m.busy {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		cmd, ok := m.selectedEditConfigCommand()
		if !ok {
			return m, nil
		}
		return m, tea.Exec(cmd, func(err error) tea.Msg { return editDoneMsg{err: err} })
	case "l":
		selected := m.selected()
		if selected.Name == "" {
			return m, nil
		}
		return m, m.loadLogFilesCmd(selected.Name)
	case "c":
		selected := m.selected()
		if selected.Name == "" {
			return m, nil
		}
		m.view = ViewConsole
		m.selectedLog = "console"
		return m, m.loadConsoleLogCmd()
	case "s":
		if m.busy {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		selected := m.selected()
		if selected.Name == "" {
			return m, nil
		}
		if !selected.IsRunning() {
			m.setStatus("Instance is already stopped")
			return m, nil
		}
		m.view = ViewConfirmStop
		return m, nil
	case "S":
		if m.busy {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		selected := m.selected()
		if selected.Name == "" {
			return m, nil
		}
		if selected.IsRunning() {
			m.setStatus("Instance is already running")
			return m, nil
		}
		m.loading = true
		m.busy = true
		m.operation = "Starting " + selected.Name + "..."
		return m, m.startCmd(selected.Name)
	case "d":
		if m.busy {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		selected := m.selected()
		if selected.Name == "" {
			return m, nil
		}
		if !selected.IsStopped() {
			m.setError("Instance must be stopped before it can be deleted")
			return m, nil
		}
		m.view = ViewConfirmDelete
		return m, nil
	case "R":
		m.pickerIndex = indexRemote(m.remotes, m.service.CurrentRemote())
		m.view = ViewRemotePicker
		return m, m.loadRemotesCmd()
	case "p":
		m.pickerIndex = indexProject(m.projects, m.service.CurrentProject())
		m.view = ViewProjectPicker
		return m, m.loadProjectsCmd()
	}

	return m, nil
}

func (m Model) handleViewportKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc":
		m.view = ViewInstances
		return m, nil
	case "r":
		if m.view == ViewConsole {
			return m, m.loadConsoleLogCmd()
		}
		if m.view == ViewLogs && m.selectedLog != "" {
			return m, m.loadLogCmd(m.selectedLog)
		}
	case "d":
		if m.view == ViewConsole {
			m.view = ViewConfirmClearConsole
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(tea.KeyMsg{Type: msgTypeForKey(key)})
	return m, cmd
}

func (m Model) handlePickerKey(key string) (tea.Model, tea.Cmd) {
	count := m.pickerCount()
	switch key {
	case "esc":
		m.view = ViewInstances
		return m, nil
	case "down", "j":
		if count > 0 && m.pickerIndex < count-1 {
			m.pickerIndex++
		}
	case "up", "k":
		if m.pickerIndex > 0 {
			m.pickerIndex--
		}
	case "enter":
		switch m.view {
		case ViewRemotePicker:
			if len(m.remotes) == 0 {
				return m, nil
			}
			name := m.remotes[m.pickerIndex].Name
			m.loading = true
			m.busy = true
			m.operation = "Switching remote..."
			return m, m.switchRemoteCmd(name)
		case ViewProjectPicker:
			if len(m.projects) == 0 {
				return m, nil
			}
			name := m.projects[m.pickerIndex].Name
			m.loading = true
			m.busy = true
			m.operation = "Switching project..."
			return m, m.switchProjectCmd(name)
		case ViewLogPicker:
			if len(m.logFiles) == 0 {
				return m, nil
			}
			m.selectedLog = m.logFiles[m.pickerIndex]
			m.view = ViewLogs
			return m, m.loadLogCmd(m.selectedLog)
		}
	}
	return m, nil
}

func (m Model) handleConfirmKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "n", "esc", "q":
		m.view = ViewInstances
		return m, nil
	case "y":
		if m.busy {
			m.setStatus("Operation already in progress")
			return m, nil
		}
		selected := m.selected()
		switch m.view {
		case ViewConfirmStop:
			m.view = ViewInstances
			m.loading = true
			m.busy = true
			m.operation = "Stopping " + selected.Name + "..."
			return m, m.stopCmd(selected.Name)
		case ViewConfirmDelete:
			m.view = ViewInstances
			m.loading = true
			m.busy = true
			m.operation = "Deleting " + selected.Name + "..."
			return m, m.deleteCmd(selected.Name)
		case ViewConfirmClearConsole:
			m.view = ViewConsole
			return m, m.clearConsoleCmd()
		}
	}
	return m, nil
}

func (m Model) selectedEditConfigCommand() (editConfigCommand, bool) {
	selected := m.selected()
	if selected.Name == "" {
		return editConfigCommand{}, false
	}
	return editConfigCommand{service: m.service, name: selected.Name}, true
}

func indexRemote(rows []appincus.RemoteRow, name string) int {
	for i, row := range rows {
		if row.Name == name {
			return i
		}
	}
	return 0
}

func indexProject(rows []appincus.ProjectRow, name string) int {
	for i, row := range rows {
		if row.Name == name {
			return i
		}
	}
	return 0
}

func msgTypeForKey(key string) tea.KeyType {
	switch key {
	case "up", "k":
		return tea.KeyUp
	case "down", "j":
		return tea.KeyDown
	case "pgup":
		return tea.KeyPgUp
	case "pgdown":
		return tea.KeyPgDown
	case "home", "g":
		return tea.KeyHome
	case "end", "G":
		return tea.KeyEnd
	default:
		return tea.KeyRunes
	}
}
