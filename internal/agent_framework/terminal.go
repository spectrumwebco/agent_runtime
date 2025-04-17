package rex

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/internal/terminal"
	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type TerminalRuntime struct {
	*Runtime
	terminalManager interfaces.TerminalManager
	mu              sync.RWMutex
}

func NewTerminalRuntime(agent interfaces.Agent, config RuntimeConfig) (*TerminalRuntime, error) {
	runtime, err := NewRuntime(agent, config)
	if err != nil {
		return nil, err
	}

	return &TerminalRuntime{
		Runtime:         runtime,
		terminalManager: terminal.NewTerminalManager(),
	}, nil
}

func (r *TerminalRuntime) CreateTerminal(ctx context.Context, terminalType string, id string, options map[string]interface{}) (interfaces.Terminal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	terminal, err := r.terminalManager.CreateTerminal(ctx, terminalType, id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create terminal: %w", err)
	}

	r.UpdateContext(map[string]interface{}{
		"terminals": map[string]interface{}{
			id: map[string]interface{}{
				"type":    terminalType,
				"id":      id,
				"running": terminal.IsRunning(),
			},
		},
	})

	return terminal, nil
}

func (r *TerminalRuntime) GetTerminal(id string) (interfaces.Terminal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.terminalManager.GetTerminal(id)
}

func (r *TerminalRuntime) ListTerminals() []interfaces.Terminal {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.terminalManager.ListTerminals()
}

func (r *TerminalRuntime) RemoveTerminal(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.terminalManager.RemoveTerminal(ctx, id)
	if err != nil {
		return err
	}

	context := r.GetContext()
	if terminals, ok := context["terminals"].(map[string]interface{}); ok {
		delete(terminals, id)
		r.UpdateContext(map[string]interface{}{
			"terminals": terminals,
		})
	}

	return nil
}

func (r *TerminalRuntime) ExecuteInTerminal(ctx context.Context, terminalID string, command string) (string, error) {
	terminal, err := r.GetTerminal(terminalID)
	if err != nil {
		return "", err
	}

	return terminal.Execute(ctx, command)
}

func (r *TerminalRuntime) CreateNeovimTerminal(ctx context.Context, id string, apiURL string, options map[string]interface{}) (interfaces.Terminal, error) {
	if options == nil {
		options = make(map[string]interface{})
	}
	options["api_url"] = apiURL

	return r.CreateTerminal(ctx, "neovim", id, options)
}

func (r *TerminalRuntime) CreateBulkTerminals(ctx context.Context, terminalType string, count int, options map[string]interface{}) ([]interfaces.Terminal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	terminals, err := r.terminalManager.CreateBulkTerminals(ctx, terminalType, count, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create bulk terminals: %w", err)
	}

	terminalsContext := make(map[string]interface{})
	for _, t := range terminals {
		terminalsContext[t.ID()] = map[string]interface{}{
			"type":    terminalType,
			"id":      t.ID(),
			"running": t.IsRunning(),
		}
	}

	r.UpdateContext(map[string]interface{}{
		"terminals": terminalsContext,
	})

	return terminals, nil
}

func (r *TerminalRuntime) ExecuteInAllTerminals(ctx context.Context, command string) map[string]string {
	terminals := r.ListTerminals()
	results := make(map[string]string)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, t := range terminals {
		wg.Add(1)
		go func(terminal interfaces.Terminal) {
			defer wg.Done()
			output, err := terminal.Execute(ctx, command)
			
			mu.Lock()
			defer mu.Unlock()
			
			if err != nil {
				results[terminal.ID()] = fmt.Sprintf("ERROR: %s", err.Error())
			} else {
				results[terminal.ID()] = output
			}
		}(t)
	}

	wg.Wait()
	return results
}
