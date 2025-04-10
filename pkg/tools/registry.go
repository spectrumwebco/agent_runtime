package tools

import (
	"fmt"
	"strings"
	"sync"
	"time" // Added import

	"context" // Added for FFI execution context
	"strings" // Added for argument parsing

	"github.com/spectrumwebco/agent_runtime/internal/config" // Assuming config is internal
	"github.com/spectrumwebco/agent_runtime/internal/env"   // Assuming env is internal
	"github.com/spectrumwebco/agent_runtime/internal/ffi/cpp"    // Added C++ FFI import
	"github.com/spectrumwebco/agent_runtime/internal/ffi/python" // Added Python FFI import
)

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) // Observation, error
}

type Registry struct {
	config *config.Config
	tools  map[string]Tool
	mutex  sync.RWMutex
	pythonInterpreter *python.Interpreter
	cppInterpreter    *cpp.Interpreter
}

type Handler struct {
	Config      ToolConfig // Placeholder for tool configuration
	Definitions map[string]ToolDefinition
	Registry    *Registry // Reference to the tool registry
}

type ToolConfig struct {
	ExecutionTimeout              time.Duration
	MaxConsecutiveExecutionTimeouts int
	TotalExecutionTimeout         time.Duration
	FormatErrorTemplate           string
	CppInterpreterConfig          cpp.InterpreterConfig // Added C++ config
}

type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"` // JSON schema for input
}

func NewRegistry(cfg *config.Config, toolCfg ToolConfig) (*Registry, error) {
	pyInterpreter, err := python.NewInterpreter()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize Python FFI interpreter: %v\n", err)
	}

	cppInterpreter, err := cpp.NewInterpreter(toolCfg.CppInterpreterConfig)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize C++ FFI interpreter: %v\n", err)
	}

	registry := &Registry{
		config:            cfg,
		tools:             make(map[string]Tool),
		pythonInterpreter: pyInterpreter,
		cppInterpreter:    cppInterpreter,
	}
	registry.registerBuiltinTools()
	return registry, nil
}

func (r *Registry) registerBuiltinTools() {
	r.Register(&ShellTool{
		name:        "shell",
		description: "Executes a shell command in the environment",
	})

	r.Register(&FileTool{
		name:        "file",
		description: "Performs file operations (read, write, list)",
	})

	r.Register(&HTTPTool{
		name:        "http",
		description: "Makes HTTP requests",
	})


	r.Register(&EditReplaceTool{
		name:        "edit_replace",
		description: "Replaces occurrences of a string in a file",
	})

	r.Register(&SearchTool{
		name:        "search",
		description: "Searches for a pattern in files",
	})


	r.Register(&EditAnthropicTool{
		name:        "edit_anthropic",
		description: "Placeholder for edit_anthropic tool",
	})
	r.Register(&EditLintingTool{
		name:        "edit_linting",
		description: "Placeholder for edit_linting tool",
	})
	r.Register(&EditRewriteTool{
		name:        "edit_rewrite",
		description: "Placeholder for edit_rewrite tool",
	})
	r.Register(&FilemapTool{
		name:        "filemap",
		description: "Placeholder for filemap tool",
	})
	r.Register(&ForfeitTool{
		name:        "forfeit",
		description: "Placeholder for forfeit tool",
	})
	r.Register(&ReviewOnSubmitTool{
		name:        "review_on_submit",
		description: "Placeholder for review_on_submit tool",
	})
	r.Register(&ReviewOnSubmitMTool{
		name:        "review_on_submit_m",
		description: "Placeholder for review_on_submit_m tool",
	})

	if r.pythonInterpreter != nil {
		r.Register(&PythonTool{
			name:        "python",
			description: "Executes Python code using the FFI interpreter",
			interpreter: r.pythonInterpreter,
		})
	} else {
		fmt.Println("Python interpreter not available, Python tool disabled.")
	}

	if r.cppInterpreter != nil {
		r.Register(&CppTool{
			name:        "cpp",
			description: "Compiles and executes C++ code using the FFI interpreter",
			interpreter: r.cppInterpreter,
		})
	} else {
		fmt.Println("C++ interpreter not available, Cpp tool disabled.")
	}

	r.Register(&SubmitTool{
		name:        "submit",
		description: "Submits the final solution or patch",
	})
}


