package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type ToolTier string

const (
	CoreTools      ToolTier = "core"      // Basic tools provided by the system
	ToolchainTools ToolTier = "toolchain" // Specialized engineering and data science tools
	ToolbeltTools  ToolTier = "toolbelt"  // Advanced MCP tools discovered dynamically
)

type StateType string

const (
	ExecutionState StateType = "execution" // Primary task execution state (TTL: 3600s)
	ToolState      StateType = "tool"      // Tool configuration state (TTL: 7200s)
	ContextState   StateType = "context"   // Long-term contextual state (TTL: 86400s)
)

type Loop struct {
	agent         *Agent
	modules       *modules.Registry
	tools         *tools.Registry
	state         *LoopState
	mutex         sync.RWMutex
	stopChan      chan struct{}
	stoppedChan   chan struct{}
	currentTier   ToolTier
	stateManager  map[StateType]map[string]interface{}
	nestedStates  []map[string]interface{}
	nestedDepth   int
	maxNestedDepth int
}

type LoopState struct {
	Phase         string                 `json:"phase"` // Current phase (e.g., task_initiation, analyze_requirements)
	Task          string                 `json:"task"`  // Description of the overall task
	StartTime     time.Time              `json:"startTime"`
	Events        []Event                `json:"events"` // Log of significant events
	Context       map[string]interface{} `json:"context"` // General context data
	CurrentAction string                 `json:"currentAction,omitempty"` // The action being executed
	LastError     string                 `json:"lastError,omitempty"`   // Last error encountered
	StepCount     int                    `json:"stepCount"`             // Number of steps taken
	ToolTier      ToolTier               `json:"toolTier"`              // Current tool tier being used
}

type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func NewLoop(agent *Agent, moduleRegistry *modules.Registry, toolRegistry *tools.Registry) *Loop {
	return &Loop{
		agent:         agent,
		modules:       moduleRegistry,
		tools:         toolRegistry,
		state:         &LoopState{
			Phase:     "idle",
			StartTime: time.Now(),
			Events:    make([]Event, 0),
			Context:   make(map[string]interface{}),
			StepCount: 0,
			ToolTier:  CoreTools,
		},
		stopChan:      make(chan struct{}),
		stoppedChan:   make(chan struct{}),
		currentTier:   CoreTools,
		stateManager:  make(map[StateType]map[string]interface{}),
		nestedStates:  make([]map[string]interface{}, 0),
		nestedDepth:   0,
		maxNestedDepth: 5, // Maximum depth for nested states
	}
}

func (l *Loop) Start(ctx context.Context, task string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.state.Phase != "idle" {
		return fmt.Errorf("agent loop is already running")
	}
	
	l.stateManager[ExecutionState] = make(map[string]interface{})
	l.stateManager[ToolState] = make(map[string]interface{})
	l.stateManager[ContextState] = make(map[string]interface{})
	
	l.state.Phase = "task_initiation" // Start with task initiation phase
	l.state.Task = task
	l.state.StartTime = time.Now()
	l.state.Events = make([]Event, 0)
	l.state.Context = make(map[string]interface{})
	l.state.ToolTier = CoreTools // Start with core tools
	
	l.stateManager[ContextState]["task"] = task
	l.stateManager[ContextState]["agent_name"] = "samsepi0l"
	
	go l.run(ctx)
	
	return nil
}

func (l *Loop) Stop() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.state.Phase == "idle" {
		return fmt.Errorf("agent loop is not running")
	}
	
	l.AddEvent("state_cleanup_initiated", map[string]interface{}{
		"reason": "manual_stop",
	})
	
	close(l.stopChan)
	
	<-l.stoppedChan
	
	l.stopChan = make(chan struct{})
	l.stoppedChan = make(chan struct{})
	
	l.stateManager = make(map[StateType]map[string]interface{})
	l.nestedStates = make([]map[string]interface{}, 0)
	l.nestedDepth = 0
	
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
	
	if eventType == "task_started" || eventType == "task_completed" || 
	   eventType == "task_failed" || eventType == "tool_transition" ||
	   eventType == "state_change" {
		l.stateManager[ExecutionState][fmt.Sprintf("event_%d", len(l.state.Events))] = event
	}
}

