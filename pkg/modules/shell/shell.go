package shell

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type ShellModule struct {
	*modules.BaseModule
}

func NewShellModule() *ShellModule {
	baseModule := modules.NewBaseModule("shell", "Shell command execution module")
	
	baseModule.AddTool(tools.Tool{
		Name:        "execute_shell",
		Description: "Executes a shell command",
		Parameters: map[string]tools.Parameter{
			"command": {
				Type:        "string",
				Description: "Shell command to execute",
			},
			"working_dir": {
				Type:        "string",
				Description: "Working directory for the command",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Command output",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "chain_commands",
		Description: "Chains multiple shell commands with && operator",
		Parameters: map[string]tools.Parameter{
			"commands": {
				Type:        "array",
				Description: "List of shell commands to chain",
			},
			"working_dir": {
				Type:        "string",
				Description: "Working directory for the commands",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Combined command output",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "pipe_commands",
		Description: "Pipes output between shell commands",
		Parameters: map[string]tools.Parameter{
			"command_pipeline": {
				Type:        "array",
				Description: "List of shell commands to pipe",
			},
			"working_dir": {
				Type:        "string",
				Description: "Working directory for the commands",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "string",
			Description: "Pipeline output",
		},
	})

	return &ShellModule{
		BaseModule: baseModule,
	}
}

func (m *ShellModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *ShellModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
