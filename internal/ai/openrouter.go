package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	OpenRouterEndpoint = "https://openrouter.ai/api/v1/chat/completions"
	
	ModelOpenAI_GPT4o = "openai/gpt-4o"
	ModelAnthropicClaude3 = "anthropic/claude-3-opus"
)

type OpenRouterProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenRouterProvider() *OpenRouterProvider {
	return &OpenRouterProvider{
		apiKey: os.Getenv("OPENROUTER_API_KEY"),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type OpenRouterRequest struct {
	Model    string             `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *OpenRouterProvider) CompletionWithLLM(prompt string, options ...Option) (string, error) {
	opts := newOptions(options...)
	model := opts.Model
	if model == "" {
		model = ModelOpenAI_GPT4o
	}
	
	ctx := context.Background()
	if opts.Context != nil {
		ctx = opts.Context
	}
	
	return p.completionWithLLM(ctx, prompt, model)
}

func (p *OpenRouterProvider) completionWithLLM(ctx context.Context, prompt, model string) (string, error) {
	if p.apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}
	
	req := OpenRouterRequest{
		Model: model,
		Messages: []OpenRouterMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		OpenRouterEndpoint,
		bytes.NewReader(reqData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://agent-runtime.spectrumwebco.com")
	httpReq.Header.Set("X-Title", "Agent Runtime")
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var apiResp OpenRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}
	
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	return apiResp.Choices[0].Message.Content, nil
}

func (p *OpenRouterProvider) CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error) {
	opts := newOptions(options...)
	model := opts.Model
	if model == "" {
		model = ModelOpenAI_GPT4o
	}
	
	ctx := context.Background()
	if opts.Context != nil {
		ctx = opts.Context
	}
	
	return p.completionWithSystemPrompt(ctx, systemPrompt, prompt, model)
}

func (p *OpenRouterProvider) completionWithSystemPrompt(ctx context.Context, systemPrompt, prompt, model string) (string, error) {
	if p.apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}
	
	req := OpenRouterRequest{
		Model: model,
		Messages: []OpenRouterMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		OpenRouterEndpoint,
		bytes.NewReader(reqData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://agent-runtime.spectrumwebco.com")
	httpReq.Header.Set("X-Title", "Agent Runtime")
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var apiResp OpenRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}
	
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	return apiResp.Choices[0].Message.Content, nil
}

func (p *OpenRouterProvider) GetModelForTask(task string) string {
	switch task {
	case "reasoning", "complex_coding", "architecture":
		return ModelAnthropicClaude3
	default:
		return ModelOpenAI_GPT4o
	}
}
