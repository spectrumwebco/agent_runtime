package file

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type FileModule struct {
	*modules.BaseModule
}

func NewFileModule() *FileModule {
	baseModule := modules.NewBaseModule("file", "File operations module for reading, writing, and editing files")
	
	baseModule.AddTool(tools.Tool{
		Name:        "read_file",
		Description: "Reads content from a file",
		Parameters: map[string]tools.Parameter{
			"path": {
				Type:        "string",
				Description: "Path to the file to read",
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Content of the file",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "write_file",
		Description: "Writes content to a file",
		Parameters: map[string]tools.Parameter{
			"path": {
				Type:        "string",
				Description: "Path to the file to write",
			},
			"content": {
				Type:        "string",
				Description: "Content to write to the file",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the write operation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "append_file",
		Description: "Appends content to a file",
		Parameters: map[string]tools.Parameter{
			"path": {
				Type:        "string",
				Description: "Path to the file to append to",
			},
			"content": {
				Type:        "string",
				Description: "Content to append to the file",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the append operation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "edit_file",
		Description: "Edits content in a file",
		Parameters: map[string]tools.Parameter{
			"path": {
				Type:        "string",
				Description: "Path to the file to edit",
			},
			"line_start": {
				Type:        "integer",
				Description: "Starting line number to edit",
			},
			"line_end": {
				Type:        "integer",
				Description: "Ending line number to edit",
			},
			"new_content": {
				Type:        "string",
				Description: "New content to replace the specified lines",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the edit operation",
		},
	})

	return &FileModule{
		BaseModule: baseModule,
	}
}

func (m *FileModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *FileModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
