package orchestration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
)

type AgentOrchestrator struct {
	graphManager    *langgraph.GraphManager
	contextMonitor  *langgraph.ContextMonitor
	workflows       map[string]*Workflow
	activeWorkflows map[string]bool
	mu              sync.RWMutex
}

type Workflow struct {
	ID          string
	Name        string
	Description string
	Graph       *langgraph.Graph
	Config      *WorkflowConfig
	Status      WorkflowStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkflowConfig struct {
	AutoTrigger       bool
	TriggerConditions []TriggerCondition
	MaxIterations     int
	Timeout           time.Duration
	RetryCount        int
	RetryDelay        time.Duration
}

type TriggerCondition struct {
	Type      TriggerType
	Pattern   string
	Priority  int
	Threshold float64
}

type TriggerType string

const (
	ContextTrigger TriggerType = "context"
	TimeTrigger TriggerType = "time"
	EventTrigger TriggerType = "event"
	ManualTrigger TriggerType = "manual"
)

type WorkflowStatus string

const (
	WorkflowStatusIdle WorkflowStatus = "idle"
	WorkflowStatusRunning WorkflowStatus = "running"
	WorkflowStatusPaused WorkflowStatus = "paused"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed WorkflowStatus = "failed"
)

func NewAgentOrchestrator(graphManager *langgraph.GraphManager, contextMonitor *langgraph.ContextMonitor) *AgentOrchestrator {
	return &AgentOrchestrator{
		graphManager:    graphManager,
		contextMonitor:  contextMonitor,
		workflows:       make(map[string]*Workflow),
		activeWorkflows: make(map[string]bool),
	}
}

func (ao *AgentOrchestrator) RegisterWorkflow(workflow *Workflow) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	if _, exists := ao.workflows[workflow.ID]; exists {
		return fmt.Errorf("workflow with ID %s already exists", workflow.ID)
	}

	ao.workflows[workflow.ID] = workflow
	return nil
}

func (ao *AgentOrchestrator) GetWorkflow(id string) (*Workflow, error) {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	workflow, exists := ao.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow with ID %s not found", id)
	}

	return workflow, nil
}

func (ao *AgentOrchestrator) ListWorkflows() []*Workflow {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	workflows := make([]*Workflow, 0, len(ao.workflows))
	for _, workflow := range ao.workflows {
		workflows = append(workflows, workflow)
	}

	return workflows
}

func (ao *AgentOrchestrator) StartWorkflow(id string, initialState map[string]interface{}) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	workflow, exists := ao.workflows[id]
	if !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}

	if workflow.Status == WorkflowStatusRunning {
		return fmt.Errorf("workflow with ID %s is already running", id)
	}

	workflow.Status = WorkflowStatusRunning
	workflow.UpdatedAt = time.Now()

	ao.activeWorkflows[id] = true

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), workflow.Config.Timeout)
		defer cancel()

		result, err := workflow.Graph.Run(ctx, initialState)
		
		ao.mu.Lock()
		defer ao.mu.Unlock()

		if err != nil {
			workflow.Status = WorkflowStatusFailed
		} else {
			workflow.Status = WorkflowStatusCompleted
		}

		workflow.UpdatedAt = time.Now()

		delete(ao.activeWorkflows, id)
	}()

	return nil
}

func (ao *AgentOrchestrator) StopWorkflow(id string) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	workflow, exists := ao.workflows[id]
	if !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}

	if workflow.Status != WorkflowStatusRunning {
		return fmt.Errorf("workflow with ID %s is not running", id)
	}

	workflow.Status = WorkflowStatusPaused
	workflow.UpdatedAt = time.Now()

	delete(ao.activeWorkflows, id)

	return nil
}

