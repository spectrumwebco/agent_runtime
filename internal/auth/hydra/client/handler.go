package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
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
	clients := router.Group("/clients")
	{
		clients.GET("", h.ListClients)
		clients.POST("", h.CreateClient)
		clients.GET("/:id", h.GetClient)
		clients.PUT("/:id", h.UpdateClient)
		clients.DELETE("/:id", h.DeleteClient)
	}
}

func (h *Handler) ListClients(c *gin.Context) {
	ctx := c.Request.Context()
	clients, err := h.manager.ListClients(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clients)
}

func (h *Handler) CreateClient(c *gin.Context) {
	var client Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.manager.CreateClient(ctx, &client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, client)
}

func (h *Handler) GetClient(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	client, err := h.manager.GetClient(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, client)
}

func (h *Handler) UpdateClient(c *gin.Context) {
	id := c.Param("id")
	var client Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client.ID = id

	ctx := c.Request.Context()
	if err := h.manager.UpdateClient(ctx, &client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, client)
}

func (h *Handler) DeleteClient(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	if err := h.manager.DeleteClient(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
