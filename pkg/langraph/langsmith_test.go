package langraph

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/langsmith"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiAgentSystemWithLangSmith(t *testing.T) {
	if os.Getenv("CI") != "true" && os.Getenv("TEST_LANGSMITH") != "true" {
		t.Skip("Skipping LangSmith integration test. Set TEST_LANGSMITH=true to run this test.")
	}

	config := &langsmith.LangSmithConfig{
		Enabled:     true,
		APIKey:      os.Getenv("LANGCHAIN_API_KEY"),
		APIUrl:      getEnvOrDefault("LANGCHAIN_ENDPOINT", "https://api.smith.langchain.com"),
		ProjectName: "test-multi-agent-system",
		SelfHosted:  os.Getenv("LANGSMITH_SELF_HOSTED") == "true",
	}

	integration, err := langsmith.NewLangGraphIntegration(config)
	require.NoError(t, err)
	require.NotNil(t, integration)

	tracer := integration.CreateTracer()
	require.NotNil(t, tracer)

	executor := NewExecutor(WithTracer(tracer))
	require.NotNil(t, executor)

	systemConfig := &MultiAgentSystemConfig{
		Name: "Test Multi-Agent System",
		Agents: []*AgentConfig{
			{
				ID:   "frontend",
				Name: "Frontend Agent",
				Type: "frontend",
				Handler: func(ctx context.Context, task *Task) (*TaskResult, error) {
					return &TaskResult{
						AgentID: "frontend",
						TaskID:  task.ID,
						Status:  "completed",
						Output:  map[string]interface{}{"message": "Frontend task completed"},
					}, nil
				},
			},
			{
				ID:   "codegen",
				Name: "Codegen Agent",
				Type: "codegen",
				Handler: func(ctx context.Context, task *Task) (*TaskResult, error) {
					return &TaskResult{
						AgentID: "codegen",
						TaskID:  task.ID,
						Status:  "completed",
						Output: map[string]interface{}{
							"code": "function hello() { return 'Hello, world!'; }",
						},
					}, nil
				},
			},
			{
				ID:   "engineering",
				Name: "Engineering Agent",
				Type: "engineering",
				Handler: func(ctx context.Context, task *Task) (*TaskResult, error) {
					return &TaskResult{
						AgentID: "engineering",
						TaskID:  task.ID,
						Status:  "completed",
						Output: map[string]interface{}{
							"design": "Architecture design document",
						},
					}, nil
				},
			},
		},
	}

	system, err := NewMultiAgentSystem(systemConfig, executor)
	require.NoError(t, err)
	require.NotNil(t, system)

	systemID, err := integration.RegisterMultiAgentSystem(context.Background(), system)
	require.NoError(t, err)
	require.NotEmpty(t, systemID)

	task := &Task{
		ID:          "task-1",
		Description: "Build a web application",
		AssignedTo:  "frontend",
		Input: map[string]interface{}{
			"requirements": "Create a simple web application",
		},
	}

	ctx := context.Background()
	
	taskCtx, taskRunID, err := integration.TraceAgentTask(
		ctx,
		task.AssignedTo,
		task.ID,
		"task",
		task.Description,
		task.Input,
	)
	require.NoError(t, err)
	require.NotEmpty(t, taskRunID)

	result, err := system.ExecuteTask(taskCtx, task)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "frontend", result.AgentID)
	assert.Equal(t, "task-1", result.TaskID)
	assert.Equal(t, "completed", result.Status)
	assert.Contains(t, result.Output, "message")

	err = integration.EndAgentTask(
		taskCtx,
		task.ID,
		result.Output,
		nil,
	)
	require.NoError(t, err)

	messageCtx, actionRunID, err := integration.TraceAgentAction(
		taskCtx,
		"frontend",
		"send_message",
		map[string]interface{}{
			"receiver": "codegen",
			"content":  "Please generate code for the frontend",
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, actionRunID)

	messageID, err := integration.TraceAgentMessage(
		messageCtx,
		"frontend",
		"codegen",
		"request",
		"Please generate code for the frontend",
		map[string]interface{}{"priority": "high"},
	)
	require.NoError(t, err)
	require.NotEmpty(t, messageID)

	err = integration.EndAgentAction(
		messageCtx,
		"frontend:send_message:"+actionRunID,
		map[string]interface{}{"status": "sent"},
		nil,
	)
	require.NoError(t, err)

	codegenTask := &Task{
		ID:          "task-2",
		Description: "Generate code for the frontend",
		AssignedTo:  "codegen",
		Input: map[string]interface{}{
			"requirements": "Generate JavaScript code for the frontend",
		},
	}

	codegenCtx, codegenRunID, err := integration.TraceAgentTask(
		ctx,
		codegenTask.AssignedTo,
		codegenTask.ID,
		"task",
		codegenTask.Description,
		codegenTask.Input,
	)
	require.NoError(t, err)
	require.NotEmpty(t, codegenRunID)

	codegenResult, err := system.ExecuteTask(codegenCtx, codegenTask)
	require.NoError(t, err)
	require.NotNil(t, codegenResult)
	assert.Equal(t, "codegen", codegenResult.AgentID)
	assert.Equal(t, "task-2", codegenResult.TaskID)
	assert.Equal(t, "completed", codegenResult.Status)
	assert.Contains(t, codegenResult.Output, "code")

	err = integration.EndAgentTask(
		codegenCtx,
		codegenTask.ID,
		codegenResult.Output,
		nil,
	)
	require.NoError(t, err)

	responseCtx, responseActionID, err := integration.TraceAgentAction(
		codegenCtx,
		"codegen",
		"send_message",
		map[string]interface{}{
			"receiver": "frontend",
			"content":  "Here's the code you requested",
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, responseActionID)

	responseMessageID, err := integration.TraceAgentMessage(
		responseCtx,
		"codegen",
		"frontend",
		"response",
		"Here's the code you requested",
		map[string]interface{}{
			"code": "function hello() { return 'Hello, world!'; }",
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, responseMessageID)

	err = integration.EndAgentAction(
		responseCtx,
		"codegen:send_message:"+responseActionID,
		map[string]interface{}{"status": "sent"},
		nil,
	)
	require.NoError(t, err)

	err = integration.CreateFeedback(
		context.Background(),
		codegenRunID,
		"code_quality",
		0.9,
		"The code is well-structured and follows best practices",
	)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
