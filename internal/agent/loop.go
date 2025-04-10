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
	Phase       string                 `json:"phase"`
	Task        string                 `json:"task"`
	StartTime   time.Time              `json:"startTime"`
	Events      []Event                `json:"events"`
	Context     map[string]interface{} `json:"context"`
	CurrentTool string                 `json:"currentTool"`
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
	l.state.Phase = "task_initiation"
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
	case "task_initiation":
		return l.executeTaskInitiation(ctx)
	case "analyze_requirements":
		return l.executeAnalyzeRequirements(ctx)
	case "tool_discovery":
		return l.executeToolDiscovery(ctx)
	case "planning_phase":
		return l.executePlanningPhase(ctx)
	case "execution_phase":
		return l.executeExecutionPhase(ctx)
	case "tool_transition_evaluation":
		return l.executeToolTransitionEvaluation(ctx)
	case "toolbelt_activation":
		return l.executeToolbeltActivation(ctx)
	case "completion_verification":
		return l.executeCompletionVerification(ctx)
	case "state_cleanup":
		return l.executeStateCleanup(ctx)
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

func (l *Loop) executeTaskInitiation(ctx context.Context) (string, error) {
	
	return "analyze_requirements", nil
}

func (l *Loop) executeAnalyzeRequirements(ctx context.Context) (string, error) {
	
	return "tool_discovery", nil
}

func (l *Loop) executeToolDiscovery(ctx context.Context) (string, error) {
	
	return "planning_phase", nil
}

func (l *Loop) executePlanningPhase(ctx context.Context) (string, error) {
	
	return "execution_phase", nil
}

func (l *Loop) executeExecutionPhase(ctx context.Context) (string, error) {
	
	return "tool_transition_evaluation", nil
}

func (l *Loop) executeToolTransitionEvaluation(ctx context.Context) (string, error) {
	
	return "completion_verification", nil
}

func (l *Loop) executeToolbeltActivation(ctx context.Context) (string, error) {
	
	return "execution_phase", nil
}

func (l *Loop) executeCompletionVerification(ctx context.Context) (string, error) {
	
	return "state_cleanup", nil
}

func (l *Loop) executeStateCleanup(ctx context.Context) (string, error) {
	
	return "idle", nil
}
