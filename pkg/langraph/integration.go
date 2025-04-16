package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type MultiAgentSystem struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	Graph         *Graph                 `json:"graph"`
	Executor      *Executor              `json:"executor"`
	AgentFactory  *AgentFactory          `json:"agent_factory"`
	DjangoFactory *DjangoAgentFactory    `json:"django_factory,omitempty"`
	EventStream   EventStream            `json:"-"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	agents        map[string]*Agent      `json:"-"`
	djangoAgents  map[string]*DjangoAgent `json:"-"`
	agentsMutex   sync.RWMutex           `json:"-"`
}

func NewMultiAgentSystem(name, description string, eventStream EventStream) *MultiAgentSystem {
	graph := NewGraph(name, description)
	executor := NewExecutor(graph)
	agentFactory := NewAgentFactory(eventStream)

	return &MultiAgentSystem{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Graph:       graph,
		Executor:    executor,
		AgentFactory: agentFactory,
		EventStream: eventStream,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		agents:      make(map[string]*Agent),
		djangoAgents: make(map[string]*DjangoAgent),
	}
}

func (s *MultiAgentSystem) SetDjangoFactory(djangoFactory *DjangoAgentFactory) {
	s.DjangoFactory = djangoFactory
}

func (s *MultiAgentSystem) AddAgent(agent *Agent) (*Node, error) {
	s.agentsMutex.Lock()
	defer s.agentsMutex.Unlock()

	s.agents[agent.Config.ID] = agent

	node := s.Graph.AddAgentNode(
		AgentType(agent.Config.Role),
		agent.Config.Name,
		agent.Config.Description,
		func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			return agent.Process(ctx, inputs)
		},
	)

	node.Metadata["agent_id"] = agent.Config.ID
	node.Metadata["agent_role"] = string(agent.Config.Role)
	node.Metadata["agent_capabilities"] = agent.Config.Capabilities

	s.UpdatedAt = time.Now().UTC()

	return node, nil
}

func (s *MultiAgentSystem) AddDjangoAgent(agent *DjangoAgent) (*Node, error) {
	s.agentsMutex.Lock()
	defer s.agentsMutex.Unlock()

	s.djangoAgents[agent.Config.AgentConfig.ID] = agent

	node := s.Graph.AddAgentNode(
		AgentType(agent.Config.AgentConfig.Role),
		agent.Config.AgentConfig.Name,
		agent.Config.AgentConfig.Description,
		func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			return agent.Process(ctx, inputs)
		},
	)

	node.Metadata["agent_id"] = agent.Config.AgentConfig.ID
	node.Metadata["agent_role"] = string(agent.Config.AgentConfig.Role)
	node.Metadata["agent_type"] = "django"

	s.UpdatedAt = time.Now().UTC()

	return node, nil
}

func (s *MultiAgentSystem) GetAgent(id string) (*Agent, error) {
	s.agentsMutex.RLock()
	defer s.agentsMutex.RUnlock()

	agent, exists := s.agents[id]
	if !exists {
		djangoAgent, exists := s.djangoAgents[id]
		if !exists {
			return nil, fmt.Errorf("agent %s does not exist", id)
		}
		
		return nil, fmt.Errorf("agent %s is a Django agent and cannot be directly accessed", id)
	}

	return agent, nil
}

func (s *MultiAgentSystem) GetDjangoAgent(id string) (*DjangoAgent, error) {
	s.agentsMutex.RLock()
	defer s.agentsMutex.RUnlock()

	agent, exists := s.djangoAgents[id]
	if !exists {
		return nil, fmt.Errorf("Django agent %s does not exist", id)
	}

	return agent, nil
}

func (s *MultiAgentSystem) ConnectAgents(sourceAgentID, targetAgentID, name, description string) (*Edge, error) {
	var sourceNodeID NodeID
	var targetNodeID NodeID
	
	for _, node := range s.Graph.Nodes {
		if agentID, ok := node.Metadata["agent_id"].(string); ok && agentID == sourceAgentID {
			sourceNodeID = node.ID
			break
		}
	}
	
	for _, node := range s.Graph.Nodes {
		if agentID, ok := node.Metadata["agent_id"].(string); ok && agentID == targetAgentID {
			targetNodeID = node.ID
			break
		}
	}
	
	if sourceNodeID == "" {
		return nil, fmt.Errorf("source agent %s not found in graph", sourceAgentID)
	}
	
	if targetNodeID == "" {
		return nil, fmt.Errorf("target agent %s not found in graph", targetAgentID)
	}

	edge, err := s.Graph.AddEdge(sourceNodeID, targetNodeID, name, description)
	if err != nil {
		return nil, err
	}

	s.UpdatedAt = time.Now().UTC()

	return edge, nil
}

func (s *MultiAgentSystem) Execute(ctx context.Context, startAgentID string, input map[string]interface{}) (*Execution, error) {
	var startNodeID NodeID
	
	for _, node := range s.Graph.Nodes {
		if agentID, ok := node.Metadata["agent_id"].(string); ok && agentID == startAgentID {
			startNodeID = node.ID
			break
		}
	}
	
	if startNodeID == "" {
		return nil, fmt.Errorf("start agent %s not found in graph", startAgentID)
	}

	metadata := map[string]interface{}{
		"system_id":   s.ID,
		"system_name": s.Name,
		"start_agent": startAgentID,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	execution, err := s.Executor.Execute(ctx, startNodeID, input, metadata)
	if err != nil {
		return nil, err
	}

	return execution, nil
}

func CreateStandardAgentSystem(name, description string, eventStream EventStream, djangoBaseURL, djangoAPIKey string) (*MultiAgentSystem, error) {
	system := NewMultiAgentSystem(name, description, eventStream)

	agentFactory := system.AgentFactory

	if djangoBaseURL != "" && djangoAPIKey != "" {
		djangoFactory := NewDjangoAgentFactory(djangoBaseURL, djangoAPIKey, eventStream)
		system.SetDjangoFactory(djangoFactory)
	}

	frontendAgent, err := agentFactory.CreateFrontendAgent("Frontend Agent", "Handles UI/UX development")
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend agent: %w", err)
	}
	frontendNode, err := system.AddAgent(frontendAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to add frontend agent to system: %w", err)
	}

	appBuilderAgent, err := agentFactory.CreateAppBuilderAgent("App Builder Agent", "Handles application assembly")
	if err != nil {
		return nil, fmt.Errorf("failed to create app builder agent: %w", err)
	}
	appBuilderNode, err := system.AddAgent(appBuilderAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to add app builder agent to system: %w", err)
	}

	codegenAgent, err := agentFactory.CreateCodegenAgent("Codegen Agent", "Handles code generation")
	if err != nil {
		return nil, fmt.Errorf("failed to create codegen agent: %w", err)
	}
	codegenNode, err := system.AddAgent(codegenAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to add codegen agent to system: %w", err)
	}

	engineeringAgent, err := agentFactory.CreateEngineeringAgent("Engineering Agent", "Handles core development tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to create engineering agent: %w", err)
	}
	engineeringNode, err := system.AddAgent(engineeringAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to add engineering agent to system: %w", err)
	}

	orchestratorAgent, err := agentFactory.CreateOrchestratorAgent("Orchestrator Agent", "Handles task orchestration")
	if err != nil {
		return nil, fmt.Errorf("failed to create orchestrator agent: %w", err)
	}
	orchestratorNode, err := system.AddAgent(orchestratorAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to add orchestrator agent to system: %w", err)
	}

	_, err = system.Graph.AddEdge(orchestratorNode.ID, frontendNode.ID, "Orchestrator -> Frontend", "Orchestrator to Frontend communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to frontend: %w", err)
	}
	
	_, err = system.Graph.AddEdge(orchestratorNode.ID, appBuilderNode.ID, "Orchestrator -> App Builder", "Orchestrator to App Builder communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to app builder: %w", err)
	}
	
	_, err = system.Graph.AddEdge(orchestratorNode.ID, codegenNode.ID, "Orchestrator -> Codegen", "Orchestrator to Codegen communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to codegen: %w", err)
	}
	
	_, err = system.Graph.AddEdge(orchestratorNode.ID, engineeringNode.ID, "Orchestrator -> Engineering", "Orchestrator to Engineering communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to engineering: %w", err)
	}

	_, err = system.Graph.AddEdge(frontendNode.ID, appBuilderNode.ID, "Frontend -> App Builder", "Frontend to App Builder communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect frontend to app builder: %w", err)
	}
	
	_, err = system.Graph.AddEdge(frontendNode.ID, codegenNode.ID, "Frontend -> Codegen", "Frontend to Codegen communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect frontend to codegen: %w", err)
	}

	_, err = system.Graph.AddEdge(appBuilderNode.ID, codegenNode.ID, "App Builder -> Codegen", "App Builder to Codegen communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect app builder to codegen: %w", err)
	}
	
	_, err = system.Graph.AddEdge(appBuilderNode.ID, engineeringNode.ID, "App Builder -> Engineering", "App Builder to Engineering communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect app builder to engineering: %w", err)
	}

	_, err = system.Graph.AddEdge(codegenNode.ID, engineeringNode.ID, "Codegen -> Engineering", "Codegen to Engineering communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect codegen to engineering: %w", err)
	}

	_, err = system.Graph.AddEdge(engineeringNode.ID, frontendNode.ID, "Engineering -> Frontend", "Engineering to Frontend communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect engineering to frontend: %w", err)
	}
	
	_, err = system.Graph.AddEdge(engineeringNode.ID, appBuilderNode.ID, "Engineering -> App Builder", "Engineering to App Builder communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect engineering to app builder: %w", err)
	}
	
	_, err = system.Graph.AddEdge(engineeringNode.ID, codegenNode.ID, "Engineering -> Codegen", "Engineering to Codegen communication")
	if err != nil {
		return nil, fmt.Errorf("failed to connect engineering to codegen: %w", err)
	}

	return system, nil
}
