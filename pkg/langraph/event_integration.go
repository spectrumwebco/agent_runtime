package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type EventIntegration struct {
	EventStream       EventStream
	AgentSubscriptions map[string][]func(*models.Event)
	SystemSubscriptions map[string][]func(*models.Event)
	lock              sync.RWMutex
}

func NewEventIntegration(eventStream EventStream) *EventIntegration {
	return &EventIntegration{
		EventStream:        eventStream,
		AgentSubscriptions: make(map[string][]func(*models.Event)),
		SystemSubscriptions: make(map[string][]func(*models.Event)),
	}
}

func (e *EventIntegration) RegisterAgent(agent *Agent) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	
	agentEventHandler := func(event *models.Event) {
		if event.Source == models.EventSourceAgent {
			if agentID, ok := event.Tags["agent_id"]; ok && agentID == agent.Config.ID {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					
					_, err := agent.Process(ctx, map[string]interface{}{
						"event_type": event.Type,
						"event_source": event.Source,
						"event_data": event.Data,
						"event_tags": event.Tags,
						"event_timestamp": event.Timestamp,
					})
					
					if err != nil {
						fmt.Printf("Error processing event with agent %s: %v\n", agent.Config.ID, err)
					}
				}()
			}
		}
	}
	
	if err := e.EventStream.Subscribe(models.EventTypeAction, agentEventHandler); err != nil {
		return fmt.Errorf("failed to subscribe agent to action events: %v", err)
	}
	
	if err := e.EventStream.Subscribe(models.EventTypeTask, agentEventHandler); err != nil {
		return fmt.Errorf("failed to subscribe agent to task events: %v", err)
	}
	
	if err := e.EventStream.Subscribe(models.EventTypeState, agentEventHandler); err != nil {
		return fmt.Errorf("failed to subscribe agent to state events: %v", err)
	}
	
	e.AgentSubscriptions[agent.Config.ID] = []func(*models.Event){agentEventHandler}
	
	event := models.NewEvent(
		models.EventTypeSystem,
		models.EventSourceSystem,
		map[string]interface{}{
			"action":     "agent_registered_with_event_stream",
			"agent_id":   agent.Config.ID,
			"agent_name": agent.Config.Name,
			"agent_role": agent.Config.Role,
		},
		map[string]string{
			"agent_id": agent.Config.ID,
		},
	)
	
	if err := e.EventStream.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add event to stream: %v", err)
	}
	
	originalProcess := agent.Process
	agent.Process = func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		outputs, err := originalProcess(ctx, inputs)
		if err != nil {
			return nil, err
		}
		
		if events, ok := outputs["events"].([]map[string]interface{}); ok {
			for _, eventData := range events {
				eventType := models.EventTypeAction
				if t, ok := eventData["type"].(string); ok {
					eventType = models.EventType(t)
				}
				
				eventSource := models.EventSourceAgent
				if s, ok := eventData["source"].(string); ok {
					eventSource = models.EventSource(s)
				}
				
				data := eventData["data"]
				if data == nil {
					data = make(map[string]interface{})
				}
				
				tags := make(map[string]string)
				if t, ok := eventData["tags"].(map[string]string); ok {
					tags = t
				}
				
				if _, ok := tags["agent_id"]; !ok {
					tags["agent_id"] = agent.Config.ID
				}
				
				event := models.NewEvent(
					eventType,
					eventSource,
					data,
					tags,
				)
				
				if err := e.EventStream.AddEvent(event); err != nil {
					fmt.Printf("Failed to add event to stream: %v\n", err)
				}
			}
			
			delete(outputs, "events")
		}
		
		return outputs, nil
	}
	
	return nil
}

