package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
)

func init() {
	stateCmd := &cobra.Command{
		Use:   "state",
		Short: "Start the state management service",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting State Management service...")
			
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				log.Fatalf("Failed to load configuration: %v", err)
			}
			
			stateManager, err := statemanager.NewStateManager(cfg)
			if err != nil {
				log.Fatalf("Failed to initialize state manager: %v", err)
			}
			defer stateManager.Close()
			
			stateManager.StartLifoTaskProcessor()
			
			fmt.Println("State Management service started successfully")
			
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			
			fmt.Println("Shutting down State Management service...")
		},
	}
	
	rootCmd.AddCommand(stateCmd)
}
