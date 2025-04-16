package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type AgentFactory struct {
	eventStream  EventStream
	stateManager StateManager
	agents       map[string]*Agent
	mutex        sync.RWMutex
}

func NewAgentFactory(eventStream EventStream, stateManager StateManager) *AgentFactory {
	return &AgentFactory{
		eventStream:  eventStream,
		stateManager: stateManager,
		agents:       make(map[string]*Agent),
		mutex:        sync.RWMutex{},
	}
}

func (f *AgentFactory) CreateAgent(id, name, role string, capabilities []string) (*Agent, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, exists := f.agents[id]; exists {
		return nil, fmt.Errorf("agent with ID %s already exists", id)
	}

	agent := NewAgent(id, name, role, capabilities, f.eventStream, f.stateManager)

	f.agents[id] = agent

	return agent, nil
}

func (f *AgentFactory) GetAgent(id string) (*Agent, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}

	return agent, nil
}

func (f *AgentFactory) ListAgents() []*Agent {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	agents := make([]*Agent, 0, len(f.agents))
	for _, agent := range f.agents {
		agents = append(agents, agent)
	}

	return agents
}

func (f *AgentFactory) DeleteAgent(id string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	agent, exists := f.agents[id]
	if !exists {
		return fmt.Errorf("agent with ID %s not found", id)
	}

	if agent.isRunning {
		if err := agent.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent: %v", err)
		}
	}

	delete(f.agents, id)

	return nil
}

func (f *AgentFactory) StartAgent(ctx context.Context, id string) error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return fmt.Errorf("agent with ID %s not found", id)
	}

	return agent.Start(ctx)
}

func (f *AgentFactory) StopAgent(id string) error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return fmt.Errorf("agent with ID %s not found", id)
	}

	return agent.Stop()
}

func (f *AgentFactory) StartAllAgents(ctx context.Context) error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for id, agent := range f.agents {
		if err := agent.Start(ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %v", id, err)
		}
	}

	return nil
}

func (f *AgentFactory) StopAllAgents() error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for id, agent := range f.agents {
		if err := agent.Stop(); err != nil {
			return fmt.Errorf("failed to stop agent %s: %v", id, err)
		}
	}

	return nil
}

func (f *AgentFactory) GetAgentsByRole(role string) []*Agent {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	var agents []*Agent
	for _, agent := range f.agents {
		if agent.Role == role {
			agents = append(agents, agent)
		}
	}

	return agents
}

func (f *AgentFactory) GetAgentsByCapability(capability string) []*Agent {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	var agents []*Agent
	for _, agent := range f.agents {
		if agent.HasCapability(capability) {
			agents = append(agents, agent)
		}
	}

	return agents
}

func (f *AgentFactory) UpdateAgentState(id string, state *statemodels.AgentState) error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return fmt.Errorf("agent with ID %s not found", id)
	}

	return agent.UpdateState(state)
}

func (f *AgentFactory) GetAgentState(id string) (*statemodels.AgentState, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}

	return agent.GetState(), nil
}
