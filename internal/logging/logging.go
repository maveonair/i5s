package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Setup configures debug logging and returns the opened log file when enabled.
func Setup(debug bool) (*os.File, error) {
	if !debug {
		log.SetOutput(io.Discard)
		return nil, nil
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(cacheDir, "i5s")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(filepath.Join(dir, "i5s.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	return f, nil
}
