package writing

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type WritingModule struct {
	*modules.BaseModule
}

func NewWritingModule() *WritingModule {
	baseModule := modules.NewBaseModule("writing", "Content creation module")
	
	baseModule.AddTool(tools.Tool{
		Name:        "write_content",
		Description: "Writes content in continuous paragraphs",
		Parameters: map[string]tools.Parameter{
			"topic": {
				Type:        "string",
				Description: "Topic to write about",
			},
			"length": {
				Type:        "string",
				Description: "Desired length of the content",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Generated content",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "cite_references",
		Description: "Cites references in writing",
		Parameters: map[string]tools.Parameter{
			"content": {
				Type:        "string",
				Description: "Content to add citations to",
			},
			"references": {
				Type:        "array",
				Description: "List of references to cite",
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Content with citations",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "compile_document",
		Description: "Compiles multiple sections into a document",
		Parameters: map[string]tools.Parameter{
			"sections": {
				Type:        "array",
				Description: "List of section files to compile",
			},
			"output_path": {
				Type:        "string",
				Description: "Path to save the compiled document",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the compilation",
		},
	})

	return &WritingModule{
		BaseModule: baseModule,
	}
}

func (m *WritingModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *WritingModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
