package langraph

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

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

func TestMultiAgentSystem(t *testing.T) {
	eventStream := NewMockEventStream()

	system, err := CreateStandardMultiAgentSystem("test-system", "Test multi-agent system", eventStream)
	assert.NoError(t, err, "Should create standard multi-agent system without error")
	assert.NotNil(t, system, "System should not be nil")

	agents := system.ListAgents()
	assert.Equal(t, 5, len(agents), "System should have 5 agents")

	var foundFrontend, foundAppBuilder, foundCodegen, foundEngineering, foundOrchestrator bool
	for _, agent := range agents {
		switch agent.Config.Role {
		case AgentRoleFrontend:
			foundFrontend = true
		case AgentRoleAppBuilder:
			foundAppBuilder = true
		case AgentRoleCodegen:
			foundCodegen = true
		case AgentRoleEngineering:
			foundEngineering = true
		case AgentRoleOrchestrator:
			foundOrchestrator = true
		}
	}
	assert.True(t, foundFrontend, "System should have a frontend agent")
	assert.True(t, foundAppBuilder, "System should have an app builder agent")
	assert.True(t, foundCodegen, "System should have a codegen agent")
	assert.True(t, foundEngineering, "System should have an engineering agent")
	assert.True(t, foundOrchestrator, "System should have an orchestrator agent")

	var orchestratorAgent *Agent
	for _, agent := range agents {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	assert.NotNil(t, orchestratorAgent, "Orchestrator agent should not be nil")

	ctx := context.Background()
	execution, err := system.Execute(ctx, orchestratorAgent.Config.ID, map[string]interface{}{
		"task": "Test task",
	})
	assert.NoError(t, err, "Should execute system without error")
	assert.NotNil(t, execution, "Execution should not be nil")

	time.Sleep(100 * time.Millisecond)

	assert.Greater(t, len(eventStream.Events), 0, "Events should have been generated")
}

func TestAgentCapabilities(t *testing.T) {
	eventStream := NewMockEventStream()

	system, err := CreateStandardMultiAgentSystem("test-system", "Test multi-agent system", eventStream)
	assert.NoError(t, err, "Should create standard multi-agent system without error")
	assert.NotNil(t, system, "System should not be nil")

	agents := system.ListAgents()
	
	for _, agent := range agents {
		switch agent.Config.Role {
		case AgentRoleFrontend:
			assert.True(t, agent.HasCapability(AgentCapabilityUIDesign), "Frontend agent should have UI design capability")
			assert.True(t, agent.HasCapability(AgentCapabilityCodeGeneration), "Frontend agent should have code generation capability")
		case AgentRoleAppBuilder:
			assert.True(t, agent.HasCapability(AgentCapabilityAPIDesign), "App builder agent should have API design capability")
			assert.True(t, agent.HasCapability(AgentCapabilityDatabaseDesign), "App builder agent should have database design capability")
		case AgentRoleCodegen:
			assert.True(t, agent.HasCapability(AgentCapabilityCodeGeneration), "Codegen agent should have code generation capability")
			assert.True(t, agent.HasCapability(AgentCapabilityCodeReview), "Codegen agent should have code review capability")
		case AgentRoleEngineering:
			assert.True(t, agent.HasCapability(AgentCapabilityTesting), "Engineering agent should have testing capability")
			assert.True(t, agent.HasCapability(AgentCapabilityDeployment), "Engineering agent should have deployment capability")
		case AgentRoleOrchestrator:
			assert.True(t, agent.HasCapability(AgentCapabilityOrchestration), "Orchestrator agent should have orchestration capability")
			assert.True(t, agent.HasCapability(AgentCapabilityPlanning), "Orchestrator agent should have planning capability")
		}
	}
}

func TestAgentCommunication(t *testing.T) {
	eventStream := NewMockEventStream()

	system, err := CreateStandardMultiAgentSystem("test-system", "Test multi-agent system", eventStream)
	assert.NoError(t, err, "Should create standard multi-agent system without error")
	assert.NotNil(t, system, "System should not be nil")

	agents := system.ListAgents()
	
	var orchestratorAgent *Agent
	for _, agent := range agents {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	assert.NotNil(t, orchestratorAgent, "Orchestrator agent should not be nil")

	var frontendAgent *Agent
	for _, agent := range agents {
		if agent.Config.Role == AgentRoleFrontend {
			frontendAgent = agent
			break
		}
	}
	assert.NotNil(t, frontendAgent, "Frontend agent should not be nil")

	ctx := context.Background()
	execution, err := system.Execute(ctx, orchestratorAgent.Config.ID, map[string]interface{}{
		"task": "Create a new UI component",
		"target_agent": frontendAgent.Config.ID,
	})
	assert.NoError(t, err, "Should execute system without error")
	assert.NotNil(t, execution, "Execution should not be nil")

	time.Sleep(100 * time.Millisecond)

	assert.Greater(t, len(eventStream.Events), 0, "Events should have been generated")

	var frontendAgentReceived bool
	for _, event := range eventStream.Events {
		if event.Type == models.EventTypeAction && event.Source == models.EventSourceAgent {
			data, ok := event.Data.(map[string]interface{})
			if !ok {
				continue
			}
			
			agentID, ok := data["agent_id"].(string)
			if !ok {
				continue
			}
			
			if agentID == frontendAgent.Config.ID {
				frontendAgentReceived = true
				break
			}
		}
	}
	assert.True(t, frontendAgentReceived, "Frontend agent should have received the task")
}

func TestIntegrationWithDjango(t *testing.T) {
	t.Skip("Skipping Django integration test as it requires actual Django backend")


	eventStream := NewMockEventStream()

	system := NewMultiAgentSystem("django-integration-system", "Multi-agent system with Django integration", eventStream)
	assert.NotNil(t, system, "System should not be nil")

	djangoFactory := NewDjangoAgentFactory("http://localhost:8000", "dummy-api-key", eventStream)
	assert.NotNil(t, djangoFactory, "Django factory should not be nil")

	system.SetDjangoFactory(djangoFactory)

	//
	//
}

func TestToolTriggeringOnContextCreation(t *testing.T) {
	eventStream := NewMockEventStream()

	system, err := CreateStandardMultiAgentSystem("test-system", "Test multi-agent system", eventStream)
	assert.NoError(t, err, "Should create standard multi-agent system without error")
	assert.NotNil(t, system, "System should not be nil")

	agents := system.ListAgents()
	
	var orchestratorAgent *Agent
	for _, agent := range agents {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	assert.NotNil(t, orchestratorAgent, "Orchestrator agent should not be nil")

	knowledgeToolTriggered := false
	orchestratorAgent.AddTool("knowledge-tool", "Knowledge tool that is triggered automatically on context creation", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		knowledgeToolTriggered = true
		return map[string]interface{}{
			"knowledge": "Knowledge about: " + params["query"].(string),
		}, nil
	}, map[string]interface{}{
		"query": "string",
	})

	ctx := context.Background()
	execution, err := system.Execute(ctx, orchestratorAgent.Config.ID, map[string]interface{}{
		"task": "Create a new context",
		"query": "test query",
	})
	assert.NoError(t, err, "Should execute system without error")
	assert.NotNil(t, execution, "Execution should not be nil")

	time.Sleep(100 * time.Millisecond)

	result, err := orchestratorAgent.ExecuteTool(ctx, "knowledge-tool", map[string]interface{}{
		"query": "test query",
	})
	assert.NoError(t, err, "Should execute knowledge tool without error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, "Knowledge about: test query", result["knowledge"], "Knowledge tool should return expected result")

	assert.True(t, knowledgeToolTriggered, "Knowledge tool should be triggered")
}
