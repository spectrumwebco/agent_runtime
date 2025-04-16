package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/tmc/langchaingo/llms"
)

type MultiAgentSystem struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Graph       *Graph                 `json:"graph"`
	Executor    *Executor              `json:"executor"`
	Agents      map[string]*Agent      `json:"agents"`
	DjangoAgents map[string]*DjangoAgent `json:"django_agents"`
	EventStream EventStream            `json:"event_stream"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	mutex       sync.RWMutex           `json:"-"`
}

func NewMultiAgentSystem(name, description string, eventStream EventStream) *MultiAgentSystem {
	id := uuid.New().String()
	graph := NewGraph(fmt.Sprintf("%s-graph", name), fmt.Sprintf("Graph for %s", name))
	executor := NewExecutor(graph)

	return &MultiAgentSystem{
		ID:          id,
		Name:        name,
		Description: description,
		Graph:       graph,
		Executor:    executor,
		Agents:      make(map[string]*Agent),
		DjangoAgents: make(map[string]*DjangoAgent),
		EventStream: eventStream,
		Metadata:    make(map[string]interface{}),
	}
}

func (s *MultiAgentSystem) AddAgent(agent *Agent) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	nodeID := fmt.Sprintf("agent-%s", agent.Config.ID)
	node, err := s.Graph.AddNode(nodeID, agent.Config.Name, agent.Config.Description, func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return agent.Process(ctx, input)
	})

	if err != nil {
		return "", fmt.Errorf("failed to add node for agent: %w", err)
	}

	s.Agents[agent.Config.ID] = agent

	return node.ID, nil
}

func (s *MultiAgentSystem) AddDjangoAgent(agent *DjangoAgent) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	nodeID := fmt.Sprintf("django-agent-%s", agent.Config.AgentConfig.ID)
	node, err := s.Graph.AddNode(nodeID, agent.Config.AgentConfig.Name, agent.Config.AgentConfig.Description, func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
		return agent.Process(ctx, input)
	})

	if err != nil {
		return "", fmt.Errorf("failed to add node for Django agent: %w", err)
	}

	s.DjangoAgents[agent.Config.AgentConfig.ID] = agent

	return node.ID, nil
}

func (s *MultiAgentSystem) GetAgent(id string) (*Agent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	agent, exists := s.Agents[id]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", id)
	}

	return agent, nil
}

func (s *MultiAgentSystem) GetDjangoAgent(id string) (*DjangoAgent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	agent, exists := s.DjangoAgents[id]
	if !exists {
		return nil, fmt.Errorf("Django agent %s not found", id)
	}

	return agent, nil
}

func (s *MultiAgentSystem) GetAgentByRole(role AgentRole) (*Agent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, agent := range s.Agents {
		if agent.Config.Role == role {
			return agent, nil
		}
	}

	return nil, fmt.Errorf("no agent found for role %s", role)
}

func (s *MultiAgentSystem) ListAgents() []*Agent {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var agents []*Agent
	for _, agent := range s.Agents {
		agents = append(agents, agent)
	}

	return agents
}

func (s *MultiAgentSystem) ListDjangoAgents() []*DjangoAgent {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var agents []*DjangoAgent
	for _, agent := range s.DjangoAgents {
		agents = append(agents, agent)
	}

	return agents
}

func (s *MultiAgentSystem) RemoveAgent(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.Agents[id]; !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	nodeID := fmt.Sprintf("agent-%s", id)
	if err := s.Graph.RemoveNode(nodeID); err != nil {
		return fmt.Errorf("failed to remove node for agent: %w", err)
	}

	delete(s.Agents, id)

	return nil
}

func (s *MultiAgentSystem) RemoveDjangoAgent(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.DjangoAgents[id]; !exists {
		return fmt.Errorf("Django agent %s not found", id)
	}

	nodeID := fmt.Sprintf("django-agent-%s", id)
	if err := s.Graph.RemoveNode(nodeID); err != nil {
		return fmt.Errorf("failed to remove node for Django agent: %w", err)
	}

	delete(s.DjangoAgents, id)

	return nil
}

func (s *MultiAgentSystem) ConnectAgents(sourceID, targetID, name, description string) (*Edge, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sourceAgent, exists := s.Agents[sourceID]
	if !exists {
		if _, exists := s.DjangoAgents[sourceID]; !exists {
			return nil, fmt.Errorf("source agent %s not found", sourceID)
		}
		sourceNodeID := fmt.Sprintf("django-agent-%s", sourceID)
		targetNodeID := fmt.Sprintf("agent-%s", targetID)
		return s.Graph.AddEdge(sourceNodeID, targetNodeID, name, description)
	}

	targetAgent, exists := s.Agents[targetID]
	if !exists {
		if _, exists := s.DjangoAgents[targetID]; !exists {
			return nil, fmt.Errorf("target agent %s not found", targetID)
		}
		sourceNodeID := fmt.Sprintf("agent-%s", sourceID)
		targetNodeID := fmt.Sprintf("django-agent-%s", targetID)
		return s.Graph.AddEdge(sourceNodeID, targetNodeID, name, description)
	}

	sourceNodeID := fmt.Sprintf("agent-%s", sourceAgent.Config.ID)
	targetNodeID := fmt.Sprintf("agent-%s", targetAgent.Config.ID)
	return s.Graph.AddEdge(sourceNodeID, targetNodeID, name, description)
}

func (s *MultiAgentSystem) Execute(ctx context.Context, agentID string, input map[string]interface{}) (*Execution, error) {
	agent, err := s.GetAgent(agentID)
	if err != nil {
		djangoAgent, err := s.GetDjangoAgent(agentID)
		if err != nil {
			return nil, fmt.Errorf("agent %s not found", agentID)
		}
		nodeID := fmt.Sprintf("django-agent-%s", djangoAgent.Config.AgentConfig.ID)
		return s.Executor.Execute(ctx, nodeID, input)
	}

	nodeID := fmt.Sprintf("agent-%s", agent.Config.ID)
	return s.Executor.Execute(ctx, nodeID, input)
}

func (s *MultiAgentSystem) SetDjangoFactory(factory *DjangoAgentFactory) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Metadata["django_factory"] = factory
}

func (s *MultiAgentSystem) GetDjangoFactory() (*DjangoAgentFactory, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	factory, exists := s.Metadata["django_factory"].(*DjangoAgentFactory)
	if !exists {
		return nil, fmt.Errorf("Django factory not found")
	}

	return factory, nil
}

func CreateMultiAgentSystemWithLangChain(name, description string, llm llms.LLM, eventStream EventStream) (*MultiAgentSystem, *LangChainBridge, error) {
	system := NewMultiAgentSystem(name, description, eventStream)

	bridge := NewLangChainBridge(system.Graph, system.Executor, llm, eventStream)

	return system, bridge, nil
}

func CreateStandardMultiAgentSystem(name, description string, eventStream EventStream) (*MultiAgentSystem, error) {
	system := NewMultiAgentSystem(name, description, eventStream)

	frontendAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Frontend Agent",
		Description: "Agent responsible for frontend development tasks",
		Role:        AgentRoleFrontend,
		Capabilities: []AgentCapability{
			AgentCapabilityUIDesign,
			AgentCapabilityCodeGeneration,
		},
		Metadata:    make(map[string]interface{}),
		ModelConfig: make(map[string]interface{}),
		ToolConfig:  make(map[string]interface{}),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	frontendAgent := NewAgent(frontendAgentConfig, func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"message": "Frontend agent processed input",
			"input":   input,
		}, nil
	}, eventStream)

	if _, err := system.AddAgent(frontendAgent); err != nil {
		return nil, fmt.Errorf("failed to add frontend agent: %w", err)
	}

	appBuilderAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "App Builder Agent",
		Description: "Agent responsible for application building tasks",
		Role:        AgentRoleAppBuilder,
		Capabilities: []AgentCapability{
			AgentCapabilityAPIDesign,
			AgentCapabilityDatabaseDesign,
			AgentCapabilityCodeGeneration,
		},
		Metadata:    make(map[string]interface{}),
		ModelConfig: make(map[string]interface{}),
		ToolConfig:  make(map[string]interface{}),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	appBuilderAgent := NewAgent(appBuilderAgentConfig, func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"message": "App builder agent processed input",
			"input":   input,
		}, nil
	}, eventStream)

	if _, err := system.AddAgent(appBuilderAgent); err != nil {
		return nil, fmt.Errorf("failed to add app builder agent: %w", err)
	}

	codegenAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Codegen Agent",
		Description: "Agent responsible for code generation tasks",
		Role:        AgentRoleCodegen,
		Capabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
		},
		Metadata:    make(map[string]interface{}),
		ModelConfig: make(map[string]interface{}),
		ToolConfig:  make(map[string]interface{}),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	codegenAgent := NewAgent(codegenAgentConfig, func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"message": "Codegen agent processed input",
			"input":   input,
		}, nil
	}, eventStream)

	if _, err := system.AddAgent(codegenAgent); err != nil {
		return nil, fmt.Errorf("failed to add codegen agent: %w", err)
	}

	engineeringAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Engineering Agent",
		Description: "Agent responsible for software engineering tasks",
		Role:        AgentRoleEngineering,
		Capabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
			AgentCapabilityTesting,
			AgentCapabilityDeployment,
			AgentCapabilityDocumentation,
		},
		Metadata:    make(map[string]interface{}),
		ModelConfig: make(map[string]interface{}),
		ToolConfig:  make(map[string]interface{}),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	engineeringAgent := NewAgent(engineeringAgentConfig, func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"message": "Engineering agent processed input",
			"input":   input,
		}, nil
	}, eventStream)

	if _, err := system.AddAgent(engineeringAgent); err != nil {
		return nil, fmt.Errorf("failed to add engineering agent: %w", err)
	}

	orchestratorAgentConfig := AgentConfig{
		ID:          uuid.New().String(),
		Name:        "Orchestrator Agent",
		Description: "Agent responsible for orchestrating other agents",
		Role:        AgentRoleOrchestrator,
		Capabilities: []AgentCapability{
			AgentCapabilityOrchestration,
			AgentCapabilityPlanning,
		},
		Metadata:    make(map[string]interface{}),
		ModelConfig: make(map[string]interface{}),
		ToolConfig:  make(map[string]interface{}),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	orchestratorAgent := NewAgent(orchestratorAgentConfig, func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"message": "Orchestrator agent processed input",
			"input":   input,
		}, nil
	}, eventStream)

	if _, err := system.AddAgent(orchestratorAgent); err != nil {
		return nil, fmt.Errorf("failed to add orchestrator agent: %w", err)
	}

	if _, err := system.ConnectAgents(orchestratorAgent.Config.ID, frontendAgent.Config.ID, "Orchestrator -> Frontend", "Orchestrator to Frontend communication"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to frontend agent: %w", err)
	}

	if _, err := system.ConnectAgents(orchestratorAgent.Config.ID, appBuilderAgent.Config.ID, "Orchestrator -> App Builder", "Orchestrator to App Builder communication"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to app builder agent: %w", err)
	}

	if _, err := system.ConnectAgents(orchestratorAgent.Config.ID, codegenAgent.Config.ID, "Orchestrator -> Codegen", "Orchestrator to Codegen communication"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to codegen agent: %w", err)
	}

	if _, err := system.ConnectAgents(orchestratorAgent.Config.ID, engineeringAgent.Config.ID, "Orchestrator -> Engineering", "Orchestrator to Engineering communication"); err != nil {
		return nil, fmt.Errorf("failed to connect orchestrator to engineering agent: %w", err)
	}

	if _, err := system.ConnectAgents(frontendAgent.Config.ID, codegenAgent.Config.ID, "Frontend -> Codegen", "Frontend to Codegen communication"); err != nil {
		return nil, fmt.Errorf("failed to connect frontend to codegen agent: %w", err)
	}

	if _, err := system.ConnectAgents(appBuilderAgent.Config.ID, codegenAgent.Config.ID, "App Builder -> Codegen", "App Builder to Codegen communication"); err != nil {
		return nil, fmt.Errorf("failed to connect app builder to codegen agent: %w", err)
	}

	if _, err := system.ConnectAgents(codegenAgent.Config.ID, engineeringAgent.Config.ID, "Codegen -> Engineering", "Codegen to Engineering communication"); err != nil {
		return nil, fmt.Errorf("failed to connect codegen to engineering agent: %w", err)
	}

	return system, nil
}
