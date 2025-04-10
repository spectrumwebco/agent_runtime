package mcp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

func RegisterFilesystemResources(mcpServer *server.MCPServer) error {
	mcpServer.RegisterResourceHandler("file://", handleFileResource)
	
	mcpServer.RegisterResourceHandler("dir://", handleDirectoryResource)
	
	return nil
}

func handleFileResource(uri string) ([]byte, error) {
	path := strings.TrimPrefix(uri, "file://")
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return data, nil
}

func handleDirectoryResource(uri string) ([]byte, error) {
	path := strings.TrimPrefix(uri, "dir://")
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory not found: %s", path)
	}
	
	entries, err := os.ReadDir(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	var result strings.Builder
	
	for _, entry := range entries {
		entryType := "file"
		if entry.IsDir() {
			entryType = "dir"
		}
		
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}
		
		result.WriteString(fmt.Sprintf("%s\t%s\t%d\t%s\n", entryType, entry.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:05")))
	}
	
	return []byte(result.String()), nil
}

func ListFilesystemResources(path string) ([]string, error) {
	var resources []string
	
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			resources = append(resources, fmt.Sprintf("dir://%s", path))
		} else {
			resources = append(resources, fmt.Sprintf("file://%s", path))
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	return resources, nil
}
