package tools

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/python"
)

type PythonTool struct {
	name        string
	description string
	interpreter *python.Interpreter
}

func NewPythonTool(name, description string) (*PythonTool, error) {
	interpreter, err := python.NewInterpreter()
	if err != nil {
		return nil, fmt.Errorf("failed to create Python interpreter: %w", err)
	}
	
	return &PythonTool{
		name:        name,
		description: description,
		interpreter: interpreter,
	}, nil
}

func (t *PythonTool) Name() string {
	return t.name
}

func (t *PythonTool) Description() string {
	return t.description
}

func (t *PythonTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	if file, ok := params["file"].(string); ok && file != "" {
		return t.interpreter.ExecFile(file)
	}
	
	if function, ok := params["function"].(string); ok && function != "" {
		module, _ := params["module"].(string)
		if module == "" {
			return nil, fmt.Errorf("module parameter is required for function execution")
		}
		
		var args []interface{}
		if argsParam, ok := params["args"].([]interface{}); ok {
			args = argsParam
		}
		
		return t.interpreter.ExecFunction(module, function, args...)
	}
	
	code, ok := params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code parameter is required")
	}
	
	return t.interpreter.Exec(code)
}

func (t *PythonTool) Close() error {
	return t.interpreter.Close()
}
