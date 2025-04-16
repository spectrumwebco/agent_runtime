package langgraph

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config         *config.Config
	agentIntegration *langgraph.AgentIntegration
	djangoBaseURL  string
}

func NewModule(cfg *config.Config, djangoBaseURL string) *Module {
	return &Module{
		config:         cfg,
		agentIntegration: langgraph.NewAgentIntegration(cfg, djangoBaseURL),
		djangoBaseURL:  djangoBaseURL,
	}
}

func (m *Module) Initialize() error {
	_, err := m.agentIntegration.BuildAgentGraph()
	if err != nil {
		return fmt.Errorf("failed to build agent graph: %w", err)
	}
	
	return nil
}

func (m *Module) ExecuteAgentWorkflow(ctx context.Context, input string, tools []map[string]interface{}) (map[string]interface{}, error) {
	return m.agentIntegration.ExecuteAgentWorkflow(ctx, input, tools)
}

func (m *Module) GetAgentIntegration() *langgraph.AgentIntegration {
	return m.agentIntegration
}

func (m *Module) BuildGraph(entryPoint langgraph.NodeID, nodes map[langgraph.NodeID]langgraph.NodeFn, edges map[langgraph.NodeID][]langgraph.Edge) (*langgraph.Graph, error) {
	graph := langgraph.NewGraph(entryPoint)
	
	for id, processFn := range nodes {
		err := graph.AddNode(id, processFn, map[string]interface{}{
			"description": fmt.Sprintf("Node %s", id),
		})
		if err != nil {
			return nil, err
		}
	}
	
	for source, edgeList := range edges {
		for _, edge := range edgeList {
			err := graph.AddEdge(source, edge.Target, edge.Condition)
			if err != nil {
				return nil, err
			}
		}
	}
	
	err := graph.Compile()
	if err != nil {
		return nil, err
	}
	
	return graph, nil
}

func (m *Module) ExecuteGraph(ctx context.Context, graph *langgraph.Graph, initialState map[string]interface{}) (map[string]interface{}, error) {
	return graph.Execute(ctx, initialState)
}

func (m *Module) CreateStateManager(initialState map[string]interface{}, schema map[string]string) *langgraph.StateManager {
	return langgraph.NewStateManager(initialState, schema)
}

func (m *Module) CreateNodeBuilder(id langgraph.NodeID, nodeType langgraph.NodeType) *langgraph.NodeBuilder {
	return langgraph.NewNodeBuilder(id, nodeType)
}

func (m *Module) CreateProcessorNode(id langgraph.NodeID, processFn langgraph.NodeFn) *langgraph.Node {
	return langgraph.ProcessorNodeBuilder(id).WithProcess(processFn).Build()
}

func (m *Module) CreateToolNode(id langgraph.NodeID, processFn langgraph.NodeFn) *langgraph.Node {
	return langgraph.ToolNodeBuilder(id).WithProcess(processFn).Build()
}

func (m *Module) CreateAgentNode(id langgraph.NodeID, processFn langgraph.NodeFn) *langgraph.Node {
	return langgraph.AgentNodeBuilder(id).WithProcess(processFn).Build()
}

func (m *Module) CreateConditionNode(id langgraph.NodeID, processFn langgraph.NodeFn) *langgraph.Node {
	return langgraph.ConditionNodeBuilder(id).WithProcess(processFn).Build()
}

func (m *Module) CreateDefaultEdge(target langgraph.NodeID) langgraph.EdgeCondition {
	return langgraph.DefaultEdge(target)
}

func (m *Module) CreateConditionalEdge(condition func(ctx context.Context, state map[string]interface{}) bool, target langgraph.NodeID) langgraph.EdgeCondition {
	return langgraph.ConditionalEdge(condition, target)
}

func (m *Module) CreateBranchEdge(key string, branches map[string]langgraph.NodeID, defaultTarget langgraph.NodeID) langgraph.EdgeCondition {
	return langgraph.BranchEdge(key, branches, defaultTarget)
}

func (m *Module) ChainStateUpdaters(updaters ...langgraph.StateUpdater) langgraph.StateUpdater {
	return langgraph.ChainStateUpdaters(updaters...)
}

func (m *Module) CreateConditionalStateUpdater(condition func(ctx context.Context, state map[string]interface{}) bool, updater langgraph.StateUpdater) langgraph.StateUpdater {
	return langgraph.ConditionalStateUpdater(condition, updater)
}

func (m *Module) CreateKeyValueStateUpdater(key string, value interface{}) langgraph.StateUpdater {
	return langgraph.KeyValueStateUpdater(key, value)
}

func (m *Module) CreateDynamicKeyValueStateUpdater(key string, valueFn func(ctx context.Context, state map[string]interface{}) (interface{}, error)) langgraph.StateUpdater {
	return langgraph.DynamicKeyValueStateUpdater(key, valueFn)
}
