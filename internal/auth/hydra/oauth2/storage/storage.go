package storage

import (
	"context"
	"time"

	"github.com/ory/fosite"
)

type Storage interface {
	fosite.Storage
	ClientStorage
	AuthorizeCodeStorage
	AccessTokenStorage
	RefreshTokenStorage
	PKCEStorage
	OpenIDConnectStorage
}

type ClientStorage interface {
	GetClient(ctx context.Context, id string) (fosite.Client, error)
	
	CreateClient(ctx context.Context, client fosite.Client) error
	
	UpdateClient(ctx context.Context, client fosite.Client) error
	
	DeleteClient(ctx context.Context, id string) error
	
	ListClients(ctx context.Context) ([]fosite.Client, error)
}

type AuthorizeCodeStorage interface {
	CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) error
	
	GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (fosite.Requester, error)
	
	InvalidateAuthorizeCodeSession(ctx context.Context, code string) error
}

type AccessTokenStorage interface {
	CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) error
	
	GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
	
	DeleteAccessTokenSession(ctx context.Context, signature string) error
}

type RefreshTokenStorage interface {
	CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) error
	
	GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
	
	DeleteRefreshTokenSession(ctx context.Context, signature string) error
}

type PKCEStorage interface {
	CreatePKCERequestSession(ctx context.Context, signature string, request fosite.Requester) error
	
	GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
	
	DeletePKCERequestSession(ctx context.Context, signature string) error
}

type OpenIDConnectStorage interface {
	CreateOpenIDConnectSession(ctx context.Context, signature string, request fosite.Requester) error
	
	GetOpenIDConnectSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
	
	DeleteOpenIDConnectSession(ctx context.Context, signature string) error
}
