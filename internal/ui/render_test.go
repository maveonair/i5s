package ui

import (
	"errors"
	"strings"
	"testing"
	"time"

	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestMainScreenRendersInstancesFullHeight(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu"), stopped("debian")}
	m := newHarness(t, svc, 120, 30)

	for _, want := range []string{"i5s", "INSTANCES", "remote local", "project default", "NAME", "STATE", "IPV4", "IPV6", "TYPE", "SNAPSHOTS", "ubuntu", "RUNNING", "10.0.3.15", "fd42::1", "container", "2", "config"} {
		requireVisible(t, m, want)
	}

	if got := len(strings.Split(m.View(), "\n")); got != 30 {
		t.Fatalf("expected 30 rendered lines, got %d", got)
	}
	requireVisible(t, m, "╰")
}

func TestHelpDocumentsConfigEdit(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "?")
	requireVisible(t, m, "edit instance config")
}

func TestResponsiveColumnsHideSnapshotsOnNarrowWidth(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 65, 20)

	requireVisible(t, m, "NAME")
	requireVisible(t, m, "STATE")
	requireVisible(t, m, "IPV4")
	requireNotVisible(t, m, "SNAPSHOTS")
}

func TestEmptyStateKeepsFullHeight(t *testing.T) {
	svc := newFakeService()
	m := newHarness(t, svc, 120, 20)

	requireVisible(t, m, "No instances found")
	if got := len(strings.Split(m.View(), "\n")); got != 20 {
		t.Fatalf("expected 20 rendered lines, got %d", got)
	}
}

func TestInitialRefreshErrorShowsEmptyTableAndRemoteError(t *testing.T) {
	svc := newFakeService()
	svc.listInstancesErr = errors.New("connection refused")
	m := New(svc, time.Hour)
	m = resize(t, m, 120, 20)

	msg := runCommand(t, m.refreshCmd())
	m = applyMessage(t, m, msg)

	requireVisible(t, m, "No instances found")
	requireVisible(t, m, "could not find any instances for the current remote")
	requireVisible(t, m, "connection refused")
}
