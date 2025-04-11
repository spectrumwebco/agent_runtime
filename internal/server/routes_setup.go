package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
)

func SetupRoutes(router *gin.Engine) {
	agentInstance, _ := agent.NewDefaultAgent()
	
	configuredRouter := SetupRouter(agentInstance)
	
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"version": "1.0.0",
		})
	})
	
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "Agent Runtime",
			"description": "Go implementation of SWE-Agent and SWE-ReX frameworks",
			"version": "1.0.0",
		})
	})
}
