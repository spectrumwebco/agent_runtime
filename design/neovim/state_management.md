# Neovim State Management for Agent Runtime

## Overview

This document outlines the state management system for Neovim integration in the Agent Runtime. The system is designed to track and maintain Neovim states across the agent lifecycle and integrate with the existing shared frontend-backend state system.

## State Types

### 1. Neovim Execution State (TTL: 3600s)

Tracks the current execution state of Neovim operations:

```json
{
  "state_type": "neovim_execution_state",
  "state_id": "neovim-exec-{uuid}",
  "data": {
    "status": "active|inactive|transitioning",
    "current_mode": "normal|insert|visual|command",
    "current_file": "/path/to/file",
    "cursor_position": {
      "line": 10,
      "column": 5
    },
    "last_command": "w",
    "command_history": ["i", "w", "q"],
    "last_updated": "2025-04-17T04:35:00Z"
  }
}
```

### 2. Neovim Tool State (TTL: 7200s)

Maintains the state of Neovim tools and their configurations:

```json
{
  "state_type": "neovim_tool_state",
  "state_id": "neovim-tool-{uuid}",
  "data": {
    "active_plugins": [
      {
        "name": "telescope.nvim",
        "status": "loaded",
        "config": {}
      },
      {
        "name": "nvim-lspconfig",
        "status": "loaded",
        "config": {
          "active_servers": ["gopls", "pyright"]
        }
      }
    ],
    "terminal_integration": {
      "active_terminals": [
        {
          "id": "term-1",
          "command": "go build",
          "status": "running|completed|failed",
          "buffer_id": 3
        }
      ]
    },
    "last_updated": "2025-04-17T04:35:00Z"
  }
}
```

### 3. Neovim Context State (TTL: 86400s)

Handles contextual information for long-running Neovim processes:

```json
{
  "state_type": "neovim_context_state",
  "state_id": "neovim-context-{uuid}",
  "data": {
    "workspace": {
      "root": "/workspace/project",
      "open_files": [
        "/workspace/project/main.go",
        "/workspace/project/utils.go"
      ],
      "file_history": [
        "/workspace/project/config.go",
        "/workspace/project/main.go"
      ]
    },
    "session": {
      "id": "session-{uuid}",
      "created_at": "2025-04-17T04:00:00Z",
      "last_active": "2025-04-17T04:35:00Z"
    },
    "last_updated": "2025-04-17T04:35:00Z"
  }
}
```

### 4. Neovim Buffer State (TTL: 1800s)

Tracks the state of individual buffers in Neovim:

```json
{
  "state_type": "neovim_buffer_state",
  "state_id": "neovim-buffer-{buffer_id}",
  "data": {
    "buffer_id": 2,
    "file_path": "/workspace/project/main.go",
    "language": "go",
    "modified": true,
    "content_hash": "sha256:abc123...",
    "diagnostics": [
      {
        "line": 15,
        "column": 10,
        "severity": "error",
        "message": "undefined: foo"
      }
    ],
    "last_updated": "2025-04-17T04:35:00Z"
  }
}
```

### 5. Neovim UI State (TTL: 1800s)

Maintains the state of the Neovim UI:

```json
{
  "state_type": "neovim_ui_state",
  "state_id": "neovim-ui-{uuid}",
  "data": {
    "layout": {
      "windows": [
        {
          "id": 1,
          "buffer_id": 2,
          "position": {
            "row": 0,
            "col": 0
          },
          "size": {
            "width": 80,
            "height": 30
          }
        }
      ],
      "current_window": 1
    },
    "colorscheme": "tokyonight",
    "font": {
      "name": "JetBrainsMono Nerd Font",
      "size": 14
    },
    "last_updated": "2025-04-17T04:35:00Z"
  }
}
```

## State Operations

### Store Neovim State

```go
func (sm *StateManager) StoreNeovimState(stateType string, stateID string, data map[string]interface{}, taskID string) error {
    // Add timestamp
    data["last_updated"] = time.Now().UTC().Format(time.RFC3339)
    
    // Store in Redis with appropriate TTL
    ttl := getTTLForStateType(stateType)
    key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)
    
    // Store in Redis
    err := sm.redisClient.HMSet(key, data)
    if err != nil {
        return err
    }
    
    // Set expiration
    err = sm.redisClient.Expire(key, ttl)
    if err != nil {
        return err
    }
    
    // Publish state update to RocketMQ
    return sm.publishStateUpdate(taskID, stateType, stateID, data)
}
```

### Retrieve Neovim State

