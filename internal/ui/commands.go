package ui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const serviceReadTimeout = 30 * time.Second

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func serviceContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), serviceReadTimeout)
}

func (m Model) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		rows, err := m.service.ListInstances(ctx)
		if err != nil {
			return opDoneMsg{err: fmt.Errorf("could not find any instances for the current remote: %w", err)}
		}
		return instancesLoadedMsg{instances: rows}
	}
}

func (m Model) loadRemotesCmd() tea.Cmd {
	return func() tea.Msg {
		rows, err := m.service.ListRemotes()
		if err != nil {
			return opDoneMsg{err: err}
		}
		return remotesLoadedMsg{remotes: rows}
	}
}

func (m Model) loadProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		rows, err := m.service.ListProjects(ctx)
		if err != nil {
			return opDoneMsg{err: fmt.Errorf("could not find any instances for the current remote: %w", err)}
		}
		return projectsLoadedMsg{projects: rows}
	}
}

func (m Model) loadSelectedStateCmd() tea.Cmd {
	selected := m.selected()
	if selected.Name == "" || !selected.IsRunning() {
		return nil
	}
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		state, _ := m.service.GetInstanceState(ctx, selected.Name)
		return stateLoadedMsg{name: selected.Name, state: state}
	}
}

func (m Model) loadLogFilesCmd(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		files, err := m.service.ListLogs(ctx, name)
		if err != nil {
			return opDoneMsg{err: err}
		}
		return logFilesLoadedMsg{files: files}
	}
}

func (m Model) loadLogCmd(logName string) tea.Cmd {
	selected := m.selected()
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		content, err := m.service.GetLog(ctx, selected.Name, logName)
		return logLoadedMsg{content: content, err: err}
	}
}

func (m Model) loadConsoleLogCmd() tea.Cmd {
	selected := m.selected()
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		content, err := m.service.GetConsoleLog(ctx, selected.Name)
		return logLoadedMsg{content: content, err: err}
	}
}

func (m Model) startCmd(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), lifecycleTimeout)
		defer cancel()
		return opDoneMsg{message: "Started " + name, err: m.service.StartInstance(ctx, name), refresh: true}
	}
}

func (m Model) stopCmd(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), lifecycleTimeout)
		defer cancel()
		return opDoneMsg{message: "Stopped " + name, err: m.service.StopInstance(ctx, name), refresh: true}
	}
}

func (m Model) deleteCmd(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), lifecycleTimeout)
		defer cancel()
		return opDoneMsg{message: "Deleted " + name, err: m.service.DeleteInstance(ctx, name), refresh: true}
	}
}

func (m Model) clearConsoleCmd() tea.Cmd {
	selected := m.selected()
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		err := m.service.ClearConsoleLog(ctx, selected.Name)
		if err != nil {
			return opDoneMsg{err: err}
		}
		return opDoneMsg{message: "Cleared console log"}
	}
}

func (m Model) switchRemoteCmd(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		return switchDoneMsg{message: "Switched remote to " + name, err: m.service.SwitchRemote(ctx, name)}
	}
}

func (m Model) switchProjectCmd(name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := serviceContext()
		defer cancel()

		return switchDoneMsg{message: "Switched project to " + name, err: m.service.SwitchProject(ctx, name)}
	}
}
