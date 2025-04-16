package oauth2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/oauth2/storage"
)

type Handler struct {
	registry Registry
	provider fosite.OAuth2Provider
}

func NewHandler(registry Registry) *Handler {
	return &Handler{
		registry: registry,
		provider: registry.OAuth2Provider(),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	oauth2 := router.Group("/oauth2")
	{
		oauth2.POST("/token", h.TokenEndpoint)
		oauth2.GET("/auth", h.AuthEndpoint)
		oauth2.POST("/revoke", h.RevokeEndpoint)
		oauth2.POST("/introspect", h.IntrospectEndpoint)
	}
}

func (h *Handler) TokenEndpoint(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := h.provider.NewAccessRequest(ctx, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.provider.FinishAccessRequest(ctx, resp, c.Request, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) AuthEndpoint(c *gin.Context) {
	ctx := c.Request.Context()
	ar, err := h.provider.NewAuthorizeRequest(ctx, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.provider.FinishAuthorizeRequest(ctx, ar, c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) RevokeEndpoint(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.provider.NewRevocationRequest(ctx, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) IntrospectEndpoint(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := h.provider.IntrospectToken(ctx, c.Request, fosite.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
