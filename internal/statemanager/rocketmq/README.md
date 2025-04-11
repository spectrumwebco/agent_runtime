# RocketMQ State Management Integration

This directory contains the Go client implementation for interacting with Apache RocketMQ, specifically for publishing and consuming agent state updates.

## Overview

RocketMQ is used as a message queue to decouple state producers (e.g., the agent loop, Kubernetes event listeners) from state consumers (e.g., monitoring systems, UI updates, potentially other agents). This allows for asynchronous communication and better scalability.

## Components

-   `producer.go`: Implements `StateProducer` for publishing state update messages to a specific RocketMQ topic.
-   `consumer.go`: Implements `StateConsumer` for subscribing to a RocketMQ topic and processing received state update messages using a provided handler function.
-   `README.md`: This file.

## Configuration

Both the producer and consumer require configuration, typically provided via environment variables or a configuration file (`internal/config`):

-   `ROCKETMQ_NAMESERVER_ADDRS`: Comma-separated list of RocketMQ NameServer addresses (e.g., `127.0.0.1:9876`).
-   `ROCKETMQ_STATE_TOPIC`: The topic name used for state updates (e.g., `agent_state_updates`).
-   `ROCKETMQ_PRODUCER_GROUP`: Group name for the producer instance.
-   `ROCKETMQ_CONSUMER_GROUP`: Group name for the consumer instance (different consumers processing the same state should ideally be in the same group for load balancing/failover).

## Usage

1.  **Initialize Producer/Consumer**: Create instances of `StateProducer` and/or `StateConsumer` during application startup, providing the necessary configuration and a handler function for the consumer.
2.  **Publish State**: Call `StateProducer.PublishStateUpdate` whenever a significant state change occurs (e.g., agent task status change, sandbox environment update). Pass the relevant state data (e.g., marshalled JSON) and optional keys.
3.  **Consume State**: The `StateConsumer` will automatically receive messages and invoke the provided handler function. Implement the handler logic to process the state update (e.g., update a database, notify the event stream).
4.  **Shutdown**: Ensure `Close()` is called on producer/consumer instances during graceful application shutdown.

## Integration Points

-   **Agent Loop (`internal/agent/loop.go`)**: Publish state updates at key points in the agent's lifecycle (task start/end, action execution, observation received).
-   **Event Stream (`internal/eventstream/stream.go`)**: Potentially consume state updates from RocketMQ to trigger cache invalidations or publish internal events.
-   **Kubernetes/Kata Lifecycle Hooks**: Implement logic within Kubernetes controllers or Kata container setup/teardown scripts to publish relevant lifecycle events via the producer. This might involve a small sidecar or direct integration if feasible.
-   **Main Application (`cmd/agent_runtime/main.go`)**: Initialize and manage the lifecycle of the producer/consumer instances.
