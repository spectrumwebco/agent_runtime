package langchain

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

type MemoryType string

const (
	MemoryTypeBuffer MemoryType = "buffer"
	MemoryTypeBufferWindow MemoryType = "buffer_window"
	MemoryTypeSummary MemoryType = "summary"
)

type MemoryConfig struct {
	Type         MemoryType
	WindowSize   int
	LLM          *Client
	MemoryKey    string
	InputKey     string
	OutputKey    string
	HumanPrefix  string
	AIPrefix     string
	ReturnMessages bool
}

func NewMemory(config MemoryConfig) (memory.Memory, error) {
	if config.MemoryKey == "" {
		config.MemoryKey = "chat_history"
	}
	if config.InputKey == "" {
		config.InputKey = "input"
	}
	if config.OutputKey == "" {
		config.OutputKey = "output"
	}
	if config.HumanPrefix == "" {
		config.HumanPrefix = "Human"
	}
	if config.AIPrefix == "" {
		config.AIPrefix = "AI"
	}

	switch config.Type {
	case MemoryTypeBuffer:
		return memory.NewConversationBuffer(
			memory.WithMemoryKey(config.MemoryKey),
			memory.WithInputKey(config.InputKey),
			memory.WithOutputKey(config.OutputKey),
			memory.WithHumanPrefix(config.HumanPrefix),
			memory.WithAIPrefix(config.AIPrefix),
			memory.WithReturnMessages(config.ReturnMessages),
		), nil
	case MemoryTypeBufferWindow:
		if config.WindowSize <= 0 {
			config.WindowSize = 5
		}
		return memory.NewConversationBufferWindow(
			memory.WithHumanPrefix(config.HumanPrefix),
			memory.WithAIPrefix(config.AIPrefix),
			memory.WithMemoryKey(config.MemoryKey),
			memory.WithInputKey(config.InputKey),
			memory.WithOutputKey(config.OutputKey),
			memory.WithReturnMessages(config.ReturnMessages),
			memory.WithK(config.WindowSize),
		), nil
	case MemoryTypeSummary:
		if config.LLM == nil {
			return nil, fmt.Errorf("LLM is required for summary memory")
		}
		return memory.NewConversationSummary(
			config.LLM.LLM,
			memory.WithHumanPrefix(config.HumanPrefix),
			memory.WithAIPrefix(config.AIPrefix),
			memory.WithMemoryKey(config.MemoryKey),
			memory.WithInputKey(config.InputKey),
			memory.WithOutputKey(config.OutputKey),
			memory.WithReturnMessages(config.ReturnMessages),
		), nil
	default:
		return nil, fmt.Errorf("unsupported memory type: %s", config.Type)
	}
}

func SaveMemory(ctx context.Context, mem memory.Memory, input, output string) error {
	inputMap := map[string]any{
		"input":  input,
		"output": output,
	}

	return mem.SaveContext(ctx, inputMap, nil)
}

func LoadMemory(ctx context.Context, mem memory.Memory) ([]schema.ChatMessage, error) {
	inputMap := map[string]any{
		"input": "",
	}

	outputMap, err := mem.LoadMemoryVariables(ctx, inputMap)
	if err != nil {
		return nil, fmt.Errorf("failed to load memory: %w", err)
	}

	chatHistory, ok := outputMap["chat_history"].([]schema.ChatMessage)
	if !ok {
		return nil, fmt.Errorf("invalid memory format")
	}

	return chatHistory, nil
}
