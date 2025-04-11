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
	"github.com/spectrumwebco/agent_runtime/internal/mcp"
	"github.com/spectrumwebco/agent_runtime/internal/parser"
	"github.com/spectrumwebco/agent_runtime/internal/server"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

var (
	configPath    string
	outputDir     string
	verbose       bool
	serverPort    int
	toolsDir      string
	mcpServerPort int
	pythonPath    string
	modelName     string
	parserType    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "agent_runtime",
		Short: "Agent Runtime - Go implementation of SWE-Agent and SWE-ReX",
		Long: `Agent Runtime is a Go implementation of SWE-Agent and SWE-ReX frameworks.
It provides a high-performance agent runtime with Python FFI capabilities.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Println("Verbose logging enabled")
			} else {
				log.SetFlags(log.LstdFlags)
			}

			if pythonPath != "" {
				currentPath := os.Getenv("PYTHONPATH")
				if currentPath != "" {
					os.Setenv("PYTHONPATH", fmt.Sprintf("%s:%s", pythonPath, currentPath))
				} else {
					os.Setenv("PYTHONPATH", pythonPath)
				}
				log.Println("PYTHONPATH set to:", os.Getenv("PYTHONPATH"))
			}
		},
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
			defer environment.Close()
			
			scriptRunner, err := initializePythonFFI()
			if err != nil {
				log.Printf("Warning: Failed to initialize Python FFI: %v", err)
			}
			
			toolRegistry, err := initializeToolRegistry(scriptRunner)
			if err != nil {
				log.Fatalf("Failed to initialize tool registry: %v", err)
			}
			
			// Create problem statement
			problemStmt := &agent.ProblemStatement{
				ID:               fmt.Sprintf("task-%d", time.Now().Unix()),
				ProblemStatement: args[0],
			}
			
			agentInstance, err := agent.NewDefaultAgent(
				agent.WithTools(toolRegistry),
				agent.WithParser(parser.ParserFactory(parserType)),
				agent.WithModelName(modelName),
			)
			if err != nil {
				log.Fatalf("Failed to create agent: %v", err)
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
			
			scriptRunner, err := initializePythonFFI()
			if err != nil {
				log.Printf("Warning: Failed to initialize Python FFI: %v", err)
			}
			
			toolRegistry, err := initializeToolRegistry(scriptRunner)
			if err != nil {
				log.Fatalf("Failed to initialize tool registry: %v", err)
			}
			
			// Initialize Gin router
			router := gin.Default()
			
			// Setup routes
			server.SetupRoutes(router, toolRegistry)
			
			// Start server
			if err := router.Run(fmt.Sprintf(":%d", serverPort)); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		},
	}
	
	mcpServerCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start the MCP server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting MCP server on port %d...\n", mcpServerPort)
			
			scriptRunner, err := initializePythonFFI()
			if err != nil {
				log.Printf("Warning: Failed to initialize Python FFI: %v", err)
			}
			
			toolRegistry, err := initializeToolRegistry(scriptRunner)
			if err != nil {
				log.Fatalf("Failed to initialize tool registry: %v", err)
			}
			
			mcpManager, err := mcp.NewManager(toolRegistry)
			if err != nil {
				log.Fatalf("Failed to initialize MCP manager: %v", err)
			}
			
			if err := mcpManager.StartServer(mcpServerPort); err != nil {
				log.Fatalf("Failed to start MCP server: %v", err)
			}
		},
	}
	
	// Add flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "configs/default.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&outputDir, "output", "output", "Directory for output files")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&toolsDir, "tools-dir", "tools", "Directory containing tools")
	rootCmd.PersistentFlags().StringVar(&pythonPath, "python-path", "", "Additional Python path for FFI")
	rootCmd.PersistentFlags().StringVar(&modelName, "model", "gpt-4", "Model name to use for agent")
	rootCmd.PersistentFlags().StringVar(&parserType, "parser", "thought_action", "Parser type (action, thought_action, xml_thought_action, function_calling, json)")
	
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Server port")
	mcpServerCmd.Flags().IntVar(&mcpServerPort, "port", 8081, "MCP server port")
	
	// Add commands to root
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(mcpServerCmd)
	
	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initializePythonFFI() (*python.ScriptRunner, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	
	interpreter, err := python.NewInterpreter()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Python interpreter: %w", err)
	}
	
	toolsPath := filepath.Join(projectRoot, toolsDir)
	scriptRunner, err := python.NewScriptRunner(toolsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize script runner: %w", err)
	}
	
	log.Println("Python FFI initialized successfully")
	return scriptRunner, nil
}

func initializeToolRegistry(scriptRunner *python.ScriptRunner) (*tools.ToolRegistry, error) {
	toolConfig := &tools.ToolConfig{
		ExecutionTimeout: 60 * time.Second,
		MaxOutputSize:    10000,
		ToolsDir:         toolsDir,
		EnableBashTool:   true,
	}
	
	registry, err := tools.NewToolRegistry(toolConfig, scriptRunner)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tool registry: %w", err)
	}
	
	log.Println("Tool registry initialized successfully")
	return registry, nil
}
