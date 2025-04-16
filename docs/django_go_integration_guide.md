# Django-Go Integration Guide

This guide explains how to integrate the Go codebase with Django applications in the Agent Runtime system.

## Overview

The integration between Django and Go is implemented using gRPC, which provides a high-performance, language-agnostic communication protocol. This integration allows Django applications to:

1. Execute tasks in the Go runtime
2. Execute agent tasks in the multi-agent system
3. Access and modify shared state
4. Subscribe to and publish events
5. Utilize the LangGraph and LangChain integrations

## Architecture

The integration consists of the following components:

1. **Go Components**:
   - `pkg/djangobridge/integration.go`: Main integration module
   - `pkg/djangobridge/agent_service.proto`: gRPC service definition
   - `pkg/djangobridge/bridge.go`: Bridge implementation for Go-Django communication

2. **Python Components**:
   - `backend/apps/python_agent/go_integration.py`: Python client for Go runtime
   - `backend/apps/python_agent/grpc_client/`: Generated gRPC client code
   - `backend/apps/python_agent/django_go_example.py`: Example Django views

3. **Shared State and Event Stream**:
   - Maintained through the Go runtime
   - Accessible from both Go and Django

## Setup

### Prerequisites

1. Go 1.24.1 or later
2. Python 3.8 or later
3. Django 4.0 or later
4. gRPC tools for Python

### Installation

1. Install the required Python packages:

```bash
pip install grpcio grpcio-tools django
```

2. Generate the gRPC client code:

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. backend/apps/python_agent/grpc_client/agent_service.proto
```

3. Configure Django settings:

```python
# In settings.py
GO_RUNTIME = {
    'host': 'localhost',
    'port': 50051,
    'timeout': 30,
    'reconnection_attempts': 5,
    'reconnection_delay': 5,
}
```

## Usage

### Initializing the Go Runtime Integration

```python
from backend.apps.python_agent.go_integration import get_go_runtime_integration

# Get the singleton instance
go_runtime = get_go_runtime_integration()

# Connect to the Go runtime
go_runtime.connect()
```

### Executing Tasks

```python
# Execute a simple task
result = go_runtime.execute_task(
    task_type="example_task",
    input_data={"param1": "value1", "param2": "value2"},
    agent_id="agent1",
    description="Example task execution"
)

# Execute a task using a specific agent
result = go_runtime.execute_agent_task(
    agent_id="frontend_agent",
    task_type="generate_component",
    input_data={"component_type": "button", "props": {"label": "Click me"}},
    description="Generate a React button component"
)
```

### Managing State

```python
# Get a value from the shared state
value = go_runtime.get_state("user_preferences")

# Set a value in the shared state
success = go_runtime.set_state("user_preferences", {"theme": "dark", "language": "en"})

# Delete a value from the shared state
success = go_runtime.delete_state("temporary_data")
```

### Working with Events

```python
# Subscribe to events
def handle_user_action(event_data):
    print(f"User action: {event_data}")

subscription_id = go_runtime.subscribe_to_events("user_action", handle_user_action)

# Publish an event
success = go_runtime.publish_event(
    event_type="system_notification",
    data={"message": "System update completed", "status": "success"},
    source="django_app",
    metadata={"importance": "high"}
)

# Unsubscribe from events
go_runtime.unsubscribe_from_events(subscription_id)
```

### Django Integration

See `backend/apps/python_agent/django_go_example.py` for examples of Django views that integrate with the Go runtime.

## Multi-Agent System Integration

The Django-Go integration provides access to the multi-agent system implemented in Go. This system includes:

1. **Frontend Agent**: Specialized in UI/UX tasks
2. **App Builder Agent**: Focused on application architecture
3. **Codegen Agent**: Specialized in code generation
4. **Software Engineering Agent**: General-purpose software engineering

### Using the Multi-Agent System

```python
# Execute a task with the Frontend Agent
result = go_runtime.execute_agent_task(
    agent_id="frontend_agent",
    task_type="design_component",
    input_data={"component_type": "form", "fields": ["name", "email", "message"]},
    description="Design a contact form component"
)

# Execute a task with the Codegen Agent
result = go_runtime.execute_agent_task(
    agent_id="codegen_agent",
    task_type="generate_code",
    input_data={"language": "python", "task": "Implement a function to calculate Fibonacci numbers"},
    description="Generate Python code for Fibonacci calculation"
)
```

## LangGraph and LangChain Integration

The Django-Go integration provides access to LangGraph and LangChain capabilities implemented in Go:

```python
# Execute a LangGraph task
result = go_runtime.execute_task(
    task_type="langraph_task",
    input_data={"graph_id": "reasoning_graph", "input": "How do I implement a binary search tree?"},
    description="Execute a reasoning task using LangGraph"
)

# Execute a LangChain task
result = go_runtime.execute_task(
    task_type="langchain_task",
    input_data={"chain_id": "research_chain", "query": "Latest advancements in quantum computing"},
    description="Execute a research task using LangChain"
)
```

## LangSmith Integration

The Django-Go integration provides access to the self-hosted LangSmith deployment for tracing and monitoring:

```python
# Execute a task with LangSmith tracing
result = go_runtime.execute_task(
    task_type="traced_task",
    input_data={"query": "Explain the concept of recursion"},
    metadata={"trace": True, "project": "recursion_explanation"},
    description="Generate an explanation of recursion with LangSmith tracing"
)
```

## Troubleshooting

### Connection Issues

If you encounter connection issues:

1. Ensure the Go runtime server is running
2. Check the host and port configuration
3. Verify network connectivity between Django and Go services

### Error Handling

The Go integration client includes comprehensive error handling:

```python
try:
    result = go_runtime.execute_task(...)
except ConnectionError as e:
    print(f"Connection error: {e}")
except TimeoutError as e:
    print(f"Timeout error: {e}")
except Exception as e:
    print(f"Unexpected error: {e}")
```

## Best Practices

1. **Connection Management**: Connect to the Go runtime during Django application startup and disconnect during shutdown
2. **Error Handling**: Implement proper error handling for all Go runtime interactions
3. **State Management**: Use the shared state for data that needs to be accessible from both Go and Django
4. **Event-Driven Communication**: Use events for asynchronous communication between Go and Django
5. **Task Execution**: Use the appropriate task execution method based on the task requirements

## Conclusion

The Django-Go integration provides a powerful way to combine the strengths of both languages and frameworks. Django provides a robust web framework with excellent ORM capabilities, while Go provides high-performance concurrent processing and multi-agent system capabilities.
