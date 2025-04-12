package error_handling

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type ErrorHandlingModule struct {
	*modules.BaseModule
}

func NewErrorHandlingModule() *ErrorHandlingModule {
	baseModule := modules.NewBaseModule("error_handling", "Error handling module for tool execution failures")
	
	baseModule.AddTool(tools.Tool{
		Name:        "verify_tool",
		Description: "Verifies tool names and arguments",
		Parameters: map[string]tools.Parameter{
			"tool_name": {
				Type:        "string",
				Description: "Name of the tool to verify",
			},
			"arguments": {
				Type:        "object",
				Description: "Arguments to verify",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Verification result with any issues found",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "fix_error",
		Description: "Attempts to fix issues based on error messages",
		Parameters: map[string]tools.Parameter{
			"error_message": {
				Type:        "string",
				Description: "Error message to analyze",
			},
			"tool_name": {
				Type:        "string",
				Description: "Name of the tool that produced the error",
			},
			"arguments": {
				Type:        "object",
				Description: "Arguments that were used",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Fix result with corrected arguments or alternative approach",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "report_failure",
		Description: "Reports failure reasons to the user and requests assistance",
		Parameters: map[string]tools.Parameter{
			"tool_name": {
				Type:        "string",
				Description: "Name of the tool that failed",
			},
			"error_message": {
				Type:        "string",
				Description: "Error message from the failure",
			},
			"attempts": {
				Type:        "array",
				Description: "Previous attempts to fix the issue",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the failure report",
		},
	})

	return &ErrorHandlingModule{
		BaseModule: baseModule,
	}
}

func (m *ErrorHandlingModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *ErrorHandlingModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
