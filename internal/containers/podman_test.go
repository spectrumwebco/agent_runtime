package containers

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	execPkg "github.com/spectrumwebco/agent_runtime/internal/exec"
)

type MockExecer struct {
	mock.Mock
}

func (m *MockExecer) CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.Command("echo", "dummy")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	m.Called(ctx, name, args)
	return cmd
}

var _ execPkg.Execer = (*MockExecer)(nil)

type MockCmd struct {
	*exec.Cmd
	mock.Mock
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func NewMockCmd() *MockCmd {
	cmd := exec.Command("echo", "dummy")
	return &MockCmd{
		Cmd: cmd,
	}
}

func (m *MockCmd) CombinedOutput() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCmd) Output() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCmd) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCmd) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCmd) Wait() error {
	args := m.Called()
	return args.Error(0)
}

func TestPodman_List(t *testing.T) {
	mockExecer := new(MockExecer)
	
	mockCmd := NewMockCmd()
	
	mockExecer.On("CommandContext", mock.Anything, "podman", mock.Anything).Return(mockCmd)
	
	sampleOutput := `[
		{
			"Id": "container1",
			"Name": "container1",
			"Created": "2021-03-30T12:34:56Z",
			"Config": {
				"Image": "ubuntu:20.04",
				"Labels": {"app": "test"}
			},
			"State": {
				"Running": true,
				"Status": "running"
			},
			"Mounts": [
				{
					"Source": "/host/path",
					"Destination": "/container/path"
				}
			],
			"NetworkSettings": {
				"Ports": {
					"80/tcp": {
						"ContainerPort": "80",
						"Protocol": "tcp",
						"Description": "HTTP",
						"HostBindings": [
							{
								"HostIP": "0.0.0.0",
								"HostPort": "8080"
							}
						]
					}
				}
			}
		}
	]`
	mockCmd.On("Run").Return(nil)
	mockCmd.On("CombinedOutput").Return([]byte(sampleOutput), nil)
	
	lister := NewPodman(mockExecer)
	
	result, err := lister.List(context.Background())
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Containers, 1)
	assert.Equal(t, "container1", result.Containers[0].ID)
	assert.Equal(t, "ubuntu:20.04", result.Containers[0].Image)
	assert.True(t, result.Containers[0].Running)
	assert.Equal(t, "running", result.Containers[0].Status)
	
	mockExecer.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
}

func TestPodman_CreateContainer(t *testing.T) {
	mockExecer := new(MockExecer)
	
	mockCmd := NewMockCmd()
	
	mockExecer.On("CommandContext", mock.Anything, "podman", mock.Anything).Return(mockCmd)
	
	mockCmd.On("CombinedOutput").Return([]byte("container1"), nil)
	
	lister := NewPodman(mockExecer)
	
	options := map[string]interface{}{
		"image": "ubuntu:20.04",
		"env": map[string]string{
			"VAR1": "value1",
		},
		"volumes": map[string]string{
			"/host/path": "/container/path",
		},
		"ports": map[string]string{
			"8080": "80",
		},
		"cmd": []string{"/bin/bash"},
	}
	
	id, err := lister.CreateContainer(context.Background(), "test-container", "ubuntu:20.04", options)
	
	assert.NoError(t, err)
	assert.Equal(t, "container1", id)
	
	mockExecer.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
}

func TestPodman_StopContainer(t *testing.T) {
	mockExecer := new(MockExecer)
	
	mockCmd := NewMockCmd()
	
	mockExecer.On("CommandContext", mock.Anything, "podman", mock.Anything).Return(mockCmd)
	
	mockCmd.On("CombinedOutput").Return([]byte(""), nil)
	
	lister := NewPodman(mockExecer)
	
	err := lister.StopContainer(context.Background(), "container1")
	
	assert.NoError(t, err)
	
	mockExecer.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
}

func TestPodman_RemoveContainer(t *testing.T) {
	mockExecer := new(MockExecer)
	
	mockCmd := NewMockCmd()
	
	mockExecer.On("CommandContext", mock.Anything, "podman", mock.Anything).Return(mockCmd)
	
	mockCmd.On("CombinedOutput").Return([]byte(""), nil)
	
	lister := NewPodman(mockExecer)
	
	err := lister.RemoveContainer(context.Background(), "container1")
	
	assert.NoError(t, err)
	
	mockExecer.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
}

func TestPodman_ExecInContainer(t *testing.T) {
	mockExecer := new(MockExecer)
	
	mockCmd := NewMockCmd()
	
	mockExecer.On("CommandContext", mock.Anything, "podman", mock.Anything).Return(mockCmd)
	
	mockCmd.On("CombinedOutput").Return([]byte("command output"), nil)
	
	lister := NewPodman(mockExecer)
	
	output, err := lister.ExecInContainer(context.Background(), "container1", []string{"ls", "-la"})
	
	assert.NoError(t, err)
	assert.Equal(t, "command output", output)
	
	mockExecer.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
}
