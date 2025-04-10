package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type FilesystemServer struct {
	server         *server.MCPServer
	allowedDirs    []string
	defaultRootDir string
}

func NewFilesystemServer(allowedDirs []string, defaultRootDir string) (*FilesystemServer, error) {
	mcpServer := server.NewMCPServer(
		"agent-runtime/filesystem",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
	)
	
	fs := &FilesystemServer{
		server:         mcpServer,
		allowedDirs:    allowedDirs,
		defaultRootDir: defaultRootDir,
	}
	
	fs.registerTools()
	
	return fs, nil
}

func (fs *FilesystemServer) registerTools() {
	fs.server.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Reads a file from the filesystem"),
		mcp.WithString("path",
			mcp.Description("Path to the file to read"),
			mcp.Required(),
		),
	), fs.handleReadFileTool)
	
	fs.server.AddTool(mcp.NewTool("write_file",
		mcp.WithDescription("Writes content to a file in the filesystem"),
		mcp.WithString("path",
			mcp.Description("Path to the file to write"),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Content to write to the file"),
			mcp.Required(),
		),
	), fs.handleWriteFileTool)
	
	fs.server.AddTool(mcp.NewTool("list_directory",
		mcp.WithDescription("Lists the contents of a directory"),
		mcp.WithString("path",
			mcp.Description("Path to the directory to list"),
			mcp.Required(),
		),
	), fs.handleListDirectoryTool)
	
	fs.server.AddTool(mcp.NewTool("create_directory",
		mcp.WithDescription("Creates a new directory"),
		mcp.WithString("path",
			mcp.Description("Path to the directory to create"),
			mcp.Required(),
		),
	), fs.handleCreateDirectoryTool)
}

func (fs *FilesystemServer) isPathAllowed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	for _, dir := range fs.allowedDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		
		if filepath.HasPrefix(absPath, absDir) {
			return true
		}
	}
	
	return false
}

func (fs *FilesystemServer) handleReadFileTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}
	
	if !fs.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed: %s", path)
	}
	
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(content),
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleWriteFileTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	path, ok1 := request.Params.Arguments["path"].(string)
	content, ok2 := request.Params.Arguments["content"].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid arguments")
	}
	
	if !fs.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed: %s", path)
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(filepath.Clean(path), []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("File written successfully: %s", path),
			},
		},
	}, nil
}

func (fs *FilesystemServer) handleListDirectoryTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}
	
	if !fs.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed: %s", path)
	}
	
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	var result string
	for _, entry := range entries {
		if entry.IsDir() {
			result += fmt.Sprintf("d %s\n", entry.Name())
		} else {
			result += fmt.Sprintf("f %s\n", entry.Name())
		}
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

func (fs *FilesystemServer) handleCreateDirectoryTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	path, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid path argument")
	}
	
	if !fs.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed: %s", path)
	}
	
	if err := os.MkdirAll(filepath.Clean(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Directory created successfully: %s", path),
			},
		},
	}, nil
}
