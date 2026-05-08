package ui

import (
	"context"

	"github.com/lxc/incus/v7/shared/api"
	appincus "github.com/maveonair/i5s/internal/incus"
)

type fakeService struct {
	remote  string
	project string

	instances []appincus.InstanceRow
	remotes   []appincus.RemoteRow
	projects  []appincus.ProjectRow

	logFiles   map[string][]string
	logs       map[string]string
	consoleLog map[string]string

	listInstancesCalls int
	getLogCalls        []string
	getConsoleCalls    []string

	listInstancesErr error
	startErr         error
	stopErr          error
	deleteErr        error
	shellErr         error
	editErr          error
	switchRemoteErr  error
	switchProjectErr error

	startCalls         []string
	stopCalls          []string
	deleteCalls        []string
	shellCalls         []string
	editCalls          []string
	switchRemoteCalls  []string
	switchProjectCalls []string
	clearConsoleCalls  []string
}

func newFakeService() *fakeService {
	return &fakeService{
		remote:     "local",
		project:    "default",
		logFiles:   map[string][]string{},
		logs:       map[string]string{},
		consoleLog: map[string]string{},
	}
}

func (f *fakeService) CurrentRemote() string  { return f.remote }
func (f *fakeService) CurrentProject() string { return f.project }

func (f *fakeService) ListRemotes() ([]appincus.RemoteRow, error) {
	return f.remotes, nil
}

func (f *fakeService) SwitchRemote(_ context.Context, name string) error {
	f.switchRemoteCalls = append(f.switchRemoteCalls, name)
	if f.switchRemoteErr != nil {
		return f.switchRemoteErr
	}
	f.remote = name
	return nil
}

func (f *fakeService) ListProjects(context.Context) ([]appincus.ProjectRow, error) {
	return f.projects, nil
}

func (f *fakeService) SwitchProject(_ context.Context, name string) error {
	f.switchProjectCalls = append(f.switchProjectCalls, name)
	if f.switchProjectErr != nil {
		return f.switchProjectErr
	}
	f.project = name
	return nil
}

func (f *fakeService) ListInstances(context.Context) ([]appincus.InstanceRow, error) {
	f.listInstancesCalls++
	if f.listInstancesErr != nil {
		return nil, f.listInstancesErr
	}
	return f.instances, nil
}

func (f *fakeService) GetInstanceState(context.Context, string) (*api.InstanceState, error) {
	return &api.InstanceState{}, nil
}

func (f *fakeService) ListLogs(_ context.Context, name string) ([]string, error) {
	return f.logFiles[name], nil
}

func (f *fakeService) GetLog(_ context.Context, name string, logName string) (string, error) {
	f.getLogCalls = append(f.getLogCalls, name+"/"+logName)
	return f.logs[name+"/"+logName], nil
}

func (f *fakeService) GetConsoleLog(_ context.Context, name string) (string, error) {
	f.getConsoleCalls = append(f.getConsoleCalls, name)
	return f.consoleLog[name], nil
}

func (f *fakeService) ClearConsoleLog(_ context.Context, name string) error {
	f.clearConsoleCalls = append(f.clearConsoleCalls, name)
	return nil
}

func (f *fakeService) StartInstance(_ context.Context, name string) error {
	f.startCalls = append(f.startCalls, name)
	return f.startErr
}

func (f *fakeService) StopInstance(_ context.Context, name string) error {
	f.stopCalls = append(f.stopCalls, name)
	return f.stopErr
}

func (f *fakeService) DeleteInstance(_ context.Context, name string) error {
	f.deleteCalls = append(f.deleteCalls, name)
	return f.deleteErr
}

func (f *fakeService) ExecShell(_ context.Context, name string) error {
	f.shellCalls = append(f.shellCalls, name)
	return f.shellErr
}

func (f *fakeService) EditInstanceConfig(_ context.Context, name string) error {
	f.editCalls = append(f.editCalls, name)
	return f.editErr
}

func running(name string) appincus.InstanceRow {
	return appincus.InstanceRow{Name: name, Type: "container", Status: "RUNNING", StatusCode: api.Running, IPv4: "10.0.3.15", IPv6: "fd42::1", Snapshots: 2}
}

func stopped(name string) appincus.InstanceRow {
	return appincus.InstanceRow{Name: name, Type: "virtual-machine", Status: "STOPPED", StatusCode: api.Stopped, IPv4: "-", IPv6: "-"}
}
