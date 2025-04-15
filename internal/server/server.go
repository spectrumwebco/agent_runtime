package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Server struct {
	router     *gin.Engine
	config     *config.Config
	mcp        *mcp.Manager
	agent      *agent.Agent
	grpcServer *GRPCServer
}

func New(cfg *config.Config) (*Server, error) {
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()
	
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	
	mcpManager, err := mcp.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP manager: %w", err)
	}
	
	agentInstance, err := agent.New(cfg, mcpManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}
	
	grpcServer, err := NewGRPCServer(cfg, agentInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC server: %w", err)
	}
	
	server := &Server{
		router:     router,
		config:     cfg,
		mcp:        mcpManager,
		agent:      agentInstance,
		grpcServer: grpcServer,
	}
	
	server.registerRoutes()
	
	return server, nil
}

func (s *Server) Start() error {
	if err := s.grpcServer.Start(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}
	
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %s\n", err)
			os.Exit(1)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	s.grpcServer.Stop()
	
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	
	return nil
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api")
	{
		api.GET("/health", s.healthCheck)
		
		agent := api.Group("/agent")
		{
			agent.POST("/execute", s.executeAgent)
			agent.GET("/status", s.getAgentStatus)
		}
		
		mcp := api.Group("/mcp")
		{
			mcp.GET("/servers", s.listMCPServers)
			mcp.GET("/servers/:name", s.getMCPServer)
		}
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) executeAgent(c *gin.Context) {
	var req struct {
		Task string `json:"task" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := s.agent.Execute(c.Request.Context(), req.Task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

func (s *Server) getAgentStatus(c *gin.Context) {
	status := s.agent.Status()
	
	c.JSON(http.StatusOK, status)
}

func (s *Server) listMCPServers(c *gin.Context) {
	servers := s.mcp.ListServers()
	
	c.JSON(http.StatusOK, servers)
}

func (s *Server) getMCPServer(c *gin.Context) {
	name := c.Param("name")
	
	server, err := s.mcp.GetServer(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, server)
}
