// Package commands provides CLI commands for the agent runtime
package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

// NewRunCommand creates a new command for running tasks with Kled
func NewRunCommand() *cobra.Command {
	var configPath string
	var task string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a task with Kled",
		Long:  `Run a task with Kled, the Senior Software Engineering Lead & Technical Authority for AI/ML.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			mcpManager, err := mcp.NewManager(cfg)
			if err != nil {
				return fmt.Errorf("failed to create MCP manager: %w", err)
			}

			if err := mcpManager.StartServers(); err != nil {
				return fmt.Errorf("failed to start MCP servers: %w", err)
			}
			defer mcpManager.StopServers()

			agent, err := agent.New(cfg, mcpManager)
			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}

			ctx := context.Background()
			result, err := agent.Execute(ctx, task)
			if err != nil {
				return fmt.Errorf("failed to execute task: %w", err)
			}

			fmt.Printf("Task executed successfully: %s\n", result.Message)
			if verbose {
				fmt.Printf("Result data: %v\n", result.Data)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&task, "task", "t", "", "Task to execute")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	cmd.MarkFlagRequired("task")

	return cmd
}

func loadConfig(configPath string) (*config.Config, error) {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		configPath = filepath.Join(homeDir, ".sam", "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = filepath.Join(homeDir, ".config", "sam", "config.yaml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				configPath = "config.yaml"
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return nil, fmt.Errorf("configuration file not found")
				}
			}
		}
	}

	return config.LoadConfig(configPath)
}
