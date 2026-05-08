//go:build !windows

package incus

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/websocket"
	client "github.com/lxc/incus/v7/client"
	"github.com/lxc/incus/v7/shared/api"
	"github.com/lxc/incus/v7/shared/termios"
	"golang.org/x/sys/unix"
)

// ExecShell opens an interactive login shell inside the instance.
func (s *IncusService) ExecShell(ctx context.Context, name string) error {
	server, err := s.instanceServer(ctx)
	if err != nil {
		return err
	}

	stdinFD := int(os.Stdin.Fd())
	stdoutFD := int(os.Stdout.Fd())
	stdinTerminal := termios.IsTerminal(stdinFD)
	stdoutTerminal := termios.IsTerminal(stdoutFD)

	var oldState *termios.State
	if stdinTerminal {
		oldState, err = termios.MakeRaw(stdinFD)
		if err != nil {
			return err
		}
		defer func() { _ = termios.Restore(stdinFD, oldState) }()
	}

	width, height := 0, 0
	if stdoutTerminal {
		width, height, err = termios.GetSize(stdoutFD)
		if err != nil {
			return err
		}
	}

	env := map[string]string{}
	if term, ok := os.LookupEnv("TERM"); ok {
		env["TERM"] = term
	}

	req := api.InstanceExecPost{
		Command:     []string{"su", "-l"},
		WaitForWS:   true,
		Interactive: true,
		Environment: env,
		Width:       width,
		Height:      height,
	}

	done := make(chan struct{})
	defer close(done)

	args := client.InstanceExecArgs{
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		DataDone: make(chan bool),
		Control: func(control *websocket.Conn) {
			controlSocketHandler(control, true, done)
		},
	}

	op, err := server.ExecInstance(name, req, &args)
	if err != nil {
		return err
	}

	if err := op.WaitContext(ctx); err != nil {
		return err
	}

	select {
	case <-args.DataDone:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func controlSocketHandler(control *websocket.Conn, interactive bool, done <-chan struct{}) {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch,
		unix.SIGWINCH,
		unix.SIGTERM,
		unix.SIGHUP,
		unix.SIGINT,
		unix.SIGQUIT,
		unix.SIGABRT,
		unix.SIGTSTP,
		unix.SIGTTIN,
		unix.SIGTTOU,
		unix.SIGUSR1,
		unix.SIGUSR2,
		unix.SIGSEGV,
		unix.SIGCONT,
	)
	defer signal.Stop(ch)
	defer func() {
		_ = control.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}()

	for {
		select {
		case <-done:
			return
		case sig := <-ch:
			switch sig {
			case unix.SIGWINCH:
				if interactive {
					_ = sendTermSize(control)
				}
			case unix.SIGTERM:
				_ = forwardSignal(control, unix.SIGTERM)
			case unix.SIGHUP:
				_ = forwardSignal(control, unix.SIGHUP)
			case unix.SIGINT:
				_ = forwardSignal(control, unix.SIGINT)
			case unix.SIGQUIT:
				_ = forwardSignal(control, unix.SIGQUIT)
			case unix.SIGABRT:
				_ = forwardSignal(control, unix.SIGABRT)
			case unix.SIGTSTP:
				_ = forwardSignal(control, unix.SIGTSTP)
			case unix.SIGTTIN:
				_ = forwardSignal(control, unix.SIGTTIN)
			case unix.SIGTTOU:
				_ = forwardSignal(control, unix.SIGTTOU)
			case unix.SIGUSR1:
				_ = forwardSignal(control, unix.SIGUSR1)
			case unix.SIGUSR2:
				_ = forwardSignal(control, unix.SIGUSR2)
			case unix.SIGSEGV:
				_ = forwardSignal(control, unix.SIGSEGV)
			case unix.SIGCONT:
				_ = forwardSignal(control, unix.SIGCONT)
			}
		}
	}
}

func sendTermSize(control *websocket.Conn) error {
	width, height, err := termios.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}

	return control.WriteJSON(api.InstanceExecControl{
		Command: "window-resize",
		Args: map[string]string{
			"width":  strconv.Itoa(width),
			"height": strconv.Itoa(height),
		},
	})
}

func forwardSignal(control *websocket.Conn, sig unix.Signal) error {
	return control.WriteJSON(api.InstanceExecControl{Command: "signal", Signal: int(sig)})
}
