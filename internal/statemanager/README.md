# State Management System

This directory contains the implementation of the State Management System for the Agent Runtime.

## Overview

The State Management System is responsible for:

1. **State Tracking**: Maintaining the state of tasks, agents, and their lifecycle.
2. **State Distribution**: Using RocketMQ to distribute state updates across the system.
3. **State Persistence**: Using Supabase (PostgreSQL) for persistent storage of states.
4. **LIFO Task Processing**: Supporting a Last-In-First-Out task queue mechanism.

## Components

### RocketMQ Integration

The `rocketmq` subdirectory contains the Go client implementation for interacting with Apache RocketMQ:

- `producer.go`: Implements `StateProducer` for publishing state updates.
- `consumer.go`: Implements `StateConsumer` for consuming state updates.
- `README.md`: Documentation for the RocketMQ integration.

### Supabase Integration

The `supabase` subdirectory contains the client implementation for interacting with Supabase:

- `client.go`: Implements `SupabaseClient` for basic Supabase operations.
- `state_manager.go`: Implements `SupabaseStateManager` for state persistence with three-Supabase architecture.

### State Manager

The `manager.go` file implements the main `StateManager` that coordinates between RocketMQ and Supabase:

- Publishes state updates to RocketMQ.
- Consumes state updates from RocketMQ.
- Persists states to Supabase.
- Implements LIFO task processing.

## Three-Supabase Architecture

The state management system uses a three-Supabase architecture:

1. **Main Supabase**: Primary database for active state storage.
2. **Read-only Supabase**: Replica for read operations to reduce load on main instance.
3. **Rollback Supabase**: Maintains state data one step behind the main database to enable recovery in case of state confusion.

## LIFO Task Queue

The LIFO (Last-In-First-Out) task queue mechanism allows the agent to address the most recent task first, then move backwards to complete the original task. This is implemented in the `StateManager.StartLifoTaskProcessor` method.

## Kubernetes and Kata Container Integration

The state management system integrates with Kubernetes and Kata Containers through:

- Kubernetes lifecycle hooks that publish state updates to RocketMQ.
- Kata Container lifecycle events that are captured and published to RocketMQ.
- DragonflyDB for caching state data and providing fast access.

## Configuration

Configuration for the state management system is defined in the `internal/config/config.go` file:

- `RocketMQConfig`: Configuration for RocketMQ connection and topics.
- `SupabaseConfig`: Configuration for Supabase connection.

## Usage

To start the state management service:

```bash
./agent_runtime state
```

To use the state management system in your code:

```go
import "github.com/spectrumwebco/agent_runtime/internal/statemanager"

// Initialize state manager
stateManager, err := statemanager.NewStateManager(cfg)
if err != nil {
    log.Fatalf("Failed to initialize state manager: %v", err)
}
defer stateManager.Close()

// Update state
err = stateManager.UpdateState("task-123", "task", taskData)
if err != nil {
    log.Printf("Failed to update state: %v", err)
}

// Get state
stateData, exists := stateManager.GetState("task-123")
if exists {
    // Use state data
}
```
