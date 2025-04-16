package kledframework

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/agent"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/types"
)

type MultiAgentManager struct {
	config         *MultiAgentManagerConfig
	agents         map[string]*agent.Agent
	agentsMutex    sync.RWMutex
	graphExecutor  *langraph.Executor
	multiAgentSys  *langraph.MultiAgentSystem
	isRunning      bool
	statusMutex    sync.RWMutex
	stopCh         chan struct{}
	eventListeners []AgentEventListener
	listenerMutex  sync.RWMutex
}

type MultiAgentManagerConfig struct {
	DefaultAgentConfig map[string]interface{}
	AgentTypes         []types.AgentType
	GraphConfig        *langraph.MultiAgentSystemConfig
	EnableOrchestrator bool
	MaxConcurrentTasks int
	TaskTimeout        time.Duration
}

type AgentEventListener func(event *types.AgentEvent)

func DefaultMultiAgentManagerConfig() *MultiAgentManagerConfig {
	return &MultiAgentManagerConfig{
		DefaultAgentConfig: map[string]interface{}{
			"max_retries": 3,
			"timeout":     30 * time.Second,
		},
		AgentTypes: []types.AgentType{
			types.AgentTypeFrontend,
			types.AgentTypeAppBuilder,
			types.AgentTypeCodegen,
			types.AgentTypeEngineering,
		},
		GraphConfig: &langraph.MultiAgentSystemConfig{
			MaxConcurrentNodes: 10,
			DefaultTimeout:     60 * time.Second,
			EnableLogging:      true,
		},
		EnableOrchestrator: true,
		MaxConcurrentTasks: 5,
		TaskTimeout:        5 * time.Minute,
	}
}

func NewMultiAgentManager(config *MultiAgentManagerConfig) (*MultiAgentManager, error) {
	if config == nil {
		config = DefaultMultiAgentManagerConfig()
	}

	graphExecutor := langraph.NewExecutor()

	manager := &MultiAgentManager{
		config:        config,
		agents:        make(map[string]*agent.Agent),
		graphExecutor: graphExecutor,
		isRunning:     false,
		stopCh:        make(chan struct{}),
	}

	multiAgentSys, err := langraph.NewMultiAgentSystem(config.GraphConfig, graphExecutor)
	if err != nil {
		return nil, fmt.Errorf("failed to create multi-agent system: %w", err)
	}
	manager.multiAgentSys = multiAgentSys

	return manager, nil
}

func (m *MultiAgentManager) Start(ctx context.Context) error {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	if m.isRunning {
		return fmt.Errorf("multi-agent manager is already running")
	}

	if err := m.graphExecutor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start graph executor: %w", err)
	}

	if err := m.multiAgentSys.Start(ctx); err != nil {
		return fmt.Errorf("failed to start multi-agent system: %w", err)
	}

	m.agentsMutex.RLock()
	for _, agent := range m.agents {
		if err := agent.Start(ctx); err != nil {
			m.agentsMutex.RUnlock()
			return fmt.Errorf("failed to start agent %s: %w", agent.ID(), err)
		}
	}
	m.agentsMutex.RUnlock()

	m.isRunning = true
	go m.run(ctx)

	log.Printf("Multi-agent manager started with %d agents", len(m.agents))

	return nil
}

func (m *MultiAgentManager) Stop() error {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	if !m.isRunning {
		return nil
	}

	close(m.stopCh)

	m.agentsMutex.RLock()
	for _, agent := range m.agents {
		if err := agent.Stop(); err != nil {
			log.Printf("Failed to stop agent %s: %v", agent.ID(), err)
		}
	}
	m.agentsMutex.RUnlock()

	if err := m.multiAgentSys.Stop(); err != nil {
		log.Printf("Failed to stop multi-agent system: %v", err)
	}

	if err := m.graphExecutor.Stop(); err != nil {
		log.Printf("Failed to stop graph executor: %v", err)
	}

	m.isRunning = false

	log.Printf("Multi-agent manager stopped")

	return nil
}

