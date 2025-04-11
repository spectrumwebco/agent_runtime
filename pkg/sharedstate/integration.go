package sharedstate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager/supabase"
)

type Integration struct {
	sharedStateManager *SharedStateManager
	eventStream        *eventstream.Stream
	stateManager       *statemanager.StateManager
	supabaseClient     *supabase.Client
	ctx                context.Context
	cancel             context.CancelFunc
	mutex              sync.RWMutex
}

type IntegrationConfig struct {
	EventStream        *eventstream.Stream
	StateManager       *statemanager.StateManager
	SupabaseClient     *supabase.Client
	SharedStateManager *SharedStateManager
}

func NewIntegration(cfg IntegrationConfig) (*Integration, error) {
	if cfg.EventStream == nil {
		return nil, fmt.Errorf("event stream is required")
	}
	if cfg.StateManager == nil {
		return nil, fmt.Errorf("state manager is required")
	}
	if cfg.SharedStateManager == nil {
		return nil, fmt.Errorf("shared state manager is required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	integration := &Integration{
		sharedStateManager: cfg.SharedStateManager,
		eventStream:        cfg.EventStream,
		stateManager:       cfg.StateManager,
		supabaseClient:     cfg.SupabaseClient,
		ctx:                ctx,
		cancel:             cancel,
	}

	eventChan := make(chan *eventstream.Event, 100)
	cfg.EventStream.Subscribe(eventstream.EventTypeStateUpdate, eventChan)
	go integration.handleStateUpdateEvents(eventChan)

	clientEventChan := make(chan map[string]interface{}, 100)
	integration.registerClientEventHandler(clientEventChan)
	go integration.handleClientEvents(clientEventChan)

	return integration, nil
}

func (i *Integration) registerClientEventHandler(eventChan chan<- map[string]interface{}) {
	handleClientEvent = func(eventData map[string]interface{}) {
		select {
		case eventChan <- eventData:
		default:
			log.Printf("Client event channel full, dropping event")
		}
	}
}

func (i *Integration) handleStateUpdateEvents(eventChan <-chan *eventstream.Event) {
	for {
		select {
		case <-i.ctx.Done():
			return
		case event := <-eventChan:
			if event.Type == eventstream.EventTypeStateUpdate {
				if data, ok := event.Data.(map[string]interface{}); ok {
					if stateType, ok := data["state_type"].(string); ok {
						if stateID, ok := data["state_id"].(string); ok {
							i.sharedStateManager.UpdateState(StateType(stateType), stateID, data)
						}
					}
				}
			}
		}
	}
}

func (i *Integration) handleClientEvents(eventChan <-chan map[string]interface{}) {
	for {
		select {
		case <-i.ctx.Done():
			return
		case eventData := <-eventChan:
			i.processClientEvent(eventData)
		}
	}
}

func (i *Integration) processClientEvent(eventData map[string]interface{}) {
	eventType, ok := eventData["type"].(string)
	if !ok {
		log.Printf("Client event missing type: %v", eventData)
		return
	}

	switch eventType {
	case "state_update":
		i.handleClientStateUpdate(eventData)
	default:
		i.forwardEventToEventStream(eventType, eventData)
	}
}

func (i *Integration) handleClientStateUpdate(eventData map[string]interface{}) {
	stateType, ok := eventData["state_type"].(string)
	if !ok {
		log.Printf("State update event missing state_type: %v", eventData)
		return
	}

	stateID, ok := eventData["state_id"].(string)
	if !ok {
		log.Printf("State update event missing state_id: %v", eventData)
		return
	}

	data, ok := eventData["data"].(map[string]interface{})
	if !ok {
		log.Printf("State update event missing data: %v", eventData)
		return
	}

	data["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	data["updated_by"] = "client"

	stateData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal state data: %v", err)
		return
	}

	if err := i.stateManager.UpdateState(stateID, stateType, stateData); err != nil {
		log.Printf("Failed to update state: %v", err)
		return
	}

	i.eventStream.Emit(&eventstream.Event{
		Type: eventstream.EventTypeStateUpdate,
		Data: map[string]interface{}{
			"state_type": stateType,
			"state_id":   stateID,
			"data":       data,
			"source":     "client",
		},
	})
}

func (i *Integration) forwardEventToEventStream(eventType string, eventData map[string]interface{}) {
	eventData["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	eventData["source"] = "client"

	i.eventStream.Emit(&eventstream.Event{
		Type: eventType,
		Data: eventData,
	})
}

func (i *Integration) Close() error {
	i.cancel()
	return nil
}

var handleClientEvent func(eventData map[string]interface{})
