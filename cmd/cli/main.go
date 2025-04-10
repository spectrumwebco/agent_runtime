package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/cmd/cli/commands"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sam",
		Short: "Sam Sepiol - Autonomous Software Engineering Agent",
		Long: `Sam Sepiol is an autonomous software engineering agent that can help you with various software engineering tasks.
It is a Go port of the SWE-Agent and SWE-ReX frameworks, providing a high-performance, extensible system for autonomous software engineering.`,
		Version: "0.1.0",
	}

	rootCmd.AddCommand(commands.NewRunCommand())
	rootCmd.AddCommand(commands.NewServeCommand())
	rootCmd.AddCommand(commands.NewToolsCommand())
	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewVersionCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
