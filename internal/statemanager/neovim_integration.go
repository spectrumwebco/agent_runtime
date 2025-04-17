package statemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/statemanager/neovim"
)

type NeovimStateManager struct {
	stateManager *StateManager
	neovimManager *neovim.StateManager
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewNeovimStateManager(stateManager *StateManager, neovimManager *neovim.StateManager) *NeovimStateManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &NeovimStateManager{
		stateManager: stateManager,
		neovimManager: neovimManager,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (nsm *NeovimStateManager) StoreNeovimState(stateType neovim.StateType, stateID string, stateData map[string]interface{}, taskID string) error {
	err := nsm.neovimManager.StoreNeovimState(stateType, stateID, stateData, taskID)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("failed to marshal Neovim state data: %w", err)
	}

	stateKey := fmt.Sprintf("neovim:%s:%s", stateType, stateID)
	err = nsm.stateManager.UpdateState(stateKey, string(stateType), jsonData)
	if err != nil {
		return fmt.Errorf("failed to update Neovim state in main state manager: %w", err)
	}

	log.Printf("Neovim state stored: Type=%s, ID=%s, TaskID=%s", stateType, stateID, taskID)
	return nil
}

func (nsm *NeovimStateManager) RetrieveNeovimState(stateType neovim.StateType, stateID string, taskID string) (map[string]interface{}, bool) {
	stateData, exists := nsm.neovimManager.RetrieveNeovimState(stateType, stateID, taskID)
	if exists {
		return stateData, true
	}

	stateKey := fmt.Sprintf("neovim:%s:%s", stateType, stateID)
	jsonData, exists := nsm.stateManager.GetState(stateKey)
	if !exists {
		return nil, false
	}

	var stateData map[string]interface{}
	err := json.Unmarshal(jsonData, &stateData)
	if err != nil {
		log.Printf("Failed to unmarshal Neovim state data: %v", err)
		return nil, false
	}

	err = nsm.neovimManager.StoreNeovimState(stateType, stateID, stateData, taskID)
	if err != nil {
		log.Printf("Failed to store retrieved Neovim state in Neovim state manager: %v", err)
	}

	return stateData, true
}

func (nsm *NeovimStateManager) UpdateNeovimState(stateType neovim.StateType, stateID string, updates map[string]interface{}, taskID string) error {
	err := nsm.neovimManager.UpdateNeovimState(stateType, stateID, updates, taskID)
	if err != nil {
		return err
	}

	stateData, exists := nsm.neovimManager.RetrieveNeovimState(stateType, stateID, taskID)
	if !exists {
		return fmt.Errorf("failed to retrieve updated Neovim state")
	}

	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated Neovim state data: %w", err)
	}

	stateKey := fmt.Sprintf("neovim:%s:%s", stateType, stateID)
	err = nsm.stateManager.UpdateState(stateKey, string(stateType), jsonData)
	if err != nil {
		return fmt.Errorf("failed to update Neovim state in main state manager: %w", err)
	}

	log.Printf("Neovim state updated: Type=%s, ID=%s, TaskID=%s", stateType, stateID, taskID)
	return nil
}

func (nsm *NeovimStateManager) MergeNeovimStates(stateType neovim.StateType, targetStateID string, sourceStateIDs []string, taskID string) error {
	err := nsm.neovimManager.MergeNeovimStates(stateType, targetStateID, sourceStateIDs, taskID)
	if err != nil {
		return err
	}

	stateData, exists := nsm.neovimManager.RetrieveNeovimState(stateType, targetStateID, taskID)
	if !exists {
		return fmt.Errorf("failed to retrieve merged Neovim state")
	}

	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("failed to marshal merged Neovim state data: %w", err)
	}

	stateKey := fmt.Sprintf("neovim:%s:%s", stateType, targetStateID)
	err = nsm.stateManager.UpdateState(stateKey, string(stateType), jsonData)
	if err != nil {
		return fmt.Errorf("failed to update merged Neovim state in main state manager: %w", err)
	}

	log.Printf("Neovim states merged: Type=%s, TargetID=%s, TaskID=%s", stateType, targetStateID, taskID)
	return nil
}

