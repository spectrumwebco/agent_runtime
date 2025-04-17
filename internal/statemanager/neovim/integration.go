package neovim

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type FrontendStateManager interface {
	UpdateState(taskID string, stateType string, data map[string]interface{}) error
	GetState(taskID string, stateType string) (map[string]interface{}, bool)
}

type BackendStateManager interface {
	UpdateState(taskID string, stateType string, data map[string]interface{}) error
	GetState(taskID string, stateType string) (map[string]interface{}, bool)
}

type StateIntegration struct {
	stateManager       *StateManager
	frontendStateManager FrontendStateManager
	backendStateManager  BackendStateManager
	eventStream        EventStream
}

func NewStateIntegration(
	stateManager *StateManager,
	frontendStateManager FrontendStateManager,
	backendStateManager BackendStateManager,
	eventStream EventStream,
) *StateIntegration {
	return &StateIntegration{
		stateManager:       stateManager,
		frontendStateManager: frontendStateManager,
		backendStateManager:  backendStateManager,
		eventStream:        eventStream,
	}
}

func (si *StateIntegration) SyncNeovimStateWithFrontend(taskID string) error {
	states, err := si.stateManager.GetAllNeovimStates(taskID)
	if err != nil {
		return err
	}

	frontendState := map[string]interface{}{
		"neovim":       states,
		"last_updated": time.Now().UTC().Format(time.RFC3339),
	}

	return si.frontendStateManager.UpdateState(taskID, "neovim", frontendState)
}

func (si *StateIntegration) SyncFrontendStateWithNeovim(taskID string) error {
	frontendState, exists := si.frontendStateManager.GetState(taskID, "neovim")
	if !exists {
		return fmt.Errorf("frontend state not found for task: %s", taskID)
	}

	neovimStatesJSON, ok := frontendState["neovim"]
	if !ok {
		return fmt.Errorf("neovim states not found in frontend state")
	}

	neovimStatesMap, ok := neovimStatesJSON.(map[string]map[string]interface{})
	if !ok {
		neovimStatesStr, ok := neovimStatesJSON.(string)
		if !ok {
			return fmt.Errorf("invalid neovim states format")
		}

		var neovimStates map[string]map[string]interface{}
		err := json.Unmarshal([]byte(neovimStatesStr), &neovimStates)
		if err != nil {
			return err
		}

		neovimStatesMap = neovimStates
	}

	for stateTypeStr, stateMap := range neovimStatesMap {
		for stateID, stateData := range stateMap {
			stateDataMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}

			stateType := StateType(stateTypeStr)
			err := si.stateManager.StoreNeovimState(stateType, stateID, stateDataMap, taskID)
			if err != nil {
				log.Printf("Failed to store Neovim state: %v", err)
			}
		}
	}

	return nil
}

func (si *StateIntegration) SyncNeovimStateWithBackend(taskID string) error {
	states, err := si.stateManager.GetAllNeovimStates(taskID)
	if err != nil {
		return err
	}

	backendState := map[string]interface{}{
		"neovim":       states,
		"last_updated": time.Now().UTC().Format(time.RFC3339),
	}

	return si.backendStateManager.UpdateState(taskID, "neovim", backendState)
}

func (si *StateIntegration) SyncBackendStateWithNeovim(taskID string) error {
	backendState, exists := si.backendStateManager.GetState(taskID, "neovim")
	if !exists {
		return fmt.Errorf("backend state not found for task: %s", taskID)
	}

	neovimStatesJSON, ok := backendState["neovim"]
	if !ok {
		return fmt.Errorf("neovim states not found in backend state")
	}

	neovimStatesMap, ok := neovimStatesJSON.(map[string]map[string]interface{})
	if !ok {
		neovimStatesStr, ok := neovimStatesJSON.(string)
		if !ok {
			return fmt.Errorf("invalid neovim states format")
		}

		var neovimStates map[string]map[string]interface{}
		err := json.Unmarshal([]byte(neovimStatesStr), &neovimStates)
		if err != nil {
			return err
		}

		neovimStatesMap = neovimStates
	}

	for stateTypeStr, stateMap := range neovimStatesMap {
		for stateID, stateData := range stateMap {
			stateDataMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}

			stateType := StateType(stateTypeStr)
			err := si.stateManager.StoreNeovimState(stateType, stateID, stateDataMap, taskID)
			if err != nil {
				log.Printf("Failed to store Neovim state: %v", err)
			}
		}
	}

	return nil
}

