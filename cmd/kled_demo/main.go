package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/kledframework"
	"github.com/spectrumwebco/agent_runtime/pkg/types"
)

func main() {
	log.Println("Starting Kled.io Framework Demo")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	config := &kledframework.Config{
		Name:      "KledDemo",
		Version:   "1.0.0",
		LogLevel:  "info",
		EventStreamConfig: map[string]interface{}{
			"type":    "memory",
			"enabled": true,
		},
		StateManagerConfig: map[string]interface{}{
			"type":    "memory",
			"enabled": true,
		},
		DjangoBridgeConfig: map[string]interface{}{
			"grpc_address": "localhost:50051",
			"http_address": ":8082",
			"timeout":      "30s",
		},
		ReactBridgeConfig: map[string]interface{}{
			"address":    ":8081",
			"path":       "/ws",
			"heartbeat":  "30s",
			"state_path": "/api/state",
		},
	}

	framework, err := kledframework.NewKledFramework(config)
	if err != nil {
		log.Fatalf("Failed to create framework: %v", err)
	}

	if err := framework.Start(ctx); err != nil {
		log.Fatalf("Failed to start framework: %v", err)
	}
	defer framework.Stop()

	managerConfig := kledframework.DefaultMultiAgentManagerConfig()

	manager, err := kledframework.NewMultiAgentManager(managerConfig)
	if err != nil {
		log.Fatalf("Failed to create multi-agent manager: %v", err)
	}

	if err := manager.Start(ctx); err != nil {
		log.Fatalf("Failed to start multi-agent manager: %v", err)
	}
	defer manager.Stop()

	frontendAgent, err := manager.CreateAgent(types.AgentTypeFrontend, "Frontend Agent", map[string]interface{}{
		"capabilities": []string{"ui_design", "user_interaction", "frontend_development"},
	})
	if err != nil {
		log.Fatalf("Failed to create frontend agent: %v", err)
	}
	log.Printf("Created frontend agent: %s", frontendAgent.ID())

	appBuilderAgent, err := manager.CreateAgent(types.AgentTypeAppBuilder, "App Builder Agent", map[string]interface{}{
		"capabilities": []string{"app_scaffolding", "component_integration", "deployment"},
	})
	if err != nil {
		log.Fatalf("Failed to create app builder agent: %v", err)
	}
	log.Printf("Created app builder agent: %s", appBuilderAgent.ID())

	codegenAgent, err := manager.CreateAgent(types.AgentTypeCodegen, "Codegen Agent", map[string]interface{}{
		"capabilities": []string{"code_generation", "refactoring", "optimization"},
	})
	if err != nil {
		log.Fatalf("Failed to create codegen agent: %v", err)
	}
	log.Printf("Created codegen agent: %s", codegenAgent.ID())

	engineeringAgent, err := manager.CreateAgent(types.AgentTypeEngineering, "Engineering Agent", map[string]interface{}{
		"capabilities": []string{"architecture_design", "system_integration", "performance_tuning"},
	})
	if err != nil {
		log.Fatalf("Failed to create engineering agent: %v", err)
	}
	log.Printf("Created engineering agent: %s", engineeringAgent.ID())

	manager.AddEventListener(func(event *types.AgentEvent) {
		log.Printf("Event: %s, Agent: %s, Type: %s", event.ID, event.AgentID, event.Type)
	})

	task := &types.AgentTask{
		ID:          "task-1",
		AgentID:     frontendAgent.ID(),
		Type:        "ui_design",
		Description: "Design a user interface for the Kled.io Framework demo",
		Status:      "pending",
		Priority:    1,
		CreatedAt:   time.Now(),
		Parameters: map[string]interface{}{
			"theme":       "dark",
			"components":  []string{"header", "sidebar", "main", "footer"},
			"responsive":  true,
			"target_user": "developer",
		},
	}

	go func() {
		taskCtx, taskCancel := context.WithTimeout(ctx, 5*time.Minute)
		defer taskCancel()

		log.Printf("Executing task: %s", task.ID)
		result, err := manager.ExecuteTask(taskCtx, task)
		if err != nil {
			log.Printf("Failed to execute task: %v", err)
			return
		}

		log.Printf("Task completed: %s, Result: %v", task.ID, result)
	}()

	<-sigCh
	log.Println("Received signal, shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := manager.Stop(); err != nil {
		log.Printf("Failed to stop multi-agent manager: %v", err)
	}

	if err := framework.Stop(); err != nil {
		log.Printf("Failed to stop framework: %v", err)
	}

	log.Println("Kled.io Framework Demo stopped")
}
