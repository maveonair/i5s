package ui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestShellBehaviorForRunningAndStoppedInstances(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu"), stopped("debian")}
	m := newHarness(t, svc, 120, 25)

	m, cmd := pressSpecialKey(t, m, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("expected shell command for running instance")
	}
	m = completeShell(t, m, nil)
	requireVisible(t, m, "Shell exited")

	m, _ = pressSpecialKey(t, m, tea.KeyDown)
	m, _ = pressSpecialKey(t, m, tea.KeyEnter)
	requireVisible(t, m, "Instance must be running to open shell")
}

func TestShellExecCommandInvokesService(t *testing.T) {
	svc := newFakeService()
	cmd := shellCommand{service: svc, name: "ubuntu"}
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected shell command error: %v", err)
	}
	if len(svc.shellCalls) != 1 || svc.shellCalls[0] != "ubuntu" {
		t.Fatalf("expected shell call for ubuntu, got %v", svc.shellCalls)
	}
}

func TestShellFailureShowsError(t *testing.T) {
	svc := newFakeService()
	svc.shellErr = errors.New("agent unavailable")
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m = completeShell(t, m, svc.shellErr)
	requireVisible(t, m, "Failed to open shell: agent unavailable")
}

func TestEditConfigUsesSelectedInstance(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu"), running("debian")}
	m := newHarness(t, svc, 120, 25)
	m, _ = pressSpecialKey(t, m, tea.KeyDown)

	cmd, ok := m.selectedEditConfigCommand()
	if !ok {
		t.Fatal("expected edit config command for selected instance")
	}
	if cmd.name != "debian" {
		t.Fatalf("expected edit command for debian, got %q", cmd.name)
	}

	m, teaCmd := pressKey(t, m, "e")
	if teaCmd == nil {
		t.Fatal("expected Bubble Tea exec command")
	}

	m = applyMessage(t, m, editDoneMsg{})
	requireVisible(t, m, "Config updated")
}

func TestEditConfigNoopsWithoutSelection(t *testing.T) {
	svc := newFakeService()
	m := newHarness(t, svc, 120, 25)

	_, cmd := pressKey(t, m, "e")
	if cmd != nil {
		t.Fatal("expected no edit command without a selected instance")
	}
	requireNoCalls(t, svc.editCalls)
}

func TestEditConfigBusyStateDoesNotCallService(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)
	m.busy = true

	m, cmd := pressKey(t, m, "e")
	if cmd != nil {
		t.Fatal("expected no edit command while busy")
	}
	requireNoCalls(t, svc.editCalls)
	requireVisible(t, m, "Operation already in progress")
}

func TestEditConfigCommandInvokesService(t *testing.T) {
	svc := newFakeService()
	cmd := editConfigCommand{service: svc, name: "ubuntu"}
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected edit command error: %v", err)
	}
	if len(svc.editCalls) != 1 || svc.editCalls[0] != "ubuntu" {
		t.Fatalf("expected edit call for ubuntu, got %v", svc.editCalls)
	}
}

func TestEditConfigFailureShowsError(t *testing.T) {
	svc := newFakeService()
	svc.editErr = errors.New("invalid config")
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m = applyMessage(t, m, editDoneMsg{err: svc.editErr})
	requireVisible(t, m, "Failed to edit config: invalid config")
}
