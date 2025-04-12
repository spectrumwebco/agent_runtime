package ai

import (
	"context"
	"fmt"
)

type Client struct {
	provider Provider
}

func NewClient(provider Provider) *Client {
	return &Client{
		provider: provider,
	}
}

func NewDefaultClient() (*Client, error) {
	provider, err := DefaultProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create default provider: %w", err)
	}
	
	return NewClient(provider), nil
}

func (c *Client) CompletionWithLLM(prompt string, options ...Option) (string, error) {
	return c.provider.CompletionWithLLM(prompt, options...)
}

func (c *Client) CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error) {
	return c.provider.CompletionWithSystemPrompt(systemPrompt, prompt, options...)
}

func (c *Client) GetModelForTask(task string) string {
	return c.provider.GetModelForTask(task)
}

func (c *Client) WithContext(ctx context.Context) Option {
	return WithContext(ctx)
}

func (c *Client) WithModel(model string) Option {
	return WithModel(model)
}
