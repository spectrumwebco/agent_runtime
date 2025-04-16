package langraph

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockEventEmitter struct {
	Events []models.Event
}

func NewMockEventEmitter() *MockEventEmitter {
	return &MockEventEmitter{
		Events: []models.Event{},
	}
}

func (m *MockEventEmitter) EmitEvent(event models.Event) error {
	m.Events = append(m.Events, event)
	return nil
}

type MockAgent struct {
	ID          string
	Role        string
	Capabilities []string
	EventEmitter *MockEventEmitter
	Messages     []string
}

func NewMockAgent(id, role string, capabilities []string, emitter *MockEventEmitter) *MockAgent {
	return &MockAgent{
		ID:           id,
		Role:         role,
		Capabilities: capabilities,
		EventEmitter: emitter,
		Messages:     []string{},
	}
}

func (a *MockAgent) ReceiveMessage(message string) {
	a.Messages = append(a.Messages, message)
	
	a.EventEmitter.EmitEvent(models.Event{
		ID:        "msg-" + time.Now().String(),
		Type:      models.EventTypeMessage,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": a.ID,
			"role":     a.Role,
			"message":  message,
		},
	})
}

func (a *MockAgent) SendMessage(targetAgentID, message string) {
	a.EventEmitter.EmitEvent(models.Event{
		ID:        "msg-" + time.Now().String(),
		Type:      models.EventTypeMessage,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id":       a.ID,
			"role":           a.Role,
			"target_agent_id": targetAgentID,
			"message":        message,
		},
	})
}

