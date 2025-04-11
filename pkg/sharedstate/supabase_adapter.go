package sharedstate

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spectrumwebco/agent_runtime/internal/statemanager/supabase"
)

type SupabaseAdapter struct {
	supabaseStateManager *supabase.SupabaseStateManager
}

func NewSupabaseAdapter(supabaseStateManager *supabase.SupabaseStateManager) *SupabaseAdapter {
	return &SupabaseAdapter{
		supabaseStateManager: supabaseStateManager,
	}
}

func (a *SupabaseAdapter) GetState(stateType StateType, stateID string) (map[string]interface{}, error) {
	stateData, err := a.supabaseStateManager.GetState(supabase.StateType(stateType), stateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from Supabase: %w", err)
	}

	var state map[string]interface{}
	if err := json.Unmarshal(stateData, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state data: %w", err)
	}

	return state, nil
}

func (a *SupabaseAdapter) UpdateState(stateType StateType, stateID string, data map[string]interface{}) error {
	stateData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal state data: %w", err)
	}

	if err := a.supabaseStateManager.SaveState(supabase.StateType(stateType), stateID, stateData); err != nil {
		return fmt.Errorf("failed to save state to Supabase: %w", err)
	}

	return nil
}

func (a *SupabaseAdapter) RollbackState(stateType StateType, stateID string) error {
	if err := a.supabaseStateManager.RollbackState(supabase.StateType(stateType), stateID); err != nil {
		return fmt.Errorf("failed to rollback state in Supabase: %w", err)
	}

	log.Printf("Successfully rolled back %s state with ID %s", stateType, stateID)
	return nil
}

func (a *SupabaseAdapter) DeleteState(stateType StateType, stateID string) error {
	if err := a.supabaseStateManager.DeleteState(supabase.StateType(stateType), stateID); err != nil {
		return fmt.Errorf("failed to delete state from Supabase: %w", err)
	}

	log.Printf("Successfully deleted %s state with ID %s", stateType, stateID)
	return nil
}
