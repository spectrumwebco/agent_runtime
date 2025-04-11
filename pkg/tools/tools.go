package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Command struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Arguments   []Argument `json:"arguments"`
	Category    string     `json:"category"`
	Timeout     int        `json:"timeout"`
}

type Argument struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type ToolConfig struct {
	ExecutionTimeout time.Duration `json:"execution_timeout"`
	MaxOutputSize    int          `json:"max_output_size"`
	ToolsDir        string        `json:"tools_dir"`
}

type ToolRegistry struct {
	config    *ToolConfig
	commands  []*Command
	toolsPath string
}

func NewToolRegistry(config *ToolConfig) (*ToolRegistry, error) {
	if config == nil {
		config = &ToolConfig{
			ExecutionTimeout: 60 * time.Second,
			MaxOutputSize:    10000,
			ToolsDir:        "tools",
		}
	}

	registry := &ToolRegistry{
		config:   config,
		commands: make([]*Command, 0),
	}

	return registry, nil
}

func (r *ToolRegistry) LoadTools(toolsPath string) error {
	r.toolsPath = toolsPath
	
	configPath := filepath.Join(toolsPath, "tools.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read tools config: %w", err)
	}

	var tools struct {
		Commands []*Command `json:"commands"`
	}
	if err := json.Unmarshal(data, &tools); err != nil {
		return fmt.Errorf("failed to parse tools config: %w", err)
	}

	r.commands = tools.Commands
	return nil
}

func (r *ToolRegistry) GetCommand(name string) (*Command, error) {
	for _, cmd := range r.commands {
		if cmd.Name == name {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("command not found: %s", name)
}

func (r *ToolRegistry) ListCommands() []*Command {
	return r.commands
}

func (r *ToolRegistry) ExecuteAction(action string, env interface{}) (string, error) {
	cmd, args, err := r.parseAction(action)
	if err != nil {
		return "", fmt.Errorf("failed to parse action: %w", err)
	}

	command, err := r.GetCommand(cmd)
	if err != nil {
		return "", err
	}

	if err := r.validateArguments(command, args); err != nil {
		return "", err
	}

	scriptPath := filepath.Join(r.toolsPath, cmd, "run.py")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("command implementation not found: %s", cmd)
	}

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("failed to marshal arguments: %w", err)
	}

	return fmt.Sprintf("Executed %s with args %s", cmd, string(argsJSON)), nil
}

func (r *ToolRegistry) parseAction(action string) (string, map[string]interface{}, error) {
	var result struct {
		Command string                 `json:"command"`
		Args    map[string]interface{} `json:"args"`
	}

	if err := json.Unmarshal([]byte(action), &result); err != nil {
		return "", nil, fmt.Errorf("invalid action format: %w", err)
	}

	return result.Command, result.Args, nil
}

func (r *ToolRegistry) validateArguments(cmd *Command, args map[string]interface{}) error {
	for _, arg := range cmd.Arguments {
		if arg.Required {
			if _, ok := args[arg.Name]; !ok {
				return fmt.Errorf("required argument missing: %s", arg.Name)
			}
		}
	}
	return nil
}
