package sharedstate

import (
	"fmt"
	"log"
	"time"
)


type MockEvent struct {
	Type      string
	Data      interface{}
	Timestamp time.Time
}

const (
	MockEventTypeStateUpdate = "state_update"
	MockEventTypeUserAction  = "user_action"
	MockEventTypeAgentAction = "agent_action"
)

type MockStream struct {
	subscribers map[string][]chan *MockEvent
}

func NewMockStream() *MockStream {
	return &MockStream{
		subscribers: make(map[string][]chan *MockEvent),
	}
}

func (s *MockStream) Subscribe(eventType string, ch chan *MockEvent) {
	if _, ok := s.subscribers[eventType]; !ok {
		s.subscribers[eventType] = make([]chan *MockEvent, 0)
	}
	s.subscribers[eventType] = append(s.subscribers[eventType], ch)
}

func (s *MockStream) Emit(event *MockEvent) {
	if channels, ok := s.subscribers[event.Type]; ok {
		for _, ch := range channels {
			select {
			case ch <- event:
			default:
				log.Printf("Channel full, dropping event of type %s", event.Type)
			}
		}
	}
}

type MockSupabaseClient struct{}

func NewMockSupabaseClient() *MockSupabaseClient {
	return &MockSupabaseClient{}
}

func (c *MockSupabaseClient) GetState(stateType, stateID string) ([]byte, error) {
	mockData := fmt.Sprintf(`{"state_type":"%s","state_id":"%s","status":"mock","timestamp":"%s"}`,
		stateType, stateID, time.Now().UTC().Format(time.RFC3339))
	return []byte(mockData), nil
}

func (c *MockSupabaseClient) SaveState(stateType, stateID string, data []byte) error {
	log.Printf("Mock Supabase: Saved state %s:%s", stateType, stateID)
	return nil
}

type MockStateManager struct {
	states map[string][]byte
}

func NewMockStateManager() *MockStateManager {
	return &MockStateManager{
		states: make(map[string][]byte),
	}
}

func (m *MockStateManager) GetState(stateType, stateID string) ([]byte, bool) {
	key := fmt.Sprintf("%s:%s", stateType, stateID)
	data, exists := m.states[key]
	return data, exists
}

func (m *MockStateManager) UpdateState(stateID, stateType string, data []byte) error {
	key := fmt.Sprintf("%s:%s", stateType, stateID)
	m.states[key] = data
	log.Printf("Mock StateManager: Updated state %s", key)
	return nil
}

func (m *MockStateManager) Close() error {
	return nil
}
