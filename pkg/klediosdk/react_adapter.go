package klediosdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ReactRuntimeAdapter struct {
	config *RuntimeAdapterConfig

	connected bool

	httpClient *http.Client

	wsConn *websocket.Conn

	baseURL string

	authToken string

	eventHandlers map[string]map[string]EventHandler

	eventChan chan *Event

	ctx    context.Context
	cancel context.CancelFunc

	mu sync.RWMutex
}

func NewReactRuntimeAdapter() *ReactRuntimeAdapter {
	ctx, cancel := context.WithCancel(context.Background())
	return &ReactRuntimeAdapter{
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		eventHandlers: make(map[string]map[string]EventHandler),
		eventChan:     make(chan *Event, 100),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (a *ReactRuntimeAdapter) Initialize(ctx context.Context, config *RuntimeAdapterConfig) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if config == nil {
		return errors.New("config cannot be nil")
	}

	a.config = config

	baseURL, ok := config.ConnectionConfig["base_url"].(string)
	if !ok || baseURL == "" {
		return errors.New("base_url is required in connection config")
	}
	a.baseURL = baseURL

	authToken, ok := config.AuthConfig["token"].(string)
	if !ok || authToken == "" {
		return errors.New("auth token is required in auth config")
	}
	a.authToken = authToken

	return nil
}

func (a *ReactRuntimeAdapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.connected {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.baseURL+"/api/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.authToken)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to React runtime API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("React runtime API health check failed with status code: %d", resp.StatusCode)
	}

	wsURL := a.baseURL + "/api/events/ws"
	wsURL = "ws" + wsURL[4:] // Replace http with ws
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	header := http.Header{}
	header.Set("Authorization", "Bearer "+a.authToken)
	wsConn, _, err := dialer.DialContext(ctx, wsURL, header)
	if err != nil {
		return fmt.Errorf("failed to connect to React runtime WebSocket: %w", err)
	}

	a.wsConn = wsConn
	a.connected = true

	go a.listenForEvents()

	return nil
}

func (a *ReactRuntimeAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return nil
	}

	a.cancel()

	if a.wsConn != nil {
		err := a.wsConn.Close()
		if err != nil {
			return fmt.Errorf("failed to close WebSocket connection: %w", err)
		}
		a.wsConn = nil
	}

	a.connected = false
	return nil
}

func (a *ReactRuntimeAdapter) IsConnected() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.connected
}

func (a *ReactRuntimeAdapter) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return nil, errors.New("adapter is not connected")
	}

	payload := map[string]interface{}{
		"task_id":      task.ID,
		"task_type":    task.Type,
		"description":  task.Description,
		"input":        task.Input,
		"agent_id":     task.AgentID,
		"timeout":      task.Timeout.Milliseconds(),
		"metadata":     task.Metadata,
		"runtime_type": string(RuntimeTypeReact),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/tasks", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create task execution request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.authToken)
	req.Body = http.NoBody

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute task in React runtime: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("task execution failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		TaskID        string                 `json:"task_id"`
		AgentID       string                 `json:"agent_id"`
		Status        string                 `json:"status"`
		Output        map[string]interface{} `json:"output"`
		Error         string                 `json:"error"`
		ExecutionTime int64                  `json:"execution_time"`
		Metadata      map[string]interface{} `json:"metadata"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode task result: %w", err)
	}

	taskResult := &TaskResult{
		TaskID:        result.TaskID,
		AgentID:       result.AgentID,
		Status:        result.Status,
		Output:        result.Output,
		Error:         result.Error,
		Runtime:       RuntimeTypeReact,
		ExecutionTime: time.Duration(result.ExecutionTime) * time.Millisecond,
		Metadata:      result.Metadata,
	}

	return taskResult, nil
}

func (a *ReactRuntimeAdapter) GetState(ctx context.Context, key string) (interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return nil, errors.New("adapter is not connected")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.baseURL+"/api/state/"+key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create state request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.authToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from React runtime: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("state retrieval failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state result: %w", err)
	}

	return result.Value, nil
}

func (a *ReactRuntimeAdapter) SetState(ctx context.Context, key string, value interface{}) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return errors.New("adapter is not connected")
	}

	payload := map[string]interface{}{
		"key":   key,
		"value": value,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal state payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", a.baseURL+"/api/state", nil)
	if err != nil {
		return fmt.Errorf("failed to create state request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.authToken)
	req.Body = http.NoBody

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to set state in React runtime: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("state update failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *ReactRuntimeAdapter) SubscribeToEvents(ctx context.Context, eventType string, handler EventHandler) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return "", errors.New("adapter is not connected")
	}

	subscriptionID := uuid.New().String()

	if _, ok := a.eventHandlers[eventType]; !ok {
		a.eventHandlers[eventType] = make(map[string]EventHandler)

		payload := map[string]interface{}{
			"event_type":      eventType,
			"subscription_id": subscriptionID,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal subscription payload: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/events/subscribe", nil)
		if err != nil {
			return "", fmt.Errorf("failed to create subscription request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+a.authToken)
		req.Body = http.NoBody

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to subscribe to events in React runtime: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("event subscription failed with status code: %d", resp.StatusCode)
		}
	}

	a.eventHandlers[eventType][subscriptionID] = handler

	return subscriptionID, nil
}

func (a *ReactRuntimeAdapter) UnsubscribeFromEvents(ctx context.Context, subscriptionID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.connected {
		return errors.New("adapter is not connected")
	}

	var eventType string
	for et, handlers := range a.eventHandlers {
		if _, ok := handlers[subscriptionID]; ok {
			eventType = et
			break
		}
	}

	if eventType == "" {
		return errors.New("subscription not found")
	}

	delete(a.eventHandlers[eventType], subscriptionID)

	if len(a.eventHandlers[eventType]) == 0 {
		delete(a.eventHandlers, eventType)

		payload := map[string]interface{}{
			"event_type":      eventType,
			"subscription_id": subscriptionID,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal unsubscription payload: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/events/unsubscribe", nil)
		if err != nil {
			return fmt.Errorf("failed to create unsubscription request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+a.authToken)
		req.Body = http.NoBody

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to unsubscribe from events in React runtime: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("event unsubscription failed with status code: %d", resp.StatusCode)
		}
	}

	return nil
}

func (a *ReactRuntimeAdapter) PublishEvent(ctx context.Context, event *Event) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.connected {
		return errors.New("adapter is not connected")
	}

	payload := map[string]interface{}{
		"id":        event.ID,
		"type":      event.Type,
		"source":    event.Source,
		"timestamp": event.Timestamp.Format(time.RFC3339),
		"data":      event.Data,
		"runtime":   string(RuntimeTypeReact),
		"metadata":  event.Metadata,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/events/publish", nil)
	if err != nil {
		return fmt.Errorf("failed to create event publish request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.authToken)
	req.Body = http.NoBody

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to publish event to React runtime: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("event publication failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *ReactRuntimeAdapter) listenForEvents() {
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			_, message, err := a.wsConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					a.mu.Lock()
					a.connected = false
					a.wsConn = nil
					a.mu.Unlock()

					time.Sleep(5 * time.Second)
					err := a.Connect(context.Background())
					if err != nil {
						continue
					}
				}
				continue
			}

			var event Event
			err = json.Unmarshal(message, &event)
			if err != nil {
				continue
			}

			a.mu.RLock()
			handlers, ok := a.eventHandlers[event.Type]
			a.mu.RUnlock()

			if ok {
				for _, handler := range handlers {
					go func(h EventHandler) {
						_ = h(context.Background(), &event)
					}(handler)
				}
			}
		}
	}
}
