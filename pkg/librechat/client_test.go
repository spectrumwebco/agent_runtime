package librechat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ExecuteCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/code/execute", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
		
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		assert.NoError(t, err)
		
		assert.Equal(t, "python", requestBody["language"])
		assert.Equal(t, "print('Hello, World!')", requestBody["code"])
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "output": "Hello, World!"}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-api-key")
	
	output, err := client.ExecuteCode(context.Background(), "python", "print('Hello, World!')")
	
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", output)
}

func TestClient_ExecuteCode_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": false, "error": "Execution failed"}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-api-key")
	
	_, err := client.ExecuteCode(context.Background(), "python", "print('Hello, World!')")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Execution failed")
}

func TestClient_ExecuteMultiLanguageCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/code/execute-multi", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
		
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		assert.NoError(t, err)
		
		requests, ok := requestBody["requests"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, requests, 2)
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"success": true,
			"responses": [
				{"id": "1", "success": true, "output": "Hello from Python!"},
				{"id": "2", "success": true, "output": "Hello from JavaScript!"}
			]
		}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-api-key")
	
	requests := []CodeExecutionRequest{
		{
			ID:       "1",
			Language: "python",
			Code:     "print('Hello from Python!')",
		},
		{
			ID:       "2",
			Language: "javascript",
			Code:     "console.log('Hello from JavaScript!')",
		},
	}
	
	responses, err := client.ExecuteMultiLanguageCode(context.Background(), requests)
	
	assert.NoError(t, err)
	assert.Len(t, responses, 2)
	assert.Equal(t, "1", responses[0].ID)
	assert.Equal(t, "Hello from Python!", responses[0].Output)
	assert.Equal(t, "2", responses[1].ID)
	assert.Equal(t, "Hello from JavaScript!", responses[1].Output)
}

func TestClient_ExecuteMultiLanguageCode_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": false, "error": "Execution failed"}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-api-key")
	
	requests := []CodeExecutionRequest{
		{
			ID:       "1",
			Language: "python",
			Code:     "print('Hello from Python!')",
		},
	}
	
	_, err := client.ExecuteMultiLanguageCode(context.Background(), requests)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Execution failed")
}

func TestClient_GetSupportedLanguages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/code/languages", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"success": true,
			"languages": ["python", "javascript", "go", "rust", "c++", "c#", "php"]
		}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-api-key")
	
	languages, err := client.GetSupportedLanguages(context.Background())
	
	assert.NoError(t, err)
	assert.Len(t, languages, 7)
	assert.Contains(t, languages, "python")
	assert.Contains(t, languages, "javascript")
	assert.Contains(t, languages, "go")
	assert.Contains(t, languages, "rust")
	assert.Contains(t, languages, "c++")
	assert.Contains(t, languages, "c#")
	assert.Contains(t, languages, "php")
}

func TestClient_GetSupportedLanguages_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": false, "error": "API error"}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL, "test-api-key")
	
	_, err := client.GetSupportedLanguages(context.Background())
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
}
