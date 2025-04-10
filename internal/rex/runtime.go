package rex

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Runtime struct {
	Agent *agent.Agent
	
	Trajectories map[string]*Trajectory
	
	mutex sync.RWMutex
	
	Context map[string]interface{}
	
	running bool
	
	stopCh chan struct{}
}

type RuntimeConfig struct {
	Name string `json:"name"`
	
	Description string `json:"description"`
	
	TrajectoryPath string `json:"trajectory_path"`
}

func NewRuntime(agent *agent.Agent, config RuntimeConfig) (*Runtime, error) {
	runtime := &Runtime{
		Agent:        agent,
		Trajectories: make(map[string]*Trajectory),
		Context:      make(map[string]interface{}),
		stopCh:       make(chan struct{}),
	}
	
	if err := runtime.LoadTrajectories(config.TrajectoryPath); err != nil {
		return nil, fmt.Errorf("failed to load trajectories: %w", err)
	}
	
	return runtime, nil
}

func (r *Runtime) LoadTrajectories(path string) error {
	return nil
}

func (r *Runtime) Start(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if r.running {
		return fmt.Errorf("runtime is already running")
	}
	
	r.running = true
	
	go func() {
		r.run(ctx)
	}()
	
	return nil
}

func (r *Runtime) Stop() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if !r.running {
		return fmt.Errorf("runtime is not running")
	}
	
	close(r.stopCh)
	
	r.running = false
	
	return nil
}

func (r *Runtime) run(ctx context.Context) {
}

func (r *Runtime) CreateTrajectory(id string, config TrajectoryConfig) (*Trajectory, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.Trajectories[id]; exists {
		return nil, fmt.Errorf("trajectory already exists: %s", id)
	}
	
	trajectory := NewTrajectory(r, id, config)
	r.Trajectories[id] = trajectory
	
	return trajectory, nil
}

func (r *Runtime) GetTrajectory(id string) (*Trajectory, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	trajectory, exists := r.Trajectories[id]
	if !exists {
		return nil, fmt.Errorf("trajectory not found: %s", id)
	}
	
	return trajectory, nil
}

func (r *Runtime) ListTrajectories() []*Trajectory {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	trajectories := make([]*Trajectory, 0, len(r.Trajectories))
	for _, trajectory := range r.Trajectories {
		trajectories = append(trajectories, trajectory)
	}
	
	return trajectories
}

func (r *Runtime) DeleteTrajectory(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.Trajectories[id]; !exists {
		return fmt.Errorf("trajectory not found: %s", id)
	}
	
	delete(r.Trajectories, id)
	
	return nil
}

func (r *Runtime) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	return r.Agent.ExecuteTool(ctx, name, params)
}

func (r *Runtime) GetContext() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	contextCopy := make(map[string]interface{})
	for k, v := range r.Context {
		contextCopy[k] = v
	}
	
	return contextCopy
}

func (r *Runtime) SetContext(ctx map[string]interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.Context = ctx
}

func (r *Runtime) UpdateContext(updates map[string]interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for k, v := range updates {
		r.Context[k] = v
	}
}

func (r *Runtime) SaveState() ([]byte, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	state := map[string]interface{}{
		"context":      r.Context,
		"trajectories": r.Trajectories,
	}
	
	return json.Marshal(state)
}

func (r *Runtime) LoadState(data []byte) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}
	
	if context, ok := state["context"].(map[string]interface{}); ok {
		r.Context = context
	}
	
	
	return nil
}

func (r *Runtime) WaitForTrajectory(id string, timeout time.Duration) error {
	trajectory, err := r.GetTrajectory(id)
	if err != nil {
		return err
	}
	
	return trajectory.Wait(timeout)
}

func (r *Runtime) RegisterTool(tool tools.Tool) error {
	return r.Agent.RegisterTool(tool)
}

func (r *Runtime) GetTool(name string) (tools.Tool, error) {
	return r.Agent.GetTool(name)
}

func (r *Runtime) ListTools() []tools.Tool {
	return r.Agent.ListTools()
}
