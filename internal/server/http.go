package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type HTTPServer struct {
	router     *gin.Engine
	server     *http.Server
	agent      *agent.Agent
	mcpManager *mcp.Manager
	config     *config.Config
}

func NewHTTPServer(cfg *config.Config, agent *agent.Agent, mcpManager *mcp.Manager) (*HTTPServer, error) {
	router := gin.Default()

	server := &HTTPServer{
		router:     router,
		agent:      agent,
		mcpManager: mcpManager,
		config:     cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Handler: router,
		},
	}

	server.setupRoutes()

	return server, nil
}

func (s *HTTPServer) setupRoutes() {
	api := s.router.Group("/api")
	{
		agentRoutes := api.Group("/agent")
		{
			agentRoutes.GET("/status", s.getAgentStatus)
			agentRoutes.POST("/start", s.startAgent)
			agentRoutes.POST("/stop", s.stopAgent)
			agentRoutes.POST("/execute", s.executeTask)
		}

		mcpRoutes := api.Group("/mcp")
		{
			mcpRoutes.GET("/servers", s.listMCPServers)
			mcpRoutes.GET("/servers/:name", s.getMCPServer)
			mcpRoutes.GET("/resources", s.listMCPResources)
			mcpRoutes.GET("/resources/:uri", s.getMCPResource)
			mcpRoutes.GET("/tools", s.listMCPTools)
			mcpRoutes.POST("/tools/:name/execute", s.executeMCPTool)
		}

		toolsRoutes := api.Group("/tools")
		{
			toolsRoutes.GET("/", s.listTools)
			toolsRoutes.GET("/:name", s.getTool)
			toolsRoutes.POST("/:name/execute", s.executeTool)
		}
	}

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"agent":  s.config.Agent.Name,
			"version": s.config.Agent.Version,
		})
	})
}

func (s *HTTPServer) Start() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	return nil
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) getAgentStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":  s.config.Agent.Name,
		"state": s.agent.GetState(),
	})
}

func (s *HTTPServer) startAgent(c *gin.Context) {
	if err := s.agent.Start(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func (s *HTTPServer) stopAgent(c *gin.Context) {
	if err := s.agent.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "stopped",
	})
}

func (s *HTTPServer) executeTask(c *gin.Context) {
	var request struct {
		Task string `json:"task" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	go func() {
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"status": "executing",
		"task":   request.Task,
	})
}

func (s *HTTPServer) listMCPServers(c *gin.Context) {
	servers := s.mcpManager.ListServers()

	c.JSON(http.StatusOK, gin.H{
		"servers": servers,
	})
}

func (s *HTTPServer) getMCPServer(c *gin.Context) {
	name := c.Param("name")

	server, err := s.mcpManager.GetServer(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, server)
}

func (s *HTTPServer) listMCPResources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"resources": []string{},
	})
}

func (s *HTTPServer) getMCPResource(c *gin.Context) {
	uri := c.Param("uri")

	c.JSON(http.StatusOK, gin.H{
		"uri":     uri,
		"content": "",
	})
}

func (s *HTTPServer) listMCPTools(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tools": []string{},
	})
}

func (s *HTTPServer) executeMCPTool(c *gin.Context) {
	name := c.Param("name")

	var request struct {
		Params map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":   name,
		"params": request.Params,
		"result": nil,
	})
}

func (s *HTTPServer) listTools(c *gin.Context) {
	tools := s.agent.ListTools()

	c.JSON(http.StatusOK, gin.H{
		"tools": tools,
	})
}

func (s *HTTPServer) getTool(c *gin.Context) {
	name := c.Param("name")

	tool, err := s.agent.GetTool(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tool": tool,
	})
}

func (s *HTTPServer) executeTool(c *gin.Context) {
	name := c.Param("name")

	var request struct {
		Params map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := s.agent.ExecuteTool(c.Request.Context(), name, request.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
