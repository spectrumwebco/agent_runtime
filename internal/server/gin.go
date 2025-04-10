package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
)

type GinServer struct {
	Agent *agent.Agent
	
	Router *gin.Engine
	
	Server *http.Server
	
	Addr string
}

func NewGinServer(agent *agent.Agent, addr string) (*GinServer, error) {
	router := gin.Default()
	
	server := &GinServer{
		Agent:  agent,
		Router: router,
		Addr:   addr,
		Server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
	
	server.setupRoutes()
	
	return server, nil
}

func (s *GinServer) setupRoutes() {
	api := s.Router.Group("/api")
	{
		agent := api.Group("/agent")
		{
			agent.GET("/status", s.getAgentStatus)
			agent.POST("/start", s.startAgent)
			agent.POST("/stop", s.stopAgent)
			agent.POST("/execute", s.executeTask)
		}
		
		tools := api.Group("/tools")
		{
			tools.GET("/", s.listTools)
			tools.GET("/:name", s.getTool)
			tools.POST("/:name/execute", s.executeTool)
		}
		
		mcp := api.Group("/mcp")
		{
			mcp.GET("/servers", s.listMCPServers)
			mcp.GET("/resources", s.listMCPResources)
			mcp.GET("/resources/:uri", s.getMCPResource)
			mcp.GET("/tools", s.listMCPTools)
			mcp.POST("/tools/:name/execute", s.executeMCPTool)
		}
	}
}

func (s *GinServer) Start() error {
	return s.Server.ListenAndServe()
}

func (s *GinServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}

func (s *GinServer) getAgentStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":  s.Agent.Config.Name,
		"state": s.Agent.GetState(),
	})
}

func (s *GinServer) startAgent(c *gin.Context) {
	if err := s.Agent.Start(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func (s *GinServer) stopAgent(c *gin.Context) {
	if err := s.Agent.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "stopped",
	})
}

func (s *GinServer) executeTask(c *gin.Context) {
	var request struct {
		Task string `json:"task" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	
	c.JSON(http.StatusOK, gin.H{
		"status": "executing",
		"task":   request.Task,
	})
}

func (s *GinServer) listTools(c *gin.Context) {
	tools := s.Agent.ListTools()
	
	c.JSON(http.StatusOK, gin.H{
		"tools": tools,
	})
}

func (s *GinServer) getTool(c *gin.Context) {
	name := c.Param("name")
	
	tool, err := s.Agent.GetTool(name)
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

func (s *GinServer) executeTool(c *gin.Context) {
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
	
	result, err := s.Agent.ExecuteTool(c.Request.Context(), name, request.Params)
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

func (s *GinServer) listMCPServers(c *gin.Context) {
	
	c.JSON(http.StatusOK, gin.H{
		"servers": []string{},
	})
}

func (s *GinServer) listMCPResources(c *gin.Context) {
	
	c.JSON(http.StatusOK, gin.H{
		"resources": []string{},
	})
}

func (s *GinServer) getMCPResource(c *gin.Context) {
	uri := c.Param("uri")
	
	
	c.JSON(http.StatusOK, gin.H{
		"uri":     uri,
		"content": "",
	})
}

func (s *GinServer) listMCPTools(c *gin.Context) {
	
	c.JSON(http.StatusOK, gin.H{
		"tools": []string{},
	})
}

func (s *GinServer) executeMCPTool(c *gin.Context) {
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
