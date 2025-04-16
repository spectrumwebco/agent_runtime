package langraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type ModelProvider string

const (
	ModelProviderKlusterAI ModelProvider = "kluster_ai"
	ModelProviderOpenAI    ModelProvider = "openai"
	ModelProviderGemini    ModelProvider = "gemini"
	ModelProviderLlama     ModelProvider = "llama"
	ModelProviderAnthropic ModelProvider = "anthropic"
)

type ModelRole string

const (
	ModelRoleCoding    ModelRole = "coding"     // For code generation tasks (Gemini 2.5 Pro)
	ModelRoleReasoning ModelRole = "reasoning"  // For reasoning tasks (Llama 4)
	ModelRoleGeneral   ModelRole = "general"    // For general purpose tasks
)

type ModelConfig struct {
	Provider    ModelProvider          `json:"provider"`
	ModelName   string                 `json:"model_name"`
	Role        ModelRole              `json:"role"`
	APIKey      string                 `json:"api_key,omitempty"`
	BaseURL     string                 `json:"base_url,omitempty"`
	Temperature float64                `json:"temperature"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

type ModelRegistry struct {
	models     map[string]llms.LLM
	configs    map[string]ModelConfig
	defaultLLM string
	mutex      sync.RWMutex
}

func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models:  make(map[string]llms.LLM),
		configs: make(map[string]ModelConfig),
	}
}

func (r *ModelRegistry) RegisterModel(name string, llm llms.LLM, config ModelConfig) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.models[name] = llm
	r.configs[name] = config
	
	if r.defaultLLM == "" || config.Role == ModelRoleCoding {
		r.defaultLLM = name
	}
}

func (r *ModelRegistry) GetModel(name string) (llms.LLM, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	llm, exists := r.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s not found", name)
	}
	
	return llm, nil
}

func (r *ModelRegistry) GetModelByRole(role ModelRole) (llms.LLM, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for name, config := range r.configs {
		if config.Role == role {
			return r.models[name], nil
		}
	}
	
	return nil, fmt.Errorf("no model found for role %s", role)
}

func (r *ModelRegistry) GetDefaultModel() (llms.LLM, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.defaultLLM == "" {
		return nil, fmt.Errorf("no default model set")
	}
	
	return r.models[r.defaultLLM], nil
}

func (r *ModelRegistry) SetDefaultModel(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.models[name]; !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	r.defaultLLM = name
	return nil
}

func (r *ModelRegistry) ListModels() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var models []string
	for name := range r.models {
		models = append(models, name)
	}
	
	return models
}

func (r *ModelRegistry) GetModelConfig(name string) (ModelConfig, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	config, exists := r.configs[name]
	if !exists {
		return ModelConfig{}, fmt.Errorf("model %s not found", name)
	}
	
	return config, nil
}

type KlusterAILLM struct {
	Config     ModelConfig
	HTTPClient *http.Client
}

func NewKlusterAILLM(config ModelConfig) (*KlusterAILLM, error) {
	if config.Provider != ModelProviderKlusterAI {
		return nil, fmt.Errorf("invalid provider for Kluster.AI LLM: %s", config.Provider)
	}
	
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for Kluster.AI LLM")
	}
	
	if config.BaseURL == "" {
		config.BaseURL = "https://api.kluster.ai/v1"
	}
	
	return &KlusterAILLM{
		Config: config,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (k *KlusterAILLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	callOptions := llms.CallOptions{}
	for _, opt := range options {
		opt(&callOptions)
	}
	
	requestBody := map[string]interface{}{
		"model": k.Config.ModelName,
		"prompt": prompt,
		"temperature": k.Config.Temperature,
	}
	
	if k.Config.MaxTokens > 0 {
		requestBody["max_tokens"] = k.Config.MaxTokens
	}
	
	for key, value := range k.Config.Options {
		requestBody[key] = value
	}
	
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/completions", k.Config.BaseURL), bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.Config.APIKey))
	
	resp, err := k.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return "", fmt.Errorf("request failed with status code %d", resp.StatusCode)
		}
		return "", fmt.Errorf("request failed: %v", errorResponse)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format")
	}
	
	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}
	
	text, ok := choice["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text format")
	}
	
	return text, nil
}

func (k *KlusterAILLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) (*llms.Generation, error) {
	var generations []schema.Generation
	
	for _, prompt := range prompts {
		text, err := k.Call(ctx, prompt, options...)
		if err != nil {
			return nil, err
		}
		
		generations = append(generations, schema.Generation{
			Text: text,
		})
	}
	
	return &llms.Generation{
		Generations: generations,
	}, nil
}

func CreateStandardModelRegistry(klusterAIAPIKey string) (*ModelRegistry, error) {
	registry := NewModelRegistry()
	
	llama4Config := ModelConfig{
		Provider:    ModelProviderKlusterAI,
		ModelName:   "llama-4",
		Role:        ModelRoleReasoning,
		APIKey:      klusterAIAPIKey,
		Temperature: 0.7,
	}
	
	llama4, err := NewKlusterAILLM(llama4Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Llama 4 model: %w", err)
	}
	
	gemini25Config := ModelConfig{
		Provider:    ModelProviderKlusterAI,
		ModelName:   "gemini-2.5-pro",
		Role:        ModelRoleCoding,
		APIKey:      klusterAIAPIKey,
		Temperature: 0.2,
	}
	
	gemini25, err := NewKlusterAILLM(gemini25Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini 2.5 Pro model: %w", err)
	}
	
	registry.RegisterModel("llama-4", llama4, llama4Config)
	registry.RegisterModel("gemini-2.5-pro", gemini25, gemini25Config)
	
	return registry, nil
}

func UpdateLangChainBridgeWithModelRegistry(bridge *LangChainBridge, registry *ModelRegistry) error {
	defaultModel, err := registry.GetDefaultModel()
	if err != nil {
		return err
	}
	
	bridge.LLM = defaultModel
	
	bridge.LCAgentFactory = NewLangChainAgentFactory(defaultModel, bridge.EventStream)
	
	return nil
}

func CreateMultiAgentSystemWithModelRegistry(name, description string, registry *ModelRegistry, eventStream EventStream) (*MultiAgentSystem, *LangChainBridge, error) {
	defaultModel, err := registry.GetDefaultModel()
	if err != nil {
		return nil, nil, err
	}
	
	system, bridge, err := CreateMultiAgentSystemWithLangChain(name, description, defaultModel, eventStream)
	if err != nil {
		return nil, nil, err
	}
	
	system.Metadata["model_registry"] = registry
	
	return system, bridge, nil
}
