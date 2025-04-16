package coder

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/devenv/coder"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config      *config.Config
	client      *coder.Client
	provisioner *coder.Provisioner
	terraform   *coder.TerraformProvider
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
	}
}

func (m *Module) Name() string {
	return "coder"
}

func (m *Module) Description() string {
	return "Self-hosted cloud development environments"
}

func (m *Module) Initialize(ctx context.Context) error {
	client, err := coder.NewClient(m.config)
	if err != nil {
		return err
	}
	m.client = client

	provisioner, err := coder.NewProvisioner(m.config)
	if err != nil {
		return err
	}
	m.provisioner = provisioner

	m.terraform = coder.NewTerraformProvider(m.config)

	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/coder")
	{
		group.POST("/workspaces", m.handleCreateWorkspace)
		group.GET("/workspaces", m.handleListWorkspaces)
		group.GET("/workspaces/:id", m.handleGetWorkspace)
		group.DELETE("/workspaces/:id", m.handleDeleteWorkspace)
		group.POST("/workspaces/:id/start", m.handleStartWorkspace)
		group.POST("/workspaces/:id/stop", m.handleStopWorkspace)
		group.GET("/templates", m.handleListTemplates)
		group.GET("/templates/:id", m.handleGetTemplate)
		group.POST("/templates", m.handleCreateTemplate)
	}
}

func (m *Module) handleCreateWorkspace(c *gin.Context) {
	var req struct {
		Name       string                 `json:"name" binding:"required"`
		TemplateID string                 `json:"template_id" binding:"required"`
		Params     map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := m.provisioner.ProvisionWorkspace(c.Request.Context(), req.Name, req.TemplateID, req.Params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (m *Module) handleListWorkspaces(c *gin.Context) {
	workspaces, err := m.client.ListWorkspaces(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, workspaces)
}

func (m *Module) handleGetWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	workspace, err := m.client.GetWorkspace(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, workspace)
}

func (m *Module) handleDeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	result, err := m.provisioner.DestroyWorkspace(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (m *Module) handleStartWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	err := m.client.StartWorkspace(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "started"})
}

func (m *Module) handleStopWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	err := m.client.StopWorkspace(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "stopped"})
}

func (m *Module) handleListTemplates(c *gin.Context) {
	templates, err := m.client.ListTemplates(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, templates)
}

func (m *Module) handleGetTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	template, err := m.client.GetTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, template)
}

func (m *Module) handleCreateTemplate(c *gin.Context) {
	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Variables   map[string]interface{} `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	template, err := m.terraform.CreateTemplate(c.Request.Context(), req.Name, req.Description, req.Variables)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, template)
}

func (m *Module) GetClient() *coder.Client {
	return m.client
}

func (m *Module) GetProvisioner() *coder.Provisioner {
	return m.provisioner
}

func (m *Module) GetTerraformProvider() *coder.TerraformProvider {
	return m.terraform
}
