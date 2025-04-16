package langgraph

import (
	"context"
	"fmt"
	"time"
)

func ExampleAgentWorkflow() {
	const (
		StartNode      NodeID = "start"
		ProcessNode    NodeID = "process"
		DecisionNode   NodeID = "decision"
		ToolNode       NodeID = "tool"
		FinishNode     NodeID = "finish"
	)

	graph := NewGraph(StartNode)

	graph.DefineStateSchema(map[string]string{
		"input":     "string",
		"processed": "string",
		"decision":  "string",
		"tool_output": "string",
		"output":    "string",
	})

	graph.AddNode(StartNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Starting agent workflow...")
		return state, nil
	}, map[string]interface{}{
		"description": "Initialize the agent workflow",
	})

	graph.AddNode(ProcessNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		input, _ := state["input"].(string)
		processed := fmt.Sprintf("Processed: %s", input)
		fmt.Println(processed)
		
		state["processed"] = processed
		return state, nil
	}, map[string]interface{}{
		"description": "Process the input",
	})

	graph.AddNode(DecisionNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		processed, _ := state["processed"].(string)
		
		decision := "use_tool"
		if len(processed) > 20 {
			decision = "finish"
		}
		
		fmt.Printf("Decision: %s\n", decision)
		state["decision"] = decision
		return state, nil
	}, map[string]interface{}{
		"description": "Make a decision based on the processed input",
	})

	graph.AddNode(ToolNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		processed, _ := state["processed"].(string)
		
		toolOutput := fmt.Sprintf("Tool output for: %s", processed)
		fmt.Println(toolOutput)
		
		state["tool_output"] = toolOutput
		return state, nil
	}, map[string]interface{}{
		"description": "Execute a tool",
	})

	graph.AddNode(FinishNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		processed, _ := state["processed"].(string)
		toolOutput, _ := state["tool_output"].(string)
		
		var output string
		if toolOutput != "" {
			output = fmt.Sprintf("Final output: %s with tool: %s", processed, toolOutput)
		} else {
			output = fmt.Sprintf("Final output: %s", processed)
		}
		
		fmt.Println(output)
		state["output"] = output
		return state, nil
	}, map[string]interface{}{
		"description": "Finalize the agent workflow",
	})

	graph.AddEdge(StartNode, ProcessNode, DefaultEdge(ProcessNode))
	graph.AddEdge(ProcessNode, DecisionNode, DefaultEdge(DecisionNode))
	
	graph.AddEdge(DecisionNode, ToolNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			decision, _ := state["decision"].(string)
			return decision == "use_tool"
		},
		ToolNode,
	))
	
	graph.AddEdge(DecisionNode, FinishNode, ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			decision, _ := state["decision"].(string)
			return decision == "finish"
		},
		FinishNode,
	))
	
	graph.AddEdge(ToolNode, FinishNode, DefaultEdge(FinishNode))

	graph.SetExitPoint(FinishNode)

	err := graph.Compile()
	if err != nil {
		fmt.Printf("Error compiling graph: %v\n", err)
		return
	}

	initialState := map[string]interface{}{
		"input": "Hello, LangGraph-Go!",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := graph.Execute(ctx, initialState)
	if err != nil {
		fmt.Printf("Error executing graph: %v\n", err)
		return
	}

	fmt.Printf("Result: %v\n", result)
}

func ExampleStateManagement() {
	schema := map[string]string{
		"input":    "string",
		"counter":  "int",
		"finished": "bool",
	}

	initialState := map[string]interface{}{
		"input":    "Hello, StateManager!",
		"counter":  0,
		"finished": false,
	}

	stateManager := NewStateManager(initialState, schema)

	currentState := stateManager.GetState()
	fmt.Printf("Initial state: %v\n", currentState)

	err := stateManager.UpdateState(map[string]interface{}{
		"counter": 1,
	})
	if err != nil {
		fmt.Printf("Error updating state: %v\n", err)
		return
	}

	updatedState := stateManager.GetState()
	fmt.Printf("Updated state: %v\n", updatedState)

	updater := ChainStateUpdaters(
		KeyValueStateUpdater("counter", 2),
		KeyValueStateUpdater("finished", true),
	)

	ctx := context.Background()
	finalState, err := updater(ctx, updatedState)
	if err != nil {
		fmt.Printf("Error applying state updaters: %v\n", err)
		return
	}

	err = stateManager.UpdateState(finalState)
	if err != nil {
		fmt.Printf("Error updating state: %v\n", err)
		return
	}

	finalStateFromManager := stateManager.GetState()
	fmt.Printf("Final state: %v\n", finalStateFromManager)

	history := stateManager.GetHistory()
	fmt.Printf("State history: %v\n", history)

	trajectory, err := stateManager.GetTrajectory()
	if err != nil {
		fmt.Printf("Error getting trajectory: %v\n", err)
		return
	}
	fmt.Printf("Trajectory: %s\n", trajectory)
}
