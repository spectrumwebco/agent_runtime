package kled

import (
	"context"
	"testing"
)

func TestNewFramework(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	if framework == nil {
		t.Fatal("Framework is nil")
	}

	if framework.Config.Name != config.Name {
		t.Errorf("Expected framework name %s, got %s", config.Name, framework.Config.Name)
	}

	if framework.Config.Version != config.Version {
		t.Errorf("Expected framework version %s, got %s", config.Version, framework.Config.Version)
	}

	if framework.Config.Debug != config.Debug {
		t.Errorf("Expected framework debug %v, got %v", config.Debug, framework.Config.Debug)
	}
}

func TestFrameworkStartStop(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}


	if err := framework.Start(); err != nil {
		t.Fatalf("Failed to start framework: %v", err)
	}

	if err := framework.Stop(); err != nil {
		t.Fatalf("Failed to stop framework: %v", err)
	}
}

func TestFrameworkRegisterIntegration(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	mockIntegration := &mockIntegration{
		name: "mock",
	}

	err = framework.RegisterIntegration(mockIntegration)
	if err != nil {
		t.Fatalf("Failed to register integration: %v", err)
	}

	integration, err := framework.GetIntegration("mock")
	if err != nil {
		t.Fatalf("Failed to get integration: %v", err)
	}
	
	if integration != mockIntegration {
		t.Errorf("Expected integration to be %v, got %v", mockIntegration, integration)
	}
}

func TestFrameworkExecuteTool(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	err = framework.RegisterTool("test_tool", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "test_result", nil
	})
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	result, err := framework.ExecuteTool(context.Background(), "test_tool", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to execute tool: %v", err)
	}

	if result != "test_result" {
		t.Errorf("Expected result to be %v, got %v", "test_result", result)
	}
}

type mockIntegration struct {
	name string
}

func (m *mockIntegration) Name() string {
	return m.name
}

func (m *mockIntegration) Start() error {
	return nil
}

func (m *mockIntegration) Stop() error {
	return nil
}

func (m *mockIntegration) SaveState() ([]byte, error) {
	return []byte{}, nil
}

func (m *mockIntegration) LoadState(data []byte) error {
	return nil
}
