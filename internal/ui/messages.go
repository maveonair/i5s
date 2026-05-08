package ui

import (
	"time"

	"github.com/lxc/incus/v7/shared/api"
	appincus "github.com/maveonair/i5s/internal/incus"
)

type tickMsg time.Time
type instancesLoadedMsg struct{ instances []appincus.InstanceRow }
type remotesLoadedMsg struct{ remotes []appincus.RemoteRow }
type projectsLoadedMsg struct{ projects []appincus.ProjectRow }

type stateLoadedMsg struct {
	name  string
	state *api.InstanceState
}

type logFilesLoadedMsg struct{ files []string }

type logLoadedMsg struct {
	content string
	err     error
}

type opDoneMsg struct {
	message string
	err     error
	refresh bool
}

type shellDoneMsg struct{ err error }
type editDoneMsg struct{ err error }

type switchDoneMsg struct {
	message string
	err     error
}
