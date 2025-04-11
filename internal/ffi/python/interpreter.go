package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Interpreter struct {
	pythonPath string
	venvPath   string
}

func NewInterpreter() (*Interpreter, error) {
	pythonPaths := []string{
		"python3.11",
		"python3.10",
		"python3.9",
		"python3",
		"/usr/bin/python3",
		"/usr/local/bin/python3",
	}

	var pythonPath string
	for _, path := range pythonPaths {
		if _, err := exec.LookPath(path); err == nil {
			pythonPath = path
			break
		}
	}

	if pythonPath == "" {
		return nil, fmt.Errorf("no suitable Python interpreter found")
	}

	venvPath := os.Getenv("VIRTUAL_ENV")
	if venvPath == "" {
		projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(os.Args[0]))))
		venvPath = filepath.Join(projectRoot, ".venv")
	}

	return &Interpreter{
		pythonPath: pythonPath,
		venvPath:   venvPath,
	}, nil
}

func (i *Interpreter) ExecutePythonScript(scriptPath string, args []string) (string, error) {
	if !filepath.IsAbs(scriptPath) {
		return "", fmt.Errorf("script path must be absolute")
	}

	var cmd *exec.Cmd
	if i.venvPath != "" {
		pythonPath := filepath.Join(i.venvPath, "bin", "python")
		if _, err := os.Stat(pythonPath); err == nil {
			cmd = exec.Command(pythonPath, append([]string{scriptPath}, args...)...)
		}
	}

	if cmd == nil {
		cmd = exec.Command(i.pythonPath, append([]string{scriptPath}, args...)...)
	}

	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(scriptPath))))
	pythonPath := os.Getenv("PYTHONPATH")
	if pythonPath == "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHONPATH=%s", projectRoot))
	} else {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHONPATH=%s:%s", projectRoot, pythonPath))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute Python script: %v\nOutput: %s", err, output)
	}

	return strings.TrimSpace(string(output)), nil
}

func (i *Interpreter) Close() error {
	return nil
}
