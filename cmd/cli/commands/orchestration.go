package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
	"github.com/spectrumwebco/agent_runtime/internal/orchestration"
)

func NewOrchestrationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestration",
		Short: "Agent workflow orchestration",
		Long:  `Manage agent workflow orchestration using LangGraph-Go.`,
	}

	cmd.AddCommand(newOrchestrationCreateCommand())
	cmd.AddCommand(newOrchestrationListCommand())
	cmd.AddCommand(newOrchestrationStartCommand())
	cmd.AddCommand(newOrchestrationStopCommand())
	cmd.AddCommand(newOrchestrationMonitorCommand())
	cmd.AddCommand(newOrchestrationExtendCommand())

	return cmd
}

func newOrchestrationCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workflow",
		Long:  `Create a new agent workflow with the specified name and options.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("workflow name is required")
			}

			name := args[0]
			description, _ := cmd.Flags().GetString("description")
			autoTrigger, _ := cmd.Flags().GetBool("auto-trigger")
			maxIterations, _ := cmd.Flags().GetInt("max-iterations")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			retryCount, _ := cmd.Flags().GetInt("retry-count")
			retryDelay, _ := cmd.Flags().GetDuration("retry-delay")
			triggerPattern, _ := cmd.Flags().GetString("trigger-pattern")
			triggerPriority, _ := cmd.Flags().GetInt("trigger-priority")
			triggerThreshold, _ := cmd.Flags().GetFloat64("trigger-threshold")

			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			graphManager := langgraph.NewGraphManager()

			contextMonitor := langgraph.NewContextMonitor()

			orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

			workflowConfig := &orchestration.WorkflowConfig{
				AutoTrigger:   autoTrigger,
				MaxIterations: maxIterations,
				Timeout:       timeout,
				RetryCount:    retryCount,
				RetryDelay:    retryDelay,
			}

			if autoTrigger && triggerPattern != "" {
				workflowConfig.TriggerConditions = []orchestration.TriggerCondition{
					{
						Type:      orchestration.ContextTrigger,
						Pattern:   triggerPattern,
						Priority:  triggerPriority,
						Threshold: triggerThreshold,
					},
				}
			}

			workflow, err := orchestrator.CreateWorkflow(name, description, workflowConfig)
			if err != nil {
				return fmt.Errorf("failed to create workflow: %w", err)
			}

			fmt.Printf("Workflow %s created successfully\n", workflow.ID)
			fmt.Printf("Name: %s\n", workflow.Name)
			fmt.Printf("Description: %s\n", workflow.Description)
			fmt.Printf("Status: %s\n", workflow.Status)
			fmt.Printf("Created at: %s\n", workflow.CreatedAt.Format(time.RFC3339))

			return nil
		},
	}

	cmd.Flags().String("description", "", "Workflow description")
	cmd.Flags().Bool("auto-trigger", false, "Whether to auto-trigger the workflow")
	cmd.Flags().Int("max-iterations", 10, "Maximum number of iterations")
	cmd.Flags().Duration("timeout", time.Minute*10, "Workflow timeout")
	cmd.Flags().Int("retry-count", 3, "Number of retries")
	cmd.Flags().Duration("retry-delay", time.Second*5, "Delay between retries")
	cmd.Flags().String("trigger-pattern", ".*", "Trigger pattern")
	cmd.Flags().Int("trigger-priority", 1, "Trigger priority")
	cmd.Flags().Float64("trigger-threshold", 0.7, "Trigger threshold")

	return cmd
}

func newOrchestrationListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workflows",
		Long:  `List all registered agent workflows.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := cmd.Flags().GetString("status")

			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			graphManager := langgraph.NewGraphManager()

			contextMonitor := langgraph.NewContextMonitor()

			orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

			var workflows []*orchestration.Workflow
			if status != "" {
				workflows = orchestrator.GetWorkflowsByStatus(orchestration.WorkflowStatus(status))
			} else {
				workflows = orchestrator.ListWorkflows()
			}

			if len(workflows) == 0 {
				fmt.Println("No workflows found")
				return nil
			}

			fmt.Println("Registered workflows:")
			for _, workflow := range workflows {
				fmt.Printf("- %s (ID: %s, Status: %s)\n", workflow.Name, workflow.ID, workflow.Status)
			}

			return nil
		},
	}

	cmd.Flags().String("status", "", "Filter by status (idle, running, paused, completed, failed)")

	return cmd
}

func newOrchestrationStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a workflow",
		Long:  `Start an agent workflow with the specified ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("workflow ID is required")
			}

			id := args[0]
			contextFile, _ := cmd.Flags().GetString("context-file")

			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			graphManager := langgraph.NewGraphManager()

			contextMonitor := langgraph.NewContextMonitor()

			orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

			workflow, err := orchestrator.GetWorkflow(id)
			if err != nil {
				return fmt.Errorf("failed to get workflow: %w", err)
			}

			var initialState map[string]interface{}
			if contextFile != "" {
				data, err := os.ReadFile(contextFile)
				if err != nil {
					return fmt.Errorf("failed to read context file: %w", err)
				}

				err = json.Unmarshal(data, &initialState)
				if err != nil {
					return fmt.Errorf("failed to parse context file: %w", err)
				}
			} else {
				initialState = make(map[string]interface{})
			}

			err = orchestrator.StartWorkflow(id, initialState)
			if err != nil {
				return fmt.Errorf("failed to start workflow: %w", err)
			}

			fmt.Printf("Workflow %s started successfully\n", workflow.Name)

			return nil
		},
	}

	cmd.Flags().String("context-file", "", "Path to context file")

	return cmd
}

func newOrchestrationStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a workflow",
		Long:  `Stop an agent workflow with the specified ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("workflow ID is required")
			}

			id := args[0]

			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			graphManager := langgraph.NewGraphManager()

			contextMonitor := langgraph.NewContextMonitor()

			orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

			workflow, err := orchestrator.GetWorkflow(id)
			if err != nil {
				return fmt.Errorf("failed to get workflow: %w", err)
			}

			err = orchestrator.StopWorkflow(id)
			if err != nil {
				return fmt.Errorf("failed to stop workflow: %w", err)
			}

			fmt.Printf("Workflow %s stopped successfully\n", workflow.Name)

			return nil
		},
	}

	return cmd
}

func newOrchestrationMonitorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor for context triggers",
		Long:  `Start monitoring for context triggers to automatically start workflows.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			graphManager := langgraph.NewGraphManager()

			contextMonitor := langgraph.NewContextMonitor()

			orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

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
				return fmt.Errorf("failed to start monitoring: %w", err)
			}

			fmt.Println("Started monitoring for context triggers")
			fmt.Println("Press Ctrl+C to stop")

			<-ctx.Done()
			fmt.Println("Stopped monitoring")

			return nil
		},
	}

	return cmd
}

func newOrchestrationExtendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extend",
		Short: "Extend agent loop",
		Long:  `Extend the agent loop with a workflow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			graphManager := langgraph.NewGraphManager()

			contextMonitor := langgraph.NewContextMonitor()

			orchestrator := orchestration.NewAgentOrchestrator(graphManager, contextMonitor)

			agentLoopExtension := langgraph.NewAgentLoopExtension()

			err = orchestrator.ExtendAgentLoop(agentLoopExtension)
			if err != nil {
				return fmt.Errorf("failed to extend agent loop: %w", err)
			}

			fmt.Println("Agent loop extended successfully")

			return nil
		},
	}

	return cmd
}