func (si *StateIntegration) PublishNeovimEvent(taskID string, eventType string, eventData map[string]interface{}) error {
	event := map[string]interface{}{
		"source":    "neovim",
		"type":      eventType,
		"data":      eventData,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return si.eventStream.PublishEvent(taskID, event)
}

func (si *StateIntegration) HandleNeovimCommand(taskID string, command string, params map[string]interface{}) error {
	event := map[string]interface{}{
		"source":    "frontend",
		"type":      "neovim_command",
		"command":   command,
		"params":    params,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	return si.eventStream.PublishEvent(taskID, event)
}

func (si *StateIntegration) HandleNeovimEvent(taskID string, event map[string]interface{}) error {
	eventType, ok := event["type"].(string)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	switch eventType {
	case "buffer_update":
		return si.handleBufferUpdateEvent(taskID, event)
	case "mode_change":
		return si.handleModeChangeEvent(taskID, event)
	case "cursor_move":
		return si.handleCursorMoveEvent(taskID, event)
	case "command_executed":
		return si.handleCommandExecutedEvent(taskID, event)
	default:
		return si.eventStream.PublishEvent(taskID, event)
	}
}

func (si *StateIntegration) handleBufferUpdateEvent(taskID string, event map[string]interface{}) error {
	eventData, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}

	bufferID, ok := eventData["buffer_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid buffer ID")
	}

	bufferStateID := fmt.Sprintf("buffer-%d", int(bufferID))
	bufferState, exists := si.stateManager.RetrieveNeovimState(BufferState, bufferStateID, taskID)
	if !exists {
		bufferState = map[string]interface{}{
			"buffer_id": int(bufferID),
		}
	}

	for k, v := range eventData {
		bufferState[k] = v
	}

	return si.stateManager.StoreNeovimState(BufferState, bufferStateID, bufferState, taskID)
}

func (si *StateIntegration) handleModeChangeEvent(taskID string, event map[string]interface{}) error {
	eventData, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}

	mode, ok := eventData["mode"].(string)
	if !ok {
		return fmt.Errorf("invalid mode")
	}

	executionStateID := "execution"
	executionState, exists := si.stateManager.RetrieveNeovimState(ExecutionState, executionStateID, taskID)
	if !exists {
		executionState = map[string]interface{}{
			"status": string(Active),
		}
	}

	executionState["current_mode"] = mode

	return si.stateManager.StoreNeovimState(ExecutionState, executionStateID, executionState, taskID)
}

func (si *StateIntegration) handleCursorMoveEvent(taskID string, event map[string]interface{}) error {
	eventData, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}

	line, ok := eventData["line"].(float64)
	if !ok {
		return fmt.Errorf("invalid line")
	}

	column, ok := eventData["column"].(float64)
	if !ok {
		return fmt.Errorf("invalid column")
	}

	executionStateID := "execution"
	executionState, exists := si.stateManager.RetrieveNeovimState(ExecutionState, executionStateID, taskID)
	if !exists {
		executionState = map[string]interface{}{
			"status": string(Active),
		}
	}

	executionState["cursor_position"] = map[string]interface{}{
		"line":   int(line),
		"column": int(column),
	}

	return si.stateManager.StoreNeovimState(ExecutionState, executionStateID, executionState, taskID)
}

func (si *StateIntegration) handleCommandExecutedEvent(taskID string, event map[string]interface{}) error {
	eventData, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}

	command, ok := eventData["command"].(string)
	if !ok {
		return fmt.Errorf("invalid command")
	}

	executionStateID := "execution"
	executionState, exists := si.stateManager.RetrieveNeovimState(ExecutionState, executionStateID, taskID)
	if !exists {
		executionState = map[string]interface{}{
			"status": string(Active),
		}
	}

	executionState["last_command"] = command

	commandHistory, ok := executionState["command_history"].([]string)
	if !ok {
		commandHistory = []string{}
	}
	commandHistory = append(commandHistory, command)
	if len(commandHistory) > 10 {
		commandHistory = commandHistory[len(commandHistory)-10:]
	}
	executionState["command_history"] = commandHistory

	return si.stateManager.StoreNeovimState(ExecutionState, executionStateID, executionState, taskID)
}
