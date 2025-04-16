package coder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	config     *config.Config
	token      string
}

func NewClient(cfg *config.Config) (*Client, error) {
	baseURL, err := url.Parse(cfg.Coder.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Coder URL: %v", err)
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		config: cfg,
	}, nil
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	c.token = "placeholder-token"
	return nil
}

func (c *Client) CreateWorkspace(ctx context.Context, name, templateID string, params map[string]interface{}) (*Workspace, error) {
	return &Workspace{
		ID:         "ws-" + name,
		Name:       name,
		TemplateID: templateID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}, nil
}

func (c *Client) GetWorkspace(ctx context.Context, id string) (*Workspace, error) {
	return &Workspace{
		ID:         id,
		Name:       "example-workspace",
		TemplateID: "template-123",
		Status:     "running",
		CreatedAt:  time.Now().Add(-24 * time.Hour),
	}, nil
}

func (c *Client) ListWorkspaces(ctx context.Context) ([]*Workspace, error) {
	return []*Workspace{
		{
			ID:         "ws-123",
			Name:       "example-workspace-1",
			TemplateID: "template-123",
			Status:     "running",
			CreatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         "ws-456",
			Name:       "example-workspace-2",
			TemplateID: "template-456",
			Status:     "stopped",
			CreatedAt:  time.Now().Add(-48 * time.Hour),
		},
	}, nil
}

func (c *Client) DeleteWorkspace(ctx context.Context, id string) error {
	return nil
}

func (c *Client) StartWorkspace(ctx context.Context, id string) error {
	return nil
}

func (c *Client) StopWorkspace(ctx context.Context, id string) error {
	return nil
}

func (c *Client) GetTemplate(ctx context.Context, id string) (*Template, error) {
	return &Template{
		ID:          id,
		Name:        "example-template",
		Description: "An example template",
		CreatedAt:   time.Now().Add(-7 * 24 * time.Hour),
	}, nil
}

func (c *Client) ListTemplates(ctx context.Context) ([]*Template, error) {
	return []*Template{
		{
			ID:          "template-123",
			Name:        "example-template-1",
			Description: "An example template",
			CreatedAt:   time.Now().Add(-7 * 24 * time.Hour),
		},
		{
			ID:          "template-456",
			Name:        "example-template-2",
			Description: "Another example template",
			CreatedAt:   time.Now().Add(-14 * 24 * time.Hour),
		},
	}, nil
}

type Workspace struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	TemplateID string    `json:"template_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	URL        string    `json:"url,omitempty"`
}

type Template struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
