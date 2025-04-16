package gorm

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/database/gorm"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config   *config.Config
	database *gorm.Database
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
	}
}

func (m *Module) Initialize() error {
	dbConfig := &gorm.Config{
		Type:         m.config.GetString("database.type"),
		Host:         m.config.GetString("database.host"),
		Port:         m.config.GetInt("database.port"),
		Username:     m.config.GetString("database.username"),
		Password:     m.config.GetString("database.password"),
		Database:     m.config.GetString("database.name"),
		SSLMode:      m.config.GetString("database.ssl_mode"),
		MaxOpenConns: m.config.GetInt("database.max_open_conns"),
		MaxIdleConns: m.config.GetInt("database.max_idle_conns"),
		TablePrefix:  m.config.GetString("database.table_prefix"),
		Debug:        m.config.GetBool("database.debug"),
	}

	if dbConfig.Type == "" {
		dbConfig.Type = "sqlite"
	}

	if dbConfig.Database == "" && dbConfig.Type == "sqlite" {
		dbConfig.Database = "kled.db"
	}

	if dbConfig.MaxOpenConns == 0 {
		dbConfig.MaxOpenConns = 10
	}

	if dbConfig.MaxIdleConns == 0 {
		dbConfig.MaxIdleConns = 5
	}

	db, err := gorm.NewDatabase(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	m.database = db

	return nil
}

func (m *Module) Start(ctx context.Context) error {
	err := m.database.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	err = m.database.Migrate(gorm.GetModels()...)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func (m *Module) Stop(ctx context.Context) error {
	err := m.database.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

func (m *Module) GetDatabase() *gorm.Database {
	return m.database
}

func (m *Module) RunMigrations() error {
	return m.database.Migrate(gorm.GetModels()...)
}

func (m *Module) CreateUser(username, email, passwordHash, firstName, lastName string) (*gorm.User, error) {
	user := &gorm.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
	}

	err := m.database.DB.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (m *Module) GetUserByID(id uint) (*gorm.User, error) {
	var user gorm.User
	err := m.database.DB.First(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (m *Module) GetUserByUsername(username string) (*gorm.User, error) {
	var user gorm.User
	err := m.database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (m *Module) GetUserByEmail(email string) (*gorm.User, error) {
	var user gorm.User
	err := m.database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (m *Module) UpdateUser(user *gorm.User) error {
	err := m.database.DB.Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (m *Module) DeleteUser(id uint) error {
	err := m.database.DB.Delete(&gorm.User{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (m *Module) CreateAgent(name, description, agentType, config string, userID uint) (*gorm.Agent, error) {
	agent := &gorm.Agent{
		Name:        name,
		Description: description,
		Type:        agentType,
		Config:      config,
		UserID:      userID,
		IsActive:    true,
	}

	err := m.database.DB.Create(agent).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agent, nil
}

func (m *Module) GetAgentByID(id uint) (*gorm.Agent, error) {
	var agent gorm.Agent
	err := m.database.DB.First(&agent, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	return &agent, nil
}

func (m *Module) GetAgentsByUserID(userID uint) ([]gorm.Agent, error) {
	var agents []gorm.Agent
	err := m.database.DB.Where("user_id = ?", userID).Find(&agents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	return agents, nil
}

func (m *Module) UpdateAgent(agent *gorm.Agent) error {
	err := m.database.DB.Save(agent).Error
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return nil
}

func (m *Module) DeleteAgent(id uint) error {
	err := m.database.DB.Delete(&gorm.Agent{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	return nil
}

func (m *Module) CreateTask(title, description, status, priority string, userID, agentID uint) (*gorm.Task, error) {
	task := &gorm.Task{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		UserID:      userID,
		AgentID:     agentID,
	}

	err := m.database.DB.Create(task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (m *Module) GetTaskByID(id uint) (*gorm.Task, error) {
	var task gorm.Task
	err := m.database.DB.First(&task, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (m *Module) GetTasksByUserID(userID uint) ([]gorm.Task, error) {
	var tasks []gorm.Task
	err := m.database.DB.Where("user_id = ?", userID).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (m *Module) GetTasksByAgentID(agentID uint) ([]gorm.Task, error) {
	var tasks []gorm.Task
	err := m.database.DB.Where("agent_id = ?", agentID).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (m *Module) UpdateTask(task *gorm.Task) error {
	err := m.database.DB.Save(task).Error
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (m *Module) DeleteTask(id uint) error {
	err := m.database.DB.Delete(&gorm.Task{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (m *Module) RunTransaction(fn func(tx *gorm.Database) error) error {
	return m.database.Transaction(func(tx *gorm.Database) error {
		return fn(tx)
	})
}
