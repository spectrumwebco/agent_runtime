package framework

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/framework"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config    *config.Config
	framework *framework.Framework
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
	}
}

func (m *Module) Initialize() error {
	opts := framework.FrameworkOptions{
		ConfigPath: m.config.GetString("framework.config_path"),
		ModelPath:  m.config.GetString("framework.model_path"),
		PolicyPath: m.config.GetString("framework.policy_path"),
		ModelText:  m.config.GetString("framework.model_text"),
	}

	fw, err := framework.NewFramework(opts)
	if err != nil {
		return fmt.Errorf("failed to create framework: %w", err)
	}

	m.framework = fw

	err = m.framework.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize framework: %w", err)
	}

	return nil
}

func (m *Module) Start(ctx context.Context) error {
	return m.framework.Start(ctx)
}

func (m *Module) Stop(ctx context.Context) error {
	return m.framework.Stop(ctx)
}

func (m *Module) GetFramework() *framework.Framework {
	return m.framework
}

func (m *Module) RunMigrations() error {
	return m.framework.RunMigrations()
}

func (m *Module) CreateService(name, version, address string, metadata map[string]string) error {
	_, err := m.framework.CreateService(name, version, address, metadata)
	return err
}

func (m *Module) CreateGraph(name string) error {
	_, err := m.framework.CreateGraph(name)
	return err
}

func (m *Module) CreateActor(id string, behavior interface{}, state map[string]interface{}) error {
	_, err := m.framework.CreateActor(id, behavior, state)
	return err
}

func (m *Module) Enforce(params ...interface{}) (bool, error) {
	return m.framework.Enforce(params...)
}

func (m *Module) WithPythonAgentIntegration(apiURL string) error {
	
	
	return nil
}

func (m *Module) WithSharedState() error {
	
	
	return nil
}
