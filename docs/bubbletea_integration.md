# BubbleTea Terminal UI Framework Integration

This document describes the integration of the [Charm BubbleTea](https://github.com/charmbracelet/bubbletea) terminal UI framework into the Kled.io Framework.

## Overview

BubbleTea is a powerful Go framework for building terminal user interfaces based on The Elm Architecture. This integration provides a seamless way to create rich, interactive terminal UIs within the Kled.io Framework.

## Directory Structure

The integration follows the established pattern for the agent_runtime repository:

```
agent_runtime/
├── cmd/cli/commands/
│   └── bubbletea.go         # CLI commands for BubbleTea
├── internal/ui/bubbletea/
│   ├── client.go            # Client wrapper for BubbleTea
│   ├── components.go        # UI components (spinner, progress bar, list)
│   ├── model.go             # Base models and utilities
│   └── examples/            # Example implementations
└── pkg/modules/bubbletea/
    └── module.go            # Module integration for the framework
```

## Features

The BubbleTea integration provides the following components:

1. **Base Models**
   - `BaseModel`: Foundation for all UI models
   - `SequentialModel`: For multi-step UI flows

2. **UI Components**
   - `SpinnerModel`: Animated loading indicator
   - `ProgressBarModel`: Customizable progress visualization
   - `ListModel`: Interactive list selection

3. **CLI Commands**
   - `kled bubbletea example`: Run a simple example
   - `kled bubbletea spinner`: Display a loading spinner
   - `kled bubbletea progress`: Show a progress bar
   - `kled bubbletea list`: Present an interactive list

## Usage

### Basic Example

```go
package main

import (
    "github.com/spectrumwebco/agent_runtime/internal/ui/bubbletea"
    "github.com/spectrumwebco/agent_runtime/pkg/config"
    bubbleteaModule "github.com/spectrumwebco/agent_runtime/pkg/modules/bubbletea"
)

func main() {
    cfg := config.NewDefaultConfig()
    module := bubbleteaModule.NewModule(cfg)
    
    // Create a spinner
    spinner := bubbletea.NewSpinnerModel()
    spinner.Message = "Loading..."
    
    // Run the model
    module.RunModel(spinner)
}
```

### Creating Custom Models

You can create custom models by embedding the `BaseModel`:

```go
type MyCustomModel struct {
    bubbletea.BaseModel
    // Add your custom fields here
}

// Implement the tea.Model interface methods
func (m MyCustomModel) Init() tea.Cmd {
    // Initialize your model
    return nil
}

func (m MyCustomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages and update your model
    return m, nil
}

func (m MyCustomModel) View() string {
    // Render your model
    return "My custom view"
}
```

## Integration with Agent Workflows

The BubbleTea UI components can be integrated with agent workflows to provide rich terminal interfaces for:

1. **Progress Visualization**: Display the progress of long-running agent tasks
2. **Interactive Selection**: Allow users to select options during agent execution
3. **Status Updates**: Show real-time status updates with spinners and messages

## Dependencies

- github.com/charmbracelet/bubbletea v0.25.0
- github.com/charmbracelet/lipgloss (for styling)

## Future Enhancements

1. Add more UI components (tables, forms, etc.)
2. Implement theming support
3. Create integration with agent state visualization
4. Add support for interactive debugging interfaces
