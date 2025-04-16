package clair

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/security/clair"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config  *config.Config
	scanner *clair.Scanner
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
	}
}

func (m *Module) Name() string {
	return "clair"
}

func (m *Module) Description() string {
	return "Container vulnerability scanning"
}

func (m *Module) Initialize(ctx context.Context) error {
	scanner, err := clair.NewScanner(m.config)
	if err != nil {
		return err
	}

	m.scanner = scanner
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/clair")
	{
		group.POST("/scan", m.handleScan)
		group.GET("/scan/:id", m.handleGetScan)
		group.POST("/scan/async", m.handleScanAsync)
		group.GET("/scan/async/:id", m.handleGetScanAsync)
	}
}

func (m *Module) handleScan(c *gin.Context) {
	var req struct {
		ImageRef string `json:"image_ref" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := m.scanner.ScanImage(c.Request.Context(), req.ImageRef)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (m *Module) handleGetScan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	c.JSON(404, gin.H{"error": "scan not found"})
}

func (m *Module) handleScanAsync(c *gin.Context) {
	var req struct {
		ImageRef string `json:"image_ref" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	job, err := m.scanner.ScanImageAsync(c.Request.Context(), req.ImageRef)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(202, job)
}

func (m *Module) handleGetScanAsync(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	job, err := m.scanner.GetScanJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, job)
}

func (m *Module) GetScanner() *clair.Scanner {
	return m.scanner
}
