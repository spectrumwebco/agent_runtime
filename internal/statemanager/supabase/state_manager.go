package supabase

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type StateType string

const (
	StateTypeTask      StateType = "task"
	StateTypeAgent     StateType = "agent"
	StateTypeLifecycle StateType = "lifecycle"
)

type SupabaseStateManager struct {
	mainClient     *SupabaseClient
	readonlyClient *SupabaseClient
	rollbackClient *SupabaseClient
	stateLock      sync.RWMutex
}

type SupabaseStateConfig struct {
	MainURL       string
	ReadonlyURL   string
	RollbackURL   string
	APIKey        string
	AuthToken     string
}

func NewSupabaseStateManager(cfg SupabaseStateConfig) *SupabaseStateManager {
	mainClient := NewSupabaseClient(SupabaseConfig{
		URL:       cfg.MainURL,
		APIKey:    cfg.APIKey,
		AuthToken: cfg.AuthToken,
	})
	
	readonlyClient := NewSupabaseClient(SupabaseConfig{
		URL:       cfg.ReadonlyURL,
		APIKey:    cfg.APIKey,
		AuthToken: cfg.AuthToken,
	})
	
	rollbackClient := NewSupabaseClient(SupabaseConfig{
		URL:       cfg.RollbackURL,
		APIKey:    cfg.APIKey,
		AuthToken: cfg.AuthToken,
	})
	
	return &SupabaseStateManager{
		mainClient:     mainClient,
		readonlyClient: readonlyClient,
		rollbackClient: rollbackClient,
	}
}

func (sm *SupabaseStateManager) SaveState(stateType StateType, stateID string, stateData []byte) error {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()
	
	var stateMap map[string]interface{}
	if err := json.Unmarshal(stateData, &stateMap); err != nil {
		return fmt.Errorf("failed to unmarshal state data: %w", err)
	}
	
	stateMap["id"] = stateID
	stateMap["state_type"] = string(stateType)
	stateMap["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	
	currentState, err := sm.readonlyClient.GetState(string(stateType), stateID)
	if err == nil {
		if err := sm.rollbackClient.UpsertState(string(stateType), currentState); err != nil {
			log.Printf("Warning: Failed to save state to rollback instance: %v\n", err)
		}
	}
	
	if err := sm.mainClient.UpsertState(string(stateType), stateMap); err != nil {
		return fmt.Errorf("failed to save state to main instance: %w", err)
	}
	
	log.Printf("Successfully saved %s state with ID %s\n", stateType, stateID)
	return nil
}

func (sm *SupabaseStateManager) GetState(stateType StateType, stateID string) ([]byte, error) {
	sm.stateLock.RLock()
	defer sm.stateLock.RUnlock()
	
	stateMap, err := sm.readonlyClient.GetState(string(stateType), stateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from readonly instance: %w", err)
	}
	
	stateData, err := json.Marshal(stateMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state data: %w", err)
	}
	
	return stateData, nil
}

func (sm *SupabaseStateManager) RollbackState(stateType StateType, stateID string) error {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()
	
	stateMap, err := sm.rollbackClient.GetState(string(stateType), stateID)
	if err != nil {
		return fmt.Errorf("failed to get state from rollback instance: %w", err)
	}
	
	if err := sm.mainClient.UpsertState(string(stateType), stateMap); err != nil {
		return fmt.Errorf("failed to save rollback state to main instance: %w", err)
	}
	
	log.Printf("Successfully rolled back %s state with ID %s\n", stateType, stateID)
	return nil
}

func (sm *SupabaseStateManager) DeleteState(stateType StateType, stateID string) error {
	sm.stateLock.Lock()
	defer sm.stateLock.Unlock()
	
	if err := sm.mainClient.DeleteState(string(stateType), stateID); err != nil {
		log.Printf("Warning: Failed to delete state from main instance: %v\n", err)
	}
	
	if err := sm.readonlyClient.DeleteState(string(stateType), stateID); err != nil {
		log.Printf("Warning: Failed to delete state from readonly instance: %v\n", err)
	}
	
	if err := sm.rollbackClient.DeleteState(string(stateType), stateID); err != nil {
		log.Printf("Warning: Failed to delete state from rollback instance: %v\n", err)
	}
	
	log.Printf("Successfully deleted %s state with ID %s\n", stateType, stateID)
	return nil
}
