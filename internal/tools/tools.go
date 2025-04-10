package tools

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type ToolFilterConfig struct {
	BlocklistErrorTemplate string   `json:"blocklist_error_template"`
	Blocklist             []string `json:"blocklist"`
	BlocklistStandalone   []string `json:"blocklist_standalone"`
	BlockUnlessRegex      map[string]string `json:"block_unless_regex"`
}

func DefaultToolFilterConfig() *ToolFilterConfig {
	return &ToolFilterConfig{
		BlocklistErrorTemplate: "Operation '{{action}}' is not supported by this environment.",
		Blocklist: []string{
			"vim",
			"vi",
			"emacs",
			"nano",
			"nohup",
			"gdb",
			"less",
			"tail -f",
			"python -m venv",
			"make",
		},
		BlocklistStandalone: []string{
			"python",
			"python3",
			"ipython",
			"bash",
			"sh",
			"/bin/bash",
			"/bin/sh",
			"nohup",
			"vi",
			"vim",
			"emacs",
			"nano",
			"su",
		},
		BlockUnlessRegex: map[string]string{
			"radare2": `\b(?:radare2)\b.*\s+-c\s+.*`,
			"r2":      `\b(?:radare2)\b.*\s+-c\s+.*`,
		},
	}
}

type ToolConfig struct {
	Filter              *ToolFilterConfig    `json:"filter"`
	Bundles             []Bundle             `json:"bundles"`
	EnvVariables        map[string]interface{} `json:"env_variables"`
	RegistryVariables   map[string]interface{} `json:"registry_variables"`
	SubmitCommand       string               `json:"submit_command"`
	EnableBashTool      bool                 `json:"enable_bash_tool"`
	FormatErrorTemplate string               `json:"format_error_template"`
	CommandDocs         string               `json:"command_docs"`
	MultiLineCommandEndings map[string]string `json:"multi_line_command_endings"`
	SubmitCommandEndName string               `json:"submit_command_end_name"`
	ResetCommands       []interface{}        `json:"reset_commands"`
	ExecutionTimeout    int                  `json:"execution_timeout"`
	InstallTimeout      int                  `json:"install_timeout"`
	TotalExecutionTimeout int                `json:"total_execution_timeout"`
	MaxConsecutiveExecutionTimeouts int      `json:"max_consecutive_execution_timeouts"`
	
	commands            []Command
	tools               []map[string]interface{}
}

func NewToolConfig() *ToolConfig {
	return &ToolConfig{
		Filter:              DefaultToolFilterConfig(),
		Bundles:             []Bundle{},
		EnvVariables:        map[string]interface{}{},
		RegistryVariables:   map[string]interface{}{},
		SubmitCommand:       "submit",
		EnableBashTool:      true,
		MultiLineCommandEndings: map[string]string{},
		ResetCommands:       []interface{}{},
		ExecutionTimeout:    30,
		InstallTimeout:      300,
		TotalExecutionTimeout: 1800,
		MaxConsecutiveExecutionTimeouts: 3,
	}
}

func (c *ToolConfig) UseFunctionCalling() bool {
	return true
}

func (c *ToolConfig) StateCommands() []string {
	var commands []string
	for _, bundle := range c.Bundles {
		if bundle.StateCommand != "" {
			commands = append(commands, bundle.StateCommand)
		}
	}
	return commands
}

func (c *ToolConfig) Commands() []Command {
	if c.commands != nil {
		return c.commands
	}
	
	commands := []Command{}
	toolSources := map[string]string{} // Track which file each tool comes from
	
	if c.EnableBashTool {
		commands = append(commands, BashCommand)
		toolSources[BashCommand.Name] = "<builtin>"
	}
	
	for _, bundle := range c.Bundles {
		for _, command := range bundle.Commands {
			if _, exists := toolSources[command.Name]; exists {
				existingSource := toolSources[command.Name]
				msg := fmt.Sprintf(
					"Tool '%s' is defined multiple times:\n"+
					"  - First definition in: %s\n"+
					"  - Duplicate definition in: %s",
					command.Name, existingSource, bundle.Path,
				)
				panic(msg)
			}
			commands = append(commands, command)
			toolSources[command.Name] = bundle.Path
		}
	}
	
	c.commands = commands
	return commands
}

func (c *ToolConfig) Tools() []map[string]interface{} {
	if c.tools != nil {
		return c.tools
	}
	
	tools := []map[string]interface{}{}
	for _, command := range c.Commands() {
		tools = append(tools, command.GetFunctionCallingTool())
	}
	
	c.tools = tools
	return tools
}

func (c *ToolConfig) Initialize() {
	commands := c.Commands()
	
	multiLineCommandEndings := map[string]string{}
	for _, command := range commands {
		if command.EndName != "" {
			multiLineCommandEndings[command.Name] = command.EndName
		}
	}
	c.MultiLineCommandEndings = multiLineCommandEndings
	
	c.CommandDocs = GenerateCommandDocs(commands, []string{}, c.EnvVariables)
	
	for _, command := range commands {
		if command.Name == c.SubmitCommand {
			c.SubmitCommandEndName = command.EndName
			break
		}
	}
}

