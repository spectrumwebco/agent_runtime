package rex

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type TrajectoryState string

const (
	StateIdle TrajectoryState = "idle"
	StateRunning TrajectoryState = "running"
	StatePaused TrajectoryState = "paused"
	StateCompleted TrajectoryState = "completed"
	StateFailed TrajectoryState = "failed"
)

type TrajectoryConfig struct {
	Name string `json:"name"`
	
	Description string `json:"description"`
	
	Steps []TrajectoryStep `json:"steps"`
	
	MaxRetries int `json:"max_retries"`
	
	RetryDelay time.Duration `json:"retry_delay"`
}

type TrajectoryStep struct {
	ID string `json:"id"`
	
	Name string `json:"name"`
	
	Description string `json:"description"`
	
	Tool string `json:"tool"`
	
	Params map[string]interface{} `json:"params"`
	
	NextSteps []string `json:"next_steps"`
	
	Condition string `json:"condition"`
}

type Trajectory struct {
	Runtime *Runtime
	
	ID string
	
	Config TrajectoryConfig
	
	State TrajectoryState
	
	CurrentStep string
	
	Context map[string]interface{}
	
	Results map[string]interface{}
	
	Error error
	
	RetryCount int
	
	mutex sync.RWMutex
	
	doneCh chan struct{}
}

func NewTrajectory(runtime *Runtime, id string, config TrajectoryConfig) *Trajectory {
	return &Trajectory{
		Runtime:     runtime,
		ID:          id,
		Config:      config,
		State:       StateIdle,
		Context:     make(map[string]interface{}),
		Results:     make(map[string]interface{}),
		doneCh:      make(chan struct{}),
	}
}

func (t *Trajectory) Start(ctx context.Context) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	if t.State == StateRunning {
		return fmt.Errorf("trajectory is already running")
	}
	
	t.State = StateRunning
	t.CurrentStep = t.getInitialStep()
	
	go func() {
		t.run(ctx)
	}()
	
	return nil
}

func (t *Trajectory) Stop() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	if t.State != StateRunning && t.State != StatePaused {
		return fmt.Errorf("trajectory is not running or paused")
	}
	
	t.State = StateIdle
	
	close(t.doneCh)
	
	return nil
}

func (t *Trajectory) Pause() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	if t.State != StateRunning {
		return fmt.Errorf("trajectory is not running")
	}
	
	t.State = StatePaused
	
	return nil
}

func (t *Trajectory) Resume(ctx context.Context) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	if t.State != StatePaused {
		return fmt.Errorf("trajectory is not paused")
	}
	
	t.State = StateRunning
	
	go func() {
		t.run(ctx)
	}()
	
	return nil
}

func (t *Trajectory) Wait(timeout time.Duration) error {
	select {
	case <-t.doneCh:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for trajectory to complete")
	}
}

func (t *Trajectory) GetState() TrajectoryState {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	return t.State
}

func (t *Trajectory) GetCurrentStep() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	return t.CurrentStep
}

func (t *Trajectory) GetResults() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	resultsCopy := make(map[string]interface{})
	for k, v := range t.Results {
		resultsCopy[k] = v
	}
	
	return resultsCopy
}

func (t *Trajectory) GetError() error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	return t.Error
}

func (t *Trajectory) GetContext() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	contextCopy := make(map[string]interface{})
	for k, v := range t.Context {
		contextCopy[k] = v
	}
	
	return contextCopy
}

func (t *Trajectory) SetContext(ctx map[string]interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	t.Context = ctx
}

func (t *Trajectory) UpdateContext(updates map[string]interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	for k, v := range updates {
		t.Context[k] = v
	}
}

func (t *Trajectory) SaveState() ([]byte, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	
	state := map[string]interface{}{
		"id":           t.ID,
		"state":        t.State,
		"current_step": t.CurrentStep,
		"context":      t.Context,
		"results":      t.Results,
		"retry_count":  t.RetryCount,
	}
	
	return json.Marshal(state)
}

func (t *Trajectory) LoadState(data []byte) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	
	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}
	
	if id, ok := state["id"].(string); ok {
		t.ID = id
	}
	
	if stateStr, ok := state["state"].(string); ok {
		t.State = TrajectoryState(stateStr)
	}
	
	if currentStep, ok := state["current_step"].(string); ok {
		t.CurrentStep = currentStep
	}
	
	if context, ok := state["context"].(map[string]interface{}); ok {
		t.Context = context
	}
	
	if results, ok := state["results"].(map[string]interface{}); ok {
		t.Results = results
	}
	
	if retryCount, ok := state["retry_count"].(float64); ok {
		t.RetryCount = int(retryCount)
	}
	
	return nil
}

func (t *Trajectory) run(ctx context.Context) {
	
	close(t.doneCh)
}

func (t *Trajectory) getInitialStep() string {
	if len(t.Config.Steps) == 0 {
		return ""
	}
	
	return t.Config.Steps[0].ID
}

func (t *Trajectory) getStep(id string) (*TrajectoryStep, error) {
	for _, step := range t.Config.Steps {
		if step.ID == id {
			return &step, nil
		}
	}
	
	return nil, fmt.Errorf("step not found: %s", id)
}

func (t *Trajectory) executeStep(ctx context.Context, step *TrajectoryStep) (interface{}, error) {
	result, err := t.Runtime.ExecuteTool(ctx, step.Tool, step.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tool %s: %w", step.Tool, err)
	}
	
	return result, nil
}

func (t *Trajectory) getNextStep(currentStep *TrajectoryStep, result interface{}) (string, error) {
	if len(currentStep.NextSteps) == 0 {
		return "", nil
	}
	
	if len(currentStep.NextSteps) == 1 {
		return currentStep.NextSteps[0], nil
	}
	
	if currentStep.Condition != "" {
		return currentStep.NextSteps[0], nil
	}
	
	return currentStep.NextSteps[0], nil
}
