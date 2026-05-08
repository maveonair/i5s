package incus

import (
	"context"
	"io"
	"sort"

	client "github.com/lxc/incus/v7/client"
)

// ListLogs returns available log file names for one instance.
func (s *IncusService) ListLogs(ctx context.Context, name string) ([]string, error) {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return nil, err
	}

	logs, err := server.GetInstanceLogfiles(name)
	if err != nil {
		return nil, err
	}
	sort.Strings(logs)
	return logs, nil
}

// GetLog returns one instance log file as text.
func (s *IncusService) GetLog(ctx context.Context, name string, logName string) (string, error) {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return "", err
	}

	rc, err := server.GetInstanceLogfile(name, logName)
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetConsoleLog returns the instance console log, falling back to console.log when needed.
func (s *IncusService) GetConsoleLog(ctx context.Context, name string) (string, error) {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return "", err
	}

	rc, err := server.GetInstanceConsoleLog(name, &client.InstanceConsoleLogArgs{})
	if err != nil {
		return s.GetLog(ctx, name, "console.log")
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ClearConsoleLog clears the instance console log.
func (s *IncusService) ClearConsoleLog(ctx context.Context, name string) error {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return err
	}

	return server.DeleteInstanceConsoleLog(name, &client.InstanceConsoleLogArgs{})
}
