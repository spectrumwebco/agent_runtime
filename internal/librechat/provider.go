package librechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	
	"github.com/spectrumwebco/agent_runtime/pkg/librechat"
)

type MCPProvider struct {
	client librechat.LibreChatClient
}

func NewMCPProvider(baseURL string) *MCPProvider {
	apiKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
	return &MCPProvider{
		client: librechat.NewClient(baseURL, apiKey),
	}
}

func (p *MCPProvider) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/mcp/librechat/execute", p.handleExecute)
	mux.HandleFunc("/mcp/librechat/execute-multi", p.handleExecuteMulti)
	mux.HandleFunc("/mcp/librechat/languages", p.handleGetLanguages)
	mux.HandleFunc("/mcp/librechat/docs", p.handleGetDocs)
}

func (p *MCPProvider) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var params struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}
	
	output, err := p.client.ExecuteCode(r.Context(), params.Language, params.Code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Code execution failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"output":  output,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *MCPProvider) handleExecuteMulti(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var params struct {
		Requests []librechat.CodeExecutionRequest `json:"requests"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}
	
	responses, err := p.client.ExecuteMultiLanguageCode(r.Context(), params.Requests)
	if err != nil {
		http.Error(w, fmt.Sprintf("Multi-language code execution failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":   true,
		"responses": responses,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *MCPProvider) handleGetLanguages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	languages, err := p.client.GetSupportedLanguages(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Get supported languages failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":   true,
		"languages": languages,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *MCPProvider) handleGetDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	docs := generateAPIDocumentation()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func generateAPIDocumentation() map[string]interface{} {
	return map[string]interface{}{
		"name":        "LibreChat Code Interpreter API",
		"description": "API for executing code in multiple programming languages",
		"version":     "1.0.0",
		"auth": map[string]interface{}{
			"type":        "apiKey",
			"name":        "x-api-key",
			"in":          "header",
			"description": "API key for authentication",
		},
		"endpoints": []map[string]interface{}{
			{
				"path":        "/mcp/librechat/execute",
				"method":      "POST",
				"description": "Execute code in a single programming language",
				"parameters": map[string]interface{}{
					"language": "Programming language to execute code in",
					"code":     "Code to execute",
				},
				"responses": map[string]interface{}{
					"success": "Whether the execution was successful",
					"output":  "Output of the code execution",
					"error":   "Error message if execution failed",
				},
				"examples": []map[string]interface{}{
					{
						"language": "python",
						"request": map[string]interface{}{
							"language": "python",
							"code":     "print('Hello, world!')",
						},
						"response": map[string]interface{}{
							"success": true,
							"output":  "Hello, world!",
						},
					},
					{
						"language": "javascript",
						"request": map[string]interface{}{
							"language": "javascript",
							"code":     "console.log('Hello, world!');",
						},
						"response": map[string]interface{}{
							"success": true,
							"output":  "Hello, world!",
						},
					},
				},
			},
			{
				"path":        "/mcp/librechat/execute-multi",
				"method":      "POST",
				"description": "Execute code in multiple programming languages",
				"parameters": map[string]interface{}{
					"requests": "Array of code execution requests",
				},
				"responses": map[string]interface{}{
					"success":   "Whether the execution was successful",
					"responses": "Array of code execution responses",
					"error":     "Error message if execution failed",
				},
				"examples": []map[string]interface{}{
					{
						"request": map[string]interface{}{
							"requests": []map[string]interface{}{
								{
									"id":       "python-1",
									"language": "python",
									"code":     "print('Hello from Python!')",
								},
								{
									"id":       "js-1",
									"language": "javascript",
									"code":     "console.log('Hello from JavaScript!');",
								},
							},
						},
						"response": map[string]interface{}{
							"success": true,
							"responses": []map[string]interface{}{
								{
									"id":      "python-1",
									"success": true,
									"output":  "Hello from Python!",
								},
								{
									"id":      "js-1",
									"success": true,
									"output":  "Hello from JavaScript!",
								},
							},
						},
					},
				},
			},
			{
				"path":        "/mcp/librechat/languages",
				"method":      "GET",
				"description": "Get the list of supported programming languages",
				"parameters":  map[string]interface{}{},
				"responses": map[string]interface{}{
					"success":   "Whether the request was successful",
					"languages": "Array of supported programming languages",
					"error":     "Error message if request failed",
				},
				"examples": []map[string]interface{}{
					{
						"response": map[string]interface{}{
							"success": true,
							"languages": []string{
								"python",
								"javascript",
								"typescript",
								"go",
								"rust",
								"c",
								"cpp",
								"csharp",
								"php",
							},
						},
					},
				},
			},
		},
		"supported_languages": []string{
			"python",
			"javascript",
			"typescript",
			"go",
			"rust",
			"c",
			"cpp",
			"csharp",
			"php",
		},
		"integration_examples": map[string]interface{}{
			"python": `
import requests
import json

def execute_code(api_key, language, code):
    url = "http://localhost:8080/mcp/librechat/execute"
    headers = {
        "Content-Type": "application/json",
        "x-api-key": api_key
    }
    payload = {
        "language": language,
        "code": code
    }
    response = requests.post(url, headers=headers, json=payload)
    return response.json()

# Example usage
api_key = "your-api-key"
result = execute_code(api_key, "python", "print('Hello, world!')")
print(result)
`,
			"javascript": `
async function executeCode(apiKey, language, code) {
  const url = "http://localhost:8080/mcp/librechat/execute";
  const headers = {
    "Content-Type": "application/json",
    "x-api-key": apiKey
  };
  const payload = {
    language,
    code
  };
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(payload)
  });
  return await response.json();
}

const apiKey = "your-api-key";
executeCode(apiKey, "javascript", "console.log('Hello, world!');")
  .then(result => console.log(result));
`,
			"go": `
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func executeCode(apiKey, language, code string) (map[string]interface{}, error) {
	url := "http://localhost:8080/mcp/librechat/execute"
	payload := map[string]interface{}{
		"language": language,
		"code":     code,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func main() {
	apiKey := "your-api-key"
	result, err := executeCode(apiKey, "go", "fmt.Println(\"Hello, world!\")")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Result: %v\n", result)
}
`,
		},
		"error_handling": map[string]interface{}{
			"description": "Error handling guidance for the LibreChat Code Interpreter API",
			"common_errors": []map[string]interface{}{
				{
					"status_code": 400,
					"description": "Bad Request - The request was invalid or cannot be served",
					"resolution":  "Check the request parameters and ensure they are valid",
				},
				{
					"status_code": 401,
					"description": "Unauthorized - Authentication failed",
					"resolution":  "Check the API key and ensure it is valid",
				},
				{
					"status_code": 500,
					"description": "Internal Server Error - The server encountered an error",
					"resolution":  "Check the error message and try again later",
				},
			},
			"best_practices": []string{
				"Always check the 'success' field in the response",
				"Handle errors gracefully and provide meaningful error messages to users",
				"Implement retry logic for transient errors",
				"Use timeouts to prevent hanging requests",
			},
		},
	}
}
