package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

const (
	ProviderLangChain ProviderType = "langchain"
)

type LangChainProvider struct {
	client  llms.LLM
	options *llms.CallOptions
	mutex   sync.RWMutex
}

func NewLangChainProvider() *LangChainProvider {
	options := &llms.CallOptions{
		Temperature: 0.7,
	}

	return &LangChainProvider{
		options: options,
	}
}

func (p *LangChainProvider) initClient() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.client != nil {
		return nil
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	client, err := openai.New(openai.WithToken(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	p.client = client
	return nil
}

func (p *LangChainProvider) CompletionWithLLM(prompt string, options ...Option) (string, error) {
	if err := p.initClient(); err != nil {
		return "", err
	}

	callOptions := p.applyOptions(options...)

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	result, err := p.client.Call(context.Background(), prompt, callOptions)
	if err != nil {
		return "", fmt.Errorf("failed to generate completion: %w", err)
	}

	return result, nil
}

func (p *LangChainProvider) CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error) {
	if err := p.initClient(); err != nil {
		return "", err
	}

	callOptions := p.applyOptions(options...)

	messages := []schema.ChatMessage{
		schema.SystemChatMessage{Content: systemPrompt},
		schema.HumanChatMessage{Content: prompt},
	}

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	chatLLM, ok := p.client.(llms.ChatLLM)
	if !ok {
		combinedPrompt := fmt.Sprintf("System: %s\n\nUser: %s", systemPrompt, prompt)
		result, err := p.client.Call(context.Background(), combinedPrompt, callOptions)
		if err != nil {
			return "", fmt.Errorf("failed to generate completion: %w", err)
		}
		return result, nil
	}

	response, err := chatLLM.GenerateChat(context.Background(), messages, callOptions)
	if err != nil {
		return "", fmt.Errorf("failed to generate chat completion: %w", err)
	}

	if len(response.Generations) == 0 || len(response.Generations[0].ChatMessages) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return response.Generations[0].ChatMessages[0].GetContent(), nil
}

func (p *LangChainProvider) GetModelForTask(task string) string {
	task = strings.ToLower(task)

	if strings.Contains(task, "complex") || 
	   strings.Contains(task, "reasoning") || 
	   strings.Contains(task, "planning") {
		return "gpt-4"
	}

	return "gpt-3.5-turbo"
}

func (p *LangChainProvider) applyOptions(options ...Option) *llms.CallOptions {
	callOptions := &llms.CallOptions{
		Temperature: p.options.Temperature,
	}

	for _, option := range options {
		option((*optionsHolder)(callOptions))
	}

	return callOptions
}

func init() {
	originalProviderFactory := ProviderFactory
	ProviderFactory = func(providerType ProviderType) (Provider, error) {
		if providerType == ProviderLangChain {
			return NewLangChainProvider(), nil
		}
		return originalProviderFactory(providerType)
	}
}
