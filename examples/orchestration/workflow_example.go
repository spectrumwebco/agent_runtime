package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
	"github.com/spectrumwebco/agent_runtime/internal/orchestration"
)

func main() {
	graphManager := langgraph.NewGraphManager()

	contextMonitor := langgraph.NewContextMonitor()

	orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

	workflowConfig := &orchestration.WorkflowConfig{
		AutoTrigger:   true,
		MaxIterations: 5,
		Timeout:       time.Minute * 5,
		RetryCount:    3,
		RetryDelay:    time.Second * 2,
		TriggerConditions: []orchestration.TriggerCondition{
			{
				Type:      orchestration.ContextTrigger,
				Pattern:   "knowledge",
				Priority:  1,
				Threshold: 0.7,
			},
		},
	}

	workflow, err := orchestrator.CreateWorkflow("KnowledgeWorkflow", "Workflow for knowledge tool triggering", workflowConfig)
	if err != nil {
		fmt.Printf("Failed to create workflow: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created workflow: %s (ID: %s)\n", workflow.Name, workflow.ID)

	graph := workflow.Graph

	graph.AddNode("start", func(state map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Starting knowledge workflow...")
		state["status"] = "started"
		state["iterations"] = 0
		state["max_iterations"] = workflowConfig.MaxIterations
		return state, nil
	})

	graph.AddNode("detect_context", func(state map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Detecting context...")
		state["context_detected"] = true
		state["context_type"] = "knowledge"
		state["context_confidence"] = 0.85
		return state, nil
	})

	graph.AddNode("process_knowledge", func(state map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Processing knowledge...")
		state["knowledge_processed"] = true
		return state, nil
	})

	graph.AddNode("update_agent", func(state map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Updating agent with new knowledge...")
		state["agent_updated"] = true
		iterations := state["iterations"].(int)
		state["iterations"] = iterations + 1
		return state, nil
	})

	graph.AddNode("end", func(state map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Ending knowledge workflow...")
		state["status"] = "completed"
		return state, nil
	})

	graph.AddEdge("start", "detect_context")
	graph.AddEdge("detect_context", "process_knowledge", langgraph.WithCondition(func(state map[string]interface{}) bool {
		contextDetected, ok := state["context_detected"].(bool)
		if !ok {
			return false
		}
		contextType, ok := state["context_type"].(string)
		if !ok {
			return false
		}
		contextConfidence, ok := state["context_confidence"].(float64)
		if !ok {
			return false
		}
		return contextDetected && contextType == "knowledge" && contextConfidence >= 0.7
	}))
	graph.AddEdge("detect_context", "end", langgraph.WithCondition(func(state map[string]interface{}) bool {
		contextDetected, ok := state["context_detected"].(bool)
		if !ok {
			return true
		}
		contextType, ok := state["context_type"].(string)
		if !ok {
			return true
		}
		contextConfidence, ok := state["context_confidence"].(float64)
		if !ok {
			return true
		}
		return !contextDetected || contextType != "knowledge" || contextConfidence < 0.7
	}))
	graph.AddEdge("process_knowledge", "update_agent")
	graph.AddEdge("update_agent", "detect_context", langgraph.WithCondition(func(state map[string]interface{}) bool {
		iterations, ok := state["iterations"].(int)
		if !ok {
			return false
		}
		maxIterations, ok := state["max_iterations"].(int)
		if !ok {
			return false
		}
		return iterations < maxIterations
	}))
	graph.AddEdge("update_agent", "end", langgraph.WithCondition(func(state map[string]interface{}) bool {
		iterations, ok := state["iterations"].(int)
		if !ok {
			return true
		}
		maxIterations, ok := state["max_iterations"].(int)
		if !ok {
			return true
		}
		return iterations >= maxIterations
	}))

	graph.SetEntryPoint("start")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal %s, shutting down...\n", sig)
		cancel()
	}()

	err = orchestrator.StartMonitoring(ctx)
	if err != nil {
		fmt.Printf("Failed to start monitoring: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Started monitoring for context triggers")
	fmt.Println("Press Ctrl+C to stop")

	go func() {
		time.Sleep(time.Second * 2)

		fmt.Println("Updating context with knowledge trigger...")
		err := orchestrator.UpdateContext(map[string]interface{}{
			"type":       "knowledge",
			"content":    "New knowledge about Go programming",
			"confidence": 0.85,
		})
		if err != nil {
			fmt.Printf("Failed to update context: %s\n", err)
			return
		}

		time.Sleep(time.Second * 10)

		fmt.Println("Updating context with another knowledge trigger...")
		err = orchestrator.UpdateContext(map[string]interface{}{
			"type":       "knowledge",
			"content":    "Additional knowledge about LangGraph-Go",
			"confidence": 0.9,
		})
		if err != nil {
			fmt.Printf("Failed to update context: %s\n", err)
			return
		}

		time.Sleep(time.Second * 10)

		fmt.Println("Test completed, stopping...")
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Stopped monitoring")

	workflows := orchestrator.ListWorkflows()
	for _, wf := range workflows {
		fmt.Printf("Workflow: %s (ID: %s)\n", wf.Name, wf.ID)
		fmt.Printf("Status: %s\n", wf.Status)
		fmt.Printf("Created at: %s\n", wf.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Updated at: %s\n", wf.UpdatedAt.Format(time.RFC3339))
	}
}
