package langgraph

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type AgentLoopExtension struct {
	integration     *AgentIntegration
	contextMonitor  *ContextMonitor
	triggerManager  *ContextTriggerManager
	stateHistory    []*StateManager
	maxHistorySize  int
	continuationMap map[string]string
	mu              sync.RWMutex
}

func NewAgentLoopExtension(integration *AgentIntegration) *AgentLoopExtension {
	triggerManager := NewContextTriggerManager(integration)
	contextMonitor := NewContextMonitor(triggerManager, 100)
	
	return &AgentLoopExtension{
		integration:     integration,
		contextMonitor:  contextMonitor,
		triggerManager:  triggerManager,
		stateHistory:    make([]*StateManager, 0),
		maxHistorySize:  10,
		continuationMap: make(map[string]string),
	}
}

func (ale *AgentLoopExtension) Initialize() error {
	err := ale.triggerManager.InitializeDefaultTriggers()
	if err != nil {
		return fmt.Errorf("failed to initialize default triggers: %w", err)
	}
	
	ctx := context.Background()
	ale.contextMonitor.Start(ctx)
	
	return nil
}

func (ale *AgentLoopExtension) ExecuteExtendedWorkflow(ctx context.Context, input string, tools []map[string]interface{}, options map[string]interface{}) (map[string]interface{}, error) {
	result, err := ale.integration.ExecuteAgentWorkflow(ctx, input, tools)
	if err != nil {
		return nil, err
	}
	
	ale.addStateToHistory(ale.integration.stateManager)
	
	shouldContinue, _ := options["continue_workflow"].(bool)
	maxIterations, _ := options["max_iterations"].(int)
	if maxIterations <= 0 {
		maxIterations = 3
	}
	
	iteration := 1
	for shouldContinue && iteration < maxIterations {
		output, _ := result["output"].(string)
		
		continuationInput := fmt.Sprintf("Continuing from previous output: %s", output)
		
		continuationResult, err := ale.integration.ExecuteAgentWorkflow(ctx, continuationInput, tools)
		if err != nil {
			return nil, err
		}
		
		ale.addStateToHistory(ale.integration.stateManager)
		
		result = continuationResult
		
		shouldContinue = ale.shouldContinueWorkflow(result)
		iteration++
	}
	
	ale.contextMonitor.AddContext(result)
	
	return result, nil
}

func (ale *AgentLoopExtension) ExecuteContextTriggeredWorkflow(ctx context.Context, contextData map[string]interface{}, tools []map[string]interface{}) ([]map[string]interface{}, error) {
	triggeredResults, err := ale.triggerManager.ProcessContext(ctx, contextData)
	if err != nil {
		return nil, err
	}
	
	results := make([]map[string]interface{}, 0)
	
	for _, result := range triggeredResults {
		toolName, _ := result["tool"].(string)
		toolInput, _ := result["input"].(map[string]interface{})
		
		inputStr := fmt.Sprintf("Execute tool %s with input: %v", toolName, toolInput)
		
		workflowResult, err := ale.integration.ExecuteAgentWorkflow(ctx, inputStr, tools)
		if err != nil {
			continue
		}
		
		ale.addStateToHistory(ale.integration.stateManager)
		
		results = append(results, workflowResult)
	}
	
	return results, nil
}

func (ale *AgentLoopExtension) RegisterContinuation(workflowID string, continuationPoint string) {
	ale.mu.Lock()
	defer ale.mu.Unlock()
	
	ale.continuationMap[workflowID] = continuationPoint
}

func (ale *AgentLoopExtension) GetContinuation(workflowID string) string {
	ale.mu.RLock()
	defer ale.mu.RUnlock()
	
	return ale.continuationMap[workflowID]
}

func (ale *AgentLoopExtension) AddTrigger(trigger *ContextTrigger) {
	ale.triggerManager.AddTrigger(trigger)
}

func (ale *AgentLoopExtension) GetStateHistory() []*StateManager {
	ale.mu.RLock()
	defer ale.mu.RUnlock()
	
	history := make([]*StateManager, len(ale.stateHistory))
	copy(history, ale.stateHistory)
	
	return history
}

func (ale *AgentLoopExtension) GetTrajectory() (string, error) {
	ale.mu.RLock()
	defer ale.mu.RUnlock()
	
	if len(ale.stateHistory) == 0 {
		return "", fmt.Errorf("no state history available")
	}
	
	trajectory := "Agent Execution Trajectory:\n"
	
	for i, stateManager := range ale.stateHistory {
		stateTrajectory, err := stateManager.GetTrajectory()
		if err != nil {
			continue
		}
		
		trajectory += fmt.Sprintf("Iteration %d:\n%s\n", i+1, stateTrajectory)
	}
	
	return trajectory, nil
}

