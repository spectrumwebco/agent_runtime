package neovim

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type StateManager struct {
	redisClient      RedisClient
	supabaseClient   SupabaseClient
	mainSupabase     SupabaseClient
	readOnlySupabase SupabaseClient
	rollbackSupabase SupabaseClient
	eventStream      EventStream
}

type RedisClient interface {
	HMSet(key string, data map[string]interface{}) error
	HGetAll(key string) (map[string]interface{}, error)
	Del(key string) error
	Expire(key string, ttl int) error
}

type SupabaseClient interface {
	StoreState(taskID string, stateType string, stateID string, data map[string]interface{}) error
	GetState(taskID string, stateType string, stateID string) (map[string]interface{}, bool)
	ExpireState(taskID string, stateType string, stateID string) error
}

type EventStream interface {
	PublishEvent(taskID string, event map[string]interface{}) error
}

func NewStateManager(
	redisClient RedisClient,
	supabaseClient SupabaseClient,
	mainSupabase SupabaseClient,
	readOnlySupabase SupabaseClient,
	rollbackSupabase SupabaseClient,
	eventStream EventStream,
) *StateManager {
	return &StateManager{
		redisClient:      redisClient,
		supabaseClient:   supabaseClient,
		mainSupabase:     mainSupabase,
		readOnlySupabase: readOnlySupabase,
		rollbackSupabase: rollbackSupabase,
		eventStream:      eventStream,
	}
}

func (sm *StateManager) StoreNeovimState(stateType StateType, stateID string, data map[string]interface{}, taskID string) error {
	data["last_updated"] = time.Now().UTC().Format(time.RFC3339)

	ttl := TTLMap[stateType]
	key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)

	err := sm.redisClient.HMSet(key, data)
	if err != nil {
		return err
	}

	err = sm.redisClient.Expire(key, ttl)
	if err != nil {
		return err
	}

	err = sm.storeNeovimStateInSupabase(taskID, string(stateType), stateID, data)
	if err != nil {
		log.Printf("Failed to store state in Supabase: %v", err)
	}

	return sm.publishStateUpdate(taskID, string(stateType), stateID, data)
}

func (sm *StateManager) RetrieveNeovimState(stateType StateType, stateID string, taskID string) (map[string]interface{}, bool) {
	key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)

	data, err := sm.redisClient.HGetAll(key)
	if err == nil && len(data) > 0 {
		return data, true
	}

	data, exists := sm.retrieveNeovimStateFromSupabase(taskID, string(stateType), stateID)
	if exists {
		sm.redisClient.HMSet(key, data)
		ttl := TTLMap[stateType]
		sm.redisClient.Expire(key, ttl)
		return data, true
	}

	return nil, false
}

func (sm *StateManager) UpdateNeovimState(stateType StateType, stateID string, updates map[string]interface{}, taskID string) error {
	key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)

	currentState, exists := sm.RetrieveNeovimState(stateType, stateID, taskID)
	if !exists {
		return fmt.Errorf("state not found: %s", key)
	}

	for k, v := range updates {
		currentState[k] = v
	}

	currentState["last_updated"] = time.Now().UTC().Format(time.RFC3339)

	return sm.StoreNeovimState(stateType, stateID, currentState, taskID)
}

func (sm *StateManager) MergeNeovimStates(sourceStateType StateType, sourceStateID string, targetStateType StateType, targetStateID string, taskID string) error {
	sourceState, exists := sm.RetrieveNeovimState(sourceStateType, sourceStateID, taskID)
	if !exists {
		return fmt.Errorf("source state not found: %s:%s", sourceStateType, sourceStateID)
	}

	targetState, exists := sm.RetrieveNeovimState(targetStateType, targetStateID, taskID)
	if !exists {
		return fmt.Errorf("target state not found: %s:%s", targetStateType, targetStateID)
	}

	for k, v := range sourceState {
		if k == "last_updated" || k == "state_type" || k == "state_id" {
			continue
		}
		targetState[k] = v
	}

	targetState["last_updated"] = time.Now().UTC().Format(time.RFC3339)

	return sm.StoreNeovimState(targetStateType, targetStateID, targetState, taskID)
}

func (sm *StateManager) ExpireNeovimState(stateType StateType, stateID string, taskID string) error {
	key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)

	err := sm.redisClient.Del(key)
	if err != nil {
		return err
	}

	return sm.supabaseClient.ExpireState(taskID, string(stateType), stateID)
}

func (sm *StateManager) GetAllNeovimStates(taskID string) (map[string]map[string]interface{}, error) {
	return map[string]map[string]interface{}{}, nil
}

