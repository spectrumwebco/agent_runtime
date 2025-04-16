package langraph

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	statemodels "github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type MockStateProvider struct {
	States map[string]map[string]interface{}
}

func NewMockStateProvider() *MockStateProvider {
	return &MockStateProvider{
		States: make(map[string]map[string]interface{}),
	}
}

func (m *MockStateProvider) GetState(ctx context.Context, key string) (map[string]interface{}, error) {
	state, exists := m.States[key]
	if !exists {
		return make(map[string]interface{}), nil
	}
	return state, nil
}

func (m *MockStateProvider) SetState(ctx context.Context, key string, state map[string]interface{}) error {
	m.States[key] = state
	return nil
}

func (m *MockStateProvider) UpdateState(ctx context.Context, key string, updates map[string]interface{}) error {
	state, exists := m.States[key]
	if !exists {
		state = make(map[string]interface{})
	}
	
	for k, v := range updates {
		state[k] = v
	}
	
	m.States[key] = state
	return nil
}

func (m *MockStateProvider) DeleteState(ctx context.Context, key string) error {
	delete(m.States, key)
	return nil
}

func (m *MockStateProvider) WatchState(ctx context.Context, key string, callback func(map[string]interface{})) error {
	state, exists := m.States[key]
	if !exists {
		state = make(map[string]interface{})
	}
	
	go callback(state)
	return nil
}

