package ui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestRemotePickerSwitchesRemote(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.remotes = []appincus.RemoteRow{{Name: "local"}, {Name: "prod"}}
	m := newHarness(t, svc, 120, 25)
	m = receiveRemotes(t, m, svc.remotes)

	m, _ = pressKey(t, m, "R")
	requireVisible(t, m, "Switch Remote")
	m, _ = pressSpecialKey(t, m, tea.KeyDown)
	_, cmd := pressSpecialKey(t, m, tea.KeyEnter)
	runCommandIgnoringMessage(t, cmd)
	if len(svc.switchRemoteCalls) != 1 || svc.switchRemoteCalls[0] != "prod" {
		t.Fatalf("expected switch remote prod, got %v", svc.switchRemoteCalls)
	}
}

func TestProjectPickerSwitchesProject(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.projects = []appincus.ProjectRow{{Name: "default"}, {Name: "dev"}}
	m := newHarness(t, svc, 120, 25)
	m = receiveProjects(t, m, svc.projects)

	m, _ = pressKey(t, m, "p")
	requireVisible(t, m, "Switch Project")
	m, _ = pressSpecialKey(t, m, tea.KeyDown)
	_, cmd := pressSpecialKey(t, m, tea.KeyEnter)
	runCommandIgnoringMessage(t, cmd)
	if len(svc.switchProjectCalls) != 1 || svc.switchProjectCalls[0] != "dev" {
		t.Fatalf("expected switch project dev, got %v", svc.switchProjectCalls)
	}
}

func TestPickerCancelDoesNotCallService(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.remotes = []appincus.RemoteRow{{Name: "local"}, {Name: "prod"}}
	m := newHarness(t, svc, 120, 25)
	m = receiveRemotes(t, m, svc.remotes)

	m, _ = pressKey(t, m, "R")
	m, _ = pressSpecialKey(t, m, tea.KeyEsc)
	requireNotVisible(t, m, "Switch Remote")
	requireNoCalls(t, svc.switchRemoteCalls)
}

func TestSwitchFailureShowsError(t *testing.T) {
	svc := newFakeService()
	svc.switchRemoteErr = errors.New("remote unavailable")
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.remotes = []appincus.RemoteRow{{Name: "local"}, {Name: "prod"}}
	m := newHarness(t, svc, 120, 25)
	m = receiveRemotes(t, m, svc.remotes)

	m, _ = pressKey(t, m, "R")
	m, _ = pressSpecialKey(t, m, tea.KeyDown)
	m, msg := pressSpecialKeyAndRun(t, m, tea.KeyEnter)
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "remote unavailable")
}
