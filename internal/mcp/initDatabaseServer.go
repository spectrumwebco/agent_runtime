package mcp

import (
	"fmt"
	
	"github.com/mark3labs/mcp-go/server"
)

func (m *Manager) initDatabaseServer() error {
	enabled := false
	for _, s := range m.config.MCP.Servers {
		if s.Name == "database" {
			enabled = s.Enabled
			break
		}
	}
	
	mcpServer := server.NewMCPServer(
		"gitops-go/database",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)
	
	m.servers["database"] = &Server{
		Name:    "database",
		Server:  mcpServer,
		Enabled: enabled,
	}
	
	return nil
}
