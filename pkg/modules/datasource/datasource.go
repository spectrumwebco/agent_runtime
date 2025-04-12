package datasource

import (
	"context"

	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type DatasourceModule struct {
	*modules.BaseModule
}

func NewDatasourceModule() *DatasourceModule {
	baseModule := modules.NewBaseModule("datasource", "Data API module for accessing authoritative datasources")
	
	baseModule.AddTool(tools.Tool{
		Name:        "list_data_apis",
		Description: "Lists available data APIs and their documentation",
		Parameters: map[string]tools.Parameter{
			"category": {
				Type:        "string",
				Description: "Optional category to filter APIs",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "array",
			Description: "List of available data APIs with documentation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "get_api_documentation",
		Description: "Retrieves detailed documentation for a specific data API",
		Parameters: map[string]tools.Parameter{
			"api_name": {
				Type:        "string",
				Description: "Name of the data API to get documentation for",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Detailed API documentation including endpoints and parameters",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "generate_api_code",
		Description: "Generates Python code template for using a specific data API",
		Parameters: map[string]tools.Parameter{
			"api_name": {
				Type:        "string",
				Description: "Name of the data API to generate code for",
			},
			"endpoint": {
				Type:        "string",
				Description: "Specific API endpoint to use",
			},
			"parameters": {
				Type:        "object",
				Description: "Parameters to include in the API call",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Python code template for using the specified API",
		},
	})

	return &DatasourceModule{
		BaseModule: baseModule,
	}
}

func (m *DatasourceModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *DatasourceModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
