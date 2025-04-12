package ai

import (
	"fmt"
	"os"
)

const (
	ModelAnthropicClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
	ModelAnthropicClaude3Haiku = "anthropic.claude-3-haiku-20240307-v1:0"
)

type BedrockProvider struct {
	region string
}

func NewBedrockProvider() *BedrockProvider {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
	}
	
	return &BedrockProvider{
		region: region,
	}
}

func (p *BedrockProvider) CompletionWithLLM(prompt string, options ...Option) (string, error) {
	return fmt.Sprintf("Bedrock completion for: %s", prompt), nil
}

func (p *BedrockProvider) CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error) {
	return fmt.Sprintf("Bedrock completion with system prompt for: %s", prompt), nil
}

func (p *BedrockProvider) GetModelForTask(task string) string {
	switch task {
	case "reasoning", "complex_coding", "architecture":
		return ModelAnthropicClaude3Sonnet
	default:
		return ModelAnthropicClaude3Haiku
	}
}
