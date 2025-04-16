package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type OrchestratorConfig struct {
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	MaxConcurrentTasks  int                    `json:"max_concurrent_tasks"`
	DefaultTimeout      time.Duration          `json:"default_timeout"`
	EventStream         EventStream            `json:"-"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type Task struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	AgentID     string                 `json:"agent_id,omitempty"`
	AgentRole   AgentRole              `json:"agent_role,omitempty"`
	Status      TaskStatus             `json:"status"`
	Priority    int                    `json:"priority"`
	Inputs      map[string]interface{} `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs,omitempty"`
	Dependencies []string              `json:"dependencies,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   time.Time              `json:"started_at,omitempty"`
	CompletedAt time.Time              `json:"completed_at,omitempty"`
	Timeout     time.Duration          `json:"timeout,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type Orchestrator struct {
	Config          OrchestratorConfig       `json:"config"`
	Tasks           map[string]*Task         `json:"tasks"`
	Agents          map[string]*Agent        `json:"agents"`
	AgentsByRole    map[AgentRole][]*Agent   `json:"agents_by_role"`
	Graph           *Graph                   `json:"graph"`
	Executor        *Executor                `json:"executor"`
	CommunicationMgr *CommunicationManager   `json:"communication_mgr"`
	taskLock        sync.RWMutex             `json:"-"`
	agentLock       sync.RWMutex             `json:"-"`
	runningTasks    int                      `json:"-"`
	runningLock     sync.Mutex               `json:"-"`
}

func NewOrchestrator(config OrchestratorConfig) *Orchestrator {
	if config.MaxConcurrentTasks <= 0 {
		config.MaxConcurrentTasks = 10
	}
	
	if config.DefaultTimeout <= 0 {
		config.DefaultTimeout = 5 * time.Minute
	}
	
	graph := NewGraph(config.Name, config.Description)
	executor := NewExecutor(graph)
	communicationMgr := NewCommunicationManager(config.EventStream)
	
	return &Orchestrator{
		Config:          config,
		Tasks:           make(map[string]*Task),
		Agents:          make(map[string]*Agent),
		AgentsByRole:    make(map[AgentRole][]*Agent),
		Graph:           graph,
		Executor:        executor,
		CommunicationMgr: communicationMgr,
	}
}

func (o *Orchestrator) RegisterAgent(agent *Agent) {
	o.agentLock.Lock()
	defer o.agentLock.Unlock()
	
	o.Agents[agent.Config.ID] = agent
	
	if _, exists := o.AgentsByRole[agent.Config.Role]; !exists {
		o.AgentsByRole[agent.Config.Role] = []*Agent{}
	}
	
	o.AgentsByRole[agent.Config.Role] = append(o.AgentsByRole[agent.Config.Role], agent)
	
	o.CommunicationMgr.RegisterAgent(agent)
	
	o.Graph.AddAgentNode(AgentType(agent.Config.Role), agent.Config.Name, agent.Config.Description, func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
		return agent.Process(ctx, inputs)
	})
	
	if o.Config.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeSystem,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "agent_registered",
				"agent_id":   agent.Config.ID,
				"agent_name": agent.Config.Name,
				"agent_role": agent.Config.Role,
			},
			map[string]string{
				"agent_id": agent.Config.ID,
			},
		)
		
		if err := o.Config.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
}

func (o *Orchestrator) CreateTask(name, description string, agentRole AgentRole, inputs map[string]interface{}, priority int, dependencies []string, timeout time.Duration, metadata map[string]interface{}) (*Task, error) {
	o.taskLock.Lock()
	defer o.taskLock.Unlock()
	
	if timeout <= 0 {
		timeout = o.Config.DefaultTimeout
	}
	
	task := &Task{
		ID:           uuid.New().String(),
		Name:         name,
		Description:  description,
		AgentRole:    agentRole,
		Status:       TaskStatusPending,
		Priority:     priority,
		Inputs:       inputs,
		Dependencies: dependencies,
		CreatedAt:    time.Now().UTC(),
		Timeout:      timeout,
		Metadata:     metadata,
	}
	
	o.Tasks[task.ID] = task
	
	if o.Config.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeTask,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":      "task_created",
				"task_id":     task.ID,
				"task_name":   task.Name,
				"agent_role":  task.AgentRole,
				"priority":    task.Priority,
				"inputs":      task.Inputs,
				"dependencies": task.Dependencies,
			},
			map[string]string{
				"task_id": task.ID,
			},
		)
		
		if err := o.Config.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	if len(task.Dependencies) == 0 {
		go o.scheduleTask(task.ID)
	}
	
	return task, nil
}

