package langraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type DjangoAgentConfig struct {
	BaseURL     string                 `json:"base_url"`
	APIKey      string                 `json:"api_key"`
	Timeout     time.Duration          `json:"timeout"`
	Headers     map[string]string      `json:"headers,omitempty"`
	AgentConfig AgentConfig            `json:"agent_config"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type DjangoAgent struct {
	Config      DjangoAgentConfig      `json:"config"`
	State       map[string]interface{} `json:"state,omitempty"`
	Tools       []Tool                 `json:"-"`
	EventStream EventStream            `json:"-"`
	httpClient  *http.Client           `json:"-"`
	stateLock   sync.RWMutex           `json:"-"`
}

type DjangoAgentRequest struct {
	AgentID   string                 `json:"agent_id"`
	Input     map[string]interface{} `json:"input"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	TaskID    string                 `json:"task_id,omitempty"`
}

type DjangoAgentResponse struct {
	AgentID   string                 `json:"agent_id"`
	Output    map[string]interface{} `json:"output"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Error     string                 `json:"error,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	TaskID    string                 `json:"task_id,omitempty"`
}

func NewDjangoAgent(config DjangoAgentConfig, eventStream EventStream) *DjangoAgent {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &DjangoAgent{
		Config:      config,
		State:       make(map[string]interface{}),
		Tools:       []Tool{},
		EventStream: eventStream,
		httpClient:  httpClient,
	}
}

func (a *DjangoAgent) Process(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"input":      input,
			"agent_id":   a.Config.AgentConfig.ID,
			"agent_name": a.Config.AgentConfig.Name,
			"agent_role": a.Config.AgentConfig.Role,
		},
		map[string]string{
			"agent_id": a.Config.AgentConfig.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	request := DjangoAgentRequest{
		AgentID:  a.Config.AgentConfig.ID,
		Input:    input,
		Metadata: a.Config.Metadata,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/agent/process/", a.Config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Config.APIKey))
	for key, value := range a.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response DjangoAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("agent error: %s", response.Error)
	}

	resultEvent := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"result":     response.Output,
			"metadata":   response.Metadata,
			"agent_id":   a.Config.AgentConfig.ID,
			"agent_name": a.Config.AgentConfig.Name,
			"agent_role": a.Config.AgentConfig.Role,
		},
		map[string]string{
			"agent_id": a.Config.AgentConfig.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(resultEvent); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	if response.Metadata != nil {
		a.UpdateState(response.Metadata)
	}

	return response.Output, nil
}

func (a *DjangoAgent) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (map[string]interface{}, error) {
	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool":       toolName,
			"parameters": params,
			"agent_id":   a.Config.AgentConfig.ID,
			"agent_name": a.Config.AgentConfig.Name,
			"agent_role": a.Config.AgentConfig.Role,
		},
		map[string]string{
			"agent_id": a.Config.AgentConfig.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	request := DjangoAgentRequest{
		AgentID: a.Config.AgentConfig.ID,
		Input: map[string]interface{}{
			"tool":       toolName,
			"parameters": params,
		},
		Metadata: a.Config.Metadata,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/agent/execute_tool/", a.Config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Config.APIKey))
	for key, value := range a.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response DjangoAgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("agent error: %s", response.Error)
	}

	resultEvent := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool":       toolName,
			"result":     response.Output,
			"metadata":   response.Metadata,
			"agent_id":   a.Config.AgentConfig.ID,
			"agent_name": a.Config.AgentConfig.Name,
			"agent_role": a.Config.AgentConfig.Role,
		},
		map[string]string{
			"agent_id": a.Config.AgentConfig.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(resultEvent); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	if response.Metadata != nil {
		a.UpdateState(response.Metadata)
	}

	return response.Output, nil
}

func (a *DjangoAgent) SetState(state map[string]interface{}) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	a.State = state
}

func (a *DjangoAgent) UpdateState(updates map[string]interface{}) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	
	if a.State == nil {
		a.State = make(map[string]interface{})
	}
	
	for k, v := range updates {
		a.State[k] = v
	}
}

func (a *DjangoAgent) GetState() map[string]interface{} {
	a.stateLock.RLock()
	defer a.stateLock.RUnlock()
	
	stateCopy := make(map[string]interface{})
	for k, v := range a.State {
		stateCopy[k] = v
	}
	
	return stateCopy
}

type DjangoAgentFactory struct {
	BaseURL     string
	APIKey      string
	Headers     map[string]string
	Timeout     time.Duration
	EventStream EventStream
	agents      map[string]*DjangoAgent
	agentsMutex sync.RWMutex
}

func NewDjangoAgentFactory(baseURL, apiKey string, eventStream EventStream) *DjangoAgentFactory {
	return &DjangoAgentFactory{
		BaseURL:     baseURL,
		APIKey:      apiKey,
		Headers:     make(map[string]string),
		Timeout:     30 * time.Second,
		EventStream: eventStream,
		agents:      make(map[string]*DjangoAgent),
	}
}

func (f *DjangoAgentFactory) CreateAgent(config AgentConfig) (*DjangoAgent, error) {
	djangoConfig := DjangoAgentConfig{
		BaseURL:     f.BaseURL,
		APIKey:      f.APIKey,
		Timeout:     f.Timeout,
		Headers:     f.Headers,
		AgentConfig: config,
		Metadata:    make(map[string]interface{}),
	}

	agent := NewDjangoAgent(djangoConfig, f.EventStream)

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[config.ID] = agent

	return agent, nil
}

func (f *DjangoAgentFactory) GetAgent(id string) (*DjangoAgent, error) {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent %s does not exist", id)
	}

	return agent, nil
}

func (f *DjangoAgentFactory) ListAgents() []*DjangoAgent {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	var agents []*DjangoAgent
	for _, agent := range f.agents {
		agents = append(agents, agent)
	}

	return agents
}

func (f *DjangoAgentFactory) RemoveAgent(id string) error {
	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()

	if _, exists := f.agents[id]; !exists {
		return fmt.Errorf("agent %s does not exist", id)
	}

	delete(f.agents, id)
	return nil
}
