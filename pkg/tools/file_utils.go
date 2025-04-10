package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	return os.WriteFile(path, []byte(content), 0644)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	return string(data), nil
}
