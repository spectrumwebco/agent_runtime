package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spectrumwebco/agent_runtime/internal/agent"
	"github.com/spectrumwebco/agent_runtime/pkg/prompts"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools/retriever"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/memory"
)

type CorkiAgent struct {
	ID            string
	Name          string
	Description   string
	AgentInstance *agent.Agent
	ToolRegistry  *tools.Registry
	Config        *CorkiConfig
}

type CorkiConfig struct {
	PromptPath    string
	ToolsPath     string
	ModulesPath   string
	LLMProvider   string
	LLMModel      string
	MemoryType    string
	MaxTokens     int
	Temperature   float64
	EnableLogging bool
}

func NewCorkiAgent(config *CorkiConfig) (*CorkiAgent, error) {
	if config == nil {
		config = &CorkiConfig{
			PromptPath:    "temp_agents/corki_agent/prompt/prompts.txt",
			ToolsPath:     "temp_agents/corki_agent/tools.json",
			ModulesPath:   "temp_agents/corki_agent",
			LLMProvider:   "openai",
			LLMModel:      "gpt-4",
			MemoryType:    "buffer",
			MaxTokens:     4096,
			Temperature:   0.7,
			EnableLogging: true,
		}
	}

	agentID := uuid.New().String()

	corki := &CorkiAgent{
		ID:          agentID,
		Name:        "Corki",
		Description: "Backup Agent for 100% Code Coverage",
		Config:      config,
	}

	toolRegistry, err := loadToolRegistry(config.ToolsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load tool registry: %w", err)
	}
	corki.ToolRegistry = toolRegistry

	agentInstance, err := initializeAgent(corki, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize agent: %w", err)
	}
	corki.AgentInstance = agentInstance

	return corki, nil
}

func loadToolRegistry(toolsPath string) (*tools.Registry, error) {
	toolsFile, err := os.ReadFile(toolsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tools file: %w", err)
	}

	var toolsConfig struct {
		Tools []tools.ToolConfig `json:"tools"`
	}
	if err := json.Unmarshal(toolsFile, &toolsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse tools JSON: %w", err)
	}

	registry := tools.NewRegistry()

	for _, toolConfig := range toolsConfig.Tools {
		tool, err := createToolFromConfig(toolConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create tool %s: %w", toolConfig.Name, err)
		}

		if err := registry.RegisterTool(tool); err != nil {
			return nil, fmt.Errorf("failed to register tool %s: %w", toolConfig.Name, err)
		}
	}

	return registry, nil
}

func createToolFromConfig(config tools.ToolConfig) (tools.Tool, error) {
	
	return &tools.GenericTool{
		ToolName:        config.Name,
		ToolDescription: config.Description,
		ToolFunction: func(ctx context.Context, args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s executed with args: %v", config.Name, args), nil
		},
	}, nil
}

func initializeAgent(corki *CorkiAgent, config *CorkiConfig) (*agent.Agent, error) {
	promptTemplate, err := loadPromptTemplate(config.PromptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt template: %w", err)
	}

	llm, err := createLLM(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}

	mem, err := createMemory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	opts := []agent.Option{
		agent.WithPromptTemplate(promptTemplate),
		agent.WithLLM(llm),
		agent.WithMemory(mem),
		agent.WithToolRegistry(corki.ToolRegistry),
		agent.WithMaxTokens(config.MaxTokens),
		agent.WithTemperature(config.Temperature),
	}

	agentInstance, err := agent.NewAgent(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return agentInstance, nil
}

func loadPromptTemplate(promptPath string) (string, error) {
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %w", err)
	}

	return string(promptBytes), nil
}

func createLLM(config *CorkiConfig) (llms.LLM, error) {
	
	return &mockLLM{}, nil
}

func createMemory(config *CorkiConfig) (schema.Memory, error) {
	
	return memory.NewBuffer(), nil
}

type mockLLM struct{}

func (m *mockLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return fmt.Sprintf("Mock response to: %s", prompt), nil
}

func (m *mockLLM) GenerateContent(ctx context.Context, messages []schema.ChatMessage, options ...llms.CallOption) (*schema.ContentResponse, error) {
	return &schema.ContentResponse{
		Choices: []*schema.ContentChoice{
			{
				Content: fmt.Sprintf("Mock response to: %v", messages),
			},
		},
	}, nil
}

func (c *CorkiAgent) Run(ctx context.Context, input string) (string, error) {
	if c.AgentInstance == nil {
		return "", fmt.Errorf("agent not initialized")
	}

	response, err := c.AgentInstance.Run(ctx, input)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	return response, nil
}

func (c *CorkiAgent) RegisterPythonModules() error {
	moduleFiles, err := filepath.Glob(filepath.Join(c.Config.ModulesPath, "*.py"))
	if err != nil {
		return fmt.Errorf("failed to list Python modules: %w", err)
	}

	for _, moduleFile := range moduleFiles {
		moduleName := filepath.Base(moduleFile)
		
		if moduleName == "__init__.py" || moduleName == "__pycache__" {
			continue
		}

		if err := c.registerPythonModule(moduleName); err != nil {
			return fmt.Errorf("failed to register Python module %s: %w", moduleName, err)
		}
	}

	return nil
}

func (c *CorkiAgent) registerPythonModule(moduleName string) error {
	
	log.Printf("Registered Python module: %s", moduleName)
	return nil
}

func (c *CorkiAgent) Shutdown() error {
	return nil
}
