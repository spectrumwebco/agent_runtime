package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type StateProvider interface {
	GetState(ctx context.Context, key string) (map[string]interface{}, error)
	SetState(ctx context.Context, key string, state map[string]interface{}) error
	UpdateState(ctx context.Context, key string, updates map[string]interface{}) error
	DeleteState(ctx context.Context, key string) error
	WatchState(ctx context.Context, key string, callback func(map[string]interface{})) error
}

type StateIntegration struct {
	Provider      StateProvider
	EventStream   EventStream
	AgentStates   map[string]string // Maps agent ID to state key
	SystemStates  map[string]string // Maps system ID to state key
	lock          sync.RWMutex
}

func NewStateIntegration(provider StateProvider, eventStream EventStream) *StateIntegration {
	return &StateIntegration{
		Provider:     provider,
		EventStream:  eventStream,
		AgentStates:  make(map[string]string),
		SystemStates: make(map[string]string),
	}
}

func (s *StateIntegration) RegisterAgent(agent *Agent) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	
	stateKey := fmt.Sprintf("agent:%s", agent.Config.ID)
	s.AgentStates[agent.Config.ID] = stateKey
	
	initialState := map[string]interface{}{
		"id":          agent.Config.ID,
		"name":        agent.Config.Name,
		"role":        agent.Config.Role,
		"capabilities": agent.Config.Capabilities,
		"created_at":  time.Now().UTC(),
		"updated_at":  time.Now().UTC(),
		"status":      "ready",
		"memory":      make(map[string]interface{}),
		"context":     make(map[string]interface{}),
		"tools":       make(map[string]interface{}),
	}
	
	if err := s.Provider.SetState(context.Background(), stateKey, initialState); err != nil {
		return fmt.Errorf("failed to initialize agent state: %v", err)
	}
	
	if s.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeState,
			models.EventSourceAgent,
			map[string]interface{}{
				"action":     "state_initialized",
				"agent_id":   agent.Config.ID,
				"agent_name": agent.Config.Name,
				"agent_role": agent.Config.Role,
				"state_key":  stateKey,
			},
			map[string]string{
				"agent_id":  agent.Config.ID,
				"state_key": stateKey,
			},
		)
		
		if err := s.EventStream.AddEvent(event); err != nil {
			return fmt.Errorf("failed to add event to stream: %v", err)
		}
	}
	
	return nil
}

func (s *StateIntegration) RegisterSystem(system *MultiAgentSystem) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	
	stateKey := fmt.Sprintf("system:%s", system.ID)
	s.SystemStates[system.ID] = stateKey
	
	initialState := map[string]interface{}{
		"id":          system.ID,
		"name":        system.Name,
		"description": system.Description,
		"created_at":  time.Now().UTC(),
		"updated_at":  time.Now().UTC(),
		"status":      "ready",
		"agents":      make([]string, 0),
		"metadata":    system.Metadata,
	}
	
	agentIDs := make([]string, 0)
	for _, agent := range system.ListAgents() {
		agentIDs = append(agentIDs, agent.Config.ID)
	}
	initialState["agents"] = agentIDs
	
	if err := s.Provider.SetState(context.Background(), stateKey, initialState); err != nil {
		return fmt.Errorf("failed to initialize system state: %v", err)
	}
	
	if s.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeState,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "state_initialized",
				"system_id":  system.ID,
				"system_name": system.Name,
				"state_key":  stateKey,
			},
			map[string]string{
				"system_id": system.ID,
				"state_key": stateKey,
			},
		)
		
		if err := s.EventStream.AddEvent(event); err != nil {
			return fmt.Errorf("failed to add event to stream: %v", err)
		}
	}
	
	for _, agent := range system.ListAgents() {
		if err := s.RegisterAgent(agent); err != nil {
			return fmt.Errorf("failed to register agent %s: %v", agent.Config.ID, err)
		}
	}
	
	return nil
}

func (s *StateIntegration) GetAgentState(ctx context.Context, agentID string) (map[string]interface{}, error) {
	s.lock.RLock()
	stateKey, exists := s.AgentStates[agentID]
	s.lock.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("agent %s not registered with state integration", agentID)
	}
	
	return s.Provider.GetState(ctx, stateKey)
}

