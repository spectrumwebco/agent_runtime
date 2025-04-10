package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Loop struct {
	agent       *Agent
	modules     *modules.Registry
	tools       *tools.Registry
	state       *LoopState
	mutex       sync.RWMutex
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

type LoopState struct {
	Phase         string                 `json:"phase"` // Current phase (e.g., query_model, execute_action)
	Task          string                 `json:"task"`  // Description of the overall task
	StartTime     time.Time              `json:"startTime"`
	Events        []Event                `json:"events"` // Log of significant events
	Context       map[string]interface{} `json:"context"` // General context data
	CurrentAction string                 `json:"currentAction,omitempty"` // The action being executed
	LastError     string                 `json:"lastError,omitempty"`   // Last error encountered
	StepCount     int                    `json:"stepCount"`             // Number of steps taken
}

type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func NewLoop(agent *Agent, moduleRegistry *modules.Registry, toolRegistry *tools.Registry) *Loop {
	return &Loop{
		agent:       agent,
		modules:     moduleRegistry,
		tools:       toolRegistry,
		state:       &LoopState{
			Phase:     "idle",
			StartTime: time.Now(),
			Events:    make([]Event, 0),
			Context:   make(map[string]interface{}),
			StepCount: 0,
		},
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}
}

func (l *Loop) Start(ctx context.Context, task string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.state.Phase != "idle" {
		return fmt.Errorf("agent loop is already running")
	}
	
	l.state.Phase = "starting"
	l.state.Task = task
	l.state.StartTime = time.Now()
	l.state.Events = make([]Event, 0)
	l.state.Context = make(map[string]interface{})
	
	go l.run(ctx)
	
	return nil
}

func (l *Loop) Stop() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.state.Phase == "idle" {
		return fmt.Errorf("agent loop is not running")
	}
	
	close(l.stopChan)
	
	<-l.stoppedChan
	
	l.stopChan = make(chan struct{})
	l.stoppedChan = make(chan struct{})
	
	return nil
}

func (l *Loop) Status() *LoopState {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	
	return l.state
}

func (l *Loop) AddEvent(eventType string, data map[string]interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
	
	l.state.Events = append(l.state.Events, event)
}

func (l *Loop) run(ctx context.Context) {
	defer close(l.stoppedChan)
	
	l.mutex.Lock()
	l.state.Phase = "setup" // Start with setup phase
	l.mutex.Unlock()
	
	l.AddEvent("task_started", map[string]interface{}{
		"task": l.state.Task,
	})
	
	for {
		select {
		case <-l.stopChan:
			l.mutex.Lock()
			l.state.Phase = "idle"
			l.mutex.Unlock()
			
			l.AddEvent("task_stopped", map[string]interface{}{
				"reason": "manual_stop",
			})
			
			return
		case <-ctx.Done():
			l.mutex.Lock()
			l.state.Phase = "idle"
			l.mutex.Unlock()
			
			l.AddEvent("task_stopped", map[string]interface{}{
				"reason": "context_cancelled",
			})
			
			return
		default:
		}
		
		nextPhase, err := l.executePhase(ctx)
		if err != nil {
			l.mutex.Lock()
			l.state.Phase = "error_handling"
			l.mutex.Unlock()
			
			l.mutex.Lock()
			l.state.LastError = err.Error()
			l.mutex.Unlock()
			l.AddEvent("phase_error", map[string]interface{}{
				"phase": l.state.Phase,
				"error": err.Error(),
			})
			
			nextPhase, err = l.handleError(ctx, err)
			if err != nil {
				l.mutex.Lock()
				l.state.Phase = "idle"
				l.mutex.Unlock()
				
				l.AddEvent("task_failed", map[string]interface{}{
					"error": err.Error(),
				})
				
				return
			}
		}
		
		if nextPhase == "idle" {
			l.mutex.Lock()
			l.state.Phase = "idle"
			l.mutex.Unlock()
			
			l.AddEvent("task_completed", map[string]interface{}{})
			
			return
		}
		
		l.mutex.Lock()
		l.state.Phase = nextPhase
		l.mutex.Unlock()
		
		l.AddEvent("phase_transition", map[string]interface{}{
			"from": l.state.Phase,
			"to":   nextPhase,
		})
	}
}

