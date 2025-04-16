package langgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/langchain"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type AgentIntegration struct {
	graph       *Graph
	stateManager *StateManager
	config      *config.Config
	httpClient  *http.Client
	djangoBaseURL string
}

func NewAgentIntegration(cfg *config.Config, djangoBaseURL string) *AgentIntegration {
	return &AgentIntegration{
		config:      cfg,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		djangoBaseURL: djangoBaseURL,
	}
}

func (ai *AgentIntegration) BuildAgentGraph() (*Graph, error) {
	const (
		StartNode      NodeID = "start"
		PlanNode       NodeID = "plan"
		ToolSelectionNode NodeID = "tool_selection"
		ToolExecutionNode NodeID = "tool_execution"
		ObservationNode NodeID = "observation"
		ReflectionNode  NodeID = "reflection"
		ActionNode      NodeID = "action"
		FinishNode      NodeID = "finish"
	)

	graph := NewGraph(StartNode)

	graph.DefineStateSchema(map[string]string{
		"input":        "string",
		"agent_state":  "map",
		"tools":        "array",
		"observations": "array",
		"plan":         "string",
		"selected_tool": "string",
		"tool_input":   "string",
		"tool_output":  "string",
		"reflection":   "string",
		"action":       "string",
		"output":       "string",
	})

	err := graph.AddNode(StartNode, ai.startNodeFn(), map[string]interface{}{
		"description": "Initialize the agent workflow",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(PlanNode, ai.planNodeFn(), map[string]interface{}{
		"description": "Create a plan based on the input",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(ToolSelectionNode, ai.toolSelectionNodeFn(), map[string]interface{}{
		"description": "Select a tool to execute",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(ToolExecutionNode, ai.toolExecutionNodeFn(), map[string]interface{}{
		"description": "Execute the selected tool",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(ObservationNode, ai.observationNodeFn(), map[string]interface{}{
		"description": "Process observations from tool execution",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(ReflectionNode, ai.reflectionNodeFn(), map[string]interface{}{
		"description": "Reflect on the observations and update the plan",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(ActionNode, ai.actionNodeFn(), map[string]interface{}{
		"description": "Determine the next action to take",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddNode(FinishNode, ai.finishNodeFn(), map[string]interface{}{
		"description": "Finalize the agent workflow",
	})
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(StartNode, PlanNode, DefaultEdge(PlanNode))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(PlanNode, ToolSelectionNode, DefaultEdge(ToolSelectionNode))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(ToolSelectionNode, ToolExecutionNode, DefaultEdge(ToolExecutionNode))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(ToolExecutionNode, ObservationNode, DefaultEdge(ObservationNode))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(ObservationNode, ReflectionNode, DefaultEdge(ReflectionNode))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(ReflectionNode, ActionNode, DefaultEdge(ActionNode))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(ActionNode, ToolSelectionNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			action, ok := state["action"].(string)
			return ok && action == "continue"
		},
		ToolSelectionNode,
	))
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge(ActionNode, FinishNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			action, ok := state["action"].(string)
			return ok && action == "finish"
		},
		FinishNode,
	))
	if err != nil {
		return nil, err
	}

	err = graph.SetExitPoint(FinishNode)
	if err != nil {
		return nil, err
	}

	err = graph.Compile()
	if err != nil {
		return nil, err
	}

	ai.graph = graph
	return graph, nil
}

func (ai *AgentIntegration) InitializeState(input string, tools []map[string]interface{}) (*StateManager, error) {
	initialState := map[string]interface{}{
		"input":        input,
		"agent_state":  map[string]interface{}{},
		"tools":        tools,
		"observations": []interface{}{},
		"plan":         "",
		"selected_tool": "",
		"tool_input":   "",
		"tool_output":  "",
		"reflection":   "",
		"action":       "",
		"output":       "",
	}

	schema := map[string]string{
		"input":        "string",
		"agent_state":  "map",
		"tools":        "array",
		"observations": "array",
		"plan":         "string",
		"selected_tool": "string",
		"tool_input":   "string",
		"tool_output":  "string",
		"reflection":   "string",
		"action":       "string",
		"output":       "string",
	}

	stateManager := NewStateManager(initialState, schema)
	ai.stateManager = stateManager
	return stateManager, nil
}

func (ai *AgentIntegration) ExecuteAgentWorkflow(ctx context.Context, input string, tools []map[string]interface{}) (map[string]interface{}, error) {
	if ai.graph == nil {
		_, err := ai.BuildAgentGraph()
		if err != nil {
			return nil, err
		}
	}

	stateManager, err := ai.InitializeState(input, tools)
	if err != nil {
		return nil, err
	}

	initialState := stateManager.GetState()
	result, err := ai.graph.Execute(ctx, initialState)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ai *AgentIntegration) callDjangoAgent(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/agent/%s", ai.djangoBaseURL, endpoint)
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ai.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Django agent API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}


func (ai *AgentIntegration) startNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		agentState := map[string]interface{}{
			"initialized": true,
			"timestamp":   time.Now().Unix(),
		}
		
		state["agent_state"] = agentState
		return state, nil
	}
}

func (ai *AgentIntegration) planNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		input, _ := state["input"].(string)
		
		payload := map[string]interface{}{
			"input": input,
			"state": state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("plan", payload)
		if err != nil {
			return nil, err
		}
		
		plan, _ := result["plan"].(string)
		state["plan"] = plan
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["plan"] = plan
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (ai *AgentIntegration) toolSelectionNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		input, _ := state["input"].(string)
		plan, _ := state["plan"].(string)
		tools, _ := state["tools"].([]interface{})
		
		payload := map[string]interface{}{
			"input": input,
			"plan":  plan,
			"tools": tools,
			"state": state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("select_tool", payload)
		if err != nil {
			return nil, err
		}
		
		selectedTool, _ := result["selected_tool"].(string)
		toolInput, _ := result["tool_input"].(string)
		
		state["selected_tool"] = selectedTool
		state["tool_input"] = toolInput
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["selected_tool"] = selectedTool
		agentState["tool_input"] = toolInput
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (ai *AgentIntegration) toolExecutionNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		selectedTool, _ := state["selected_tool"].(string)
		toolInput, _ := state["tool_input"].(string)
		
		payload := map[string]interface{}{
			"tool":  selectedTool,
			"input": toolInput,
			"state": state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("execute_tool", payload)
		if err != nil {
			return nil, err
		}
		
		toolOutput, _ := result["output"].(string)
		state["tool_output"] = toolOutput
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["tool_output"] = toolOutput
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (ai *AgentIntegration) observationNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		selectedTool, _ := state["selected_tool"].(string)
		toolInput, _ := state["tool_input"].(string)
		toolOutput, _ := state["tool_output"].(string)
		
		observation := map[string]interface{}{
			"tool":   selectedTool,
			"input":  toolInput,
			"output": toolOutput,
			"timestamp": time.Now().Unix(),
		}
		
		observations, _ := state["observations"].([]interface{})
		observations = append(observations, observation)
		state["observations"] = observations
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["observations"] = observations
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (ai *AgentIntegration) reflectionNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		observations, _ := state["observations"].([]interface{})
		plan, _ := state["plan"].(string)
		
		payload := map[string]interface{}{
			"observations": observations,
			"plan":         plan,
			"state":        state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("reflect", payload)
		if err != nil {
			return nil, err
		}
		
		reflection, _ := result["reflection"].(string)
		state["reflection"] = reflection
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["reflection"] = reflection
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (ai *AgentIntegration) actionNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		reflection, _ := state["reflection"].(string)
		observations, _ := state["observations"].([]interface{})
		
		payload := map[string]interface{}{
			"reflection":   reflection,
			"observations": observations,
			"state":        state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("decide_action", payload)
		if err != nil {
			return nil, err
		}
		
		action, _ := result["action"].(string)
		state["action"] = action
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["action"] = action
		state["agent_state"] = agentState
		
		return state, nil
	}
}

func (ai *AgentIntegration) finishNodeFn() NodeFn {
	return func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		observations, _ := state["observations"].([]interface{})
		
		payload := map[string]interface{}{
			"observations": observations,
			"state":        state["agent_state"],
		}
		
		result, err := ai.callDjangoAgent("finalize", payload)
		if err != nil {
			return nil, err
		}
		
		output, _ := result["output"].(string)
		state["output"] = output
		
		agentState, _ := state["agent_state"].(map[string]interface{})
		agentState["output"] = output
		agentState["completed"] = true
		agentState["completion_timestamp"] = time.Now().Unix()
		state["agent_state"] = agentState
		
		return state, nil
	}
}
