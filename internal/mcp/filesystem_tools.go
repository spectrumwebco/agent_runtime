package mcp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"
)

func RegisterFilesystemTools(mcpServer *server.MCPServer) error {
	mcpServer.RegisterTool("file_read", "Reads a file", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the file",
			Required:    true,
		},
	}, handleFileRead)
	
	mcpServer.RegisterTool("file_write", "Writes to a file", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the file",
			Required:    true,
		},
		{
			Name:        "content",
			Type:        "string",
			Description: "Content to write",
			Required:    true,
		},
		{
			Name:        "append",
			Type:        "boolean",
			Description: "Whether to append to the file",
			Required:    false,
		},
	}, handleFileWrite)
	
	mcpServer.RegisterTool("file_delete", "Deletes a file", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the file",
			Required:    true,
		},
	}, handleFileDelete)
	
	mcpServer.RegisterTool("dir_create", "Creates a directory", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the directory",
			Required:    true,
		},
	}, handleDirCreate)
	
	mcpServer.RegisterTool("dir_list", "Lists directory contents", []server.ToolParameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to the directory",
			Required:    true,
		},
	}, handleDirList)
	
	return nil
}

func handleFileRead(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return string(data), nil
}

func handleFileWrite(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content parameter is required")
	}
	
	append, _ := params["append"].(bool)
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}
	
	var err error
	if append {
		f, err := os.OpenFile(filepath.Clean(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		
		if _, err := f.WriteString(content); err != nil {
			return nil, fmt.Errorf("failed to write to file: %w", err)
		}
	} else {
		err = os.WriteFile(filepath.Clean(path), []byte(content), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write to file: %w", err)
		}
	}
	
	return true, nil
}

func handleFileDelete(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	
	if err := os.Remove(filepath.Clean(path)); err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}
	
	return true, nil
}

func handleDirCreate(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	if err := os.MkdirAll(filepath.Clean(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	return true, nil
}

func handleDirList(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory not found: %s", path)
	}
	
	entries, err := os.ReadDir(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	var result []string
	
	for _, entry := range entries {
		result = append(result, entry.Name())
	}
	
	return result, nil
}
