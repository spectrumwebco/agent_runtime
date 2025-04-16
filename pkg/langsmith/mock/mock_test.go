package mock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockLangSmithClient is a mock implementation of the LangSmith client
type MockLangSmithClient struct {
	Enabled bool
	Runs    map[string]map[string]interface{}
}

// NewMockClient creates a new mock LangSmith client
func NewMockClient() *MockLangSmithClient {
	return &MockLangSmithClient{
		Enabled: true,
		Runs:    make(map[string]map[string]interface{}),
	}
}

// CreateRun creates a mock run
func (m *MockLangSmithClient) CreateRun(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	runID := "run-" + time.Now().Format(time.RFC3339)
	m.Runs[runID] = req
	return map[string]interface{}{
		"id":     runID,
		"status": "running",
	}, nil
}

// UpdateRun updates a mock run
func (m *MockLangSmithClient) UpdateRun(ctx context.Context, runID string, req map[string]interface{}) (map[string]interface{}, error) {
	if run, ok := m.Runs[runID]; ok {
		for k, v := range req {
			run[k] = v
		}
		m.Runs[runID] = run
		return run, nil
	}
	return nil, nil
}

// TestMockLangSmithClient tests the mock LangSmith client
func TestMockLangSmithClient(t *testing.T) {
	client := NewMockClient()
	assert.NotNil(t, client)
	
	run, err := client.CreateRun(context.Background(), map[string]interface{}{
		"name":   "Test Run",
		"status": "running",
	})
	assert.NoError(t, err)
	assert.NotNil(t, run)
	assert.Contains(t, run, "id")
	
	runID := run["id"].(string)
	updatedRun, err := client.UpdateRun(context.Background(), runID, map[string]interface{}{
		"status": "completed",
	})
	assert.NoError(t, err)
	assert.NotNil(t, updatedRun)
	assert.Equal(t, "completed", updatedRun["status"])
}
