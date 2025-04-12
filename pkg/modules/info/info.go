package info

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type InfoModule struct {
	*modules.BaseModule
}

func NewInfoModule() *InfoModule {
	baseModule := modules.NewBaseModule("info", "Information retrieval and prioritization module")
	
	baseModule.AddTool(tools.Tool{
		Name:        "search_web",
		Description: "Searches the web for information",
		Parameters: map[string]tools.Parameter{
			"query": {
				Type:        "string",
				Description: "Search query",
			},
			"num_results": {
				Type:        "integer",
				Description: "Number of results to return",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "array",
			Description: "Search results with URLs and snippets",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "prioritize_info",
		Description: "Prioritizes information based on source",
		Parameters: map[string]tools.Parameter{
			"info_items": {
				Type:        "array",
				Description: "Information items to prioritize",
			},
			"criteria": {
				Type:        "string",
				Description: "Criteria for prioritization",
				Enum:        []string{"reliability", "relevance", "recency"},
			},
		},
		Returns: tools.ReturnType{
			Type:        "array",
			Description: "Prioritized information items",
		},
	})

	return &InfoModule{
		BaseModule: baseModule,
	}
}

func (m *InfoModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *InfoModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
