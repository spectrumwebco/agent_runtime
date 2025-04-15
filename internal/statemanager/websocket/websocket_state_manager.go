package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager/interface"
)

type WebSocketStateManager struct {
	supabaseManager _interface.StateManager
	clients         map[string][]*websocket.Conn
	clientsMutex    sync.RWMutex
	eventChan       chan _interface.StateChangeEvent
	done            chan struct{}
}

func NewWebSocketStateManager(supabaseManager _interface.StateManager) *WebSocketStateManager {
	wsm := &WebSocketStateManager{
		supabaseManager: supabaseManager,
		clients:         make(map[string][]*websocket.Conn),
		eventChan:       make(chan _interface.StateChangeEvent, 100),
		done:            make(chan struct{}),
	}

	go wsm.eventLoop()

	return wsm
}

func (wsm *WebSocketStateManager) SaveState(stateType _interface.StateType, stateID string, stateData []byte) error {
	if err := wsm.supabaseManager.SaveState(stateType, stateID, stateData); err != nil {
		return err
	}

	wsm.PublishStateChange(stateType, stateID, stateData)

	return nil
}

func (wsm *WebSocketStateManager) GetState(stateType _interface.StateType, stateID string) ([]byte, error) {
	return wsm.supabaseManager.GetState(stateType, stateID)
}

func (wsm *WebSocketStateManager) DeleteState(stateType _interface.StateType, stateID string) error {
	if err := wsm.supabaseManager.DeleteState(stateType, stateID); err != nil {
		return err
	}

	wsm.PublishStateChange(stateType, stateID, nil)

	return nil
}

func (wsm *WebSocketStateManager) RegisterClient(stateType _interface.StateType, stateID string, conn *websocket.Conn) {
	wsm.clientsMutex.Lock()
	defer wsm.clientsMutex.Unlock()

	key := fmt.Sprintf("%s:%s", stateType, stateID)
	wsm.clients[key] = append(wsm.clients[key], conn)

	log.Printf("Registered WebSocket client for %s state with ID %s\n", stateType, stateID)
}

func (wsm *WebSocketStateManager) UnregisterClient(stateType _interface.StateType, stateID string, conn *websocket.Conn) {
	wsm.clientsMutex.Lock()
	defer wsm.clientsMutex.Unlock()

	key := fmt.Sprintf("%s:%s", stateType, stateID)
	clients := wsm.clients[key]
	for i, client := range clients {
		if client == conn {
			wsm.clients[key] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	log.Printf("Unregistered WebSocket client for %s state with ID %s\n", stateType, stateID)
}

func (wsm *WebSocketStateManager) SubscribeToStateChanges(ctx context.Context, stateType _interface.StateType, stateID string) (<-chan []byte, error) {
	ch := make(chan []byte, 10)

	go func() {
		defer close(ch)

		subscriptionID := fmt.Sprintf("%s:%s:%d", stateType, stateID, time.Now().UnixNano())

		eventCh := make(chan _interface.StateChangeEvent, 10)

		wsm.clientsMutex.Lock()
		key := fmt.Sprintf("%s:%s", stateType, stateID)
		wsm.clients[key] = append(wsm.clients[key], nil) // Use nil as a placeholder for non-WebSocket subscribers
		wsm.clientsMutex.Unlock()

		defer func() {
			wsm.clientsMutex.Lock()
			clients := wsm.clients[key]
			for i, client := range clients {
				if client == nil {
					wsm.clients[key] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			wsm.clientsMutex.Unlock()
		}()

		initialState, err := wsm.GetState(stateType, stateID)
		if err == nil && initialState != nil {
			select {
			case ch <- initialState:
			case <-ctx.Done():
				return
			}
		}

		for {
			select {
			case event := <-eventCh:
				if event.StateType == stateType && event.StateID == stateID {
					select {
					case ch <- event.Data:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

func (wsm *WebSocketStateManager) PublishStateChange(stateType _interface.StateType, stateID string, stateData []byte) error {
	event := _interface.StateChangeEvent{
		StateType: stateType,
		StateID:   stateID,
		Data:      stateData,
		Timestamp: time.Now(),
	}

	select {
	case wsm.eventChan <- event:
		return nil
	default:
		return fmt.Errorf("event channel is full")
	}
}

func (wsm *WebSocketStateManager) eventLoop() {
	for {
		select {
		case event := <-wsm.eventChan:
			wsm.notifySubscribers(event)
		case <-wsm.done:
			return
		}
	}
}

func (wsm *WebSocketStateManager) notifySubscribers(event _interface.StateChangeEvent) {
	wsm.clientsMutex.RLock()
	defer wsm.clientsMutex.RUnlock()

	key := fmt.Sprintf("%s:%s", event.StateType, event.StateID)
	clients := wsm.clients[key]

	if len(clients) == 0 {
		return
	}

	message, err := json.Marshal(map[string]interface{}{
		"type":       "state_change",
		"state_type": event.StateType,
		"state_id":   event.StateID,
		"data":       json.RawMessage(event.Data),
		"timestamp":  event.Timestamp,
	})

	if err != nil {
		log.Printf("Error marshaling state change event: %v\n", err)
		return
	}

	for _, client := range clients {
		if client != nil {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error sending state change to WebSocket client: %v\n", err)
			}
		}
	}
}

func (wsm *WebSocketStateManager) Close() error {
	close(wsm.done)
	return wsm.supabaseManager.Close()
}
