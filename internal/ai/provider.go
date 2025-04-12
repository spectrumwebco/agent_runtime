package ai

import (
	"fmt"
	"os"
)

type ProviderType string

const (
	ProviderOpenRouter ProviderType = "openrouter"
	ProviderKlusterAI  ProviderType = "klusterai"
	ProviderBedrock    ProviderType = "bedrock"
	ProviderGemini     ProviderType = "gemini"
)

type Provider interface {
	CompletionWithLLM(prompt string, options ...Option) (string, error)
	CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error)
	GetModelForTask(task string) string
}

func ProviderFactory(providerType ProviderType) (Provider, error) {
	switch providerType {
	case ProviderKlusterAI:
		return NewKlusterAIProvider(), nil
	case ProviderOpenRouter:
		return NewOpenRouterProvider(), nil
	case ProviderBedrock:
		return NewBedrockProvider(), nil
	case ProviderGemini:
		return NewGeminiProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

func DefaultProvider() (Provider, error) {
	providerType := os.Getenv("AI_PROVIDER")
	if providerType == "" {
		providerType = string(ProviderKlusterAI) // Default to KlusterAI
	}
	
	return ProviderFactory(ProviderType(providerType))
}
