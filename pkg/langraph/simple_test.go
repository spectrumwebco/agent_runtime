package langraph

import (
	"testing"
	"time"
)

type SimpleAgent struct {
	ID          string
	Role        string
	Capabilities []string
	Messages     []string
}

func NewSimpleAgent(id, role string, capabilities []string) *SimpleAgent {
	return &SimpleAgent{
		ID:           id,
		Role:         role,
		Capabilities: capabilities,
		Messages:     []string{},
	}
}

func (a *SimpleAgent) ReceiveMessage(message string) {
	a.Messages = append(a.Messages, message)
}

func (a *SimpleAgent) HasCapability(capability string) bool {
	for _, cap := range a.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

type SimpleMultiAgentSystem struct {
	Agents map[string]*SimpleAgent
}

func NewSimpleMultiAgentSystem() *SimpleMultiAgentSystem {
	return &SimpleMultiAgentSystem{
		Agents: make(map[string]*SimpleAgent),
	}
}

func (s *SimpleMultiAgentSystem) AddAgent(agent *SimpleAgent) {
	s.Agents[agent.ID] = agent
}

func (s *SimpleMultiAgentSystem) GetAgent(id string) *SimpleAgent {
	return s.Agents[id]
}

func (s *SimpleMultiAgentSystem) SendMessageToAgent(sourceAgentID, targetAgentID, message string) error {
	targetAgent := s.Agents[targetAgentID]
	
	if targetAgent == nil {
		return nil
	}
	
	targetAgent.ReceiveMessage(message)
	return nil
}

func TestSimpleMultiAgentCommunication(t *testing.T) {
	system := NewSimpleMultiAgentSystem()
	
	frontendAgent := NewSimpleAgent("frontend-1", "frontend", []string{"ui_design", "code_generation"})
	appBuilderAgent := NewSimpleAgent("appbuilder-1", "appbuilder", []string{"api_design", "database_design"})
	codegenAgent := NewSimpleAgent("codegen-1", "codegen", []string{"code_generation", "code_review"})
	engineeringAgent := NewSimpleAgent("engineering-1", "engineering", []string{"testing", "deployment"})
	
	system.AddAgent(frontendAgent)
	system.AddAgent(appBuilderAgent)
	system.AddAgent(codegenAgent)
	system.AddAgent(engineeringAgent)
	
	if !frontendAgent.HasCapability("ui_design") {
		t.Errorf("Frontend agent should have UI design capability")
	}
	
	if !appBuilderAgent.HasCapability("api_design") {
		t.Errorf("App builder agent should have API design capability")
	}
	
	if !codegenAgent.HasCapability("code_generation") {
		t.Errorf("Codegen agent should have code generation capability")
	}
	
	if !engineeringAgent.HasCapability("testing") {
		t.Errorf("Engineering agent should have testing capability")
	}
	
	err := system.SendMessageToAgent("frontend-1", "codegen-1", "Generate code for UI component")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	if len(codegenAgent.Messages) != 1 {
		t.Errorf("Codegen agent should have received 1 message, got %d", len(codegenAgent.Messages))
	}
	
	if codegenAgent.Messages[0] != "Generate code for UI component" {
		t.Errorf("Codegen agent should have received the correct message, got %s", codegenAgent.Messages[0])
	}
	
	err = system.SendMessageToAgent("codegen-1", "engineering-1", "Implement tests for generated code")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	err = system.SendMessageToAgent("engineering-1", "appbuilder-1", "Create API for tested component")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	if len(engineeringAgent.Messages) != 1 {
		t.Errorf("Engineering agent should have received 1 message, got %d", len(engineeringAgent.Messages))
	}
	
	if engineeringAgent.Messages[0] != "Implement tests for generated code" {
		t.Errorf("Engineering agent should have received the correct message, got %s", engineeringAgent.Messages[0])
	}
	
	if len(appBuilderAgent.Messages) != 1 {
		t.Errorf("App builder agent should have received 1 message, got %d", len(appBuilderAgent.Messages))
	}
	
	if appBuilderAgent.Messages[0] != "Create API for tested component" {
		t.Errorf("App builder agent should have received the correct message, got %s", appBuilderAgent.Messages[0])
	}
}

func TestNonLinearWorkflow(t *testing.T) {
	system := NewSimpleMultiAgentSystem()
	
	frontendAgent := NewSimpleAgent("frontend-1", "frontend", []string{"ui_design", "code_generation"})
	appBuilderAgent := NewSimpleAgent("appbuilder-1", "appbuilder", []string{"api_design", "database_design"})
	codegenAgent := NewSimpleAgent("codegen-1", "codegen", []string{"code_generation", "code_review"})
	engineeringAgent := NewSimpleAgent("engineering-1", "engineering", []string{"testing", "deployment"})
	orchestratorAgent := NewSimpleAgent("orchestrator-1", "orchestrator", []string{"orchestration", "planning"})
	
	system.AddAgent(frontendAgent)
	system.AddAgent(appBuilderAgent)
	system.AddAgent(codegenAgent)
	system.AddAgent(engineeringAgent)
	system.AddAgent(orchestratorAgent)
	
	
	err := system.SendMessageToAgent("orchestrator-1", "frontend-1", "Design a user profile page")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	err = system.SendMessageToAgent("frontend-1", "appbuilder-1", "Design API for user profile")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	err = system.SendMessageToAgent("frontend-1", "codegen-1", "Generate code for user profile UI")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	err = system.SendMessageToAgent("appbuilder-1", "codegen-1", "Implement API for user profile")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	err = system.SendMessageToAgent("codegen-1", "engineering-1", "Test user profile UI and API")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	err = system.SendMessageToAgent("engineering-1", "orchestrator-1", "User profile page with API integration is complete")
	if err != nil {
		t.Errorf("Should send message without error: %v", err)
	}
	
	if len(frontendAgent.Messages) != 1 {
		t.Errorf("Frontend agent should have received 1 message, got %d", len(frontendAgent.Messages))
	}
	
	if len(appBuilderAgent.Messages) != 1 {
		t.Errorf("App builder agent should have received 1 message, got %d", len(appBuilderAgent.Messages))
	}
	
	if len(codegenAgent.Messages) != 2 {
		t.Errorf("Codegen agent should have received 2 messages, got %d", len(codegenAgent.Messages))
	}
	
	if len(engineeringAgent.Messages) != 1 {
		t.Errorf("Engineering agent should have received 1 message, got %d", len(engineeringAgent.Messages))
	}
	
	if len(orchestratorAgent.Messages) != 1 {
		t.Errorf("Orchestrator agent should have received 1 message, got %d", len(orchestratorAgent.Messages))
	}
}
