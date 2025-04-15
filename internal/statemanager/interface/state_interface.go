package stateinterface

import (
	"context"
	"time"
)

type StateType string

const (
	StateTypeTask StateType = "task"
	
	StateTypeAgent StateType = "agent"
	
	StateTypeLifecycle StateType = "lifecycle"
	
	StateTypeShared StateType = "shared"
)

type StateManager interface {
	SaveState(stateType StateType, stateID string, stateData []byte) error
	
	GetState(stateType StateType, stateID string) ([]byte, error)
	
	DeleteState(stateType StateType, stateID string) error
	
	SubscribeToStateChanges(ctx context.Context, stateType StateType, stateID string) (<-chan []byte, error)
	
	PublishStateChange(stateType StateType, stateID string, stateData []byte) error
	
	Close() error
}

type StateEntry struct {
	ID        string    `json:"id"`
	StateType string    `json:"state_type"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StateChangeEvent struct {
	StateType StateType `json:"state_type"`
	StateID   string    `json:"state_id"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
