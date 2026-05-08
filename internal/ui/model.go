package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lxc/incus/v7/shared/api"
	appincus "github.com/maveonair/i5s/internal/incus"
)

const lifecycleTimeout = 2 * time.Minute

// ViewMode identifies the active screen or modal in the TUI.
type ViewMode int

const (
	// ViewInstances shows the main instance table.
	ViewInstances ViewMode = iota
	// ViewLogs shows an instance log file.
	ViewLogs
	// ViewConsole shows the selected instance console log.
	ViewConsole
	// ViewRemotePicker shows the runtime remote picker.
	ViewRemotePicker
	// ViewProjectPicker shows the runtime project picker.
	ViewProjectPicker
	// ViewLogPicker shows the instance log-file picker.
	ViewLogPicker
	// ViewHelp shows the help screen.
	ViewHelp
	// ViewConfirmStop asks for stop confirmation.
	ViewConfirmStop
	// ViewConfirmDelete asks for delete confirmation.
	ViewConfirmDelete
	// ViewConfirmClearConsole asks for console-log clear confirmation.
	ViewConfirmClearConsole
)

// Model is the root Bubble Tea model for the i5s TUI.
type Model struct {
	service appincus.Service
	styles  Styles

	width  int
	height int

	view ViewMode

	instances     []appincus.InstanceRow
	selectedIndex int
	selectedName  string

	remotes     []appincus.RemoteRow
	projects    []appincus.ProjectRow
	pickerIndex int

	filter    string
	filtering bool

	loading       bool
	busy          bool
	operation     string
	errorMessage  string
	statusMessage string
	lastRefresh   time.Time
	refreshEvery  time.Duration

	state        *api.InstanceState
	stateForName string

	logFiles    []string
	selectedLog string
	viewport    viewport.Model
}

// New creates a Model bound to an Incus service and refresh interval.
func New(service appincus.Service, refreshEvery time.Duration) Model {
	vp := viewport.New(0, 0)
	return Model{service: service, styles: NewStyles(), refreshEvery: refreshEvery, viewport: vp, loading: true, operation: "Loading..."}
}

// Init starts the initial load commands for the TUI.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.refreshCmd(), m.loadRemotesCmd(), m.loadProjectsCmd(), tickCmd(m.refreshEvery))
}

// Update handles Bubble Tea messages and returns the next model state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeViewport()
		return m, nil
	case tickMsg:
		if m.view == ViewInstances && !m.loading {
			m.loading = true
			m.operation = "Refreshing..."
			return m, tea.Batch(m.refreshCmd(), tickCmd(m.refreshEvery))
		}
		return m, tickCmd(m.refreshEvery)
	case instancesLoadedMsg:
		m.loading = false
		if !m.busy {
			m.operation = ""
		}
		m.errorMessage = ""
		m.instances = msg.instances
		m.lastRefresh = time.Now()
		m.restoreSelection()
		return m, m.loadSelectedStateCmd()
	case remotesLoadedMsg:
		m.remotes = msg.remotes
		return m, nil
	case projectsLoadedMsg:
		m.projects = msg.projects
		return m, nil
	case stateLoadedMsg:
		if msg.name == m.selected().Name {
			m.state = msg.state
			m.stateForName = msg.name
		}
		return m, nil
	case logFilesLoadedMsg:
		m.logFiles = msg.files
		if len(msg.files) == 0 {
			m.setStatus("No logs found")
			m.view = ViewInstances
			return m, nil
		}
		if len(msg.files) == 1 {
			m.selectedLog = msg.files[0]
			m.view = ViewLogs
			return m, m.loadLogCmd(msg.files[0])
		}
		m.pickerIndex = 0
		m.view = ViewLogPicker
		return m, nil
	case logLoadedMsg:
		if msg.err != nil {
			m.setError(msg.err.Error())
			return m, nil
		}
		m.viewport.SetContent(msg.content)
		m.viewport.GotoBottom()
		return m, nil
	case opDoneMsg:
		m.loading = false
		m.busy = false
		m.operation = ""
		if msg.err != nil {
			m.setError(msg.err.Error())
			return m, nil
		}
		m.setStatus(msg.message)
		if msg.refresh {
			m.loading = true
			m.operation = "Refreshing..."
			return m, m.refreshCmd()
		}
		return m, nil
	case shellDoneMsg:
		if msg.err != nil {
			m.setError(fmt.Sprintf("Failed to open shell: %v", msg.err))
		} else {
			m.setStatus("Shell exited")
		}
		return m, m.refreshCmd()
	case editDoneMsg:
		if msg.err != nil {
			m.setError(fmt.Sprintf("Failed to edit config: %v", msg.err))
			return m, nil
		}
		m.setStatus("Config updated")
		return m, m.refreshCmd()
	case switchDoneMsg:
		m.loading = false
		m.busy = false
		m.operation = ""
		if msg.err != nil {
			m.setError(msg.err.Error())
			return m, nil
		}
		m.setStatus(msg.message)
		m.view = ViewInstances
		return m, tea.Batch(m.loadProjectsCmd(), m.refreshCmd())
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.view == ViewLogs || m.view == ViewConsole || m.view == ViewHelp {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) setError(message string) {
	m.errorMessage = message
	m.statusMessage = ""
}

func (m *Model) setStatus(message string) {
	m.statusMessage = message
	m.errorMessage = ""
}

func (m *Model) resizeViewport() {
	m.viewport.Width = m.width
	m.viewport.Height = max(1, m.height-3)
}

func (m Model) pickerCount() int {
	switch m.view {
	case ViewRemotePicker:
		return len(m.remotes)
	case ViewProjectPicker:
		return len(m.projects)
	case ViewLogPicker:
		return len(m.logFiles)
	default:
		return 0
	}
}