func NewHandler(config ToolConfig, definitions []ToolDefinition, registry *Registry) (*Handler, error) {
	if registry == nil {
		return nil, fmt.Errorf("registry cannot be nil for Handler")
	}
	defsMap := make(map[string]ToolDefinition)
	for _, def := range definitions {
		defsMap[def.Name] = def
	}
	return &Handler{
		Config:      config,
		Definitions: defsMap,
		Registry:    registry,
	}, nil
}

func (r *Registry) Register(tool Tool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tools[tool.Name()]; exists {
		return fmt.Errorf("tool already registered: %s", tool.Name())
	}

	r.tools[tool.Name()] = tool
	fmt.Printf("Registered tool: %s\n", tool.Name())
	return nil
}

func (r *Registry) Get(name string) (Tool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool, nil
}

func (h *Handler) ParseActions(output string) (thought string, action string, err error) {
	message := output

	toolBlockMarker := "```tool"
	bashBlockMarker := "```bash" // SWE-Agent often uses bash blocks for shell commands

	if strings.Contains(message, toolBlockMarker) {
		parts := strings.SplitN(message, toolBlockMarker, 2)
		thought = strings.TrimSpace(parts[0])
		actionPart := strings.TrimSpace(parts[1])
		action = strings.TrimSuffix(actionPart, "```")
		action = strings.TrimSpace(action)
	} else if strings.Contains(message, bashBlockMarker) {
		parts := strings.SplitN(message, bashBlockMarker, 2)
		thought = strings.TrimSpace(parts[0])
		actionPart := strings.TrimSpace(parts[1])
		action = strings.TrimSuffix(actionPart, "```")
		action = strings.TrimSpace(action)
		action = "shell " + action
	} else {
		if strings.Contains(strings.ToLower(message), "submit") {
			thought = message
			action = "submit" // Assume submit action if keyword found
		} else {
			thought = message
			action = "" // No specific action identified
		}
	}

	fmt.Printf("Parsed Thought: %s\nParsed Action: %s\n", thought, action)
	return thought, action, nil // Return nil error for now
}

func (h *Handler) ExecuteAction(ctx context.Context, action string, environment *env.SWEEnv) (observation string, err error) {
	fmt.Printf("Executing action: %s\n", action)
	if environment == nil {
		return "", fmt.Errorf("environment is not initialized")
	}
	if h.Registry == nil {
		return "", fmt.Errorf("tool registry is not initialized")
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), h.Config.ExecutionTimeout) // Use configured timeout
		defer cancel()
	}


	parts := strings.Fields(action)
	if len(parts) == 0 {
		return "No action specified.", nil // Or error?
	}
	toolName := parts[0]
	argsStr := ""
	if len(parts) > 1 {
		argsStr = strings.Join(parts[1:], " ")
	}

	tool, err := h.Registry.Get(toolName)
	if err != nil {
		fmt.Printf("Tool '%s' not found, attempting shell execution.\n", toolName)
		shellTool, shellErr := h.Registry.Get("shell")
		if shellErr != nil {
			return "", fmt.Errorf("tool '%s' not found and shell tool is unavailable: %w", toolName, shellErr)
		}
		argsMap := map[string]interface{}{"command": action}
		return shellTool.Execute(ctx, argsMap, environment)
	}

	argsMap := make(map[string]interface{})

	switch toolName {
	case "submit":
		argsMap["submission"] = argsStr // Assumes submission text follows the command
	case "python", "cpp":
		argsMap["code"] = argsStr // Assumes code follows the command
	case "shell":
		argsMap["command"] = argsStr // The rest of the string is the command
	case "file":
		parts := strings.Fields(argsStr)
		if len(parts) >= 2 {
			argsMap["operation"] = parts[0]
			argsMap["path"] = parts[1]
			if parts[0] == "write" && len(parts) > 2 {
				argsMap["content"] = strings.Join(parts[2:], " ")
			}
		} else {
			fmt.Printf("Warning: Could not parse args for file tool: %s\n", argsStr)
		}
	case "edit_replace":
		parts := strings.Fields(argsStr) // Basic split, might fail with quoted strings
		if len(parts) >= 3 {
			argsMap["path"] = parts[0]
			argsMap["old_string"] = parts[1]
			argsMap["new_string"] = strings.Join(parts[2:], " ") // Assume new_string is the rest
		} else {
			fmt.Printf("Warning: Could not parse args for edit_replace tool: %s\n", argsStr)
		}
	case "search":
		parts := strings.Fields(argsStr)
		if len(parts) >= 1 {
			argsMap["query"] = parts[0] // Basic, assumes query is first word
			if len(parts) > 1 {
				argsMap["path"] = strings.Join(parts[1:], " ") // Assume path is the rest
			}
		} else {
			fmt.Printf("Warning: Could not parse args for search tool: %s\n", argsStr)
		}
	case "http":
		parts := strings.Fields(argsStr)
		if len(parts) >= 2 {
			argsMap["method"] = parts[0]
			argsMap["url"] = parts[1]
			if len(parts) > 2 {
				argsMap["body"] = strings.Join(parts[2:], " ")
			}
		} else {
			fmt.Printf("Warning: Could not parse args for http tool: %s\n", argsStr)
		}
	default:
		argsMap["raw_args"] = argsStr
		fmt.Printf("Using default raw_args parsing for tool %s: %s\n", toolName, argsStr)

	}

	return tool.Execute(ctx, argsMap, environment)
}

