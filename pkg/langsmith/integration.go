package langsmith

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
)

type LangSmithTracer struct {
	client      *Client
	projectName string
	runID       string
	parentRunID string
	disabled    bool
}

func NewLangSmithTracer(client *Client, projectName string) *LangSmithTracer {
	if client == nil {
		client = NewClient(nil)
	}

	if projectName == "" {
		projectName = client.config.ProjectName
	}

	return &LangSmithTracer{
		client:      client,
		projectName: projectName,
		disabled:    client.config.Disabled,
	}
}

func (t *LangSmithTracer) StartRun(ctx context.Context, name string, runType string, inputs map[string]interface{}, tags []string) (string, error) {
	if t.disabled {
		runID := uuid.New().String()
		t.runID = runID
		return runID, nil
	}

	startTime := time.Now()
	req := &CreateRunRequest{
		Name:        name,
		RunType:     runType,
		StartTime:   startTime,
		Inputs:      inputs,
		ProjectName: t.projectName,
		Status:      "running",
		Tags:        tags,
		ParentRunID: t.parentRunID,
	}

	run, err := t.client.CreateRun(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create run: %w", err)
	}

	t.runID = run.ID
	return run.ID, nil
}

func (t *LangSmithTracer) EndRun(ctx context.Context, outputs map[string]interface{}, error string) error {
	if t.disabled {
		return nil
	}

	endTime := time.Now()
	status := "completed"
	if error != "" {
		status = "failed"
	}

	update := map[string]interface{}{
		"end_time": endTime,
		"outputs":  outputs,
		"status":   status,
	}

	if error != "" {
		update["error"] = error
	}

	_, err := t.client.UpdateRun(ctx, t.runID, update)
	if err != nil {
		return fmt.Errorf("failed to update run: %w", err)
	}

	return nil
}

func (t *LangSmithTracer) TraceEvent(ctx context.Context, name string, eventType string, data map[string]interface{}) error {
	if t.disabled {
		return nil
	}

	childRunID, err := t.StartRun(ctx, name, eventType, data, nil)
	if err != nil {
		return fmt.Errorf("failed to start child run: %w", err)
	}

	_, err = t.client.UpdateRun(ctx, childRunID, map[string]interface{}{
		"end_time": time.Now(),
		"outputs":  data,
		"status":   "completed",
	})
	if err != nil {
		return fmt.Errorf("failed to update child run: %w", err)
	}

	return nil
}

func (t *LangSmithTracer) WithParentRun(parentRunID string) *LangSmithTracer {
	return &LangSmithTracer{
		client:      t.client,
		projectName: t.projectName,
		parentRunID: parentRunID,
		disabled:    t.disabled,
	}
}

type LangGraphTracer struct {
	tracer *LangSmithTracer
}

func NewLangGraphTracer(client *Client, projectName string) *LangGraphTracer {
	return &LangGraphTracer{
		tracer: NewLangSmithTracer(client, projectName),
	}
}

func (t *LangGraphTracer) OnGraphStart(ctx context.Context, graph *langraph.Graph, inputs map[string]interface{}) (context.Context, error) {
	runID, err := t.tracer.StartRun(ctx, graph.Name, "graph", inputs, []string{"graph"})
	if err != nil {
		return ctx, fmt.Errorf("failed to start graph run: %w", err)
	}

	ctx = context.WithValue(ctx, contextKeyRunID, runID)
	return ctx, nil
}

func (t *LangGraphTracer) OnGraphEnd(ctx context.Context, graph *langraph.Graph, outputs map[string]interface{}, err error) error {
	runID, ok := ctx.Value(contextKeyRunID).(string)
	if !ok {
		return fmt.Errorf("run ID not found in context")
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	childTracer := t.tracer.WithParentRun(runID)
	return childTracer.EndRun(ctx, outputs, errStr)
}

func (t *LangGraphTracer) OnNodeStart(ctx context.Context, node *langraph.Node, inputs map[string]interface{}) (context.Context, error) {
	parentRunID, ok := ctx.Value(contextKeyRunID).(string)
	if !ok {
		return ctx, fmt.Errorf("parent run ID not found in context")
	}

	childTracer := t.tracer.WithParentRun(parentRunID)
	runID, err := childTracer.StartRun(ctx, node.Name, "node", inputs, []string{"node", node.Type})
	if err != nil {
		return ctx, fmt.Errorf("failed to start node run: %w", err)
	}

	ctx = context.WithValue(ctx, contextKeyNodeRunID(node.ID), runID)
	return ctx, nil
}

func (t *LangGraphTracer) OnNodeEnd(ctx context.Context, node *langraph.Node, outputs map[string]interface{}, err error) error {
	nodeRunIDKey := contextKeyNodeRunID(node.ID)
	runID, ok := ctx.Value(nodeRunIDKey).(string)
	if !ok {
		return fmt.Errorf("node run ID not found in context for node %s", node.ID)
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	childTracer := t.tracer.WithParentRun(runID)
	return childTracer.EndRun(ctx, outputs, errStr)
}

func (t *LangGraphTracer) OnAgentAction(ctx context.Context, agentID string, action string, inputs map[string]interface{}) (context.Context, error) {
	parentRunID, ok := ctx.Value(contextKeyRunID).(string)
	if !ok {
		return ctx, fmt.Errorf("parent run ID not found in context")
	}

	childTracer := t.tracer.WithParentRun(parentRunID)
	runID, err := childTracer.StartRun(ctx, fmt.Sprintf("Agent %s: %s", agentID, action), "agent_action", inputs, []string{"agent_action", agentID})
	if err != nil {
		return ctx, fmt.Errorf("failed to start agent action run: %w", err)
	}

	actionKey := fmt.Sprintf("%s:%s", agentID, action)
	ctx = context.WithValue(ctx, contextKeyActionRunID(actionKey), runID)
	return ctx, nil
}

func (t *LangGraphTracer) OnAgentActionEnd(ctx context.Context, agentID string, action string, outputs map[string]interface{}, err error) error {
	actionKey := fmt.Sprintf("%s:%s", agentID, action)
	runIDKey := contextKeyActionRunID(actionKey)
	runID, ok := ctx.Value(runIDKey).(string)
	if !ok {
		return fmt.Errorf("action run ID not found in context for agent %s action %s", agentID, action)
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	childTracer := t.tracer.WithParentRun(runID)
	return childTracer.EndRun(ctx, outputs, errStr)
}

type contextKey string

const (
	contextKeyRunID contextKey = "langsmith_run_id"
)

func contextKeyNodeRunID(nodeID string) contextKey {
	return contextKey(fmt.Sprintf("langsmith_node_run_id:%s", nodeID))
}

func contextKeyActionRunID(actionKey string) contextKey {
	return contextKey(fmt.Sprintf("langsmith_action_run_id:%s", actionKey))
}
