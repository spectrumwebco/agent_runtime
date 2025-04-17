package tools

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type NeovimTerminalTool struct {
	name           string
	description    string
	terminalRuntime interfaces.TerminalManager
	apiURL         string
}

func NewNeovimTerminalTool(terminalRuntime interfaces.TerminalManager, apiURL string) *NeovimTerminalTool {
	return &NeovimTerminalTool{
		name:           "neovim_terminal",
		description:    "Execute commands in Neovim terminals",
		terminalRuntime: terminalRuntime,
		apiURL:         apiURL,
	}
}

func (t *NeovimTerminalTool) Name() string {
	return t.name
}

func (t *NeovimTerminalTool) Description() string {
	return t.description
}

func (t *NeovimTerminalTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	action, ok := params["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action parameter is required")
	}

	switch action {
	case "create":
		return t.createTerminal(ctx, params)
	case "execute":
		return t.executeCommand(ctx, params)
	case "stop":
		return t.stopTerminal(ctx, params)
	case "list":
		return t.listTerminals(ctx)
	case "create_bulk":
		return t.createBulkTerminals(ctx, params)
	case "execute_all":
		return t.executeInAllTerminals(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}

func (t *NeovimTerminalTool) createTerminal(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id parameter is required")
	}

	options := make(map[string]interface{})
	if opts, ok := params["options"].(map[string]interface{}); ok {
		options = opts
	}
	options["api_url"] = t.apiURL

	terminal, err := t.terminalRuntime.CreateTerminal(ctx, "neovim", id, options)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":      terminal.ID(),
		"type":    terminal.GetType(),
		"running": terminal.IsRunning(),
	}, nil
}

func (t *NeovimTerminalTool) executeCommand(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id parameter is required")
	}

	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}

	terminal, err := t.terminalRuntime.GetTerminal(id)
	if err != nil {
		return nil, err
	}

	output, err := terminal.Execute(ctx, command)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":     terminal.ID(),
		"output": output,
	}, nil
}

func (t *NeovimTerminalTool) stopTerminal(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id parameter is required")
	}

	err := t.terminalRuntime.RemoveTerminal(ctx, id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"id":      id,
	}, nil
}

func (t *NeovimTerminalTool) listTerminals(ctx context.Context) (interface{}, error) {
	terminals := t.terminalRuntime.ListTerminals()
	result := make([]map[string]interface{}, 0, len(terminals))

	for _, terminal := range terminals {
		result = append(result, map[string]interface{}{
			"id":      terminal.ID(),
			"type":    terminal.GetType(),
			"running": terminal.IsRunning(),
		})
	}

	return result, nil
}

func (t *NeovimTerminalTool) createBulkTerminals(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	count, ok := params["count"].(float64)
	if !ok {
		return nil, fmt.Errorf("count parameter is required")
	}

	options := make(map[string]interface{})
	if opts, ok := params["options"].(map[string]interface{}); ok {
		options = opts
	}
	options["api_url"] = t.apiURL

	terminals, err := t.terminalRuntime.CreateBulkTerminals(ctx, "neovim", int(count), options)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(terminals))
	for _, terminal := range terminals {
		result = append(result, map[string]interface{}{
			"id":      terminal.ID(),
			"type":    terminal.GetType(),
			"running": terminal.IsRunning(),
		})
	}

	return result, nil
}

func (t *NeovimTerminalTool) executeInAllTerminals(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}

	results := make(map[string]string)
	terminals := t.terminalRuntime.ListTerminals()

	for _, terminal := range terminals {
		output, err := terminal.Execute(ctx, command)
		if err != nil {
			results[terminal.ID()] = fmt.Sprintf("ERROR: %s", err.Error())
		} else {
			results[terminal.ID()] = output
		}
	}

	return results, nil
}