func (h *Handler) GetState(environment *env.SWEEnv) (map[string]interface{}, error) {
	fmt.Println("Getting environment state...")
	if environment == nil {
		return nil, fmt.Errorf("environment is not initialized")
	}
	state := make(map[string]interface{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Short timeout for state gathering
	defer cancel()

	cwd, err := environment.CommunicateWithContext(ctx, "pwd", "")
	if err == nil {
		state["cwd"] = strings.TrimSpace(cwd)
	} else {
		fmt.Printf("Warning: Failed to get cwd for state: %v\n", err)
	}

	ls, err := environment.CommunicateWithContext(ctx, "ls -la", "")
	if err == nil {
		state["ls"] = ls
	} else {
		fmt.Printf("Warning: Failed to get ls for state: %v\n", err)
	}

	state["open_file"] = "placeholder/file.py"
	state["open_file_content"] = "Placeholder file content..."
	return state, nil
}

func (h *Handler) Install(environment *env.SWEEnv) error {
	fmt.Println("Installing tools in environment...")
	if environment == nil {
		return fmt.Errorf("environment is not initialized")
	}
	return nil
}


type ShellTool struct {
	name        string
	description string
}

func (t *ShellTool) Name() string        { return t.name }
func (t *ShellTool) Description() string { return t.description }
func (t *ShellTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	command, ok := args["command"].(string)
	if !ok || command == "" {
		return "Error: Empty command provided to shell tool.", nil
		// return "", fmt.Errorf("invalid or missing 'command' argument for shell tool")
	}
	fmt.Printf("Executing shell command: %s\n", command)
	output, err := environment.CommunicateWithContext(ctx, command, "")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Sprintf("Error: Command '%s' timed out.", command), nil
		}
		return fmt.Sprintf("Error executing command '%s': %v\nOutput:\n%s", command, err, output), nil
	}
	return "Output:\n" + output, nil
}


type FileTool struct {
	name        string
	description string
}

