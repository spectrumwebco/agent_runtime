package python

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Interpreter struct {
	pythonPath string
	env        []string
	mutex      sync.Mutex
}

func NewInterpreter() (*Interpreter, error) {
	pythonPath, err := exec.LookPath("python3")
	if err != nil {
		return nil, fmt.Errorf("failed to find Python interpreter: %w", err)
	}
	
	interpreter := &Interpreter{
		pythonPath: pythonPath,
		env:        os.Environ(),
	}
	
	return interpreter, nil
}

func (i *Interpreter) WithPath(path string) *Interpreter {
	i.pythonPath = path
	return i
}

func (i *Interpreter) WithEnv(env []string) *Interpreter {
	i.env = env
	return i
}

func (i *Interpreter) Exec(code string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	tmpFile, err := os.CreateTemp("", "python-*.py")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(code); err != nil {
		return "", fmt.Errorf("failed to write code to temporary file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}
	
	cmd := exec.Command(i.pythonPath, tmpFile.Name())
	cmd.Env = i.env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute Python code: %w\nOutput: %s", err, output)
	}
	
	return string(output), nil
}

func (i *Interpreter) ExecFile(path string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", path)
	}
	
	cmd := exec.Command(i.pythonPath, path)
	cmd.Env = i.env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute Python file: %w\nOutput: %s", err, output)
	}
	
	return string(output), nil
}

func (i *Interpreter) ExecFunction(module, function string, args ...interface{}) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	var code strings.Builder
	code.WriteString(fmt.Sprintf("import %s\n", module))
	code.WriteString("import json\n")
	code.WriteString("import sys\n\n")
	
	code.WriteString("result = ")
	code.WriteString(module)
	code.WriteString(".")
	code.WriteString(function)
	code.WriteString("(")
	
	for j, arg := range args {
		if j > 0 {
			code.WriteString(", ")
		}
		
		switch v := arg.(type) {
		case string:
			code.WriteString(fmt.Sprintf("%q", v))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			code.WriteString(fmt.Sprintf("%v", v))
		default:
			code.WriteString(fmt.Sprintf("json.loads(%q)", fmt.Sprintf("%v", v)))
		}
	}
	
	code.WriteString(")\n")
	code.WriteString("print(json.dumps(result))")
	
	tmpFile, err := os.CreateTemp("", "python-*.py")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(code.String()); err != nil {
		return "", fmt.Errorf("failed to write code to temporary file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}
	
	cmd := exec.Command(i.pythonPath, tmpFile.Name())
	cmd.Env = i.env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute Python function: %w\nOutput: %s", err, output)
	}
	
	return string(output), nil
}

func (i *Interpreter) InstallPackage(pkg string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	
	cmd := exec.Command(i.pythonPath, "-m", "pip", "install", pkg)
	cmd.Env = i.env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package: %w\nOutput: %s", err, output)
	}
	
	return nil
}

func (i *Interpreter) Close() error {
	return nil
}
