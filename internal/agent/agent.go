package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Agent struct {
	config    *config.Config
	mcpManager *mcp.Manager
	tools     *tools.Registry
	state     *State
}

type State struct {
	Status    string    `json:"status"`
	StartTime time.Time `json:"startTime"`
	Task      string    `json:"task"`
}

type ExecutionResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func New(cfg *config.Config, mcpManager *mcp.Manager) (*Agent, error) {
	toolRegistry, err := tools.NewRegistry(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tool registry: %w", err)
	}
	
	agent := &Agent{
		config:    cfg,
		mcpManager: mcpManager,
		tools:     toolRegistry,
		state: &State{
			Status:    "idle",
			StartTime: time.Now(),
		},
	}
	
	return agent, nil
}

func (a *Agent) Execute(ctx context.Context, task string) (*ExecutionResult, error) {
	a.state.Status = "running"
	a.state.Task = task
	
	
	result := &ExecutionResult{
		Success: true,
		Message: "Task executed successfully",
		Data: map[string]interface{}{
			"task": task,
			"time": time.Now().Format(time.RFC3339),
		},
	}
	
	a.state.Status = "idle"
	a.state.Task = ""
	
	return result, nil
}

func (a *Agent) Status() *State {
	return a.state
}
