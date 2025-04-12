package browser

import (
	"context"
	
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type BrowserModule struct {
	*modules.BaseModule
}

func NewBrowserModule() *BrowserModule {
	baseModule := modules.NewBaseModule("browser", "Browser interaction module for web navigation and interaction")
	
	baseModule.AddTool(tools.Tool{
		Name:        "navigate_browser",
		Description: "Navigates to a URL in the browser",
		Parameters: map[string]tools.Parameter{
			"url": {
				Type:        "string",
				Description: "URL to navigate to",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the navigation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "view_browser",
		Description: "Views the current browser page",
		Parameters: map[string]tools.Parameter{
			"reload": {
				Type:        "boolean",
				Description: "Whether to reload the page before viewing",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Page content with visible elements",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "click_browser",
		Description: "Clicks on an element in the browser",
		Parameters: map[string]tools.Parameter{
			"element_index": {
				Type:        "integer",
				Description: "Index of the element to click",
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the click operation",
		},
	})

	baseModule.AddTool(tools.Tool{
		Name:        "scroll_browser",
		Description: "Scrolls the browser page",
		Parameters: map[string]tools.Parameter{
			"direction": {
				Type:        "string",
				Description: "Direction to scroll",
				Enum:        []string{"up", "down"},
			},
			"amount": {
				Type:        "integer",
				Description: "Amount to scroll in pixels",
				Required:    false,
			},
		},
		Returns: tools.ReturnType{
			Type:        "object",
			Description: "Status of the scroll operation",
		},
	})

	return &BrowserModule{
		BaseModule: baseModule,
	}
}

func (m *BrowserModule) Initialize(ctx context.Context) error {
	return m.BaseModule.Initialize(ctx)
}

func (m *BrowserModule) Cleanup() error {
	return m.BaseModule.Cleanup()
}
