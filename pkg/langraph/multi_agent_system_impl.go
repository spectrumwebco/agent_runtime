package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type EventStream interface {
	AddEvent(event *models.Event) error
	Subscribe(eventType models.EventType, callback func(*models.Event)) error
	Unsubscribe(eventType models.EventType, callback func(*models.Event)) error
}

type MultiAgentSystem struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Orchestrator    *Orchestrator          `json:"orchestrator"`
	Graph           *Graph                 `json:"graph"`
	Executor        *Executor              `json:"executor"`
	EventStream     EventStream            `json:"-"`
	DjangoFactory   *DjangoAgentFactory    `json:"-"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	lock            sync.RWMutex           `json:"-"`
}

func NewMultiAgentSystem(name, description string, eventStream EventStream) *MultiAgentSystem {
	graph := NewGraph(name, description)
	executor := NewExecutor(graph)
	
	orchestratorConfig := OrchestratorConfig{
		Name:               name,
		Description:        description,
		MaxConcurrentTasks: 10,
		DefaultTimeout:     5 * time.Minute,
		EventStream:        eventStream,
		Metadata:           make(map[string]interface{}),
	}
	
	orchestrator := NewOrchestrator(orchestratorConfig)
	
	return &MultiAgentSystem{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Orchestrator: orchestrator,
		Graph:       graph,
		Executor:    executor,
		EventStream: eventStream,
		Metadata:    make(map[string]interface{}),
	}
}

func (m *MultiAgentSystem) SetDjangoFactory(factory *DjangoAgentFactory) {
	m.lock.Lock()
	defer m.lock.Unlock()
	
	m.DjangoFactory = factory
}

func (m *MultiAgentSystem) AddAgent(agent *Agent) {
	m.lock.Lock()
	defer m.lock.Unlock()
	
	m.Orchestrator.RegisterAgent(agent)
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "agent_added",
				"agent_id":   agent.Config.ID,
				"agent_name": agent.Config.Name,
				"agent_role": agent.Config.Role,
				"system_id":  m.ID,
			},
			map[string]string{
				"agent_id":  agent.Config.ID,
				"system_id": m.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
}

func (m *MultiAgentSystem) AddDjangoAgent(role AgentRole, name, description string) (*Agent, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	
	if m.DjangoFactory == nil {
		return nil, fmt.Errorf("Django agent factory not set")
	}
	
	agent, err := m.DjangoFactory.CreateAgent(role, name, description)
	if err != nil {
		return nil, err
	}
	
	m.Orchestrator.RegisterAgent(agent)
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "django_agent_added",
				"agent_id":   agent.Config.ID,
				"agent_name": agent.Config.Name,
				"agent_role": agent.Config.Role,
				"system_id":  m.ID,
			},
			map[string]string{
				"agent_id":  agent.Config.ID,
				"system_id": m.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return agent, nil
}

func (m *MultiAgentSystem) GetAgent(agentID string) (*Agent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	
	return m.Orchestrator.GetAgent(agentID)
}

func (m *MultiAgentSystem) GetAgentByRole(role AgentRole) (*Agent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	
	agents, err := m.Orchestrator.GetAgentsByRole(role)
	if err != nil {
		return nil, err
	}
	
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents found for role %s", role)
	}
	
	return agents[0], nil
}

func (m *MultiAgentSystem) ListAgents() []*Agent {
	m.lock.RLock()
	defer m.lock.RUnlock()
	
	var agents []*Agent
	for _, agent := range m.Orchestrator.Agents {
		agents = append(agents, agent)
	}
	
	return agents
}

func (m *MultiAgentSystem) ConnectAgents(sourceAgentID, targetAgentID string, name, description string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	
	sourceAgent, err := m.Orchestrator.GetAgent(sourceAgentID)
	if err != nil {
		return err
	}
	
	targetAgent, err := m.Orchestrator.GetAgent(targetAgentID)
	if err != nil {
		return err
	}
	
	channel, err := m.Orchestrator.CommunicationMgr.CreateChannel(context.Background(), sourceAgent, targetAgent, name, description)
	if err != nil {
		return err
	}
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":           "agents_connected",
				"source_agent_id":  sourceAgent.Config.ID,
				"target_agent_id":  targetAgent.Config.ID,
				"channel_id":       channel.ID,
				"channel_name":     channel.Name,
				"system_id":        m.ID,
			},
			map[string]string{
				"source_agent_id": sourceAgent.Config.ID,
				"target_agent_id": targetAgent.Config.ID,
				"system_id":       m.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return nil
}

func (m *MultiAgentSystem) Execute(ctx context.Context, agentID string, inputs map[string]interface{}) (*Execution, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	
	agent, err := m.Orchestrator.GetAgent(agentID)
	if err != nil {
		return nil, err
	}
	
	var agentNodeID NodeID
	for id, node := range m.Graph.Nodes {
		if node.Type == NodeTypeAgent && node.Metadata["agent_id"] == agent.Config.ID {
			agentNodeID = id
			break
		}
	}
	
	if agentNodeID == "" {
		return nil, fmt.Errorf("node for agent %s not found", agentID)
	}
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "execution_started",
				"agent_id":   agent.Config.ID,
				"agent_name": agent.Config.Name,
				"agent_role": agent.Config.Role,
				"system_id":  m.ID,
				"inputs":     inputs,
			},
			map[string]string{
				"agent_id":  agent.Config.ID,
				"system_id": m.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	execution, err := m.Executor.Execute(ctx, agentNodeID, inputs, map[string]interface{}{
		"agent_id":   agent.Config.ID,
		"agent_name": agent.Config.Name,
		"agent_role": agent.Config.Role,
		"system_id":  m.ID,
	})
	
	if err != nil {
		if m.EventStream != nil {
			event := models.NewEvent(
				models.EventTypeSystem,
				models.EventSourceSystem,
				map[string]interface{}{
					"action":     "execution_failed",
					"agent_id":   agent.Config.ID,
					"agent_name": agent.Config.Name,
					"agent_role": agent.Config.Role,
					"system_id":  m.ID,
					"error":      err.Error(),
				},
				map[string]string{
					"agent_id":  agent.Config.ID,
					"system_id": m.ID,
				},
			)
			
			if err := m.EventStream.AddEvent(event); err != nil {
				fmt.Printf("Failed to add event to stream: %v\n", err)
			}
		}
		
		return nil, err
	}
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":       "execution_completed",
				"agent_id":     agent.Config.ID,
				"agent_name":   agent.Config.Name,
				"agent_role":   agent.Config.Role,
				"system_id":    m.ID,
				"execution_id": execution.ID,
			},
			map[string]string{
				"agent_id":     agent.Config.ID,
				"system_id":    m.ID,
				"execution_id": execution.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return execution, nil
}

func CreateStandardMultiAgentSystem(name, description string, eventStream EventStream) (*MultiAgentSystem, error) {
	system := NewMultiAgentSystem(name, description, eventStream)
	
	if err := system.Orchestrator.CreateStandardAgents(); err != nil {
		return nil, fmt.Errorf("failed to create standard agents: %v", err)
	}
	
	for _, agent := range system.Orchestrator.Agents {
		node := system.Graph.AddAgentNode(AgentType(agent.Config.Role), agent.Config.Name, agent.Config.Description, func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			return agent.Process(ctx, inputs)
		})
		
		node.Metadata["agent_id"] = agent.Config.ID
	}
	
	for _, channel := range system.Orchestrator.CommunicationMgr.Channels {
		var sourceNodeID, targetNodeID NodeID
		
		for id, node := range system.Graph.Nodes {
			if agentID, ok := node.Metadata["agent_id"].(string); ok {
				if agentID == channel.SourceAgent.Config.ID {
					sourceNodeID = id
				} else if agentID == channel.TargetAgent.Config.ID {
					targetNodeID = id
				}
			}
		}
		
		if sourceNodeID != "" && targetNodeID != "" {
			_, err := system.Graph.AddEdge(sourceNodeID, targetNodeID, channel.Name, channel.Description)
			if err != nil {
				return nil, fmt.Errorf("failed to add edge between agents: %v", err)
			}
		}
	}
	
	return system, nil
}

func (m *MultiAgentSystem) GetAgentsByCapability(capability AgentCapability) []*Agent {
	m.lock.RLock()
	defer m.lock.RUnlock()
	
	var agents []*Agent
	for _, agent := range m.Orchestrator.Agents {
		if agent.HasCapability(capability) {
			agents = append(agents, agent)
		}
	}
	
	return agents
}

func (m *MultiAgentSystem) CreateWorkflow(name, description string) (*Graph, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	
	workflow := NewGraph(name, description)
	workflow.Metadata["system_id"] = m.ID
	workflow.Metadata["type"] = "workflow"
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":       "workflow_created",
				"workflow_id":  workflow.ID,
				"workflow_name": name,
				"system_id":    m.ID,
			},
			map[string]string{
				"workflow_id": workflow.ID,
				"system_id":   m.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return workflow, nil
}

func (m *MultiAgentSystem) ExecuteWorkflow(ctx context.Context, workflow *Graph, startNodeID NodeID, inputs map[string]interface{}) (*Execution, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	
	executor := NewExecutor(workflow)
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":       "workflow_execution_started",
				"workflow_id":  workflow.ID,
				"workflow_name": workflow.Name,
				"system_id":    m.ID,
				"inputs":       inputs,
			},
			map[string]string{
				"workflow_id": workflow.ID,
				"system_id":   m.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	execution, err := executor.Execute(ctx, startNodeID, inputs, map[string]interface{}{
		"workflow_id":  workflow.ID,
		"workflow_name": workflow.Name,
		"system_id":    m.ID,
	})
	
	if err != nil {
		if m.EventStream != nil {
			event := models.NewEvent(
				models.EventTypeSystem,
				models.EventSourceSystem,
				map[string]interface{}{
					"action":       "workflow_execution_failed",
					"workflow_id":  workflow.ID,
					"workflow_name": workflow.Name,
					"system_id":    m.ID,
					"error":        err.Error(),
				},
				map[string]string{
					"workflow_id": workflow.ID,
					"system_id":   m.ID,
				},
			)
			
			if err := m.EventStream.AddEvent(event); err != nil {
				fmt.Printf("Failed to add event to stream: %v\n", err)
			}
		}
		
		return nil, err
	}
	
	if m.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":       "workflow_execution_completed",
				"workflow_id":  workflow.ID,
				"workflow_name": workflow.Name,
				"system_id":    m.ID,
				"execution_id": execution.ID,
			},
			map[string]string{
				"workflow_id":  workflow.ID,
				"system_id":    m.ID,
				"execution_id": execution.ID,
			},
		)
		
		if err := m.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return execution, nil
}

func (m *MultiAgentSystem) CreateStructuredOutputHandler(schema map[string]interface{}, pollingInterval time.Duration) func(ctx context.Context, agent *Agent, inputs map[string]interface{}) (map[string]interface{}, error) {
	return func(ctx context.Context, agent *Agent, inputs map[string]interface{}) (map[string]interface{}, error) {
		outputs, err := agent.Process(ctx, inputs)
		if err != nil {
			return nil, err
		}
		
		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-ticker.C:
				if structuredOutput, ok := outputs["structured_output"].(map[string]interface{}); ok {
					return structuredOutput, nil
				}
			}
		}
	}
}
