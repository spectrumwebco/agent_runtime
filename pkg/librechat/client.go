package librechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LibreChatClient interface {
	ExecuteCode(ctx context.Context, language string, code string) (string, error)
	ExecuteMultiLanguageCode(ctx context.Context, requests []CodeExecutionRequest) ([]CodeExecutionResponse, error)
	GetSupportedLanguages(ctx context.Context) ([]string, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ExecuteCode(ctx context.Context, language string, code string) (string, error) {
	payload := map[string]interface{}{
		"language": language,
		"code":     code,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal execute code payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", 
			fmt.Sprintf("%s/api/code/execute", c.baseURL), bytes.NewReader(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("create execute code request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute code request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("execute code failed: status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Success bool   `json:"success"`
		Output  string `json:"output"`
		Error   string `json:"error"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode execute code response: %w", err)
	}
	
	if !result.Success {
		return "", fmt.Errorf("execute code failed: %s", result.Error)
	}
	
	return result.Output, nil
}

func (c *Client) ExecuteMultiLanguageCode(ctx context.Context, requests []CodeExecutionRequest) ([]CodeExecutionResponse, error) {
	payload := map[string]interface{}{
		"requests": requests,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal execute multi-language code payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", 
			fmt.Sprintf("%s/api/code/execute-multi", c.baseURL), bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("create execute multi-language code request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute multi-language code request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("execute multi-language code failed: status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Success   bool                   `json:"success"`
		Responses []CodeExecutionResponse `json:"responses"`
		Error     string                 `json:"error"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode execute multi-language code response: %w", err)
	}
	
	if !result.Success {
		return nil, fmt.Errorf("execute multi-language code failed: %s", result.Error)
	}
	
	return result.Responses, nil
}

func (c *Client) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", 
			fmt.Sprintf("%s/api/code/languages", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("create get supported languages request: %w", err)
	}
	
	req.Header.Set("x-api-key", c.apiKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get supported languages request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get supported languages failed: status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Success   bool     `json:"success"`
		Languages []string `json:"languages"`
		Error     string   `json:"error"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode get supported languages response: %w", err)
	}
	
	if !result.Success {
		return nil, fmt.Errorf("get supported languages failed: %s", result.Error)
	}
	
	return result.Languages, nil
}

type CodeExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	ID       string `json:"id,omitempty"`
}

type CodeExecutionResponse struct {
	ID      string `json:"id,omitempty"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}