func (ale *AgentLoopExtension) addStateToHistory(stateManager *StateManager) {
	ale.mu.Lock()
	defer ale.mu.Unlock()
	
	stateCopy := stateManager.Clone()
	
	ale.stateHistory = append(ale.stateHistory, stateCopy)
	
	if len(ale.stateHistory) > ale.maxHistorySize {
		ale.stateHistory = ale.stateHistory[len(ale.stateHistory)-ale.maxHistorySize:]
	}
}

func (ale *AgentLoopExtension) shouldContinueWorkflow(result map[string]interface{}) bool {
	action, ok := result["action"].(string)
	if ok && action == "continue" {
		return true
	}
	
	continuationFlag, ok := result["should_continue"].(bool)
	if ok && continuationFlag {
		return true
	}
	
	output, ok := result["output"].(string)
	if ok && len(output) > 0 {
		continuationPhrases := []string{
			"more analysis needed",
			"further investigation required",
			"continue processing",
			"next steps",
			"to be continued",
		}
		
		for _, phrase := range continuationPhrases {
			if containsIgnoreCase(output, phrase) {
				return true
			}
		}
	}
	
	return false
}

func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}

type ExtendedAgentGraph struct {
	baseGraph      *Graph
	contextNode    NodeID
	knowledgeNode  NodeID
	feedbackNode   NodeID
	continuationNode NodeID
}

func NewExtendedAgentGraph(baseGraph *Graph) *ExtendedAgentGraph {
	return &ExtendedAgentGraph{
		baseGraph:      baseGraph,
		contextNode:    "context",
		knowledgeNode:  "knowledge",
		feedbackNode:   "feedback",
		continuationNode: "continuation",
	}
}

func (eag *ExtendedAgentGraph) BuildExtendedGraph(ai *AgentIntegration) (*Graph, error) {
	graph := eag.baseGraph.Clone()
	
	err := graph.AddNode(eag.contextNode, eag.contextNodeFn(ai), map[string]interface{}{
		"description": "Process context and trigger tools",
	})
	if err != nil {
		return nil, err
	}
	
	err = graph.AddNode(eag.knowledgeNode, eag.knowledgeNodeFn(ai), map[string]interface{}{
		"description": "Access knowledge and provide information",
	})
	if err != nil {
		return nil, err
	}
	
	err = graph.AddNode(eag.feedbackNode, eag.feedbackNodeFn(ai), map[string]interface{}{
		"description": "Process feedback and update agent behavior",
	})
	if err != nil {
		return nil, err
	}
	
	err = graph.AddNode(eag.continuationNode, eag.continuationNodeFn(ai), map[string]interface{}{
		"description": "Continue the agent workflow",
	})
	if err != nil {
		return nil, err
	}
	
	err = graph.AddEdge(eag.contextNode, eag.knowledgeNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			contextType, _ := state["context_type"].(string)
			return contextType == "knowledge"
		},
		eag.knowledgeNode,
	))
	if err != nil {
		return nil, err
	}
	
	actionNode := NodeID("action")
	err = graph.AddEdge(actionNode, eag.continuationNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			action, _ := state["action"].(string)
			return action == "continue"
		},
		eag.continuationNode,
	))
	if err != nil {
		return nil, err
	}
	
	toolSelectionNode := NodeID("tool_selection")
	err = graph.AddEdge(eag.continuationNode, toolSelectionNode, DefaultEdge(toolSelectionNode))
	if err != nil {
		return nil, err
	}
	
	finishNode := NodeID("finish")
	err = graph.AddEdge(finishNode, eag.feedbackNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			feedbackEnabled, _ := state["feedback_enabled"].(bool)
			return feedbackEnabled
		},
		eag.feedbackNode,
	))
	if err != nil {
		return nil, err
	}
	
	err = graph.SetExitPoint(eag.feedbackNode)
	if err != nil {
		return nil, err
	}
	
	err = graph.Compile()
	if err != nil {
		return nil, err
	}
	
	return graph, nil
}

