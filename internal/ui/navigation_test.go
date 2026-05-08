package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestFilterInputFiltersVisibleRows(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu"), stopped("debian")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "/")
	m, _ = pressKey(t, m, "u")
	m, _ = pressKey(t, m, "b")
	m, _ = pressSpecialKey(t, m, tea.KeyEnter)

	requireVisible(t, m, "ubuntu")
	requireNotVisible(t, m, "debian")
}

func TestNavigationBoundaries(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("alpha"), running("bravo")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressSpecialKey(t, m, tea.KeyUp)
	requireVisible(t, m, "› alpha")
	m, _ = pressSpecialKey(t, m, tea.KeyDown)
	requireVisible(t, m, "› bravo")
	m, _ = pressSpecialKey(t, m, tea.KeyDown)
	requireVisible(t, m, "› bravo")
}

func TestFilterEscClearsFilter(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu"), stopped("debian")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "/")
	m, _ = pressKey(t, m, "u")
	m, _ = pressSpecialKey(t, m, tea.KeyEsc)

	requireVisible(t, m, "ubuntu")
	requireVisible(t, m, "debian")
}

func TestHelpViewAndBack(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "?")
	requireVisible(t, m, "HELP")
	requireVisible(t, m, "Navigation")
	m, _ = pressSpecialKey(t, m, tea.KeyEsc)
	requireVisible(t, m, "INSTANCES")
}

func TestCtrlCQuitsFromModal(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)
	m, _ = pressKey(t, m, "s")

	_, msg := pressSpecialKeyAndRun(t, m, tea.KeyCtrlC)
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}
