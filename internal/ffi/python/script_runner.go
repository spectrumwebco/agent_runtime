package python

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ScriptRunner struct {
	interpreter *Interpreter
	toolsPath   string
	pythonPath  string
	agentPath   string
}

func NewScriptRunner(toolsPath string) (*ScriptRunner, error) {
	interpreter, err := NewInterpreter()
	if err != nil {
		return nil, fmt.Errorf("failed to create Python interpreter: %w", err)
	}

	if _, err := os.Stat(toolsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tools path %s does not exist", toolsPath)
	}

	agentPath := filepath.Dir(toolsPath)

	return &ScriptRunner{
		interpreter: interpreter,
		toolsPath:   toolsPath,
		pythonPath:  os.Getenv("PYTHONPATH"),
		agentPath:   agentPath,
	}, nil
}

func (r *ScriptRunner) RunToolScript(toolName string, args map[string]interface{}) (interface{}, error) {
	toolLocations := []string{
		filepath.Join(r.toolsPath, "defaults", "bin", toolName),
		filepath.Join(r.toolsPath, toolName, "bin", toolName),
		filepath.Join(r.toolsPath, "registry", "bin", toolName),
		filepath.Join(r.toolsPath, "defaults", "lib", toolName + ".py"),
	}

	var toolScript string
	for _, loc := range toolLocations {
		if _, err := os.Stat(loc); err == nil {
			toolScript = loc
			break
		}
	}

	if toolScript == "" {
		return r.runToolModule(toolName, args)
	}

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	output, err := r.interpreter.ExecutePythonScript(toolScript, []string{string(argsJSON)})
	if err != nil {
		return nil, fmt.Errorf("failed to execute tool script: %w", err)
	}

	return r.parseOutput(output)
}

func (r *ScriptRunner) runToolModule(toolName string, args map[string]interface{}) (interface{}, error) {
	tempDir, err := os.MkdirTemp("", "agent-runtime-tool")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	scriptPath := filepath.Join(tempDir, "run_tool.py")
	script := fmt.Sprintf(`
import json
import sys
import os
from pathlib import Path

# Add agent path to Python path
agent_path = "%s"
if agent_path not in sys.path:
    sys.path.insert(0, agent_path)

try:
    from agent.tools.tools import ToolConfig, ToolHandler
    from agent.environment.swe_env import SWEEnv
except ImportError:
    # Try with sweagent prefix if agent import fails
    try:
        from sweagent.tools.tools import ToolConfig, ToolHandler
        from sweagent.environment.swe_env import SWEEnv
    except ImportError:
        print(json.dumps({"error": "Failed to import required modules"}))
        sys.exit(1)

def main():
    args = json.loads('''%s''')
    tool_name = "%s"
    
    # Initialize environment and tool handler
    env = SWEEnv()
    config = ToolConfig()
    handler = ToolHandler(config)
    
    # Find the tool
    tool = None
    for cmd in handler.config.commands:
        if cmd.name == tool_name:
            tool = cmd
            break
    
    if tool is None:
        print(json.dumps({"error": f"Tool {tool_name} not found"}))
        return
    
    # Execute the tool
    try:
        result = tool.execute(env, **args)
        print(json.dumps({"result": result}))
    except Exception as e:
        print(json.dumps({"error": str(e)}))

if __name__ == "__main__":
    main()
`, r.agentPath, argsJSON, toolName)

	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp script: %w", err)
	}

	output, err := r.interpreter.ExecutePythonScript(scriptPath, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to execute tool module: %w", err)
	}

	return r.parseOutput(output)
}

func (r *ScriptRunner) parseOutput(output string) (interface{}, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err == nil {
		if errMsg, ok := result["error"]; ok {
			return nil, fmt.Errorf("tool execution error: %v", errMsg)
		}
		if res, ok := result["result"]; ok {
			return res, nil
		}
		return result, nil
	}

	return output, nil
}

func (r *ScriptRunner) Close() error {
	return r.interpreter.Close()
}
