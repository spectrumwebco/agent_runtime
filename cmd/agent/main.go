package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

var (
	configPath string
	port       int
	host       string
	logLevel   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent Runtime - A Go-based autonomous agent framework",
		Long: `Agent Runtime is a high-performance Go-based backend framework for building 
autonomous software engineering agents. It leverages the Model Context Protocol (MCP) 
to create a powerful, extensible system that can integrate with various LLMs and tools.`,
	}
	
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to configuration file")
	rootCmd.PersistentFlags().IntVar(&port, "port", 8080, "Port to run the server on")
	rootCmd.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "Host to run the server on")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	
	rootCmd.AddCommand(serveCmd())
	rootCmd.AddCommand(versionCmd())
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Agent Runtime server",
		Long:  `Start the Agent Runtime server to handle agent requests.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(configPath)
			if err != nil {
				fmt.Printf("Failed to load configuration: %v\n", err)
				os.Exit(1)
			}
			
			if port != 8080 {
				cfg.Server.Port = port
			}
			if host != "0.0.0.0" {
				cfg.Server.Host = host
			}
			if logLevel != "info" {
				cfg.Logging.Level = logLevel
			}
			
			fmt.Printf("Starting server on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
		},
	}
	
	return cmd
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version information of the Agent Runtime.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Agent Runtime v0.1.0")
		},
	}
	
	return cmd
}
