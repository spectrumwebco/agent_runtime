package example

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro/v2/registry"
	"github.com/spectrumwebco/agent_runtime/internal/database/models"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/service"
)

type AgentService struct {
	service *service.Service
}

type AgentRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Prompt      string            `json:"prompt"`
	Config      map[string]string `json:"config"`
	WorkspaceID uint              `json:"workspace_id"`
}

type AgentResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Prompt      string            `json:"prompt"`
	Config      map[string]string `json:"config"`
	WorkspaceID uint              `json:"workspace_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func NewAgentService(registryAddrs []string) (*AgentService, error) {
	reg := registry.NewRegistry(
		registry.Addrs(registryAddrs...),
		registry.Timeout(time.Second*5),
	)

	srv, err := service.NewService(
		service.WithName("kled.service.agent"),
		service.WithVersion("latest"),
		service.WithRegistry(reg),
		service.WithMetadata(map[string]string{
			"type": "kled.io",
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	agentService := &AgentService{
		service: srv,
	}

	if err := srv.Server().Handle(
		srv.Server().NewHandler(
			&AgentHandler{
				service: agentService,
			},
		),
	); err != nil {
		return nil, fmt.Errorf("failed to register handler: %w", err)
	}

	return agentService, nil
}

func (s *AgentService) Run() error {
	return s.service.Run()
}

type AgentHandler struct {
	service *AgentService
}

func (h *AgentHandler) Create(ctx context.Context, req *AgentRequest, rsp *AgentResponse) error {
	agent := &models.Agent{
		Name:        req.Name,
		Description: req.Description,
		Prompt:      req.Prompt,
		Config:      fmt.Sprintf("%v", req.Config),
		WorkspaceID: req.WorkspaceID,
	}


	rsp.ID = 1 // Mock ID
	rsp.Name = agent.Name
	rsp.Description = agent.Description
	rsp.Prompt = agent.Prompt
	rsp.Config = req.Config
	rsp.WorkspaceID = agent.WorkspaceID
	rsp.CreatedAt = time.Now()
	rsp.UpdatedAt = time.Now()

	return nil
}

func (h *AgentHandler) Get(ctx context.Context, req *struct{ ID uint }, rsp *AgentResponse) error {

	rsp.ID = req.ID
	rsp.Name = "Example Agent"
	rsp.Description = "An example agent"
	rsp.Prompt = "This is an example prompt"
	rsp.Config = map[string]string{"key": "value"}
	rsp.WorkspaceID = 1
	rsp.CreatedAt = time.Now()
	rsp.UpdatedAt = time.Now()

	return nil
}

func (h *AgentHandler) List(ctx context.Context, req *struct{}, rsp *[]*AgentResponse) error {

	*rsp = []*AgentResponse{
		{
			ID:          1,
			Name:        "Example Agent 1",
			Description: "An example agent",
			Prompt:      "This is an example prompt",
			Config:      map[string]string{"key": "value"},
			WorkspaceID: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "Example Agent 2",
			Description: "Another example agent",
			Prompt:      "This is another example prompt",
			Config:      map[string]string{"key": "value"},
			WorkspaceID: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return nil
}

func (h *AgentHandler) Update(ctx context.Context, req *AgentRequest, rsp *AgentResponse) error {

	rsp.ID = 1
	rsp.Name = req.Name
	rsp.Description = req.Description
	rsp.Prompt = req.Prompt
	rsp.Config = req.Config
	rsp.WorkspaceID = req.WorkspaceID
	rsp.CreatedAt = time.Now()
	rsp.UpdatedAt = time.Now()

	return nil
}

func (h *AgentHandler) Delete(ctx context.Context, req *struct{ ID uint }, rsp *struct{}) error {
	return nil
}
