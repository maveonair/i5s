package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestMultipleLogsOpenPickerThenSelectedLog(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.logFiles["ubuntu"] = []string{"lxc.log", "console.log"}
	svc.logs["ubuntu/lxc.log"] = "hello from lxc"
	m := newHarness(t, svc, 120, 25)

	m, msg := pressKeyAndRun(t, m, "l")
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "Select Log")
	requireVisible(t, m, "lxc.log")

	m, msg = pressSpecialKeyAndRun(t, m, tea.KeyEnter)
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "LOGS")
	requireVisible(t, m, "hello from lxc")
}

func TestSingleLogOpensDirectlyAndNoLogsShowsStatus(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.logFiles["ubuntu"] = []string{"lxc.log"}
	svc.logs["ubuntu/lxc.log"] = "single log content"
	m := newHarness(t, svc, 120, 25)

	m, msg := pressKeyAndRun(t, m, "l")
	m, msg = applyMessageAndRun(t, m, msg)
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "LOGS")
	requireVisible(t, m, "single log content")

	svc = newFakeService()
	svc.instances = []appincus.InstanceRow{running("empty")}
	m = newHarness(t, svc, 120, 25)
	m, msg = pressKeyAndRun(t, m, "l")
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "No logs found")
}

func TestLogAndConsoleRefreshCallServiceAgain(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.logFiles["ubuntu"] = []string{"lxc.log"}
	svc.logs["ubuntu/lxc.log"] = "log content"
	svc.consoleLog["ubuntu"] = "console content"
	m := newHarness(t, svc, 120, 25)

	m, msg := pressKeyAndRun(t, m, "l")
	m, msg = applyMessageAndRun(t, m, msg)
	m = applyMessage(t, m, msg)
	_, cmd := pressKey(t, m, "r")
	runCommandIgnoringMessage(t, cmd)
	if len(svc.getLogCalls) != 2 {
		t.Fatalf("expected two log reads, got %v", svc.getLogCalls)
	}

	m = newHarness(t, svc, 120, 25)
	m, msg = pressKeyAndRun(t, m, "c")
	m = applyMessage(t, m, msg)
	_, cmd = pressKey(t, m, "r")
	runCommandIgnoringMessage(t, cmd)
	if len(svc.getConsoleCalls) != 2 {
		t.Fatalf("expected two console reads, got %v", svc.getConsoleCalls)
	}
}

func TestConsoleLogAndClearConfirmation(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	svc.consoleLog["ubuntu"] = "console output"
	m := newHarness(t, svc, 120, 25)

	m, msg := pressKeyAndRun(t, m, "c")
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "CONSOLE")
	requireVisible(t, m, "console output")

	m, _ = pressKey(t, m, "d")
	requireVisible(t, m, "Clear Console Log")
	_, cmd := pressKey(t, m, "y")
	runCommandIgnoringMessage(t, cmd)
	if len(svc.clearConsoleCalls) != 1 || svc.clearConsoleCalls[0] != "ubuntu" {
		t.Fatalf("expected clear console call for ubuntu, got %v", svc.clearConsoleCalls)
	}
}
