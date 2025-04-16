package models

import (
	"time"

	"github.com/google/uuid"
)

type StateType string

const (
	StateTypeAgent StateType = "agent"
	StateTypeUser StateType = "user"
	StateTypeSession StateType = "session"
	StateTypeSystem StateType = "system"
	StateTypeModule StateType = "module"
	StateTypeRuntime StateType = "runtime"
	StateTypeComponent StateType = "component"
	StateTypeUI StateType = "ui"
	StateTypeAction StateType = "action"
	StateTypeTool StateType = "tool"
)

type State struct {
	ID        string                 `json:"id"`
	Type      StateType              `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Version   int64                  `json:"version"`
}

func NewState(stateType StateType, data map[string]interface{}, metadata map[string]string) *State {
	now := time.Now().UTC()
	return &State{
		ID:        uuid.New().String(),
		Type:      stateType,
		Data:      data,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}
}

func (s *State) Update(data map[string]interface{}) {
	s.Data = data
	s.UpdatedAt = time.Now().UTC()
	s.Version++
}

func (s *State) Merge(data map[string]interface{}) {
	for k, v := range data {
		s.Data[k] = v
	}
	s.UpdatedAt = time.Now().UTC()
	s.Version++
}

func (s *State) GetValue(key string) (interface{}, bool) {
	value, exists := s.Data[key]
	return value, exists
}

func (s *State) SetValue(key string, value interface{}) {
	if s.Data == nil {
		s.Data = make(map[string]interface{})
	}
	s.Data[key] = value
	s.UpdatedAt = time.Now().UTC()
	s.Version++
}

func (s *State) AddMetadata(key, value string) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}
	s.Metadata[key] = value
}

func (s *State) Clone() *State {
	dataCopy := make(map[string]interface{})
	for k, v := range s.Data {
		dataCopy[k] = v
	}

	metadataCopy := make(map[string]string)
	for k, v := range s.Metadata {
		metadataCopy[k] = v
	}

	return &State{
		ID:        s.ID,
		Type:      s.Type,
		Data:      dataCopy,
		Metadata:  metadataCopy,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		Version:   s.Version,
	}
}
