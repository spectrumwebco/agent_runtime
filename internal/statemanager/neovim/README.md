# Neovim State Management

This package implements a state management system for Neovim integration in the Agent Runtime. It provides a robust framework for tracking and maintaining Neovim states across the agent lifecycle and integrates with the existing shared frontend-backend state system.

## Overview

The Neovim state management system is designed to:

1. Track and maintain Neovim states across the agent lifecycle
2. Integrate with the existing shared frontend-backend state system
3. Support the three-Supabase architecture for state persistence
4. Enable the agent to maintain context across long-running operations
5. Provide snapshot and rollback capabilities for state transitions

## State Types

The system defines five state types:

1. **Execution State (TTL: 3600s)**: Tracks the current execution state of Neovim operations
2. **Tool State (TTL: 7200s)**: Maintains the state of Neovim tools and their configurations
3. **Context State (TTL: 86400s)**: Handles contextual information for long-running Neovim processes
4. **Buffer State (TTL: 1800s)**: Tracks the state of individual buffers in Neovim
5. **UI State (TTL: 1800s)**: Maintains the state of the Neovim UI

## Usage

### Creating a State Manager

```go
import (
    "github.com/spectrumwebco/agent_runtime/internal/statemanager/neovim"
)

// Create a new state manager
stateManager := neovim.NewStateManager(
    redisClient,
    supabaseClient,
    mainSupabase,
    readOnlySupabase,
    rollbackSupabase,
    eventStream,
)
```

### Storing a Neovim State

```go
// Create a new execution state
executionState := map[string]interface{}{
    "status": string(neovim.Active),
    "current_mode": string(neovim.NormalMode),
    "current_file": "/workspace/project/main.go",
    "cursor_position": map[string]interface{}{
        "line": 10,
        "column": 5,
    },
    "last_command": "w",
    "command_history": []string{"i", "w", "q"},
}

// Store the execution state
err := stateManager.StoreNeovimState(neovim.ExecutionState, "execution", executionState, taskID)
if err != nil {
    log.Printf("Failed to store execution state: %v", err)
}
```

### Retrieving a Neovim State

```go
// Retrieve the execution state
executionState, exists := stateManager.RetrieveNeovimState(neovim.ExecutionState, "execution", taskID)
if !exists {
    log.Printf("Execution state not found")
    return
}

// Access state properties
status := executionState["status"].(string)
currentMode := executionState["current_mode"].(string)
```

### Updating a Neovim State

```go
// Update the execution state
updates := map[string]interface{}{
    "current_mode": string(neovim.InsertMode),
    "cursor_position": map[string]interface{}{
        "line": 15,
        "column": 10,
    },
}

err := stateManager.UpdateNeovimState(neovim.ExecutionState, "execution", updates, taskID)
if err != nil {
    log.Printf("Failed to update execution state: %v", err)
}
```

### Creating a State Snapshot

```go
// Create a snapshot of all Neovim states
snapshotID, err := stateManager.CreateNeovimStateSnapshot(taskID)
if err != nil {
    log.Printf("Failed to create snapshot: %v", err)
    return
}

// Store the snapshot ID for later use
log.Printf("Created snapshot: %s", snapshotID)
```

### Rolling Back to a Snapshot

```go
// Roll back to a previous snapshot
err := stateManager.RollbackToNeovimSnapshot(taskID, snapshotID)
if err != nil {
    log.Printf("Failed to roll back to snapshot: %v", err)
    return
}

log.Printf("Successfully rolled back to snapshot: %s", snapshotID)
```

## Integration with Shared Frontend-Backend State

The Neovim state management system integrates with the shared frontend-backend state through the `StateIntegration` struct:

```go
// Create a new state integration
stateIntegration := neovim.NewStateIntegration(
    stateManager,
    frontendStateManager,
    backendStateManager,
    eventStream,
)

// Synchronize Neovim states with the frontend
err := stateIntegration.SyncNeovimStateWithFrontend(taskID)
if err != nil {
    log.Printf("Failed to sync Neovim state with frontend: %v", err)
}

// Synchronize frontend state with Neovim
err = stateIntegration.SyncFrontendStateWithNeovim(taskID)
if err != nil {
    log.Printf("Failed to sync frontend state with Neovim: %v", err)
}
```

## Event Handling

The state integration system can handle Neovim events and commands:

```go
// Handle a Neovim event
event := map[string]interface{}{
    "source": "neovim",
    "type": "buffer_update",
    "data": map[string]interface{}{
        "buffer_id": 2,
        "file_path": "/workspace/project/main.go",
        "modified": true,
    },
    "timestamp": time.Now().UTC().Format(time.RFC3339),
}

err := stateIntegration.HandleNeovimEvent(taskID, event)
if err != nil {
    log.Printf("Failed to handle Neovim event: %v", err)
}

// Handle a Neovim command from the frontend
err = stateIntegration.HandleNeovimCommand(taskID, "write", map[string]interface{}{
    "file_path": "/workspace/project/main.go",
})
if err != nil {
    log.Printf("Failed to handle Neovim command: %v", err)
}
```

## Three-Supabase Architecture Integration

The Neovim state management system integrates with the three-Supabase architecture:

1. **Main Supabase**: Primary database for active Neovim state storage
2. **Read-only Supabase**: Replica for read operations to reduce load on the main instance
3. **Rollback Supabase**: Maintains Neovim state data one step behind the main database to enable recovery in case of state confusion

This architecture provides a robust framework for state persistence and recovery.

## Conclusion

The Neovim state management system provides a comprehensive framework for tracking and maintaining Neovim states across the agent lifecycle. It integrates with the existing shared frontend-backend state system and supports the three-Supabase architecture for state persistence. The system enables the agent to maintain context across long-running operations, transition between different tool systems, and recover gracefully from errors.
