package klediosdk

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GoRuntimeAdapter struct {
	config *RuntimeAdapterConfig

	connected bool

	eventHandlers map[string]map[string]EventHandler

	mu sync.RWMutex
}

func NewGoRuntimeAdapter() *GoRuntimeAdapter {
	return &GoRuntimeAdapter{
		eventHandlers: make(map[string]map[string]EventHandler),
	}
}

func (a *GoRuntimeAdapter) Initialize(ctx context.Context, config *RuntimeAdapterConfig) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if config == nil {
		return errors.New("config cannot be nil")
	}

	a.config = config
	return nil
}

func (a *GoRuntimeAdapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.connected = true
	return nil
}

func (a *GoRuntimeAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.connected = false
	return nil
}

func (a *GoRuntimeAdapter) IsConnected() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.connected
}

func (a *GoRuntimeAdapter) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return nil, errors.New("adapter is not connected")
	}


	startTime := time.Now()

	time.Sleep(100 * time.Millisecond)

	result := &TaskResult{
		TaskID:        task.ID,
		AgentID:       task.AgentID,
		Status:        "completed",
		Output:        make(map[string]interface{}),
		Runtime:       RuntimeTypeGo,
		ExecutionTime: time.Since(startTime),
		Metadata:      make(map[string]interface{}),
	}

	result.Output["message"] = "Task executed successfully in Go runtime"
	result.Output["timestamp"] = time.Now().Format(time.RFC3339)

	result.Metadata["runtime_version"] = "go1.24.1"
	result.Metadata["execution_environment"] = "native"

	return result, nil
}

func (a *GoRuntimeAdapter) GetState(ctx context.Context, key string) (interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return nil, errors.New("adapter is not connected")
	}

	return "go-state-value-for-" + key, nil
}

func (a *GoRuntimeAdapter) SetState(ctx context.Context, key string, value interface{}) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return errors.New("adapter is not connected")
	}

	return nil
}

func (a *GoRuntimeAdapter) SubscribeToEvents(ctx context.Context, eventType string, handler EventHandler) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return "", errors.New("adapter is not connected")
	}

	subscriptionID := uuid.New().String()

	if _, ok := a.eventHandlers[eventType]; !ok {
		a.eventHandlers[eventType] = make(map[string]EventHandler)
	}

	a.eventHandlers[eventType][subscriptionID] = handler

	return subscriptionID, nil
}

func (a *GoRuntimeAdapter) UnsubscribeFromEvents(ctx context.Context, subscriptionID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return errors.New("adapter is not connected")
	}

	for eventType, handlers := range a.eventHandlers {
		if _, ok := handlers[subscriptionID]; ok {
			delete(a.eventHandlers[eventType], subscriptionID)
			return nil
		}
	}

	return errors.New("subscription not found")
}

func (a *GoRuntimeAdapter) PublishEvent(ctx context.Context, event *Event) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return errors.New("adapter is not connected")
	}

	handlers, ok := a.eventHandlers[event.Type]
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			_ = h(context.Background(), event)
		}(handler)
	}

	return nil
}
