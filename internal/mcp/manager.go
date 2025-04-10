package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Manager struct {
	config  *config.Config
	servers map[string]*Server
}

type Server struct {
	Name    string
	Server  *server.MCPServer
	Enabled bool
}

type ServerInfo struct {
	Name        string            `json:"name"`
	Enabled     bool              `json:"enabled"`
	Version     string            `json:"version"`
	Capabilities map[string]bool  `json:"capabilities"`
}

func NewManager(cfg *config.Config) (*Manager, error) {
	manager := &Manager{
		config:  cfg,
		servers: make(map[string]*Server),
	}
	
	if err := manager.initServers(); err != nil {
		return nil, err
	}
	
	return manager, nil
}

func (m *Manager) initServers() error {
	if err := m.initFilesystemServer(); err != nil {
		return err
	}
	
	if err := m.initToolsServer(); err != nil {
		return err
	}
	
	if err := m.initRuntimeServer(); err != nil {
		return err
	}
	
	return nil
}

func (m *Manager) initFilesystemServer() error {
	enabled := false
	for _, s := range m.config.MCP.Servers {
		if s.Name == "filesystem" {
			enabled = s.Enabled
			break
		}
	}
	
	mcpServer := server.NewMCPServer(
		"agent-runtime/filesystem",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
	)
	
	m.servers["filesystem"] = &Server{
		Name:    "filesystem",
		Server:  mcpServer,
		Enabled: enabled,
	}
	
	return nil
}

func (m *Manager) initToolsServer() error {
	enabled := false
	for _, s := range m.config.MCP.Servers {
		if s.Name == "tools" {
			enabled = s.Enabled
			break
		}
	}
	
	mcpServer := server.NewMCPServer(
		"agent-runtime/tools",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	
	m.servers["tools"] = &Server{
		Name:    "tools",
		Server:  mcpServer,
		Enabled: enabled,
	}
	
	return nil
}

func (m *Manager) initRuntimeServer() error {
	enabled := false
	for _, s := range m.config.MCP.Servers {
		if s.Name == "runtime" {
			enabled = s.Enabled
			break
		}
	}
	
	mcpServer := server.NewMCPServer(
		"agent-runtime/runtime",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	
	m.servers["runtime"] = &Server{
		Name:    "runtime",
		Server:  mcpServer,
		Enabled: enabled,
	}
	
	return nil
}

func (m *Manager) ListServers() []ServerInfo {
	servers := make([]ServerInfo, 0, len(m.servers))
	
	for _, s := range m.servers {
		servers = append(servers, ServerInfo{
			Name:    s.Name,
			Enabled: s.Enabled,
			Version: "1.0.0",
			Capabilities: map[string]bool{
				"resources": true,
				"tools":     true,
			},
		})
	}
	
	return servers
}

func (m *Manager) GetServer(name string) (*ServerInfo, error) {
	server, ok := m.servers[name]
	if !ok {
		return nil, fmt.Errorf("server not found: %s", name)
	}
	
	return &ServerInfo{
		Name:    server.Name,
		Enabled: server.Enabled,
		Version: "1.0.0",
		Capabilities: map[string]bool{
			"resources": true,
			"tools":     true,
		},
	}, nil
}

func (m *Manager) StartServers() error {
	for _, s := range m.servers {
		if s.Enabled {
		}
	}
	
	return nil
}

func (m *Manager) StopServers() error {
	for _, s := range m.servers {
		if s.Enabled {
		}
	}
	
	return nil
}
