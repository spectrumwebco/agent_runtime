package todo

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type TodoModule struct {
	*modules.BaseModule
}

func NewTodoModule() *TodoModule {
	baseModule := modules.NewBaseModule("todo", "Todo list management module for task tracking")
	
	baseModule.AddTool(tools.Tool{
		Name:        "create_todo",
		Description: "Creates a new todo list based on task planning",
		Parameters: map[string]tools.Parameter{
			"task_plan": {
				Type:        "string",
				Description: "Task plan to create todo list from",
			},
			"file_path": {
				Type:        "string",
				Description: "Path to save the todo.md file",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of todo list creation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "update_todo",
		Description: "Updates a todo item's status",
		Parameters: map[string]tools.Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the todo.md file",
			},
			"item_number": {
				Type:        "integer",
				Description: "Number of the todo item to update",
			},
			"status": {
				Type:        "string",
				Description: "New status of the todo item (completed, in_progress, skipped)",
				Enum:        []string{"completed", "in_progress", "skipped"},
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of todo item update",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "rebuild_todo",
		Description: "Rebuilds the todo list when task planning changes",
		Parameters: map[string]tools.Parameter{
			"task_plan": {
				Type:        "string",
				Description: "Updated task plan",
			},
			"file_path": {
				Type:        "string",
				Description: "Path to the todo.md file",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of todo list rebuild",
		},
	})

	return &TodoModule{
		BaseModule: baseModule,
	}
}

func (m *TodoModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *TodoModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
