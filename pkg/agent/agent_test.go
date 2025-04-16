package agent

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventStream struct {
	mock.Mock
	Events []*models.Event
}

func (m *MockEventStream) AddEvent(event *models.Event) error {
	m.Events = append(m.Events, event)
	args := m.Called(event)
	return args.Error(0)
}

type MockStateManager struct {
	mock.Mock
	States map[string]*statemodels.AgentState
}

func (m *MockStateManager) GetState(agentID string) (*statemodels.AgentState, error) {
	args := m.Called(agentID)
	return args.Get(0).(*statemodels.AgentState), args.Error(1)
}

func (m *MockStateManager) UpdateState(agentID string, state *statemodels.AgentState) error {
	m.States[agentID] = state
	args := m.Called(agentID, state)
	return args.Error(0)
}

type MockTool struct {
	mock.Mock
	ToolName        string
	ToolDescription string
}

func (m *MockTool) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockTool) Description() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockTool) Execute(ctx context.Context, args map[string]interface{}) (map[string]interface{}, error) {
	callArgs := m.Called(ctx, args)
	return callArgs.Get(0).(map[string]interface{}), callArgs.Error(1)
}

func TestNewAgent(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	agent := NewAgent("test-agent", "Test Agent", "test", []string{"capability1", "capability2"}, eventStream, stateManager)

	assert.Equal(t, "test-agent", agent.ID)
	assert.Equal(t, "Test Agent", agent.Name)
	assert.Equal(t, "test", agent.Role)
	assert.Equal(t, []string{"capability1", "capability2"}, agent.Capabilities)
	assert.Equal(t, "initialized", agent.State.Status)
	assert.NotNil(t, agent.Configuration)
	assert.NotNil(t, agent.State)
	assert.Equal(t, eventStream, agent.eventStream)
	assert.Equal(t, stateManager, agent.stateManager)
	assert.Empty(t, agent.tools)
	assert.False(t, agent.isRunning)
}

func TestAgentStartStop(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	eventStream.On("AddEvent", mock.Anything).Return(nil)
	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	agent := NewAgent("test-agent", "Test Agent", "test", []string{"capability1", "capability2"}, eventStream, stateManager)

	ctx := context.Background()
	err := agent.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, agent.isRunning)
	assert.Equal(t, "running", agent.State.Status)

	err = agent.Start(ctx)
	assert.Error(t, err)

	err = agent.Stop()
	assert.NoError(t, err)
	assert.False(t, agent.isRunning)
	assert.Equal(t, "stopped", agent.State.Status)

	err = agent.Stop()
	assert.Error(t, err)

	assert.GreaterOrEqual(t, len(eventStream.Events), 2)
	assert.Equal(t, "agent_started", eventStream.Events[0].Data.(map[string]interface{})["action"])
	assert.Equal(t, "agent_stopped", eventStream.Events[1].Data.(map[string]interface{})["action"])
}

func TestAgentTools(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	tool := &MockTool{
		ToolName:        "test-tool",
		ToolDescription: "Test tool",
	}
	tool.On("Name").Return("test-tool")
	tool.On("Description").Return("Test tool")
	tool.On("Execute", mock.Anything, mock.Anything).Return(map[string]interface{}{
		"result": "success",
	}, nil)

	eventStream.On("AddEvent", mock.Anything).Return(nil)
	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	agent := NewAgent("test-agent", "Test Agent", "test", []string{"capability1", "capability2"}, eventStream, stateManager)

	agent.AddTool(tool)

	tools := agent.GetTools()
	assert.Len(t, tools, 1)
	assert.Equal(t, tool, tools[0])

	ctx := context.Background()
	result, err := agent.ExecuteTool(ctx, "test-tool", map[string]interface{}{
		"param1": "value1",
	})
	assert.NoError(t, err)
	assert.Equal(t, "success", result["result"])

	result, err = agent.ExecuteTool(ctx, "non-existent-tool", map[string]interface{}{})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAgentConfiguration(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	eventStream.On("AddEvent", mock.Anything).Return(nil)
	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	agent := NewAgent("test-agent", "Test Agent", "test", []string{"capability1", "capability2"}, eventStream, stateManager)

	agent.UpdateConfiguration(map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	})

	config := agent.GetConfiguration()
	assert.Equal(t, "value1", config["key1"])
	assert.Equal(t, 42, config["key2"])

	agent.UpdateConfiguration(map[string]interface{}{
		"key2": 43,
		"key3": true,
	})

	config = agent.GetConfiguration()
	assert.Equal(t, "value1", config["key1"])
	assert.Equal(t, 43, config["key2"])
	assert.Equal(t, true, config["key3"])
}

func TestAgentState(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	eventStream.On("AddEvent", mock.Anything).Return(nil)
	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	agent := NewAgent("test-agent", "Test Agent", "test", []string{"capability1", "capability2"}, eventStream, stateManager)

	state := agent.GetState()
	assert.Equal(t, "test-agent", state.ID)
	assert.Equal(t, "initialized", state.Status)

	newState := &statemodels.AgentState{
		ID:        "test-agent",
		Status:    "processing",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"progress": 50,
		},
	}
	err := agent.UpdateState(newState)
	assert.NoError(t, err)

	state = agent.GetState()
	assert.Equal(t, "test-agent", state.ID)
	assert.Equal(t, "processing", state.Status)
	assert.Equal(t, 50, state.Data["progress"])
}

func TestAgentCapabilities(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	agent := NewAgent("test-agent", "Test Agent", "test", []string{"capability1", "capability2"}, eventStream, stateManager)

	assert.True(t, agent.HasCapability("capability1"))
	assert.True(t, agent.HasCapability("capability2"))
	assert.False(t, agent.HasCapability("capability3"))
}

func TestAgentFactory(t *testing.T) {
	eventStream := &MockEventStream{}
	stateManager := &MockStateManager{
		States: make(map[string]*statemodels.AgentState),
	}

	eventStream.On("AddEvent", mock.Anything).Return(nil)
	stateManager.On("UpdateState", mock.Anything, mock.Anything).Return(nil)

	factory := NewAgentFactory(eventStream, stateManager)

	agent1, err := factory.CreateAgent("agent1", "Agent 1", "role1", []string{"capability1"})
	assert.NoError(t, err)
	assert.NotNil(t, agent1)

	agent2, err := factory.CreateAgent("agent2", "Agent 2", "role2", []string{"capability2"})
	assert.NoError(t, err)
	assert.NotNil(t, agent2)

	agent3, err := factory.CreateAgent("agent1", "Agent 3", "role3", []string{"capability3"})
	assert.Error(t, err)
	assert.Nil(t, agent3)

	agent, err := factory.GetAgent("agent1")
	assert.NoError(t, err)
	assert.Equal(t, agent1, agent)

	agent, err = factory.GetAgent("non-existent-agent")
	assert.Error(t, err)
	assert.Nil(t, agent)

	agents := factory.ListAgents()
	assert.Len(t, agents, 2)

	agents = factory.GetAgentsByRole("role1")
	assert.Len(t, agents, 1)
	assert.Equal(t, agent1, agents[0])

	agents = factory.GetAgentsByCapability("capability2")
	assert.Len(t, agents, 1)
	assert.Equal(t, agent2, agents[0])

	err = factory.DeleteAgent("agent1")
	assert.NoError(t, err)

	agents = factory.ListAgents()
	assert.Len(t, agents, 1)
	assert.Equal(t, agent2, agents[0])

	err = factory.DeleteAgent("non-existent-agent")
	assert.Error(t, err)
}
