package kledframework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/djangobridge"
)

type DjangoBridge struct {
	config     *DjangoBridgeConfig
	client     *djangobridge.Client
	httpServer *http.Server
	isRunning  bool
	mutex      sync.RWMutex
	stopCh     chan struct{}
}

type DjangoBridgeConfig struct {
	GRPCAddress    string
	HTTPAddress    string
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	AllowedOrigins []string
}

func DefaultDjangoBridgeConfig() *DjangoBridgeConfig {
	return &DjangoBridgeConfig{
		GRPCAddress:    "localhost:50051",
		HTTPAddress:    ":8082",
		Timeout:        30 * time.Second,
		MaxRetries:     3,
		RetryDelay:     time.Second,
		AllowedOrigins: []string{"*"},
	}
}

func NewDjangoBridge(config *DjangoBridgeConfig) (*DjangoBridge, error) {
	if config == nil {
		config = DefaultDjangoBridgeConfig()
	}

	client, err := djangobridge.NewClient(config.GRPCAddress, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Django bridge client: %w", err)
	}

	bridge := &DjangoBridge{
		config:    config,
		client:    client,
		isRunning: false,
		mutex:     sync.RWMutex{},
		stopCh:    make(chan struct{}),
	}

	return bridge, nil
}

func (db *DjangoBridge) Start(ctx context.Context) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.isRunning {
		return fmt.Errorf("Django bridge is already running")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/django/execute", db.handleExecutePythonCode)
	mux.HandleFunc("/api/django/query", db.handleQueryDjangoModel)
	mux.HandleFunc("/api/django/create", db.handleCreateDjangoModel)
	mux.HandleFunc("/api/django/update", db.handleUpdateDjangoModel)
	mux.HandleFunc("/api/django/delete", db.handleDeleteDjangoModel)
	mux.HandleFunc("/api/django/command", db.handleExecuteDjangoManagementCommand)
	mux.HandleFunc("/api/django/settings", db.handleGetDjangoSettings)
	mux.HandleFunc("/api/django/agent", db.handleCreateDjangoAgent)

	db.httpServer = &http.Server{
		Addr:    db.config.HTTPAddress,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting Django bridge HTTP server on %s", db.config.HTTPAddress)
		if err := db.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Django bridge HTTP server error: %v", err)
		}
	}()

	db.isRunning = true

	return nil
}

func (db *DjangoBridge) Stop() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if !db.isRunning {
		return nil
	}

	close(db.stopCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := db.client.Close(); err != nil {
		return fmt.Errorf("failed to close gRPC client: %w", err)
	}

	db.isRunning = false

	return nil
}

func (db *DjangoBridge) ExecutePythonCode(ctx context.Context, code string, timeout time.Duration) (map[string]interface{}, error) {
	if timeout == 0 {
		timeout = db.config.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var result map[string]interface{}
	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		result, err = db.client.ExecutePythonCode(ctx, code)
		if err == nil {
			return result, nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return nil, fmt.Errorf("failed to execute Python code after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) QueryDjangoModel(ctx context.Context, modelName string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var results []map[string]interface{}
	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		results, err = db.client.QueryDjangoModel(ctx, modelName, filters)
		if err == nil {
			return results, nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return nil, fmt.Errorf("failed to query Django model after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) CreateDjangoModel(ctx context.Context, modelName string, data map[string]interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var result map[string]interface{}
	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		result, err = db.client.CreateDjangoModel(ctx, modelName, data)
		if err == nil {
			return result, nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return nil, fmt.Errorf("failed to create Django model after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) UpdateDjangoModel(ctx context.Context, modelName string, id int, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		err = db.client.UpdateDjangoModel(ctx, modelName, id, data)
		if err == nil {
			return nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return fmt.Errorf("failed to update Django model after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) DeleteDjangoModel(ctx context.Context, modelName string, id int) error {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		err = db.client.DeleteDjangoModel(ctx, modelName, id)
		if err == nil {
			return nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return fmt.Errorf("failed to delete Django model after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) ExecuteDjangoManagementCommand(ctx context.Context, command string, args []string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var result string
	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		result, err = db.client.ExecuteDjangoManagementCommand(ctx, command, args)
		if err == nil {
			return result, nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return "", fmt.Errorf("failed to execute Django management command after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) GetDjangoSettings(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var settings map[string]interface{}
	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		settings, err = db.client.GetDjangoSettings(ctx)
		if err == nil {
			return settings, nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return nil, fmt.Errorf("failed to get Django settings after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) CreateDjangoAgent(ctx context.Context, id, name, role string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, db.config.Timeout)
	defer cancel()

	var agent map[string]interface{}
	var err error

	for i := 0; i <= db.config.MaxRetries; i++ {
		agent, err = db.client.CreateDjangoAgent(ctx, id, name, role)
		if err == nil {
			return agent, nil
		}

		if i < db.config.MaxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(db.config.RetryDelay):
			}
		}
	}

	return nil, fmt.Errorf("failed to create Django agent after %d retries: %w", db.config.MaxRetries, err)
}

func (db *DjangoBridge) Close() error {
	return db.Stop()
}


func (db *DjangoBridge) handleExecutePythonCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Code    string        `json:"code"`
		Timeout time.Duration `json:"timeout,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	result, err := db.ExecutePythonCode(r.Context(), request.Code, request.Timeout)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute Python code: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Failed to encode result: %v", err)
	}
}

func (db *DjangoBridge) handleQueryDjangoModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ModelName string                 `json:"model_name"`
		Filters   map[string]interface{} `json:"filters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	results, err := db.QueryDjangoModel(r.Context(), request.ModelName, request.Filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query Django model: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Failed to encode results: %v", err)
	}
}

func (db *DjangoBridge) handleCreateDjangoModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ModelName string                 `json:"model_name"`
		Data      map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	result, err := db.CreateDjangoModel(r.Context(), request.ModelName, request.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Django model: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Failed to encode result: %v", err)
	}
}

func (db *DjangoBridge) handleUpdateDjangoModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ModelName string                 `json:"model_name"`
		ID        int                    `json:"id"`
		Data      map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if err := db.UpdateDjangoModel(r.Context(), request.ModelName, request.ID, request.Data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update Django model: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (db *DjangoBridge) handleDeleteDjangoModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ModelName string `json:"model_name"`
		ID        int    `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if err := db.DeleteDjangoModel(r.Context(), request.ModelName, request.ID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete Django model: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (db *DjangoBridge) handleExecuteDjangoManagementCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	result, err := db.ExecuteDjangoManagementCommand(r.Context(), request.Command, request.Args)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute Django management command: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"result": result}); err != nil {
		log.Printf("Failed to encode result: %v", err)
	}
}

func (db *DjangoBridge) handleGetDjangoSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	settings, err := db.GetDjangoSettings(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get Django settings: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		log.Printf("Failed to encode settings: %v", err)
	}
}

func (db *DjangoBridge) handleCreateDjangoAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	agent, err := db.CreateDjangoAgent(r.Context(), request.ID, request.Name, request.Role)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Django agent: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(agent); err != nil {
		log.Printf("Failed to encode agent: %v", err)
	}
}
