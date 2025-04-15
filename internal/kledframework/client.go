package kledframework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a client for a microservice
type Client struct {
	Name       string
	Registry   *Registry
	HTTPClient *http.Client
}

// ClientConfig contains configuration for a client
type ClientConfig struct {
	Name     string
	Registry *Registry
	Timeout  time.Duration
}

// NewClient creates a new client
func NewClient(config ClientConfig) *Client {
	// Set default values
	if config.Timeout == 0 {
		config.Timeout = time.Second * 5
	}

	// Create the HTTP client
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		Name:       config.Name,
		Registry:   config.Registry,
		HTTPClient: httpClient,
	}
}

// Call calls a method on a service
func (c *Client) Call(serviceName, endpoint string, req, rsp interface{}) error {
	// Get the service from the registry
	service, err := c.Registry.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("http://%s%s", service.Address, endpoint)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	httpRsp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer httpRsp.Body.Close()

	// Check the response status
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("service returned status %d", httpRsp.StatusCode)
	}

	// Read the response body
	rspBody, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(rspBody, rsp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
