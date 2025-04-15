package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

// FrameworkConfig defines the configuration for the Kled.io Framework
type FrameworkConfig struct {
	Name        string
	Description string
	Debug       bool
}

// Framework represents the Kled.io Framework
type Framework struct {
	Config FrameworkConfig
}

// NewFramework creates a new instance of the Kled.io Framework
func NewFramework(config FrameworkConfig) (*Framework, error) {
	return &Framework{
		Config: config,
	}, nil
}

// Start starts the Kled.io Framework
func (f *Framework) Start(ctx context.Context) error {
	return nil
}

// Stop stops the Kled.io Framework
func (f *Framework) Stop() error {
	return nil
}

// NewStartCommand creates a new command for starting the Kled.io Framework
func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Kled.io Framework",
		Long:  `Start the Kled.io Framework with the specified configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")
			if configPath == "" {
				fmt.Println("Error: config path is required")
				os.Exit(1)
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			frameworkConfig := FrameworkConfig{
				Name:        cfg.Agent.Name,
				Description: cfg.Agent.Description,
				Debug:       cfg.Logging.Level == "debug",
			}

			framework, err := NewFramework(frameworkConfig)
			if err != nil {
				fmt.Printf("Error creating framework: %v\n", err)
				os.Exit(1)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := framework.Start(ctx); err != nil {
				fmt.Printf("Error starting framework: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Kled.io Framework started successfully")

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh

			fmt.Println("Stopping Kled.io Framework...")

			if err := framework.Stop(); err != nil {
				fmt.Printf("Error stopping framework: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Kled.io Framework stopped successfully")
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to configuration file")
	cmd.MarkFlagRequired("config")

	return cmd
}