func TestFullyIntegratedMultiAgentSystem(t *testing.T) {
	eventStream := NewMockEventStream()
	stateProvider := NewMockStateProvider()
	
	system, err := CreateFullyIntegratedMultiAgentSystem("test-system", "Test fully integrated multi-agent system", eventStream, stateProvider)
	assert.NoError(t, err, "Should create fully integrated multi-agent system without error")
	assert.NotNil(t, system, "System should not be nil")
	
	agents := system.ListAgents()
	assert.Equal(t, 5, len(agents), "System should have 5 agents")
	
	var frontendAgent *Agent
	for _, agent := range agents {
		if agent.Config.Role == AgentRoleFrontend {
			frontendAgent = agent
			break
		}
	}
	assert.NotNil(t, frontendAgent, "Frontend agent should not be nil")
	
	stateIntegration, ok := system.Metadata["state_integration"].(*StateIntegration)
	assert.True(t, ok, "System should have state integration in metadata")
	assert.NotNil(t, stateIntegration, "State integration should not be nil")
	
	eventIntegration, ok := system.Metadata["event_integration"].(*EventIntegration)
	assert.True(t, ok, "System should have event integration in metadata")
	assert.NotNil(t, eventIntegration, "Event integration should not be nil")
	
	ctx := context.Background()
	agentState, err := stateIntegration.GetAgentState(ctx, frontendAgent.Config.ID)
	assert.NoError(t, err, "Should get agent state without error")
	assert.NotNil(t, agentState, "Agent state should not be nil")
	assert.Equal(t, frontendAgent.Config.ID, agentState["id"], "Agent state should have correct ID")
	assert.Equal(t, frontendAgent.Config.Name, agentState["name"], "Agent state should have correct name")
	assert.Equal(t, frontendAgent.Config.Role, agentState["role"], "Agent state should have correct role")
	
	err = stateIntegration.UpdateAgentState(ctx, frontendAgent.Config.ID, map[string]interface{}{
		"status": "busy",
		"task":   "Creating UI component",
	})
	assert.NoError(t, err, "Should update agent state without error")
	
	agentState, err = stateIntegration.GetAgentState(ctx, frontendAgent.Config.ID)
	assert.NoError(t, err, "Should get agent state without error")
	assert.Equal(t, "busy", agentState["status"], "Agent state should have updated status")
	assert.Equal(t, "Creating UI component", agentState["task"], "Agent state should have updated task")
	
	var orchestratorAgent *Agent
	for _, agent := range agents {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	assert.NotNil(t, orchestratorAgent, "Orchestrator agent should not be nil")
	
	execution, err := system.Execute(ctx, orchestratorAgent.Config.ID, map[string]interface{}{
		"task": "Create a new UI component",
	})
	assert.NoError(t, err, "Should execute system without error")
	assert.NotNil(t, execution, "Execution should not be nil")
	
	time.Sleep(100 * time.Millisecond)
	
	assert.Greater(t, len(eventStream.Events), 0, "Events should have been generated")
	
	err = eventIntegration.EmitEvent(models.EventTypeAction, models.EventSourceAgent, map[string]interface{}{
		"action":   "ui_component_created",
		"agent_id": frontendAgent.Config.ID,
		"component": map[string]interface{}{
			"name":    "Button",
			"type":    "button",
			"variant": "primary",
		},
	}, map[string]string{
		"agent_id": frontendAgent.Config.ID,
	})
	assert.NoError(t, err, "Should emit event without error")
	
	var eventFound bool
	for _, event := range eventStream.Events {
		if event.Type == models.EventTypeAction && event.Source == models.EventSourceAgent {
			data, ok := event.Data.(map[string]interface{})
			if !ok {
				continue
			}
			
			action, ok := data["action"].(string)
			if !ok {
				continue
			}
			
			if action == "ui_component_created" {
				eventFound = true
				break
			}
		}
	}
	assert.True(t, eventFound, "Event should have been added to stream")
}

func TestNonLinearAgentCommunication(t *testing.T) {
	eventStream := NewMockEventStream()
	stateProvider := NewMockStateProvider()
	
	system, err := CreateFullyIntegratedMultiAgentSystem("test-system", "Test fully integrated multi-agent system", eventStream, stateProvider)
	assert.NoError(t, err, "Should create fully integrated multi-agent system without error")
	
	var frontendAgent, appBuilderAgent, codegenAgent, engineeringAgent, orchestratorAgent *Agent
	for _, agent := range system.ListAgents() {
		switch agent.Config.Role {
		case AgentRoleFrontend:
			frontendAgent = agent
		case AgentRoleAppBuilder:
			appBuilderAgent = agent
		case AgentRoleCodegen:
			codegenAgent = agent
		case AgentRoleEngineering:
			engineeringAgent = agent
		case AgentRoleOrchestrator:
			orchestratorAgent = agent
		}
	}
	
	assert.NotNil(t, frontendAgent, "Frontend agent should not be nil")
	assert.NotNil(t, appBuilderAgent, "App builder agent should not be nil")
	assert.NotNil(t, codegenAgent, "Codegen agent should not be nil")
	assert.NotNil(t, engineeringAgent, "Engineering agent should not be nil")
	assert.NotNil(t, orchestratorAgent, "Orchestrator agent should not be nil")
	
	
	err = system.ConnectAgents(frontendAgent.Config.ID, codegenAgent.Config.ID, "frontend-codegen-custom", "Custom communication channel from frontend to codegen agent")
	assert.NoError(t, err, "Should connect frontend to codegen without error")
	
	err = system.ConnectAgents(codegenAgent.Config.ID, engineeringAgent.Config.ID, "codegen-engineering-custom", "Custom communication channel from codegen to engineering agent")
	assert.NoError(t, err, "Should connect codegen to engineering without error")
	
	err = system.ConnectAgents(engineeringAgent.Config.ID, appBuilderAgent.Config.ID, "engineering-appbuilder-custom", "Custom communication channel from engineering to app builder agent")
	assert.NoError(t, err, "Should connect engineering to app builder without error")
	
	workflow, err := system.CreateWorkflow("non-linear-workflow", "Non-linear communication workflow")
	assert.NoError(t, err, "Should create workflow without error")
	
	frontendNode := workflow.AddAgentNode(AgentTypeFrontend, "Frontend Node", "Frontend agent node", func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
		return frontendAgent.Process(ctx, inputs)
	})
	frontendNode.Metadata["agent_id"] = frontendAgent.Config.ID
	
	codegenNode := workflow.AddAgentNode(AgentTypeCodegen, "Codegen Node", "Codegen agent node", func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
		return codegenAgent.Process(ctx, inputs)
	})
	codegenNode.Metadata["agent_id"] = codegenAgent.Config.ID
	
	engineeringNode := workflow.AddAgentNode(AgentTypeEngineering, "Engineering Node", "Engineering agent node", func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
		return engineeringAgent.Process(ctx, inputs)
	})
	engineeringNode.Metadata["agent_id"] = engineeringAgent.Config.ID
	
	appBuilderNode := workflow.AddAgentNode(AgentTypeAppBuilder, "App Builder Node", "App builder agent node", func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
		return appBuilderAgent.Process(ctx, inputs)
	})
	appBuilderNode.Metadata["agent_id"] = appBuilderAgent.Config.ID
	
	_, err = workflow.AddEdge(frontendNode.ID, codegenNode.ID, "frontend-codegen", "Edge from frontend to codegen")
	assert.NoError(t, err, "Should add edge from frontend to codegen without error")
	
	_, err = workflow.AddEdge(codegenNode.ID, engineeringNode.ID, "codegen-engineering", "Edge from codegen to engineering")
	assert.NoError(t, err, "Should add edge from codegen to engineering without error")
	
	_, err = workflow.AddEdge(engineeringNode.ID, appBuilderNode.ID, "engineering-appbuilder", "Edge from engineering to app builder")
	assert.NoError(t, err, "Should add edge from engineering to app builder without error")
	
	ctx := context.Background()
	execution, err := system.ExecuteWorkflow(ctx, workflow, frontendNode.ID, map[string]interface{}{
		"task": "Create a new UI component with API integration",
	})
	assert.NoError(t, err, "Should execute workflow without error")
	assert.NotNil(t, execution, "Execution should not be nil")
	
	time.Sleep(100 * time.Millisecond)
	
	assert.Greater(t, len(eventStream.Events), 0, "Events should have been generated")
	
	var workflowStartEvent, workflowCompleteEvent bool
	for _, event := range eventStream.Events {
		if event.Type == models.EventTypeSystem && event.Source == models.EventSourceSystem {
			data, ok := event.Data.(map[string]interface{})
			if !ok {
				continue
			}
			
			action, ok := data["action"].(string)
			if !ok {
				continue
			}
			
			if action == "workflow_execution_started" {
				workflowStartEvent = true
			} else if action == "workflow_execution_completed" {
				workflowCompleteEvent = true
			}
		}
	}
	
	assert.True(t, workflowStartEvent, "Workflow execution started event should have been generated")
	assert.True(t, workflowCompleteEvent, "Workflow execution completed event should have been generated")
}

