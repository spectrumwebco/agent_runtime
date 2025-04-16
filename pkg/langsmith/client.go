package langsmith

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ClientConfig struct {
	APIKey string

	APIUrl string

	ProjectName string

	Disabled bool

	Timeout time.Duration
}

func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		APIKey:      os.Getenv("LANGCHAIN_API_KEY"),
		APIUrl:      getEnvOrDefault("LANGCHAIN_ENDPOINT", "https://api.smith.langchain.com"),
		ProjectName: getEnvOrDefault("LANGCHAIN_PROJECT", "default"),
		Disabled:    os.Getenv("LANGCHAIN_TRACING_V2") != "true",
		Timeout:     30 * time.Second,
	}
}

type Client struct {
	config     *ClientConfig
	httpClient *http.Client
}

func NewClient(config *ClientConfig) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

type Run struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	RunType       string                 `json:"run_type"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       *time.Time             `json:"end_time,omitempty"`
	ExtraData     map[string]interface{} `json:"extra,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Inputs        map[string]interface{} `json:"inputs"`
	Outputs       map[string]interface{} `json:"outputs,omitempty"`
	ParentRunID   string                 `json:"parent_run_id,omitempty"`
	ProjectName   string                 `json:"project_name"`
	Status        string                 `json:"status"`
	Tags          []string               `json:"tags,omitempty"`
	FeedbackStats map[string]interface{} `json:"feedback_stats,omitempty"`
}

type CreateRunRequest struct {
	Name        string                 `json:"name"`
	RunType     string                 `json:"run_type"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	ExtraData   map[string]interface{} `json:"extra,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Inputs      map[string]interface{} `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs,omitempty"`
	ParentRunID string                 `json:"parent_run_id,omitempty"`
	ProjectName string                 `json:"project_name"`
	Status      string                 `json:"status"`
	Tags        []string               `json:"tags,omitempty"`
}

func (c *Client) CreateRun(ctx context.Context, req *CreateRunRequest) (*Run, error) {
	if c.config.Disabled {
		return &Run{
			ID:          uuid.New().String(),
			Name:        req.Name,
			RunType:     req.RunType,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			ExtraData:   req.ExtraData,
			Error:       req.Error,
			Inputs:      req.Inputs,
			Outputs:     req.Outputs,
			ParentRunID: req.ParentRunID,
			ProjectName: req.ProjectName,
			Status:      req.Status,
			Tags:        req.Tags,
		}, nil
	}

	url := fmt.Sprintf("%s/runs", c.config.APIUrl)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var run Run
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &run, nil
}

func (c *Client) UpdateRun(ctx context.Context, runID string, req map[string]interface{}) (*Run, error) {
	if c.config.Disabled {
		return nil, nil
	}

	url := fmt.Sprintf("%s/runs/%s", c.config.APIUrl, runID)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var run Run
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &run, nil
}

func (c *Client) GetRun(ctx context.Context, runID string) (*Run, error) {
	if c.config.Disabled {
		return nil, nil
	}

	url := fmt.Sprintf("%s/runs/%s", c.config.APIUrl, runID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var run Run
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &run, nil
}

func (c *Client) ListRuns(ctx context.Context, projectName string, limit int, offset int) ([]*Run, error) {
	if c.config.Disabled {
		return nil, nil
	}

	url := fmt.Sprintf("%s/runs?project_name=%s&limit=%d&offset=%d", c.config.APIUrl, projectName, limit, offset)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var runs []*Run
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return runs, nil
}

func (c *Client) CreateFeedback(ctx context.Context, runID string, key string, value interface{}, comment string) error {
	if c.config.Disabled {
		return nil
	}

	url := fmt.Sprintf("%s/feedback", c.config.APIUrl)
	body, err := json.Marshal(map[string]interface{}{
		"run_id":  runID,
		"key":     key,
		"value":   value,
		"comment": comment,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
