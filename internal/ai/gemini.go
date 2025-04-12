package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	GeminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models"
	
	ModelGemini15Pro = "gemini-1.5-pro"
	ModelGemini15Flash = "gemini-1.5-flash"
)

type GeminiProvider struct {
	apiKey string
	client *http.Client
}

func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{
		apiKey: os.Getenv("GEMINI_API_KEY"),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Role  string `json:"role,omitempty"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *GeminiProvider) CompletionWithLLM(prompt string, options ...Option) (string, error) {
	opts := newOptions(options...)
	model := opts.Model
	if model == "" {
		model = ModelGemini15Pro
	}
	
	ctx := context.Background()
	if opts.Context != nil {
		ctx = opts.Context
	}
	
	return p.completionWithLLM(ctx, prompt, model)
}

func (p *GeminiProvider) completionWithLLM(ctx context.Context, prompt, model string) (string, error) {
	if p.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}
	
	req := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{
						Text: prompt,
					},
				},
			},
		},
	}
	
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", GeminiEndpoint, model, p.apiKey)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		strings.NewReader(string(reqData)),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	
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
	
	var apiResp GeminiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}
	
	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	
	return apiResp.Candidates[0].Content.Parts[0].Text, nil
}

func (p *GeminiProvider) CompletionWithSystemPrompt(systemPrompt, prompt string, options ...Option) (string, error) {
	opts := newOptions(options...)
	model := opts.Model
	if model == "" {
		model = ModelGemini15Pro
	}
	
	ctx := context.Background()
	if opts.Context != nil {
		ctx = opts.Context
	}
	
	combinedPrompt := fmt.Sprintf("%s\n\n%s", systemPrompt, prompt)
	return p.completionWithLLM(ctx, combinedPrompt, model)
}

func (p *GeminiProvider) GetModelForTask(task string) string {
	switch task {
	case "coding", "complex_coding":
		return ModelGemini15Pro
	default:
		return ModelGemini15Flash
	}
}