func (l *Loop) PushState() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.nestedDepth >= l.maxNestedDepth {
		return fmt.Errorf("maximum nested state depth reached (%d)", l.maxNestedDepth)
	}
	
	stateCopy := make(map[string]interface{})
	for k, v := range l.state.Context {
		stateCopy[k] = v
	}
	
	l.nestedStates = append(l.nestedStates, stateCopy)
	l.nestedDepth++
	
	l.AddEvent("state_pushed", map[string]interface{}{
		"depth": l.nestedDepth,
	})
	
	return nil
}

func (l *Loop) PopState() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.nestedDepth <= 0 {
		return fmt.Errorf("no nested state to pop")
	}
	
	previousState := l.nestedStates[len(l.nestedStates)-1]
	l.nestedStates = l.nestedStates[:len(l.nestedStates)-1]
	l.nestedDepth--
	
	for k, v := range l.state.Context {
		if k == "last_observation" || k == "last_model_output" || k == "current_thought" {
			previousState[k] = v
		}
	}
	
	l.state.Context = previousState
	
	l.AddEvent("state_popped", map[string]interface{}{
		"depth": l.nestedDepth,
	})
	
	return nil
}

func (l *Loop) TransitionToolTier(targetTier ToolTier) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.state.ToolTier == targetTier {
		return nil // Already in the target tier
	}
	
	previousTier := l.state.ToolTier
	l.state.ToolTier = targetTier
	l.currentTier = targetTier
	
	l.stateManager[ToolState]["previous_tier"] = previousTier
	
	l.AddEvent("tool_transition", map[string]interface{}{
		"from": previousTier,
		"to":   targetTier,
	})
	
	return nil
}

