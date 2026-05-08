package config

import (
	"flag"
	"time"
)

// Runtime contains command-line configuration for one i5s session.
type Runtime struct {
	Remote       string
	Project      string
	RefreshEvery time.Duration
	Debug        bool
	ShowVersion  bool
}

// Parse reads command-line flags into a Runtime configuration.
func Parse() Runtime {
	cfg := Runtime{RefreshEvery: 5 * time.Second}
	flag.StringVar(&cfg.Remote, "remote", "", "Incus remote to use for this session")
	flag.StringVar(&cfg.Project, "project", "", "Incus project to use for this session")
	flag.DurationVar(&cfg.RefreshEvery, "refresh", cfg.RefreshEvery, "Auto-refresh interval")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Print version and exit")
	flag.Parse()
	return cfg
}