func (nsm *NeovimStateManager) ExpireNeovimState(stateType neovim.StateType, stateID string, taskID string) error {
	err := nsm.neovimManager.ExpireNeovimState(stateType, stateID, taskID)
	if err != nil {
		return err
	}

	log.Printf("Neovim state expired: Type=%s, ID=%s, TaskID=%s", stateType, stateID, taskID)
	return nil
}

func (nsm *NeovimStateManager) CreateNeovimStateSnapshot(taskID string) (string, error) {
	snapshotID, err := nsm.neovimManager.CreateNeovimStateSnapshot(taskID)
	if err != nil {
		return "", err
	}

	states, err := nsm.neovimManager.GetAllNeovimStates(taskID)
	if err != nil {
		return "", fmt.Errorf("failed to get all Neovim states for snapshot: %w", err)
	}

	jsonData, err := json.Marshal(states)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Neovim states for snapshot: %w", err)
	}

	snapshotKey := fmt.Sprintf("neovim:snapshot:%s", snapshotID)
	err = nsm.stateManager.UpdateState(snapshotKey, "neovim_snapshot", jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to store Neovim snapshot in main state manager: %w", err)
	}

	log.Printf("Neovim state snapshot created: ID=%s, TaskID=%s", snapshotID, taskID)
	return snapshotID, nil
}

func (nsm *NeovimStateManager) RollbackToNeovimSnapshot(taskID string, snapshotID string) error {
	err := nsm.neovimManager.RollbackToNeovimSnapshot(taskID, snapshotID)
	if err != nil {
		return err
	}

	states, err := nsm.neovimManager.GetAllNeovimStates(taskID)
	if err != nil {
		return fmt.Errorf("failed to get all Neovim states after rollback: %w", err)
	}

	for stateTypeStr, stateMap := range states {
		for stateID, stateData := range stateMap {
			stateDataMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}

			jsonData, err := json.Marshal(stateDataMap)
			if err != nil {
				log.Printf("Failed to marshal Neovim state data after rollback: %v", err)
				continue
			}

			stateKey := fmt.Sprintf("neovim:%s:%s", stateTypeStr, stateID)
			err = nsm.stateManager.UpdateState(stateKey, stateTypeStr, jsonData)
			if err != nil {
				log.Printf("Failed to update Neovim state in main state manager after rollback: %v", err)
			}
		}
	}

	log.Printf("Rolled back to Neovim state snapshot: ID=%s, TaskID=%s", snapshotID, taskID)
	return nil
}

func (nsm *NeovimStateManager) SyncNeovimStateWithFrontend(taskID string) error {
	states, err := nsm.neovimManager.GetAllNeovimStates(taskID)
	if err != nil {
		return fmt.Errorf("failed to get all Neovim states for frontend sync: %w", err)
	}

	jsonData, err := json.Marshal(states)
	if err != nil {
		return fmt.Errorf("failed to marshal Neovim states for frontend sync: %w", err)
	}

	syncKey := fmt.Sprintf("neovim:frontend_sync:%s", taskID)
	err = nsm.stateManager.UpdateState(syncKey, "neovim_frontend_sync", jsonData)
	if err != nil {
		return fmt.Errorf("failed to store Neovim states for frontend sync: %w", err)
	}

	log.Printf("Neovim states synchronized with frontend: TaskID=%s", taskID)
	return nil
}

func (nsm *NeovimStateManager) SyncNeovimStateWithBackend(taskID string) error {
	states, err := nsm.neovimManager.GetAllNeovimStates(taskID)
	if err != nil {
		return fmt.Errorf("failed to get all Neovim states for backend sync: %w", err)
	}

	jsonData, err := json.Marshal(states)
	if err != nil {
		return fmt.Errorf("failed to marshal Neovim states for backend sync: %w", err)
	}

	syncKey := fmt.Sprintf("neovim:backend_sync:%s", taskID)
	err = nsm.stateManager.UpdateState(syncKey, "neovim_backend_sync", jsonData)
	if err != nil {
		return fmt.Errorf("failed to store Neovim states for backend sync: %w", err)
	}

	log.Printf("Neovim states synchronized with backend: TaskID=%s", taskID)
	return nil
}

func (nsm *NeovimStateManager) Close() error {
	nsm.cancel()
	return nil
}
