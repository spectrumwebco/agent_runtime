package mcp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type RuntimeServer struct {
	server      *server.MCPServer
	sandboxDir  string
	environment map[string]string
}

func NewRuntimeServer(sandboxDir string, environment map[string]string) (*RuntimeServer, error) {
	mcpServer := server.NewMCPServer(
		"agent-runtime/runtime",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	
	rs := &RuntimeServer{
		server:      mcpServer,
		sandboxDir:  sandboxDir,
		environment: environment,
	}
	
	rs.registerTools()
	
	return rs, nil
}

func (rs *RuntimeServer) registerTools() {
	rs.server.AddTool(mcp.NewTool("execute_command",
		mcp.WithDescription("Executes a command in the runtime environment"),
		mcp.WithString("command",
			mcp.Description("Command to execute"),
			mcp.Required(),
		),
		mcp.WithBoolean("sandbox",
			mcp.Description("Whether to execute the command in a sandbox"),
			mcp.Default(true),
		),
	), rs.handleExecuteCommandTool)
	
	rs.server.AddTool(mcp.NewTool("create_session",
		mcp.WithDescription("Creates a new session in the runtime environment"),
		mcp.WithString("name",
			mcp.Description("Name of the session"),
			mcp.Required(),
		),
	), rs.handleCreateSessionTool)
	
	rs.server.AddTool(mcp.NewTool("run_in_session",
		mcp.WithDescription("Runs a command in a session"),
		mcp.WithString("session",
			mcp.Description("Name of the session"),
			mcp.Required(),
		),
		mcp.WithString("command",
			mcp.Description("Command to execute"),
			mcp.Required(),
		),
	), rs.handleRunInSessionTool)
}

func (rs *RuntimeServer) handleExecuteCommandTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	command, ok1 := request.Params.Arguments["command"].(string)
	sandbox, ok2 := request.Params.Arguments["sandbox"].(bool)
	if !ok1 {
		return nil, fmt.Errorf("invalid command argument")
	}
	if !ok2 {
		sandbox = true
	}
	
	var cmd *exec.Cmd
	if sandbox {
		if err := os.MkdirAll(rs.sandboxDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create sandbox directory: %w", err)
		}
		
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
		cmd.Dir = rs.sandboxDir
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
	}
	
	cmd.Env = os.Environ()
	for key, value := range rs.environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w\nOutput: %s", err, output)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

func (rs *RuntimeServer) handleCreateSessionTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name argument")
	}
	
	sessionDir := filepath.Join(rs.sandboxDir, "sessions", name)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Session created: %s", name),
			},
		},
	}, nil
}

func (rs *RuntimeServer) handleRunInSessionTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	session, ok1 := request.Params.Arguments["session"].(string)
	command, ok2 := request.Params.Arguments["command"].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid arguments")
	}
	
	sessionDir := filepath.Join(rs.sandboxDir, "sessions", session)
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("session not found: %s", session)
	}
	
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = sessionDir
	
	cmd.Env = os.Environ()
	for key, value := range rs.environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w\nOutput: %s", err, output)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}
