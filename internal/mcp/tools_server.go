package mcp

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

func RegisterToolsServerTools(mcpServer *server.MCPServer) error {
	mcpServer.RegisterTool("shell_exec", "Executes a shell command", []server.ToolParameter{
		{
			Name:        "command",
			Type:        "string",
			Description: "Command to execute",
			Required:    true,
		},
	}, handleShellExec)
	
	mcpServer.RegisterTool("git_clone", "Clones a git repository", []server.ToolParameter{
		{
			Name:        "url",
			Type:        "string",
			Description: "URL of the repository",
			Required:    true,
		},
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to clone to",
			Required:    true,
		},
	}, handleGitClone)
	
	mcpServer.RegisterTool("go_build", "Builds a Go project", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the project",
			Required:    true,
		},
		{
			Name:        "output",
			Type:        "string",
			Description: "Output file",
			Required:    false,
		},
	}, handleGoBuild)
	
	mcpServer.RegisterTool("go_test", "Runs Go tests", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the project",
			Required:    true,
		},
		{
			Name:        "verbose",
			Type:        "boolean",
			Description: "Verbose output",
			Required:    false,
		},
	}, handleGoTest)
	
	return nil
}

func handleShellExec(params map[string]interface{}) (interface{}, error) {
	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}
	
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %w\nOutput: %s", err, output)
	}
	
	return string(output), nil
}

func handleGitClone(params map[string]interface{}) (interface{}, error) {
	url, ok := params["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url parameter is required")
	}
	
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	cmd := exec.Command("git", "clone", url, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git clone failed: %w\nOutput: %s", err, output)
	}
	
	return true, nil
}

func handleGoBuild(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	args := []string{"build"}
	
	if output, ok := params["output"].(string); ok && output != "" {
		args = append(args, "-o", output)
	}
	
	args = append(args, "./...")
	
	cmd := exec.Command("go", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("go build failed: %w\nOutput: %s", err, output)
	}
	
	return true, nil
}

func handleGoTest(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	args := []string{"test"}
	
	if verbose, ok := params["verbose"].(bool); ok && verbose {
		args = append(args, "-v")
	}
	
	args = append(args, "./...")
	
	cmd := exec.Command("go", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	
	return string(output), nil
}
