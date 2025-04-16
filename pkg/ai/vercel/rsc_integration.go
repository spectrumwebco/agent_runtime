package vercel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/djangobridge"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/klediosdk"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type RSCIntegration struct {
	rscAdapter         *RSCAdapter
	sharedStateManager *sharedstate.SharedStateManager
	eventStream        *eventstream.Stream
	djangoIntegration  *djangobridge.DjangoIntegration
	multiAgentSystem   *langraph.MultiAgentSystem
	sdk                *klediosdk.SDK

	agentComponentMap map[string][]string // Maps agent IDs to component IDs
	toolComponentMap  map[string][]string // Maps tool IDs to component IDs

	mu sync.RWMutex
}

func NewRSCIntegration(
	rscAdapter *RSCAdapter,
	sharedStateManager *sharedstate.SharedStateManager,
	eventStream *eventstream.Stream,
	djangoIntegration *djangobridge.DjangoIntegration,
	multiAgentSystem *langraph.MultiAgentSystem,
	sdk *klediosdk.SDK,
) *RSCIntegration {
	integration := &RSCIntegration{
		rscAdapter:         rscAdapter,
		sharedStateManager: sharedStateManager,
		eventStream:        eventStream,
		djangoIntegration:  djangoIntegration,
		multiAgentSystem:   multiAgentSystem,
		sdk:                sdk,
		agentComponentMap:  make(map[string][]string),
		toolComponentMap:   make(map[string][]string),
	}

	integration.setupEventListeners()

	return integration
}

func (i *RSCIntegration) setupEventListeners() {
	if i.eventStream == nil {
		return
	}

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
			component, err := i.rscAdapter.GenerateComponentFromAgentAction(
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

			i.mapComponentToAgent(actionEvent.AgentID, component.ID)
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
			component, err := i.rscAdapter.GenerateComponentFromToolUsage(
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

			i.mapComponentToAgent(toolEvent.AgentID, component.ID)
			i.mapComponentToTool(toolEvent.ToolID, component.ID)
		}()
	})
}

func (i *RSCIntegration) mapComponentToAgent(agentID string, componentID string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.agentComponentMap[agentID]; !exists {
		i.agentComponentMap[agentID] = make([]string, 0)
	}

	i.agentComponentMap[agentID] = append(i.agentComponentMap[agentID], componentID)
}

func (i *RSCIntegration) mapComponentToTool(toolID string, componentID string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.toolComponentMap[toolID]; !exists {
		i.toolComponentMap[toolID] = make([]string, 0)
	}

	i.toolComponentMap[toolID] = append(i.toolComponentMap[toolID], componentID)
}

func (i *RSCIntegration) GetComponentsByAgent(agentID string) ([]GeneratedComponent, error) {
	i.mu.RLock()
	componentIDs, exists := i.agentComponentMap[agentID]
	i.mu.RUnlock()

	if !exists {
		return []GeneratedComponent{}, nil
	}

	components := make([]GeneratedComponent, 0, len(componentIDs))
	for _, componentID := range componentIDs {
		component, err := i.rscAdapter.GetComponent(context.Background(), componentID)
		if err != nil {
			continue
		}

		components = append(components, component)
	}

	return components, nil
}

func (i *RSCIntegration) GetComponentsByTool(toolID string) ([]GeneratedComponent, error) {
	i.mu.RLock()
	componentIDs, exists := i.toolComponentMap[toolID]
	i.mu.RUnlock()

	if !exists {
		return []GeneratedComponent{}, nil
	}

	components := make([]GeneratedComponent, 0, len(componentIDs))
	for _, componentID := range componentIDs {
		component, err := i.rscAdapter.GetComponent(context.Background(), componentID)
		if err != nil {
			continue
		}

		components = append(components, component)
	}

	return components, nil
}

func (i *RSCIntegration) GenerateComponent(ctx context.Context, componentType ComponentType, props map[string]interface{}) (GeneratedComponent, error) {
	return i.rscAdapter.GenerateComponent(ctx, componentType, props)
}

func (i *RSCIntegration) GetComponent(ctx context.Context, componentID string) (GeneratedComponent, error) {
	return i.rscAdapter.GetComponent(ctx, componentID)
}

func (i *RSCIntegration) ListComponents(ctx context.Context) ([]GeneratedComponent, error) {
	return i.rscAdapter.GetComponentHistory(), nil
}

func (i *RSCIntegration) RegisterDjangoIntegration(djangoIntegration *djangobridge.DjangoIntegration) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.djangoIntegration = djangoIntegration
}

func (i *RSCIntegration) RegisterMultiAgentSystem(multiAgentSystem *langraph.MultiAgentSystem) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.multiAgentSystem = multiAgentSystem
}

func (i *RSCIntegration) RegisterSDK(sdk *klediosdk.SDK) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.sdk = sdk
}