func (l *Loop) executePhase(ctx context.Context) (string, error) {
	l.mutex.RLock()
	phase := l.state.Phase
	l.mutex.RUnlock()
	
	switch phase {
	case "setup":
		return l.executeSetup(ctx)
	case "query_model":
		return l.executeQueryModel(ctx)
	case "parse_action":
		return l.executeParseAction(ctx)
	case "execute_action":
		return l.executeExecuteAction(ctx)
	case "handle_observation":
		return l.executeHandleObservation(ctx)
	case "error_handling":
		return l.state.Phase, nil // Stay in error handling until resolved
	case "idle":
		return "idle", nil // Should not execute phase when idle
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
}

func (l *Loop) handleError(ctx context.Context, err error) (string, error) {
	l.AddEvent("error_handling", map[string]interface{}{
		"error": err.Error(),
	})
	
	
	return "idle", nil
}


func (l *Loop) executeSetup(ctx context.Context) (string, error) {
	l.AddEvent("setup_started", nil)
	l.agent.History = []Message{} // Reset history
	l.AddEvent("setup_completed", nil)
	return "query_model", nil // Move to querying the model first
}

func (l *Loop) executeQueryModel(ctx context.Context) (string, error) {
	l.AddEvent("query_model_started", nil)
	l.state.Context["last_model_output"] = map[string]string{
		"message": "placeholder model output with thought and action", // Placeholder
	}
	l.AddEvent("query_model_completed", map[string]interface{}{"output_length": len("placeholder...")})
	return "parse_action", nil
}

func (l *Loop) executeParseAction(ctx context.Context) (string, error) {
	l.AddEvent("parse_action_started", nil)
	modelOutput := l.state.Context["last_model_output"].(map[string]string)["message"]
	thought, action, err := l.tools.ParseAction(modelOutput) // Assuming ParseAction exists
	if err != nil {
		l.AddEvent("parse_action_failed", map[string]interface{}{"error": err.Error()})
		return "error_handling", fmt.Errorf("parsing action failed: %w", err) // Let main loop handle error phase
	}
	l.state.Context["current_thought"] = thought
	l.state.CurrentAction = action
	l.AddEvent("parse_action_completed", map[string]interface{}{"action": action})
	return "execute_action", nil
}

func (l *Loop) executeExecuteAction(ctx context.Context) (string, error) {
	l.mutex.Lock()
	action := l.state.CurrentAction
	l.mutex.Unlock()

	l.AddEvent("execute_action_started", map[string]interface{}{"action": action})

	if action == "exit" {
		l.AddEvent("exit_command_received", nil)
		return "idle", nil // Task completed by exit command
	}

	observation, err := l.tools.Execute(ctx, action, l.agent.Env) // Assuming Execute exists
	if err != nil {
		l.AddEvent("execute_action_failed", map[string]interface{}{"action": action, "error": err.Error()})
		return "error_handling", fmt.Errorf("executing action '%s' failed: %w", action, err)
	}

	l.mutex.Lock()
	l.state.Context["last_observation"] = observation
	l.state.StepCount++ // Increment step count after successful execution
	l.mutex.Unlock()

	l.AddEvent("execute_action_completed", map[string]interface{}{"action": action, "observation_length": len(observation)})
	return "handle_observation", nil
}

func (l *Loop) executeHandleObservation(ctx context.Context) (string, error) {
	l.AddEvent("handle_observation_started", nil)

	done := false // Check step output or other conditions

	l.AddEvent("handle_observation_completed", map[string]interface{}{"task_done": done})
	if done {
		return "idle", nil
	}
	return "query_model", nil // Loop back to query model for next step
}
