package langchain

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func Example() {
	client, err := NewClient(ClientConfig{
		Type:        ModelTypeOpenAI,
		ModelName:   "gpt-4o",
		Temperature: 0.7,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mem, err := NewMemory(MemoryConfig{
		Type:       MemoryTypeBuffer,
		HumanPrefix: "User",
		AIPrefix:    "Kled",
	})
	if err != nil {
		log.Fatalf("Failed to create memory: %v", err)
	}

	chain, err := NewChain(ChainConfig{
		Client:       client,
		Memory:       mem,
		SystemPrompt: "You are Kled, a Senior Software Engineering Lead & Technical Authority for AI/ML.",
	})
	if err != nil {
		log.Fatalf("Failed to create chain: %v", err)
	}

	ctx := context.Background()
	response, err := chain.Run(ctx, "What can you help me with?")
	if err != nil {
		log.Fatalf("Failed to run chain: %v", err)
	}

	fmt.Println("Response:", response)

	shellTool := NewShellTool()
	searchTool := NewSearchTool()
	calculatorTool := NewCalculatorTool()
	fileTool := NewFileTool()
	httpTool := NewHTTPTool()
	codeTool := NewCodeTool()

	agent, err := NewAgent(AgentConfig{
		Client:       client,
		Memory:       mem,
		Tools:        []tools.Tool{shellTool, searchTool, calculatorTool, fileTool, httpTool, codeTool},
		SystemPrompt: "You are Kled, a Senior Software Engineering Lead & Technical Authority for AI/ML.",
		MaxIterations: 5,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	response, err = agent.Run(ctx, "What is the current directory?")
	if err != nil {
		log.Fatalf("Failed to run agent: %v", err)
	}

	fmt.Println("Agent Response:", response)

	toolDefs := []schema.FunctionDefinition{
		ToolToFunctionDefinition(shellTool),
		ToolToFunctionDefinition(calculatorTool),
	}

	toolMap := map[string]schema.FunctionCallable{
		shellTool.Name():      ToolToFunctionCallable(shellTool),
		calculatorTool.Name(): ToolToFunctionCallable(calculatorTool),
	}

	response, err = client.ExecuteWithTools(ctx, "Calculate 2+2 and then list the files in the current directory", toolDefs, toolMap)
	if err != nil {
		log.Fatalf("Failed to execute with tools: %v", err)
	}

	fmt.Println("Tool Execution Response:", response)
}

func ExampleIntegration() {
	client, err := NewClient(ClientConfig{
		Type:        ModelTypeOpenAI,
		ModelName:   "gpt-4o",
		Temperature: 0.7,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mem, err := NewMemory(MemoryConfig{
		Type:       MemoryTypeBuffer,
		HumanPrefix: "User",
		AIPrefix:    "Kled",
	})
	if err != nil {
		log.Fatalf("Failed to create memory: %v", err)
	}

	shellTool := NewShellTool()
	searchTool := NewSearchTool()
	calculatorTool := NewCalculatorTool()
	fileTool := NewFileTool()
	httpTool := NewHTTPTool()
	codeTool := NewCodeTool()

	agent, err := NewAgent(AgentConfig{
		Client:       client,
		Memory:       mem,
		Tools:        []tools.Tool{shellTool, searchTool, calculatorTool, fileTool, httpTool, codeTool},
		SystemPrompt: "You are Kled, a Senior Software Engineering Lead & Technical Authority for AI/ML.",
		MaxIterations: 5,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	ctx := context.Background()
	
	userRequest := "I need help setting up a Kubernetes cluster"
	
	response, err := agent.Run(ctx, userRequest)
	if err != nil {
		log.Fatalf("Failed to process request: %v", err)
	}
	
	fmt.Println("Kled's Response:", response)
	
	codeRequest := "Write a Go function to read a file"
	
	codeResponse, err := agent.Run(ctx, codeRequest)
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}
	
	fmt.Println("Generated Code:", codeResponse)
}
