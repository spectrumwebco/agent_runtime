package langchain

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

type Agent struct {
	LLM    *Client
	Memory memory.Memory
	Agent  agents.Agent
	Tools  []tools.Tool
}

type AgentConfig struct {
	Client       *Client
	Memory       memory.Memory
	Tools        []tools.Tool
	SystemPrompt string
	MaxIterations int
}

func NewAgent(config AgentConfig) (*Agent, error) {
	if config.MaxIterations == 0 {
		config.MaxIterations = 10
	}

	opts := []agents.Option{
		agents.WithMaxIterations(config.MaxIterations),
	}

	if config.Memory != nil {
		opts = append(opts, agents.WithMemory(config.Memory))
	}

	if config.SystemPrompt != "" {
		opts = append(opts, agents.WithSystemMessage(config.SystemPrompt))
	}

	agent, err := agents.Initialize(
		config.Client.LLM,
		config.Tools,
		agents.ChatConversational,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return &Agent{
		LLM:    config.Client,
		Memory: config.Memory,
		Agent:  agent,
		Tools:  config.Tools,
	}, nil
}

func (a *Agent) Run(ctx context.Context, input string) (string, error) {
	inputMap := map[string]any{
		"input": input,
	}

	output, err := agents.Run(ctx, a.Agent, inputMap)
	if err != nil {
		return "", fmt.Errorf("failed to run agent: %w", err)
	}

	return output, nil
}

func (a *Agent) GetTools() []tools.Tool {
	return a.Tools
}

func (a *Agent) AddTool(tool tools.Tool) error {
	for _, t := range a.Tools {
		if t.Name() == tool.Name() {
			return fmt.Errorf("tool %s already exists", tool.Name())
		}
	}

	a.Tools = append(a.Tools, tool)

	agent, err := agents.Initialize(
		a.LLM.LLM,
		a.Tools,
		agents.ChatConversational,
		agents.WithMemory(a.Memory),
	)
	if err != nil {
		return fmt.Errorf("failed to reinitialize agent: %w", err)
	}

	a.Agent = agent
	return nil
}

func (a *Agent) RemoveTool(name string) error {
	var tools []tools.Tool
	found := false
	for _, tool := range a.Tools {
		if tool.Name() != name {
			tools = append(tools, tool)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("tool %s not found", name)
	}

	a.Tools = tools

	agent, err := agents.Initialize(
		a.LLM.LLM,
		a.Tools,
		agents.ChatConversational,
		agents.WithMemory(a.Memory),
	)
	if err != nil {
		return fmt.Errorf("failed to reinitialize agent: %w", err)
	}

	a.Agent = agent
	return nil
}

func (a *Agent) ClearMemory() error {
	if a.Memory == nil {
		return nil
	}
	return a.Memory.Clear()
}
