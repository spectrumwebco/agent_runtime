package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type AgentRole string

const (
	AgentRoleFrontend AgentRole = "frontend"
	AgentRoleAppBuilder AgentRole = "app_builder"
	AgentRoleCodegen AgentRole = "codegen"
	AgentRoleEngineering AgentRole = "engineering"
	AgentRoleOrchestrator AgentRole = "orchestrator"
)

type AgentCapability string

const (
	AgentCapabilityCodeGeneration AgentCapability = "code_generation"
	AgentCapabilityCodeReview AgentCapability = "code_review"
	AgentCapabilityUIDesign AgentCapability = "ui_design"
	AgentCapabilityAPIDesign AgentCapability = "api_design"
	AgentCapabilityDatabaseDesign AgentCapability = "database_design"
	AgentCapabilityTesting AgentCapability = "testing"
	AgentCapabilityDeployment AgentCapability = "deployment"
	AgentCapabilityDocumentation AgentCapability = "documentation"
	AgentCapabilityPlanning AgentCapability = "planning"
	AgentCapabilityOrchestration AgentCapability = "orchestration"
)

type AgentConfig struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Role         AgentRole              `json:"role"`
	Capabilities []AgentCapability      `json:"capabilities"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ModelConfig  map[string]interface{} `json:"model_config,omitempty"`
	ToolConfig   map[string]interface{} `json:"tool_config,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type Agent struct {
	Config      AgentConfig             `json:"config"`
	State       map[string]interface{}  `json:"state,omitempty"`
	Tools       []Tool                  `json:"-"`
	Handler     AgentHandler            `json:"-"`
	EventStream EventStream             `json:"-"`
	stateLock   sync.RWMutex            `json:"-"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Handler     ToolHandler            `json:"-"`
}

type ToolHandler func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error)

type AgentHandler func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error)

type EventStream interface {
	AddEvent(event *models.Event) error
	Subscribe(eventType models.EventType, callback func(*models.Event)) error
	Unsubscribe(eventType models.EventType, callback func(*models.Event)) error
}

func NewAgent(config AgentConfig, handler AgentHandler, eventStream EventStream) *Agent {
	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now().UTC()
	}
	
	config.UpdatedAt = time.Now().UTC()

	return &Agent{
		Config:      config,
		State:       make(map[string]interface{}),
		Tools:       []Tool{},
		Handler:     handler,
		EventStream: eventStream,
	}
}

func (a *Agent) AddTool(name, description string, handler ToolHandler, parameters map[string]interface{}) {
	a.Tools = append(a.Tools, Tool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
		Handler:     handler,
	})
}

func (a *Agent) GetTool(name string) (*Tool, error) {
	for _, tool := range a.Tools {
		if tool.Name == name {
			return &tool, nil
		}
	}
	return nil, fmt.Errorf("tool %s not found", name)
}

func (a *Agent) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (map[string]interface{}, error) {
	tool, err := a.GetTool(toolName)
	if err != nil {
		return nil, err
	}

	if tool.Handler == nil {
		return nil, fmt.Errorf("tool %s has no handler", toolName)
	}

	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool":       toolName,
			"parameters": params,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	result, err := tool.Handler(ctx, a, params)
	
	resultEvent := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool":       toolName,
			"result":     result,
			"error":      err != nil,
			"error_msg":  err,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(resultEvent); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	return result, err
}

func (a *Agent) Process(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	if a.Handler == nil {
		return nil, fmt.Errorf("agent %s has no handler", a.Config.Name)
	}

	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"input":      input,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	result, err := a.Handler(ctx, a, input)
	
	resultEvent := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"result":     result,
			"error":      err != nil,
			"error_msg":  err,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(resultEvent); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	return result, err
}

func (a *Agent) SetState(state map[string]interface{}) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	a.State = state
}

func (a *Agent) UpdateState(updates map[string]interface{}) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	
	if a.State == nil {
		a.State = make(map[string]interface{})
	}
	
	for k, v := range updates {
		a.State[k] = v
	}
}

func (a *Agent) GetState() map[string]interface{} {
	a.stateLock.RLock()
	defer a.stateLock.RUnlock()
	
	stateCopy := make(map[string]interface{})
	for k, v := range a.State {
		stateCopy[k] = v
	}
	
	return stateCopy
}

func (a *Agent) HasCapability(capability AgentCapability) bool {
	for _, cap := range a.Config.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}
