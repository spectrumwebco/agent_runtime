package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spectrumwebco/agent_runtime/internal/server"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	port := flag.Int("port", 8080, "Port to run the server on")
	host := flag.String("host", "0.0.0.0", "Host to run the server on")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	
	flag.Parse()
	
	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}
	
	subcommand := flag.Arg(0)
	
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	if *port != 8080 {
		cfg.Server.Port = *port
	}
	if *host != "0.0.0.0" {
		cfg.Server.Host = *host
	}
	if *logLevel != "info" {
		cfg.Logging.Level = *logLevel
	}
	
	switch subcommand {
	case "serve":
		serve(cfg)
	case "version":
		printVersion()
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func serve(cfg *config.Config) {
	s, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := s.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func printVersion() {
	fmt.Println("Agent Runtime v0.1.0")
}

func printUsage() {
	fmt.Println("Usage: agent_runtime [options] <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  serve    Start the Agent Runtime server")
	fmt.Println("  version  Print the version information")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}
