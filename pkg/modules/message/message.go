package message

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type MessageModule struct {
	*modules.BaseModule
}

func NewMessageModule() *MessageModule {
	baseModule := modules.NewBaseModule("message", "User communication module for notifications and questions")
	
	baseModule.AddTool(tools.Tool{
		Name:        "notify_user",
		Description: "Sends a non-blocking message to the user",
		Parameters: map[string]tools.Parameter{
			"message": {
				Type:        "string",
				Description: "Message content to send to the user",
			},
			"attachments": {
				Type:        "array",
				Description: "Optional file attachments",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of message delivery",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "ask_user",
		Description: "Sends a blocking message to the user that requires a reply",
		Parameters: map[string]tools.Parameter{
			"message": {
				Type:        "string",
				Description: "Question to ask the user",
			},
			"attachments": {
				Type:        "array",
				Description: "Optional file attachments",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "User's reply to the question",
		},
	})

	return &MessageModule{
		BaseModule: baseModule,
	}
}

func (m *MessageModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *MessageModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
