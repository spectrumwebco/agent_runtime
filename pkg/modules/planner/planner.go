package planner

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type PlannerModule struct {
	*modules.BaseModule
}

func NewPlannerModule() *PlannerModule {
	baseModule := modules.NewBaseModule("planner", "Task planning module for overall task planning and execution tracking")
	
	baseModule.AddTool(tools.Tool{
		Name:        "create_plan",
		Description: "Creates a new task execution plan with numbered pseudocode steps",
		Parameters: map[string]tools.Parameter{
			"task": {
				Type:        "string",
				Description: "Task description to plan for",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Task plan with numbered steps",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "update_plan",
		Description: "Updates the current task execution plan with progress information",
		Parameters: map[string]tools.Parameter{
			"step_number": {
				Type:        "integer",
				Description: "Current step number being executed",
			},
			"status": {
				Type:        "string",
				Description: "Status of the current step (in_progress, completed, failed)",
				Enum:        []string{"in_progress", "completed", "failed"},
			},
			"reflection": {
				Type:        "string",
				Description: "Reflection on the current step execution",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Updated task plan",
		},
	})

	return &PlannerModule{
		BaseModule: baseModule,
	}
}

func (m *PlannerModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *PlannerModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