func (t *FileTool) Name() string        { return t.name }
func (t *FileTool) Description() string { return t.description }
func (t *FileTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	operation, opOK := args["operation"].(string)
	path, pathOK := args["path"].(string)
	if !opOK || !pathOK {
		return "Error: Invalid arguments for file tool. Requires 'operation' and 'path'.", nil
	}

	switch operation {
	case "read":
		fmt.Printf("Reading file: %s\n", path)
		content, err := environment.ReadFile(ctx, path, "utf-8", "ignore")
		if err != nil {
			return fmt.Sprintf("Error reading file %s: %v", path, err), nil
		}
		return fmt.Sprintf("File %s content:\n%s", path, content), nil
	case "write":
		content, contentOK := args["content"].(string)
		if !contentOK {
			return "Error: Missing 'content' for file write operation.", nil
		}
		fmt.Printf("Writing to file: %s\n", path)
		err := environment.WriteFile(ctx, path, content)
		if err != nil {
			return fmt.Sprintf("Error writing to file %s: %v", path, err), nil
		}
		return fmt.Sprintf("Successfully wrote to file %s", path), nil
	default:
		return fmt.Sprintf("Error: Unsupported file operation '%s'. Supported: read, write", operation), nil
	}
}


type EditReplaceTool struct {
	name        string
	description string
}

func (t *EditReplaceTool) Name() string        { return t.name }
func (t *EditReplaceTool) Description() string { return t.description }
func (t *EditReplaceTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	path, pathOK := args["path"].(string)
	oldStr, oldOK := args["old_string"].(string)
	newStr, newOK := args["new_string"].(string)
	if !pathOK || !oldOK || !newOK {
		return "Error: Invalid arguments for edit_replace tool. Requires 'path', 'old_string', and 'new_string'.", nil
	}
	fmt.Printf("Executing edit_replace: Replace '%s' with '%s' in file %s\n", oldStr, newStr, path)

	currentContent, err := environment.ReadFile(ctx, path, "utf-8", "ignore")
	if err != nil {
		return fmt.Sprintf("Error reading file %s for edit_replace: %v", path, err), nil
	}

	if !strings.Contains(currentContent, oldStr) {
		return fmt.Sprintf("Error: Pattern '%s' not found in file %s.", oldStr, path), nil
	}
	newContent := strings.ReplaceAll(currentContent, oldStr, newStr)

	err = environment.WriteFile(ctx, path, newContent)
	if err != nil {
		return fmt.Sprintf("Error writing modified content to file %s: %v", path, err), nil
	}

	return fmt.Sprintf("Successfully replaced '%s' with '%s' in file %s", oldStr, newStr, path), nil
}


type SearchTool struct {
	name        string
	description string
}

func (t *SearchTool) Name() string        { return t.name }
func (t *SearchTool) Description() string { return t.description }
func (t *SearchTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	query, queryOK := args["query"].(string)
	if !queryOK {
		return "Error: Invalid arguments for search tool. Requires 'query'.", nil
	}
	path, _ := args["path"].(string) // Optional path

	searchCmd := fmt.Sprintf("grep -rn '%s'", query) // Basic grep example
	if path != "" {
		searchCmd += " " + path
	}

	fmt.Printf("Executing search command: %s\n", searchCmd)
	output, err := environment.CommunicateWithContext(ctx, searchCmd, "")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Sprintf("Error: Search command '%s' timed out.", searchCmd), nil
		}
		return fmt.Sprintf("Error executing search command '%s': %v\nOutput:\n%s", searchCmd, err, output), nil
	}
	if output == "" {
		return fmt.Sprintf("No results found for query '%s'.", query), nil
	}
	return "Output:\n" + output, nil
}


type EditAnthropicTool struct {
	name        string
	description string
}

func (t *EditAnthropicTool) Name() string        { return t.name }
func (t *EditAnthropicTool) Description() string { return t.description }
func (t *EditAnthropicTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing edit_anthropic (placeholder): %+v\n", args)
	return "Placeholder: Edit Anthropic tool executed.", nil
}


type EditLintingTool struct {
	name        string
	description string
}

func (t *EditLintingTool) Name() string        { return t.name }
func (t *EditLintingTool) Description() string { return t.description }
func (t *EditLintingTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing edit_linting (placeholder): %+v\n", args)
	return "Placeholder: Edit Linting tool executed.", nil
}


type EditRewriteTool struct {
	name        string
	description string
}

func (t *EditRewriteTool) Name() string        { return t.name }
func (t *EditRewriteTool) Description() string { return t.description }
func (t *EditRewriteTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing edit_rewrite (placeholder): %+v\n", args)
	return "Placeholder: Edit Rewrite tool executed.", nil
}


