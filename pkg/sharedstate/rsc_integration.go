package sharedstate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/ai/vercel"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/adapters"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type RSCIntegration struct {
	sharedStateManager *SharedStateManager
	rscAdapter         *adapters.RSCAdapter
	eventStream        *eventstream.Stream
	rscClient          *vercel.RSCAdapter
	rscIntegration     *vercel.RSCIntegration

	componentCallbacks map[string]func(component vercel.GeneratedComponent)
	mu                 sync.RWMutex
}

func NewRSCIntegration(
	sharedStateManager *SharedStateManager,
	rscAdapter *adapters.RSCAdapter,
	eventStream *eventstream.Stream,
	rscClient *vercel.RSCAdapter,
	rscIntegration *vercel.RSCIntegration,
) *RSCIntegration {
	integration := &RSCIntegration{
		sharedStateManager:  sharedStateManager,
		rscAdapter:          rscAdapter,
		eventStream:         eventStream,
		rscClient:           rscClient,
		rscIntegration:      rscIntegration,
		componentCallbacks:  make(map[string]func(component vercel.GeneratedComponent)),
	}

	integration.setupEventListeners()

	return integration
}

func (i *RSCIntegration) setupEventListeners() {
	if i.eventStream == nil {
		return
	}

	i.eventStream.Subscribe("state_update", func(event eventstream.Event) {
		var stateEvent struct {
			StateType string                 `json:"state_type"`
			StateID   string                 `json:"state_id"`
			Data      map[string]interface{} `json:"data"`
		}

		err := json.Unmarshal([]byte(event.Payload), &stateEvent)
		if err != nil {
			log.Printf("Error unmarshaling state update event: %v", err)
			return
		}

		if stateEvent.StateType == string(models.StateTypeComponent) {
			i.handleComponentStateUpdate(stateEvent.StateID, stateEvent.Data)
		}
	})

	i.eventStream.Subscribe("agent_action", func(event eventstream.Event) {
		var actionEvent struct {
			AgentID    string                 `json:"agent_id"`
			ActionID   string                 `json:"action_id"`
			ActionType string                 `json:"action_type"`
			ActionData map[string]interface{} `json:"action_data"`
		}

		err := json.Unmarshal([]byte(event.Payload), &actionEvent)
		if err != nil {
			log.Printf("Error unmarshaling agent action event: %v", err)
			return
		}

		go func() {
			componentID, err := i.rscAdapter.GenerateComponentFromAgentAction(
				context.Background(),
				actionEvent.AgentID,
				actionEvent.ActionID,
				actionEvent.ActionType,
				actionEvent.ActionData,
			)
			if err != nil {
				log.Printf("Error generating component from agent action: %v", err)
				return
			}

			i.notifyComponentGenerated(componentID)
		}()
	})

	i.eventStream.Subscribe("tool_usage", func(event eventstream.Event) {
		var toolEvent struct {
			AgentID    string                 `json:"agent_id"`
			ToolID     string                 `json:"tool_id"`
			ToolName   string                 `json:"tool_name"`
			ToolInput  map[string]interface{} `json:"tool_input"`
			ToolOutput map[string]interface{} `json:"tool_output"`
		}

		err := json.Unmarshal([]byte(event.Payload), &toolEvent)
		if err != nil {
			log.Printf("Error unmarshaling tool usage event: %v", err)
			return
		}

		go func() {
			componentID, err := i.rscAdapter.GenerateComponentFromToolUsage(
				context.Background(),
				toolEvent.AgentID,
				toolEvent.ToolID,
				toolEvent.ToolName,
				toolEvent.ToolInput,
				toolEvent.ToolOutput,
			)
			if err != nil {
				log.Printf("Error generating component from tool usage: %v", err)
				return
			}

			i.notifyComponentGenerated(componentID)
		}()
	})
}

func (i *RSCIntegration) handleComponentStateUpdate(stateID string, data map[string]interface{}) {
	componentID := stateID
	if len(componentID) > 10 && componentID[:10] == "component:" {
		componentID = componentID[10:]
	}

	component, err := i.rscAdapter.GetComponent(context.Background(), componentID)
	if err != nil {
		log.Printf("Error getting component: %v", err)
		return
	}

	i.notifyComponentUpdated(*component)
}

