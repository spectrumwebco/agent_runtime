package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/env"
	"github.com/spectrumwebco/agent_runtime/internal/ffi/python"
	"github.com/spectrumwebco/agent_runtime/pkg/tools/parsing"
)

type Command struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Arguments     []Argument `json:"arguments"`
	Category      string     `json:"category"`
	Timeout       int        `json:"timeout"`
	EndName       string     `json:"end_name,omitempty"`
	InvokeFormat  string     `json:"invoke_format,omitempty"`
	ArgumentFormat string    `json:"argument_format,omitempty"`
}

type Argument struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	Required      bool   `json:"required"`
	ArgumentFormat string `json:"argument_format,omitempty"`
}

type ToolFilterConfig struct {
	BlocklistErrorTemplate string   `json:"blocklist_error_template"`
	Blocklist             []string `json:"blocklist"`
	BlocklistStandalone   []string `json:"blocklist_standalone"`
	BlockUnlessRegex      map[string]string `json:"block_unless_regex"`
}

type Bundle struct {
	Path         string
	StateCommand string
	Commands     []*Command
}

type ToolConfig struct {
	ExecutionTimeout           time.Duration `json:"execution_timeout"`
	MaxOutputSize              int           `json:"max_output_size"`
	ToolsDir                   string        `json:"tools_dir"`
	Filter                     ToolFilterConfig `json:"filter"`
	Bundles                    []Bundle      `json:"bundles"`
	EnvVariables               map[string]interface{} `json:"env_variables"`
	RegistryVariables          map[string]interface{} `json:"registry_variables"`
	SubmitCommand              string        `json:"submit_command"`
	EnableBashTool             bool          `json:"enable_bash_tool"`
	FormatErrorTemplate        string        `json:"format_error_template"`
	MultiLineCommandEndings    map[string]string `json:"multi_line_command_endings"`
	SubmitCommandEndName       string        `json:"submit_command_end_name"`
	ResetCommands              []string      `json:"reset_commands"`
	InstallTimeout             int           `json:"install_timeout"`
	TotalExecutionTimeout      int           `json:"total_execution_timeout"`
	MaxConsecutiveTimeouts     int           `json:"max_consecutive_timeouts"`
}

type ToolRegistry struct {
	config       *ToolConfig
	commands     []*Command
	toolsPath    string
	scriptRunner *python.ScriptRunner
	commandPatterns map[string]*regexp.Regexp
	mockState    map[string]interface{}
}

