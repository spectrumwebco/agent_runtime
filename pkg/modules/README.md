# Agent Runtime Module System

The Agent Runtime Module System provides a modular architecture for extending the agent's capabilities through specialized modules. Each module implements a specific set of functionality and provides tools that the agent can use during task execution.

## Module Architecture

The Agent Runtime uses a 3-tiered architecture for modules:

1. **Core Modules**: Essential modules for basic agent functionality (planner, knowledge, datasource)
2. **Specialized Modules**: Task-specific modules for enhanced capabilities (todo, message, file, info)
3. **Interaction Modules**: Modules for external system interaction (browser, shell, coding, deploy, writing)

Each module follows a consistent event-based communication pattern and supports the agent loop execution model:
- Analyze Events: Process user needs and current state
- Select Tools: Choose appropriate tools based on context
- Wait for Execution: Allow sandbox to execute actions
- Iterate: Repeat steps until task completion
- Submit Results: Provide deliverables to users

## Module Interface

All modules implement the Module interface defined in `module.go`, which includes:

- `Name()`: Returns the module name
- `Description()`: Returns the module description
- `Tools()`: Returns a list of tools provided by the module
- `Initialize(context)`: Initializes the module with execution context
- `Cleanup()`: Cleans up module resources

## Available Modules

### Core Modules

#### Planner Module

The Planner module is responsible for overall task planning and execution tracking. It provides tools for:

- Creating task execution plans with numbered pseudocode steps
- Updating plans with progress information
- Tracking the current step number, status, and reflection

#### Knowledge Module

The Knowledge module provides best practice references and memory capabilities. It offers tools for:

- Retrieving task-relevant knowledge and best practices
- Storing new knowledge for future reference
- Managing knowledge with different scopes (general, task-specific, domain-specific)

#### Datasource Module

The Datasource module enables access to authoritative data sources. It provides tools for:

- Listing available data APIs and their documentation
- Retrieving detailed documentation for specific APIs
- Generating Python code templates for using data APIs

### Task Management Modules

#### Todo Module

The Todo module manages task tracking and completion status. It provides tools for:

- Creating todo lists based on task planning
- Updating todo item statuses (completed, in_progress, skipped)
- Rebuilding todo lists when task planning changes

#### Message Module

The Message module handles user communication. It offers tools for:

- Sending non-blocking notifications to users
- Asking questions that require user responses
- Attaching files to messages

#### File Module

The File module manages file operations. It provides tools for:

- Reading content from files
- Writing content to files
- Appending content to existing files
- Editing specific lines in files

#### Info Module

The Info module handles information retrieval and prioritization. It offers tools for:

- Searching the web for information
- Prioritizing information based on reliability, relevance, or recency

### Interaction Modules

#### Browser Module

The Browser module enables web interaction. It provides tools for:

- Navigating to URLs
- Viewing page content
- Clicking on page elements
- Scrolling through pages

#### Shell Module

The Shell module manages command-line operations. It offers tools for:

- Executing individual shell commands
- Chaining multiple commands with && operator
- Piping output between commands

#### Coding Module

The Coding module handles code generation and execution. It provides tools for:

- Writing code to files
- Executing code from files
- Searching for code solutions

#### Deploy Module

The Deploy module manages deployment operations. It offers tools for:

- Exposing ports for external access
- Deploying static websites
- Deploying applications (web, API, worker)

#### Writing Module

The Writing module assists with content creation. It provides tools for:

- Writing content on specific topics
- Citing references in writing
- Compiling multiple sections into documents

#### Error Handling Module

The Error Handling module manages tool execution failures. It offers tools for:

- Verifying tool names and arguments
- Fixing issues based on error messages
- Reporting failures to users when automatic fixes fail

## Module Registration

Modules must be registered with the Registry to be available to the agent:

```go
registry := modules.NewRegistry()
// Core Modules
registry.Register(planner.NewPlannerModule())
registry.Register(knowledge.NewKnowledgeModule())
registry.Register(datasource.NewDatasourceModule())

// Task Management Modules
registry.Register(todo.NewTodoModule())
registry.Register(message.NewMessageModule())
registry.Register(file.NewFileModule())
registry.Register(info.NewInfoModule())

// Interaction Modules
registry.Register(browser.NewBrowserModule())
registry.Register(shell.NewShellModule())
registry.Register(coding.NewCodingModule())
registry.Register(deploy.NewDeployModule())
registry.Register(writing.NewWritingModule())
registry.Register(error_handling.NewErrorHandlingModule())
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
