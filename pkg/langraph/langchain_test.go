package langraph

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

type MockLLM struct{}

func (m *MockLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return "This is a mock response for: " + prompt, nil
}

func (m *MockLLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) (*llms.Generation, error) {
	var generations []schema.Generation
	for _, prompt := range prompts {
		generations = append(generations, schema.Generation{
			Text: "This is a mock response for: " + prompt,
		})
	}
	return &llms.Generation{
		Generations: generations,
	}, nil
}

type MockEventStream struct {
	Events []*models.Event
}

func NewMockEventStream() *MockEventStream {
	return &MockEventStream{
		Events: []*models.Event{},
	}
}

func (m *MockEventStream) AddEvent(event *models.Event) error {
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockEventStream) Subscribe(eventType models.EventType, callback func(*models.Event)) error {
	return nil
}

func (m *MockEventStream) Unsubscribe(eventType models.EventType, callback func(*models.Event)) error {
	return nil
}

func TestLangChainBridge(t *testing.T) {
	llm := &MockLLM{}
	eventStream := NewMockEventStream()

	graph := NewGraph("test-graph", "Test graph for LangChain integration")
	executor := NewExecutor(graph)

	bridge := NewLangChainBridge(graph, executor, llm, eventStream)
	assert.NotNil(t, bridge, "Bridge should not be nil")

	toolNodeID, err := bridge.CreateLangChainToolNode("test-tool", "Test tool", func(ctx context.Context, input string) (string, error) {
		return "Tool result: " + input, nil
	})
	assert.NoError(t, err, "Should create tool node without error")
	assert.NotEmpty(t, toolNodeID, "Tool node ID should not be empty")

	chainNodeID, err := bridge.CreateLangChainChainNode("test-chain", "Test chain", "This is a test prompt with {{.input}}")
	assert.NoError(t, err, "Should create chain node without error")
	assert.NotEmpty(t, chainNodeID, "Chain node ID should not be empty")

	edge, err := bridge.ConnectLangChainNodes(toolNodeID, chainNodeID, "test-edge", "Test edge")
	assert.NoError(t, err, "Should connect nodes without error")
	assert.NotNil(t, edge, "Edge should not be nil")

	ctx := context.Background()
	execution, err := bridge.ExecuteLangChainGraph(ctx, toolNodeID, map[string]interface{}{
		"input": "test input",
	})
	assert.NoError(t, err, "Should execute graph without error")
	assert.NotNil(t, execution, "Execution should not be nil")
}

func TestLangChainGraphBuilder(t *testing.T) {
	llm := &MockLLM{}
	eventStream := NewMockEventStream()

	system := NewMultiAgentSystem("test-system", "Test system for LangChain integration", eventStream)
	assert.NotNil(t, system, "System should not be nil")

	builder := NewLangChainGraphBuilder(system, llm, eventStream)
	assert.NotNil(t, builder, "Builder should not be nil")

	builder.AddAgent("test-agent", "Test agent", "react", nil)
	builder.AddTool("test-tool", "Test tool", func(ctx context.Context, input string) (string, error) {
		return "Tool result: " + input, nil
	})
	builder.AddChain("test-chain", "Test chain", "This is a test prompt with {{.input}}")

	builder.Connect("test-agent", "test-tool", "Agent -> Tool", "Agent to Tool communication")
	builder.Connect("test-tool", "test-chain", "Tool -> Chain", "Tool to Chain communication")

	builtSystem := builder.Build()
	assert.NotNil(t, builtSystem, "Built system should not be nil")
	assert.Equal(t, system, builtSystem, "Built system should be the same as the original system")
}

