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
	KlusterAIEndpoint = "https://api.kluster.ai/v1/chat/completions"
	
	ModelLlama4Scout = "llama-4-scout"
	
	ModelLlama4Maverick = "llama-4-maverick"
)

type KlusterAIProvider struct {
	apiKey string
	client *http.Client
}

func NewKlusterAIProvider() *KlusterAIProvider {
	return &KlusterAIProvider{
		apiKey: os.Getenv("KLUSTERAI_API_KEY"),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *KlusterAIProvider) CompletionWithLLM(prompt string, options ...Option) (string, error) {
	opts := newOptions(options...)
	model := opts.Model
	if model == "" {
		model = ModelLlama4Scout
	}
	
	ctx := context.Background()
	if opts.Context != nil {
		ctx = opts.Context
	}
	
	return p.completionWithLLM(ctx, prompt, model)
}

func (p *KlusterAIProvider) completionWithLLM(ctx context.Context, prompt, model string) (string, error) {
	if p.apiKey == "" {
		return "", fmt.Errorf("KLUSTERAI_API_KEY environment variable not set")
	}
	
	if model != ModelLlama4Scout && model != ModelLlama4Maverick && model != "llama-4" {
		model = ModelLlama4Scout
	}
	
	req := CompletionRequest{
		Model: model,
		Messages: []Message{
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
		KlusterAIEndpoint,
		bytes.NewReader(reqData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://agent-runtime.spectrumwebco.com")
	
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
	
	var apiResp CompletionResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	return apiResp.Choices[0].Message.Content, nil
}

func (p *KlusterAIProvider) CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error) {
	opts := newOptions(options...)
	model := opts.Model
	if model == "" {
		model = ModelLlama4Maverick
	}
	
	ctx := context.Background()
	if opts.Context != nil {
		ctx = opts.Context
	}
	
	return p.completionWithSystemPrompt(ctx, systemPrompt, prompt, model)
}

func (p *KlusterAIProvider) completionWithSystemPrompt(ctx context.Context, systemPrompt, prompt, model string) (string, error) {
	if p.apiKey == "" {
		return "", fmt.Errorf("KLUSTERAI_API_KEY environment variable not set")
	}
	
	if model != ModelLlama4Scout && model != ModelLlama4Maverick && model != "llama-4" {
		model = ModelLlama4Maverick
	}
	
	req := CompletionRequest{
		Model: model,
		Messages: []Message{
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
		KlusterAIEndpoint,
		bytes.NewReader(reqData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://agent-runtime.spectrumwebco.com")
	
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
	
	var apiResp CompletionResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	return apiResp.Choices[0].Message.Content, nil
}

func (p *KlusterAIProvider) GetModelForTask(task string) string {
	switch task {
	case "reasoning":
		return ModelLlama4Maverick
	case "standard":
		return ModelLlama4Scout
	default:
		return ModelLlama4Scout
	}
}
