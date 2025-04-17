package terminal

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type TerminalManager struct {
	terminals map[string]interfaces.Terminal
	mu        sync.RWMutex
}

func NewTerminalManager() *TerminalManager {
	return &TerminalManager{
		terminals: make(map[string]interfaces.Terminal),
	}
}

func (m *TerminalManager) CreateTerminal(ctx context.Context, terminalType string, id string, options map[string]interface{}) (interfaces.Terminal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.terminals[id]; exists {
		return nil, fmt.Errorf("terminal with ID %s already exists", id)
	}

	var terminal interfaces.Terminal

	switch terminalType {
	case "neovim":
		apiURL, ok := options["api_url"].(string)
		if !ok {
			return nil, fmt.Errorf("api_url is required for neovim terminal")
		}
		terminal = NewNeovimTerminal(id, apiURL, options)
	default:
		return nil, fmt.Errorf("unsupported terminal type: %s", terminalType)
	}

	if err := terminal.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start terminal: %w", err)
	}

	m.terminals[id] = terminal
	return terminal, nil
}

func (m *TerminalManager) GetTerminal(id string) (interfaces.Terminal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	terminal, exists := m.terminals[id]
	if !exists {
		return nil, fmt.Errorf("terminal with ID %s not found", id)
	}

	return terminal, nil
}

func (m *TerminalManager) ListTerminals() []interfaces.Terminal {
	m.mu.RLock()
	defer m.mu.RUnlock()

	terminals := make([]interfaces.Terminal, 0, len(m.terminals))
	for _, terminal := range m.terminals {
		terminals = append(terminals, terminal)
	}

	return terminals
}

func (m *TerminalManager) RemoveTerminal(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	terminal, exists := m.terminals[id]
	if !exists {
		return fmt.Errorf("terminal with ID %s not found", id)
	}

	if err := terminal.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop terminal: %w", err)
	}

	delete(m.terminals, id)
	return nil
}

func (m *TerminalManager) CreateBulkTerminals(ctx context.Context, terminalType string, count int, options map[string]interface{}) ([]interfaces.Terminal, error) {
	terminals := make([]interfaces.Terminal, 0, count)

	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%s_%d", terminalType, i)
		if idPrefix, ok := options["id_prefix"].(string); ok {
			id = fmt.Sprintf("%s%d", idPrefix, i)
		}

		terminal, err := m.CreateTerminal(ctx, terminalType, id, options)
		if err != nil {
			for _, t := range terminals {
				_ = m.RemoveTerminal(ctx, t.ID())
			}
			return nil, fmt.Errorf("failed to create terminal %s: %w", id, err)
		}

		terminals = append(terminals, terminal)
	}

	return terminals, nil
}
