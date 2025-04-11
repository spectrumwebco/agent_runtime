package python

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type ScriptRunner struct {
	interpreter    *Interpreter
	toolsPath      string
	pythonPath     string
	agentPath      string
	rocketProducer rocketmq.Producer
	eventTopic     string
	moduleCache    map[string]bool
}

func NewScriptRunner(toolsPath string, rocketMQAddrs []string, eventTopic string) (*ScriptRunner, error) {
	interpreter, err := NewInterpreter()
	if err != nil {
		return nil, fmt.Errorf("failed to create Python interpreter: %w", err)
	}

	if _, err := os.Stat(toolsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tools path %s does not exist", toolsPath)
	}

	agentPath := filepath.Dir(toolsPath)
	
	runner := &ScriptRunner{
		interpreter: interpreter,
		toolsPath:   toolsPath,
		pythonPath:  os.Getenv("PYTHONPATH"),
		agentPath:   agentPath,
		eventTopic:  eventTopic,
		moduleCache: make(map[string]bool),
	}
	
	if len(rocketMQAddrs) > 0 && eventTopic != "" {
		log.Printf("Initializing RocketMQ producer with addresses: %v", rocketMQAddrs)
		p, err := rocketmq.NewProducer(
			producer.WithNameServer(rocketMQAddrs),
			producer.WithRetry(3),
			producer.WithGroupName("agent-runtime-python-ffi"),
		)
		if err != nil {
			log.Printf("Warning: Failed to create RocketMQ producer: %v", err)
		} else {
			if err := p.Start(); err != nil {
				log.Printf("Warning: Failed to start RocketMQ producer: %v", err)
			} else {
				runner.rocketProducer = p
				log.Println("RocketMQ producer started successfully")
			}
		}
	}

	return runner, nil
}

func (r *ScriptRunner) RunToolScript(toolName string, args map[string]interface{}) (interface{}, error) {
	log.Printf("Running tool script: %s", toolName)
	
	startTime := time.Now()
	
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
			log.Printf("Found tool script at: %s", toolScript)
			break
		}
	}

	if toolScript == "" {
		log.Printf("Tool script not found in standard locations, trying module approach")
		return r.runToolModule(toolName, args)
	}

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	output, err := r.interpreter.ExecutePythonScript(toolScript, []string{string(argsJSON)})
	if err != nil {
		r.publishToolEvent(toolName, "error", map[string]interface{}{
			"error":     err.Error(),
			"tool_name": toolName,
			"duration":  time.Since(startTime).Milliseconds(),
		})
		return nil, fmt.Errorf("failed to execute tool script: %w", err)
	}

	result, err := r.parseOutput(output)
	
	r.publishToolEvent(toolName, "success", map[string]interface{}{
		"tool_name": toolName,
		"duration":  time.Since(startTime).Milliseconds(),
	})
	
	return result, err
}

func (r *ScriptRunner) runToolModule(toolName string, args map[string]interface{}) (interface{}, error) {
	log.Printf("Running tool module: %s", toolName)
	startTime := time.Now()
	
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
import traceback
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
        print(json.dumps({"error": "Failed to import required modules", "traceback": traceback.format_exc()}))
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
        print(json.dumps({
            "error": str(e),
            "traceback": traceback.format_exc()
        }))

if __name__ == "__main__":
    main()
`, r.agentPath, argsJSON, toolName)

	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		r.publishToolEvent(toolName, "error", map[string]interface{}{
			"error":     fmt.Sprintf("failed to write temp script: %v", err),
			"tool_name": toolName,
			"duration":  time.Since(startTime).Milliseconds(),
		})
		return nil, fmt.Errorf("failed to write temp script: %w", err)
	}

	output, err := r.interpreter.ExecutePythonScript(scriptPath, []string{})
	if err != nil {
		r.publishToolEvent(toolName, "error", map[string]interface{}{
			"error":     fmt.Sprintf("failed to execute tool module: %v", err),
			"tool_name": toolName,
			"duration":  time.Since(startTime).Milliseconds(),
		})
		return nil, fmt.Errorf("failed to execute tool module: %w", err)
	}

	result, err := r.parseOutput(output)
	
	if err != nil {
		r.publishToolEvent(toolName, "error", map[string]interface{}{
			"error":     err.Error(),
			"tool_name": toolName,
			"duration":  time.Since(startTime).Milliseconds(),
		})
	} else {
		r.publishToolEvent(toolName, "success", map[string]interface{}{
			"tool_name": toolName,
			"duration":  time.Since(startTime).Milliseconds(),
		})
	}
	
	return result, err
}

func (r *ScriptRunner) parseOutput(output string) (interface{}, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err == nil {
		if errMsg, ok := result["error"]; ok {
			errStr := fmt.Sprintf("%v", errMsg)
			if traceback, ok := result["traceback"]; ok {
				log.Printf("Tool execution error traceback: %v", traceback)
			}
			return nil, fmt.Errorf("tool execution error: %v", errStr)
		}
		if res, ok := result["result"]; ok {
			return res, nil
		}
		return result, nil
	}

	return output, nil
}

func (r *ScriptRunner) publishToolEvent(toolName, eventType string, data map[string]interface{}) {
	if r.rocketProducer == nil || r.eventTopic == "" {
		return // RocketMQ not configured
	}
	
	data["timestamp"] = time.Now().UnixNano() / int64(time.Millisecond)
	data["event_type"] = eventType
	data["source"] = "python-ffi"
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling tool event data: %v", err)
		return
	}
	
	msg := &primitive.Message{
		Topic: r.eventTopic,
		Body:  jsonData,
	}
	
	msg.WithTag(fmt.Sprintf("tool-%s", toolName))
	msg.WithKeys([]string{toolName, eventType})
	
	r.rocketProducer.SendAsync(context.Background(), 
		func(ctx context.Context, result *primitive.SendResult, err error) {
			if err != nil {
				log.Printf("Error sending tool event to RocketMQ: %v", err)
			} else {
				log.Printf("Tool event sent to RocketMQ: %s, MessageID: %s", 
					toolName, result.MsgID)
			}
		}, msg)
}

func (r *ScriptRunner) Close() error {
	if r.rocketProducer != nil {
		if err := r.rocketProducer.Shutdown(); err != nil {
			log.Printf("Error shutting down RocketMQ producer: %v", err)
		}
	}
	
	return r.interpreter.Close()
}
