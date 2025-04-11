package python

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ScriptRunner struct {
	interpreter *Interpreter
	scriptDir   string
}

func NewScriptRunner(scriptDir string) (*ScriptRunner, error) {
	interpreter, err := NewInterpreter()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Python interpreter: %w", err)
	}

	return &ScriptRunner{
		interpreter: interpreter,
		scriptDir:   scriptDir,
	}, nil
}

func (r *ScriptRunner) RunToolScript(toolName string, args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool arguments: %w", err)
	}

	scriptPath := filepath.Join(r.scriptDir, "tools", toolName, "run.py")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tool script not found: %s", scriptPath)
	}

	output, err := r.interpreter.ExecutePythonScript(scriptPath, []string{string(argsJSON)})
	if err != nil {
		return nil, fmt.Errorf("failed to run tool script: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse tool output: %w", err)
	}

	return result, nil
}

func (r *ScriptRunner) Close() error {
	return r.interpreter.Close()
}
