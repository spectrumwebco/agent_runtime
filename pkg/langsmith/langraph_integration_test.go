package langsmith

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLangGraphIntegration(t *testing.T) {
	if os.Getenv("CI") != "true" && os.Getenv("TEST_LANGSMITH") != "true" {
		t.Skip("Skipping LangSmith integration test. Set TEST_LANGSMITH=true to run this test.")
	}

	config := &LangSmithConfig{
		Enabled:     true,
		APIKey:      os.Getenv("LANGCHAIN_API_KEY"),
		APIUrl:      getEnvOrDefault("LANGCHAIN_ENDPOINT", "https://api.smith.langchain.com"),
		ProjectName: "test-langraph-integration",
		SelfHosted:  os.Getenv("LANGSMITH_SELF_HOSTED") == "true",
	}

	integration, err := NewLangGraphIntegration(config)
	require.NoError(t, err)
	require.NotNil(t, integration)

	tracer := integration.CreateTracer()
	require.NotNil(t, tracer)

	graph := langraph.NewGraph("Test Graph")
	
	nodeA := langraph.NewNode("NodeA", "test", func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"output": "Hello from NodeA"}, nil
	})
	
	nodeB := langraph.NewNode("NodeB", "test", func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		inputValue, _ := inputs["input"].(string)
		return map[string]interface{}{"output": inputValue + " and NodeB"}, nil
	})
	
	graph.AddNode(nodeA)
	graph.AddNode(nodeB)
	
	graph.AddEdge(nodeA, nodeB, func(outputs map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"input": outputs["output"]}
	})

	graphID, err := integration.RegisterGraph(context.Background(), graph)
	require.NoError(t, err)
	require.NotEmpty(t, graphID)

	executor := langraph.NewExecutor(langraph.WithTracer(tracer))
	
	ctx := context.Background()
	inputs := map[string]interface{}{"input": "Hello, LangSmith!"}
	
	outputs, err := executor.ExecuteGraph(ctx, graph, inputs)
	require.NoError(t, err)
	assert.Contains(t, outputs, "output")
	assert.Contains(t, outputs["output"], "Hello from NodeA and NodeB")

	agentConfig := &langraph.MultiAgentSystemConfig{
		Name: "Test Multi-Agent System",
		Agents: []*langraph.AgentConfig{
			{
				ID:   "agent1",
				Name: "Agent 1",
				Type: "test",
			},
			{
				ID:   "agent2",
				Name: "Agent 2",
				Type: "test",
			},
		},
	}

	system, err := langraph.NewMultiAgentSystem(agentConfig, executor)
	require.NoError(t, err)

	systemID, err := integration.RegisterMultiAgentSystem(context.Background(), system)
	require.NoError(t, err)
	require.NotEmpty(t, systemID)

	taskCtx, taskID, err := integration.TraceAgentTask(
		context.Background(),
		"agent1",
		"task1",
		"test",
		"Test task",
		map[string]interface{}{"input": "Test input"},
	)
	require.NoError(t, err)
	require.NotEmpty(t, taskID)

	actionCtx, actionID, err := integration.TraceAgentAction(
		taskCtx,
		"agent1",
		"test_action",
		map[string]interface{}{"input": "Test action input"},
	)
	require.NoError(t, err)
	require.NotEmpty(t, actionID)

	err = integration.EndAgentAction(
		actionCtx,
		"agent1:test_action:"+actionID,
		map[string]interface{}{"output": "Test action output"},
		nil,
	)
	require.NoError(t, err)

	messageID, err := integration.TraceAgentMessage(
		taskCtx,
		"agent1",
		"agent2",
		"test_message",
		"Hello from Agent 1",
		map[string]interface{}{"priority": "high"},
	)
	require.NoError(t, err)
	require.NotEmpty(t, messageID)

	err = integration.EndAgentTask(
		taskCtx,
		"task1",
		map[string]interface{}{"output": "Test task output"},
		nil,
	)
	require.NoError(t, err)

	err = integration.CreateFeedback(
		context.Background(),
		taskID,
		"correctness",
		1.0,
		"The task was completed correctly",
	)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func TestLangGraphIntegrationWithSelfHosted(t *testing.T) {
	if os.Getenv("TEST_LANGSMITH_SELF_HOSTED") != "true" {
		t.Skip("Skipping self-hosted LangSmith integration test. Set TEST_LANGSMITH_SELF_HOSTED=true to run this test.")
	}

	config := &LangSmithConfig{
		Enabled:     true,
		APIKey:      os.Getenv("LANGSMITH_LICENSE_KEY"),
		APIUrl:      getEnvOrDefault("LANGSMITH_URL", "http://langsmith-api.langsmith.svc.cluster.local:8000"),
		ProjectName: "test-langraph-integration-self-hosted",
		SelfHosted:  true,
		SelfHostedConfig: &SelfHostedConfig{
			URL:           getEnvOrDefault("LANGSMITH_URL", "http://langsmith-api.langsmith.svc.cluster.local:8000"),
			AdminUsername: getEnvOrDefault("LANGSMITH_ADMIN_USERNAME", "admin"),
			AdminPassword: os.Getenv("LANGSMITH_ADMIN_PASSWORD"),
			LicenseKey:    os.Getenv("LANGSMITH_LICENSE_KEY"),
		},
	}

	integration, err := NewLangGraphIntegration(config)
	require.NoError(t, err)
	require.NotNil(t, integration)

	tracer := integration.CreateTracer()
	require.NotNil(t, tracer)

	graph := langraph.NewGraph("Self-Hosted Test Graph")
	
	node := langraph.NewNode("TestNode", "test", func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"output": "Hello from self-hosted LangSmith!"}, nil
	})
	
	graph.AddNode(node)

	graphID, err := integration.RegisterGraph(context.Background(), graph)
	require.NoError(t, err)
	require.NotEmpty(t, graphID)

	executor := langraph.NewExecutor(langraph.WithTracer(tracer))
	
	ctx := context.Background()
	inputs := map[string]interface{}{"input": "Self-hosted test"}
	
	outputs, err := executor.ExecuteGraph(ctx, graph, inputs)
	require.NoError(t, err)
	assert.Contains(t, outputs, "output")
	assert.Contains(t, outputs["output"], "Hello from self-hosted LangSmith!")

	time.Sleep(1 * time.Second)
}
