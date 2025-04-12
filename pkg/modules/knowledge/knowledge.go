package knowledge

import (
	"context"

	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type KnowledgeModule struct {
	*modules.BaseModule
}

func NewKnowledgeModule() *KnowledgeModule {
	baseModule := modules.NewBaseModule("knowledge", "Knowledge and memory module for best practice references")
	
	baseModule.AddTool(tools.Tool{
		Name:        "retrieve_knowledge",
		Description: "Retrieves task-relevant knowledge and best practices",
		Parameters: map[string]tools.Parameter{
			"topic": {
				Type:        "string",
				Description: "Topic to retrieve knowledge for",
			},
			"scope": {
				Type:        "string",
				Description: "Scope of the knowledge (e.g., 'general', 'task-specific')",
				Enum:        []string{"general", "task-specific", "domain-specific"},
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Knowledge item with scope and content",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "store_knowledge",
		Description: "Stores new knowledge for future reference",
		Parameters: map[string]tools.Parameter{
			"topic": {
				Type:        "string",
				Description: "Topic of the knowledge",
			},
			"content": {
				Type:        "string",
				Description: "Content of the knowledge",
			},
			"scope": {
				Type:        "string",
				Description: "Scope of the knowledge (e.g., 'general', 'task-specific')",
				Enum:        []string{"general", "task-specific", "domain-specific"},
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the knowledge storage operation",
		},
	})

	return &KnowledgeModule{
		BaseModule: baseModule,
	}
}

func (m *KnowledgeModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *KnowledgeModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