```go
func (sm *StateManager) RetrieveNeovimState(stateType string, stateID string, taskID string) (map[string]interface{}, bool) {
    key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)
    
    // Try to get from Redis first
    data, err := sm.redisClient.HGetAll(key)
    if err == nil && len(data) > 0 {
        return data, true
    }
    
    // If not in Redis, try to get from Supabase
    data, exists := sm.supabaseClient.GetState(taskID, stateType, stateID)
    if exists {
        // Store back in Redis for future access
        sm.redisClient.HMSet(key, data)
        ttl := getTTLForStateType(stateType)
        sm.redisClient.Expire(key, ttl)
        return data, true
    }
    
    return nil, false
}
```

### Update Neovim State

```go
func (sm *StateManager) UpdateNeovimState(stateType string, stateID string, updates map[string]interface{}, taskID string) error {
    key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)
    
    // Get current state
    currentState, exists := sm.RetrieveNeovimState(stateType, stateID, taskID)
    if !exists {
        return fmt.Errorf("state not found: %s", key)
    }
    
    // Update state
    for k, v := range updates {
        currentState[k] = v
    }
    
    // Update timestamp
    currentState["last_updated"] = time.Now().UTC().Format(time.RFC3339)
    
    // Store updated state
    return sm.StoreNeovimState(stateType, stateID, currentState, taskID)
}
```

### Merge Neovim States

```go
func (sm *StateManager) MergeNeovimStates(sourceStateType string, sourceStateID string, targetStateType string, targetStateID string, taskID string) error {
    // Get source state
    sourceState, exists := sm.RetrieveNeovimState(sourceStateType, sourceStateID, taskID)
    if !exists {
        return fmt.Errorf("source state not found: %s:%s", sourceStateType, sourceStateID)
    }
    
    // Get target state
    targetState, exists := sm.RetrieveNeovimState(targetStateType, targetStateID, taskID)
    if !exists {
        return fmt.Errorf("target state not found: %s:%s", targetStateType, targetStateID)
    }
    
    // Merge states
    for k, v := range sourceState {
        // Skip metadata fields
        if k == "last_updated" || k == "state_type" || k == "state_id" {
            continue
        }
        targetState[k] = v
    }
    
    // Update timestamp
    targetState["last_updated"] = time.Now().UTC().Format(time.RFC3339)
    
    // Store merged state
    return sm.StoreNeovimState(targetStateType, targetStateID, targetState, taskID)
}
```

### Expire Neovim State

```go
func (sm *StateManager) ExpireNeovimState(stateType string, stateID string, taskID string) error {
    key := fmt.Sprintf("neovim:%s:%s:%s", taskID, stateType, stateID)
    
    // Delete from Redis
    err := sm.redisClient.Del(key)
    if err != nil {
        return err
    }
    
    // Mark as expired in Supabase
    return sm.supabaseClient.ExpireState(taskID, stateType, stateID)
}
```

## Integration with Shared Frontend-Backend State

The Neovim state management system integrates with the shared frontend-backend state through:

1. **State Synchronization**: Neovim states are synchronized with the frontend through the shared state system.
2. **Event Propagation**: Neovim events are propagated to the frontend through the event stream.
3. **Command Execution**: Frontend commands are executed in Neovim through the shared state system.

### Frontend-Backend State Synchronization

```go
func (sm *StateManager) SyncNeovimStateWithFrontend(taskID string) error {
    // Get all Neovim states for the task
    states, err := sm.GetAllNeovimStates(taskID)
    if err != nil {
        return err
    }
    
    // Create frontend state update
    frontendState := map[string]interface{}{
        "neovim": states,
        "last_updated": time.Now().UTC().Format(time.RFC3339),
    }
    
    // Publish to frontend state topic
    return sm.publishFrontendStateUpdate(taskID, frontendState)
}
```

### Neovim Event Stream Integration

```go
func (sm *StateManager) PublishNeovimEvent(taskID string, eventType string, eventData map[string]interface{}) error {
    event := map[string]interface{}{
        "source": "neovim",
        "type": eventType,
        "data": eventData,
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }
    
    // Publish to event stream
    return sm.eventStream.PublishEvent(taskID, event)
}
```

## State Transition Management

When transitioning between Neovim and other tools, the state management system ensures state preservation:

1. **Pre-transition Snapshot**: Captures the current Neovim state before transitioning.
2. **State Maintenance**: Maintains the state during the transition.
3. **Post-transition Verification**: Ensures continuity after the transition.
4. **Rollback Capability**: Provides the ability to roll back to a previous state if the transition fails.

### Pre-transition Snapshot

