package main

import (
	"fmt"
	"os"

	"github.com/maveonair/i5s/internal/app"
	"github.com/maveonair/i5s/internal/config"
	"github.com/maveonair/i5s/internal/logging"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	cfg := config.Parse()
	if cfg.ShowVersion {
		fmt.Printf("Version: %s\nCommit: %s\n", version, commit)
		return
	}

	logFile, err := logging.Setup(cfg.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logging: %v\n", err)
		os.Exit(1)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	if err := app.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
