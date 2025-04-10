package agent_test

import (
	"context"
	"testing"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgent(t *testing.T) {
	config := agent.AgentConfig{
		Name:        "test-agent",
		Description: "Test agent",
		ToolsPath:   "../../pkg/tools/tools.json",
		PromptsPath: "../../pkg/prompts/prompts.txt",
		ModulesPath: "../../pkg/modules",
		MCPServers:  []string{},
	}
	
	agent, err := agent.NewAgent(config)
	require.NoError(t, err)
	
	t.Run("agent_state", func(t *testing.T) {
		assert.Equal(t, agent.StateIdle, agent.GetState())
		
		ctx := context.Background()
		err := agent.Start(ctx)
		require.NoError(t, err)
		
		assert.Equal(t, agent.StateRunning, agent.GetState())
		
		err = agent.Stop()
		require.NoError(t, err)
		
		assert.Equal(t, agent.StateIdle, agent.GetState())
	})
	
	t.Run("agent_context", func(t *testing.T) {
		agent.SetContext(map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		})
		
		context := agent.GetContext()
		assert.Equal(t, "value1", context["key1"])
		assert.Equal(t, 42, context["key2"])
		
		agent.UpdateContext(map[string]interface{}{
			"key2": 43,
			"key3": true,
		})
		
		context = agent.GetContext()
		assert.Equal(t, "value1", context["key1"])
		assert.Equal(t, 43, context["key2"])
		assert.Equal(t, true, context["key3"])
	})
	
	t.Run("tool_registration_and_execution", func(t *testing.T) {
		testTool := tools.NewBaseTool(
			"test-tool",
			"Test tool",
			tools.FileOperations,
			tools.CoreTool,
			[]tools.ToolParameter{
				{
					Name:        "param1",
					Type:        "string",
					Description: "Parameter 1",
					Required:    true,
				},
			},
			func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
				param1, _ := params["param1"].(string)
				return "Result: " + param1, nil
			},
		)
		
		err := agent.RegisterTool(testTool)
		require.NoError(t, err)
		
		tool, err := agent.GetTool("test-tool")
		require.NoError(t, err)
		assert.Equal(t, "test-tool", tool.Name())
		
		ctx := context.Background()
		result, err := agent.ExecuteTool(ctx, "test-tool", map[string]interface{}{
			"param1": "test",
		})
		require.NoError(t, err)
		assert.Equal(t, "Result: test", result)
	})
}

func TestLoop(t *testing.T) {
	config := agent.AgentConfig{
		Name:        "test-agent",
		Description: "Test agent",
		ToolsPath:   "../../pkg/tools/tools.json",
		PromptsPath: "../../pkg/prompts/prompts.txt",
		ModulesPath: "../../pkg/modules",
		MCPServers:  []string{},
	}
	
	agent, err := agent.NewAgent(config)
	require.NoError(t, err)
	
	loop := agent.Loop
	
	t.Run("loop_state", func(t *testing.T) {
		ctx := context.Background()
		err := loop.Start(ctx)
		require.NoError(t, err)
		
		time.Sleep(100 * time.Millisecond)
		
		err = loop.Stop()
		require.NoError(t, err)
	})
}
