# React Server Components (RSC) Integration

This document describes the integration of React Server Components (RSC) with the agent runtime system, enabling Agentic Generative UI based on agent actions and tool usage.

## Overview

The RSC integration allows the agent runtime to generate dynamic UI components based on agent actions, tool usage, and shared state. This creates a more visually appealing and interactive agent chat experience, similar to CoPilotKit CoAgents but built specifically for our system.

## Architecture

The RSC integration consists of the following components:

1. **RSC Adapter**: Generates React Server Components based on agent actions and tool usage.
2. **Shared State Integration**: Stores and manages component state in the shared state system.
3. **Event Stream Integration**: Publishes and subscribes to events related to component generation and updates.
4. **Django Integration**: Provides REST API endpoints for component generation and retrieval.
5. **Frontend Components**: Renders the generated components in the UI.

## Component Types

The RSC integration supports the following component types:

- **Button**: Interactive button component
- **Card**: Container component for displaying content
- **CodeBlock**: Component for displaying code with syntax highlighting
- **ProgressBar**: Component for displaying progress
- **Terminal**: Component for displaying terminal output
- **Markdown**: Component for rendering markdown content
- **AgentAction**: Component for displaying agent actions
- **ToolOutput**: Component for displaying tool usage results

## Integration with Shared State

The RSC integration uses the shared state system to store and manage component state. Components are stored in the shared state with the following state types:

- `StateTypeComponent`: Stores component data
- `StateTypeUI`: Stores UI-specific data
- `StateTypeAction`: Stores agent action data
- `StateTypeTool`: Stores tool usage data

## Integration with Event Stream

The RSC integration uses the event stream to publish and subscribe to events related to component generation and updates. The following event types are used:

- `EventTypeComponentGenerated`: Published when a new component is generated
- `EventTypeComponentUpdated`: Published when a component is updated
- `EventTypeAgentAction`: Subscribed to for generating components from agent actions
- `EventTypeToolUsage`: Subscribed to for generating components from tool usage

## Django Integration

The Django integration provides REST API endpoints for component generation and retrieval:

- `POST /api/rsc/components`: Generate a component
- `POST /api/rsc/components/agent-action`: Generate a component from an agent action
- `POST /api/rsc/components/tool-usage`: Generate a component from tool usage
- `GET /api/rsc/components/<component_id>`: Get a component by ID
- `GET /api/rsc/components`: List all components
- `GET /api/rsc/components/agent/<agent_id>`: Get components for an agent
- `GET /api/rsc/components/tool/<tool_id>`: Get components for a tool
- `GET /api/rsc/stream`: Stream components as Server-Sent Events (SSE)

## Frontend Integration

The frontend integration provides React hooks and components for rendering the generated components:

- `useComponentStream`: Hook for subscribing to component updates
- `useAgentState`: Hook for accessing agent state
- `AgentUI`: Component for rendering the agent UI
- `AgentAction`: Component for rendering agent actions
- `ToolOutput`: Component for rendering tool usage results

## Usage

### Generating Components from Agent Actions

```go
// Go code
component, err := rscManager.GenerateComponentFromAgentAction(
    ctx,
    agentID,
    actionID,
    "thinking",
    map[string]interface{}{
        "thought": "I need to analyze this code...",
    },
)
```

```python
# Python code
component_id = rsc_integration.generate_component_from_agent_action(
    agent_id="agent1",
    action_id="action1",
    action_type="thinking",
    action_data={
        "thought": "I need to analyze this code...",
    },
)
```

### Generating Components from Tool Usage

```go
// Go code
component, err := rscManager.GenerateComponentFromToolUsage(
    ctx,
    agentID,
    toolID,
    "code_generator",
    map[string]interface{}{
        "language": "python",
        "task": "Generate a function to calculate Fibonacci numbers",
    },
    map[string]interface{}{
        "code": "def fibonacci(n):\n    if n <= 1:\n        return n\n    return fibonacci(n-1) + fibonacci(n-2)",
    },
)
```

```python
# Python code
component_id = rsc_integration.generate_component_from_tool_usage(
    agent_id="agent1",
    tool_id="tool1",
    tool_name="code_generator",
    tool_input={
        "language": "python",
        "task": "Generate a function to calculate Fibonacci numbers",
    },
    tool_output={
        "code": "def fibonacci(n):\n    if n <= 1:\n        return n\n    return fibonacci(n-1) + fibonacci(n-2)",
    },
)
```

### Rendering Components in the UI

```tsx
// React code
import { AgentUI } from './components/AgentUI';

function App() {
  return (
    <div className="app">
      <AgentUI agentId="agent1" />
    </div>
  );
}
```

## Conclusion

The RSC integration enables the agent runtime to generate dynamic UI components based on agent actions and tool usage, creating a more visually appealing and interactive agent chat experience. This integration is similar to CoPilotKit CoAgents but built specifically for our system, providing a more customized and integrated solution.
