package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/ai/vercel"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type RSCAdapter struct {
	rscClient      *vercel.RSCAdapter
	stateMap       map[string]*models.State
	componentMap   map[string]vercel.GeneratedComponent
	mu             sync.RWMutex
	eventCallbacks map[string]func(event interface{})
}

func NewRSCAdapter(rscClient *vercel.RSCAdapter) *RSCAdapter {
	return &RSCAdapter{
		rscClient:      rscClient,
		stateMap:       make(map[string]*models.State),
		componentMap:   make(map[string]vercel.GeneratedComponent),
		eventCallbacks: make(map[string]func(event interface{})),
	}
}

func (a *RSCAdapter) GetState(ctx context.Context, stateType models.StateType, stateID string) (*models.State, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", stateType, stateID)
	state, exists := a.stateMap[key]
	if !exists {
		return nil, fmt.Errorf("state not found: %s", key)
	}

	return state.Clone(), nil
}

func (a *RSCAdapter) SetState(ctx context.Context, stateType models.StateType, stateID string, data map[string]interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := fmt.Sprintf("%s:%s", stateType, stateID)
	state, exists := a.stateMap[key]
	if !exists {
		state = models.NewState(stateType, data, nil)
		a.stateMap[key] = state
	} else {
		state.Update(data)
	}

	if stateType == models.StateTypeComponent && data != nil {
		componentJSON, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshaling component data: %w", err)
		}

		var component vercel.GeneratedComponent
		err = json.Unmarshal(componentJSON, &component)
		if err != nil {
			return fmt.Errorf("error unmarshaling component data: %w", err)
		}

		a.componentMap[stateID] = component
	}

	return nil
}

func (a *RSCAdapter) UpdateState(ctx context.Context, stateType models.StateType, stateID string, data map[string]interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := fmt.Sprintf("%s:%s", stateType, stateID)
	state, exists := a.stateMap[key]
	if !exists {
		state = models.NewState(stateType, data, nil)
		a.stateMap[key] = state
	} else {
		state.Merge(data)
	}

	if stateType == models.StateTypeComponent && data != nil {
		componentID := stateID
		component, exists := a.componentMap[componentID]
		if !exists {
			componentJSON, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("error marshaling component data: %w", err)
			}

			err = json.Unmarshal(componentJSON, &component)
			if err != nil {
				return fmt.Errorf("error unmarshaling component data: %w", err)
			}
		} else {
			for k, v := range data {
				switch k {
				case "props":
					if props, ok := v.(map[string]interface{}); ok {
						for pk, pv := range props {
							component.Props[pk] = pv
						}
					}
				case "children":
					if children, ok := v.([]interface{}); ok {
						childComponents := make([]vercel.GeneratedComponent, 0, len(children))
						for _, child := range children {
							if childData, ok := child.(map[string]interface{}); ok {
								childJSON, err := json.Marshal(childData)
								if err != nil {
									continue
								}

								var childComponent vercel.GeneratedComponent
								err = json.Unmarshal(childJSON, &childComponent)
								if err != nil {
									continue
								}

								childComponents = append(childComponents, childComponent)
							}
						}
						component.Children = childComponents
					}
				case "updated_at":
					if timestamp, ok := v.(int64); ok {
						component.UpdatedAt = timestamp
					}
				}
			}
		}

		component.UpdatedAt = time.Now().Unix()
		a.componentMap[componentID] = component
	}

	return nil
}

func (a *RSCAdapter) DeleteState(ctx context.Context, stateType models.StateType, stateID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := fmt.Sprintf("%s:%s", stateType, stateID)
	delete(a.stateMap, key)

	if stateType == models.StateTypeComponent {
		delete(a.componentMap, stateID)
	}

	return nil
}

func (a *RSCAdapter) ListStates(ctx context.Context, stateType models.StateType) ([]*models.State, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	prefix := fmt.Sprintf("%s:", stateType)
	states := make([]*models.State, 0)

	for key, state := range a.stateMap {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			states = append(states, state.Clone())
		}
	}

	return states, nil
}

func (a *RSCAdapter) GenerateComponent(ctx context.Context, componentType vercel.ComponentType, props map[string]interface{}) (string, error) {
	if a.rscClient == nil {
		return "", fmt.Errorf("RSC client not initialized")
	}

	component, err := a.rscClient.GenerateComponent(ctx, componentType, props)
	if err != nil {
		return "", fmt.Errorf("error generating component: %w", err)
	}

	if component.ID == "" {
		component.ID = uuid.New().String()
	}

	componentID := component.ID
	stateID := fmt.Sprintf("component:%s", componentID)

	componentJSON, err := json.Marshal(component)
	if err != nil {
		return "", fmt.Errorf("error marshaling component: %w", err)
	}

	var componentData map[string]interface{}
	err = json.Unmarshal(componentJSON, &componentData)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling component: %w", err)
	}

	err = a.SetState(ctx, models.StateTypeComponent, stateID, componentData)
	if err != nil {
		return "", fmt.Errorf("error setting component state: %w", err)
	}

	return componentID, nil
}

