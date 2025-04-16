package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/server/middleware"
	"github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

// SetupAPIRoutes sets up the API routes
func SetupAPIRoutes(router *gin.Engine, enforcer *auth.Enforcer, authToken string) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	
	// Public routes
	v1.GET("/health", HealthCheck)
	v1.GET("/version", GetVersion)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(authToken))
	protected.Use(middleware.AuthorizeMiddleware(enforcer))

	// Agent routes
	agents := protected.Group("/agents")
	agents.GET("", ListAgents)
	agents.POST("", CreateAgent)
	agents.GET("/:id", GetAgent)
	agents.PUT("/:id", UpdateAgent)
	agents.DELETE("/:id", DeleteAgent)
	agents.POST("/:id/execute", ExecuteAgent)

	// Tool routes
	tools := protected.Group("/tools")
	tools.GET("", ListTools)
	tools.POST("", CreateTool)
	tools.GET("/:id", GetTool)
	tools.PUT("/:id", UpdateTool)
	tools.DELETE("/:id", DeleteTool)
	tools.POST("/:id/execute", ExecuteTool)

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authToken))
	admin.Use(middleware.AdminMiddleware())

	admin.GET("/users", ListUsers)
	admin.POST("/users", CreateUser)
	admin.GET("/users/:id", GetUser)
	admin.PUT("/users/:id", UpdateUser)
	admin.DELETE("/users/:id", DeleteUser)
}

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// GetVersion handles the version endpoint
func GetVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": "0.1.0",
	})
}

// ListAgents handles the list agents endpoint
func ListAgents(c *gin.Context) {
	c.JSON(200, gin.H{
		"agents": []string{},
	})
}

// CreateAgent handles the create agent endpoint
func CreateAgent(c *gin.Context) {
	c.JSON(201, gin.H{
		"id": "1",
	})
}

// GetAgent handles the get agent endpoint
func GetAgent(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
	})
}

// UpdateAgent handles the update agent endpoint
func UpdateAgent(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
	})
}

// DeleteAgent handles the delete agent endpoint
func DeleteAgent(c *gin.Context) {
	c.JSON(204, nil)
}

// ExecuteAgent handles the execute agent endpoint
func ExecuteAgent(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
		"result": "success",
	})
}

// ListTools handles the list tools endpoint
func ListTools(c *gin.Context) {
	c.JSON(200, gin.H{
		"tools": []string{},
	})
}

// CreateTool handles the create tool endpoint
func CreateTool(c *gin.Context) {
	c.JSON(201, gin.H{
		"id": "1",
	})
}

// GetTool handles the get tool endpoint
func GetTool(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
	})
}

// UpdateTool handles the update tool endpoint
func UpdateTool(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
	})
}

// DeleteTool handles the delete tool endpoint
func DeleteTool(c *gin.Context) {
	c.JSON(204, nil)
}

// ExecuteTool handles the execute tool endpoint
func ExecuteTool(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
		"result": "success",
	})
}

// ListUsers handles the list users endpoint
func ListUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"users": []string{},
	})
}

// CreateUser handles the create user endpoint
func CreateUser(c *gin.Context) {
	c.JSON(201, gin.H{
		"id": "1",
	})
}

// GetUser handles the get user endpoint
func GetUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
	})
}

// UpdateUser handles the update user endpoint
func UpdateUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": c.Param("id"),
	})
}

// DeleteUser handles the delete user endpoint
func DeleteUser(c *gin.Context) {
	c.JSON(204, nil)
}
