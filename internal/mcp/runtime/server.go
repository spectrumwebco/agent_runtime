package runtime

import (
	"fmt"
	"log"

	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/env"
	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Server struct {
	agent       *agent.Agent
	environment *env.SWEEnv
	toolManager *tools.Registry
	mcpServer   *mcp.RuntimeServer
}

func NewServer(agentInstance *agent.Agent, environment *env.SWEEnv, toolManager *tools.Registry) (*Server, error) {
	mcpServer, err := mcp.NewRuntimeServer(agentInstance, environment, toolManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP runtime server: %w", err)
	}
	
	return &Server{
		agent:       agentInstance,
		environment: environment,
		toolManager: toolManager,
		mcpServer:   mcpServer,
	}, nil
}

func (s *Server) Run() {
	log.Println("Runtime MCP Server running...")
	select {} // Keep running
}
