//go:build windows

package incus

import (
	"context"
	"errors"
)

// ExecShell reports that interactive shell exec is unsupported on Windows.
func (s *IncusService) ExecShell(_ context.Context, _ string) error {
	return errors.New("native shell exec is not implemented on Windows yet")
}
