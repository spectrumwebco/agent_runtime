package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type NeovimTerminal struct {
	id        string
	apiURL    string
	running   bool
	mu        sync.RWMutex
	client    *http.Client
	options   map[string]interface{}
}

func NewNeovimTerminal(id string, apiURL string, options map[string]interface{}) *NeovimTerminal {
	return &NeovimTerminal{
		id:      id,
		apiURL:  apiURL,
		running: false,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		options: options,
	}
}

func (t *NeovimTerminal) ID() string {
	return t.id
}

func (t *NeovimTerminal) Execute(ctx context.Context, command string) (string, error) {
	t.mu.RLock()
	if !t.running {
		t.mu.RUnlock()
		return "", fmt.Errorf("terminal %s is not running", t.id)
	}
	t.mu.RUnlock()

	neovimCommand := command
	if !strings.HasPrefix(command, ":") {
		neovimCommand = fmt.Sprintf(":terminal %s<CR>", command)
	}

	payload := map[string]interface{}{
		"id":      t.id,
		"command": neovimCommand,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal command payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/execute", t.apiURL), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success bool   `json:"success"`
		Output  string `json:"output"`
		Error   string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("command execution failed: %s", result.Error)
	}

	return result.Output, nil
}

func (t *NeovimTerminal) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return fmt.Errorf("terminal %s is already running", t.id)
	}

	payload := map[string]interface{}{
		"id": t.id,
	}

	if t.options != nil {
		for k, v := range t.options {
			payload[k] = v
		}
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal start payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/start", t.apiURL), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start terminal: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("terminal start failed: %s", result.Error)
	}

	t.running = true
	return nil
}

func (t *NeovimTerminal) Stop(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return nil
	}

	payload := map[string]interface{}{
		"id": t.id,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal stop payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/stop", t.apiURL), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to stop terminal: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("terminal stop failed: %s", result.Error)
	}

	t.running = false
	return nil
}

func (t *NeovimTerminal) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

func (t *NeovimTerminal) GetType() string {
	return "neovim"
}
