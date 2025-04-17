package eventstream

import "time"

type EventType string

type EventSource string

const (
	EventTypeMessage    EventType = "message"
	EventTypeAction     EventType = "action"
	EventTypeObservation EventType = "observation"
	EventTypePlan       EventType = "plan"
	EventTypeKnowledge  EventType = "knowledge"
	EventTypeDatasource EventType = "datasource"
	EventTypeStateUpdate EventType = "state_update" // For RocketMQ
	EventTypeCacheUpdate EventType = "cache_update" // For DragonflyDB invalidation
	EventTypeNeovim     EventType = "neovim"        // For Neovim operations

	EventSourceUser    EventSource = "user"
	EventSourceAgent   EventSource = "agent"
	EventSourceSystem  EventSource = "system"
	EventSourceModule  EventSource = "module"
	EventSourceCICD    EventSource = "ci_cd"
	EventSourceK8s     EventSource = "kubernetes"
	EventSourceSandbox EventSource = "sandbox"
	EventSourceNeovim  EventSource = "neovim"       // For Neovim events
)

type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Source    EventSource `json:"source"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"` // Flexible data payload
	Metadata  map[string]string `json:"metadata,omitempty"` // Optional metadata
}

func NewEvent(eventType EventType, source EventSource, data interface{}, metadata map[string]string) *Event {
	return &Event{
		ID:        generateEventID(), // Implement ID generation (e.g., UUID)
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now().UTC(),
		Data:      data,
		Metadata:  metadata,
	}
}

func generateEventID() string {
	return time.Now().Format(time.RFC3339Nano)
}
