package interfaces

import (
	"context"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Agent interface {
	ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
	
	RegisterTool(tool tools.Tool) error
	
	GetTool(name string) (tools.Tool, error)
	
	ListTools() []tools.Tool
}
