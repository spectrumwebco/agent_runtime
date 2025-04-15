package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/server"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

var (
	grpcPort int
)

func init() {
	grpcServerCmd := &cobra.Command{
		Use:   "grpc-server",
		Short: "Start the gRPC server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting gRPC server on port %d...\n", grpcPort)
			
			// Load configuration
			cfg := &config.Config{}
			cfg.Server.GRPCPort = grpcPort
			cfg.Server.GRPCHost = "0.0.0.0"
			
			// Create agent instance
			agentInstance, err := agent.NewDefaultAgent()
			if err != nil {
				log.Fatalf("Failed to create agent: %v", err)
			}
			
			// Create server
			srv, err := server.New(cfg)
			if err != nil {
				log.Fatalf("Failed to create server: %v", err)
			}
			
			// Start server
			if err := srv.Start(); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
			
			fmt.Printf("gRPC server started on port %d\n", grpcPort)
			
			// Wait for interrupt signal
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
			<-sigCh
			
			// Stop server
			if err := srv.Stop(); err != nil {
				log.Fatalf("Failed to stop server: %v", err)
			}
		},
	}
	
	grpcServerCmd.Flags().IntVar(&grpcPort, "port", 50051, "gRPC server port")
	
	rootCmd.AddCommand(grpcServerCmd)
}
