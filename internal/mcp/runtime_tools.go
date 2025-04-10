package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spectrumwebco/agent_runtime/internal/cpp"
	"github.com/spectrumwebco/agent_runtime/internal/python"
)

func RegisterRuntimeTools(mcpServer *server.MCPServer) error {
	mcpServer.RegisterTool("python_exec", "Executes Python code", []server.ToolParameter{
		{
			Name:        "code",
			Type:        "string",
			Description: "Python code to execute",
			Required:    true,
		},
		{
			Name:        "input",
			Type:        "string",
			Description: "Input to provide to the code",
			Required:    false,
		},
	}, handlePythonExec)
	
	mcpServer.RegisterTool("cpp_exec", "Executes C++ code", []server.ToolParameter{
		{
			Name:        "code",
			Type:        "string",
			Description: "C++ code to execute",
			Required:    true,
		},
		{
			Name:        "input",
			Type:        "string",
			Description: "Input to provide to the code",
			Required:    false,
		},
		{
			Name:        "flags",
			Type:        "array",
			Description: "Compiler flags",
			Required:    false,
		},
		{
			Name:        "include_dirs",
			Type:        "array",
			Description: "Include directories",
			Required:    false,
		},
		{
			Name:        "libraries",
			Type:        "array",
			Description: "Libraries to link",
			Required:    false,
		},
		{
			Name:        "library",
			Type:        "string",
			Description: "Compile as a library with this name",
			Required:    false,
		},
	}, handleCppExec)
	
	mcpServer.RegisterTool("sandbox_exec", "Executes a command in a sandbox", []server.ToolParameter{
		{
			Name:        "command",
			Type:        "string",
			Description: "Command to execute",
			Required:    true,
		},
		{
			Name:        "timeout",
			Type:        "number",
			Description: "Timeout in seconds",
			Required:    false,
		},
	}, handleSandboxExec)
	
	return nil
}

func handlePythonExec(params map[string]interface{}) (interface{}, error) {
	code, ok := params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code parameter is required")
	}
	
	interpreter := python.NewInterpreter()
	defer interpreter.Close()
	
	if input, ok := params["input"].(string); ok && input != "" {
		return interpreter.ExecWithInput(code, input)
	}
	
	return interpreter.Exec(code)
}

func handleCppExec(params map[string]interface{}) (interface{}, error) {
	code, ok := params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("code parameter is required")
	}
	
	config := cpp.InterpreterConfig{
		Flags: []string{"-std=c++17", "-O2"},
	}
	
	if flagsParam, ok := params["flags"].([]interface{}); ok {
		flags := make([]string, 0, len(flagsParam))
		for _, flag := range flagsParam {
			if flagStr, ok := flag.(string); ok {
				flags = append(flags, flagStr)
			}
		}
		config.Flags = flags
	}
	
	if includeDirsParam, ok := params["include_dirs"].([]interface{}); ok {
		includeDirs := make([]string, 0, len(includeDirsParam))
		for _, dir := range includeDirsParam {
			if dirStr, ok := dir.(string); ok {
				includeDirs = append(includeDirs, dirStr)
			}
		}
		config.IncludeDirs = includeDirs
	}
	
	if librariesParam, ok := params["libraries"].([]interface{}); ok {
		libraries := make([]string, 0, len(librariesParam))
		for _, lib := range librariesParam {
			if libStr, ok := lib.(string); ok {
				libraries = append(libraries, libStr)
			}
		}
		config.Libraries = libraries
	}
	
	interpreter, err := cpp.NewInterpreter(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create C++ interpreter: %w", err)
	}
	defer interpreter.Close()
	
	ctx := context.Background()
	
	if library, ok := params["library"].(string); ok && library != "" {
		return interpreter.CompileLibrary(ctx, code, library)
	}
	
	if input, ok := params["input"].(string); ok && input != "" {
		return interpreter.ExecWithInput(ctx, code, input)
	}
	
	return interpreter.Exec(ctx, code)
}

func handleSandboxExec(params map[string]interface{}) (interface{}, error) {
	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}
	
	timeout := 30.0
	if timeoutParam, ok := params["timeout"].(float64); ok {
		timeout = timeoutParam
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command execution timed out after %.1f seconds", timeout)
		}
		return nil, fmt.Errorf("command execution failed: %w\nOutput: %s", err, output)
	}
	
	return string(output), nil
}
