package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
)

type SimpleEventStream struct {
	Events []*models.Event
}

func NewSimpleEventStream() *SimpleEventStream {
	return &SimpleEventStream{
		Events: []*models.Event{},
	}
}

func (s *SimpleEventStream) AddEvent(event *models.Event) error {
	s.Events = append(s.Events, event)
	fmt.Printf("Event added: %s - %s\n", event.Type, event.Source)
	return nil
}

func (s *SimpleEventStream) Subscribe(eventType models.EventType, callback func(*models.Event)) error {
	return nil
}

func (s *SimpleEventStream) Unsubscribe(eventType models.EventType, callback func(*models.Event)) error {
	return nil
}

func main() {
	fmt.Println("LangGraph and LangChain Integration Demo")
	fmt.Println("=======================================")

	eventStream := NewSimpleEventStream()

	system, err := langraph.CreateStandardMultiAgentSystem("demo-system", "Demo multi-agent system", eventStream)
	if err != nil {
		fmt.Printf("Error creating multi-agent system: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Multi-agent system created successfully")
	fmt.Println("Agents:")
	for _, agent := range system.ListAgents() {
		fmt.Printf("- %s (%s): %s\n", agent.Config.Name, agent.Config.Role, agent.Config.Description)
		fmt.Printf("  Capabilities: %v\n", agent.Config.Capabilities)
	}

	var orchestratorAgent *langraph.Agent
	for _, agent := range system.ListAgents() {
		if agent.Config.Role == langraph.AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}

	if orchestratorAgent == nil {
		fmt.Println("Orchestrator agent not found")
		os.Exit(1)
	}

	fmt.Println("\nExecuting system with orchestrator agent...")
	ctx := context.Background()
	execution, err := system.Execute(ctx, orchestratorAgent.Config.ID, map[string]interface{}{
		"task": "Create a new UI component",
	})

	if err != nil {
		fmt.Printf("Error executing system: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Execution started: %s\n", execution.ID)

	time.Sleep(500 * time.Millisecond)

	fmt.Println("\nEvents generated:")
	for i, event := range eventStream.Events {
		fmt.Printf("%d. Type: %s, Source: %s\n", i+1, event.Type, event.Source)
	}

	fmt.Println("\nLangGraph and LangChain integration verified successfully!")
}
