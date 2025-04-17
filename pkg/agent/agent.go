package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type Agent struct {
	ID            string
	Name          string
	Role          string
	Capabilities  []string
	Configuration map[string]interface{}
	State         *statemodels.AgentState
	eventStream   EventStream
	stateManager  StateManager
	tools         []Tool
	neovimTools   []NeovimTool
	activeNeovim  map[string]bool
	mutex         sync.RWMutex
	isRunning     bool
	stopCh        chan struct{}
}

type EventStream interface {
	AddEvent(event *models.Event) error
}

type StateManager interface {
	GetState(agentID string) (*statemodels.AgentState, error)
	UpdateState(agentID string, state *statemodels.AgentState) error
}

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error)
}

type NeovimTool interface {
	Tool
	IsNeovimTool() bool
}

func NewAgent(id, name, role string, capabilities []string, eventStream EventStream, stateManager StateManager) *Agent {
	if id == "" {
		id = uuid.New().String()
	}

	hasNeovimCapability := false
	for _, cap := range capabilities {
		if cap == "neovim" {
			hasNeovimCapability = true
			break
		}
	}
	if !hasNeovimCapability {
		capabilities = append(capabilities, "neovim")
	}

	return &Agent{
		ID:            id,
		Name:          name,
		Role:          role,
		Capabilities:  capabilities,
		Configuration: make(map[string]interface{}),
		State: &statemodels.AgentState{
			ID:        id,
			Status:    "initialized",
			Timestamp: time.Now(),
			Data:      make(map[string]interface{}),
		},
		eventStream:  eventStream,
		stateManager: stateManager,
		tools:        []Tool{},
		neovimTools:  []NeovimTool{},
		activeNeovim: make(map[string]bool),
		mutex:        sync.RWMutex{},
		isRunning:    false,
		stopCh:       make(chan struct{}),
	}
}

func (a *Agent) Start(ctx context.Context) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.isRunning {
		return fmt.Errorf("agent %s is already running", a.ID)
	}

	a.isRunning = true
	a.State.Status = "running"
	a.State.Timestamp = time.Now()

	if err := a.stateManager.UpdateState(a.ID, a.State); err != nil {
		log.Printf("Failed to update agent state: %v", err)
	}

	a.emitEvent(models.EventTypeSystem, models.EventSourceAgent, "agent_started", map[string]string{
		"agent_id":   a.ID,
		"agent_name": a.Name,
		"agent_role": a.Role,
	})

	go a.run(ctx)

	return nil
}

func (a *Agent) Stop() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if !a.isRunning {
		return fmt.Errorf("agent %s is not running", a.ID)
	}

	close(a.stopCh)
	a.isRunning = false
	a.State.Status = "stopped"
	a.State.Timestamp = time.Now()

	if err := a.stateManager.UpdateState(a.ID, a.State); err != nil {
		log.Printf("Failed to update agent state: %v", err)
	}

	a.emitEvent(models.EventTypeSystem, models.EventSourceAgent, "agent_stopped", map[string]string{
		"agent_id":   a.ID,
		"agent_name": a.Name,
		"agent_role": a.Role,
	})

	return nil
}

func (a *Agent) AddTool(tool Tool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if neovimTool, ok := tool.(NeovimTool); ok && neovimTool.IsNeovimTool() {
		a.neovimTools = append(a.neovimTools, neovimTool)
		a.emitEvent(models.EventTypeNeovim, models.EventSourceNeovim, "neovim_tool_added", map[string]string{
			"agent_id":  a.ID,
			"tool_name": tool.Name(),
		})
	} else {
		a.tools = append(a.tools, tool)
	}
}

func (a *Agent) GetTools() []Tool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.tools
}

func (a *Agent) GetNeovimTools() []NeovimTool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.neovimTools
}