func (s *StateIntegration) UpdateAgentState(ctx context.Context, agentID string, updates map[string]interface{}) error {
	s.lock.RLock()
	stateKey, exists := s.AgentStates[agentID]
	s.lock.RUnlock()
	
	if !exists {
		return fmt.Errorf("agent %s not registered with state integration", agentID)
	}
	
	updates["updated_at"] = time.Now().UTC()
	
	if err := s.Provider.UpdateState(ctx, stateKey, updates); err != nil {
		return fmt.Errorf("failed to update agent state: %v", err)
	}
	
	if s.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeState,
			models.EventSourceAgent,
			map[string]interface{}{
				"action":     "state_updated",
				"agent_id":   agentID,
				"state_key":  stateKey,
				"updates":    updates,
			},
			map[string]string{
				"agent_id":  agentID,
				"state_key": stateKey,
			},
		)
		
		if err := s.EventStream.AddEvent(event); err != nil {
			return fmt.Errorf("failed to add event to stream: %v", err)
		}
	}
	
	return nil
}

func (s *StateIntegration) GetSystemState(ctx context.Context, systemID string) (map[string]interface{}, error) {
	s.lock.RLock()
	stateKey, exists := s.SystemStates[systemID]
	s.lock.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("system %s not registered with state integration", systemID)
	}
	
	return s.Provider.GetState(ctx, stateKey)
}

func (s *StateIntegration) UpdateSystemState(ctx context.Context, systemID string, updates map[string]interface{}) error {
	s.lock.RLock()
	stateKey, exists := s.SystemStates[systemID]
	s.lock.RUnlock()
	
	if !exists {
		return fmt.Errorf("system %s not registered with state integration", systemID)
	}
	
	updates["updated_at"] = time.Now().UTC()
	
	if err := s.Provider.UpdateState(ctx, stateKey, updates); err != nil {
		return fmt.Errorf("failed to update system state: %v", err)
	}
	
	if s.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeState,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "state_updated",
				"system_id":  systemID,
				"state_key":  stateKey,
				"updates":    updates,
			},
			map[string]string{
				"system_id": systemID,
				"state_key": stateKey,
			},
		)
		
		if err := s.EventStream.AddEvent(event); err != nil {
			return fmt.Errorf("failed to add event to stream: %v", err)
		}
	}
	
	return nil
}

func (s *StateIntegration) WatchAgentState(ctx context.Context, agentID string, callback func(map[string]interface{})) error {
	s.lock.RLock()
	stateKey, exists := s.AgentStates[agentID]
	s.lock.RUnlock()
	
	if !exists {
		return fmt.Errorf("agent %s not registered with state integration", agentID)
	}
	
	return s.Provider.WatchState(ctx, stateKey, callback)
}

func (s *StateIntegration) WatchSystemState(ctx context.Context, systemID string, callback func(map[string]interface{})) error {
	s.lock.RLock()
	stateKey, exists := s.SystemStates[systemID]
	s.lock.RUnlock()
	
	if !exists {
		return fmt.Errorf("system %s not registered with state integration", systemID)
	}
	
	return s.Provider.WatchState(ctx, stateKey, callback)
}

func (s *StateIntegration) CreateStateAwareAgent(config AgentConfig, stateProvider StateProvider) (*Agent, error) {
	agent, err := NewAgent(config)
	if err != nil {
		return nil, err
	}
	
	if err := s.RegisterAgent(agent); err != nil {
		return nil, err
	}
	
	originalProcess := agent.Process
	agent.Process = func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		state, err := s.GetAgentState(ctx, agent.Config.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent state: %v", err)
		}
		
		stateInputs := make(map[string]interface{})
		for k, v := range inputs {
			stateInputs[k] = v
		}
		stateInputs["state"] = state
		
		outputs, err := originalProcess(ctx, stateInputs)
		if err != nil {
			return nil, err
		}
		
		if stateUpdates, ok := outputs["state_updates"].(map[string]interface{}); ok {
			if err := s.UpdateAgentState(ctx, agent.Config.ID, stateUpdates); err != nil {
				return nil, fmt.Errorf("failed to update agent state: %v", err)
			}
			
			delete(outputs, "state_updates")
		}
		
		return outputs, nil
	}
	
	return agent, nil
}