func (e *EventIntegration) RegisterSystem(system *MultiAgentSystem) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	
	systemEventHandler := func(event *models.Event) {
		if event.Source == models.EventSourceSystem {
			if systemID, ok := event.Tags["system_id"]; ok && systemID == system.ID {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					
					var orchestratorAgent *Agent
					for _, agent := range system.ListAgents() {
						if agent.Config.Role == AgentRoleOrchestrator {
							orchestratorAgent = agent
							break
						}
					}
					
					if orchestratorAgent == nil {
						fmt.Printf("Orchestrator agent not found for system %s\n", system.ID)
						return
					}
					
					_, err := orchestratorAgent.Process(ctx, map[string]interface{}{
						"event_type": event.Type,
						"event_source": event.Source,
						"event_data": event.Data,
						"event_tags": event.Tags,
						"event_timestamp": event.Timestamp,
					})
					
					if err != nil {
						fmt.Printf("Error processing event with orchestrator agent: %v\n", err)
					}
				}()
			}
		}
	}
	
	if err := e.EventStream.Subscribe(models.EventTypeSystem, systemEventHandler); err != nil {
		return fmt.Errorf("failed to subscribe system to system events: %v", err)
	}
	
	if err := e.EventStream.Subscribe(models.EventTypeTask, systemEventHandler); err != nil {
		return fmt.Errorf("failed to subscribe system to task events: %v", err)
	}
	
	e.SystemSubscriptions[system.ID] = []func(*models.Event){systemEventHandler}
	
	event := models.NewEvent(
		models.EventTypeSystem,
		models.EventSourceSystem,
		map[string]interface{}{
			"action":      "system_registered_with_event_stream",
			"system_id":   system.ID,
			"system_name": system.Name,
		},
		map[string]string{
			"system_id": system.ID,
		},
	)
	
	if err := e.EventStream.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add event to stream: %v", err)
	}
	
	for _, agent := range system.ListAgents() {
		if err := e.RegisterAgent(agent); err != nil {
			return fmt.Errorf("failed to register agent %s: %v", agent.Config.ID, err)
		}
	}
	
	return nil
}

func (e *EventIntegration) UnregisterAgent(agentID string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	
	handlers, exists := e.AgentSubscriptions[agentID]
	if !exists {
		return fmt.Errorf("agent %s not registered with event integration", agentID)
	}
	
	for _, handler := range handlers {
		if err := e.EventStream.Unsubscribe(models.EventTypeAction, handler); err != nil {
			return fmt.Errorf("failed to unsubscribe agent from action events: %v", err)
		}
		
		if err := e.EventStream.Unsubscribe(models.EventTypeTask, handler); err != nil {
			return fmt.Errorf("failed to unsubscribe agent from task events: %v", err)
		}
		
		if err := e.EventStream.Unsubscribe(models.EventTypeState, handler); err != nil {
			return fmt.Errorf("failed to unsubscribe agent from state events: %v", err)
		}
	}
	
	delete(e.AgentSubscriptions, agentID)
	
	event := models.NewEvent(
		models.EventTypeSystem,
		models.EventSourceSystem,
		map[string]interface{}{
			"action":   "agent_unregistered_from_event_stream",
			"agent_id": agentID,
		},
		map[string]string{
			"agent_id": agentID,
		},
	)
	
	if err := e.EventStream.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add event to stream: %v", err)
	}
	
	return nil
}

func (e *EventIntegration) UnregisterSystem(systemID string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	
	handlers, exists := e.SystemSubscriptions[systemID]
	if !exists {
		return fmt.Errorf("system %s not registered with event integration", systemID)
	}
	
	for _, handler := range handlers {
		if err := e.EventStream.Unsubscribe(models.EventTypeSystem, handler); err != nil {
			return fmt.Errorf("failed to unsubscribe system from system events: %v", err)
		}
		
		if err := e.EventStream.Unsubscribe(models.EventTypeTask, handler); err != nil {
			return fmt.Errorf("failed to unsubscribe system from task events: %v", err)
		}
	}
	
	delete(e.SystemSubscriptions, systemID)
	
	event := models.NewEvent(
		models.EventTypeSystem,
		models.EventSourceSystem,
		map[string]interface{}{
			"action":    "system_unregistered_from_event_stream",
			"system_id": systemID,
		},
		map[string]string{
			"system_id": systemID,
		},
	)
	
	if err := e.EventStream.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add event to stream: %v", err)
	}
	
	return nil
}

func (e *EventIntegration) EmitEvent(eventType models.EventType, eventSource models.EventSource, data map[string]interface{}, tags map[string]string) error {
	event := models.NewEvent(
		eventType,
		eventSource,
		data,
		tags,
	)
	
	return e.EventStream.AddEvent(event)
}