func (l *Loop) run(ctx context.Context) {
	defer close(l.stoppedChan)
	
	for _, stateType := range []StateType{ExecutionState, ToolState, ContextState} {
		if l.stateManager[stateType] == nil {
			l.stateManager[stateType] = make(map[string]interface{})
		}
	}
	
	l.mutex.Lock()
	l.state.Phase = "task_initiation" // Start with task initiation phase
	l.mutex.Unlock()
	
	l.AddEvent("task_started", map[string]interface{}{
		"task": l.state.Task,
		"agent": "samsepi0l",
	})
	
	for {
		select {
		case <-l.stopChan:
			l.mutex.Lock()
			l.state.Phase = "state_cleanup"
			l.mutex.Unlock()
			
			_, cleanupErr := l.executePhase(ctx)
			if cleanupErr != nil {
				l.AddEvent("cleanup_error", map[string]interface{}{
					"error": cleanupErr.Error(),
				})
			}
			
			l.mutex.Lock()
			l.state.Phase = "idle"
			l.mutex.Unlock()
			
			l.AddEvent("task_stopped", map[string]interface{}{
				"reason": "manual_stop",
			})
			
			return
		case <-ctx.Done():
			l.mutex.Lock()
			l.state.Phase = "state_cleanup"
			l.mutex.Unlock()
			
			_, cleanupErr := l.executePhase(ctx)
			if cleanupErr != nil {
				l.AddEvent("cleanup_error", map[string]interface{}{
					"error": cleanupErr.Error(),
				})
			}
			
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
			l.state.LastError = err.Error()
			l.mutex.Unlock()
			
			l.AddEvent("phase_error", map[string]interface{}{
				"phase": l.state.Phase,
				"error": err.Error(),
			})
			
			l.stateManager[ExecutionState]["last_error"] = err.Error()
			l.stateManager[ExecutionState]["error_phase"] = l.state.Phase
			
			nextPhase, err = l.handleError(ctx, err)
			if err != nil {
				l.mutex.Lock()
				l.state.Phase = "state_cleanup"
				l.mutex.Unlock()
				
				_, cleanupErr := l.executePhase(ctx)
				if cleanupErr != nil {
					l.AddEvent("cleanup_error", map[string]interface{}{
						"error": cleanupErr.Error(),
					})
				}
				
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
			l.state.Phase = "state_cleanup"
			l.mutex.Unlock()
			
			_, cleanupErr := l.executePhase(ctx)
			if cleanupErr != nil {
				l.AddEvent("cleanup_error", map[string]interface{}{
					"error": cleanupErr.Error(),
				})
			}
			
			l.mutex.Lock()
			l.state.Phase = "idle"
			l.mutex.Unlock()
			
			l.AddEvent("task_completed", map[string]interface{}{
				"agent": "samsepi0l",
				"steps_taken": l.state.StepCount,
			})
			
			return
		}
		
		if nextPhase == "tool_transition_evaluation" {
			l.evaluateToolTransition(ctx)
		}
		
		l.mutex.Lock()
		currentPhase := l.state.Phase
		l.state.Phase = nextPhase
		l.mutex.Unlock()
		
		l.AddEvent("phase_transition", map[string]interface{}{
			"from": currentPhase,
			"to":   nextPhase,
		})
	}
}

func (l *Loop) evaluateToolTransition(ctx context.Context) {
	l.mutex.RLock()
	currentTier := l.state.ToolTier
	task := l.state.Task
	l.mutex.RUnlock()
	
	isDataScienceTask := false
	dataScienceTriggers := []string{
		"data manipulation", "data science", "machine learning",
		"statistical analysis", "data visualization", "model training",
	}
	
	for _, trigger := range dataScienceTriggers {
		if strings.Contains(strings.ToLower(task), strings.ToLower(trigger)) {
			isDataScienceTask = true
			break
		}
	}
	
	if isDataScienceTask && currentTier == CoreTools {
		l.TransitionToolTier(ToolchainTools)
		l.AddEvent("tool_tier_transition_reason", map[string]interface{}{
			"reason": "data_science_task_detected",
			"triggers_matched": dataScienceTriggers,
		})
	}
	
	needsSpecializedTools := false
	if l.state.StepCount > 5 && isDataScienceTask {
		needsSpecializedTools = true
	}
	
	if needsSpecializedTools && currentTier == ToolchainTools {
		l.TransitionToolTier(ToolbeltTools)
		l.AddEvent("tool_tier_transition_reason", map[string]interface{}{
			"reason": "specialized_tools_needed",
			"step_count": l.state.StepCount,
		})
	}
}

func (l *Loop) executeStep(ctx context.Context) (nextPhase string, err error) {
	l.mutex.Lock()
	l.state.StepCount++
	stepNum := l.state.StepCount
	currentTier := l.state.ToolTier
	l.mutex.Unlock()

	l.AddEvent("step_started", map[string]interface{}{
		"step": stepNum,
		"agent": "samsepi0l",
		"tool_tier": currentTier,
	})

	if stepNum > 1 && stepNum%3 == 0 {
		return "tool_transition_evaluation", nil
	}

	l.AddEvent("query_model_started", map[string]interface{}{
		"step": stepNum,
		"tool_tier": currentTier,
	})
	
	modelOutput, queryErr := l.agent.QueryModel(ctx, l.agent.History)
	if queryErr != nil {
		l.AddEvent("query_model_failed", map[string]interface{}{
			"step": stepNum,
			"error": queryErr.Error(),
		})
		l.mutex.Lock()
		l.state.LastError = fmt.Sprintf("querying model failed: %v", queryErr)
		l.mutex.Unlock()
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
	
	l.stateManager[ExecutionState][fmt.Sprintf("model_output_%d", stepNum)] = modelOutput
	
	l.AddEvent("query_model_completed", map[string]interface{}{
		"step": stepNum,
		"output_length": len(modelOutput["message"].(string)),
	})

	l.AddEvent("parse_action_started", map[string]interface{}{
		"step": stepNum,
	})
	
	thought, action, parseErr := l.tools.ParseAction(modelOutput["message"].(string))
	if parseErr != nil {
		l.AddEvent("parse_action_failed", map[string]interface{}{
			"step": stepNum,
			"error": parseErr.Error(),
		})
		l.mutex.Lock()
		l.state.LastError = fmt.Sprintf("parsing action failed: %v", parseErr)
		l.mutex.Unlock()
		return "error_handling", parseErr
	}
	
	l.mutex.Lock()
	l.state.Context["current_thought"] = thought
	l.state.CurrentAction = action
	l.mutex.Unlock()
	
	l.stateManager[ExecutionState][fmt.Sprintf("thought_%d", stepNum)] = thought
	l.stateManager[ExecutionState][fmt.Sprintf("action_%d", stepNum)] = action
	
	l.AddEvent("parse_action_completed", map[string]interface{}{
		"step": stepNum,
		"action": action,
	})

	l.AddEvent("execute_action_started", map[string]interface{}{
		"step": stepNum,
		"action": action,
		"tool_tier": currentTier,
	})

	if action == "exit" {
		l.AddEvent("exit_command_received", map[string]interface{}{
			"step": stepNum,
		})
		return "completion_verification", nil
	}
	
	if action == "submit" {
		l.AddEvent("submit_command_received", map[string]interface{}{
			"step": stepNum,
		})
		l.mutex.Lock()
		l.state.Context["submission"] = modelOutput["message"]
		l.mutex.Unlock()
		l.stateManager[ExecutionState]["submission"] = modelOutput["message"]
		return "completion_verification", nil
	}
	
	if strings.HasPrefix(action, "toolchain.") && currentTier == CoreTools {
		l.TransitionToolTier(ToolchainTools)
		action = strings.TrimPrefix(action, "toolchain.")
	} else if strings.HasPrefix(action, "toolbelt.") && currentTier != ToolbeltTools {
		l.TransitionToolTier(ToolbeltTools)
		action = strings.TrimPrefix(action, "toolbelt.")
	}

	execCtx, cancel := context.WithTimeout(ctx, l.tools.Config.ExecutionTimeout)
	defer cancel()
	
	l.PushState()
	
	observation, execErr := l.tools.ExecuteAction(execCtx, action, l.agent.Env)
	if execErr != nil {
		l.PopState()
		
		l.AddEvent("execute_action_failed", map[string]interface{}{
			"step": stepNum,
			"action": action,
			"error": execErr.Error(),
		})
		l.mutex.Lock()
		l.state.LastError = fmt.Sprintf("executing action '%s' failed: %v", action, execErr)
		l.mutex.Unlock()
		return "error_handling", execErr
	}
	
	l.mutex.Lock()
	l.state.Context["last_observation"] = observation
	l.mutex.Unlock()
	
	l.stateManager[ExecutionState][fmt.Sprintf("observation_%d", stepNum)] = observation
	
	l.PopState()
	
	l.AddEvent("execute_action_completed", map[string]interface{}{
		"step": stepNum,
		"action": action,
		"observation_length": len(observation),
	})

	l.AddEvent("handle_observation_started", map[string]interface{}{
		"step": stepNum,
	})

	l.agent.AddStepToHistory(thought, action, modelOutput, observation)

	taskComplete := false
	
	if strings.Contains(strings.ToLower(observation), "task complete") ||
	   strings.Contains(strings.ToLower(thought), "task is complete") {
		taskComplete = true
	}
	
	if stepNum >= 20 {
		return "completion_verification", nil
	}

	l.AddEvent("step_completed", map[string]interface{}{
		"step": stepNum,
		"task_done": taskComplete,
	})
	
	if taskComplete {
		return "completion_verification", nil
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
	case "execute_step":
		return l.executeStep(ctx)
		
	case "tool_transition_evaluation":
		return l.executeToolTransitionEvaluation(ctx)
	case "toolbelt_activation":
		return l.executeToolbeltActivation(ctx)
		
	case "completion_verification":
		return l.executeCompletionVerification(ctx)
	case "state_cleanup":
		return l.executeStateCleanup(ctx)
		
	case "error_handling":
		l.mutex.RLock()
		lastErr := l.state.LastError
		l.mutex.RUnlock()
		return l.handleError(ctx, fmt.Errorf(lastErr))
	case "fallback_processing":
		return l.executeFallbackProcessing(ctx)
		
	case "idle":
		return "idle", nil
		
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
}

func (l *Loop) executeTaskInitiation(ctx context.Context) (string, error) {
	l.AddEvent("task_initiation_started", map[string]interface{}{
		"task": l.state.Task,
		"agent": "samsepi0l",
	})
	
	l.stateManager[ExecutionState]["task_id"] = fmt.Sprintf("task_%d", time.Now().Unix())
	l.stateManager[ExecutionState]["start_time"] = time.Now().Unix()
	
	l.stateManager[ContextState]["agent_name"] = "samsepi0l"
	l.stateManager[ContextState]["agent_role"] = "Senior Software Engineering Lead & Technical Authority AI/ML"
	
	l.AddEvent("task_initiation_completed", nil)
	return "analyze_requirements", nil
}

func (l *Loop) executeAnalyzeRequirements(ctx context.Context) (string, error) {
	l.AddEvent("analyze_requirements_started", nil)
	
	task := l.state.Task
	
	isDataScienceTask := false
	dataScienceTriggers := []string{
		"data manipulation", "data science", "machine learning",
		"statistical analysis", "data visualization", "model training",
	}
	
	for _, trigger := range dataScienceTriggers {
		if strings.Contains(strings.ToLower(task), strings.ToLower(trigger)) {
			isDataScienceTask = true
			break
		}
	}
	
	l.mutex.Lock()
	l.state.Context["is_data_science_task"] = isDataScienceTask
	l.mutex.Unlock()
	
	l.stateManager[ExecutionState]["is_data_science_task"] = isDataScienceTask
	
	l.AddEvent("analyze_requirements_completed", map[string]interface{}{
		"is_data_science_task": isDataScienceTask,
	})
	
	return "tool_discovery", nil
}

func (l *Loop) executeToolDiscovery(ctx context.Context) (string, error) {
	l.AddEvent("tool_discovery_started", nil)
	
	l.mutex.RLock()
	isDataScienceTask, _ := l.state.Context["is_data_science_task"].(bool)
	l.mutex.RUnlock()
	
	initialTier := CoreTools
	if isDataScienceTask {
		initialTier = ToolchainTools
		l.TransitionToolTier(ToolchainTools)
	}
	
	l.stateManager[ToolState]["initial_tier"] = string(initialTier)
	l.stateManager[ToolState]["available_tools"] = l.tools.ListTools()
	
	l.AddEvent("tool_discovery_completed", map[string]interface{}{
		"initial_tier": initialTier,
	})
	
	return "planning_phase", nil
}

func (l *Loop) executePlanningPhase(ctx context.Context) (string, error) {
	l.AddEvent("planning_phase_started", nil)
	
	planOutput, err := l.agent.QueryModel(ctx, l.agent.History)
	if err != nil {
		return "error_handling", fmt.Errorf("planning phase failed: %v", err)
	}
	
	l.mutex.Lock()
	l.state.Context["execution_plan"] = planOutput
	l.mutex.Unlock()
	
	l.stateManager[ExecutionState]["has_plan"] = true
	
	l.AddEvent("planning_phase_completed", map[string]interface{}{
		"plan_generated": true,
	})
	
	return "execution_phase", nil
}

func (l *Loop) executeExecutionPhase(ctx context.Context) (string, error) {
	l.AddEvent("execution_phase_started", nil)
	
	if l.state.StepCount == 0 {
		l.mutex.Lock()
		l.state.StepCount = 1
		l.mutex.Unlock()
	}
	
	return "execute_step", nil
}

func (l *Loop) executeToolTransitionEvaluation(ctx context.Context) (string, error) {
	l.AddEvent("tool_transition_evaluation_started", nil)
	
	l.evaluateToolTransition(ctx)
	
	l.AddEvent("tool_transition_evaluation_completed", map[string]interface{}{
		"current_tier": l.state.ToolTier,
	})
	
	if l.state.ToolTier == ToolbeltTools {
		return "toolbelt_activation", nil
	}
	
	return "execution_phase", nil
}

func (l *Loop) executeToolbeltActivation(ctx context.Context) (string, error) {
	l.AddEvent("toolbelt_activation_started", nil)
	
	mcpTools := []string{"perplexity", "filesystem", "memory", "sequentialthinking"}
	
	l.stateManager[ToolState]["mcp_tools"] = mcpTools
	
	l.AddEvent("toolbelt_activation_completed", map[string]interface{}{
		"available_mcp_tools": mcpTools,
	})
	
	return "execution_phase", nil
}

func (l *Loop) executeCompletionVerification(ctx context.Context) (string, error) {
	l.AddEvent("completion_verification_started", nil)
	
	allStepsCompleted := true // Placeholder logic
	
	if allStepsCompleted {
		l.AddEvent("completion_verification_completed", map[string]interface{}{
			"all_steps_completed": true,
		})
		return "state_cleanup", nil
	}
	
	l.AddEvent("completion_verification_completed", map[string]interface{}{
		"all_steps_completed": false,
	})
	return "execution_phase", nil
}

func (l *Loop) executeStateCleanup(ctx context.Context) (string, error) {
	l.AddEvent("state_cleanup_started", nil)
	
	l.nestedStates = make([]map[string]interface{}, 0)
	l.nestedDepth = 0
	
	l.stateManager[ExecutionState]["cleanup_time"] = time.Now().Unix()
	
	l.AddEvent("state_cleanup_completed", nil)
	return "idle", nil
}

func (l *Loop) executeFallbackProcessing(ctx context.Context) (string, error) {
	l.AddEvent("fallback_processing_started", nil)
	
	l.mutex.RLock()
	lastErr := l.state.LastError
	l.mutex.RUnlock()
	
	l.AddEvent("fallback_attempt", map[string]interface{}{
		"original_error": lastErr,
	})
	
	l.mutex.Lock()
	l.state.Context["using_fallback"] = true
	l.mutex.Unlock()
	
	l.AddEvent("fallback_processing_completed", nil)
	return "execution_phase", nil
}

func (l *Loop) handleError(ctx context.Context, err error) (string, error) {
	l.AddEvent("error_handling_started", map[string]interface{}{
		"error": err.Error(),
		"agent": "samsepi0l",
	})

	l.stateManager[ExecutionState]["last_error"] = err.Error()
	l.stateManager[ExecutionState]["error_time"] = time.Now().Unix()
	
	shouldAttemptFallback := true // Placeholder logic
	
	if shouldAttemptFallback {
		l.AddEvent("error_handling_fallback_attempt", map[string]interface{}{
			"error": err.Error(),
		})
		return "fallback_processing", nil
	}
	
	fmt.Printf("Error encountered: %v. Stopping loop.\n", err)
	l.AddEvent("error_handling_failed", map[string]interface{}{
		"error": err.Error(),
	})
	return "idle", err // Return the original error to signal failure
}

func (l *Loop) executeSetup(ctx context.Context) (string, error) {
	l.AddEvent("setup_started", map[string]interface{}{
		"agent": "samsepi0l",
	})
	
	for _, stateType := range []StateType{ExecutionState, ToolState, ContextState} {
		if l.stateManager[stateType] == nil {
			l.stateManager[stateType] = make(map[string]interface{})
		}
	}
	
	l.mutex.Lock()
	l.agent.History = make([]Message, 0)
	l.state.StepCount = 0
	l.state.LastError = ""
	l.state.CurrentAction = ""
	l.state.ToolTier = CoreTools
	l.state.Context = make(map[string]interface{})
	l.mutex.Unlock()
	
	l.nestedStates = make([]map[string]interface{}, 0)
	l.nestedDepth = 0
	
	l.stateManager[ContextState]["agent_name"] = "samsepi0l"
	l.stateManager[ContextState]["agent_role"] = "Senior Software Engineering Lead & Technical Authority AI/ML"
	
	l.AddEvent("setup_completed", nil)
	return "task_initiation", nil // Start with task initiation phase
}
