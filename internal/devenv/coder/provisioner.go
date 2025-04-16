package coder

import (
	"context"
	"fmt"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Provisioner struct {
	client *Client
	config *config.Config
}

func NewProvisioner(cfg *config.Config) (*Provisioner, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Coder client: %v", err)
	}

	return &Provisioner{
		client: client,
		config: cfg,
	}, nil
}

func (p *Provisioner) ProvisionWorkspace(ctx context.Context, name, templateID string, params map[string]interface{}) (*ProvisionResult, error) {
	workspace, err := p.client.CreateWorkspace(ctx, name, templateID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %v", err)
	}

	result := &ProvisionResult{
		WorkspaceID: workspace.ID,
		Status:      "provisioning",
		StartedAt:   time.Now(),
	}

	result.Status = "provisioned"
	result.CompletedAt = time.Now()

	return result, nil
}

func (p *Provisioner) DestroyWorkspace(ctx context.Context, id string) (*DestroyResult, error) {
	err := p.client.DeleteWorkspace(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete workspace: %v", err)
	}

	result := &DestroyResult{
		WorkspaceID: id,
		Status:      "destroying",
		StartedAt:   time.Now(),
	}

	result.Status = "destroyed"
	result.CompletedAt = time.Now()

	return result, nil
}

func (p *Provisioner) GetProvisionStatus(ctx context.Context, id string) (*ProvisionResult, error) {
	workspace, err := p.client.GetWorkspace(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %v", err)
	}

	result := &ProvisionResult{
		WorkspaceID: workspace.ID,
		Status:      workspace.Status,
	}

	if workspace.Status == "running" {
		result.Status = "provisioned"
		result.CompletedAt = workspace.UpdatedAt
	}

	return result, nil
}

type ProvisionResult struct {
	WorkspaceID string    `json:"workspace_id"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

type DestroyResult struct {
	WorkspaceID string    `json:"workspace_id"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}
