package langraph

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type LangChainBridge struct {
	Graph         *Graph
	Executor      *Executor
	LLM           llms.LLM
	EventStream   EventStream
	LCAgentFactory *LangChainAgentFactory
	agentNodes    map[string]NodeID
	toolNodes     map[string]NodeID
	chainNodes    map[string]NodeID
	mutex         sync.RWMutex
}

func NewLangChainBridge(graph *Graph, executor *Executor, llm llms.LLM, eventStream EventStream) *LangChainBridge {
	lcFactory := NewLangChainAgentFactory(llm, eventStream)
	
	return &LangChainBridge{
		Graph:         graph,
		Executor:      executor,
		LLM:           llm,
		EventStream:   eventStream,
		LCAgentFactory: lcFactory,
		agentNodes:    make(map[string]NodeID),
		toolNodes:     make(map[string]NodeID),
		chainNodes:    make(map[string]NodeID),
	}
}

func (b *LangChainBridge) AddLangChainAgent(agent *LangChainAgent, name, description string) (NodeID, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	node := b.Graph.AddAgentNode(
		AgentType(agent.Config.Role),
		name,
		description,
		func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			return agent.Process(ctx, inputs)
		},
	)

	node.Metadata["agent_id"] = agent.Config.ID
	node.Metadata["agent_role"] = string(agent.Config.Role)
	node.Metadata["agent_type"] = "langchain"

	b.agentNodes[agent.Config.ID] = node.ID

	return node.ID, nil
}

func (b *LangChainBridge) AddLangChainTool(tool tools.Tool, name, description string) (NodeID, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	node := b.Graph.AddToolNode(
		name,
		description,
		func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			input, ok := inputs["input"].(string)
			if !ok {
				input = fmt.Sprintf("%v", inputs)
			}

			output, err := tool.Call(ctx, input)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"output": output,
			}, nil
		},
	)

	node.Metadata["tool_name"] = tool.Name()
	node.Metadata["tool_type"] = "langchain"

	b.toolNodes[tool.Name()] = node.ID

	return node.ID, nil
}

func (b *LangChainBridge) AddLangChainChain(chain chains.Chain, name, description string) (NodeID, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	node := b.Graph.AddNode(
		name,
		description,
		func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			return chains.Call(ctx, chain, inputs)
		},
	)

	node.Metadata["chain_name"] = name
	node.Metadata["chain_type"] = "langchain"

	b.chainNodes[name] = node.ID

	return node.ID, nil
}

func (b *LangChainBridge) CreateLangChainAgentNode(config LangChainConfig, name, description string) (NodeID, error) {
	agent, err := b.LCAgentFactory.CreateAgent(config)
	if err != nil {
		return "", err
	}

	return b.AddLangChainAgent(agent, name, description)
}

func (b *LangChainBridge) CreateLangChainToolNode(toolName, toolDescription string, handler func(ctx context.Context, input string) (string, error)) (NodeID, error) {
	tool := &LangChainToolWrapper{
		name:        toolName,
		description: toolDescription,
		handler:     handler,
	}

	return b.AddLangChainTool(tool, toolName, toolDescription)
}

func (b *LangChainBridge) CreateLangChainChainNode(prompt string, name, description string) (NodeID, error) {
	chain := chains.NewLLMChain(b.LLM, prompts.NewPromptTemplate(
		prompt,
		[]string{"input"},
	))

	return b.AddLangChainChain(chain, name, description)
}

