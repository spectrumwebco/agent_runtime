package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/cmd/cli/commands"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kled",
		Short: "Kled - Senior Software Engineering Lead & Technical Authority for AI/ML",
		Long: `Kled is a Senior Software Engineering Lead & Technical Authority for AI/ML that can help you with various software engineering tasks.
It is a Go port of the SWE-Agent and SWE-ReX frameworks, providing a high-performance, extensible system for autonomous software engineering.`,
		Version: "0.1.0",
	}

	rootCmd.AddCommand(commands.NewRunCommand())
	rootCmd.AddCommand(commands.NewServeCommand())
	rootCmd.AddCommand(commands.NewToolsCommand())
	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewVersionCommand())
	rootCmd.AddCommand(commands.NewLangChainCommand())
	rootCmd.AddCommand(commands.NewAICommand())
	rootCmd.AddCommand(commands.NewLangGraphCommand())
	rootCmd.AddCommand(commands.NewActorModelCommand())
	rootCmd.AddCommand(commands.NewGORMCommand())
	rootCmd.AddCommand(commands.NewGoMicroCommand())
	rootCmd.AddCommand(commands.NewOrchestrationCommand())
	rootCmd.AddCommand(commands.NewLeetCodeCommand())
	rootCmd.AddCommand(commands.NewBubbleTeaCommand())
	rootCmd.AddCommand(commands.NewK9sCommand())
	rootCmd.AddCommand(commands.NewCollyCommand())
	rootCmd.AddCommand(commands.NewGodsCommand())
	rootCmd.AddCommand(commands.NewKopsCommand())
	rootCmd.AddCommand(commands.NewHydraCommand())
	rootCmd.AddCommand(commands.NewClairCommand())
	rootCmd.AddCommand(commands.NewCoderCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
