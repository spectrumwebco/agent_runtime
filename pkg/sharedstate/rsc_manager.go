package sharedstate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/ai/vercel"
)

type RSCManager struct {
	sharedStateManager *SharedStateManager
	eventStream        *eventstream.Stream
	rscClient          *vercel.RSCAdapter
	rscIntegration     *vercel.RSCIntegration
	componentCache     map[string]vercel.GeneratedComponent
	mu                 sync.RWMutex
}

func NewRSCManager(sharedStateManager *SharedStateManager, eventStream *eventstream.Stream, rscClient *vercel.RSCAdapter, rscIntegration *vercel.RSCIntegration) *RSCManager {
	manager := &RSCManager{
		sharedStateManager: sharedStateManager,
		eventStream:        eventStream,
		rscClient:          rscClient,
		rscIntegration:     rscIntegration,
		componentCache:     make(map[string]vercel.GeneratedComponent),
	}

	manager.setupEventListeners()

	return manager
}

func (m *RSCManager) setupEventListeners() {
	if m.eventStream == nil {
		return
	}

	m.eventStream.Subscribe(eventstream.EventTypeAgentAction, func(event *eventstream.Event) {
		if event.Data == nil {
			return
		}

		data, ok := event.Data.(map[string]interface{})
		if !ok {
			return
		}

		agentID, _ := data["agent_id"].(string)
		actionID, _ := data["action_id"].(string)
		actionType, _ := data["action_type"].(string)
		actionData, _ := data["action_data"].(map[string]interface{})

		if agentID == "" || actionID == "" || actionType == "" {
			return
		}

		go func() {
			component, err := m.rscClient.GenerateComponentFromAgentAction(
				context.Background(),
				agentID,
				actionID,
				actionType,
				actionData,
			)
			if err != nil {
				log.Printf("Error generating component from agent action: %v", err)
				return
			}

			m.storeComponent(component)

			m.storeComponentInSharedState(component)

			m.publishComponentEvent(component)
		}()
	})

	m.eventStream.Subscribe(eventstream.EventTypeToolUsage, func(event *eventstream.Event) {
		if event.Data == nil {
			return
		}

		data, ok := event.Data.(map[string]interface{})
		if !ok {
			return
		}

		agentID, _ := data["agent_id"].(string)
		toolID, _ := data["tool_id"].(string)
		toolName, _ := data["tool_name"].(string)
		toolInput, _ := data["tool_input"].(map[string]interface{})
		toolOutput, _ := data["tool_output"].(map[string]interface{})

		if agentID == "" || toolID == "" || toolName == "" {
			return
		}

		go func() {
			component, err := m.rscClient.GenerateComponentFromToolUsage(
				context.Background(),
				agentID,
				toolID,
				toolName,
				toolInput,
				toolOutput,
			)
			if err != nil {
				log.Printf("Error generating component from tool usage: %v", err)
				return
			}

			m.storeComponent(component)

			m.storeComponentInSharedState(component)

			m.publishComponentEvent(component)
		}()
	})
}

func (m *RSCManager) storeComponent(component vercel.GeneratedComponent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.componentCache[component.ID] = component
}

func (m *RSCManager) storeComponentInSharedState(component vercel.GeneratedComponent) {
	componentJSON, err := json.Marshal(component)
	if err != nil {
		log.Printf("Error marshaling component: %v", err)
		return
	}

	var componentData map[string]interface{}
	err = json.Unmarshal(componentJSON, &componentData)
	if err != nil {
		log.Printf("Error unmarshaling component: %v", err)
		return
	}

	err = m.sharedStateManager.UpdateState(StateTypeUI, component.ID, componentData)
	if err != nil {
		log.Printf("Error updating state: %v", err)
	}
}

func (m *RSCManager) publishComponentEvent(component vercel.GeneratedComponent) {
	if m.eventStream == nil {
		return
	}

	componentJSON, err := json.Marshal(component)
	if err != nil {
		log.Printf("Error marshaling component: %v", err)
		return
	}

	var componentData map[string]interface{}
	err = json.Unmarshal(componentJSON, &componentData)
	if err != nil {
		log.Printf("Error unmarshaling component: %v", err)
		return
	}

	event := &eventstream.Event{
		ID:        uuid.New().String(),
		Type:      eventstream.EventTypeComponentGenerated,
		Data:      componentData,
		Timestamp: time.Now().Unix(),
	}

	err = m.eventStream.Publish(event)
	if err != nil {
		log.Printf("Error publishing component event: %v", err)
	}
}

func (m *RSCManager) GenerateComponent(ctx context.Context, componentType vercel.ComponentType, props map[string]interface{}) (vercel.GeneratedComponent, error) {
	if m.rscClient == nil {
		return vercel.GeneratedComponent{}, fmt.Errorf("RSC client not initialized")
	}

	component, err := m.rscClient.GenerateComponent(ctx, componentType, props)
	if err != nil {
		return vercel.GeneratedComponent{}, fmt.Errorf("error generating component: %w", err)
	}

	m.storeComponent(component)

	m.storeComponentInSharedState(component)

	m.publishComponentEvent(component)

	return component, nil
}

func (m *RSCManager) GetComponent(ctx context.Context, componentID string) (vercel.GeneratedComponent, error) {
	m.mu.RLock()
	component, exists := m.componentCache[componentID]
	m.mu.RUnlock()

	if exists {
		return component, nil
	}

	if m.rscIntegration != nil {
		return m.rscIntegration.GetComponent(ctx, componentID)
	}

	return vercel.GeneratedComponent{}, fmt.Errorf("component not found: %s", componentID)
}

func (m *RSCManager) ListComponents(ctx context.Context) ([]vercel.GeneratedComponent, error) {
	m.mu.RLock()
	components := make([]vercel.GeneratedComponent, 0, len(m.componentCache))
	for _, component := range m.componentCache {
		components = append(components, component)
	}
	m.mu.RUnlock()

	return components, nil
}

func (m *RSCManager) GetComponentsByAgent(ctx context.Context, agentID string) ([]vercel.GeneratedComponent, error) {
	if m.rscIntegration != nil {
		return m.rscIntegration.GetComponentsByAgent(agentID)
	}

	m.mu.RLock()
	components := make([]vercel.GeneratedComponent, 0)
	for _, component := range m.componentCache {
		if component.AgentID == agentID {
			components = append(components, component)
		}
	}
	m.mu.RUnlock()

	return components, nil
}

func (m *RSCManager) GetComponentsByTool(ctx context.Context, toolID string) ([]vercel.GeneratedComponent, error) {
	if m.rscIntegration != nil {
		return m.rscIntegration.GetComponentsByTool(toolID)
	}

	m.mu.RLock()
	components := make([]vercel.GeneratedComponent, 0)
	for _, component := range m.componentCache {
		if component.ToolID == toolID {
			components = append(components, component)
		}
	}
	m.mu.RUnlock()

	return components, nil
}
