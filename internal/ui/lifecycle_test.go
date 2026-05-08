package ui

import (
	"errors"
	"testing"

	appincus "github.com/maveonair/i5s/internal/incus"
)

func TestStopRunningInstanceRequiresConfirmation(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "s")
	requireVisible(t, m, "Stop Instance")
	requireVisible(t, m, "Stop instance ubuntu?")

	m, _ = pressKey(t, m, "n")
	requireNoCalls(t, svc.stopCalls)

	m, _ = pressKey(t, m, "s")
	m, msg := pressKeyAndRun(t, m, "y")
	if len(svc.stopCalls) != 1 || svc.stopCalls[0] != "ubuntu" {
		t.Fatalf("expected stop call for ubuntu, got %v", svc.stopCalls)
	}
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "Refreshing")
}

func TestStopStoppedInstanceShowsStatusWithoutServiceCall(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{stopped("debian")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "s")
	requireVisible(t, m, "Instance is already stopped")
	requireNoCalls(t, svc.stopCalls)
}

func TestStartStoppedInstanceCallsService(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{stopped("debian")}
	m := newHarness(t, svc, 120, 25)

	m, msg := pressKeyAndRun(t, m, "S")
	if len(svc.startCalls) != 1 || svc.startCalls[0] != "debian" {
		t.Fatalf("expected start call for debian, got %v", svc.startCalls)
	}
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "Refreshing")
}

func TestStartRunningInstanceShowsStatusWithoutServiceCall(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "S")
	requireVisible(t, m, "Instance is already running")
	requireNoCalls(t, svc.startCalls)
}

func TestDeleteRunningInstanceIsRefused(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "d")
	requireVisible(t, m, "Instance must be stopped before it can be deleted")
	requireNoCalls(t, svc.deleteCalls)
}

func TestDeleteStoppedInstanceRequiresConfirmation(t *testing.T) {
	svc := newFakeService()
	svc.instances = []appincus.InstanceRow{stopped("debian")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "d")
	requireVisible(t, m, "Delete Instance")
	requireVisible(t, m, "Delete instance debian?")
	m, _ = pressKey(t, m, "q")
	requireNoCalls(t, svc.deleteCalls)

	m, _ = pressKey(t, m, "d")
	m, msg := pressKeyAndRun(t, m, "y")
	if len(svc.deleteCalls) != 1 || svc.deleteCalls[0] != "debian" {
		t.Fatalf("expected delete call for debian, got %v", svc.deleteCalls)
	}
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "Refreshing")
}

func TestLifecycleFailureShowsError(t *testing.T) {
	svc := newFakeService()
	svc.stopErr = errors.New("stop failed")
	svc.instances = []appincus.InstanceRow{running("ubuntu")}
	m := newHarness(t, svc, 120, 25)

	m, _ = pressKey(t, m, "s")
	m, msg := pressKeyAndRun(t, m, "y")
	m = applyMessage(t, m, msg)
	requireVisible(t, m, "stop failed")
}
