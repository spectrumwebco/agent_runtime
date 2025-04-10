package cpp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Interpreter struct {
	mutex       sync.Mutex
	tempDir     string
	includeDirs []string
	libraries   []string
	compiler    string
	flags       []string
}

type InterpreterConfig struct {
	IncludeDirs []string
	Libraries   []string
	Compiler    string
	Flags       []string
}

func NewInterpreter(config InterpreterConfig) (*Interpreter, error) {
	tempDir, err := os.MkdirTemp("", "cpp-interpreter-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	compiler := config.Compiler
	if compiler == "" {
		compiler = "g++"
	}

	_, err = exec.LookPath(compiler)
	if err != nil {
		return nil, fmt.Errorf("C++ compiler not found: %w", err)
	}

	flags := config.Flags
	if len(flags) == 0 {
		flags = []string{"-std=c++17", "-O2"}
	}

	return &Interpreter{
		tempDir:     tempDir,
		includeDirs: config.IncludeDirs,
		libraries:   config.Libraries,
		compiler:    compiler,
		flags:       flags,
	}, nil
}

func (i *Interpreter) Close() error {
	return os.RemoveAll(i.tempDir)
}

func (i *Interpreter) Exec(ctx context.Context, code string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	sourceFile := filepath.Join(i.tempDir, fmt.Sprintf("program_%d.cpp", os.Getpid()))
	execFile := filepath.Join(i.tempDir, fmt.Sprintf("program_%d", os.Getpid()))

	if err := os.WriteFile(sourceFile, []byte(code), 0600); err != nil {
		return "", fmt.Errorf("failed to write source file: %w", err)
	}

	args := append(i.flags, "-o", execFile, sourceFile)
	
	for _, dir := range i.includeDirs {
		args = append(args, "-I"+dir)
	}
	
	for _, lib := range i.libraries {
		args = append(args, "-l"+lib)
	}

	cmd := exec.CommandContext(ctx, i.compiler, args...)
	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compilation failed: %s\n%w", compileOutput, err)
	}

	cmd = exec.CommandContext(ctx, execFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %s\n%w", output, err)
	}

	return string(output), nil
}

func (i *Interpreter) ExecWithInput(ctx context.Context, code, input string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	sourceFile := filepath.Join(i.tempDir, fmt.Sprintf("program_%d.cpp", os.Getpid()))
	execFile := filepath.Join(i.tempDir, fmt.Sprintf("program_%d", os.Getpid()))

	if err := os.WriteFile(sourceFile, []byte(code), 0600); err != nil {
		return "", fmt.Errorf("failed to write source file: %w", err)
	}

	args := append(i.flags, "-o", execFile, sourceFile)
	
	for _, dir := range i.includeDirs {
		args = append(args, "-I"+dir)
	}
	
	for _, lib := range i.libraries {
		args = append(args, "-l"+lib)
	}

	cmd := exec.CommandContext(ctx, i.compiler, args...)
	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compilation failed: %s\n%w", compileOutput, err)
	}

	cmd = exec.CommandContext(ctx, execFile)
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %s\n%w", output, err)
	}

	return string(output), nil
}

func (i *Interpreter) CompileLibrary(ctx context.Context, code, libraryName string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	sourceFile := filepath.Join(i.tempDir, fmt.Sprintf("%s.cpp", libraryName))
	libFile := filepath.Join(i.tempDir, fmt.Sprintf("lib%s.so", libraryName))

	if err := os.WriteFile(sourceFile, []byte(code), 0600); err != nil {
		return "", fmt.Errorf("failed to write source file: %w", err)
	}

	args := append(i.flags, "-shared", "-fPIC", "-o", libFile, sourceFile)
	
	for _, dir := range i.includeDirs {
		args = append(args, "-I"+dir)
	}
	
	for _, lib := range i.libraries {
		args = append(args, "-l"+lib)
	}

	cmd := exec.CommandContext(ctx, i.compiler, args...)
	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("library compilation failed: %s\n%w", compileOutput, err)
	}

	return libFile, nil
}

func (i *Interpreter) AddIncludeDir(dir string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	i.includeDirs = append(i.includeDirs, dir)
}

func (i *Interpreter) AddLibrary(lib string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	i.libraries = append(i.libraries, lib)
}

func (i *Interpreter) SetCompiler(compiler string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	_, err := exec.LookPath(compiler)
	if err != nil {
		return fmt.Errorf("C++ compiler not found: %w", err)
	}
	
	i.compiler = compiler
	return nil
}

func (i *Interpreter) SetFlags(flags []string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	i.flags = flags
}
