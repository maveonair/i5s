package incus

import (
	"context"
	"fmt"
	"sort"
	"sync"

	client "github.com/lxc/incus/v7/client"
	"github.com/lxc/incus/v7/shared/api"
	"github.com/lxc/incus/v7/shared/cliconfig"
)

const instanceEnrichmentConcurrency = 8

// Service defines the Incus operations used by the TUI.
type Service interface {
	// CurrentRemote returns the active runtime remote.
	CurrentRemote() string
	// CurrentProject returns the active runtime project.
	CurrentProject() string
	// ListRemotes returns configured remotes.
	ListRemotes() ([]RemoteRow, error)
	// SwitchRemote changes the runtime remote.
	SwitchRemote(ctx context.Context, name string) error
	// ListProjects returns projects for the active remote.
	ListProjects(ctx context.Context) ([]ProjectRow, error)
	// SwitchProject changes the runtime project.
	SwitchProject(ctx context.Context, name string) error
	// ListInstances returns instances for the active remote/project.
	ListInstances(ctx context.Context) ([]InstanceRow, error)
	// GetInstanceState returns runtime state for one instance.
	GetInstanceState(ctx context.Context, name string) (*api.InstanceState, error)
	// ListLogs returns log names for one instance.
	ListLogs(ctx context.Context, name string) ([]string, error)
	// GetLog returns one named instance log.
	GetLog(ctx context.Context, name string, logName string) (string, error)
	// GetConsoleLog returns the console log for one instance.
	GetConsoleLog(ctx context.Context, name string) (string, error)
	// ClearConsoleLog clears the console log for one instance.
	ClearConsoleLog(ctx context.Context, name string) error
	// StartInstance starts one instance.
	StartInstance(ctx context.Context, name string) error
	// StopInstance stops one instance.
	StopInstance(ctx context.Context, name string) error
	// DeleteInstance deletes one instance.
	DeleteInstance(ctx context.Context, name string) error
	// ExecShell opens an interactive shell in one instance.
	ExecShell(ctx context.Context, name string) error
	// EditInstanceConfig edits one instance config in the user's editor.
	EditInstanceConfig(ctx context.Context, name string) error
}

// IncusService connects i5s to one runtime Incus remote and project.
type IncusService struct {
	mu      sync.Mutex
	config  *cliconfig.Config
	remote  string
	project string
	server  client.InstanceServer
}

// New creates an IncusService using Incus CLI configuration and optional runtime overrides.
func New(remoteOverride, projectOverride string) (*IncusService, error) {
	conf, err := cliconfig.LoadConfig("")
	if err != nil {
		return nil, err
	}

	remote := remoteOverride
	if remote == "" {
		remote = conf.DefaultRemote
	}
	if remote == "" {
		remote = "local"
	}

	project := projectOverride
	if project == "" {
		if r, ok := conf.Remotes[remote]; ok && r.Project != "" {
			project = r.Project
		}
	}
	if project == "" {
		project = api.ProjectDefaultName
	}

	return &IncusService{config: conf, remote: remote, project: project}, nil
}

// CurrentRemote returns the remote selected for this i5s session.
func (s *IncusService) CurrentRemote() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remote
}

// CurrentProject returns the project selected for this i5s session.
func (s *IncusService) CurrentProject() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.project
}

func (s *IncusService) instanceServer(ctx context.Context) (client.InstanceServer, error) {
	s.mu.Lock()
	if s.server != nil {
		server := s.server
		s.mu.Unlock()
		return serverWithContext(server, ctx), nil
	}
	remote := s.remote
	project := s.project
	s.mu.Unlock()

	server, err := s.buildServer(ctx, remote, project)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	if s.remote == remote && s.project == project && s.server == nil {
		s.server = server
	}
	s.mu.Unlock()

	return serverWithContext(server, ctx), nil
}

func serverWithContext(server client.InstanceServer, ctx context.Context) client.InstanceServer {
	if ctx == nil {
		return server
	}

	// The concrete Incus protocol supports context-bound requests, but the
	// versioned InstanceServer interface does not expose that method.
	if contextual, ok := server.(interface {
		WithContext(context.Context) client.InstanceServer
	}); ok {
		return contextual.WithContext(ctx)
	}

	return server
}

func (s *IncusService) connect(remote, project string) error {
	server, err := s.buildServer(context.Background(), remote, project)
	if err != nil {
		return err
	}

	s.remote = remote
	s.project = project
	s.server = server
	return nil
}

func (s *IncusService) buildServer(ctx context.Context, remote, project string) (client.InstanceServer, error) {
	type result struct {
		server client.InstanceServer
		err    error
	}

	ch := make(chan result, 1)
	go func() {
		server, err := s.config.GetInstanceServer(remote)
		ch <- result{server: server, err: err}
	}()

	var res result
	if ctx == nil {
		res = <-ch
	} else {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res = <-ch:
		}
	}

	if res.err != nil {
		return nil, res.err
	}

	if project != "" {
		res.server = res.server.UseProject(project)
	}

	return res.server, nil
}

// ListRemotes returns remotes from the loaded Incus CLI configuration.
func (s *IncusService) ListRemotes() ([]RemoteRow, error) {
	rows := make([]RemoteRow, 0, len(s.config.Remotes))
	for name, remote := range s.config.Remotes {
		project := remote.Project
		if project == "" {
			project = api.ProjectDefaultName
		}

		rows = append(rows, RemoteRow{
			Name:           name,
			Addr:           remote.Addr,
			Protocol:       remote.Protocol,
			DefaultProject: project,
			Static:         remote.Static,
			Public:         remote.Public,
		})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows, nil
}

// SwitchRemote changes the runtime remote and resets the project for this session.
func (s *IncusService) SwitchRemote(_ context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	remote, ok := s.config.Remotes[name]
	if !ok {
		return fmt.Errorf("remote %q not found", name)
	}

	project := remote.Project
	if project == "" {
		project = api.ProjectDefaultName
	}

	s.remote = name
	s.project = project
	s.server = nil
	return nil
}

// ListProjects returns projects available on the current runtime remote.
func (s *IncusService) ListProjects(ctx context.Context) ([]ProjectRow, error) {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return nil, err
	}

	projects, err := server.GetProjects()
	if err != nil {
		return nil, err
	}

	rows := make([]ProjectRow, 0, len(projects))
	for _, project := range projects {
		rows = append(rows, ProjectRow{Name: project.Name, Description: project.Description, UsedBy: len(project.UsedBy)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows, nil
}

// SwitchProject changes the runtime project for this session.
func (s *IncusService) SwitchProject(_ context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.project = name
	if s.server != nil {
		s.server = s.server.UseProject(name)
	}
	return nil
}
