package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type SupabaseAdapter struct {
	apiURL      string
	apiKey      string
	httpClient  *http.Client
	initialized bool
}

type SupabaseConfig struct {
	APIURL     string
	APIKey     string
	HTTPClient *http.Client
	Timeout    time.Duration
}

func NewSupabaseAdapter(config SupabaseConfig) *SupabaseAdapter {
	if config.HTTPClient == nil {
		timeout := config.Timeout
		if timeout == 0 {
			timeout = time.Second * 10
		}
		config.HTTPClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &SupabaseAdapter{
		apiURL:     config.APIURL,
		apiKey:     config.APIKey,
		httpClient: config.HTTPClient,
	}
}

func (a *SupabaseAdapter) Initialize() error {
	if a.initialized {
		return nil
	}

	if a.apiURL == "" {
		return fmt.Errorf("Supabase API URL is required")
	}
	if a.apiKey == "" {
		return fmt.Errorf("Supabase API key is required")
	}

	if !strings.HasSuffix(a.apiURL, "/") {
		a.apiURL = a.apiURL + "/"
	}

	a.initialized = true
	return nil
}

func (a *SupabaseAdapter) GetState(stateType models.StateType, stateID string) (*models.State, error) {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%srest/v1/states?id=eq.%s&type=eq.%s", a.apiURL, stateID, stateType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", a.apiKey)
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("state not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var states []*models.State
	if err := json.NewDecoder(resp.Body).Decode(&states); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(states) == 0 {
		return nil, fmt.Errorf("state not found")
	}

	return states[0], nil
}

func (a *SupabaseAdapter) UpdateState(state *models.State) error {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%srest/v1/states?id=eq.%s", a.apiURL, state.ID)
	
	state.Version++
	state.UpdatedAt = time.Now().UTC()
	
	body, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", a.apiKey)
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *SupabaseAdapter) CreateState(state *models.State) error {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%srest/v1/states", a.apiURL)
	
	body, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", a.apiKey)
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *SupabaseAdapter) DeleteState(stateType models.StateType, stateID string) error {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%srest/v1/states?id=eq.%s&type=eq.%s", a.apiURL, stateID, stateType)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", a.apiKey)
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *SupabaseAdapter) ListStates(stateType models.StateType) ([]*models.State, error) {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%srest/v1/states?type=eq.%s", a.apiURL, stateType)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", a.apiKey)
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var states []*models.State
	if err := json.NewDecoder(resp.Body).Decode(&states); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return states, nil
}

func (a *SupabaseAdapter) Close() error {
	a.initialized = false
	return nil
}
