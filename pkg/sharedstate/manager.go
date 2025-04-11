package sharedstate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
)

type StateType string

const (
	StateTypeAgent StateType = "agent"
	StateTypeTask StateType = "task"
	StateTypeUI StateType = "ui"
)

type SharedStateManager struct {
	wsManager      *WebSocketManager
	eventStream    *eventstream.Stream
	stateManager   *statemanager.StateManager
	stateCache     map[string]map[string]interface{}
	stateCacheLock sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

type SharedStateConfig struct {
	EventStream  *eventstream.Stream
	StateManager *statemanager.StateManager
}

func NewSharedStateManager(cfg SharedStateConfig) *SharedStateManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	manager := &SharedStateManager{
		wsManager:    NewWebSocketManager(),
		eventStream:  cfg.EventStream,
		stateManager: cfg.StateManager,
		stateCache:   make(map[string]map[string]interface{}),
		ctx:          ctx,
		cancel:       cancel,
	}
	
	go manager.wsManager.Start(ctx)
	
	eventChan := make(chan *eventstream.Event, 100)
	manager.eventStream.Subscribe(eventstream.EventTypeStateUpdate, eventChan)
	go manager.handleEvents(eventChan)
	
	return manager
}

func (m *SharedStateManager) handleEvents(eventChan chan *eventstream.Event) {
	for {
		select {
		case <-m.ctx.Done():
			return
		case event := <-eventChan:
			if event.Type == eventstream.EventTypeStateUpdate {
				if data, ok := event.Data.(map[string]interface{}); ok {
					if stateType, ok := data["state_type"].(string); ok {
						if stateID, ok := data["state_id"].(string); ok {
							m.updateStateCache(StateType(stateType), stateID, data)
							
							m.wsManager.PublishToTopic(fmt.Sprintf("state:%s:%s", stateType, stateID), data)
						}
					}
				}
			}
		}
	}
}

func (m *SharedStateManager) updateStateCache(stateType StateType, stateID string, data map[string]interface{}) {
	m.stateCacheLock.Lock()
	defer m.stateCacheLock.Unlock()
	
	key := fmt.Sprintf("%s:%s", stateType, stateID)
	m.stateCache[key] = data
}

func (m *SharedStateManager) GetState(stateType StateType, stateID string) (map[string]interface{}, error) {
	m.stateCacheLock.RLock()
	key := fmt.Sprintf("%s:%s", stateType, stateID)
	if state, ok := m.stateCache[key]; ok {
		m.stateCacheLock.RUnlock()
		return state, nil
	}
	m.stateCacheLock.RUnlock()
	
	stateData, err := m.stateManager.GetState(string(stateType), stateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	
	var state map[string]interface{}
	if err := json.Unmarshal(stateData, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}
	
	m.updateStateCache(stateType, stateID, state)
	
	return state, nil
}

func (m *SharedStateManager) UpdateState(stateType StateType, stateID string, data map[string]interface{}) error {
	m.updateStateCache(stateType, stateID, data)
	
	stateData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	
	if err := m.stateManager.UpdateState(stateID, string(stateType), stateData); err != nil {
		return fmt.Errorf("failed to update state: %w", err)
	}
	
	m.wsManager.PublishToTopic(fmt.Sprintf("state:%s:%s", stateType, stateID), data)
	
	return nil
}

func (m *SharedStateManager) RegisterWebSocketHandler(router gin.IRouter) {
	router.GET("/ws", m.wsManager.HandleWebSocket)
}

func (m *SharedStateManager) Close() error {
	m.cancel() // Signal background tasks to stop
	return nil
}
