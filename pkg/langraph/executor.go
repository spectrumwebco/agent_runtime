package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ExecutionStatus string

const (
	ExecutionStatusPending ExecutionStatus = "pending"
	ExecutionStatusRunning ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

type ExecutionResult struct {
	NodeID  NodeID                  `json:"node_id"`
	Outputs map[string]interface{}  `json:"outputs"`
	Error   string                  `json:"error,omitempty"`
	Status  ExecutionStatus         `json:"status"`
	Metrics map[string]interface{}  `json:"metrics,omitempty"`
}

type Execution struct {
	ID           string                       `json:"id"`
	GraphID      string                       `json:"graph_id"`
	StartNodeID  NodeID                       `json:"start_node_id"`
	Status       ExecutionStatus              `json:"status"`
	Results      map[NodeID]*ExecutionResult  `json:"results"`
	StartTime    time.Time                    `json:"start_time"`
	EndTime      time.Time                    `json:"end_time,omitempty"`
	Inputs       map[string]interface{}       `json:"inputs"`
	Metadata     map[string]interface{}       `json:"metadata,omitempty"`
	Error        string                       `json:"error,omitempty"`
	resultsMutex sync.RWMutex                 `json:"-"`
}

type Executor struct {
	executions     map[string]*Execution
	executionLock  sync.RWMutex
	graph          *Graph
}

func NewExecutor(graph *Graph) *Executor {
	return &Executor{
		executions: make(map[string]*Execution),
		graph:      graph,
	}
}

func (e *Executor) Execute(ctx context.Context, startNodeID NodeID, inputs map[string]interface{}, metadata map[string]interface{}) (*Execution, error) {
	if _, err := e.graph.GetNode(startNodeID); err != nil {
		return nil, err
	}

	execution := &Execution{
		ID:          uuid.New().String(),
		GraphID:     e.graph.ID,
		StartNodeID: startNodeID,
		Status:      ExecutionStatusPending,
		Results:     make(map[NodeID]*ExecutionResult),
		StartTime:   time.Now().UTC(),
		Inputs:      inputs,
		Metadata:    metadata,
	}

	e.executionLock.Lock()
	e.executions[execution.ID] = execution
	e.executionLock.Unlock()

	go func() {
		e.executeGraph(ctx, execution)
	}()

	return execution, nil
}

func (e *Executor) executeGraph(ctx context.Context, execution *Execution) {
	execution.Status = ExecutionStatusRunning

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results, err := e.graph.TraverseAsync(ctx, execution.StartNodeID, execution.Inputs)
	
	execution.resultsMutex.Lock()
	defer execution.resultsMutex.Unlock()

	execution.EndTime = time.Now().UTC()

	if err != nil {
		execution.Status = ExecutionStatusFailed
		execution.Error = err.Error()
		return
	}

	for nodeID, outputs := range results {
		execution.Results[nodeID] = &ExecutionResult{
			NodeID:  nodeID,
			Outputs: outputs,
			Status:  ExecutionStatusCompleted,
			Metrics: map[string]interface{}{
				"execution_time": time.Since(execution.StartTime).Milliseconds(),
			},
		}
	}

	execution.Status = ExecutionStatusCompleted
}

func (e *Executor) GetExecution(id string) (*Execution, error) {
	e.executionLock.RLock()
	defer e.executionLock.RUnlock()

	execution, exists := e.executions[id]
	if !exists {
		return nil, fmt.Errorf("execution %s does not exist", id)
	}

	return execution, nil
}

func (e *Executor) CancelExecution(id string) error {
	e.executionLock.Lock()
	defer e.executionLock.Unlock()

	execution, exists := e.executions[id]
	if !exists {
		return fmt.Errorf("execution %s does not exist", id)
	}

	if execution.Status != ExecutionStatusRunning && execution.Status != ExecutionStatusPending {
		return fmt.Errorf("execution %s is not running or pending", id)
	}

	execution.Status = ExecutionStatusCancelled
	execution.EndTime = time.Now().UTC()

	return nil
}

func (e *Executor) GetExecutionResult(executionID string, nodeID NodeID) (*ExecutionResult, error) {
	execution, err := e.GetExecution(executionID)
	if err != nil {
		return nil, err
	}

	execution.resultsMutex.RLock()
	defer execution.resultsMutex.RUnlock()

	result, exists := execution.Results[nodeID]
	if !exists {
		return nil, fmt.Errorf("result for node %s does not exist in execution %s", nodeID, executionID)
	}

	return result, nil
}

func (e *Executor) GetExecutionResults(executionID string) (map[NodeID]*ExecutionResult, error) {
	execution, err := e.GetExecution(executionID)
	if err != nil {
		return nil, err
	}

	execution.resultsMutex.RLock()
	defer execution.resultsMutex.RUnlock()

	resultsCopy := make(map[NodeID]*ExecutionResult)
	for nodeID, result := range execution.Results {
		resultsCopy[nodeID] = result
	}

	return resultsCopy, nil
}

func (e *Executor) WaitForExecution(ctx context.Context, executionID string) (*Execution, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			execution, err := e.GetExecution(executionID)
			if err != nil {
				return nil, err
			}

			if execution.Status == ExecutionStatusCompleted || 
			   execution.Status == ExecutionStatusFailed || 
			   execution.Status == ExecutionStatusCancelled {
				return execution, nil
			}
		}
	}
}

func (e *Executor) ListExecutions() []*Execution {
	e.executionLock.RLock()
	defer e.executionLock.RUnlock()

	var executions []*Execution
	for _, execution := range e.executions {
		executions = append(executions, execution)
	}

	return executions
}
