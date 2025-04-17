package exec

import (
	"context"
	"os/exec"
)

type Execer interface {
	CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd
}

type defaultExecer struct{}

func (d defaultExecer) CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}

var DefaultExecer Execer = defaultExecer{}
