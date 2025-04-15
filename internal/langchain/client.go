package langchain

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

type ModelType string

const (
	ModelTypeOpenAI ModelType = "openai"
	ModelTypeOllama ModelType = "ollama"
)

type ClientConfig struct {
	Type       ModelType
	ModelName  string
	APIKey     string
	BaseURL    string
	MaxTokens  int
	Temperature float64
}

type Client struct {
	Type      ModelType
	ModelName string
	LLM       llms.Model
}

func NewClient(config ClientConfig) (*Client, error) {
	var model llms.Model
	var err error

	switch config.Type {
	case ModelTypeOpenAI:
		if config.APIKey == "" {
			config.APIKey = os.Getenv("OPENAI_API_KEY")
			if config.APIKey == "" {
				return nil, fmt.Errorf("OpenAI API key is required")
			}
		}

		if config.ModelName == "" {
			config.ModelName = openai.GPT4o
		}

		opts := []openai.Option{
			openai.WithModel(config.ModelName),
		}

		if config.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(config.BaseURL))
		}

		if config.MaxTokens > 0 {
			opts = append(opts, openai.WithMaxTokens(config.MaxTokens))
		}

		if config.Temperature > 0 {
			opts = append(opts, openai.WithTemperature(float32(config.Temperature)))
		}

		model, err = openai.New(
			openai.WithToken(config.APIKey),
			openai.WithModel(config.ModelName),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
		}

	case ModelTypeOllama:
		if config.ModelName == "" {
			config.ModelName = "llama3"
		}

		opts := []ollama.Option{
			ollama.WithModel(config.ModelName),
		}

		if config.BaseURL != "" {
			opts = append(opts, ollama.WithServerURL(config.BaseURL))
		}

		model, err = ollama.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create Ollama client: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported model type: %s", config.Type)
	}

	return &Client{
		Type:      config.Type,
		ModelName: config.ModelName,
		LLM:       model,
	}, nil
}

func (c *Client) Generate(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	completion, err := c.LLM.Call(ctx, prompt, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}
	return completion, nil
}

func (c *Client) GenerateWithTools(ctx context.Context, prompt string, tools []schema.FunctionDefinition, opts ...llms.CallOption) (*schema.ModelResponse, error) {
	options := append(opts, llms.WithFunctions(tools...))

	response, err := c.LLM.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeHuman, prompt),
	}, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text with tools: %w", err)
	}

	return response, nil
}

func (c *Client) ExecuteWithTools(ctx context.Context, prompt string, tools []schema.FunctionDefinition, toolMap map[string]schema.FunctionCallable, opts ...llms.CallOption) (string, error) {
	response, err := c.GenerateWithTools(ctx, prompt, tools, opts...)
	if err != nil {
		return "", err
	}

	if len(response.ToolCalls()) == 0 {
		return response.Content(), nil
	}

	var toolResults []llms.MessageContent
	for _, toolCall := range response.ToolCalls() {
		tool, ok := toolMap[toolCall.Name()]
		if !ok {
			return "", fmt.Errorf("tool not found: %s", toolCall.Name())
		}

		result, err := tool(ctx, toolCall.Arguments())
		if err != nil {
			return "", fmt.Errorf("failed to execute tool %s: %w", toolCall.Name(), err)
		}

		toolResults = append(toolResults, llms.ToolResult(toolCall.ID(), result))
	}

	finalResponse, err := c.LLM.GenerateContent(ctx, append(
		[]llms.MessageContent{
			llms.TextParts(schema.ChatMessageTypeHuman, prompt),
			llms.AIMessage(response.Content(), response.ToolCalls()...),
		},
		toolResults...,
	), opts...)
	if err != nil {
		return "", fmt.Errorf("failed to generate final response: %w", err)
	}

	return finalResponse.Content(), nil
}
