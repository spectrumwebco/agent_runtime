package langsmith

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
)

type LangGraphIntegration struct {
	client      *Client
	config      *LangSmithConfig
	projectName string
}

func NewLangGraphIntegration(config *LangSmithConfig) (*LangGraphIntegration, error) {
	if config == nil {
		var err error
		config, err = LoadConfig("")
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	client := NewClient(&ClientConfig{
		APIKey:      config.APIKey,
		APIUrl:      config.APIUrl,
		ProjectName: config.ProjectName,
		Disabled:    !config.Enabled,
	})

	return &LangGraphIntegration{
		client:      client,
		config:      config,
		projectName: config.ProjectName,
	}, nil
}

func (i *LangGraphIntegration) CreateTracer() langraph.Tracer {
	return NewLangGraphTracer(i.client, i.projectName)
}

func (i *LangGraphIntegration) RegisterGraph(ctx context.Context, graph *langraph.Graph) (string, error) {
	if i.client.config.Disabled {
		return uuid.New().String(), nil
	}

	run, err := i.client.CreateRun(ctx, &CreateRunRequest{
		Name:        graph.Name,
		RunType:     "graph",
		StartTime:   time.Now(),
		Inputs:      map[string]interface{}{},
		ProjectName: i.projectName,
		Status:      "running",
		Tags:        []string{"graph", "registration"},
		ExtraData: map[string]interface{}{
			"graph_id":   graph.ID,
			"node_count": len(graph.Nodes),
			"edge_count": len(graph.Edges),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to register graph: %w", err)
	}

	nodeInfo := make([]map[string]interface{}, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		nodeInfo = append(nodeInfo, map[string]interface{}{
			"id":   node.ID,
			"name": node.Name,
			"type": node.Type,
		})
	}

	edgeInfo := make([]map[string]interface{}, 0, len(graph.Edges))
	for _, edge := range graph.Edges {
		edgeInfo = append(edgeInfo, map[string]interface{}{
			"source": edge.Source.ID,
			"target": edge.Target.ID,
		})
	}

	_, err = i.client.UpdateRun(ctx, run.ID, map[string]interface{}{
		"end_time": time.Now(),
		"outputs": map[string]interface{}{
			"nodes": nodeInfo,
			"edges": edgeInfo,
		},
		"status": "completed",
	})
	if err != nil {
		return "", fmt.Errorf("failed to update graph registration: %w", err)
	}

	return run.ID, nil
}

func (i *LangGraphIntegration) RegisterMultiAgentSystem(ctx context.Context, system *langraph.MultiAgentSystem) (string, error) {
	if i.client.config.Disabled {
		return uuid.New().String(), nil
	}

	run, err := i.client.CreateRun(ctx, &CreateRunRequest{
		Name:        system.Name(),
		RunType:     "multi_agent_system",
		StartTime:   time.Now(),
		Inputs:      map[string]interface{}{},
		ProjectName: i.projectName,
		Status:      "running",
		Tags:        []string{"multi_agent_system", "registration"},
		ExtraData: map[string]interface{}{
			"system_id":   system.ID(),
			"agent_count": len(system.Agents()),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to register multi-agent system: %w", err)
	}

	agentInfo := make([]map[string]interface{}, 0, len(system.Agents()))
	for _, agent := range system.Agents() {
		agentInfo = append(agentInfo, map[string]interface{}{
			"id":   agent.ID(),
			"name": agent.Name(),
			"type": agent.Type(),
		})
	}

	_, err = i.client.UpdateRun(ctx, run.ID, map[string]interface{}{
		"end_time": time.Now(),
		"outputs": map[string]interface{}{
			"agents": agentInfo,
		},
		"status": "completed",
	})
	if err != nil {
		return "", fmt.Errorf("failed to update multi-agent system registration: %w", err)
	}

	return run.ID, nil
}

func (i *LangGraphIntegration) TraceAgentTask(ctx context.Context, agentID, taskID, taskType, description string, inputs map[string]interface{}) (context.Context, string, error) {
	if i.client.config.Disabled {
		return ctx, uuid.New().String(), nil
	}

	run, err := i.client.CreateRun(ctx, &CreateRunRequest{
		Name:        fmt.Sprintf("Task: %s", description),
		RunType:     "agent_task",
		StartTime:   time.Now(),
		Inputs:      inputs,
		ProjectName: i.projectName,
		Status:      "running",
		Tags:        []string{"agent_task", taskType, agentID},
		ExtraData: map[string]interface{}{
			"agent_id":    agentID,
			"task_id":     taskID,
			"task_type":   taskType,
			"description": description,
		},
	})
	if err != nil {
		return ctx, "", fmt.Errorf("failed to trace agent task: %w", err)
	}

	ctx = context.WithValue(ctx, contextKeyTaskRunID(taskID), run.ID)
	return ctx, run.ID, nil
}

func (i *LangGraphIntegration) EndAgentTask(ctx context.Context, taskID string, outputs map[string]interface{}, err error) error {
	if i.client.config.Disabled {
		return nil
	}

	runIDKey := contextKeyTaskRunID(taskID)
	runID, ok := ctx.Value(runIDKey).(string)
	if !ok {
		return fmt.Errorf("task run ID not found in context for task %s", taskID)
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	status := "completed"
	if err != nil {
		status = "failed"
	}

	_, err = i.client.UpdateRun(ctx, runID, map[string]interface{}{
		"end_time": time.Now(),
		"outputs":  outputs,
		"status":   status,
		"error":    errStr,
	})
	if err != nil {
		return fmt.Errorf("failed to end agent task trace: %w", err)
	}

	return nil
}

func (i *LangGraphIntegration) TraceAgentAction(ctx context.Context, agentID, actionType string, inputs map[string]interface{}) (context.Context, string, error) {
	if i.client.config.Disabled {
		return ctx, uuid.New().String(), nil
	}

	var parentRunID string
	if taskRunID, ok := ctx.Value(contextKeyTaskRunID("current")).(string); ok {
		parentRunID = taskRunID
	}

	run, err := i.client.CreateRun(ctx, &CreateRunRequest{
		Name:        fmt.Sprintf("Agent %s: %s", agentID, actionType),
		RunType:     "agent_action",
		StartTime:   time.Now(),
		Inputs:      inputs,
		ProjectName: i.projectName,
		Status:      "running",
		Tags:        []string{"agent_action", actionType, agentID},
		ParentRunID: parentRunID,
		ExtraData: map[string]interface{}{
			"agent_id":    agentID,
			"action_type": actionType,
		},
	})
	if err != nil {
		return ctx, "", fmt.Errorf("failed to trace agent action: %w", err)
	}

	actionKey := fmt.Sprintf("%s:%s:%s", agentID, actionType, uuid.New().String())
	ctx = context.WithValue(ctx, contextKeyActionRunID(actionKey), run.ID)
	return ctx, run.ID, nil
}

func (i *LangGraphIntegration) EndAgentAction(ctx context.Context, actionKey string, outputs map[string]interface{}, err error) error {
	if i.client.config.Disabled {
		return nil
	}

	runIDKey := contextKeyActionRunID(actionKey)
	runID, ok := ctx.Value(runIDKey).(string)
	if !ok {
		return fmt.Errorf("action run ID not found in context for action %s", actionKey)
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	status := "completed"
	if err != nil {
		status = "failed"
	}

	_, err = i.client.UpdateRun(ctx, runID, map[string]interface{}{
		"end_time": time.Now(),
		"outputs":  outputs,
		"status":   status,
		"error":    errStr,
	})
	if err != nil {
		return fmt.Errorf("failed to end agent action trace: %w", err)
	}

	return nil
}

func (i *LangGraphIntegration) TraceAgentMessage(ctx context.Context, senderID, receiverID, messageType, content string, metadata map[string]interface{}) (string, error) {
	if i.client.config.Disabled {
		return uuid.New().String(), nil
	}

	var parentRunID string
	if taskRunID, ok := ctx.Value(contextKeyTaskRunID("current")).(string); ok {
		parentRunID = taskRunID
	}

	run, err := i.client.CreateRun(ctx, &CreateRunRequest{
		Name:        fmt.Sprintf("Message: %s -> %s", senderID, receiverID),
		RunType:     "agent_message",
		StartTime:   time.Now(),
		Inputs: map[string]interface{}{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"type":        messageType,
			"content":     content,
			"metadata":    metadata,
		},
		ProjectName: i.projectName,
		Status:      "completed",
		Tags:        []string{"agent_message", messageType, senderID, receiverID},
		ParentRunID: parentRunID,
		ExtraData: map[string]interface{}{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"type":        messageType,
		},
		EndTime: timePtr(time.Now()),
	})
	if err != nil {
		return "", fmt.Errorf("failed to trace agent message: %w", err)
	}

	return run.ID, nil
}

func (i *LangGraphIntegration) CreateFeedback(ctx context.Context, runID, key string, value interface{}, comment string) error {
	if i.client.config.Disabled {
		return nil
	}

	return i.client.CreateFeedback(ctx, runID, key, value, comment)
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func contextKeyTaskRunID(taskID string) contextKey {
	return contextKey(fmt.Sprintf("langsmith_task_run_id:%s", taskID))
}

func contextKeyActionRunID(actionKey string) contextKey {
	return contextKey(fmt.Sprintf("langsmith_action_run_id:%s", actionKey))
}