func (i *RSCIntegration) notifyComponentGenerated(componentID string) {
	component, err := i.rscAdapter.GetComponent(context.Background(), componentID)
	if err != nil {
		log.Printf("Error getting component: %v", err)
		return
	}

	i.notifyComponentUpdated(*component)
}

func (i *RSCIntegration) notifyComponentUpdated(component vercel.GeneratedComponent) {
	i.mu.RLock()
	for _, callback := range i.componentCallbacks {
		go callback(component)
	}
	i.mu.RUnlock()

	if i.eventStream != nil {
		componentJSON, err := json.Marshal(component)
		if err != nil {
			log.Printf("Error marshaling component: %v", err)
			return
		}

		event := eventstream.Event{
			Type:    "component_updated",
			Payload: string(componentJSON),
		}

		err = i.eventStream.Publish(event)
		if err != nil {
			log.Printf("Error publishing component updated event: %v", err)
		}
	}
}

func (i *RSCIntegration) RegisterComponentCallback(id string, callback func(component vercel.GeneratedComponent)) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.componentCallbacks[id] = callback
}

func (i *RSCIntegration) UnregisterComponentCallback(id string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.componentCallbacks, id)
}

func (i *RSCIntegration) GenerateComponent(ctx context.Context, componentType vercel.ComponentType, props map[string]interface{}) (string, error) {
	return i.rscAdapter.GenerateComponent(ctx, componentType, props)
}

func (i *RSCIntegration) GetComponent(ctx context.Context, componentID string) (*vercel.GeneratedComponent, error) {
	return i.rscAdapter.GetComponent(ctx, componentID)
}

func (i *RSCIntegration) ListComponents(ctx context.Context) ([]vercel.GeneratedComponent, error) {
	return i.rscAdapter.ListComponents(ctx)
}

func (i *RSCIntegration) GetComponentsByAgent(ctx context.Context, agentID string) ([]vercel.GeneratedComponent, error) {
	if i.rscIntegration != nil {
		return i.rscIntegration.GetComponentsByAgent(agentID)
	}

	states, err := i.sharedStateManager.ListStates(models.StateTypeAction)
	if err != nil {
		return nil, fmt.Errorf("error listing action states: %w", err)
	}

	componentIDs := make([]string, 0)
	for _, state := range states {
		if state.Data == nil {
			continue
		}

		stateAgentID, ok := state.Data["agent_id"].(string)
		if !ok || stateAgentID != agentID {
			continue
		}

		componentID, ok := state.Data["component_id"].(string)
		if !ok {
			continue
		}

		componentIDs = append(componentIDs, componentID)
	}

	components := make([]vercel.GeneratedComponent, 0, len(componentIDs))
	for _, componentID := range componentIDs {
		component, err := i.rscAdapter.GetComponent(ctx, componentID)
		if err != nil {
			continue
		}

		components = append(components, *component)
	}

	return components, nil
}

func (i *RSCIntegration) GetComponentsByTool(ctx context.Context, toolID string) ([]vercel.GeneratedComponent, error) {
	if i.rscIntegration != nil {
		return i.rscIntegration.GetComponentsByTool(toolID)
	}

	states, err := i.sharedStateManager.ListStates(models.StateTypeTool)
	if err != nil {
		return nil, fmt.Errorf("error listing tool states: %w", err)
	}

	componentIDs := make([]string, 0)
	for _, state := range states {
		if state.Data == nil {
			continue
		}

		stateToolID, ok := state.Data["tool_id"].(string)
		if !ok || stateToolID != toolID {
			continue
		}

		componentID, ok := state.Data["component_id"].(string)
		if !ok {
			continue
		}

		componentIDs = append(componentIDs, componentID)
	}

	components := make([]vercel.GeneratedComponent, 0, len(componentIDs))
	for _, componentID := range componentIDs {
		component, err := i.rscAdapter.GetComponent(ctx, componentID)
		if err != nil {
			continue
		}

		components = append(components, *component)
	}

	return components, nil
}