func (a *MockAgent) HasCapability(capability string) bool {
	for _, cap := range a.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

type MockMultiAgentSystem struct {
	Agents       map[string]*MockAgent
	EventEmitter *MockEventEmitter
}

func NewMockMultiAgentSystem(emitter *MockEventEmitter) *MockMultiAgentSystem {
	return &MockMultiAgentSystem{
		Agents:       make(map[string]*MockAgent),
		EventEmitter: emitter,
	}
}

func (s *MockMultiAgentSystem) AddAgent(agent *MockAgent) {
	s.Agents[agent.ID] = agent
}

func (s *MockMultiAgentSystem) GetAgent(id string) *MockAgent {
	return s.Agents[id]
}

func (s *MockMultiAgentSystem) SendMessageToAgent(sourceAgentID, targetAgentID, message string) error {
	sourceAgent := s.Agents[sourceAgentID]
	targetAgent := s.Agents[targetAgentID]
	
	if sourceAgent == nil || targetAgent == nil {
		return nil
	}
	
	sourceAgent.SendMessage(targetAgentID, message)
	targetAgent.ReceiveMessage(message)
	
	return nil
}

func TestMockMultiAgentCommunication(t *testing.T) {
	emitter := NewMockEventEmitter()
	system := NewMockMultiAgentSystem(emitter)
	
	frontendAgent := NewMockAgent("frontend-1", "frontend", []string{"ui_design", "code_generation"}, emitter)
	appBuilderAgent := NewMockAgent("appbuilder-1", "appbuilder", []string{"api_design", "database_design"}, emitter)
	codegenAgent := NewMockAgent("codegen-1", "codegen", []string{"code_generation", "code_review"}, emitter)
	engineeringAgent := NewMockAgent("engineering-1", "engineering", []string{"testing", "deployment"}, emitter)
	
	system.AddAgent(frontendAgent)
	system.AddAgent(appBuilderAgent)
	system.AddAgent(codegenAgent)
	system.AddAgent(engineeringAgent)
	
	assert.True(t, frontendAgent.HasCapability("ui_design"), "Frontend agent should have UI design capability")
	assert.True(t, appBuilderAgent.HasCapability("api_design"), "App builder agent should have API design capability")
	assert.True(t, codegenAgent.HasCapability("code_generation"), "Codegen agent should have code generation capability")
	assert.True(t, engineeringAgent.HasCapability("testing"), "Engineering agent should have testing capability")
	
	err := system.SendMessageToAgent("frontend-1", "codegen-1", "Generate code for UI component")
	assert.NoError(t, err, "Should send message without error")
	
	assert.Equal(t, 1, len(codegenAgent.Messages), "Codegen agent should have received 1 message")
	assert.Equal(t, "Generate code for UI component", codegenAgent.Messages[0], "Codegen agent should have received the correct message")
	
	err = system.SendMessageToAgent("codegen-1", "engineering-1", "Implement tests for generated code")
	assert.NoError(t, err, "Should send message without error")
	
	err = system.SendMessageToAgent("engineering-1", "appbuilder-1", "Create API for tested component")
	assert.NoError(t, err, "Should send message without error")
	
	assert.Equal(t, 1, len(engineeringAgent.Messages), "Engineering agent should have received 1 message")
	assert.Equal(t, "Implement tests for generated code", engineeringAgent.Messages[0], "Engineering agent should have received the correct message")
	
	assert.Equal(t, 1, len(appBuilderAgent.Messages), "App builder agent should have received 1 message")
	assert.Equal(t, "Create API for tested component", appBuilderAgent.Messages[0], "App builder agent should have received the correct message")
	
	assert.Equal(t, 6, len(emitter.Events), "6 events should have been emitted (3 sends and 3 receives)")
}

func TestMockAgentAutonomy(t *testing.T) {
	emitter := NewMockEventEmitter()
	system := NewMockMultiAgentSystem(emitter)
	
	frontendAgent := NewMockAgent("frontend-1", "frontend", []string{"ui_design", "code_generation"}, emitter)
	appBuilderAgent := NewMockAgent("appbuilder-1", "appbuilder", []string{"api_design", "database_design"}, emitter)
	codegenAgent := NewMockAgent("codegen-1", "codegen", []string{"code_generation", "code_review"}, emitter)
	engineeringAgent := NewMockAgent("engineering-1", "engineering", []string{"testing", "deployment"}, emitter)
	
	system.AddAgent(frontendAgent)
	system.AddAgent(appBuilderAgent)
	system.AddAgent(codegenAgent)
	system.AddAgent(engineeringAgent)
	
	
	frontendAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "task-1",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": frontendAgent.ID,
			"action":   "create_ui_component",
			"component": map[string]interface{}{
				"name":    "UserProfile",
				"type":    "component",
				"variant": "card",
			},
		},
	})
	
	err := system.SendMessageToAgent("frontend-1", "codegen-1", "Generate code for UserProfile component")
	assert.NoError(t, err, "Should send message without error")
	
	codegenAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "task-2",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": codegenAgent.ID,
			"action":   "generate_code",
			"code":     "function UserProfile() { return <div>...</div>; }",
		},
	})
	
	err = system.SendMessageToAgent("codegen-1", "engineering-1", "Implement tests for UserProfile component")
	assert.NoError(t, err, "Should send message without error")
	
	engineeringAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "task-3",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": engineeringAgent.ID,
			"action":   "implement_tests",
			"tests":    "test('renders user profile', () => { ... });",
		},
	})
	
	err = system.SendMessageToAgent("engineering-1", "appbuilder-1", "Create API for UserProfile component")
	assert.NoError(t, err, "Should send message without error")
	
	appBuilderAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "task-4",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": appBuilderAgent.ID,
			"action":   "create_api",
			"api":      "/api/users/{id}",
		},
	})
	
	err = system.SendMessageToAgent("appbuilder-1", "frontend-1", "Integrate API with UserProfile component")
	assert.NoError(t, err, "Should send message without error")
	
	assert.Equal(t, 1, len(frontendAgent.Messages), "Frontend agent should have received 1 message")
	assert.Equal(t, 1, len(codegenAgent.Messages), "Codegen agent should have received 1 message")
	assert.Equal(t, 1, len(engineeringAgent.Messages), "Engineering agent should have received 1 message")
	assert.Equal(t, 1, len(appBuilderAgent.Messages), "App builder agent should have received 1 message")
	
	assert.Equal(t, 12, len(emitter.Events), "12 events should have been emitted")
	
	var actionEvents []models.Event
	for _, event := range emitter.Events {
		if event.Type == models.EventTypeAction {
			actionEvents = append(actionEvents, event)
		}
	}
	
	assert.Equal(t, 4, len(actionEvents), "4 action events should have been emitted")
	
	actions := []string{}
	for _, event := range actionEvents {
		data, ok := event.Data.(map[string]interface{})
		if !ok {
			continue
		}
		
		action, ok := data["action"].(string)
		if !ok {
			continue
		}
		
		actions = append(actions, action)
	}
	
	expectedActions := []string{"create_ui_component", "generate_code", "implement_tests", "create_api"}
	assert.Equal(t, expectedActions, actions, "Actions should be in the expected order")
}

