package terminal

import (
	"context"
	"fmt"
	"sync"
	
	"github.com/spectrumwebco/agent_runtime/internal/containers"
	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
	"github.com/spectrumwebco/agent_runtime/pkg/librechat"
)

type Manager struct {
	terminals     map[string]interfaces.Terminal
	mu            sync.RWMutex
	libreChat     *librechat.Client
	libreChatURL  string
	defaultAPIURL string
}

func NewManager(defaultAPIURL, libreChatURL, libreChatAPIKey string) *Manager {
	return &Manager{
		terminals:     make(map[string]interfaces.Terminal),
		mu:            sync.RWMutex{},
		libreChat:     librechat.NewClient(libreChatURL, libreChatAPIKey),
		libreChatURL:  libreChatURL,
		defaultAPIURL: defaultAPIURL,
	}
}

func (m *Manager) CreateTerminal(ctx context.Context, terminalType string, id string, options map[string]interface{}) (interfaces.Terminal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.terminals[id]; exists {
		return nil, fmt.Errorf("terminal with ID %s already exists", id)
	}
	
	var terminal interfaces.Terminal
	
	if options == nil {
		options = make(map[string]interface{})
	}
	
	if _, ok := options["api_url"]; !ok {
		options["api_url"] = m.defaultAPIURL
	}
	
	if _, ok := options["librechat_client"]; !ok {
		options["librechat_client"] = m.libreChat
	}
	
	switch terminalType {
	case "neovim":
		apiURL, ok := options["api_url"].(string)
		if !ok {
			return nil, fmt.Errorf("api_url is required for neovim terminal")
		}
		terminal = NewNeovimTerminal(id, apiURL, options)
		
	case "neovim-kata":
		apiURL, ok := options["api_url"].(string)
		if !ok {
			return nil, fmt.Errorf("api_url is required for neovim-kata terminal")
		}
		
		kataConfigPath, ok := options["kata_config_path"].(string)
		if !ok {
			kataConfigPath = "/etc/kata-containers/configuration.toml"
		}
		
		kataRuntimeDir, ok := options["kata_runtime_dir"].(string)
		if !ok {
			kataRuntimeDir = "/var/run/kata-containers"
		}
		
		kataManager := containers.NewKataManager(nil, kataConfigPath, kataRuntimeDir)
		
		terminal = NewNeovimKataTerminal(id, apiURL, options, kataManager)
		
	case "tmux":
		apiURL, ok := options["api_url"].(string)
		if !ok {
			return nil, fmt.Errorf("api_url is required for tmux terminal")
		}
		terminal = NewTmuxTerminal(id, apiURL, options)
		
	case "ssh":
		host, ok := options["host"].(string)
		if !ok {
			return nil, fmt.Errorf("host is required for ssh terminal")
		}
		
		user, ok := options["user"].(string)
		if !ok {
			return nil, fmt.Errorf("user is required for ssh terminal")
		}
		
		password, _ := options["password"].(string)
		keyPath, _ := options["key_path"].(string)
		
		terminal = NewSSHTerminal(id, host, user, password, keyPath, options)
		
	default:
		return nil, fmt.Errorf("unsupported terminal type: %s", terminalType)
	}
	
	m.terminals[id] = terminal
	return terminal, nil
}

func (m *Manager) GetTerminal(id string) (interfaces.Terminal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	terminal, exists := m.terminals[id]
	if !exists {
		return nil, fmt.Errorf("terminal with ID %s not found", id)
	}
	
	return terminal, nil
}

func (m *Manager) ListTerminals() []interfaces.Terminal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	terminals := make([]interfaces.Terminal, 0, len(m.terminals))
	for _, terminal := range m.terminals {
		terminals = append(terminals, terminal)
	}
	
	return terminals
}

func (m *Manager) RemoveTerminal(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	terminal, exists := m.terminals[id]
	if !exists {
		return fmt.Errorf("terminal with ID %s not found", id)
	}
	
	if terminal.IsRunning() {
		if err := terminal.Stop(ctx); err != nil {
			return fmt.Errorf("stop terminal: %w", err)
		}
	}
	
	delete(m.terminals, id)
	return nil
}

func (m *Manager) CreateBulkTerminals(ctx context.Context, terminalType string, count int, options map[string]interface{}) ([]interfaces.Terminal, error) {
	terminals := make([]interfaces.Terminal, 0, count)
	
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%s-%d", terminalType, i)
		terminal, err := m.CreateTerminal(ctx, terminalType, id, options)
		if err != nil {
			for _, t := range terminals {
				_ = m.RemoveTerminal(ctx, t.ID())
			}
			return nil, fmt.Errorf("create terminal %s: %w", id, err)
		}
		
		terminals = append(terminals, terminal)
	}
	
	return terminals, nil
}

func (m *Manager) GetLibreChatClient() *librechat.Client {
	return m.libreChat
}

func (m *Manager) ExecuteCodeInTerminal(ctx context.Context, terminalID, language, code string) (string, error) {
	terminal, err := m.GetTerminal(terminalID)
	if err != nil {
		return "", err
	}
	
	output, err := m.libreChat.ExecuteCode(ctx, language, code)
	if err != nil {
		return "", fmt.Errorf("execute code with LibreChat: %w", err)
	}
	
	if _, err := terminal.Execute(ctx, fmt.Sprintf("echo '%s'", output)); err != nil {
		return "", fmt.Errorf("send output to terminal: %w", err)
	}
	
	return output, nil
}