func (m *MultiAgentManager) CreateAgent(agentType types.AgentType, name string, config map[string]interface{}) (*agent.Agent, error) {
	id := uuid.New().String()

	mergedConfig := make(map[string]interface{})
	for k, v := range m.config.DefaultAgentConfig {
		mergedConfig[k] = v
	}
	for k, v := range config {
		mergedConfig[k] = v
	}

	agentConfig := &agent.Config{
		ID:           id,
		Name:         name,
		Type:         string(agentType),
		Config:       mergedConfig,
		Capabilities: []string{},
	}

	switch agentType {
	case types.AgentTypeFrontend:
		agentConfig.Capabilities = append(agentConfig.Capabilities, "ui_design", "user_interaction", "frontend_development")
	case types.AgentTypeAppBuilder:
		agentConfig.Capabilities = append(agentConfig.Capabilities, "app_scaffolding", "component_integration", "deployment")
	case types.AgentTypeCodegen:
		agentConfig.Capabilities = append(agentConfig.Capabilities, "code_generation", "refactoring", "optimization")
	case types.AgentTypeEngineering:
		agentConfig.Capabilities = append(agentConfig.Capabilities, "architecture_design", "system_integration", "performance_tuning")
	case types.AgentTypeOrchestrator:
		agentConfig.Capabilities = append(agentConfig.Capabilities, "task_coordination", "agent_management", "workflow_optimization")
	}

	agent, err := agent.NewAgent(agentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	if err := m.multiAgentSys.RegisterAgent(id, agent); err != nil {
		return nil, fmt.Errorf("failed to register agent with multi-agent system: %w", err)
	}

	m.agentsMutex.Lock()
	m.agents[id] = agent
	m.agentsMutex.Unlock()

	m.emitEvent(&types.AgentEvent{
		ID:        uuid.New().String(),
		AgentID:   id,
		Type:      "agent_created",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_type": agentType,
			"agent_name": name,
		},
	})

	return agent, nil
}

func (m *MultiAgentManager) GetAgent(id string) (*agent.Agent, error) {
	m.agentsMutex.RLock()
	defer m.agentsMutex.RUnlock()

	agent, exists := m.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", id)
	}

	return agent, nil
}

func (m *MultiAgentManager) ListAgents() []*agent.Agent {
	m.agentsMutex.RLock()
	defer m.agentsMutex.RUnlock()

	agents := make([]*agent.Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}

	return agents
}

func (m *MultiAgentManager) RemoveAgent(id string) error {
	m.agentsMutex.Lock()
	defer m.agentsMutex.Unlock()

	agent, exists := m.agents[id]
	if !exists {
		return fmt.Errorf("agent %s not found", id)
	}

	if err := agent.Stop(); err != nil {
		return fmt.Errorf("failed to stop agent: %w", err)
	}

	if err := m.multiAgentSys.UnregisterAgent(id); err != nil {
		return fmt.Errorf("failed to unregister agent from multi-agent system: %w", err)
	}

	delete(m.agents, id)

	m.emitEvent(&types.AgentEvent{
		ID:        uuid.New().String(),
		AgentID:   id,
		Type:      "agent_removed",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	})

	return nil
}

func (m *MultiAgentManager) ExecuteTask(ctx context.Context, task *types.AgentTask) (map[string]interface{}, error) {
	m.statusMutex.RLock()
	if !m.isRunning {
		m.statusMutex.RUnlock()
		return nil, fmt.Errorf("multi-agent manager is not running")
	}
	m.statusMutex.RUnlock()

	taskCtx, cancel := context.WithTimeout(ctx, m.config.TaskTimeout)
	defer cancel()

	graph, err := m.createTaskGraph(task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task graph: %w", err)
	}

	result, err := m.graphExecutor.ExecuteGraph(taskCtx, graph)
	if err != nil {
		return nil, fmt.Errorf("failed to execute task graph: %w", err)
	}

	m.emitEvent(&types.AgentEvent{
		ID:        uuid.New().String(),
		AgentID:   task.AgentID,
		Type:      "task_completed",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"task_id":     task.ID,
			"task_type":   task.Type,
			"task_result": result,
		},
	})

	return result, nil
}

func (m *MultiAgentManager) AddEventListener(listener AgentEventListener) {
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()

	m.eventListeners = append(m.eventListeners, listener)
}

func (m *MultiAgentManager) RemoveEventListener(listener AgentEventListener) {
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()

	for i, l := range m.eventListeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			m.eventListeners = append(m.eventListeners[:i], m.eventListeners[i+1:]...)
			break
		}
	}
}

func (m *MultiAgentManager) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.Stop()
			return
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.checkAgentHealth()
		}
	}
}

func (m *MultiAgentManager) checkAgentHealth() {
	m.agentsMutex.RLock()
	defer m.agentsMutex.RUnlock()

	for _, agent := range m.agents {
		if !agent.IsRunning() {
			log.Printf("Agent %s is not running, attempting to restart", agent.ID())
			if err := agent.Start(context.Background()); err != nil {
				log.Printf("Failed to restart agent %s: %v", agent.ID(), err)
			}
		}
	}
}

