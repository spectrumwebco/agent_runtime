package langchain

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

type Chain struct {
	LLM          *Client
	Memory       memory.Memory
	Chain        chains.Chain
	SystemPrompt string
}

type ChainConfig struct {
	Client       *Client
	Memory       memory.Memory
	SystemPrompt string
}

func NewChain(config ChainConfig) (*Chain, error) {
	opts := []chains.ConversationalRetrievalChainOption{}

	if config.Memory != nil {
		opts = append(opts, chains.WithMemory(config.Memory))
	}

	chain := chains.NewConversation(
		config.Client.LLM,
		opts...,
	)

	return &Chain{
		LLM:          config.Client,
		Memory:       config.Memory,
		Chain:        chain,
		SystemPrompt: config.SystemPrompt,
	}, nil
}

func (c *Chain) Run(ctx context.Context, input string) (string, error) {
	inputMap := map[string]any{
		"input": input,
	}

	output, err := chains.Call(ctx, c.Chain, inputMap, llms.WithSystemPrompt(c.SystemPrompt))
	if err != nil {
		return "", fmt.Errorf("failed to run chain: %w", err)
	}

	result, ok := output["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid output type")
	}

	return result, nil
}

func (c *Chain) RunWithTools(ctx context.Context, input string, tools []schema.FunctionDefinition, toolMap map[string]schema.FunctionCallable) (string, error) {
	return c.LLM.ExecuteWithTools(ctx, input, tools, toolMap, llms.WithSystemPrompt(c.SystemPrompt))
}

func (c *Chain) ClearMemory() error {
	if c.Memory == nil {
		return nil
	}
	return c.Memory.Clear()
}

func (c *Chain) GetMemory() memory.Memory {
	return c.Memory
}

func (c *Chain) SetMemory(mem memory.Memory) {
	c.Memory = mem
}
