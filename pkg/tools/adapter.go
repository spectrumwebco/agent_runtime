package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/spectrumwebco/agent_runtime/internal/config"
)

func LoadToolsFromConfig(toolRegistry *ToolRegistry) error {
	repoRoot, err := config.FindRepoRoot()
	if err != nil {
		return fmt.Errorf("failed to find repository root: %w", err)
	}
	
	adapter, err := config.NewConfigAdapter(config.PkgSource)
	if err != nil {
		return fmt.Errorf("failed to create config adapter: %w", err)
	}
	
	toolsConfig, err := adapter.LoadTools()
	if err != nil {
		return fmt.Errorf("failed to load tools configuration: %w", err)
	}
	
	toolsArray, ok := toolsConfig["tools"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid tools format")
	}
	
	commands := make([]*Command, 0, len(toolsArray))
	for _, toolInterface := range toolsArray {
		tool, ok := toolInterface.(map[string]interface{})
		if !ok {
			continue
		}
		
		command := &Command{
			Name:        tool["name"].(string),
			Description: tool["description"].(string),
		}
		
		if parameters, ok := tool["parameters"].(map[string]interface{}); ok {
			for paramName, paramInterface := range parameters {
				param, ok := paramInterface.(map[string]interface{})
				if !ok {
					continue
				}
				
				argument := Argument{
					Name:        paramName,
					Type:        param["type"].(string),
					Description: param["description"].(string),
					Required:    true,
				}
				
				command.Arguments = append(command.Arguments, argument)
			}
		}
		
		commands = append(commands, command)
	}
	
	toolRegistry.commands = commands
	
	return nil
}

func SaveToolsToConfig(toolRegistry *ToolRegistry) error {
	repoRoot, err := config.FindRepoRoot()
	if err != nil {
		return fmt.Errorf("failed to find repository root: %w", err)
	}
	
	adapter, err := config.NewConfigAdapter(config.PkgSource)
	if err != nil {
		return fmt.Errorf("failed to create config adapter: %w", err)
	}
	
	toolsConfig := map[string]interface{}{
		"tools": toolRegistry.commands,
	}
	
	data, err := json.MarshalIndent(toolsConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tools configuration: %w", err)
	}
	
	pkgToolsPath := filepath.Join(repoRoot, "pkg", "tools", "tools.json")
	if err := ioutil.WriteFile(pkgToolsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tools configuration: %w", err)
	}
	
	if err := adapter.UpdateGoConfigsFromPkg(); err != nil {
		return fmt.Errorf("failed to sync tools configuration to go_configs: %w", err)
	}
	
	return nil
}