func (o *Orchestrator) scheduleTask(taskID string) {
	o.taskLock.RLock()
	task, exists := o.Tasks[taskID]
	o.taskLock.RUnlock()
	
	if !exists {
		fmt.Printf("Task %s does not exist\n", taskID)
		return
	}
	
	if task.Status != TaskStatusPending {
		return
	}
	
	for _, depID := range task.Dependencies {
		o.taskLock.RLock()
		dep, exists := o.Tasks[depID]
		o.taskLock.RUnlock()
		
		if !exists || dep.Status != TaskStatusCompleted {
			return
		}
	}
	
	o.runningLock.Lock()
	if o.runningTasks >= o.Config.MaxConcurrentTasks {
		o.runningLock.Unlock()
		return
	}
	o.runningTasks++
	o.runningLock.Unlock()
	
	o.agentLock.RLock()
	agents, exists := o.AgentsByRole[task.AgentRole]
	o.agentLock.RUnlock()
	
	if !exists || len(agents) == 0 {
		fmt.Printf("No agents available for role %s\n", task.AgentRole)
		
		o.runningLock.Lock()
		o.runningTasks--
		o.runningLock.Unlock()
		
		return
	}
	
	agent := agents[0]
	
	o.taskLock.Lock()
	task.AgentID = agent.Config.ID
	task.Status = TaskStatusAssigned
	o.taskLock.Unlock()
	
	if o.Config.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeTask,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "task_assigned",
				"task_id":    task.ID,
				"task_name":  task.Name,
				"agent_id":   agent.Config.ID,
				"agent_name": agent.Config.Name,
				"agent_role": agent.Config.Role,
			},
			map[string]string{
				"task_id":  task.ID,
				"agent_id": agent.Config.ID,
			},
		)
		
		if err := o.Config.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	go o.executeTask(task.ID)
}

func (o *Orchestrator) executeTask(taskID string) {
	o.taskLock.Lock()
	task, exists := o.Tasks[taskID]
	if !exists {
		o.taskLock.Unlock()
		
		o.runningLock.Lock()
		o.runningTasks--
		o.runningLock.Unlock()
		
		return
	}
	
	task.Status = TaskStatusInProgress
	task.StartedAt = time.Now().UTC()
	o.taskLock.Unlock()
	
	if o.Config.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeTask,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "task_started",
				"task_id":    task.ID,
				"task_name":  task.Name,
				"agent_id":   task.AgentID,
			},
			map[string]string{
				"task_id":  task.ID,
				"agent_id": task.AgentID,
			},
		)
		
		if err := o.Config.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	o.agentLock.RLock()
	agent, exists := o.Agents[task.AgentID]
	o.agentLock.RUnlock()
	
	if !exists {
		o.taskLock.Lock()
		task.Status = TaskStatusFailed
		task.Error = fmt.Sprintf("Agent %s does not exist", task.AgentID)
		task.CompletedAt = time.Now().UTC()
		o.taskLock.Unlock()
		
		o.runningLock.Lock()
		o.runningTasks--
		o.runningLock.Unlock()
		
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()
	
	outputs, err := agent.Process(ctx, task.Inputs)
	
	o.taskLock.Lock()
	task.CompletedAt = time.Now().UTC()
	
	if err != nil {
		task.Status = TaskStatusFailed
		task.Error = err.Error()
	} else {
		task.Status = TaskStatusCompleted
		task.Outputs = outputs
	}
	o.taskLock.Unlock()
	
	if o.Config.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeTask,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "task_completed",
				"task_id":    task.ID,
				"task_name":  task.Name,
				"agent_id":   task.AgentID,
				"status":     task.Status,
				"outputs":    task.Outputs,
				"error":      task.Error,
			},
			map[string]string{
				"task_id":  task.ID,
				"agent_id": task.AgentID,
			},
		)
		
		if err := o.Config.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	o.runningLock.Lock()
	o.runningTasks--
	o.runningLock.Unlock()
	
	o.taskLock.RLock()
	for _, t := range o.Tasks {
		for _, depID := range t.Dependencies {
			if depID == task.ID && t.Status == TaskStatusPending {
				go o.scheduleTask(t.ID)
			}
		}
	}
	o.taskLock.RUnlock()
}

func (o *Orchestrator) GetTask(taskID string) (*Task, error) {
	o.taskLock.RLock()
	defer o.taskLock.RUnlock()
	
	task, exists := o.Tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s does not exist", taskID)
	}
	
	return task, nil
}

