package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Tool interface {
	Name() string
	
	Description() string
	
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

type ShellTool struct {
	name        string
	description string
	sandbox     bool
}

func (t *ShellTool) Name() string {
	return t.name
}

func (t *ShellTool) Description() string {
	return t.description
}

func (t *ShellTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	command, ok := args["command"].(string)
	if !ok {
		return nil, fmt.Errorf("missing command argument")
	}
	
	if t.sandbox && !isCommandAllowed(command) {
		return nil, fmt.Errorf("command not allowed in sandbox mode: %s", command)
	}
	
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w\nOutput: %s", err, output)
	}
	
	return string(output), nil
}

func isCommandAllowed(command string) bool {
	disallowed := []string{
		"rm -rf /",
		"rm -rf /*",
		"mkfs",
		"dd",
		"wget",
		"curl",
	}
	
	for _, d := range disallowed {
		if strings.Contains(command, d) {
			return false
		}
	}
	
	return true
}

type FileTool struct {
	name         string
	description  string
	sandbox      bool
	allowedPaths []string
}

func (t *FileTool) Name() string {
	return t.name
}

func (t *FileTool) Description() string {
	return t.description
}

func (t *FileTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("missing operation argument")
	}
	
	switch operation {
	case "read":
		return t.read(args)
	case "write":
		return t.write(args)
	case "list":
		return t.list(args)
	case "delete":
		return t.delete(args)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

func (t *FileTool) read(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing path argument")
	}
	
	if t.sandbox && !t.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed in sandbox mode: %s", path)
	}
	
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return string(content), nil
}

func (t *FileTool) write(args map[string]interface{}) (interface{}, error) {
	path, ok1 := args["path"].(string)
	content, ok2 := args["content"].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("missing path or content argument")
	}
	
	if t.sandbox && !t.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed in sandbox mode: %s", path)
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(filepath.Clean(path), []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	
	return fmt.Sprintf("File written successfully: %s", path), nil
}

func (t *FileTool) list(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing path argument")
	}
	
	if t.sandbox && !t.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed in sandbox mode: %s", path)
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
	
	return result, nil
}

func (t *FileTool) delete(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing path argument")
	}
	
	if t.sandbox && !t.isPathAllowed(path) {
		return nil, fmt.Errorf("path not allowed in sandbox mode: %s", path)
	}
	
	if err := os.RemoveAll(filepath.Clean(path)); err != nil {
		return nil, fmt.Errorf("failed to delete: %w", err)
	}
	
	return fmt.Sprintf("Deleted successfully: %s", path), nil
}

func (t *FileTool) isPathAllowed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	for _, allowedPath := range t.allowedPaths {
		absAllowedPath, err := filepath.Abs(allowedPath)
		if err != nil {
			continue
		}
		
		if strings.HasPrefix(absPath, absAllowedPath) {
			return true
		}
	}
	
	return false
}

type HTTPTool struct {
	name        string
	description string
}

func (t *HTTPTool) Name() string {
	return t.name
}

func (t *HTTPTool) Description() string {
	return t.description
}

func (t *HTTPTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	url, ok1 := args["url"].(string)
	method, ok2 := args["method"].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("missing url or method argument")
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	if headers, ok := args["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Add(key, strValue)
			}
		}
	}
	
	if body, ok := args["body"].(string); ok && body != "" {
		req.Body = io.NopCloser(strings.NewReader(body))
		req.ContentLength = int64(len(body))
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	result := map[string]interface{}{
		"status":  resp.StatusCode,
		"headers": resp.Header,
		"body":    string(respBody),
	}
	
	return result, nil
}