func (a *RSCAdapter) GenerateComponentFromAgentAction(ctx context.Context, agentID string, actionID string, actionType string, actionData map[string]interface{}) (string, error) {
	if a.rscClient == nil {
		return "", fmt.Errorf("RSC client not initialized")
	}

	component, err := a.rscClient.GenerateComponentFromAgentAction(ctx, agentID, actionID, actionType, actionData)
	if err != nil {
		return "", fmt.Errorf("error generating component from agent action: %w", err)
	}

	if component.ID == "" {
		component.ID = uuid.New().String()
	}

	componentID := component.ID
	stateID := fmt.Sprintf("component:%s", componentID)

	componentJSON, err := json.Marshal(component)
	if err != nil {
		return "", fmt.Errorf("error marshaling component: %w", err)
	}

	var componentData map[string]interface{}
	err = json.Unmarshal(componentJSON, &componentData)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling component: %w", err)
	}

	err = a.SetState(ctx, models.StateTypeComponent, stateID, componentData)
	if err != nil {
		return "", fmt.Errorf("error setting component state: %w", err)
	}

	actionStateID := fmt.Sprintf("action:%s", actionID)
	actionState := map[string]interface{}{
		"agent_id":     agentID,
		"action_id":    actionID,
		"action_type":  actionType,
		"action_data":  actionData,
		"component_id": componentID,
		"timestamp":    time.Now().Unix(),
	}

	err = a.SetState(ctx, models.StateTypeAction, actionStateID, actionState)
	if err != nil {
		log.Printf("Error setting action state: %v", err)
	}

	return componentID, nil
}

func (a *RSCAdapter) GenerateComponentFromToolUsage(ctx context.Context, agentID string, toolID string, toolName string, toolInput map[string]interface{}, toolOutput map[string]interface{}) (string, error) {
	if a.rscClient == nil {
		return "", fmt.Errorf("RSC client not initialized")
	}

	component, err := a.rscClient.GenerateComponentFromToolUsage(ctx, agentID, toolID, toolName, toolInput, toolOutput)
	if err != nil {
		return "", fmt.Errorf("error generating component from tool usage: %w", err)
	}

	if component.ID == "" {
		component.ID = uuid.New().String()
	}

	componentID := component.ID
	stateID := fmt.Sprintf("component:%s", componentID)

	componentJSON, err := json.Marshal(component)
	if err != nil {
		return "", fmt.Errorf("error marshaling component: %w", err)
	}

	var componentData map[string]interface{}
	err = json.Unmarshal(componentJSON, &componentData)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling component: %w", err)
	}

	err = a.SetState(ctx, models.StateTypeComponent, stateID, componentData)
	if err != nil {
		return "", fmt.Errorf("error setting component state: %w", err)
	}

	toolStateID := fmt.Sprintf("tool:%s", toolID)
	toolState := map[string]interface{}{
		"agent_id":     agentID,
		"tool_id":      toolID,
		"tool_name":    toolName,
		"tool_input":   toolInput,
		"tool_output":  toolOutput,
		"component_id": componentID,
		"timestamp":    time.Now().Unix(),
	}

	err = a.SetState(ctx, models.StateTypeTool, toolStateID, toolState)
	if err != nil {
		log.Printf("Error setting tool state: %v", err)
	}

	return componentID, nil
}

func (a *RSCAdapter) GetComponent(ctx context.Context, componentID string) (*vercel.GeneratedComponent, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	component, exists := a.componentMap[componentID]
	if !exists {
		return nil, fmt.Errorf("component not found: %s", componentID)
	}

	componentJSON, err := json.Marshal(component)
	if err != nil {
		return nil, fmt.Errorf("error marshaling component: %w", err)
	}

	var componentCopy vercel.GeneratedComponent
	err = json.Unmarshal(componentJSON, &componentCopy)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling component: %w", err)
	}

	return &componentCopy, nil
}

func (a *RSCAdapter) ListComponents(ctx context.Context) ([]vercel.GeneratedComponent, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	components := make([]vercel.GeneratedComponent, 0, len(a.componentMap))
	for _, component := range a.componentMap {
		componentJSON, err := json.Marshal(component)
		if err != nil {
			continue
		}

		var componentCopy vercel.GeneratedComponent
		err = json.Unmarshal(componentJSON, &componentCopy)
		if err != nil {
			continue
		}

		components = append(components, componentCopy)
	}

	return components, nil
}

func (a *RSCAdapter) RegisterEventCallback(id string, callback func(event interface{})) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.eventCallbacks[id] = callback
}

func (a *RSCAdapter) UnregisterEventCallback(id string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.eventCallbacks, id)
}

func (a *RSCAdapter) HandleEvent(event interface{}) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, callback := range a.eventCallbacks {
		go callback(event)
	}
}
