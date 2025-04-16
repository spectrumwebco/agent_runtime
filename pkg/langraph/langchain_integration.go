package langraph

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/retrievers"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/embeddings"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type LangChainAgent struct {
	Config      AgentConfig            `json:"config"`
	State       map[string]interface{} `json:"state,omitempty"`
	Tools       []Tool                 `json:"-"`
	EventStream EventStream            `json:"-"`
	stateLock   sync.RWMutex           `json:"-"`
	
	LLM          llms.LLM                 `json:"-"`
	Agent        agents.Agent             `json:"-"`
	Memory       memory.Memory            `json:"-"`
	LCTools      []tools.Tool             `json:"-"`
	Callbacks    []callbacks.Handler      `json:"-"`
	Retriever    retrievers.Retriever     `json:"-"`
	VectorStore  vectorstores.VectorStore `json:"-"`
	TextSplitter textsplitter.TextSplitter `json:"-"`
	Embedder     embeddings.Embedder      `json:"-"`
}

type LangChainConfig struct {
	AgentConfig  AgentConfig            `json:"agent_config"`
	LLMConfig    map[string]interface{} `json:"llm_config,omitempty"`
	MemoryConfig map[string]interface{} `json:"memory_config,omitempty"`
	ToolsConfig  map[string]interface{} `json:"tools_config,omitempty"`
	AgentType    string                 `json:"agent_type,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

func NewLangChainAgent(config LangChainConfig, llm llms.LLM, eventStream EventStream) (*LangChainAgent, error) {
	if llm == nil {
		return nil, fmt.Errorf("LLM is required for LangChain agent")
	}

	agent := &LangChainAgent{
		Config:      config.AgentConfig,
		State:       make(map[string]interface{}),
		Tools:       []Tool{},
		EventStream: eventStream,
		LLM:         llm,
		LCTools:     []tools.Tool{},
		Callbacks:   []callbacks.Handler{},
	}

	if config.MemoryConfig != nil {
		agent.Memory = memory.NewConversationBuffer()
	}

	agent.Callbacks = append(agent.Callbacks, &LangChainEventCallback{
		AgentID:      agent.Config.ID,
		AgentName:    agent.Config.Name,
		AgentRole:    string(agent.Config.Role),
		EventStream:  eventStream,
	})

	return agent, nil
}

func (a *LangChainAgent) Process(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"input":      input,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	inputText, ok := input["input"].(string)
	if !ok {
		if inputMap, ok := input["input"].(map[string]interface{}); ok {
			if text, ok := inputMap["text"].(string); ok {
				inputText = text
			}
		}
		
		if inputText == "" {
			inputText = fmt.Sprintf("%v", input)
		}
	}

	var result map[string]interface{}
	var err error

	if a.Agent != nil {
		agentResponse, err := a.Agent.Run(
			ctx,
			agents.NewDefaultPromptInput(inputText),
			agents.WithCallbacks(a.Callbacks...),
		)
		
		if err != nil {
			return nil, fmt.Errorf("agent execution failed: %w", err)
		}
		
		result = map[string]interface{}{
			"output": agentResponse,
		}
	} else {
		chain := chains.NewLLMChain(a.LLM, prompts.NewPromptTemplate(
			"You are a helpful assistant. Answer the following question:\n{{.input}}\nAnswer:",
			[]string{"input"},
		))
		
		chainResponse, err := chains.Call(ctx, chain, map[string]interface{}{
			"input": inputText,
		}, chains.WithCallbacks(a.Callbacks...))
		
		if err != nil {
			return nil, fmt.Errorf("chain execution failed: %w", err)
		}
		
		result = chainResponse
	}

	resultEvent := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"result":     result,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(resultEvent); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	return result, err
}

func (a *LangChainAgent) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (map[string]interface{}, error) {
	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool":       toolName,
			"parameters": params,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	var selectedTool tools.Tool
	for _, tool := range a.LCTools {
		if tool.Name() == toolName {
			selectedTool = tool
			break
		}
	}

	if selectedTool == nil {
		for _, tool := range a.Tools {
			if tool.Name == toolName {
				result, err := tool.Handler(ctx, a, params)
				
				resultEvent := models.NewEvent(
					models.EventTypeResponse,
					models.EventSourceAgent,
					map[string]interface{}{
						"tool":       toolName,
						"result":     result,
						"error":      err != nil,
						"error_msg":  err,
						"agent_id":   a.Config.ID,
						"agent_name": a.Config.Name,
						"agent_role": a.Config.Role,
					},
					map[string]string{
						"agent_id": a.Config.ID,
					},
				)
				
				if a.EventStream != nil {
					if err := a.EventStream.AddEvent(resultEvent); err != nil {
						fmt.Printf("Failed to add event to stream: %v\n", err)
					}
				}

				return result, err
			}
		}

		return nil, fmt.Errorf("tool %s not found", toolName)
	}

	var input string
	if text, ok := params["input"].(string); ok {
		input = text
	} else {
		input = fmt.Sprintf("%v", params)
	}

	output, err := selectedTool.Call(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	result := map[string]interface{}{
		"output": output,
	}

	resultEvent := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool":       toolName,
			"result":     result,
			"agent_id":   a.Config.ID,
			"agent_name": a.Config.Name,
			"agent_role": a.Config.Role,
		},
		map[string]string{
			"agent_id": a.Config.ID,
		},
	)
	
	if a.EventStream != nil {
		if err := a.EventStream.AddEvent(resultEvent); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}

	return result, nil
}

func (a *LangChainAgent) SetState(state map[string]interface{}) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	a.State = state
}

func (a *LangChainAgent) UpdateState(updates map[string]interface{}) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	
	if a.State == nil {
		a.State = make(map[string]interface{})
	}
	
	for k, v := range updates {
		a.State[k] = v
	}
}

func (a *LangChainAgent) GetState() map[string]interface{} {
	a.stateLock.RLock()
	defer a.stateLock.RUnlock()
	
	stateCopy := make(map[string]interface{})
	for k, v := range a.State {
		stateCopy[k] = v
	}
	
	return stateCopy
}

func (a *LangChainAgent) AddTool(name, description string, handler ToolHandler, parameters map[string]interface{}) {
	a.Tools = append(a.Tools, Tool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
		Handler:     handler,
	})
}

func (a *LangChainAgent) AddLangChainTool(tool tools.Tool) {
	a.LCTools = append(a.LCTools, tool)
}

func (a *LangChainAgent) SetAgent(agent agents.Agent) {
	a.Agent = agent
}

func (a *LangChainAgent) SetMemory(memory memory.Memory) {
	a.Memory = memory
}

func (a *LangChainAgent) SetRetriever(retriever retrievers.Retriever) {
	a.Retriever = retriever
}

func (a *LangChainAgent) SetVectorStore(vectorStore vectorstores.VectorStore) {
	a.VectorStore = vectorStore
}

func (a *LangChainAgent) SetTextSplitter(textSplitter textsplitter.TextSplitter) {
	a.TextSplitter = textSplitter
}

func (a *LangChainAgent) SetEmbedder(embedder embeddings.Embedder) {
	a.Embedder = embedder
}

type LangChainEventCallback struct {
	AgentID     string
	AgentName   string
	AgentRole   string
	EventStream EventStream
}

func (c *LangChainEventCallback) HandleLLMStart(ctx context.Context, llmStart *callbacks.LLMStartData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"llm_start":  true,
			"prompts":    llmStart.Prompts,
			"agent_id":   c.AgentID,
			"agent_name": c.AgentName,
			"agent_role": c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleLLMEnd(ctx context.Context, llmEnd *callbacks.LLMEndData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"llm_end":    true,
			"generations": llmEnd.Generations,
			"agent_id":   c.AgentID,
			"agent_name": c.AgentName,
			"agent_role": c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleChainStart(ctx context.Context, chainStart *callbacks.ChainStartData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"chain_start": true,
			"chain_name":  chainStart.ChainName,
			"inputs":      chainStart.Inputs,
			"agent_id":    c.AgentID,
			"agent_name":  c.AgentName,
			"agent_role":  c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleChainEnd(ctx context.Context, chainEnd *callbacks.ChainEndData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"chain_end":  true,
			"chain_name": chainEnd.ChainName,
			"outputs":    chainEnd.Outputs,
			"agent_id":   c.AgentID,
			"agent_name": c.AgentName,
			"agent_role": c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleToolStart(ctx context.Context, toolStart *callbacks.ToolStartData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool_start": true,
			"tool_name":  toolStart.ToolName,
			"input":      toolStart.Input,
			"agent_id":   c.AgentID,
			"agent_name": c.AgentName,
			"agent_role": c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleToolEnd(ctx context.Context, toolEnd *callbacks.ToolEndData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"tool_end":   true,
			"tool_name":  toolEnd.ToolName,
			"output":     toolEnd.Output,
			"agent_id":   c.AgentID,
			"agent_name": c.AgentName,
			"agent_role": c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleAgentAction(ctx context.Context, agentAction *callbacks.AgentActionData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeAction,
		models.EventSourceAgent,
		map[string]interface{}{
			"agent_action": true,
			"tool":         agentAction.Tool,
			"tool_input":   agentAction.ToolInput,
			"log":          agentAction.Log,
			"agent_id":     c.AgentID,
			"agent_name":   c.AgentName,
			"agent_role":   c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleAgentFinish(ctx context.Context, agentFinish *callbacks.AgentFinishData) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"agent_finish": true,
			"output":       agentFinish.Output,
			"log":          agentFinish.Log,
			"agent_id":     c.AgentID,
			"agent_name":   c.AgentName,
			"agent_role":   c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

func (c *LangChainEventCallback) HandleText(ctx context.Context, text string) error {
	if c.EventStream == nil {
		return nil
	}

	event := models.NewEvent(
		models.EventTypeMessage,
		models.EventSourceAgent,
		map[string]interface{}{
			"text":       text,
			"agent_id":   c.AgentID,
			"agent_name": c.AgentName,
			"agent_role": c.AgentRole,
		},
		map[string]string{
			"agent_id": c.AgentID,
		},
	)

	return c.EventStream.AddEvent(event)
}

type LangChainAgentFactory struct {
	LLM         llms.LLM
	EventStream EventStream
	agents      map[string]*LangChainAgent
	agentsMutex sync.RWMutex
}

func NewLangChainAgentFactory(llm llms.LLM, eventStream EventStream) *LangChainAgentFactory {
	return &LangChainAgentFactory{
		LLM:         llm,
		EventStream: eventStream,
		agents:      make(map[string]*LangChainAgent),
	}
}

func (f *LangChainAgentFactory) CreateAgent(config LangChainConfig) (*LangChainAgent, error) {
	agent, err := NewLangChainAgent(config, f.LLM, f.EventStream)
	if err != nil {
		return nil, err
	}

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[config.AgentConfig.ID] = agent

	return agent, nil
}

func (f *LangChainAgentFactory) GetAgent(id string) (*LangChainAgent, error) {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent %s does not exist", id)
	}

	return agent, nil
}

func (f *LangChainAgentFactory) ListAgents() []*LangChainAgent {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	var agents []*LangChainAgent
	for _, agent := range f.agents {
		agents = append(agents, agent)
	}

	return agents
}

func (f *LangChainAgentFactory) RemoveAgent(id string) error {
	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()

	if _, exists := f.agents[id]; !exists {
		return fmt.Errorf("agent %s does not exist", id)
	}

	delete(f.agents, id)
	return nil
}

func CreateLangChainTool(tool Tool) tools.Tool {
	return &LangChainToolAdapter{
		name:        tool.Name,
		description: tool.Description,
		handler:     tool.Handler,
	}
}

type LangChainToolAdapter struct {
	name        string
	description string
	handler     ToolHandler
}

func (t *LangChainToolAdapter) Name() string {
	return t.name
}

func (t *LangChainToolAdapter) Description() string {
	return t.description
}

func (t *LangChainToolAdapter) Call(ctx context.Context, input string) (string, error) {
	agent := &Agent{
		Config: AgentConfig{
			ID:   "temp",
			Name: "temp",
		},
		State: make(map[string]interface{}),
	}

	result, err := t.handler(ctx, agent, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return "", err
	}

	if output, ok := result["output"].(string); ok {
		return output, nil
	}

	return fmt.Sprintf("%v", result), nil
}

func CreateToolFromLangChainTool(lcTool tools.Tool) Tool {
	return Tool{
		Name:        lcTool.Name(),
		Description: lcTool.Description(),
		Parameters: map[string]interface{}{
			"input": "string",
		},
		Handler: func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
			input, ok := params["input"].(string)
			if !ok {
				input = fmt.Sprintf("%v", params)
			}

			output, err := lcTool.Call(ctx, input)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"output": output,
			}, nil
		},
	}
}

func UpdateMultiAgentSystemWithLangChain(system *MultiAgentSystem, llm llms.LLM) (*LangChainAgentFactory, error) {
	if llm == nil {
		return nil, fmt.Errorf("LLM is required for LangChain integration")
	}

	lcFactory := NewLangChainAgentFactory(llm, system.EventStream)

	system.Metadata["langchain_integrated"] = true

	return lcFactory, nil
}
