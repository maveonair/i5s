package ui

import (
	"context"
	"io"

	appincus "github.com/maveonair/i5s/internal/incus"
)

type shellCommand struct {
	service appincus.Service
	name    string
}

func (c shellCommand) Run() error          { return c.service.ExecShell(context.Background(), c.name) }
func (c shellCommand) SetStdin(io.Reader)  {}
func (c shellCommand) SetStdout(io.Writer) {}
func (c shellCommand) SetStderr(io.Writer) {}

type editConfigCommand struct {
	service appincus.Service
	name    string
}

func (c editConfigCommand) Run() error {
	return c.service.EditInstanceConfig(context.Background(), c.name)
}
func (c editConfigCommand) SetStdin(io.Reader)  {}
func (c editConfigCommand) SetStdout(io.Writer) {}
func (c editConfigCommand) SetStderr(io.Writer) {}
