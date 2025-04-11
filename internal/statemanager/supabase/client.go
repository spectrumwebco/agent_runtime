package supabase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type SupabaseClient struct {
	URL       string
	APIKey    string
	AuthToken string
	client    *http.Client
}

type SupabaseConfig struct {
	URL       string
	APIKey    string
	AuthToken string
}

func NewSupabaseClient(cfg SupabaseConfig) *SupabaseClient {
	return &SupabaseClient{
		URL:       cfg.URL,
		APIKey:    cfg.APIKey,
		AuthToken: cfg.AuthToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (sc *SupabaseClient) UpsertState(table string, stateData map[string]interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s", sc.URL, table)
	
	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("failed to marshal state data: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("apikey", sc.APIKey)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.AuthToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Prefer", "resolution=merge-duplicates")
	
	resp, err := sc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase error: status code %d", resp.StatusCode)
	}
	
	log.Printf("Successfully upserted state to %s table\n", table)
	return nil
}

func (sc *SupabaseClient) GetState(table string, id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", sc.URL, table, id)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("apikey", sc.APIKey)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.AuthToken))
	
	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase error: status code %d", resp.StatusCode)
	}
	
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("state not found")
	}
	
	return result[0], nil
}

func (sc *SupabaseClient) DeleteState(table string, id string) error {
	url := fmt.Sprintf("%s/rest/v1/%s?id=eq.%s", sc.URL, table, id)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("apikey", sc.APIKey)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.AuthToken))
	
	resp, err := sc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase error: status code %d", resp.StatusCode)
	}
	
	log.Printf("Successfully deleted state from %s table\n", table)
	return nil
}
