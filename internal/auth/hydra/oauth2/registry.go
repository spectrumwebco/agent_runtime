package oauth2

import (
	"github.com/ory/fosite"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/oauth2/storage"
)

type Registry interface {
	OAuth2Provider() fosite.OAuth2Provider
	OAuth2Storage() storage.Storage
	Config() *Config
}

type Config struct {
	Issuer             string
	AccessTokenLifespan int
	AuthCodeLifespan   int
	IDTokenLifespan    int
	SecretKey          []byte
}