type Bundle struct {
	Path         string    `json:"path"`
	Commands     []Command `json:"commands"`
	StateCommand string    `json:"state_command,omitempty"`
}

type Command struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	EndName     string                 `json:"end_name,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

func (c Command) GetFunctionCallingTool() map[string]interface{} {
	tool := map[string]interface{}{
		"type":        "function",
		"function":    map[string]interface{}{
			"name":        c.Name,
			"description": c.Description,
		},
	}
	
	if c.Schema != nil {
		tool["function"].(map[string]interface{})["parameters"] = c.Schema
	}
	
	return tool
}

var BashCommand = Command{
	Name:        "bash",
	Description: "Execute a bash command",
	Schema: map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "The bash command to execute",
			},
		},
		"required": []string{"command"},
	},
}

type ToolHandler struct {
	Config           *ToolConfig
	ResetCommands    []string
	CommandPatterns  map[string]*regexp.Regexp
	Logger           tools.Logger
	MockState        map[string]string
}

func NewToolHandler(config *ToolConfig, logger tools.Logger) *ToolHandler {
	if logger == nil {
		logger = &tools.DefaultLogger{}
	}
	
	configCopy := *config
	
	handler := &ToolHandler{
		Config:          &configCopy,
		ResetCommands:   []string{},
		CommandPatterns: map[string]*regexp.Regexp{},
		Logger:          logger,
	}
	
	handler.CommandPatterns = handler.getCommandPatterns()
	
	return handler
}

func (h *ToolHandler) Install(env Environment) error {
	if err := h.installCommands(env); err != nil {
		return err
	}
	
	return h.Reset(env)
}

func (h *ToolHandler) Reset(env Environment) error {
	h.Logger.Info("Resetting tools")
	
	if err := env.SetEnvVariables(h.Config.EnvVariables); err != nil {
		return err
	}
	
	registryVars, err := json.Marshal(h.Config.RegistryVariables)
	if err != nil {
		return fmt.Errorf("failed to marshal registry variables: %w", err)
	}
	
	if err := env.WriteFile("/root/.swe-agent-env", string(registryVars)); err != nil {
		return err
	}
	
	if err := env.WriteFile("/root/state.json", "{}"); err != nil {
		return err
	}
	
	if len(h.ResetCommands) > 0 {
		_, err := env.Communicate(strings.Join(h.ResetCommands, " && "), h.Config.InstallTimeout, "raise", "Failed to reset tools")
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (h *ToolHandler) GetState(env Environment) (map[string]string, error) {
	if h.MockState != nil {
		return h.MockState, nil
	}
	
	for _, stateCommand := range h.Config.StateCommands() {
		_, err := env.Communicate(stateCommand, h.Config.ExecutionTimeout, "warn", "Failed to execute state command")
		if err != nil {
			h.Logger.Warn("Failed to execute state command: %s", err)
		}
	}
	
	state, err := h.getState(env)
	if err != nil {
		return nil, err
	}
	
	h.Logger.Debug("Retrieved state from environment: %v", state)
	return state, nil
}

func (h *ToolHandler) getState(env Environment) (map[string]string, error) {
	stateStr, err := env.ReadFile("/root/state.json", "", "")
	if err != nil {
		h.Logger.Warn("State file not found, returning empty state")
		return map[string]string{}, nil
	}
	
	if strings.TrimSpace(stateStr) == "" {
		h.Logger.Warn("State file is empty, returning empty state")
		return map[string]string{}, nil
	}
	
	var state map[string]string
	if err := json.Unmarshal([]byte(stateStr), &state); err != nil {
		return nil, fmt.Errorf("state %q is not valid json: %w", stateStr, err)
	}
	
	return state, nil
}

func (h *ToolHandler) ShouldBlockAction(action string) bool {
	action = strings.TrimSpace(action)
	if action == "" {
		return false
	}
	
	for _, blockedPrefix := range h.Config.Filter.Blocklist {
		if strings.HasPrefix(action, blockedPrefix) {
			return true
		}
	}
	
	for _, blockedCommand := range h.Config.Filter.BlocklistStandalone {
		if action == blockedCommand {
			return true
		}
	}
	
	name := strings.Split(action, " ")[0]
	if regex, exists := h.Config.Filter.BlockUnlessRegex[name]; exists {
		matched, _ := regexp.MatchString(regex, action)
		if !matched {
			return true
		}
	}
	
	return false
}

func (h *ToolHandler) CheckForSubmissionCmd(output string) bool {
	return strings.Contains(output, "<<SWE_AGENT_SUBMISSION>>")
}

func (h *ToolHandler) ParseActions(output map[string]interface{}) (string, string, error) {
	thought := ""
	action := ""
	
	if thoughtVal, exists := output["thought"]; exists {
		if thoughtStr, ok := thoughtVal.(string); ok {
			thought = thoughtStr
		}
	}
	
	if actionVal, exists := output["action"]; exists {
		if actionStr, ok := actionVal.(string); ok {
			action = actionStr
		}
	}
	
	return thought, action, nil
}

func (h *ToolHandler) GuardMultilineInput(action string) string {
	return GuardMultilineInput(action, h.getFirstMultilineCmd)
}

func (h *ToolHandler) getFirstMultilineCmd(action string) *regexp.Regexp {
	patterns := map[string]*regexp.Regexp{}
	for k, v := range h.CommandPatterns {
		if _, exists := h.Config.MultiLineCommandEndings[k]; exists || k == h.Config.SubmitCommand {
			patterns[k] = v
		}
	}
	
	var matches []*regexp.Regexp
	for _, pat := range patterns {
		if pat.MatchString(action) {
			matches = append(matches, pat)
		}
	}
	
	if len(matches) == 0 {
		return nil
	}
	
	return matches[0]
}

func (h *ToolHandler) getCommandPatterns() map[string]*regexp.Regexp {
	patterns := map[string]*regexp.Regexp{}
	
	for _, command := range h.Config.Commands() {
		var pat *regexp.Regexp
		if command.EndName != "" {
			pat = regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(.*?)^(%s)\s*$`, regexp.QuoteMeta(command.Name), regexp.QuoteMeta(command.EndName)))
		} else {
			pat = regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(.*?)$`, regexp.QuoteMeta(command.Name)))
		}
		patterns[command.Name] = pat
	}
	
	submitPat := regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(.*?)^(%s)\s*$`, regexp.QuoteMeta(h.Config.SubmitCommand), regexp.QuoteMeta(h.Config.SubmitCommandEndName)))
	patterns[h.Config.SubmitCommand] = submitPat
	
	return patterns
}

