package langgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type StateManager struct {
	initialState map[string]interface{}
	schema       map[string]string
	history      []map[string]interface{}
	mu           sync.RWMutex
}

func NewStateManager(initialState map[string]interface{}, schema map[string]string) *StateManager {
	return &StateManager{
		initialState: initialState,
		schema:       schema,
		history:      []map[string]interface{}{},
	}
}

func (sm *StateManager) GetState() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	if len(sm.history) == 0 {
		return sm.initialState
	}
	
	return sm.history[len(sm.history)-1]
}

func (sm *StateManager) UpdateState(updates map[string]interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	currentState := sm.GetState()
	newState := make(map[string]interface{})
	
	for k, v := range currentState {
		newState[k] = v
	}
	
	for k, v := range updates {
		newState[k] = v
	}
	
	for key, expectedType := range sm.schema {
		val, exists := newState[key]
		if !exists {
			return fmt.Errorf("required state key %s is missing", key)
		}
		
		switch expectedType {
		case "string":
			if _, ok := val.(string); !ok {
				return fmt.Errorf("state key %s should be a string", key)
			}
		case "int":
			if _, ok := val.(int); !ok {
				return fmt.Errorf("state key %s should be an int", key)
			}
		case "bool":
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("state key %s should be a bool", key)
			}
		case "map":
			if _, ok := val.(map[string]interface{}); !ok {
				return fmt.Errorf("state key %s should be a map", key)
			}
		case "array":
			if _, ok := val.([]interface{}); !ok {
				return fmt.Errorf("state key %s should be an array", key)
			}
		}
	}
	
	sm.history = append(sm.history, newState)
	
	return nil
}

func (sm *StateManager) GetHistory() []map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	return sm.history
}

func (sm *StateManager) GetTrajectory() (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	trajectory, err := json.Marshal(sm.history)
	if err != nil {
		return "", err
	}
	
	return string(trajectory), nil
}

func (sm *StateManager) Reset() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.history = []map[string]interface{}{}
}

type StateUpdater func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error)

func ChainStateUpdaters(updaters ...StateUpdater) StateUpdater {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		currentState := state
		var err error
		
		for _, updater := range updaters {
			currentState, err = updater(ctx, currentState)
			if err != nil {
				return nil, err
			}
		}
		
		return currentState, nil
	}
}

func ConditionalStateUpdater(condition func(ctx context.Context, state map[string]interface{}) bool, updater StateUpdater) StateUpdater {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		if condition(ctx, state) {
			return updater(ctx, state)
		}
		return state, nil
	}
}

func KeyValueStateUpdater(key string, value interface{}) StateUpdater {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		newState := make(map[string]interface{})
		for k, v := range state {
			newState[k] = v
		}
		newState[key] = value
		return newState, nil
	}
}

func DynamicKeyValueStateUpdater(key string, valueFn func(ctx context.Context, state map[string]interface{}) (interface{}, error)) StateUpdater {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		value, err := valueFn(ctx, state)
		if err != nil {
			return nil, err
		}
		
		newState := make(map[string]interface{})
		for k, v := range state {
			newState[k] = v
		}
		newState[key] = value
		return newState, nil
	}
}