func TestMockMultiAgentWorkflow(t *testing.T) {
	
	emitter := NewMockEventEmitter()
	system := NewMockMultiAgentSystem(emitter)
	
	frontendAgent := NewMockAgent("frontend-1", "frontend", []string{"ui_design", "code_generation"}, emitter)
	appBuilderAgent := NewMockAgent("appbuilder-1", "appbuilder", []string{"api_design", "database_design"}, emitter)
	codegenAgent := NewMockAgent("codegen-1", "codegen", []string{"code_generation", "code_review"}, emitter)
	engineeringAgent := NewMockAgent("engineering-1", "engineering", []string{"testing", "deployment"}, emitter)
	orchestratorAgent := NewMockAgent("orchestrator-1", "orchestrator", []string{"orchestration", "planning"}, emitter)
	
	system.AddAgent(frontendAgent)
	system.AddAgent(appBuilderAgent)
	system.AddAgent(codegenAgent)
	system.AddAgent(engineeringAgent)
	system.AddAgent(orchestratorAgent)
	
	
	orchestratorAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "workflow-1",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": orchestratorAgent.ID,
			"action":   "receive_task",
			"task":     "Create a user profile page with API integration",
		},
	})
	
	err := system.SendMessageToAgent("orchestrator-1", "frontend-1", "Design a user profile page")
	assert.NoError(t, err, "Should send message without error")
	
	frontendAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "workflow-2",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": frontendAgent.ID,
			"action":   "design_ui",
			"ui":       "User profile page design",
		},
	})
	
	err = system.SendMessageToAgent("frontend-1", "appbuilder-1", "Design API for user profile")
	assert.NoError(t, err, "Should send message without error")
	
	err = system.SendMessageToAgent("frontend-1", "codegen-1", "Generate code for user profile UI")
	assert.NoError(t, err, "Should send message without error")
	
	appBuilderAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "workflow-3",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": appBuilderAgent.ID,
			"action":   "design_api",
			"api":      "User profile API design",
		},
	})
	
	err = system.SendMessageToAgent("appbuilder-1", "codegen-1", "Implement API for user profile")
	assert.NoError(t, err, "Should send message without error")
	
	codegenAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "workflow-4",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": codegenAgent.ID,
			"action":   "generate_code",
			"ui_code":  "User profile UI code",
			"api_code": "User profile API code",
		},
	})
	
	err = system.SendMessageToAgent("codegen-1", "engineering-1", "Test user profile UI and API")
	assert.NoError(t, err, "Should send message without error")
	
	engineeringAgent.EventEmitter.EmitEvent(models.Event{
		ID:        "workflow-5",
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id":  engineeringAgent.ID,
			"action":    "test_code",
			"ui_tests":  "User profile UI tests",
			"api_tests": "User profile API tests",
		},
	})
	
	err = system.SendMessageToAgent("engineering-1", "orchestrator-1", "User profile page with API integration is complete")
	assert.NoError(t, err, "Should send message without error")
	
	assert.Equal(t, 1, len(frontendAgent.Messages), "Frontend agent should have received 1 message")
	assert.Equal(t, 1, len(appBuilderAgent.Messages), "App builder agent should have received 1 message")
	assert.Equal(t, 2, len(codegenAgent.Messages), "Codegen agent should have received 2 messages")
	assert.Equal(t, 1, len(engineeringAgent.Messages), "Engineering agent should have received 1 message")
	assert.Equal(t, 1, len(orchestratorAgent.Messages), "Orchestrator agent should have received 1 message")
	
	
	messageMap := make(map[string]int)
	for _, event := range emitter.Events {
		if event.Type == models.EventTypeMessage {
			data, ok := event.Data.(map[string]interface{})
			if !ok {
				continue
			}
			
			sourceID, sourceOK := data["agent_id"].(string)
			targetID, targetOK := data["target_agent_id"].(string)
			
			if sourceOK && targetOK {
				key := sourceID + "->" + targetID
				messageMap[key]++
			}
		}
	}
	
	assert.Equal(t, 1, messageMap["orchestrator-1->frontend-1"], "Orchestrator should have sent 1 message to Frontend")
	assert.Equal(t, 1, messageMap["frontend-1->appbuilder-1"], "Frontend should have sent 1 message to App Builder")
	assert.Equal(t, 1, messageMap["frontend-1->codegen-1"], "Frontend should have sent 1 message to Codegen")
	assert.Equal(t, 1, messageMap["appbuilder-1->codegen-1"], "App Builder should have sent 1 message to Codegen")
	assert.Equal(t, 1, messageMap["codegen-1->engineering-1"], "Codegen should have sent 1 message to Engineering")
	assert.Equal(t, 1, messageMap["engineering-1->orchestrator-1"], "Engineering should have sent 1 message to Orchestrator")
}
