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
	config       *ToolConfig
	commands     []*Command
	toolsPath    string
	scriptRunner *python.ScriptRunner // Added script runner for FFI
}

func NewToolRegistry(config *ToolConfig, scriptRunner *python.ScriptRunner) (*ToolRegistry, error) {
	if config == nil {
		config = &ToolConfig{
			ExecutionTimeout: 60 * time.Second,
			MaxOutputSize:    10000,
			ToolsDir:        "tools", // Default, should be configurable
		}
	}
	if scriptRunner == nil {
		projectRoot, _ := os.Getwd() // Or determine project root differently
		toolsFullPath := filepath.Join(projectRoot, config.ToolsDir)
		var err error
		scriptRunner, err = python.NewScriptRunner(toolsFullPath)
		if err != nil {
			fmt.Printf("Warning: Failed to create default script runner: %v. Python tools might not work.\n", err)
			// return nil, fmt.Errorf("failed to create default script runner: %w", err)
		} else {
			fmt.Println("Warning: No script runner provided, created default.")
		}
	}

	registry := &ToolRegistry{
		config:       config,
		commands:     make([]*Command, 0),
		scriptRunner: scriptRunner, // Assign even if nil (indicates FFI issues)
	}

	if config.ToolsDir != "" {
		projectRoot, _ := os.Getwd() // Or determine project root differently
		toolsFullPath := filepath.Join(projectRoot, config.ToolsDir)
		if err := registry.LoadTools(toolsFullPath); err != nil {
			fmt.Printf("Warning: Failed to load tools from %s during initialization: %v\n", toolsFullPath, err)
		} else {
			fmt.Printf("Successfully loaded tools from %s\n", toolsFullPath)
		}
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
	parser := parsing.NewCmdParser() // Assuming CmdParser exists and is ported
	parsedCmd, err := parser.Parse(action)
	if err != nil {
		return "", fmt.Errorf("failed to parse action '%s': %w", action, err)
	}

	cmdName := parsedCmd.Name
	args := parsedCmd.Args // Args should be map[string]interface{}

	commandDef, err := r.GetCommand(cmdName)
	if err != nil {
		fmt.Printf("Warning: Command definition for '%s' not found in tools.json\n", cmdName)
	} else {
		if err := r.validateArguments(commandDef, args); err != nil {
			return "", fmt.Errorf("argument validation failed for command '%s': %w", cmdName, err)
		}
	}

	if r.scriptRunner == nil {
		return "", fmt.Errorf("script runner is not initialized or failed to initialize")
	}

	resultInterface, err := r.scriptRunner.RunToolScript(cmdName, args)
	if err != nil {
		return "", fmt.Errorf("failed to execute tool '%s' via FFI: %w", cmdName, err)
	}

	var resultStr string
	switch v := resultInterface.(type) {
	case string:
		resultStr = v
	case []byte:
		resultStr = string(v)
	default:
		jsonBytes, jsonErr := json.MarshalIndent(v, "", "  ")
		if jsonErr != nil {
			resultStr = fmt.Sprintf("%+v", v) // Fallback
		} else {
			resultStr = string(jsonBytes)
		}
	}

	if len(resultStr) > r.config.MaxOutputSize {
		resultStr = resultStr[:r.config.MaxOutputSize] + "\n... (output truncated)"
	}

	return resultStr, nil
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
