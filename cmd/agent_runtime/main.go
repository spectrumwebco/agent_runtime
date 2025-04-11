package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/internal/env"
	"github.com/spectrumwebco/agent_runtime/internal/ffi/python"
	"github.com/spectrumwebco/agent_runtime/internal/server"
)

var (
	configPath string
	outputDir  string
	verbose    bool
	serverPort int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "agent_runtime",
		Short: "Agent Runtime - Go implementation of SWE-Agent and SWE-ReX",
		Long: `Agent Runtime is a Go implementation of SWE-Agent and SWE-ReX frameworks.
It provides a high-performance agent runtime with Python FFI capabilities.`,
	}

	// Run command - executes an agent with a problem statement
	runCmd := &cobra.Command{
		Use:   "run [problem_statement]",
		Short: "Run the agent on a problem statement",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running agent with problem statement:", args[0])
			
			// Load configuration
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				log.Fatalf("Failed to load configuration: %v", err)
			}
			
			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				log.Fatalf("Failed to create output directory: %v", err)
			}
			
			// Initialize environment
			environment, err := env.NewSWEEnv()
			if err != nil {
				log.Fatalf("Failed to initialize environment: %v", err)
			}
			
			// Create problem statement
			problemStmt := &agent.ProblemStatement{
				ID:               fmt.Sprintf("task-%d", time.Now().Unix()),
				ProblemStatement: args[0],
			}
			
			// Initialize agent
			agentInstance, err := agent.NewDefaultAgent()
			if err != nil {
				log.Fatalf("Failed to create agent: %v", err)
			}
			
			// Initialize Python FFI if needed
			pythonInterp, err := python.NewInterpreter()
			if err != nil {
				log.Printf("Warning: Failed to initialize Python interpreter: %v", err)
			} else {
				defer pythonInterp.Close()
				log.Println("Python interpreter initialized successfully")
			}
			
			// Run agent
			result, err := agentInstance.Run(environment, problemStmt, outputDir)
			if err != nil {
				log.Fatalf("Agent execution failed: %v", err)
			}
			
			fmt.Printf("Agent execution completed. Result: %+v\n", result)
		},
	}
	
	// Server command - starts the HTTP server
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the agent runtime server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting Agent Runtime server on port %d...\n", serverPort)
			
			// Initialize Gin router
			router := gin.Default()
			
			// Setup routes
			server.SetupRoutes(router)
			
			// Start server
			if err := router.Run(fmt.Sprintf(":%d", serverPort)); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		},
	}
	
	// Add flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "configs/default.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&outputDir, "output", "output", "Directory for output files")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Server port")
	
	// Add commands to root
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serverCmd)
	
	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