type FilemapTool struct {
	name        string
	description string
}

func (t *FilemapTool) Name() string        { return t.name }
func (t *FilemapTool) Description() string { return t.description }
func (t *FilemapTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing filemap (placeholder): %+v\n", args)
	return "Placeholder: Filemap tool executed.", nil
}


type ForfeitTool struct {
	name        string
	description string
}

func (t *ForfeitTool) Name() string        { return t.name }
func (t *ForfeitTool) Description() string { return t.description }
func (t *ForfeitTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing forfeit (placeholder): %+v\n", args)
	return "Placeholder: Forfeit tool executed.", nil
}


type ReviewOnSubmitTool struct {
	name        string
	description string
}

func (t *ReviewOnSubmitTool) Name() string        { return t.name }
func (t *ReviewOnSubmitTool) Description() string { return t.description }
func (t *ReviewOnSubmitTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing review_on_submit (placeholder): %+v\n", args)
	return "Placeholder: Review On Submit tool executed.", nil
}


type ReviewOnSubmitMTool struct {
	name        string
	description string
}

func (t *ReviewOnSubmitMTool) Name() string        { return t.name }
func (t *ReviewOnSubmitMTool) Description() string { return t.description }
func (t *ReviewOnSubmitMTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing review_on_submit_m (placeholder): %+v\n", args)
	return "Placeholder: Review On Submit M tool executed.", nil
}


type HTTPTool struct {
	name        string
	description string
}

func (t *HTTPTool) Name() string        { return t.name }
func (t *HTTPTool) Description() string { return t.description }
func (t *HTTPTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	fmt.Printf("Executing HTTP request (placeholder): %+v\n", args)
	return "HTTP request observation placeholder", nil
}


type PythonTool struct {
	name        string
	description string
	interpreter *python.Interpreter
}

func (t *PythonTool) Name() string        { return t.name }
func (t *PythonTool) Description() string { return t.description }
func (t *PythonTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	if t.interpreter == nil {
		return "Error: Python interpreter is not available.", nil
	}
	code, ok := args["code"].(string)
	if !ok || code == "" {
		return "Error: Invalid or missing 'code' argument for python tool.", nil
	}
	fmt.Printf("Executing Python code via FFI: %s\n", code)
	output, err := t.interpreter.Exec(ctx, code)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "Error: Python execution timed out.", nil
		}
		return fmt.Sprintf("Error executing Python code: %v\nOutput:\n%s", err, output), nil
	}
	return "Output:\n" + output, nil
}


type CppTool struct {
	name        string
	description string
	interpreter *cpp.Interpreter
}

func (t *CppTool) Name() string        { return t.name }
func (t *CppTool) Description() string { return t.description }
func (t *CppTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	if t.interpreter == nil {
		return "Error: C++ interpreter is not available.", nil
	}
	code, ok := args["code"].(string)
	if !ok || code == "" {
		return "Error: Invalid or missing 'code' argument for cpp tool.", nil
	}
	fmt.Printf("Executing C++ code via FFI: %s\n", code)
	output, err := t.interpreter.Exec(ctx, code)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "Error: C++ execution timed out.", nil
		}
		return fmt.Sprintf("Error executing C++ code: %v\nOutput:\n%s", err, output), nil
	}
	return "Output:\n" + output, nil
}


type SubmitTool struct {
	name        string
	description string
}

func (t *SubmitTool) Name() string        { return t.name }
func (t *SubmitTool) Description() string { return t.description }
func (t *SubmitTool) Execute(ctx context.Context, args map[string]interface{}, environment *env.SWEEnv) (string, error) {
	submission, ok := args["submission"].(string)
	if !ok {
		rawArgs, rawOk := args["raw_args"].(string)
		if rawOk {
			submission = rawArgs
		} else {
			submission = "" // Default to empty submission if not provided or parsing failed
		}
	}
	fmt.Printf("Executing submit with submission: %s\n", submission)

	err := environment.Submit(submission)
	if err != nil {
		return fmt.Sprintf("Error executing submit: %v", err), nil
	}

	return "Submission successful.", nil
}
