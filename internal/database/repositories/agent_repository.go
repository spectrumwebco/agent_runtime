package repositories

import (
	"errors"

	"github.com/spectrumwebco/agent_runtime/internal/database"
	"github.com/spectrumwebco/agent_runtime/internal/database/models"
	"gorm.io/gorm"
)

type AgentRepository struct {
	db *database.DB
}

func NewAgentRepository(db *database.DB) *AgentRepository {
	return &AgentRepository{
		db: db,
	}
}

func (r *AgentRepository) Create(agent *models.Agent) error {
	return r.db.Create(agent).Error
}

func (r *AgentRepository) GetByID(id uint) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.First(&agent, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepository) GetByWorkspaceID(workspaceID uint) ([]models.Agent, error) {
	var agents []models.Agent
	err := r.db.Where("workspace_id = ?", workspaceID).Find(&agents).Error
	return agents, err
}

func (r *AgentRepository) List(offset, limit int) ([]models.Agent, error) {
	var agents []models.Agent
	err := r.db.Offset(offset).Limit(limit).Find(&agents).Error
	return agents, err
}

func (r *AgentRepository) Update(agent *models.Agent) error {
	return r.db.Save(agent).Error
}

func (r *AgentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Agent{}, id).Error
}

func (r *AgentRepository) AddToolToAgent(agentID, toolID uint) error {
	return r.db.Exec("INSERT INTO agent_tools (agent_id, tool_id) VALUES (?, ?)", agentID, toolID).Error
}

func (r *AgentRepository) RemoveToolFromAgent(agentID, toolID uint) error {
	return r.db.Exec("DELETE FROM agent_tools WHERE agent_id = ? AND tool_id = ?", agentID, toolID).Error
}

func (r *AgentRepository) GetAgentTools(agentID uint) ([]models.Tool, error) {
	var agent models.Agent
	err := r.db.Preload("Tools").First(&agent, agentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return agent.Tools, nil
}

func (r *AgentRepository) CreateExecution(execution *models.Execution) error {
	return r.db.Create(execution).Error
}

func (r *AgentRepository) GetExecutionByID(id uint) (*models.Execution, error) {
	var execution models.Execution
	err := r.db.First(&execution, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &execution, nil
}

func (r *AgentRepository) GetAgentExecutions(agentID uint) ([]models.Execution, error) {
	var executions []models.Execution
	err := r.db.Where("agent_id = ?", agentID).Find(&executions).Error
	return executions, err
}

func (r *AgentRepository) CreateExecutionStep(step *models.ExecutionStep) error {
	return r.db.Create(step).Error
}

func (r *AgentRepository) GetExecutionSteps(executionID uint) ([]models.ExecutionStep, error) {
	var steps []models.ExecutionStep
	err := r.db.Where("execution_id = ?", executionID).Order("order asc").Find(&steps).Error
	return steps, err
}
