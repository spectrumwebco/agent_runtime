package client

import (
	"context"
	"errors"
	"time"

	"github.com/ory/fosite"
	"github.com/spectrumwebco/agent_runtime/internal/auth/hydra/oauth2/storage"
)

type Client struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Secret            string   `json:"secret,omitempty"`
	RedirectURIs      []string `json:"redirect_uris"`
	GrantTypes        []string `json:"grant_types"`
	ResponseTypes     []string `json:"response_types"`
	Scopes            []string `json:"scopes"`
	Audience          []string `json:"audience"`
	Public            bool     `json:"public"`
	TokenEndpointAuth string   `json:"token_endpoint_auth_method"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *Client) GetGrantTypes() []string {
	return c.GrantTypes
}

func (c *Client) GetResponseTypes() []string {
	return c.ResponseTypes
}

func (c *Client) GetScopes() []string {
	return c.Scopes
}

func (c *Client) GetAudience() []string {
	return c.Audience
}

func (c *Client) IsPublic() bool {
	return c.Public
}

func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

type Manager interface {
	GetClient(ctx context.Context, id string) (*Client, error)
	
	CreateClient(ctx context.Context, client *Client) error
	
	UpdateClient(ctx context.Context, client *Client) error
	
	DeleteClient(ctx context.Context, id string) error
	
	ListClients(ctx context.Context) ([]*Client, error)
}

type DefaultManager struct {
	storage storage.ClientStorage
}

func NewManager(storage storage.ClientStorage) Manager {
	return &DefaultManager{
		storage: storage,
	}
}

func (m *DefaultManager) GetClient(ctx context.Context, id string) (*Client, error) {
	fositeClient, err := m.storage.GetClient(ctx, id)
	if err != nil {
		return nil, err
	}
	
	client, ok := fositeClient.(*Client)
	if !ok {
		return nil, errors.New("invalid client type")
	}
	
	return client, nil
}

func (m *DefaultManager) CreateClient(ctx context.Context, client *Client) error {
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()
	
	return m.storage.CreateClient(ctx, client)
}

func (m *DefaultManager) UpdateClient(ctx context.Context, client *Client) error {
	client.UpdatedAt = time.Now()
	
	return m.storage.UpdateClient(ctx, client)
}

func (m *DefaultManager) DeleteClient(ctx context.Context, id string) error {
	return m.storage.DeleteClient(ctx, id)
}

func (m *DefaultManager) ListClients(ctx context.Context) ([]*Client, error) {
	fositeClients, err := m.storage.ListClients(ctx)
	if err != nil {
		return nil, err
	}
	
	clients := make([]*Client, 0, len(fositeClients))
	for _, fositeClient := range fositeClients {
		client, ok := fositeClient.(*Client)
		if !ok {
			return nil, errors.New("invalid client type")
		}
		clients = append(clients, client)
	}
	
	return clients, nil
}