func CreateStateAwareMultiAgentSystem(name, description string, eventStream EventStream, stateProvider StateProvider) (*MultiAgentSystem, error) {
	system := NewMultiAgentSystem(name, description, eventStream)
	
	stateIntegration := NewStateIntegration(stateProvider, eventStream)
	
	if err := stateIntegration.RegisterSystem(system); err != nil {
		return nil, fmt.Errorf("failed to register system with state integration: %v", err)
	}
	
	system.Metadata["state_integration"] = stateIntegration
	
	frontendAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Frontend Agent",
		Description: "Agent responsible for UI/UX development tasks",
		Role:        AgentRoleFrontend,
		Capabilities: []AgentCapability{
			AgentCapabilityUIDesign,
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
		},
	}
	frontendAgent, err := stateIntegration.CreateStateAwareAgent(frontendAgentConfig, stateProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend agent: %v", err)
	}
	system.AddAgent(frontendAgent)
	
	appBuilderAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "App Builder Agent",
		Description: "Agent responsible for application assembly and API design",
		Role:        AgentRoleAppBuilder,
		Capabilities: []AgentCapability{
			AgentCapabilityAPIDesign,
			AgentCapabilityDatabaseDesign,
			AgentCapabilityArchitecture,
		},
	}
	appBuilderAgent, err := stateIntegration.CreateStateAwareAgent(appBuilderAgentConfig, stateProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create app builder agent: %v", err)
	}
	system.AddAgent(appBuilderAgent)
	
	codegenAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Codegen Agent",
		Description: "Agent responsible for code generation and review",
		Role:        AgentRoleCodegen,
		Capabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
			AgentCapabilityTesting,
		},
	}
	codegenAgent, err := stateIntegration.CreateStateAwareAgent(codegenAgentConfig, stateProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create codegen agent: %v", err)
	}
	system.AddAgent(codegenAgent)
	
	engineeringAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Engineering Agent",
		Description: "Agent responsible for core development tasks",
		Role:        AgentRoleEngineering,
		Capabilities: []AgentCapability{
			AgentCapabilityArchitecture,
			AgentCapabilityTesting,
			AgentCapabilityDeployment,
			AgentCapabilityDebugging,
		},
	}
	engineeringAgent, err := stateIntegration.CreateStateAwareAgent(engineeringAgentConfig, stateProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create engineering agent: %v", err)
	}
	system.AddAgent(engineeringAgent)
	
	orchestratorAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Orchestrator Agent",
		Description: "Agent responsible for coordinating other agents",
		Role:        AgentRoleOrchestrator,
		Capabilities: []AgentCapability{
			AgentCapabilityOrchestration,
			AgentCapabilityPlanning,
			AgentCapabilityTaskManagement,
		},
	}
	orchestratorAgent, err := stateIntegration.CreateStateAwareAgent(orchestratorAgentConfig, stateProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create orchestrator agent: %v", err)
	}
	system.AddAgent(orchestratorAgent)
	
	if err := system.ConnectAgents(orchestratorAgent.Config.ID, frontendAgent.Config.ID, "orchestrator-frontend", "Communication channel from orchestrator to frontend agent"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to frontend agent: %v", err)
	}
	
	if err := system.ConnectAgents(orchestratorAgent.Config.ID, appBuilderAgent.Config.ID, "orchestrator-appbuilder", "Communication channel from orchestrator to app builder agent"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to app builder agent: %v", err)
	}
	
	if err := system.ConnectAgents(orchestratorAgent.Config.ID, codegenAgent.Config.ID, "orchestrator-codegen", "Communication channel from orchestrator to codegen agent"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to codegen agent: %v", err)
	}
	
	if err := system.ConnectAgents(orchestratorAgent.Config.ID, engineeringAgent.Config.ID, "orchestrator-engineering", "Communication channel from orchestrator to engineering agent"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to engineering agent: %v", err)
	}
	
	if err := system.ConnectAgents(frontendAgent.Config.ID, codegenAgent.Config.ID, "frontend-codegen", "Communication channel from frontend agent to codegen agent"); err != nil {
		return nil, fmt.Errorf("failed to connect frontend agent to codegen agent: %v", err)
	}
	
	if err := system.ConnectAgents(appBuilderAgent.Config.ID, codegenAgent.Config.ID, "appbuilder-codegen", "Communication channel from app builder agent to codegen agent"); err != nil {
		return nil, fmt.Errorf("failed to connect app builder agent to codegen agent: %v", err)
	}
	
	if err := system.ConnectAgents(codegenAgent.Config.ID, engineeringAgent.Config.ID, "codegen-engineering", "Communication channel from codegen agent to engineering agent"); err != nil {
		return nil, fmt.Errorf("failed to connect codegen agent to engineering agent: %v", err)
	}
	
	return system, nil
}