func TestCreateStandardLangChainAgentSystem(t *testing.T) {
	llm := &MockLLM{}
	eventStream := NewMockEventStream()

	system, builder, err := CreateStandardLangChainAgentSystem("test-system", "Test system for LangChain integration", llm, eventStream)
	assert.NoError(t, err, "Should create standard LangChain agent system without error")
	assert.NotNil(t, system, "System should not be nil")
	assert.NotNil(t, builder, "Builder should not be nil")

	assert.NotEmpty(t, builder.nodes, "Builder should have nodes")
	assert.Contains(t, builder.nodes, "frontend_agent", "Builder should have frontend agent")
	assert.Contains(t, builder.nodes, "app_builder_agent", "Builder should have app builder agent")
	assert.Contains(t, builder.nodes, "codegen_agent", "Builder should have codegen agent")
	assert.Contains(t, builder.nodes, "engineering_agent", "Builder should have engineering agent")
	assert.Contains(t, builder.nodes, "orchestrator_agent", "Builder should have orchestrator agent")
}

func TestIntegrationWithKlusterAI(t *testing.T) {
	t.Skip("Skipping integration test with Kluster.AI as it requires actual API access")


	llmConfig := openai.NewConfig("dummy-api-key")
	llm, err := openai.New(llmConfig)
	assert.NoError(t, err, "Should create LLM without error")

	eventStream := NewMockEventStream()

	system, bridge, err := CreateMultiAgentSystemWithLangChain("kluster-ai-system", "Multi-agent system with Kluster.AI integration", llm, eventStream)
	assert.NoError(t, err, "Should create multi-agent system without error")
	assert.NotNil(t, system, "System should not be nil")
	assert.NotNil(t, bridge, "Bridge should not be nil")

}

func TestDjangoIntegration(t *testing.T) {
	t.Skip("Skipping Django integration test as it requires actual Django backend")


	llm := &MockLLM{}
	eventStream := NewMockEventStream()

	system := NewMultiAgentSystem("django-integration-system", "Multi-agent system with Django integration", eventStream)
	assert.NotNil(t, system, "System should not be nil")

	djangoFactory := NewDjangoAgentFactory("http://localhost:8000", "dummy-api-key", eventStream)
	assert.NotNil(t, djangoFactory, "Django factory should not be nil")

	system.SetDjangoFactory(djangoFactory)

	bridge := NewLangChainBridge(system.Graph, system.Executor, llm, eventStream)
	assert.NotNil(t, bridge, "Bridge should not be nil")

	//
}

func TestToolTriggeringOnContextCreation(t *testing.T) {
	llm := &MockLLM{}
	eventStream := NewMockEventStream()

	system := NewMultiAgentSystem("tool-triggering-system", "Multi-agent system with automatic tool triggering", eventStream)
	assert.NotNil(t, system, "System should not be nil")

	builder := NewLangChainGraphBuilder(system, llm, eventStream)
	assert.NotNil(t, builder, "Builder should not be nil")

	knowledgeToolTriggered := false
	builder.AddTool("knowledge-tool", "Knowledge tool that is triggered automatically on context creation", func(ctx context.Context, input string) (string, error) {
		knowledgeToolTriggered = true
		return "Knowledge about: " + input, nil
	})

	builder.AddAgent("context-creator-agent", "Agent that creates context", "react", nil)

	contextNodeID, err := builder.Bridge.CreateLangChainChainNode("context-node", "Context node that triggers knowledge tool", "Creating context: {{.input}}")
	assert.NoError(t, err, "Should create context node without error")
	assert.NotEmpty(t, contextNodeID, "Context node ID should not be empty")

	system.Graph.AddEdge(builder.nodes["context-creator-agent"], builder.nodes["knowledge-tool"], "Context -> Knowledge", "Automatic triggering of knowledge tool on context creation")

	ctx := context.Background()
	execution, err := builder.Execute(ctx, "context-creator-agent", map[string]interface{}{
		"input": "new context",
	})
	assert.NoError(t, err, "Should execute context creation without error")
	assert.NotNil(t, execution, "Execution should not be nil")

	time.Sleep(100 * time.Millisecond)

	assert.True(t, knowledgeToolTriggered, "Knowledge tool should be triggered automatically on context creation")
}
