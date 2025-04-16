package models

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

type EventSource string

type EventPriority int

const (
	EventTypeMessage     EventType = "message"
	EventTypeAction      EventType = "action"
	EventTypeObservation EventType = "observation"
	EventTypePlan        EventType = "plan"
	EventTypeKnowledge   EventType = "knowledge"
	EventTypeDatasource  EventType = "datasource"
	EventTypeStateUpdate EventType = "state_update"
	EventTypeCacheUpdate EventType = "cache_update"
	EventTypeSystem      EventType = "system"
	EventTypeError       EventType = "error"
	EventTypeWarning     EventType = "warning"
	EventTypeInfo        EventType = "info"
	EventTypeDebug       EventType = "debug"
	EventTypeMetric      EventType = "metric"
	EventTypeAudit       EventType = "audit"
	EventTypeCommand     EventType = "command"
	EventTypeResponse    EventType = "response"

	EventSourceUser     EventSource = "user"
	EventSourceAgent    EventSource = "agent"
	EventSourceSystem   EventSource = "system"
	EventSourceModule   EventSource = "module"
	EventSourceCICD     EventSource = "ci_cd"
	EventSourceK8s      EventSource = "kubernetes"
	EventSourceSandbox  EventSource = "sandbox"
	EventSourceDatabase EventSource = "database"
	EventSourceAPI      EventSource = "api"
	EventSourceWeb      EventSource = "web"
	EventSourceMobile   EventSource = "mobile"
	EventSourceDesktop  EventSource = "desktop"
	EventSourceCLI      EventSource = "cli"
	EventSourceExternal EventSource = "external"
	EventSourceInternal EventSource = "internal"
	EventSourceMonitor  EventSource = "monitor"
	EventSourceScheduler EventSource = "scheduler"

	EventPriorityLow    EventPriority = 0
	EventPriorityNormal EventPriority = 1
	EventPriorityHigh   EventPriority = 2
	EventPriorityCritical EventPriority = 3
)

type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Source      EventSource            `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
	Priority    EventPriority          `json:"priority"`
	CorrelationID string               `json:"correlation_id,omitempty"`
	ParentID    string                 `json:"parent_id,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	AgentID     string                 `json:"agent_id,omitempty"`
	ModuleID    string                 `json:"module_id,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

func NewEvent(eventType EventType, source EventSource, data map[string]interface{}, metadata map[string]string) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now().UTC(),
		Data:      data,
		Metadata:  metadata,
		Priority:  EventPriorityNormal,
	}
}

func (e *Event) WithPriority(priority EventPriority) *Event {
	e.Priority = priority
	return e
}

func (e *Event) WithCorrelationID(correlationID string) *Event {
	e.CorrelationID = correlationID
	return e
}

func (e *Event) WithParentID(parentID string) *Event {
	e.ParentID = parentID
	return e
}

func (e *Event) WithSessionID(sessionID string) *Event {
	e.SessionID = sessionID
	return e
}

func (e *Event) WithUserID(userID string) *Event {
	e.UserID = userID
	return e
}

func (e *Event) WithAgentID(agentID string) *Event {
	e.AgentID = agentID
	return e
}

func (e *Event) WithModuleID(moduleID string) *Event {
	e.ModuleID = moduleID
	return e
}

func (e *Event) WithTags(tags ...string) *Event {
	e.Tags = tags
	return e
}

func (e *Event) AddTag(tag string) *Event {
	if e.Tags == nil {
		e.Tags = []string{}
	}
	e.Tags = append(e.Tags, tag)
	return e
}

func (e *Event) AddMetadata(key, value string) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

func (e *Event) AddData(key string, value interface{}) *Event {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = value
	return e
}

func (e *Event) Clone() *Event {
	dataCopy := make(map[string]interface{})
	for k, v := range e.Data {
		dataCopy[k] = v
	}

	metadataCopy := make(map[string]string)
	for k, v := range e.Metadata {
		metadataCopy[k] = v
	}

	tagsCopy := make([]string, len(e.Tags))
	copy(tagsCopy, e.Tags)

	return &Event{
		ID:            e.ID,
		Type:          e.Type,
		Source:        e.Source,
		Timestamp:     e.Timestamp,
		Data:          dataCopy,
		Metadata:      metadataCopy,
		Priority:      e.Priority,
		CorrelationID: e.CorrelationID,
		ParentID:      e.ParentID,
		SessionID:     e.SessionID,
		UserID:        e.UserID,
		AgentID:       e.AgentID,
		ModuleID:      e.ModuleID,
		Tags:          tagsCopy,
	}
}