func (b *LangChainBridge) CreateLangChainAgentExecutor(agentType string, tools []tools.Tool, name, description string) (NodeID, error) {
	var agent agents.Agent
	var err error

	switch agentType {
	case "zero-shot":
		agent, err = agents.NewZeroShotAgent(b.LLM, tools)
	case "conversational":
		agent, err = agents.NewConversationalAgent(b.LLM, tools, memory.NewConversationBuffer())
	case "react":
		agent, err = agents.NewReactAgent(b.LLM, tools)
	default:
		return "", fmt.Errorf("unsupported agent type: %s", agentType)
	}

	if err != nil {
		return "", err
	}

	node := b.Graph.AddAgentNode(
		AgentType(agentType),
		name,
		description,
		func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error) {
			input, ok := inputs["input"].(string)
			if !ok {
				input = fmt.Sprintf("%v", inputs)
			}

			output, err := agent.Run(
				ctx,
				agents.NewDefaultPromptInput(input),
			)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"output": output,
			}, nil
		},
	)

	node.Metadata["agent_type"] = agentType
	node.Metadata["agent_framework"] = "langchain"

	b.agentNodes[name] = node.ID

	return node.ID, nil
}

func (b *LangChainBridge) ConnectLangChainNodes(sourceID, targetID NodeID, name, description string) (*Edge, error) {
	return b.Graph.AddEdge(sourceID, targetID, name, description)
}

func (b *LangChainBridge) ExecuteLangChainGraph(ctx context.Context, startNodeID NodeID, inputs map[string]interface{}) (*Execution, error) {
	return b.Executor.Execute(ctx, startNodeID, inputs, nil)
}

type LangChainToolWrapper struct {
	name        string
	description string
	handler     func(ctx context.Context, input string) (string, error)
}

func (t *LangChainToolWrapper) Name() string {
	return t.name
}

func (t *LangChainToolWrapper) Description() string {
	return t.description
}

func (t *LangChainToolWrapper) Call(ctx context.Context, input string) (string, error) {
	return t.handler(ctx, input)
}

func CreateMultiAgentSystemWithLangChain(name, description string, llm llms.LLM, eventStream EventStream) (*MultiAgentSystem, *LangChainBridge, error) {
	system := NewMultiAgentSystem(name, description, eventStream)

	bridge := NewLangChainBridge(system.Graph, system.Executor, llm, eventStream)

	return system, bridge, nil
}

type LangChainGraphBuilder struct {
	Bridge      *LangChainBridge
	System      *MultiAgentSystem
	LLM         llms.LLM
	EventStream EventStream
	nodes       map[string]NodeID
	mutex       sync.RWMutex
}

func NewLangChainGraphBuilder(system *MultiAgentSystem, llm llms.LLM, eventStream EventStream) *LangChainGraphBuilder {
	bridge := NewLangChainBridge(system.Graph, system.Executor, llm, eventStream)
	
	return &LangChainGraphBuilder{
		Bridge:      bridge,
		System:      system,
		LLM:         llm,
		EventStream: eventStream,
		nodes:       make(map[string]NodeID),
	}
}

func (b *LangChainGraphBuilder) AddAgent(name, description, agentType string, tools []tools.Tool) *LangChainGraphBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	nodeID, err := b.Bridge.CreateLangChainAgentExecutor(agentType, tools, name, description)
	if err != nil {
		fmt.Printf("Error creating agent node: %v\n", err)
		return b
	}

	b.nodes[name] = nodeID
	return b
}

func (b *LangChainGraphBuilder) AddTool(name, description string, handler func(ctx context.Context, input string) (string, error)) *LangChainGraphBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	nodeID, err := b.Bridge.CreateLangChainToolNode(name, description, handler)
	if err != nil {
		fmt.Printf("Error creating tool node: %v\n", err)
		return b
	}

	b.nodes[name] = nodeID
	return b
}

func (b *LangChainGraphBuilder) AddChain(name, description, prompt string) *LangChainGraphBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	nodeID, err := b.Bridge.CreateLangChainChainNode(prompt, name, description)
	if err != nil {
		fmt.Printf("Error creating chain node: %v\n", err)
		return b
	}

	b.nodes[name] = nodeID
	return b
}

