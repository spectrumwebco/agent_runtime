# Event Stream

This directory contains the implementation of the Event Stream component for the Agent Runtime.

## Overview

The Event Stream is a central component responsible for:

1.  **Distributing Events**: Relaying events (messages, actions, observations, etc.) between different parts of the agent system.
2.  **Context Caching**: Utilizing DragonflyDB to cache application context (e.g., file contents, state information) for faster access.
3.  **Cache Invalidation**: Implementing a robust system to invalidate cached data when the underlying source changes, ensuring the agent always works with up-to-date information.

## Components

-   `stream.go`: The main implementation of the `Stream` struct, handling event publishing, subscription, and interaction with DragonflyDB.
-   `pkg/eventstream/event.go`: Defines the `Event` struct and related types (`EventType`, `EventSource`).

## DragonflyDB Integration

-   The stream connects to a DragonflyDB instance (compatible with Redis API).
-   Events themselves might be stored temporarily.
-   Application context (e.g., file contents, API responses) is cached using standard Redis commands (`SET`, `GET`, `DEL`).
-   Cache keys follow a convention (e.g., `context:file:/path/to/file`, `context:state:task_id`).

## Cache Invalidation

-   The `handleCacheInvalidation` function in `stream.go` contains the core logic.
-   It listens to specific events (e.g., file changes from the sandbox, state updates from RocketMQ) and deletes corresponding cache keys from DragonflyDB.
-   The goal is to ensure that calls to `GetAppContext` return fresh data if the source has changed since the last cache write.

## Data Sources

The Event Stream needs to be connected to various data sources to receive events that trigger cache invalidation:

-   **CI/CD Pipelines**: Events indicating code changes, build status, etc.
-   **Local Production Server**: Events related to the agent's own operations or local file changes.
-   **Sandbox**: Events about file modifications, command outputs within the sandbox.
-   **Kubernetes**: Events related to cluster state changes (e.g., pod status, deployment updates).
-   **RocketMQ**: State update messages.

Integration with these sources will involve adapters or listeners that publish relevant events to this Event Stream.
