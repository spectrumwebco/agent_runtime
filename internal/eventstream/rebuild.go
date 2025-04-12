package eventstream

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type ContextRebuildFunc func(identifier string) (string, error)

type ContextRebuildRegistry map[string]ContextRebuildFunc

func NewContextRebuildRegistry() ContextRebuildRegistry {
	registry := make(ContextRebuildRegistry)
	
	registry["file"] = rebuildFileContext
	registry["k8s"] = rebuildK8sContext
	registry["git"] = rebuildGitContext
	registry["ci"] = rebuildCIContext
	registry["state"] = rebuildStateContext
	
	return registry
}

func rebuildFileContext(filePath string) (string, error) {
	log.Printf("Rebuilding file context for: %s\n", filePath)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	return string(content), nil
}

func rebuildK8sContext(resourcePath string) (string, error) {
	log.Printf("Rebuilding k8s context for: %s\n", resourcePath)
	
	parts := strings.SplitN(resourcePath, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid k8s resource path format: %s", resourcePath)
	}
	
	resourceType := parts[0]
	resourceName := parts[1]
	
	
	switch resourceType {
	case "pod":
		return fmt.Sprintf("Placeholder for pod %s information", resourceName), nil
	case "deployment":
		return fmt.Sprintf("Placeholder for deployment %s information", resourceName), nil
	case "service":
		return fmt.Sprintf("Placeholder for service %s information", resourceName), nil
	default:
		return "", fmt.Errorf("unsupported k8s resource type: %s", resourceType)
	}
}

func rebuildGitContext(gitCommand string) (string, error) {
	log.Printf("Rebuilding git context for: %s\n", gitCommand)
	
	
	switch gitCommand {
	case "status":
		return "Placeholder for git status output", nil
	case "log":
		return "Placeholder for git log output", nil
	case "diff":
		return "Placeholder for git diff output", nil
	default:
		return "", fmt.Errorf("unsupported git command: %s", gitCommand)
	}
}

func rebuildCIContext(pipelineID string) (string, error) {
	log.Printf("Rebuilding CI context for pipeline: %s\n", pipelineID)
	
	
	return fmt.Sprintf("Placeholder for CI pipeline %s information", pipelineID), nil
}

func rebuildStateContext(stateKey string) (string, error) {
	log.Printf("Rebuilding state context for: %s\n", stateKey)
	
	
	return fmt.Sprintf("Placeholder for application state %s", stateKey), nil
}

func GetCacheExpiration(contextType string) time.Duration {
	switch contextType {
	case "file":
		return 30 * time.Minute // Files might change frequently during development
	case "k8s":
		return 5 * time.Minute // K8s resources can change often
	case "git":
		return 10 * time.Minute // Git state changes with commits
	case "ci":
		return 15 * time.Minute // CI pipelines update periodically
	case "state":
		return 1 * time.Hour // Application state might be more stable
	default:
		return 30 * time.Minute // Default TTL
	}
}
