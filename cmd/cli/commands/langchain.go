package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/langchain"
	"github.com/tmc/langchaingo/tools"
)

func NewLangChainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "langchain",
		Short: "LangChain integration for AI capabilities",
		Long:  `LangChain integration providing advanced AI capabilities for the Kled.io Framework.`,
	}

	cmd.AddCommand(newLangChainGenerateCommand())
	cmd.AddCommand(newLangChainAgentCommand())
	cmd.AddCommand(newLangChainToolsCommand())

	return cmd
}

func newLangChainGenerateCommand() *cobra.Command {
	var modelName string
	var temperature float64

	cmd := &cobra.Command{
		Use:   "generate [prompt]",
		Short: "Generate text using LangChain",
		Long:  `Generate text using LangChain with the specified model and parameters.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := args[0]

			client, err := langchain.NewClient(langchain.ClientConfig{
				Type:        langchain.ModelTypeOpenAI,
				ModelName:   modelName,
				APIKey:      os.Getenv("OPENAI_API_KEY"),
				Temperature: temperature,
			})
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			ctx := context.Background()
			response, err := client.Generate(ctx, prompt)
			if err != nil {
				return fmt.Errorf("failed to generate text: %w", err)
			}

			fmt.Println(response)

			return nil
		},
	}

	cmd.Flags().StringVar(&modelName, "model", "gpt-4o", "Model name to use")
	cmd.Flags().Float64Var(&temperature, "temperature", 0.7, "Temperature for generation")

	return cmd
}

func newLangChainAgentCommand() *cobra.Command {
	var modelName string
	var temperature float64
	var systemPrompt string
	var maxIterations int

	cmd := &cobra.Command{
		Use:   "agent [input]",
		Short: "Run the LangChain agent",
		Long:  `Run the LangChain agent with the specified model and parameters.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			client, err := langchain.NewClient(langchain.ClientConfig{
				Type:        langchain.ModelTypeOpenAI,
				ModelName:   modelName,
				APIKey:      os.Getenv("OPENAI_API_KEY"),
				Temperature: temperature,
			})
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			mem, err := langchain.NewMemory(langchain.MemoryConfig{
				Type:        langchain.MemoryTypeBuffer,
				HumanPrefix: "User",
				AIPrefix:    "Kled",
			})
			if err != nil {
				return fmt.Errorf("failed to create memory: %w", err)
			}

			shellTool := langchain.NewShellTool()
			calculatorTool := langchain.NewCalculatorTool()
			fileTool := langchain.NewFileTool()

			agent, err := langchain.NewAgent(langchain.AgentConfig{
				Client:        client,
				Memory:        mem,
				Tools:         []tools.Tool{shellTool, calculatorTool, fileTool},
				SystemPrompt:  systemPrompt,
				MaxIterations: maxIterations,
			})
			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}

			ctx := context.Background()
			response, err := agent.Run(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to run agent: %w", err)
			}

			fmt.Println(response)

			return nil
		},
	}

	cmd.Flags().StringVar(&modelName, "model", "gpt-4o", "Model name to use")
	cmd.Flags().Float64Var(&temperature, "temperature", 0.7, "Temperature for generation")
	cmd.Flags().StringVar(&systemPrompt, "system-prompt", "You are Kled, a Senior Software Engineering Lead & Technical Authority for AI/ML.", "System prompt for the agent")
	cmd.Flags().IntVar(&maxIterations, "max-iterations", 5, "Maximum number of iterations for the agent")

	return cmd
}

func newLangChainToolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "List available LangChain tools",
		Long:  `List all available tools that can be used with the LangChain agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			shellTool := langchain.NewShellTool()
			searchTool := langchain.NewSearchTool()
			calculatorTool := langchain.NewCalculatorTool()
			fileTool := langchain.NewFileTool()
			httpTool := langchain.NewHTTPTool()
			codeTool := langchain.NewCodeTool()

			tools := []tools.Tool{shellTool, searchTool, calculatorTool, fileTool, httpTool, codeTool}
			fmt.Println("Available LangChain tools:")
			for _, tool := range tools {
				fmt.Printf("- %s: %s\n", tool.Name(), tool.Description())
			}

			return nil
		},
	}

	return cmd
}
