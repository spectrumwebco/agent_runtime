package deploy

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type DeployModule struct {
	*modules.BaseModule
}

func NewDeployModule() *DeployModule {
	baseModule := modules.NewBaseModule("deploy", "Deployment operations module")
	
	baseModule.AddTool(tools.Tool{
		Name:        "expose_port",
		Description: "Exposes a port for external access",
		Parameters: map[string]tools.Parameter{
			"port": {
				Type:        "integer",
				Description: "Port number to expose",
			},
			"service_name": {
				Type:        "string",
				Description: "Name of the service running on the port",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Public access URL for the exposed port",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "deploy_website",
		Description: "Deploys a static website",
		Parameters: map[string]tools.Parameter{
			"source_dir": {
				Type:        "string",
				Description: "Directory containing the website files",
			},
			"permanent": {
				Type:        "boolean",
				Description: "Whether to deploy permanently",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Deployment information with access URL",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "deploy_application",
		Description: "Deploys an application",
		Parameters: map[string]tools.Parameter{
			"source_dir": {
				Type:        "string",
				Description: "Directory containing the application files",
			},
			"app_type": {
				Type:        "string",
				Description: "Type of application to deploy",
				Enum:        []string{"web", "api", "worker"},
			},
			"permanent": {
				Type:        "boolean",
				Description: "Whether to deploy permanently",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Deployment information with access URL",
		},
	})

	return &DeployModule{
		BaseModule: baseModule,
	}
}

func (m *DeployModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *DeployModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
