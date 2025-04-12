package vercel

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type ActionHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

type ServerActionsAdapter struct {
	actions      map[string]ActionHandler
	mu           sync.RWMutex
	stateManager *sharedstate.SharedStateManager
	aiClient     *VercelAIClient
}

func NewServerActionsAdapter(stateManager *sharedstate.SharedStateManager, aiClient *VercelAIClient) *ServerActionsAdapter {
	return &ServerActionsAdapter{
		actions:      make(map[string]ActionHandler),
		stateManager: stateManager,
		aiClient:     aiClient,
	}
}

func (a *ServerActionsAdapter) RegisterAction(name string, handler ActionHandler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.actions[name] = handler
}

func (a *ServerActionsAdapter) ExecuteAction(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	a.mu.RLock()
	handler, exists := a.actions[name]
	a.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("action %s not found", name)
	}

	stateCtx := context.WithValue(ctx, "state_manager", a.stateManager)

	result, err := handler(stateCtx, params)
	if err != nil {
		return nil, fmt.Errorf("error executing action %s: %w", name, err)
	}

	if stateID, ok := params["state_id"].(string); ok && stateID != "" {
		stateData := map[string]interface{}{
			"action":     name,
			"parameters": params,
			"result":     result,
		}
		
		err = a.stateManager.UpdateState(sharedstate.ActionState, stateID, stateData)
		if err != nil {
			return nil, fmt.Errorf("error updating state: %w", err)
		}
		
		a.aiClient.HandleStateUpdate(sharedstate.ActionState, stateID, stateData)
	}

	return result, nil
}

func (a *ServerActionsAdapter) ListActions() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	actions := make([]string, 0, len(a.actions))
	for name := range a.actions {
		actions = append(actions, name)
	}
	
	return actions
}
