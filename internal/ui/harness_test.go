package ui

import (
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func newHarness(t *testing.T, service *fakeService, width int, height int) Model {
	t.Helper()
	m := New(service, time.Hour)
	m = resize(t, m, width, height)
	return receiveInstances(t, m, service.instances)
}

func resize(t *testing.T, m Model, width int, height int) Model {
	t.Helper()
	model, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: height})
	return model.(Model)
}

func receiveInstances(t *testing.T, m Model, rows []appincus.InstanceRow) Model {
	t.Helper()
	model, _ := m.Update(instancesLoadedMsg{instances: rows})
	return model.(Model)
}

func receiveRemotes(t *testing.T, m Model, rows []appincus.RemoteRow) Model {
	t.Helper()
	model, _ := m.Update(remotesLoadedMsg{remotes: rows})
	return model.(Model)
}

func receiveProjects(t *testing.T, m Model, rows []appincus.ProjectRow) Model {
	t.Helper()
	model, _ := m.Update(projectsLoadedMsg{projects: rows})
	return model.(Model)
}

func completeShell(t *testing.T, m Model, err error) Model {
	t.Helper()
	model, _ := m.Update(shellDoneMsg{err: err})
	return model.(Model)
}

func pressKey(t *testing.T, m Model, key string) (Model, tea.Cmd) {
	t.Helper()
	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	return model.(Model), cmd
}

func pressSpecialKey(t *testing.T, m Model, key tea.KeyType) (Model, tea.Cmd) {
	t.Helper()
	model, cmd := m.Update(tea.KeyMsg{Type: key})
	return model.(Model), cmd
}

func pressKeyAndRun(t *testing.T, m Model, key string) (Model, tea.Msg) {
	t.Helper()
	m, cmd := pressKey(t, m, key)
	return m, runCommand(t, cmd)
}

func pressSpecialKeyAndRun(t *testing.T, m Model, key tea.KeyType) (Model, tea.Msg) {
	t.Helper()
	m, cmd := pressSpecialKey(t, m, key)
	return m, runCommand(t, cmd)
}

func applyMessage(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	model, _ := m.Update(msg)
	return model.(Model)
}

func applyMessageAndRun(t *testing.T, m Model, msg tea.Msg) (Model, tea.Msg) {
	t.Helper()
	model, cmd := m.Update(msg)
	return model.(Model), runCommand(t, cmd)
}

func runCommand(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected command")
	}
	return cmd()
}

func runCommandIgnoringMessage(t *testing.T, cmd tea.Cmd) {
	t.Helper()
	_ = runCommand(t, cmd)
}

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func plain(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

func visible(m Model) string {
	return plain(m.View())
}

func requireVisible(t *testing.T, m Model, want string) {
	t.Helper()
	output := visible(m)
	if !strings.Contains(output, want) {
		t.Fatalf("expected output to contain %q\noutput:\n%s", want, output)
	}
}

func requireNotVisible(t *testing.T, m Model, want string) {
	t.Helper()
	output := visible(m)
	if strings.Contains(output, want) {
		t.Fatalf("expected output not to contain %q\noutput:\n%s", want, output)
	}
}

func requireNoCalls[T any](t *testing.T, calls []T) {
	t.Helper()
	if len(calls) != 0 {
		t.Fatalf("expected no calls, got %v", calls)
	}
}
