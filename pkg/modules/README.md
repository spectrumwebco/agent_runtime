# Agent Runtime Module System

The Agent Runtime Module System provides a modular architecture for extending the agent's capabilities through specialized modules. Each module implements a specific set of functionality and provides tools that the agent can use during task execution.

## Module Interface

All modules implement the Module interface defined in `module.go`, which includes:

- `Name()`: Returns the module name
- `Description()`: Returns the module description
- `Tools()`: Returns a list of tools provided by the module
- `Initialize(context)`: Initializes the module with execution context
- `Cleanup()`: Cleans up module resources

## Available Modules

### Planner Module

The Planner module is responsible for overall task planning and execution tracking. It provides tools for:

- Creating task execution plans with numbered pseudocode steps
- Updating plans with progress information
- Tracking the current step number, status, and reflection

### Knowledge Module

The Knowledge module provides best practice references and memory capabilities. It offers tools for:

- Retrieving task-relevant knowledge and best practices
- Storing new knowledge for future reference
- Managing knowledge with different scopes (general, task-specific, domain-specific)

### Datasource Module

The Datasource module enables access to authoritative data sources. It provides tools for:

- Listing available data APIs and their documentation
- Retrieving detailed documentation for specific APIs
- Generating Python code templates for using data APIs

## Module Registration

Modules must be registered with the Registry to be available to the agent:

```go
registry := modules.NewRegistry()
registry.Register(planner.NewPlannerModule())
registry.Register(knowledge.NewKnowledgeModule())
registry.Register(datasource.NewDatasourceModule())
```

## Implementing New Modules

To implement a new module:

1. Create a new package in the `pkg/modules` directory
2. Define a struct that embeds `*modules.BaseModule`
3. Implement the `Initialize` and `Cleanup` methods
4. Create a constructor function that initializes the base module and adds tools
5. Register the module with the Registry

Example:

```go
package mymodule

import (
	"context"

	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type MyModule struct {
	*modules.BaseModule
}

func NewMyModule() *MyModule {
	baseModule := modules.NewBaseModule("mymodule", "Description of my module")
	
	// Add module-specific tools
	baseModule.AddTool(tools.Tool{
		Name:        "my_tool",
		Description: "Description of my tool",
		Parameters: map[string]tools.Parameter{
			"param1": {
				Type:        "string",
				Description: "Description of parameter 1",
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Description of return value",
		},
	})

	return &MyModule{
		BaseModule: baseModule,
	}
}

func (m *MyModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *MyModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
```