func (sm *StateManager) CreateNeovimStateSnapshot(taskID string) (string, error) {
	snapshotID := uuid.New().String()

	states, err := sm.GetAllNeovimStates(taskID)
	if err != nil {
		return "", err
	}

	statesJSON, err := json.Marshal(states)
	if err != nil {
		return "", err
	}

	snapshot := map[string]interface{}{
		"id":         snapshotID,
		"task_id":    taskID,
		"states":     string(statesJSON),
		"created_at": time.Now().UTC().Format(time.RFC3339),
	}

	key := fmt.Sprintf("neovim:snapshot:%s:%s", taskID, snapshotID)
	err = sm.redisClient.HMSet(key, snapshot)
	if err != nil {
		return "", err
	}

	err = sm.redisClient.Expire(key, 24*60*60)
	if err != nil {
		return "", err
	}

	return snapshotID, nil
}

func (sm *StateManager) RollbackToNeovimSnapshot(taskID string, snapshotID string) error {
	key := fmt.Sprintf("neovim:snapshot:%s:%s", taskID, snapshotID)
	snapshot, err := sm.redisClient.HGetAll(key)
	if err != nil {
		return err
	}

	if len(snapshot) == 0 {
		return fmt.Errorf("snapshot not found: %s", snapshotID)
	}

	statesJSON, ok := snapshot["states"].(string)
	if !ok {
		return fmt.Errorf("invalid snapshot format")
	}

	var states map[string]map[string]interface{}
	err = json.Unmarshal([]byte(statesJSON), &states)
	if err != nil {
		return err
	}

	for stateTypeStr, stateMap := range states {
		for stateID, stateData := range stateMap {
			stateDataMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}

			stateType := StateType(stateTypeStr)
			err = sm.StoreNeovimState(stateType, stateID, stateDataMap, taskID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (sm *StateManager) SyncNeovimStateWithFrontend(taskID string) error {
	states, err := sm.GetAllNeovimStates(taskID)
	if err != nil {
		return err
	}

	frontendState := map[string]interface{}{
		"neovim":       states,
		"last_updated": time.Now().UTC().Format(time.RFC3339),
	}

	return sm.publishFrontendStateUpdate(taskID, frontendState)
}

func (sm *StateManager) PublishNeovimEvent(taskID string, eventType string, eventData map[string]interface{}) error {
	event := map[string]interface{}{
		"source":    "neovim",
		"type":      eventType,
		"data":      eventData,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return sm.eventStream.PublishEvent(taskID, event)
}

func (sm *StateManager) publishStateUpdate(taskID string, stateType string, stateID string, data map[string]interface{}) error {
	event := map[string]interface{}{
		"source":    "neovim",
		"type":      "state_update",
		"state_type": stateType,
		"state_id":   stateID,
		"data":       data,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	return sm.eventStream.PublishEvent(taskID, event)
}

func (sm *StateManager) publishFrontendStateUpdate(taskID string, frontendState map[string]interface{}) error {
	event := map[string]interface{}{
		"source":    "neovim",
		"type":      "frontend_state_update",
		"data":      frontendState,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return sm.eventStream.PublishEvent(taskID, event)
}

func (sm *StateManager) storeNeovimStateInSupabase(taskID string, stateType string, stateID string, data map[string]interface{}) error {
	err := sm.mainSupabase.StoreState(taskID, stateType, stateID, data)
	if err != nil {
		return err
	}

	previousState, _ := sm.mainSupabase.GetState(taskID, stateType, stateID)

	if previousState != nil {
		err = sm.rollbackSupabase.StoreState(taskID, stateType, stateID, previousState)
		if err != nil {
			log.Printf("Failed to store previous state in rollback Supabase: %v", err)
		}
	}

	return nil
}

func (sm *StateManager) retrieveNeovimStateFromSupabase(taskID string, stateType string, stateID string) (map[string]interface{}, bool) {
	data, exists := sm.readOnlySupabase.GetState(taskID, stateType, stateID)
	if exists {
		return data, true
	}

	data, exists = sm.mainSupabase.GetState(taskID, stateType, stateID)
	if exists {
		return data, true
	}

	return nil, false
}

func (sm *StateManager) rollbackNeovimState(taskID string, stateType string, stateID string) error {
	previousState, exists := sm.rollbackSupabase.GetState(taskID, stateType, stateID)
	if !exists {
		return fmt.Errorf("previous state not found in rollback Supabase")
	}

	err := sm.mainSupabase.StoreState(taskID, stateType, stateID, previousState)
	if err != nil {
		return err
	}

	err = sm.readOnlySupabase.StoreState(taskID, stateType, stateID, previousState)
	if err != nil {
		log.Printf("Failed to update read-only Supabase: %v", err)
	}

	return nil
}