func NewToolRegistry(config *ToolConfig, scriptRunner *python.ScriptRunner) (*ToolRegistry, error) {
	if config == nil {
		config = &ToolConfig{
			ExecutionTimeout: 60 * time.Second,
			MaxOutputSize:    10000,
			ToolsDir:        "tools",
			Filter: ToolFilterConfig{
				BlocklistErrorTemplate: "Operation '{{action}}' is not supported by this environment.",
				Blocklist: []string{
					"vim", "vi", "emacs", "nano", "nohup", "gdb", "less", "tail -f", "python -m venv", "make",
				},
				BlocklistStandalone: []string{
					"python", "python3", "ipython", "bash", "sh", "/bin/bash", "/bin/sh", "nohup", "vi", "vim", "emacs", "nano", "su",
				},
				BlockUnlessRegex: map[string]string{
					"radare2": `\b(?:radare2)\b.*\s+-c\s+.*`,
					"r2": `\b(?:radare2)\b.*\s+-c\s+.*`,
				},
			},
			SubmitCommand: "submit",
			EnableBashTool: true,
			InstallTimeout: 300,
			TotalExecutionTimeout: 1800,
			MaxConsecutiveTimeouts: 3,
			EnvVariables: make(map[string]interface{}),
			RegistryVariables: make(map[string]interface{}),
			MultiLineCommandEndings: make(map[string]string),
		}
	}
	
	if scriptRunner == nil {
		projectRoot, _ := os.Getwd()
		toolsFullPath := filepath.Join(projectRoot, config.ToolsDir)
		var err error
		scriptRunner, err = python.NewScriptRunner(toolsFullPath)
		if err != nil {
			fmt.Printf("Warning: Failed to create default script runner: %v. Python tools might not work.\n", err)
		} else {
			fmt.Println("Warning: No script runner provided, created default.")
		}
	}

	registry := &ToolRegistry{
		config:       config,
		commands:     make([]*Command, 0),
		scriptRunner: scriptRunner,
		commandPatterns: make(map[string]*regexp.Regexp),
		mockState:    nil,
	}

	if config.ToolsDir != "" {
		projectRoot, _ := os.Getwd()
		toolsFullPath := filepath.Join(projectRoot, config.ToolsDir)
		if err := registry.LoadTools(toolsFullPath); err != nil {
			fmt.Printf("Warning: Failed to load tools from %s during initialization: %v\n", toolsFullPath, err)
		} else {
			fmt.Printf("Successfully loaded tools from %s\n", toolsFullPath)
		}
	}

	registry.initCommandPatterns()

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
	
	if r.config.EnableBashTool {
		bashCommand := &Command{
			Name:        "bash",
			Description: "Execute a bash command",
			Arguments: []Argument{
				{
					Name:        "command",
					Type:        "string",
					Description: "The command to execute",
					Required:    true,
				},
			},
			Category: "system",
			Timeout:  30,
		}
		r.commands = append(r.commands, bashCommand)
	}
	
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
	if r.shouldBlockAction(action) {
		return "", fmt.Errorf(r.config.Filter.BlocklistErrorTemplate, map[string]string{"action": action})
	}

	parser := parsing.NewCmdParser()
	parsedCmd, err := parser.Parse(action)
	if err != nil {
		return "", fmt.Errorf("failed to parse action '%s': %w", action, err)
	}

	cmdName := parsedCmd.Name
	args := parsedCmd.Args

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

func (r *ToolRegistry) GetState(env *env.SWEEnv) (map[string]interface{}, error) {
	if r.mockState != nil {
		return r.mockState, nil
	}

	for _, bundle := range r.config.Bundles {
		if bundle.StateCommand != "" {
			_, err := env.ExecuteCommand(bundle.StateCommand)
			if err != nil {
				fmt.Printf("Warning: Failed to execute state command: %v\n", err)
			}
		}
	}

	stateFile := "/root/state.json"
	stateStr, err := env.ReadFile(stateFile)
	if err != nil {
		return make(map[string]interface{}), nil
	}

	if stateStr == "" {
		return make(map[string]interface{}), nil
	}

	var state map[string]interface{}
	if err := json.Unmarshal([]byte(stateStr), &state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}

	return state, nil
}

func (r *ToolRegistry) SetMockState(state map[string]interface{}) {
	r.mockState = state
}

func (r *ToolRegistry) Install(env *env.SWEEnv) error {
	err := env.SetEnvVariables(r.config.EnvVariables)
	if err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	registryJSON, err := json.Marshal(r.config.RegistryVariables)
	if err != nil {
		return fmt.Errorf("failed to marshal registry variables: %w", err)
	}
	err = env.CreateFile("/root/.swe-agent-env", string(registryJSON))
	if err != nil {
		return fmt.Errorf("failed to write registry variables: %w", err)
	}

	err = env.CreateFile("/root/state.json", "{}")
	if err != nil {
		return fmt.Errorf("failed to initialize state: %w", err)
	}

	for _, cmd := range r.config.ResetCommands {
		_, err := env.ExecuteCommand(cmd)
		if err != nil {
			return fmt.Errorf("failed to execute reset command: %w", err)
		}
	}

	return nil
}

func (r *ToolRegistry) Reset(env *env.SWEEnv) error {
	return r.Install(env)
}

func (r *ToolRegistry) shouldBlockAction(action string) bool {
	action = strings.TrimSpace(action)
	if action == "" {
		return false
	}
	
	for _, blocked := range r.config.Filter.Blocklist {
		if strings.HasPrefix(action, blocked) {
			return true
		}
	}
	
	if contains(r.config.Filter.BlocklistStandalone, action) {
		return true
	}
	
	name := strings.Split(action, " ")[0]
	if pattern, ok := r.config.Filter.BlockUnlessRegex[name]; ok {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Warning: Invalid regex pattern '%s': %v\n", pattern, err)
			return true
		}
		if !re.MatchString(action) {
			return true
		}
	}
	
	return false
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

func (r *ToolRegistry) initCommandPatterns() {
	for _, cmd := range r.commands {
		if cmd.EndName != "" {
			pattern := regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(.*?)^(%s)\s*$`, regexp.QuoteMeta(cmd.Name), regexp.QuoteMeta(cmd.EndName)))
			r.commandPatterns[cmd.Name] = pattern
		} else {
			pattern := regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(.*?)$`, regexp.QuoteMeta(cmd.Name)))
			r.commandPatterns[cmd.Name] = pattern
		}
	}
	
	if r.config.SubmitCommand != "" && r.config.SubmitCommandEndName != "" {
		pattern := regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(.*?)^(%s)\s*$`, 
			regexp.QuoteMeta(r.config.SubmitCommand), 
			regexp.QuoteMeta(r.config.SubmitCommandEndName)))
		r.commandPatterns[r.config.SubmitCommand] = pattern
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
