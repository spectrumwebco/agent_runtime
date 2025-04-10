package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/internal/server"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	agentInstance, err := agent.NewAgent(agent.AgentConfig{
		Name:        cfg.Agent.Name,
		Description: cfg.Agent.Description,
		ToolsPath:   "pkg/tools/tools.json",
		PromptsPath: "pkg/prompts/prompts.txt",
		ModulesPath: "pkg/modules",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	mcpManager, err := mcp.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create MCP manager: %v", err)
	}

	httpServer, err := server.NewHTTPServer(cfg, agentInstance, mcpManager)
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	if err := mcpManager.StartServers(); err != nil {
		log.Fatalf("Failed to start MCP servers: %v", err)
	}

	if err := httpServer.Start(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	fmt.Printf("Agent Runtime server started on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Agent: %s (v%s)\n", cfg.Agent.Name, cfg.Agent.Version)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down...")

	if err := httpServer.Stop(); err != nil {
		log.Printf("Error stopping HTTP server: %v", err)
	}

	if err := mcpManager.StopServers(); err != nil {
		log.Printf("Error stopping MCP servers: %v", err)
	}

	if err := agentInstance.Stop(); err != nil {
		log.Printf("Error stopping agent: %v", err)
	}

	fmt.Println("Shutdown complete")
}
