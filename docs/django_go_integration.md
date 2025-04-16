# Django-Go Integration

This document describes the integration between the Django backend and Go components in the Agent Runtime system.

## Overview

The integration enables seamless communication between the Django backend and Go components, allowing for:

1. **Shared State Management**: Synchronizing state between Django and Go components
2. **Event Streaming**: Propagating events between Django and Go components
3. **Tool Execution**: Executing Python tools from Go code and vice versa
4. **Model Access**: Accessing Django models from Go code
5. **Agent Communication**: Enabling communication between agents implemented in different languages

## Architecture

The integration is built on a gRPC-based bridge that provides a bidirectional communication channel between Django and Go components:

```
┌─────────────────┐                 ┌─────────────────┐
│                 │                 │                 │
│  Go Components  │◄───── gRPC ────►│ Django Backend  │
│                 │                 │                 │
└────────┬────────┘                 └────────┬────────┘
         │                                   │
         │                                   │
         ▼                                   ▼
┌─────────────────┐                 ┌─────────────────┐
│                 │                 │                 │
│  Shared State   │◄─── Redis ─────►│  Shared State   │
│                 │                 │                 │
└─────────────────┘                 └─────────────────┘
         │                                   │
         │                                   │
         ▼                                   ▼
┌─────────────────┐                 ┌─────────────────┐
│                 │                 │                 │
│  Event Stream   │◄── RocketMQ ───►│  Event Stream   │
│                 │                 │                 │
└─────────────────┘                 └─────────────────┘
```

## Components

### Django Bridge

The Django Bridge provides a Go interface for interacting with the Django backend:

- **DjangoBridge**: Main interface for interacting with Django
- **GRPCClient**: Low-level gRPC client for communicating with Django
- **AgentService**: gRPC service definition for agent-related operations

### Django gRPC Server

The Django gRPC server provides a Python interface for handling requests from Go components:

- **AgentService**: gRPC service implementation for agent-related operations
- **ModelService**: gRPC service implementation for model-related operations
- **CommandService**: gRPC service implementation for command-related operations

### Shared State

The shared state system enables state synchronization between Django and Go components:

- **Redis**: Used for storing shared state
- **StateProvider**: Interface for accessing shared state
- **StateIntegration**: Integration between shared state and LangGraph

### Event Stream

The event stream system enables event propagation between Django and Go components:

- **RocketMQ**: Used for event streaming
- **EventStream**: Interface for publishing and subscribing to events
- **EventIntegration**: Integration between event stream and LangGraph

## Usage

### Creating a Django Agent from Go

```go
// Create a new Django bridge
bridge, err := djangobridge.NewDjangoBridge("localhost:50051", eventStream)
if err != nil {
    // Handle error
}
defer bridge.Close()

// Create a Django agent
agent, err := bridge.CreateDjangoAgent(ctx, "agent-id", "Agent Name", "agent-role")
if err != nil {
    // Handle error
}

// Use the agent
// ...
```

### Executing Python Code from Go

```go
// Create a new Django bridge
bridge, err := djangobridge.NewDjangoBridge("localhost:50051", eventStream)
if err != nil {
    // Handle error
}
defer bridge.Close()

// Execute Python code
result, err := bridge.ExecutePythonCode(ctx, `
import numpy as np
print("Hello from Python!")
result = np.array([1, 2, 3]).sum()
`, 5*time.Second)
if err != nil {
    // Handle error
}

// Use the result
fmt.Println(result)
```

### Accessing Django Models from Go

```go
// Create a new Django bridge
bridge, err := djangobridge.NewDjangoBridge("localhost:50051", eventStream)
if err != nil {
    // Handle error
}
defer bridge.Close()

// Query Django models
users, err := bridge.QueryDjangoModel(ctx, "User", map[string]interface{}{
    "is_active": true,
})
if err != nil {
    // Handle error
}

// Use the users
for _, user := range users {
    fmt.Println(user["username"])
}
```

## Database Integration

The integration supports multiple databases:

1. **Supabase**: PostgreSQL database with authentication and functions capabilities
2. **DragonflyDB**: Redis replacement with memcached functionality
3. **RAGflow**: Vector database with deep search and understanding capabilities
4. **RocketMQ**: Go-based messaging system for state communication

Each database is configured to work with Kubernetes service discovery in production environments.

## Multi-Agent Communication

The integration enables communication between agents implemented in different languages:

1. **Go Agents**: Agents implemented in Go using LangGraph-Go
2. **Python Agents**: Agents implemented in Python using Django
3. **TypeScript Agents**: Agents implemented in TypeScript using React

Agents can communicate through:

1. **Direct Communication**: Agents can directly call other agents
2. **Event-Based Communication**: Agents can publish events that other agents can subscribe to
3. **Tool-Based Communication**: Agents can use tools provided by other agents
4. **Graph-Based Communication**: Agents can communicate through the graph execution system

## Future Work

Future improvements to the integration include:

1. **Better Performance**: Optimizing the performance of the gRPC bridge
2. **More Features**: Adding more features to the bridge
3. **Better Error Handling**: Improving error handling and recovery
4. **Better Testing**: Adding more tests for the bridge
5. **Better Documentation**: Improving documentation for the bridge
