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

func (l *Loop) executeStep(ctx context.Context) (nextPhase string, err error) {
	l.mutex.Lock()
	l.state.StepCount++
	stepNum := l.state.StepCount
	l.mutex.Unlock()

	l.AddEvent("step_started", map[string]interface{}{"step": stepNum})

	l.AddEvent("query_model_started", map[string]interface{}{"step": stepNum})
	modelOutput, queryErr := l.agent.QueryModel(ctx, l.agent.History) // Use agent's method
	if queryErr != nil {
		l.AddEvent("query_model_failed", map[string]interface{}{"step": stepNum, "error": queryErr.Error()})
		l.state.LastError = fmt.Sprintf("querying model failed: %v", queryErr)
		return "error_handling", queryErr
	}
	if modelOutput == nil {
		modelOutput = map[string]interface{}{
			"message": fmt.Sprintf("Placeholder model output for step %d with thought and action", stepNum),
		}
	}
	l.mutex.Lock()
	l.state.Context["last_model_output"] = modelOutput
	l.mutex.Unlock()
	l.AddEvent("query_model_completed", map[string]interface{}{"step": stepNum, "output_length": len(modelOutput["message"].(string))})

	l.AddEvent("parse_action_started", map[string]interface{}{"step": stepNum})
	thought, action, parseErr := l.tools.ParseAction(modelOutput["message"].(string))
	if parseErr != nil {
		l.AddEvent("parse_action_failed", map[string]interface{}{"step": stepNum, "error": parseErr.Error()})
		l.state.LastError = fmt.Sprintf("parsing action failed: %v", parseErr)
		return "error_handling", parseErr // Go to error handling for now
	}
	l.mutex.Lock()
	l.state.Context["current_thought"] = thought
	l.state.CurrentAction = action
	l.mutex.Unlock()
	l.AddEvent("parse_action_completed", map[string]interface{}{"step": stepNum, "action": action})

	l.AddEvent("execute_action_started", map[string]interface{}{"step": stepNum, "action": action})

	if action == "exit" { // Simple exit handling
		l.AddEvent("exit_command_received", map[string]interface{}{"step": stepNum})
		return "idle", nil // Task completed by exit command
	}
	if action == "submit" { // Simple submit handling
		l.AddEvent("submit_command_received", map[string]interface{}{"step": stepNum})
		l.state.Context["submission"] = "Placeholder submission content" // Placeholder
		return "idle", nil // Task completed by submit command
	}

	execCtx, cancel := context.WithTimeout(ctx, l.tools.Config.ExecutionTimeout) // Use configured timeout
	defer cancel()
	observation, execErr := l.tools.ExecuteAction(execCtx, action, l.agent.Env)
	if execErr != nil {
		l.AddEvent("execute_action_failed", map[string]interface{}{"step": stepNum, "action": action, "error": execErr.Error()})
		l.state.LastError = fmt.Sprintf("executing action '%s' failed: %v", action, execErr)
		return "error_handling", execErr
	}
	l.mutex.Lock()
	l.state.Context["last_observation"] = observation
	l.mutex.Unlock()
	l.AddEvent("execute_action_completed", map[string]interface{}{"step": stepNum, "action": action, "observation_length": len(observation)})

	l.AddEvent("handle_observation_started", map[string]interface{}{"step": stepNum})

	l.agent.AddStepToHistory(thought, action, modelOutput, observation)

	done := false // Placeholder

	l.AddEvent("step_completed", map[string]interface{}{"step": stepNum, "task_done": done})
	if done {
		return "idle", nil
	}

	return "execute_step", nil
}

func (l *Loop) executePhase(ctx context.Context) (string, error) {
	l.mutex.RLock()
	phase := l.state.Phase
	l.mutex.RUnlock()

	switch phase {
	case "setup":
		return l.executeSetup(ctx)
	case "execute_step": // New phase for the main loop cycle
		return l.executeStep(ctx)
	case "error_handling":
		l.mutex.RLock()
		lastErr := l.state.LastError
		l.mutex.RUnlock()
		return l.handleError(ctx, fmt.Errorf(lastErr)) // Pass the stored error
	case "idle":
		return "idle", nil // Stay idle
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
}

func (l *Loop) handleError(ctx context.Context, err error) (string, error) {
	l.AddEvent("error_handling_started", map[string]interface{}{
		"error": err.Error(),
	})

	fmt.Printf("Error encountered: %v. Stopping loop.\n", err)
	l.AddEvent("error_handling_failed", map[string]interface{}{"error": err.Error()})
	return "idle", err // Return the original error to signal failure
}

func (l *Loop) executeSetup(ctx context.Context) (string, error) {
	l.AddEvent("setup_started", nil)
	l.mutex.Lock()
	l.agent.History = make([]Message, 0)
	l.state.StepCount = 0
	l.state.LastError = ""
	l.state.CurrentAction = ""
	delete(l.state.Context, "last_model_output")
	delete(l.state.Context, "last_observation")
	delete(l.state.Context, "current_thought")
	delete(l.state.Context, "submission")
	l.mutex.Unlock()


	l.AddEvent("setup_completed", nil)
	return "execute_step", nil // Start the main execution loop
}
