package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/ai"
)

func NewAICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI capabilities for the Kled.io Framework",
		Long:  `AI integration providing advanced capabilities for the Kled.io Framework.`,
	}

	cmd.AddCommand(newAIGenerateCommand())
	cmd.AddCommand(newAIProvidersCommand())

	return cmd
}

func newAIGenerateCommand() *cobra.Command {
	var systemPrompt string
	var model string
	var providerType string

	cmd := &cobra.Command{
		Use:   "generate [prompt]",
		Short: "Generate text using AI",
		Long:  `Generate text using AI with the specified model and parameters.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := args[0]

			var client *ai.Client
			var err error

			if providerType != "" {
				provider, err := ai.ProviderFactory(ai.ProviderType(providerType))
				if err != nil {
					return fmt.Errorf("failed to create provider: %w", err)
				}
				client = ai.NewClient(provider)
			} else {
				client, err = ai.NewDefaultClient()
				if err != nil {
					return fmt.Errorf("failed to create default client: %w", err)
				}
			}

			var options []ai.Option
			if model != "" {
				options = append(options, client.WithModel(model))
			}

			var response string
			if systemPrompt != "" {
				response, err = client.CompletionWithSystemPrompt(systemPrompt, prompt, options...)
			} else {
				response, err = client.CompletionWithLLM(prompt, options...)
			}

			if err != nil {
				return fmt.Errorf("failed to generate text: %w", err)
			}

			fmt.Println(response)
			return nil
		},
	}

	cmd.Flags().StringVar(&systemPrompt, "system", "", "System prompt to use")
	cmd.Flags().StringVar(&model, "model", "", "Model to use")
	cmd.Flags().StringVar(&providerType, "provider", "", "Provider to use (openrouter, klusterai, bedrock, gemini, langchain)")

	return cmd
}

func newAIProvidersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "List available AI providers",
		Long:  `List all available AI providers that can be used with the AI command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Available AI providers:")
			fmt.Println("- openrouter: OpenRouter API")
			fmt.Println("- klusterai: KlusterAI (default)")
			fmt.Println("- bedrock: AWS Bedrock")
			fmt.Println("- gemini: Google Gemini")
			fmt.Println("- langchain: LangChain integration")
			
			fmt.Println("\nAPI key status:")
			
			if os.Getenv("OPENAI_API_KEY") != "" {
				fmt.Println("- OpenAI API key: ✓")
			} else {
				fmt.Println("- OpenAI API key: ✗")
			}
			
			if os.Getenv("OPENROUTER_API_KEY") != "" {
				fmt.Println("- OpenRouter API key: ✓")
			} else {
				fmt.Println("- OpenRouter API key: ✗")
			}
			
			if os.Getenv("KLUSTERAI_API_KEY") != "" {
				fmt.Println("- KlusterAI API key: ✓")
			} else {
				fmt.Println("- KlusterAI API key: ✗")
			}
			
			if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
				fmt.Println("- AWS credentials: ✓")
			} else {
				fmt.Println("- AWS credentials: ✗")
			}
			
			if os.Getenv("GEMINI_API_KEY") != "" {
				fmt.Println("- Gemini API key: ✓")
			} else {
				fmt.Println("- Gemini API key: ✗")
			}
			
			return nil
		},
	}

	return cmd
}
