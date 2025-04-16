package trust

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	trust := router.Group("/trust")
	{
		trust.GET("/grants/jwt-bearer/issuers", h.ListTrustedIssuers)
		trust.POST("/grants/jwt-bearer/issuers", h.CreateTrustedIssuer)
		trust.GET("/grants/jwt-bearer/issuers/:id", h.GetTrustedIssuer)
		trust.PUT("/grants/jwt-bearer/issuers/:id", h.UpdateTrustedIssuer)
		trust.DELETE("/grants/jwt-bearer/issuers/:id", h.DeleteTrustedIssuer)
	}
}

func (h *Handler) ListTrustedIssuers(c *gin.Context) {
	ctx := c.Request.Context()
	issuers, err := h.manager.ListTrustedIssuers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, issuers)
}

func (h *Handler) CreateTrustedIssuer(c *gin.Context) {
	var issuer TrustedIssuer
	if err := c.ShouldBindJSON(&issuer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.manager.CreateTrustedIssuer(ctx, &issuer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, issuer)
}

func (h *Handler) GetTrustedIssuer(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	issuer, err := h.manager.GetTrustedIssuer(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, issuer)
}

func (h *Handler) UpdateTrustedIssuer(c *gin.Context) {
	id := c.Param("id")
	var issuer TrustedIssuer
	if err := c.ShouldBindJSON(&issuer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	issuer.ID = id

	ctx := c.Request.Context()
	if err := h.manager.UpdateTrustedIssuer(ctx, &issuer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, issuer)
}

func (h *Handler) DeleteTrustedIssuer(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	if err := h.manager.DeleteTrustedIssuer(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
