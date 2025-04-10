package tools_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spectrumwebco/agent_runtime/internal/tools"
	pkgtools "github.com/spectrumwebco/agent_runtime/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileOperationTools(t *testing.T) {
	registry := pkgtools.NewRegistry()
	
	err := tools.RegisterCoreTools(registry)
	require.NoError(t, err)
	
	tempDir, err := os.MkdirTemp("", "tools-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	ctx := context.Background()
	
	t.Run("file_write", func(t *testing.T) {
		testFilePath := filepath.Join(tempDir, "test.txt")
		
		result, err := registry.ExecuteTool(ctx, "file_write", map[string]interface{}{
			"path":    testFilePath,
			"content": "Hello, world!",
		})
		
		require.NoError(t, err)
		assert.Equal(t, true, result)
		
		content, err := os.ReadFile(testFilePath)
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", string(content))
	})
	
	t.Run("file_read", func(t *testing.T) {
		testFilePath := filepath.Join(tempDir, "test.txt")
		
		err := os.WriteFile(testFilePath, []byte("Hello, world!"), 0644)
		require.NoError(t, err)
		
		result, err := registry.ExecuteTool(ctx, "file_read", map[string]interface{}{
			"path": testFilePath,
		})
		
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", result)
	})
	
	t.Run("file_str_replace", func(t *testing.T) {
		testFilePath := filepath.Join(tempDir, "test.txt")
		
		err := os.WriteFile(testFilePath, []byte("Hello, world!"), 0644)
		require.NoError(t, err)
		
		result, err := registry.ExecuteTool(ctx, "file_str_replace", map[string]interface{}{
			"path":    testFilePath,
			"old_str": "world",
			"new_str": "Go",
		})
		
		require.NoError(t, err)
		assert.Equal(t, true, result)
		
		content, err := os.ReadFile(testFilePath)
		require.NoError(t, err)
		assert.Equal(t, "Hello, Go!", string(content))
	})
	
	t.Run("file_exists", func(t *testing.T) {
		testFilePath := filepath.Join(tempDir, "test.txt")
		
		err := os.WriteFile(testFilePath, []byte("Hello, world!"), 0644)
		require.NoError(t, err)
		
		result, err := registry.ExecuteTool(ctx, "file_exists", map[string]interface{}{
			"path": testFilePath,
		})
		
		require.NoError(t, err)
		assert.Equal(t, true, result)
		
		result, err = registry.ExecuteTool(ctx, "file_exists", map[string]interface{}{
			"path": filepath.Join(tempDir, "non-existing.txt"),
		})
		
		require.NoError(t, err)
		assert.Equal(t, false, result)
	})
	
	t.Run("file_list", func(t *testing.T) {
		testFilePath1 := filepath.Join(tempDir, "test1.txt")
		testFilePath2 := filepath.Join(tempDir, "test2.txt")
		
		err := os.WriteFile(testFilePath1, []byte("Hello, world!"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(testFilePath2, []byte("Hello, Go!"), 0644)
		require.NoError(t, err)
		
		result, err := registry.ExecuteTool(ctx, "file_list", map[string]interface{}{
			"path": tempDir,
		})
		
		require.NoError(t, err)
		fileList, ok := result.([]string)
		require.True(t, ok)
		assert.Contains(t, fileList, "test1.txt")
		assert.Contains(t, fileList, "test2.txt")
	})
}

func TestShellOperationTools(t *testing.T) {
	registry := pkgtools.NewRegistry()
	
	err := tools.RegisterCoreTools(registry)
	require.NoError(t, err)
	
	ctx := context.Background()
	
	t.Run("shell_exec", func(t *testing.T) {
		result, err := registry.ExecuteTool(ctx, "shell_exec", map[string]interface{}{
			"command": "echo 'Hello, world!'",
		})
		
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Hello, world!")
	})
}
