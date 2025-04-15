package langchain

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/internal/langchain"
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Module struct {
	modules.BaseModule
	client *langchain.Client
	agent  *langchain.Agent
	mutex  sync.RWMutex
}

func NewModule() *Module {
	return &Module{
		BaseModule: *modules.NewBaseModule("langchain", "LangChain integration for AI capabilities"),
	}
}

func (m *Module) Initialize(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.BaseModule.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize base module: %w", err)
	}

	client, err := langchain.NewClient(langchain.ClientConfig{
		Type:        langchain.ModelTypeOpenAI,
		ModelName:   "gpt-4o",
		Temperature: 0.7,
	})
	if err != nil {
		return fmt.Errorf("failed to create LangChain client: %w", err)
	}
	m.client = client

	mem, err := langchain.NewMemory(langchain.MemoryConfig{
		Type:        langchain.MemoryTypeBuffer,
		HumanPrefix: "User",
		AIPrefix:    "Kled",
	})
	if err != nil {
		return fmt.Errorf("failed to create memory: %w", err)
	}

	m.registerTools()

	agent, err := langchain.NewAgent(langchain.AgentConfig{
		Client:        client,
		Memory:        mem,
		SystemPrompt:  "You are Kled, a Senior Software Engineering Lead & Technical Authority for AI/ML.",
		MaxIterations: 5,
	})
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	m.agent = agent

	return nil
}

func (m *Module) Cleanup(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.BaseModule.Cleanup(ctx); err != nil {
		return fmt.Errorf("failed to clean up base module: %w", err)
	}

	if m.agent != nil {
		if err := m.agent.ClearMemory(); err != nil {
			return fmt.Errorf("failed to clear agent memory: %w", err)
		}
	}

	return nil
}

func (m *Module) registerTools() {
	generateTool := tools.NewTool(
		"langchain.generate",
		"Generate text using LangChain",
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			m.mutex.RLock()
			defer m.mutex.RUnlock()

			prompt, ok := args["prompt"].(string)
			if !ok {
				return nil, fmt.Errorf("prompt must be a string")
			}

			response, err := m.client.Generate(ctx, prompt)
			if err != nil {
				return nil, fmt.Errorf("failed to generate text: %w", err)
			}

			return response, nil
		},
	)

	agentTool := tools.NewTool(
		"langchain.agent",
		"Run the LangChain agent",
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			m.mutex.RLock()
			defer m.mutex.RUnlock()

			input, ok := args["input"].(string)
			if !ok {
				return nil, fmt.Errorf("input must be a string")
			}

			response, err := m.agent.Run(ctx, input)
			if err != nil {
				return nil, fmt.Errorf("failed to run agent: %w", err)
			}

			return response, nil
		},
	)

	m.RegisterTool(generateTool)
	m.RegisterTool(agentTool)
}

func (m *Module) GetClient() *langchain.Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.client
}

func (m *Module) GetAgent() *langchain.Agent {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.agent
}

func init() {
	modules.RegisterModule("langchain", func() modules.Module {
		return NewModule()
	})
}