func (a *Agent) ExecuteTool(ctx context.Context, toolName string, args map[string]interface{}) (map[string]interface{}, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, tool := range a.tools {
		if tool.Name() == toolName {
			a.emitEvent(models.EventTypeAction, models.EventSourceAgent, "tool_execution_started", map[string]string{
				"agent_id":  a.ID,
				"tool_name": toolName,
			})

			result, err := tool.Execute(ctx, args)

			if err != nil {
				a.emitEvent(models.EventTypeAction, models.EventSourceAgent, "tool_execution_failed", map[string]string{
					"agent_id":  a.ID,
					"tool_name": toolName,
					"error":     err.Error(),
				})
				return nil, err
			}

			a.emitEvent(models.EventTypeAction, models.EventSourceAgent, "tool_execution_completed", map[string]string{
				"agent_id":  a.ID,
				"tool_name": toolName,
			})

			return result, nil
		}
	}

	for _, tool := range a.neovimTools {
		if tool.Name() == toolName {
			a.emitEvent(models.EventTypeNeovim, models.EventSourceNeovim, "neovim_tool_execution_started", map[string]string{
				"agent_id":  a.ID,
				"tool_name": toolName,
			})

			result, err := tool.Execute(ctx, args)

			if err != nil {
				a.emitEvent(models.EventTypeNeovim, models.EventSourceNeovim, "neovim_tool_execution_failed", map[string]string{
					"agent_id":  a.ID,
					"tool_name": toolName,
					"error":     err.Error(),
				})
				return nil, err
			}

			a.emitEvent(models.EventTypeNeovim, models.EventSourceNeovim, "neovim_tool_execution_completed", map[string]string{
				"agent_id":  a.ID,
				"tool_name": toolName,
			})

			return result, nil
		}
	}

	return nil, fmt.Errorf("tool %s not found", toolName)
}

func (a *Agent) UpdateConfiguration(config map[string]interface{}) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for k, v := range config {
		a.Configuration[k] = v
	}

	a.emitEvent(models.EventTypeSystem, models.EventSourceAgent, "configuration_updated", map[string]string{
		"agent_id": a.ID,
	})
}

func (a *Agent) GetConfiguration() map[string]interface{} {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.Configuration
}

func (a *Agent) UpdateState(state *statemodels.AgentState) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.State = state
	a.State.Timestamp = time.Now()

	if err := a.stateManager.UpdateState(a.ID, a.State); err != nil {
		return err
	}

	a.emitEvent(models.EventTypeSystem, models.EventSourceAgent, "state_updated", map[string]string{
		"agent_id": a.ID,
		"status":   a.State.Status,
	})

	return nil
}

func (a *Agent) GetState() *statemodels.AgentState {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.State
}

func (a *Agent) HasCapability(capability string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, cap := range a.Capabilities {
		if cap == capability {
			return true
		}
	}

	return false
}

func (a *Agent) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.Stop()
			return
		case <-a.stopCh:
			return
		case <-ticker.C:
			a.checkState()
		}
	}
}

func (a *Agent) checkState() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.State.Timestamp = time.Now()

	if err := a.stateManager.UpdateState(a.ID, a.State); err != nil {
		log.Printf("Failed to update agent state: %v", err)
	}
}

func (a *Agent) emitEvent(eventType models.EventType, eventSource models.EventSource, action string, metadata map[string]string) {
	if a.eventStream == nil {
		return
	}

	event := &models.Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    eventSource,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_id": a.ID,
			"action":   action,
		},
		Metadata: metadata,
	}

	if err := a.eventStream.AddEvent(event); err != nil {
		log.Printf("Failed to emit event: %v", err)
	}
}

func (a *Agent) ToLangGraphAgent() (*langraph.Agent, error) {
	config := &langraph.AgentConfig{
		ID:           a.ID,
		Name:         a.Name,
		Role:         a.Role,
		Capabilities: a.Capabilities,
		Description:  fmt.Sprintf("%s agent with role %s", a.Name, a.Role),
	}

	return langraph.NewAgent(config), nil
}