func (h *ToolHandler) installCommands(env Environment) error {
	if err := env.SetEnvVariables(h.Config.EnvVariables); err != nil {
		return err
	}
	
	cwd, err := env.Communicate("pwd", 10, "raise", "Failed to get current directory")
	if err != nil {
		return err
	}
	cwd = strings.TrimSpace(cwd)
	
	
	for _, bundle := range h.Config.Bundles {
		cmds := []string{
			fmt.Sprintf("export PATH=/root/tools/%s/bin:$PATH", bundle.Path),
			fmt.Sprintf("chmod +x /root/tools/%s/bin/*", bundle.Path),
		}
		
		cmds = append(cmds, fmt.Sprintf("cd /root/tools/%s && source install.sh", bundle.Path))
		cmds = append(cmds, fmt.Sprintf("chmod +x /root/tools/%s/bin/*", bundle.Path))
		
		_, err := env.Communicate(strings.Join(cmds, " && "), h.Config.InstallTimeout, "raise", "Failed to install tools")
		if err != nil {
			return err
		}
	}
	
	_, err = env.Communicate(fmt.Sprintf("cd %s", cwd), 10, "raise", "Failed to change directory")
	if err != nil {
		return err
	}
	
	path, err := env.Communicate("echo $PATH", 10, "raise", "Failed to get PATH")
	if err != nil {
		return err
	}
	path = strings.TrimSpace(path)
	
	
	return nil
}

type Environment interface {
	Communicate(input string, timeout int, check string, errorMsg string) (string, error)
	ReadFile(path string, encoding string, errors string) (string, error)
	WriteFile(path string, content string) error
	SetEnvVariables(envVariables map[string]interface{}) error
}

func GuardMultilineInput(action string, getFirstMultilineCmd func(string) *regexp.Regexp) string {
	return action
}

func GenerateCommandDocs(commands []Command, additionalDocs []string, envVariables map[string]interface{}) string {
	var docs strings.Builder
	
	docs.WriteString("# Available Commands\n\n")
	
	for _, command := range commands {
		docs.WriteString(fmt.Sprintf("## %s\n\n", command.Name))
		docs.WriteString(fmt.Sprintf("%s\n\n", command.Description))
		
		if command.Schema != nil {
			docs.WriteString("### Parameters\n\n")
			
			if properties, ok := command.Schema["properties"].(map[string]interface{}); ok {
				for name, prop := range properties {
					if propMap, ok := prop.(map[string]interface{}); ok {
						docs.WriteString(fmt.Sprintf("- `%s`: ", name))
						
						if description, ok := propMap["description"].(string); ok {
							docs.WriteString(description)
						}
						
						if propType, ok := propMap["type"].(string); ok {
							docs.WriteString(fmt.Sprintf(" (type: %s)", propType))
						}
						
						docs.WriteString("\n")
					}
				}
			}
			
			docs.WriteString("\n")
		}
	}
	
	for _, doc := range additionalDocs {
		docs.WriteString(doc)
		docs.WriteString("\n\n")
	}
	
	if len(envVariables) > 0 {
		docs.WriteString("# Environment Variables\n\n")
		
		for name, value := range envVariables {
			docs.WriteString(fmt.Sprintf("- `%s`: %v\n", name, value))
		}
	}
	
	return docs.String()
}