func (o *Orchestrator) GetAgent(agentID string) (*Agent, error) {
	o.agentLock.RLock()
	defer o.agentLock.RUnlock()
	
	agent, exists := o.Agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s does not exist", agentID)
	}
	
	return agent, nil
}

func (o *Orchestrator) GetAgentsByRole(role AgentRole) ([]*Agent, error) {
	o.agentLock.RLock()
	defer o.agentLock.RUnlock()
	
	agents, exists := o.AgentsByRole[role]
	if !exists {
		return nil, fmt.Errorf("no agents found for role %s", role)
	}
	
	return agents, nil
}

func (o *Orchestrator) CancelTask(taskID string) error {
	o.taskLock.Lock()
	defer o.taskLock.Unlock()
	
	task, exists := o.Tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s does not exist", taskID)
	}
	
	if task.Status != TaskStatusPending && task.Status != TaskStatusAssigned && task.Status != TaskStatusInProgress {
		return fmt.Errorf("task %s is not pending, assigned, or in progress", taskID)
	}
	
	task.Status = TaskStatusCancelled
	task.CompletedAt = time.Now().UTC()
	
	if o.Config.EventStream != nil {
		event := models.NewEvent(
			models.EventTypeTask,
			models.EventSourceSystem,
			map[string]interface{}{
				"action":     "task_cancelled",
				"task_id":    task.ID,
				"task_name":  task.Name,
				"agent_id":   task.AgentID,
			},
			map[string]string{
				"task_id":  task.ID,
				"agent_id": task.AgentID,
			},
		)
		
		if err := o.Config.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return nil
}

func (o *Orchestrator) CreateStandardAgents() error {
	frontendAgent, err := CreateAgentWithRoleDefinition(FrontendAgentRole, o.Config.EventStream)
	if err != nil {
		return fmt.Errorf("failed to create frontend agent: %v", err)
	}
	o.RegisterAgent(frontendAgent)
	
	appBuilderAgent, err := CreateAgentWithRoleDefinition(AppBuilderAgentRole, o.Config.EventStream)
	if err != nil {
		return fmt.Errorf("failed to create app builder agent: %v", err)
	}
	o.RegisterAgent(appBuilderAgent)
	
	codegenAgent, err := CreateAgentWithRoleDefinition(CodegenAgentRole, o.Config.EventStream)
	if err != nil {
		return fmt.Errorf("failed to create codegen agent: %v", err)
	}
	o.RegisterAgent(codegenAgent)
	
	engineeringAgent, err := CreateAgentWithRoleDefinition(EngineeringAgentRole, o.Config.EventStream)
	if err != nil {
		return fmt.Errorf("failed to create engineering agent: %v", err)
	}
	o.RegisterAgent(engineeringAgent)
	
	orchestratorAgent, err := CreateAgentWithRoleDefinition(OrchestratorAgentRole, o.Config.EventStream)
	if err != nil {
		return fmt.Errorf("failed to create orchestrator agent: %v", err)
	}
	o.RegisterAgent(orchestratorAgent)
	
	return o.CommunicationMgr.CreateStandardCommunicationChannels(context.Background(), []*Agent{
		frontendAgent,
		appBuilderAgent,
		codegenAgent,
		engineeringAgent,
		orchestratorAgent,
	})
}

func (o *Orchestrator) ExecuteWorkflow(ctx context.Context, workflowName string, inputs map[string]interface{}) (map[string]interface{}, error) {
	o.agentLock.RLock()
	var orchestratorAgent *Agent
	for _, agent := range o.Agents {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	o.agentLock.RUnlock()
	
	if orchestratorAgent == nil {
		return nil, fmt.Errorf("orchestrator agent not found")
	}
	
	workflowInputs := map[string]interface{}{
		"workflow_name": workflowName,
		"inputs":        inputs,
	}
	
	return orchestratorAgent.Process(ctx, workflowInputs)
}

func CreateMultiAgentSystem(name, description string, eventStream EventStream) (*Orchestrator, error) {
	config := OrchestratorConfig{
		Name:               name,
		Description:        description,
		MaxConcurrentTasks: 10,
		DefaultTimeout:     5 * time.Minute,
		EventStream:        eventStream,
		Metadata:           make(map[string]interface{}),
	}
	
	orchestrator := NewOrchestrator(config)
	
	if err := orchestrator.CreateStandardAgents(); err != nil {
		return nil, fmt.Errorf("failed to create standard agents: %v", err)
	}
	
	return orchestrator, nil
}