func (b *LangChainGraphBuilder) Connect(sourceName, targetName, edgeName, edgeDescription string) *LangChainGraphBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	sourceID, sourceExists := b.nodes[sourceName]
	targetID, targetExists := b.nodes[targetName]

	if !sourceExists {
		fmt.Printf("Source node %s does not exist\n", sourceName)
		return b
	}

	if !targetExists {
		fmt.Printf("Target node %s does not exist\n", targetName)
		return b
	}

	_, err := b.Bridge.ConnectLangChainNodes(sourceID, targetID, edgeName, edgeDescription)
	if err != nil {
		fmt.Printf("Error connecting nodes: %v\n", err)
		return b
	}

	return b
}

func (b *LangChainGraphBuilder) Build() *MultiAgentSystem {
	return b.System
}

func (b *LangChainGraphBuilder) Execute(ctx context.Context, startNodeName string, inputs map[string]interface{}) (*Execution, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	startNodeID, exists := b.nodes[startNodeName]
	if !exists {
		return nil, fmt.Errorf("start node %s does not exist", startNodeName)
	}

	return b.Bridge.ExecuteLangChainGraph(ctx, startNodeID, inputs)
}

func CreateStandardLangChainAgentSystem(name, description string, llm llms.LLM, eventStream EventStream) (*MultiAgentSystem, *LangChainGraphBuilder, error) {
	system := NewMultiAgentSystem(name, description, eventStream)

	builder := NewLangChainGraphBuilder(system, llm, eventStream)

	searchTool := &LangChainToolWrapper{
		name:        "search",
		description: "Search for information on the web",
		handler: func(ctx context.Context, input string) (string, error) {
			return fmt.Sprintf("Search results for: %s", input), nil
		},
	}

	calculatorTool := &LangChainToolWrapper{
		name:        "calculator",
		description: "Perform calculations",
		handler: func(ctx context.Context, input string) (string, error) {
			return fmt.Sprintf("Calculation result for: %s", input), nil
		},
	}

	builder.AddAgent("frontend_agent", "Handles UI/UX development", "react", []tools.Tool{searchTool, calculatorTool})
	builder.AddAgent("app_builder_agent", "Handles application assembly", "react", []tools.Tool{searchTool, calculatorTool})
	builder.AddAgent("codegen_agent", "Handles code generation", "react", []tools.Tool{searchTool, calculatorTool})
	builder.AddAgent("engineering_agent", "Handles core development tasks", "react", []tools.Tool{searchTool, calculatorTool})
	builder.AddAgent("orchestrator_agent", "Handles task orchestration", "react", []tools.Tool{searchTool, calculatorTool})

	builder.Connect("orchestrator_agent", "frontend_agent", "Orchestrator -> Frontend", "Orchestrator to Frontend communication")
	builder.Connect("orchestrator_agent", "app_builder_agent", "Orchestrator -> App Builder", "Orchestrator to App Builder communication")
	builder.Connect("orchestrator_agent", "codegen_agent", "Orchestrator -> Codegen", "Orchestrator to Codegen communication")
	builder.Connect("orchestrator_agent", "engineering_agent", "Orchestrator -> Engineering", "Orchestrator to Engineering communication")
	
	builder.Connect("frontend_agent", "app_builder_agent", "Frontend -> App Builder", "Frontend to App Builder communication")
	builder.Connect("frontend_agent", "codegen_agent", "Frontend -> Codegen", "Frontend to Codegen communication")
	
	builder.Connect("app_builder_agent", "codegen_agent", "App Builder -> Codegen", "App Builder to Codegen communication")
	builder.Connect("app_builder_agent", "engineering_agent", "App Builder -> Engineering", "App Builder to Engineering communication")
	
	builder.Connect("codegen_agent", "engineering_agent", "Codegen -> Engineering", "Codegen to Engineering communication")
	
	builder.Connect("engineering_agent", "frontend_agent", "Engineering -> Frontend", "Engineering to Frontend communication")
	builder.Connect("engineering_agent", "app_builder_agent", "Engineering -> App Builder", "Engineering to App Builder communication")
	builder.Connect("engineering_agent", "codegen_agent", "Engineering -> Codegen", "Engineering to Codegen communication")

	return system, builder, nil
}
