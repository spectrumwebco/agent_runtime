package jwk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3"
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
	jwks := router.Group("/jwks")
	{
		jwks.GET("", h.ListSets)
		jwks.POST("", h.CreateSet)
		jwks.GET("/:set", h.GetSet)
		jwks.PUT("/:set", h.UpdateSet)
		jwks.DELETE("/:set", h.DeleteSet)
		jwks.GET("/:set/:kid", h.GetKey)
		jwks.PUT("/:set/:kid", h.UpdateKey)
		jwks.DELETE("/:set/:kid", h.DeleteKey)
	}
}

func (h *Handler) ListSets(c *gin.Context) {
	ctx := c.Request.Context()
	sets, err := h.manager.ListSets(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sets)
}

func (h *Handler) CreateSet(c *gin.Context) {
	var req CreateSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	set, err := h.manager.CreateSet(ctx, req.SetID, req.KeyType, req.Use, req.Algorithm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, set)
}

func (h *Handler) GetSet(c *gin.Context) {
	setID := c.Param("set")
	ctx := c.Request.Context()
	set, err := h.manager.GetSet(ctx, setID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}

func (h *Handler) UpdateSet(c *gin.Context) {
	setID := c.Param("set")
	var req UpdateSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	set, err := h.manager.UpdateSet(ctx, setID, req.Keys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}

func (h *Handler) DeleteSet(c *gin.Context) {
	setID := c.Param("set")
	ctx := c.Request.Context()
	if err := h.manager.DeleteSet(ctx, setID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetKey(c *gin.Context) {
	setID := c.Param("set")
	keyID := c.Param("kid")
	ctx := c.Request.Context()
	key, err := h.manager.GetKey(ctx, setID, keyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, key)
}

func (h *Handler) UpdateKey(c *gin.Context) {
	setID := c.Param("set")
	keyID := c.Param("kid")
	var key jose.JSONWebKey
	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	updatedKey, err := h.manager.UpdateKey(ctx, setID, keyID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedKey)
}

func (h *Handler) DeleteKey(c *gin.Context) {
	setID := c.Param("set")
	keyID := c.Param("kid")
	ctx := c.Request.Context()
	if err := h.manager.DeleteKey(ctx, setID, keyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type CreateSetRequest struct {
	SetID     string `json:"set_id"`
	KeyType   string `json:"key_type"`
	Use       string `json:"use"`
	Algorithm string `json:"algorithm"`
}

type UpdateSetRequest struct {
	Keys []jose.JSONWebKey `json:"keys"`
}