func TestAgentAutonomy(t *testing.T) {
	eventStream := NewMockEventStream()
	stateProvider := NewMockStateProvider()
	
	system, err := CreateFullyIntegratedMultiAgentSystem("test-system", "Test fully integrated multi-agent system", eventStream, stateProvider)
	assert.NoError(t, err, "Should create fully integrated multi-agent system without error")
	
	var orchestratorAgent *Agent
	for _, agent := range system.ListAgents() {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	assert.NotNil(t, orchestratorAgent, "Orchestrator agent should not be nil")
	
	orchestratorAgent.AddTool("decide_next_agent", "Tool to decide which agent should handle a task next", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		task, ok := params["task"].(string)
		if !ok {
			return map[string]interface{}{
				"error": "Task parameter is required",
			}, nil
		}
		
		var nextAgentRole AgentRole
		if task == "Create UI component" {
			nextAgentRole = AgentRoleFrontend
		} else if task == "Design API" {
			nextAgentRole = AgentRoleAppBuilder
		} else if task == "Generate code" {
			nextAgentRole = AgentRoleCodegen
		} else if task == "Implement tests" {
			nextAgentRole = AgentRoleEngineering
		} else {
			nextAgentRole = AgentRoleFrontend
		}
		
		var nextAgentID string
		for _, a := range system.ListAgents() {
			if a.Config.Role == nextAgentRole {
				nextAgentID = a.Config.ID
				break
			}
		}
		
		return map[string]interface{}{
			"next_agent_id":   nextAgentID,
			"next_agent_role": nextAgentRole,
			"reason":          "Selected based on task requirements",
		}, nil
	}, map[string]interface{}{
		"task": "string",
	})
	
	ctx := context.Background()
	result, err := orchestratorAgent.ExecuteTool(ctx, "decide_next_agent", map[string]interface{}{
		"task": "Create UI component",
	})
	assert.NoError(t, err, "Should execute decide_next_agent tool without error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, AgentRoleFrontend, result["next_agent_role"], "Next agent role should be frontend")
	
	result, err = orchestratorAgent.ExecuteTool(ctx, "decide_next_agent", map[string]interface{}{
		"task": "Design API",
	})
	assert.NoError(t, err, "Should execute decide_next_agent tool without error")
	assert.Equal(t, AgentRoleAppBuilder, result["next_agent_role"], "Next agent role should be app builder")
	
	result, err = orchestratorAgent.ExecuteTool(ctx, "decide_next_agent", map[string]interface{}{
		"task": "Generate code",
	})
	assert.NoError(t, err, "Should execute decide_next_agent tool without error")
	assert.Equal(t, AgentRoleCodegen, result["next_agent_role"], "Next agent role should be codegen")
	
	result, err = orchestratorAgent.ExecuteTool(ctx, "decide_next_agent", map[string]interface{}{
		"task": "Implement tests",
	})
	assert.NoError(t, err, "Should execute decide_next_agent tool without error")
	assert.Equal(t, AgentRoleEngineering, result["next_agent_role"], "Next agent role should be engineering")
	
	execution, err := system.Execute(ctx, orchestratorAgent.Config.ID, map[string]interface{}{
		"task": "Create UI component",
	})
	assert.NoError(t, err, "Should execute system without error")
	assert.NotNil(t, execution, "Execution should not be nil")
	
	time.Sleep(100 * time.Millisecond)
	
	assert.Greater(t, len(eventStream.Events), 0, "Events should have been generated")
}

func TestLangChainIntegration(t *testing.T) {
	t.Skip("Skipping LangChain integration test as it requires actual LangChain-Go implementation")
	
}

func TestAgentPassingInformation(t *testing.T) {
	eventStream := NewMockEventStream()
	stateProvider := NewMockStateProvider()
	
	system, err := CreateFullyIntegratedMultiAgentSystem("test-system", "Test fully integrated multi-agent system", eventStream, stateProvider)
	assert.NoError(t, err, "Should create fully integrated multi-agent system without error")
	
	var frontendAgent, appBuilderAgent, codegenAgent, engineeringAgent *Agent
	for _, agent := range system.ListAgents() {
		switch agent.Config.Role {
		case AgentRoleFrontend:
			frontendAgent = agent
		case AgentRoleAppBuilder:
			appBuilderAgent = agent
		case AgentRoleCodegen:
			codegenAgent = agent
		case AgentRoleEngineering:
			engineeringAgent = agent
		}
	}
	
	assert.NotNil(t, frontendAgent, "Frontend agent should not be nil")
	assert.NotNil(t, appBuilderAgent, "App builder agent should not be nil")
	assert.NotNil(t, codegenAgent, "Codegen agent should not be nil")
	assert.NotNil(t, engineeringAgent, "Engineering agent should not be nil")
	
	stateIntegration, ok := system.Metadata["state_integration"].(*StateIntegration)
	assert.True(t, ok, "System should have state integration in metadata")
	assert.NotNil(t, stateIntegration, "State integration should not be nil")
	
	ctx := context.Background()
	
	err = stateIntegration.UpdateAgentState(ctx, frontendAgent.Config.ID, map[string]interface{}{
		"status": "busy",
		"task":   "Creating UI component",
		"output": map[string]interface{}{
			"component": map[string]interface{}{
				"name":    "UserProfile",
				"type":    "component",
				"variant": "card",
				"props": map[string]interface{}{
					"user_id": "string",
					"avatar":  "string",
					"name":    "string",
					"bio":     "string",
				},
			},
		},
	})
	assert.NoError(t, err, "Should update frontend agent state without error")
	
	err = stateIntegration.UpdateAgentState(ctx, appBuilderAgent.Config.ID, map[string]interface{}{
		"status": "busy",
		"task":   "Creating API for UI component",
		"input": map[string]interface{}{
			"component": map[string]interface{}{
				"name":    "UserProfile",
				"type":    "component",
				"variant": "card",
				"props": map[string]interface{}{
					"user_id": "string",
					"avatar":  "string",
					"name":    "string",
					"bio":     "string",
				},
			},
		},
		"output": map[string]interface{}{
			"api": map[string]interface{}{
				"endpoint": "/api/users/{user_id}",
				"method":   "GET",
				"response": map[string]interface{}{
					"user_id": "string",
					"avatar":  "string",
					"name":    "string",
					"bio":     "string",
				},
			},
		},
	})
	assert.NoError(t, err, "Should update app builder agent state without error")
	
	err = stateIntegration.UpdateAgentState(ctx, codegenAgent.Config.ID, map[string]interface{}{
		"status": "busy",
		"task":   "Generating code for UI component and API",
		"input": map[string]interface{}{
			"component": map[string]interface{}{
				"name":    "UserProfile",
				"type":    "component",
				"variant": "card",
				"props": map[string]interface{}{
					"user_id": "string",
					"avatar":  "string",
					"name":    "string",
					"bio":     "string",
				},
			},
			"api": map[string]interface{}{
				"endpoint": "/api/users/{user_id}",
				"method":   "GET",
				"response": map[string]interface{}{
					"user_id": "string",
					"avatar":  "string",
					"name":    "string",
					"bio":     "string",
				},
			},
		},
		"output": map[string]interface{}{
			"component_code": "function UserProfile({ user_id, avatar, name, bio }) { return <div>...</div>; }",
			"api_code":       "app.get('/api/users/:user_id', (req, res) => { ... });",
		},
	})
	assert.NoError(t, err, "Should update codegen agent state without error")
	
	err = stateIntegration.UpdateAgentState(ctx, engineeringAgent.Config.ID, map[string]interface{}{
		"status": "busy",
		"task":   "Implementing tests for UI component and API",
		"input": map[string]interface{}{
			"component_code": "function UserProfile({ user_id, avatar, name, bio }) { return <div>...</div>; }",
			"api_code":       "app.get('/api/users/:user_id', (req, res) => { ... });",
		},
		"output": map[string]interface{}{
			"component_test_code": "test('renders user profile', () => { ... });",
			"api_test_code":       "test('returns user data', async () => { ... });",
		},
	})
	assert.NoError(t, err, "Should update engineering agent state without error")
	
	frontendState, err := stateIntegration.GetAgentState(ctx, frontendAgent.Config.ID)
	assert.NoError(t, err, "Should get frontend agent state without error")
	assert.Equal(t, "Creating UI component", frontendState["task"], "Frontend agent task should be correct")
	
	appBuilderState, err := stateIntegration.GetAgentState(ctx, appBuilderAgent.Config.ID)
	assert.NoError(t, err, "Should get app builder agent state without error")
	assert.Equal(t, "Creating API for UI component", appBuilderState["task"], "App builder agent task should be correct")
	
	codegenState, err := stateIntegration.GetAgentState(ctx, codegenAgent.Config.ID)
	assert.NoError(t, err, "Should get codegen agent state without error")
	assert.Equal(t, "Generating code for UI component and API", codegenState["task"], "Codegen agent task should be correct")
	
	engineeringState, err := stateIntegration.GetAgentState(ctx, engineeringAgent.Config.ID)
	assert.NoError(t, err, "Should get engineering agent state without error")
	assert.Equal(t, "Implementing tests for UI component and API", engineeringState["task"], "Engineering agent task should be correct")
	
	appBuilderInput, ok := appBuilderState["input"].(map[string]interface{})
	assert.True(t, ok, "App builder agent state should have input")
	assert.NotNil(t, appBuilderInput["component"], "App builder agent input should have component")
	
	codegenInput, ok := codegenState["input"].(map[string]interface{})
	assert.True(t, ok, "Codegen agent state should have input")
	assert.NotNil(t, codegenInput["component"], "Codegen agent input should have component")
	assert.NotNil(t, codegenInput["api"], "Codegen agent input should have API")
	
	engineeringInput, ok := engineeringState["input"].(map[string]interface{})
	assert.True(t, ok, "Engineering agent state should have input")
	assert.NotNil(t, engineeringInput["component_code"], "Engineering agent input should have component code")
	assert.NotNil(t, engineeringInput["api_code"], "Engineering agent input should have API code")
}
