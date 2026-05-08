package app

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maveonair/i5s/internal/config"
	appincus "github.com/maveonair/i5s/internal/incus"
	"github.com/maveonair/i5s/internal/ui"
)

// Run starts the fullscreen Bubble Tea program.
func Run(cfg config.Runtime) error {
	service, err := appincus.New(cfg.Remote, cfg.Project)
	if err != nil {
		return fmt.Errorf("failed to connect to Incus: %w", err)
	}

	log.Printf("starting i5s remote=%s project=%s", service.CurrentRemote(), service.CurrentProject())
	program := tea.NewProgram(ui.New(service, cfg.RefreshEvery), tea.WithAltScreen())
	_, err = program.Run()
	return err
}
