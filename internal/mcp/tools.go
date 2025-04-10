package mcp

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type ToolsServer struct {
	server      *server.MCPServer
	toolManager *tools.Registry
}

func NewToolsServer(toolManager *tools.Registry) (*ToolsServer, error) {
	mcpServer := server.NewMCPServer(
		"agent-runtime/tools",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	
	ts := &ToolsServer{
		server:      mcpServer,
		toolManager: toolManager,
	}
	
	ts.registerTools()
	
	return ts, nil
}

func (ts *ToolsServer) registerTools() {
	ts.server.AddTool(mcp.NewTool("execute",
		mcp.WithDescription("Executes a tool"),
		mcp.WithString("tool",
			mcp.Description("Name of the tool to execute"),
			mcp.Required(),
		),
		mcp.WithObject("args",
			mcp.Description("Arguments for the tool"),
			mcp.Required(),
		),
	), ts.handleExecuteTool)
	
	ts.server.AddTool(mcp.NewTool("list_tools",
		mcp.WithDescription("Lists available tools"),
	), ts.handleListToolsTool)
}

func (ts *ToolsServer) handleExecuteTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	toolName, ok1 := request.Params.Arguments["tool"].(string)
	args, ok2 := request.Params.Arguments["args"].(map[string]interface{})
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid arguments")
	}
	
	tool, err := ts.toolManager.Get(toolName)
	if err != nil {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}
	
	result, err := tool.Execute(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}
	
	var resultStr string
	switch v := result.(type) {
	case string:
		resultStr = v
	case []byte:
		resultStr = string(v)
	default:
		resultStr = fmt.Sprintf("%v", v)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: resultStr,
			},
		},
	}, nil
}

func (ts *ToolsServer) handleListToolsTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	tools := ts.toolManager.List()
	
	var result string
	for _, tool := range tools {
		result += fmt.Sprintf("%s: %s\n", tool.Name(), tool.Description())
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}
