package librechat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spectrumwebco/agent_runtime/pkg/librechat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLibreChatClient struct {
	mock.Mock
}

func (m *mockLibreChatClient) ExecuteCode(ctx context.Context, language, code string) (string, error) {
	args := m.Called(ctx, language, code)
	return args.String(0), args.Error(1)
}

func (m *mockLibreChatClient) ExecuteMultiLanguageCode(ctx context.Context, requests []librechat.CodeExecutionRequest) ([]librechat.CodeExecutionResponse, error) {
	args := m.Called(ctx, requests)
	return args.Get(0).([]librechat.CodeExecutionResponse), args.Error(1)
}

func (m *mockLibreChatClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func TestMCPProvider_handleExecute(t *testing.T) {
	mockClient := new(mockLibreChatClient)
	mockClient.On("ExecuteCode", mock.Anything, "python", "print('Hello, world!')").Return("Hello, world!", nil)

	provider := &MCPProvider{
		client: mockClient,
	}

	req, err := http.NewRequest("POST", "/mcp/librechat/execute", strings.NewReader(`{"language":"python","code":"print('Hello, world!')"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(provider.handleExecute)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Hello, world!", response["output"])
}

func TestMCPProvider_handleExecuteMulti(t *testing.T) {
	mockClient := new(mockLibreChatClient)
	expectedResponses := []librechat.CodeExecutionResponse{
		{
			ID:      "python-1",
			Success: true,
			Output:  "Hello from Python!",
		},
		{
			ID:      "js-1",
			Success: true,
			Output:  "Hello from JavaScript!",
		},
	}
	mockClient.On("ExecuteMultiLanguageCode", mock.Anything, mock.Anything).Return(expectedResponses, nil)

	provider := &MCPProvider{
		client: mockClient,
	}

	req, err := http.NewRequest("POST", "/mcp/librechat/execute-multi", strings.NewReader(`{
		"requests": [
			{
				"id": "python-1",
				"language": "python",
				"code": "print('Hello from Python!')"
			},
			{
				"id": "js-1",
				"language": "javascript",
				"code": "console.log('Hello from JavaScript!')"
			}
		]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(provider.handleExecuteMulti)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	responses, ok := response["responses"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, responses, 2)
}

func TestMCPProvider_handleGetLanguages(t *testing.T) {
	mockClient := new(mockLibreChatClient)
	expectedLanguages := []string{"python", "javascript", "go"}
	mockClient.On("GetSupportedLanguages", mock.Anything).Return(expectedLanguages, nil)

	provider := &MCPProvider{
		client: mockClient,
	}

	req, err := http.NewRequest("GET", "/mcp/librechat/languages", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(provider.handleGetLanguages)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	languages, ok := response["languages"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, languages, 3)
	assert.Equal(t, "python", languages[0])
	assert.Equal(t, "javascript", languages[1])
	assert.Equal(t, "go", languages[2])
}

func TestMCPProvider_handleGetDocs(t *testing.T) {
	provider := &MCPProvider{}

	req, err := http.NewRequest("GET", "/mcp/librechat/docs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(provider.handleGetDocs)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["name"])
	assert.NotNil(t, response["description"])
	assert.NotNil(t, response["version"])
}
