package eventstream

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	Source    string      `json:"source,omitempty"`
}

func NewEvent(eventType string, data interface{}) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

const (
	EventTypeStateUpdate        = "state_update"
	EventTypeAgentAction        = "agent_action"
	EventTypeToolUsage          = "tool_usage"
	EventTypeComponentGenerated = "component_generated"
	EventTypeComponentUpdated   = "component_updated"
)

type Stream struct {
	subscribers map[string][]chan *Event
	mu          sync.RWMutex
}

func NewStream() *Stream {
	return &Stream{
		subscribers: make(map[string][]chan *Event),
	}
}

func (s *Stream) Subscribe(eventType string, ch chan *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.subscribers[eventType]; !ok {
		s.subscribers[eventType] = make([]chan *Event, 0)
	}

	s.subscribers[eventType] = append(s.subscribers[eventType], ch)
}

func (s *Stream) Publish(event *Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if subscribers, ok := s.subscribers[event.Type]; ok {
		for _, ch := range subscribers {
			select {
			case ch <- event:
			default:
			}
		}
	}

	return nil
}

func (s *Stream) SubscribeWithCallback(eventType string, callback func(*Event)) {
	ch := make(chan *Event, 100)
	s.Subscribe(eventType, ch)

	go func() {
		for event := range ch {
			callback(event)
		}
	}()
}

func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribers = make(map[string][]chan *Event)
}
