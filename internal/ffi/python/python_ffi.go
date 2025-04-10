package python

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Interpreter struct {
	pythonPath string // Path to the python executable
}

func NewInterpreter() (*Interpreter, error) {
	path, err := exec.LookPath("python3")
	if err != nil {
		path, err = exec.LookPath("python")
		if err != nil {
			return nil, fmt.Errorf("python executable not found in PATH: %w", err)
		}
	}
	fmt.Printf("Python FFI Interpreter initialized using executable: %s\n", path)
	return &Interpreter{pythonPath: path}, nil
}

func (i *Interpreter) Exec(ctx context.Context, code string) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) // 30-second timeout
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, i.pythonPath, "-c", code)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Printf("Executing Python code via FFI: %s\n", code)
	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("python execution timed out: %w", ctx.Err())
	}
	if err != nil {
		return stdoutStr, fmt.Errorf("python execution failed: %w\nStderr: %s", err, stderrStr)
	}

	if stderrStr != "" {
		fmt.Printf("Python execution produced stderr (even on success):\n%s\n", stderrStr)
		return stdoutStr + "\n--- STDERR ---\n" + stderrStr, nil
	}

	return stdoutStr, nil
}

func (i *Interpreter) Close() {
	fmt.Println("Python FFI Interpreter closed.")
}
