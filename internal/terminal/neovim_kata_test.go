package terminal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKataManager struct {
	mock.Mock
}

func (m *MockKataManager) CreateContainer(ctx context.Context, name string, options map[string]interface{}) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

func (m *MockKataManager) StopContainer(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockKataManager) RemoveContainer(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockKataManager) ExecInContainer(ctx context.Context, name string, command []string) (string, error) {
	args := m.Called(ctx, name, command)
	return args.String(0), args.Error(1)
}

func TestNeovimKataTerminal_Start(t *testing.T) {
	mockKataManager := new(MockKataManager)
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()
	
	options := map[string]interface{}{
		"cpus":   2,
		"memory": 2048,
		"debug":  false,
	}
	terminal := NewNeovimKataTerminal("test-id", server.URL, options, mockKataManager)
	
	mockKataManager.On("CreateContainer", mock.Anything, "neovim-test-id", mock.Anything).Return(nil)
	
	err := terminal.Start(context.Background())
	
	assert.NoError(t, err)
	assert.True(t, terminal.IsRunning())
	mockKataManager.AssertExpectations(t)
}

func TestNeovimKataTerminal_Stop(t *testing.T) {
	mockKataManager := new(MockKataManager)
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()
	
	options := map[string]interface{}{
		"cpus":   2,
		"memory": 2048,
		"debug":  false,
	}
	terminal := NewNeovimKataTerminal("test-id", server.URL, options, mockKataManager)
	
	terminal.running = true
	
	mockKataManager.On("StopContainer", mock.Anything, "neovim-test-id").Return(nil)
	mockKataManager.On("RemoveContainer", mock.Anything, "neovim-test-id").Return(nil)
	
	err := terminal.Stop(context.Background())
	
	assert.NoError(t, err)
	assert.False(t, terminal.IsRunning())
	mockKataManager.AssertExpectations(t)
}

func TestNeovimKataTerminal_Execute(t *testing.T) {
	mockKataManager := new(MockKataManager)
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "output": "command output"}`))
	}))
	defer server.Close()
	
	options := map[string]interface{}{
		"cpus":   2,
		"memory": 2048,
		"debug":  false,
	}
	terminal := NewNeovimKataTerminal("test-id", server.URL, options, mockKataManager)
	
	terminal.running = true
	
	mockKataManager.On("ExecInContainer", mock.Anything, "neovim-test-id", []string{"nvim", "--headless", "-c", "echo 'ls -la'"}).Return("", nil)
	
	output, err := terminal.Execute(context.Background(), "ls -la")
	
	assert.NoError(t, err)
	assert.Equal(t, "command output", output)
	mockKataManager.AssertExpectations(t)
}

func TestNeovimKataTerminal_ID(t *testing.T) {
	mockKataManager := new(MockKataManager)
	
	terminal := NewNeovimKataTerminal("test-id", "http://localhost:8080", nil, mockKataManager)
	
	id := terminal.ID()
	
	assert.Equal(t, "test-id", id)
}

func TestNeovimKataTerminal_GetType(t *testing.T) {
	mockKataManager := new(MockKataManager)
	
	terminal := NewNeovimKataTerminal("test-id", "http://localhost:8080", nil, mockKataManager)
	
	terminalType := terminal.GetType()
	
	assert.Equal(t, "neovim-kata", terminalType)
}

func TestNeovimKataTerminal_IsRunning(t *testing.T) {
	mockKataManager := new(MockKataManager)
	
	terminal := NewNeovimKataTerminal("test-id", "http://localhost:8080", nil, mockKataManager)
	
	running := terminal.IsRunning()
	
	assert.False(t, running)
	
	terminal.running = true
	
	running = terminal.IsRunning()
	
	assert.True(t, running)
}
