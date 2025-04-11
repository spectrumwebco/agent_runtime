package cpp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type InterpreterConfig struct {
	CompilerPath string   // Path to the C++ compiler (e.g., g++)
	Flags        []string // Compiler flags (e.g., -std=c++17)
	Timeout      time.Duration // Execution timeout
}

type Interpreter struct {
	config InterpreterConfig
}

func NewInterpreter(config InterpreterConfig) (*Interpreter, error) {
	if config.CompilerPath == "" {
		path, err := exec.LookPath("g++") // Default to g++
		if err != nil {
			path, err = exec.LookPath("clang++") // Fallback to clang++
			if err != nil {
				return nil, fmt.Errorf("c++ compiler (g++ or clang++) not found in PATH: %w", err)
			}
		}
		config.CompilerPath = path
	}

	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second // Default 60-second timeout
	}

	fmt.Printf("C++ FFI Interpreter initialized using compiler: %s with flags: %v\n", config.CompilerPath, config.Flags)
	return &Interpreter{config: config}, nil
}

func (i *Interpreter) Exec(ctx context.Context, code string) (string, error) {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), i.config.Timeout)
		defer cancel()
	} else {
		if _, deadlineSet := ctx.Deadline(); !deadlineSet {
			ctx, cancel = context.WithTimeout(ctx, i.config.Timeout)
			defer cancel()
		}
	}


	tempDir, err := os.MkdirTemp("", "cpp_ffi_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	sourceFilePath := filepath.Join(tempDir, "main.cpp")
	err = os.WriteFile(sourceFilePath, []byte(code), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write temporary C++ source file: %w", err)
	}

	executablePath := filepath.Join(tempDir, "main")

	compileArgs := append(i.config.Flags, sourceFilePath, "-o", executablePath)
	compileCmd := exec.CommandContext(ctx, i.config.CompilerPath, compileArgs...)
	compileCmd.Dir = tempDir // Run compiler in the temp directory

	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	fmt.Printf("Compiling C++ code: %s %v\n", i.config.CompilerPath, compileArgs)
	err = compileCmd.Run()
	compileStderrStr := compileStderr.String()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("c++ compilation timed out: %w", ctx.Err())
	}
	if err != nil {
		return "", fmt.Errorf("c++ compilation failed: %w\nCompiler Stderr: %s", err, compileStderrStr)
	}
	if compileStderrStr != "" {
		fmt.Printf("C++ compilation produced stderr (warnings):\n%s\n", compileStderrStr)
	}

	execCmd := exec.CommandContext(ctx, executablePath)
	execCmd.Dir = tempDir // Run executable in the temp directory

	var execStdout, execStderr bytes.Buffer
	execCmd.Stdout = &execStdout
	execCmd.Stderr = &execStderr

	fmt.Printf("Executing compiled C++ code: %s\n", executablePath)
	err = execCmd.Run()

	execStdoutStr := execStdout.String()
	execStderrStr := execStderr.String()

	if ctx.Err() == context.DeadlineExceeded {
		return execStdoutStr, fmt.Errorf("c++ execution timed out: %w\nStderr: %s", ctx.Err(), execStderrStr)
	}
	if err != nil {
		return execStdoutStr, fmt.Errorf("c++ execution failed: %w\nStderr: %s", err, execStderrStr)
	}

	if execStderrStr != "" {
		fmt.Printf("C++ execution produced stderr (even on success):\n%s\n", execStderrStr)
		return execStdoutStr + "\n--- STDERR ---\n" + execStderrStr, nil
	}

	return execStdoutStr, nil
}

func (i *Interpreter) Close() {
	fmt.Println("C++ FFI Interpreter closed.")
}
