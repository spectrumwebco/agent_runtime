package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func NewLangGraphCommand() *cobra.Command {
	return langgraphCmd
}

var contextTriggerCmd = &cobra.Command{
	Use:   "context-trigger",
	Short: "Manage context triggers for automatic tool activation",
	Long: `Manage context triggers for automatic tool activation.
This command provides functionality to create, list, and manage triggers that activate tools when specific context patterns are detected.`,
}

var createContextTriggerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a context trigger for automatic tool activation",
	Long: `Create a context trigger for automatic tool activation.
This command creates a trigger that activates a tool when specific context patterns are detected.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		patterns, _ := cmd.Flags().GetStringSlice("patterns")
		toolName, _ := cmd.Flags().GetString("tool")
		priority, _ := cmd.Flags().GetInt("priority")
		configFile, _ := cmd.Flags().GetString("config")
		
		if configFile != "" {
			configData, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			
			var triggerConfig struct {
				ID          string   `json:"id"`
				Name        string   `json:"name"`
				Description string   `json:"description"`
				Patterns    []string `json:"patterns"`
				ToolName    string   `json:"tool_name"`
				Priority    int      `json:"priority"`
			}
			
			err = json.Unmarshal(configData, &triggerConfig)
			if err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}
			
			id = triggerConfig.ID
			name = triggerConfig.Name
			description = triggerConfig.Description
			patterns = triggerConfig.Patterns
			toolName = triggerConfig.ToolName
			priority = triggerConfig.Priority
		}
		
		if id == "" {
			return fmt.Errorf("trigger ID is required")
		}
		
		if name == "" {
			name = id
		}
		
		if len(patterns) == 0 {
			return fmt.Errorf("at least one pattern is required")
		}
		
		if toolName == "" {
			return fmt.Errorf("tool name is required")
		}
		
		trigger, err := langgraph.NewContextTrigger(id, name, description, patterns, toolName, priority)
		if err != nil {
			return fmt.Errorf("failed to create trigger: %w", err)
		}
		
		triggerConfig := struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Patterns    []string `json:"patterns"`
			ToolName    string   `json:"tool_name"`
			Priority    int      `json:"priority"`
		}{
			ID:          trigger.ID,
			Name:        trigger.Name,
			Description: trigger.Description,
			Patterns:    trigger.ContextPatterns,
			ToolName:    trigger.ToolName,
			Priority:    trigger.Priority,
		}
		
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = fmt.Sprintf("trigger_%s.json", id)
		}
		
		triggerJSON, err := json.MarshalIndent(triggerConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal trigger config: %w", err)
		}
		
		err = os.WriteFile(outputFile, triggerJSON, 0644)
		if err != nil {
			return fmt.Errorf("failed to write trigger config: %w", err)
		}
		
		fmt.Printf("Created trigger '%s' for tool '%s'\n", name, toolName)
		fmt.Printf("Trigger config saved to %s\n", outputFile)
		
		return nil
	},
}

var startContextMonitorCmd = &cobra.Command{
	Use:   "start-monitor",
	Short: "Start the context monitor for automatic tool triggering",
	Long: `Start the context monitor for automatic tool triggering.
This command starts a monitor that continuously checks for context changes and triggers tools based on configured triggers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		djangoBaseURL, _ := cmd.Flags().GetString("django-url")
		triggersDir, _ := cmd.Flags().GetString("triggers-dir")
		
		if djangoBaseURL == "" {
			djangoBaseURL = "http://localhost:8000"
		}
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		integration := langgraph.NewAgentIntegration(cfg, djangoBaseURL)
		triggerManager := langgraph.NewContextTriggerManager(integration)
		
		if triggersDir != "" {
			entries, err := os.ReadDir(triggersDir)
			if err != nil {
				return fmt.Errorf("failed to read triggers directory: %w", err)
			}
			
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
					triggerPath := fmt.Sprintf("%s/%s", triggersDir, entry.Name())
					triggerData, err := os.ReadFile(triggerPath)
					if err != nil {
						fmt.Printf("Warning: Failed to read trigger file %s: %v\n", triggerPath, err)
						continue
					}
					
					var triggerConfig struct {
						ID          string   `json:"id"`
						Name        string   `json:"name"`
						Description string   `json:"description"`
						Patterns    []string `json:"patterns"`
						ToolName    string   `json:"tool_name"`
						Priority    int      `json:"priority"`
					}
					
					err = json.Unmarshal(triggerData, &triggerConfig)
					if err != nil {
						fmt.Printf("Warning: Failed to parse trigger file %s: %v\n", triggerPath, err)
						continue
					}
					
					trigger, err := langgraph.NewContextTrigger(
						triggerConfig.ID,
						triggerConfig.Name,
						triggerConfig.Description,
						triggerConfig.Patterns,
						triggerConfig.ToolName,
						triggerConfig.Priority,
					)
					if err != nil {
						fmt.Printf("Warning: Failed to create trigger from file %s: %v\n", triggerPath, err)
						continue
					}
					
					triggerManager.AddTrigger(trigger)
					fmt.Printf("Loaded trigger '%s' for tool '%s'\n", trigger.Name, trigger.ToolName)
				}
			}
		}
		
		if len(triggerManager.GetTriggers()) == 0 {
			err = triggerManager.InitializeDefaultTriggers()
			if err != nil {
				return fmt.Errorf("failed to initialize default triggers: %w", err)
			}
			fmt.Println("Initialized default triggers")
		}
		
		monitor := langgraph.NewContextMonitor(triggerManager, 100)
		
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		
		monitor.Start(ctx)
		fmt.Println("Context monitor started")
		fmt.Println("Press Ctrl+C to stop")
		
		initialContext := map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"status":    "initialized",
		}
		monitor.AddContext(initialContext)
		
		<-ctx.Done()
		
		return nil
	},
}

var listContextTriggersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all context triggers",
	Long: `List all context triggers.
This command lists all context triggers that are configured for automatic tool activation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		triggersDir, _ := cmd.Flags().GetString("triggers-dir")
		
		if triggersDir == "" {
			return fmt.Errorf("triggers directory is required")
		}
		
		entries, err := os.ReadDir(triggersDir)
		if err != nil {
			return fmt.Errorf("failed to read triggers directory: %w", err)
		}
		
		fmt.Println("Context Triggers:")
		fmt.Println("----------------")
		
		count := 0
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				triggerPath := fmt.Sprintf("%s/%s", triggersDir, entry.Name())
				triggerData, err := os.ReadFile(triggerPath)
				if err != nil {
					fmt.Printf("Warning: Failed to read trigger file %s: %v\n", triggerPath, err)
					continue
				}
				
				var triggerConfig struct {
					ID          string   `json:"id"`
					Name        string   `json:"name"`
					Description string   `json:"description"`
					Patterns    []string `json:"patterns"`
					ToolName    string   `json:"tool_name"`
					Priority    int      `json:"priority"`
				}
				
				err = json.Unmarshal(triggerData, &triggerConfig)
				if err != nil {
					fmt.Printf("Warning: Failed to parse trigger file %s: %v\n", triggerPath, err)
					continue
				}
				
				fmt.Printf("ID: %s\n", triggerConfig.ID)
				fmt.Printf("Name: %s\n", triggerConfig.Name)
				fmt.Printf("Description: %s\n", triggerConfig.Description)
				fmt.Printf("Tool: %s\n", triggerConfig.ToolName)
				fmt.Printf("Priority: %d\n", triggerConfig.Priority)
				fmt.Printf("Patterns: %s\n", strings.Join(triggerConfig.Patterns, ", "))
				fmt.Println("----------------")
				
				count++
			}
		}
		
		if count == 0 {
			fmt.Println("No triggers found")
		} else {
			fmt.Printf("Found %d triggers\n", count)
		}
		
		return nil
	},
}

func init() {
	createContextTriggerCmd.Flags().StringP("id", "i", "", "Unique identifier for the trigger")
	createContextTriggerCmd.Flags().StringP("name", "n", "", "Name of the trigger")
	createContextTriggerCmd.Flags().StringP("description", "d", "", "Description of the trigger")
	createContextTriggerCmd.Flags().StringSliceP("patterns", "p", []string{}, "Context patterns that activate the trigger")
	createContextTriggerCmd.Flags().StringP("tool", "t", "", "Name of the tool to trigger")
	createContextTriggerCmd.Flags().IntP("priority", "r", 0, "Priority of the trigger (higher values have higher priority)")
	createContextTriggerCmd.Flags().StringP("config", "c", "", "Path to a JSON file containing trigger configuration")
	createContextTriggerCmd.Flags().StringP("output", "o", "", "Path to save the trigger configuration")
	
	startContextMonitorCmd.Flags().StringP("django-url", "u", "http://localhost:8000", "URL of the Django agent API")
	startContextMonitorCmd.Flags().StringP("triggers-dir", "d", "", "Directory containing trigger configuration files")
	
	listContextTriggersCmd.Flags().StringP("triggers-dir", "d", "", "Directory containing trigger configuration files")
	
	contextTriggerCmd.AddCommand(createContextTriggerCmd)
	contextTriggerCmd.AddCommand(startContextMonitorCmd)
	contextTriggerCmd.AddCommand(listContextTriggersCmd)
	
	langgraphCmd.AddCommand(contextTriggerCmd)
}