func (e *EventIntegration) CreateEventAwareAgent(config AgentConfig) (*Agent, error) {
	agent, err := NewAgent(config)
	if err != nil {
		return nil, err
	}
	
	if err := e.RegisterAgent(agent); err != nil {
		return nil, err
	}
	
	return agent, nil
}

func CreateEventAwareMultiAgentSystem(name, description string, eventStream EventStream) (*MultiAgentSystem, error) {
	system := NewMultiAgentSystem(name, description, eventStream)
	
	eventIntegration := NewEventIntegration(eventStream)
	
	if err := eventIntegration.RegisterSystem(system); err != nil {
		return nil, fmt.Errorf("failed to register system with event integration: %v", err)
	}
	
	system.Metadata["event_integration"] = eventIntegration
	
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
	frontendAgent, err := eventIntegration.CreateEventAwareAgent(frontendAgentConfig)
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
	appBuilderAgent, err := eventIntegration.CreateEventAwareAgent(appBuilderAgentConfig)
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
	codegenAgent, err := eventIntegration.CreateEventAwareAgent(codegenAgentConfig)
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
	engineeringAgent, err := eventIntegration.CreateEventAwareAgent(engineeringAgentConfig)
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
	orchestratorAgent, err := eventIntegration.CreateEventAwareAgent(orchestratorAgentConfig)
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

func CreateFullyIntegratedMultiAgentSystem(name, description string, eventStream EventStream, stateProvider StateProvider) (*MultiAgentSystem, error) {
	eventIntegration := NewEventIntegration(eventStream)
	
	stateIntegration := NewStateIntegration(stateProvider, eventStream)
	
	system := NewMultiAgentSystem(name, description, eventStream)
	
	if err := eventIntegration.RegisterSystem(system); err != nil {
		return nil, fmt.Errorf("failed to register system with event integration: %v", err)
	}
	
	if err := stateIntegration.RegisterSystem(system); err != nil {
		return nil, fmt.Errorf("failed to register system with state integration: %v", err)
	}
	
	system.Metadata["event_integration"] = eventIntegration
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
	frontendAgent, err := NewAgent(frontendAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend agent: %v", err)
	}
	
	if err := eventIntegration.RegisterAgent(frontendAgent); err != nil {
		return nil, fmt.Errorf("failed to register frontend agent with event integration: %v", err)
	}
	
	if err := stateIntegration.RegisterAgent(frontendAgent); err != nil {
		return nil, fmt.Errorf("failed to register frontend agent with state integration: %v", err)
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
	appBuilderAgent, err := NewAgent(appBuilderAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create app builder agent: %v", err)
	}
	
	if err := eventIntegration.RegisterAgent(appBuilderAgent); err != nil {
		return nil, fmt.Errorf("failed to register app builder agent with event integration: %v", err)
	}
	
	if err := stateIntegration.RegisterAgent(appBuilderAgent); err != nil {
		return nil, fmt.Errorf("failed to register app builder agent with state integration: %v", err)
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
	codegenAgent, err := NewAgent(codegenAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create codegen agent: %v", err)
	}
	
	if err := eventIntegration.RegisterAgent(codegenAgent); err != nil {
		return nil, fmt.Errorf("failed to register codegen agent with event integration: %v", err)
	}
	
	if err := stateIntegration.RegisterAgent(codegenAgent); err != nil {
		return nil, fmt.Errorf("failed to register codegen agent with state integration: %v", err)
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
	engineeringAgent, err := NewAgent(engineeringAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create engineering agent: %v", err)
	}
	
	if err := eventIntegration.RegisterAgent(engineeringAgent); err != nil {
		return nil, fmt.Errorf("failed to register engineering agent with event integration: %v", err)
	}
	
	if err := stateIntegration.RegisterAgent(engineeringAgent); err != nil {
		return nil, fmt.Errorf("failed to register engineering agent with state integration: %v", err)
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
	orchestratorAgent, err := NewAgent(orchestratorAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create orchestrator agent: %v", err)
	}
	
	if err := eventIntegration.RegisterAgent(orchestratorAgent); err != nil {
		return nil, fmt.Errorf("failed to register orchestrator agent with event integration: %v", err)
	}
	
	if err := stateIntegration.RegisterAgent(orchestratorAgent); err != nil {
		return nil, fmt.Errorf("failed to register orchestrator agent with state integration: %v", err)
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
