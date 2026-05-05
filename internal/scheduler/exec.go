package scheduler

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

const defaultTimeout = 5 * time.Minute

// execCommand runs a shell command string with a default timeout.
// It splits the command on whitespace and executes it directly.
func execCommand(command string) error {
	return execCommandWithTimeout(command, defaultTimeout)
}

// execCommandWithTimeout runs a shell command with an explicit timeout.
func execCommandWithTimeout(command string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	parts := splitCommand(command)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return &ExecError{Cause: err, Output: string(output)}
		}
		return err
	}
	return nil
}

// ExecError wraps a command execution error with its combined output.
type ExecError struct {
	Cause  error
	Output string
}

func (e *ExecError) Error() string {
	return e.Cause.Error() + ": " + strings.TrimSpace(e.Output)
}

func (e *ExecError) Unwrap() error { return e.Cause }

// splitCommand splits a command string into argv, respecting quoted strings.
func splitCommand(cmd string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false

	for _, ch := range cmd {
		switch {
		case ch == '"':
			inQuote = !inQuote
		case ch == ' ' && !inQuote:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