func (ao *AgentOrchestrator) CreateWorkflow(name, description string, config *WorkflowConfig) (*Workflow, error) {
	graph := langgraph.NewGraph(name)

	ao.graphManager.RegisterGraph(graph)

	workflow := &Workflow{
		ID:          fmt.Sprintf("workflow-%d", time.Now().UnixNano()),
		Name:        name,
		Description: description,
		Graph:       graph,
		Config:      config,
		Status:      WorkflowStatusIdle,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := ao.RegisterWorkflow(workflow)
	if err != nil {
		return nil, err
	}

	if config.AutoTrigger {
		for _, condition := range config.TriggerConditions {
			if condition.Type == ContextTrigger {
				ao.contextMonitor.RegisterTrigger(&langgraph.ContextTrigger{
					ID:        workflow.ID,
					Pattern:   condition.Pattern,
					Priority:  condition.Priority,
					Threshold: condition.Threshold,
					Callback: func(ctx context.Context, matchedContext map[string]interface{}) error {
						return ao.StartWorkflow(workflow.ID, matchedContext)
					},
				})
			}
		}
	}

	return workflow, nil
}

func (ao *AgentOrchestrator) DeleteWorkflow(id string) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	workflow, exists := ao.workflows[id]
	if !exists {
		return fmt.Errorf("workflow with ID %s not found", id)
	}

	if workflow.Status == WorkflowStatusRunning {
		return fmt.Errorf("cannot delete workflow with ID %s while it is running", id)
	}

	ao.graphManager.UnregisterGraph(workflow.Graph.GetID())

	delete(ao.workflows, id)

	return nil
}

func (ao *AgentOrchestrator) StartMonitoring(ctx context.Context) error {
	return ao.contextMonitor.StartMonitoring(ctx)
}

func (ao *AgentOrchestrator) StopMonitoring() error {
	return ao.contextMonitor.StopMonitoring()
}

func (ao *AgentOrchestrator) UpdateContext(context map[string]interface{}) error {
	return ao.contextMonitor.UpdateContext(context)
}

func (ao *AgentOrchestrator) GetActiveWorkflows() []*Workflow {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	activeWorkflows := make([]*Workflow, 0, len(ao.activeWorkflows))
	for id := range ao.activeWorkflows {
		if workflow, exists := ao.workflows[id]; exists {
			activeWorkflows = append(activeWorkflows, workflow)
		}
	}

	return activeWorkflows
}

func (ao *AgentOrchestrator) GetWorkflowsByStatus(status WorkflowStatus) []*Workflow {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	workflows := make([]*Workflow, 0)
	for _, workflow := range ao.workflows {
		if workflow.Status == status {
			workflows = append(workflows, workflow)
		}
	}

	return workflows
}

func (ao *AgentOrchestrator) ExtendAgentLoop(agentLoopExtension *langgraph.AgentLoopExtension) error {
	config := &WorkflowConfig{
		AutoTrigger:   true,
		MaxIterations: 10,
		Timeout:       time.Minute * 10,
		RetryCount:    3,
		RetryDelay:    time.Second * 5,
		TriggerConditions: []TriggerCondition{
			{
				Type:      ContextTrigger,
				Pattern:   ".*",
				Priority:  1,
				Threshold: 0.7,
			},
		},
	}

	workflow, err := ao.CreateWorkflow("ExtendedAgentLoop", "Extended agent loop workflow", config)
	if err != nil {
		return err
	}

	graph := workflow.Graph

	graph.AddNode("start", agentLoopExtension.StartNode)
	graph.AddNode("process", agentLoopExtension.ProcessNode)
	graph.AddNode("decide", agentLoopExtension.DecideNode)
	graph.AddNode("execute", agentLoopExtension.ExecuteNode)
	graph.AddNode("feedback", agentLoopExtension.FeedbackNode)
	graph.AddNode("end", agentLoopExtension.EndNode)

	graph.AddEdge("start", "process")
	graph.AddEdge("process", "decide")
	graph.AddEdge("decide", "execute", langgraph.WithCondition(func(state map[string]interface{}) bool {
		decision, ok := state["decision"].(string)
		return ok && decision == "execute"
	}))
	graph.AddEdge("decide", "end", langgraph.WithCondition(func(state map[string]interface{}) bool {
		decision, ok := state["decision"].(string)
		return ok && decision == "end"
	}))
	graph.AddEdge("execute", "feedback")
	graph.AddEdge("feedback", "process", langgraph.WithCondition(func(state map[string]interface{}) bool {
		iterations, ok := state["iterations"].(int)
		maxIterations, ok2 := state["max_iterations"].(int)
		return ok && ok2 && iterations < maxIterations
	}))
	graph.AddEdge("feedback", "end", langgraph.WithCondition(func(state map[string]interface{}) bool {
		iterations, ok := state["iterations"].(int)
		maxIterations, ok2 := state["max_iterations"].(int)
		return ok && ok2 && iterations >= maxIterations
	}))

	graph.SetEntryPoint("start")

	return nil
}