func (eag *ExtendedAgentGraph) contextNodeFn(ai *AgentIntegration) NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		contextData, _ := state["context_data"].(map[string]interface{})
		
		payload := map[string]interface{}{
			"context": contextData,
			"state":   state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("process_context", payload)
		if err != nil {
			return state, nil
		}
		
		contextType, _ := result["context_type"].(string)
		state["context_type"] = contextType
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["context_type"] = contextType
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (eag *ExtendedAgentGraph) knowledgeNodeFn(ai *AgentIntegration) NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		contextData, _ := state["context_data"].(map[string]interface{})
		
		payload := map[string]interface{}{
			"context": contextData,
			"state":   state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("access_knowledge", payload)
		if err != nil {
			return state, nil
		}
		
		knowledge, _ := result["knowledge"].(string)
		state["knowledge"] = knowledge
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["knowledge"] = knowledge
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (eag *ExtendedAgentGraph) feedbackNodeFn(ai *AgentIntegration) NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		output, _ := state["output"].(string)
		
		payload := map[string]interface{}{
			"output": output,
			"state":  state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("process_feedback", payload)
		if err != nil {
			return state, nil
		}
		
		feedback, _ := result["feedback"].(string)
		state["feedback"] = feedback
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["feedback"] = feedback
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (eag *ExtendedAgentGraph) continuationNodeFn(ai *AgentIntegration) NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		agentState, _ := state["agent_state"].(map[string]interface{})
		
		continuationCount, _ := agentState["continuation_count"].(int)
		continuationCount++
		agentState["continuation_count"] = continuationCount
		
		agentState[fmt.Sprintf("continuation_%d_timestamp", continuationCount)] = time.Now().Unix()
		
		state["agent_state"] = agentState
		
		fmt.Printf("Agent workflow continuation #%d at %s\n", 
			continuationCount, 
			time.Now().Format(time.RFC3339))
		
		return state, nil
	}
}

type ExtendedStateManager struct {
	*StateManager
	continuationPoints []string
	feedbackHistory    []string
	knowledgeAccesses  []string
	mu                 sync.RWMutex
}

func NewExtendedStateManager(stateManager *StateManager) *ExtendedStateManager {
	return &ExtendedStateManager{
		StateManager:       stateManager,
		continuationPoints: make([]string, 0),
		feedbackHistory:    make([]string, 0),
		knowledgeAccesses:  make([]string, 0),
	}
}

func (esm *ExtendedStateManager) AddContinuationPoint(point string) {
	esm.mu.Lock()
	defer esm.mu.Unlock()
	
	esm.continuationPoints = append(esm.continuationPoints, point)
}

func (esm *ExtendedStateManager) AddFeedback(feedback string) {
	esm.mu.Lock()
	defer esm.mu.Unlock()
	
	esm.feedbackHistory = append(esm.feedbackHistory, feedback)
}

func (esm *ExtendedStateManager) AddKnowledgeAccess(knowledge string) {
	esm.mu.Lock()
	defer esm.mu.Unlock()
	
	esm.knowledgeAccesses = append(esm.knowledgeAccesses, knowledge)
}

func (esm *ExtendedStateManager) GetContinuationPoints() []string {
	esm.mu.RLock()
	defer esm.mu.RUnlock()
	
	points := make([]string, len(esm.continuationPoints))
	copy(points, esm.continuationPoints)
	
	return points
}

func (esm *ExtendedStateManager) GetFeedbackHistory() []string {
	esm.mu.RLock()
	defer esm.mu.RUnlock()
	
	history := make([]string, len(esm.feedbackHistory))
	copy(history, esm.feedbackHistory)
	
	return history
}

func (esm *ExtendedStateManager) GetKnowledgeAccesses() []string {
	esm.mu.RLock()
	defer esm.mu.RUnlock()
	
	accesses := make([]string, len(esm.knowledgeAccesses))
	copy(accesses, esm.knowledgeAccesses)
	
	return accesses
}

func (esm *ExtendedStateManager) GetExtendedTrajectory() (string, error) {
	baseTrajectory, err := esm.StateManager.GetTrajectory()
	if err != nil {
		return "", err
	}
	
	esm.mu.RLock()
	defer esm.mu.RUnlock()
	
	extendedTrajectory := baseTrajectory + "\n\nExtended Information:\n"
	
	if len(esm.continuationPoints) > 0 {
		extendedTrajectory += "Continuation Points:\n"
		for i, point := range esm.continuationPoints {
			extendedTrajectory += fmt.Sprintf("  %d. %s\n", i+1, point)
		}
		extendedTrajectory += "\n"
	}
	
	if len(esm.feedbackHistory) > 0 {
		extendedTrajectory += "Feedback History:\n"
		for i, feedback := range esm.feedbackHistory {
			extendedTrajectory += fmt.Sprintf("  %d. %s\n", i+1, feedback)
		}
		extendedTrajectory += "\n"
	}
	
	if len(esm.knowledgeAccesses) > 0 {
		extendedTrajectory += "Knowledge Accesses:\n"
		for i, knowledge := range esm.knowledgeAccesses {
			extendedTrajectory += fmt.Sprintf("  %d. %s\n", i+1, knowledge)
		}
	}
	
	return extendedTrajectory, nil
}
