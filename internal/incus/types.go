package incus

import (
	"fmt"
	"strings"
	"time"

	"github.com/lxc/incus/v7/shared/api"
)

// InstanceRow is the instance data rendered and searched by the TUI.
type InstanceRow struct {
	Name        string
	Type        string
	Status      string
	StatusCode  api.StatusCode
	IPv4        string
	IPv6        string
	Image       string
	Description string
	Profiles    []string
	Location    string
	Snapshots   int
	CreatedAt   time.Time
	Raw         api.Instance
}

// IsRunning reports whether the instance is currently running.
func (r InstanceRow) IsRunning() bool {
	return r.StatusCode == api.Running
}

// IsStopped reports whether the instance is currently stopped.
func (r InstanceRow) IsStopped() bool {
	return r.StatusCode == api.Stopped
}

// SearchText returns normalized text used by the instance filter.
func (r InstanceRow) SearchText() string {
	return strings.ToLower(strings.Join([]string{
		r.Name,
		r.Type,
		r.Status,
		r.IPv4,
		r.IPv6,
		fmt.Sprintf("%d", r.Snapshots),
		r.Image,
		r.Description,
		strings.Join(r.Profiles, " "),
		r.Location,
	}, " "))
}

// RemoteRow is the remote data rendered by the remote picker.
type RemoteRow struct {
	Name           string
	Addr           string
	Protocol       string
	DefaultProject string
	Static         bool
	Public         bool
}

// ProjectRow is the project data rendered by the project picker.
type ProjectRow struct {
	Name        string
	Description string
	UsedBy      int
}
