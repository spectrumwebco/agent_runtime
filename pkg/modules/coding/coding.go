package coding

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type CodingModule struct {
	*modules.BaseModule
}

func NewCodingModule() *CodingModule {
	baseModule := modules.NewBaseModule("coding", "Code generation and execution module")
	
	baseModule.AddTool(tools.Tool{
		Name:        "write_code",
		Description: "Writes code to a file",
		Parameters: map[string]tools.Parameter{
			"path": {
				Type:        "string",
				Description: "Path to save the code file",
			},
			"content": {
				Type:        "string",
				Description: "Code content to write",
			},
			"language": {
				Type:        "string",
				Description: "Programming language of the code",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the write operation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "execute_code",
		Description: "Executes code from a file",
		Parameters: map[string]tools.Parameter{
			"path": {
				Type:        "string",
				Description: "Path to the code file to execute",
			},
			"language": {
				Type:        "string",
				Description: "Programming language of the code",
			},
			"args": {
				Type:        "array",
				Description: "Command line arguments for execution",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Execution result with stdout and stderr",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "search_solutions",
		Description: "Searches for code solutions",
		Parameters: map[string]tools.Parameter{
			"problem": {
				Type:        "string",
				Description: "Problem description to search solutions for",
			},
			"language": {
				Type:        "string",
				Description: "Programming language for solutions",
			},
		},
		Returns: tools.ReturnType{
			Type:        "array",
			Description: "Code solution suggestions",
		},
	})

	return &CodingModule{
		BaseModule: baseModule,
	}
}

func (m *CodingModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *CodingModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