```go
func (sm *StateManager) CreateNeovimStateSnapshot(taskID string) (string, error) {
    snapshotID := uuid.New().String()
    
    // Get all Neovim states for the task
    states, err := sm.GetAllNeovimStates(taskID)
    if err != nil {
        return "", err
    }
    
    // Create snapshot
    snapshot := map[string]interface{}{
        "id": snapshotID,
        "task_id": taskID,
        "states": states,
        "created_at": time.Now().UTC().Format(time.RFC3339),
    }
    
    // Store snapshot
    key := fmt.Sprintf("neovim:snapshot:%s:%s", taskID, snapshotID)
    err = sm.redisClient.HMSet(key, snapshot)
    if err != nil {
        return "", err
    }
    
    // Set expiration (24 hours)
    err = sm.redisClient.Expire(key, 24*60*60)
    if err != nil {
        return "", err
    }
    
    return snapshotID, nil
}
```

### Rollback to Snapshot

```go
func (sm *StateManager) RollbackToNeovimSnapshot(taskID string, snapshotID string) error {
    // Get snapshot
    key := fmt.Sprintf("neovim:snapshot:%s:%s", taskID, snapshotID)
    snapshot, err := sm.redisClient.HGetAll(key)
    if err != nil {
        return err
    }
    
    if len(snapshot) == 0 {
        return fmt.Errorf("snapshot not found: %s", snapshotID)
    }
    
    // Extract states from snapshot
    statesJSON, ok := snapshot["states"].(string)
    if !ok {
        return fmt.Errorf("invalid snapshot format")
    }
    
    var states map[string]map[string]interface{}
    err = json.Unmarshal([]byte(statesJSON), &states)
    if err != nil {
        return err
    }
    
    // Restore states
    for stateType, stateMap := range states {
        for stateID, stateData := range stateMap {
            stateDataMap, ok := stateData.(map[string]interface{})
            if !ok {
                continue
            }
            
            err = sm.StoreNeovimState(stateType, stateID, stateDataMap, taskID)
            if err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

## Three-Supabase Architecture Integration

The Neovim state management system integrates with the three-Supabase architecture:

1. **Main Supabase**: Primary database for active Neovim state storage.
2. **Read-only Supabase**: Replica for read operations to reduce load on the main instance.
3. **Rollback Supabase**: Maintains Neovim state data one step behind the main database to enable recovery in case of state confusion.

### Store State in Three-Supabase Architecture

```go
func (sm *StateManager) storeNeovimStateInSupabase(taskID string, stateType string, stateID string, data map[string]interface{}) error {
    // Store in main Supabase
    err := sm.mainSupabase.StoreState(taskID, stateType, stateID, data)
    if err != nil {
        return err
    }
    
    // Get previous state from main Supabase before update
    previousState, _ := sm.mainSupabase.GetState(taskID, stateType, stateID)
    
    // Store previous state in rollback Supabase
    if previousState != nil {
        err = sm.rollbackSupabase.StoreState(taskID, stateType, stateID, previousState)
        if err != nil {
            log.Printf("Failed to store previous state in rollback Supabase: %v", err)
        }
    }
    
    return nil
}
```

### Retrieve State from Three-Supabase Architecture

```go
func (sm *StateManager) retrieveNeovimStateFromSupabase(taskID string, stateType string, stateID string) (map[string]interface{}, bool) {
    // Try to get from read-only Supabase first
    data, exists := sm.readOnlySupabase.GetState(taskID, stateType, stateID)
    if exists {
        return data, true
    }
    
    // If not in read-only Supabase, try to get from main Supabase
    data, exists = sm.mainSupabase.GetState(taskID, stateType, stateID)
    if exists {
        return data, true
    }
    
    return nil, false
}
```

### Rollback State in Three-Supabase Architecture

```go
func (sm *StateManager) rollbackNeovimState(taskID string, stateType string, stateID string) error {
    // Get state from rollback Supabase
    previousState, exists := sm.rollbackSupabase.GetState(taskID, stateType, stateID)
    if !exists {
        return fmt.Errorf("previous state not found in rollback Supabase")
    }
    
    // Store previous state in main Supabase
    err := sm.mainSupabase.StoreState(taskID, stateType, stateID, previousState)
    if err != nil {
        return err
    }
    
    // Update read-only Supabase
    err = sm.readOnlySupabase.StoreState(taskID, stateType, stateID, previousState)
    if err != nil {
        log.Printf("Failed to update read-only Supabase: %v", err)
    }
    
    return nil
}
```

## Conclusion

This Neovim state management system provides a robust framework for tracking and maintaining Neovim states across the agent lifecycle. It integrates with the existing shared frontend-backend state system and supports the three-Supabase architecture for state persistence. The system enables the agent to maintain context across long-running operations, transition between different tool systems, and recover gracefully from errors.
