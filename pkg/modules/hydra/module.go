package hydra

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/client"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/jwk"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/oauth2"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/trust"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config         *config.Config
	oauth2Handler  *oauth2.Handler
	clientHandler  *client.Handler
	trustHandler   *trust.Handler
	jwkHandler     *jwk.Handler
	oauth2Provider fosite.OAuth2Provider
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
	}
}

func (m *Module) Name() string {
	return "hydra"
}

func (m *Module) Description() string {
	return "OAuth 2.0 and OpenID Connect server"
}

func (m *Module) Initialize(ctx context.Context) error {
	
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) RegisterRoutes(router *gin.Engine) {
	if m.oauth2Handler != nil {
		m.oauth2Handler.RegisterRoutes(router)
	}
	
	if m.clientHandler != nil {
		m.clientHandler.RegisterRoutes(router)
	}
	
	if m.trustHandler != nil {
		m.trustHandler.RegisterRoutes(router)
	}
	
	if m.jwkHandler != nil {
		m.jwkHandler.RegisterRoutes(router)
	}
}

func (m *Module) GetOAuth2Provider() fosite.OAuth2Provider {
	return m.oauth2Provider
}

func (m *Module) GetOAuth2Handler() *oauth2.Handler {
	return m.oauth2Handler
}

func (m *Module) GetClientHandler() *client.Handler {
	return m.clientHandler
}

func (m *Module) GetTrustHandler() *trust.Handler {
	return m.trustHandler
}

func (m *Module) GetJWKHandler() *jwk.Handler {
	return m.jwkHandler
}