func (m *MultiAgentManager) createTaskGraph(task *types.AgentTask) (*langraph.Graph, error) {
	graph := langraph.NewGraph()

	m.agentsMutex.RLock()
	for id, agent := range m.agents {
		node := langraph.NewNode(id, agent)
		graph.AddNode(node)
	}
	m.agentsMutex.RUnlock()

	if m.config.EnableOrchestrator {
		orchestratorID := ""
		m.agentsMutex.RLock()
		for id, agent := range m.agents {
			if agent.Type() == string(types.AgentTypeOrchestrator) {
				orchestratorID = id
				break
			}
		}
		m.agentsMutex.RUnlock()

		if orchestratorID != "" {
			m.agentsMutex.RLock()
			for id, agent := range m.agents {
				if id != orchestratorID {
					graph.AddEdge(id, orchestratorID, map[string]interface{}{
						"weight": 1.0,
						"type":   "request",
					})

					graph.AddEdge(orchestratorID, id, map[string]interface{}{
						"weight": 1.0,
						"type":   "response",
					})
				}
			}
			m.agentsMutex.RUnlock()
		}
	} else {
		m.agentsMutex.RLock()
		for id1, agent1 := range m.agents {
			for id2, agent2 := range m.agents {
				if id1 != id2 {
					if agent1.Type() == string(types.AgentTypeFrontend) && agent2.Type() == string(types.AgentTypeAppBuilder) {
						graph.AddEdge(id1, id2, map[string]interface{}{
							"weight": 1.0,
							"type":   "ui_to_app",
						})
					}

					if agent1.Type() == string(types.AgentTypeAppBuilder) && agent2.Type() == string(types.AgentTypeCodegen) {
						graph.AddEdge(id1, id2, map[string]interface{}{
							"weight": 1.0,
							"type":   "app_to_code",
						})
					}

					if agent1.Type() == string(types.AgentTypeCodegen) && agent2.Type() == string(types.AgentTypeEngineering) {
						graph.AddEdge(id1, id2, map[string]interface{}{
							"weight": 1.0,
							"type":   "code_to_engineering",
						})
					}

					if agent1.Type() == string(types.AgentTypeEngineering) && agent2.Type() == string(types.AgentTypeFrontend) {
						graph.AddEdge(id1, id2, map[string]interface{}{
							"weight": 1.0,
							"type":   "engineering_to_ui",
						})
					}
				}
			}
		}
		m.agentsMutex.RUnlock()
	}

	entryPointID := ""
	switch task.Type {
	case "ui_design":
		m.agentsMutex.RLock()
		for id, agent := range m.agents {
			if agent.Type() == string(types.AgentTypeFrontend) {
				entryPointID = id
				break
			}
		}
		m.agentsMutex.RUnlock()
	case "code_generation":
		m.agentsMutex.RLock()
		for id, agent := range m.agents {
			if agent.Type() == string(types.AgentTypeCodegen) {
				entryPointID = id
				break
			}
		}
		m.agentsMutex.RUnlock()
	case "app_building":
		m.agentsMutex.RLock()
		for id, agent := range m.agents {
			if agent.Type() == string(types.AgentTypeAppBuilder) {
				entryPointID = id
				break
			}
		}
		m.agentsMutex.RUnlock()
	case "engineering":
		m.agentsMutex.RLock()
		for id, agent := range m.agents {
			if agent.Type() == string(types.AgentTypeEngineering) {
				entryPointID = id
				break
			}
		}
		m.agentsMutex.RUnlock()
	default:
		if m.config.EnableOrchestrator {
			m.agentsMutex.RLock()
			for id, agent := range m.agents {
				if agent.Type() == string(types.AgentTypeOrchestrator) {
					entryPointID = id
					break
				}
			}
			m.agentsMutex.RUnlock()
		}

		if entryPointID == "" {
			m.agentsMutex.RLock()
			for id := range m.agents {
				entryPointID = id
				break
			}
			m.agentsMutex.RUnlock()
		}
	}

	if entryPointID == "" {
		return nil, fmt.Errorf("failed to determine entry point for task")
	}

	graph.SetEntryPoint(entryPointID)

	graph.SetParameters(map[string]interface{}{
		"task_id":          task.ID,
		"task_type":        task.Type,
		"task_description": task.Description,
		"task_parameters":  task.Parameters,
	})

	return graph, nil
}

func (m *MultiAgentManager) emitEvent(event *types.AgentEvent) {
	m.listenerMutex.RLock()
	defer m.listenerMutex.RUnlock()

	for _, listener := range m.eventListeners {
		go listener(event)
	}
}
