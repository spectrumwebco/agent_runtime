package langraph

import (
	"time"
)

type EventType string

type EventSource string

const (
	EventTypeMessage EventType = "message"
	EventTypeAction  EventType = "action"
	EventTypeSystem  EventType = "system"
	
	EventSourceAgent  EventSource = "agent"
	EventSourceSystem EventSource = "system"
	EventSourceUser   EventSource = "user"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Source    EventSource            `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
}
