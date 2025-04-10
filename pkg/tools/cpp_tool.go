package tools

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/cpp"
)

type CppTool struct {
	name        string
	description string
	interpreter *cpp.Interpreter
}

func NewCppTool(name, description string) (*CppTool, error) {
	interpreter, err := cpp.NewInterpreter(cpp.InterpreterConfig{
		Flags: []string{"-std=c++17", "-O2"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create C++ interpreter: %w", err)
	}
	
	return &CppTool{
		name:        name,
		description: description,
		interpreter: interpreter,
	}, nil
}

func (t *CppTool) Name() string {
	return t.name
}

func (t *CppTool) Description() string {
	return t.description
}

func (t *CppTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	code, ok := params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code parameter is required")
	}
	
	if flagsParam, ok := params["flags"].([]interface{}); ok {
		flags := make([]string, 0, len(flagsParam))
		for _, flag := range flagsParam {
			if flagStr, ok := flag.(string); ok {
				flags = append(flags, flagStr)
			}
		}
		t.interpreter.SetFlags(flags)
	}
	
	if includeDirsParam, ok := params["include_dirs"].([]interface{}); ok {
		for _, dir := range includeDirsParam {
			if dirStr, ok := dir.(string); ok {
				t.interpreter.AddIncludeDir(dirStr)
			}
		}
	}
	
	if librariesParam, ok := params["libraries"].([]interface{}); ok {
		for _, lib := range librariesParam {
			if libStr, ok := lib.(string); ok {
				t.interpreter.AddLibrary(libStr)
			}
		}
	}
	
	if library, ok := params["library"].(string); ok && library != "" {
		return t.interpreter.CompileLibrary(ctx, code, library)
	}
	
	if input, ok := params["input"].(string); ok && input != "" {
		return t.interpreter.ExecWithInput(ctx, code, input)
	}
	
	return t.interpreter.Exec(ctx, code)
}

func (t *CppTool) Close() error {
	return t.interpreter.Close()
}
