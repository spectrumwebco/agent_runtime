package djangobridge

import (
	"context"
	"testing"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/stretchr/testify/assert"
)

type MockEventStream struct {
	Events []*models.Event
}

func NewMockEventStream() *MockEventStream {
	return &MockEventStream{
		Events: []*models.Event{},
	}
}

func (m *MockEventStream) AddEvent(event *models.Event) error {
	m.Events = append(m.Events, event)
	return nil
}

func TestDjangoBridge(t *testing.T) {
	t.Skip("Skipping test that requires a running Django gRPC server")

	eventStream := NewMockEventStream()

	bridge, err := NewDjangoBridge("localhost:50051", eventStream)
	assert.NoError(t, err, "Should create Django bridge without error")
	assert.NotNil(t, bridge, "Bridge should not be nil")

	defer bridge.Close()

	ctx := context.Background()
	result, err := bridge.ExecutePythonCode(ctx, "print('Hello, World!')", 5*time.Second)
	assert.NoError(t, err, "Should execute Python code without error")
	assert.NotNil(t, result, "Result should not be nil")

	models, err := bridge.QueryDjangoModel(ctx, "User", map[string]interface{}{
		"is_active": true,
	})
	assert.NoError(t, err, "Should query Django model without error")
	assert.NotNil(t, models, "Models should not be nil")

	createdModel, err := bridge.CreateDjangoModel(ctx, "User", map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
	})
	assert.NoError(t, err, "Should create Django model without error")
	assert.NotNil(t, createdModel, "Created model should not be nil")

	err = bridge.UpdateDjangoModel(ctx, "User", 1, map[string]interface{}{
		"email": "updated@example.com",
	})
	assert.NoError(t, err, "Should update Django model without error")

	err = bridge.DeleteDjangoModel(ctx, "User", 1)
	assert.NoError(t, err, "Should delete Django model without error")

	output, err := bridge.ExecuteDjangoManagementCommand(ctx, "check", []string{"--database=default"})
	assert.NoError(t, err, "Should execute Django management command without error")
	assert.NotEmpty(t, output, "Output should not be empty")

	settings, err := bridge.GetDjangoSettings(ctx)
	assert.NoError(t, err, "Should get Django settings without error")
	assert.NotNil(t, settings, "Settings should not be nil")

	agent, err := bridge.CreateDjangoAgent(ctx, "test-agent", "Test Agent", "test")
	assert.NoError(t, err, "Should create Django agent without error")
	assert.NotNil(t, agent, "Agent should not be nil")

	assert.Greater(t, len(eventStream.Events), 0, "Events should have been emitted")
}

func TestGRPCClient(t *testing.T) {
	t.Skip("Skipping test that requires a running Django gRPC server")

	client, err := NewGRPCClient("localhost:50051")
	assert.NoError(t, err, "Should create gRPC client without error")
	assert.NotNil(t, client, "Client should not be nil")

	defer client.Close()

	ctx := context.Background()
	result, err := client.ExecuteAgentTask(ctx, "test_task", map[string]interface{}{
		"param1": "value1",
		"param2": "value2",
	}, 5*time.Second)
	assert.NoError(t, err, "Should execute agent task without error")
	assert.NotNil(t, result, "Result should not be nil")

	models, err := client.QueryModel(ctx, "User", map[string]interface{}{
		"is_active": true,
	}, 5*time.Second)
	assert.NoError(t, err, "Should query model without error")
	assert.NotNil(t, models, "Models should not be nil")

	createdModel, err := client.CreateModel(ctx, "User", map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
	}, 5*time.Second)
	assert.NoError(t, err, "Should create model without error")
	assert.NotNil(t, createdModel, "Created model should not be nil")

	err = client.UpdateModel(ctx, "User", "1", map[string]interface{}{
		"email": "updated@example.com",
	}, 5*time.Second)
	assert.NoError(t, err, "Should update model without error")

	err = client.DeleteModel(ctx, "User", "1", 5*time.Second)
	assert.NoError(t, err, "Should delete model without error")

	output, err := client.ExecuteManagementCommand(ctx, "check", []string{"--database=default"}, 5*time.Second)
	assert.NoError(t, err, "Should execute management command without error")
	assert.NotEmpty(t, output, "Output should not be empty")

	settings, err := client.GetSettings(ctx, 5*time.Second)
	assert.NoError(t, err, "Should get settings without error")
	assert.NotNil(t, settings, "Settings should not be nil")

	pythonResult, err := client.ExecutePythonCode(ctx, "print('Hello, World!')", 5*time.Second)
	assert.NoError(t, err, "Should execute Python code without error")
	assert.NotNil(t, pythonResult, "Python result should not be nil")

	err = client.Debug(ctx, "Debug message", map[string]interface{}{
		"key": "value",
	})
	assert.NoError(t, err, "Should log debug message without error")

	err = client.Info(ctx, "Info message", map[string]interface{}{
		"key": "value",
	})
	assert.NoError(t, err, "Should log info message without error")

	err = client.Warning(ctx, "Warning message", map[string]interface{}{
		"key": "value",
	})
	assert.NoError(t, err, "Should log warning message without error")

	err = client.Error(ctx, "Error message", map[string]interface{}{
		"key": "value",
	})
	assert.NoError(t, err, "Should log error message without error")

	err = client.Critical(ctx, "Critical message", map[string]interface{}{
		"key": "value",
	})
	assert.NoError(t, err, "Should log critical message without error")
}
