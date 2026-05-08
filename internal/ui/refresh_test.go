package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestRefreshPreservesSelectionByName(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("alpha"), running("bravo")}
	m := newHarness(t, svc, 120, 25)
	m, _ = pressSpecialKey(t, m, tea.KeyDown)

	m = receiveInstances(t, m, []appincus.InstanceRow{running("bravo"), running("charlie")})
	if !strings.Contains(visible(m), "› bravo") {
		t.Fatalf("expected selected bravo after refresh\noutput:\n%s", visible(m))
	}
}

func TestManualRefreshCallsListInstances(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)
	before := svc.listInstancesCalls

	_, cmd := pressKey(t, m, "r")
	runCommandIgnoringMessage(t, cmd)
	if svc.listInstancesCalls != before+1 {
		t.Fatalf("expected ListInstances call on refresh, before %d after %d", before, svc.listInstancesCalls)
	}
}
