package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type NodeID string

type EdgeID string

type NodeType string

type AgentType string

const (
	AgentTypeFrontend    AgentType = "frontend"
	AgentTypeAppBuilder  AgentType = "app_builder"
	AgentTypeCodegen     AgentType = "codegen"
	AgentTypeEngineering AgentType = "engineering"

	NodeTypeAgent   NodeType = "agent"
	NodeTypeTask    NodeType = "task"
	NodeTypeData    NodeType = "data"
	NodeTypeService NodeType = "service"
)

type Node struct {
	ID          NodeID                  `json:"id"`
	Type        NodeType                `json:"type"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty"`
	State       map[string]interface{}  `json:"state,omitempty"`
	Inputs      map[string]interface{}  `json:"inputs,omitempty"`
	Outputs     map[string]interface{}  `json:"outputs,omitempty"`
	Handler     NodeHandler             `json:"-"`
	AgentType   AgentType               `json:"agent_type,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	stateLock   sync.RWMutex            `json:"-"`
}

type Edge struct {
	ID          EdgeID                 `json:"id"`
	Source      NodeID                 `json:"source"`
	Target      NodeID                 `json:"target"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Weight      float64                `json:"weight"`
	CreatedAt   time.Time              `json:"created_at"`
}

type NodeHandler func(ctx context.Context, node *Node, inputs map[string]interface{}) (map[string]interface{}, error)

type Graph struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Nodes       map[NodeID]*Node       `json:"nodes"`
	Edges       map[EdgeID]*Edge       `json:"edges"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	nodeLock    sync.RWMutex           `json:"-"`
	edgeLock    sync.RWMutex           `json:"-"`
}

func NewGraph(name, description string) *Graph {
	return &Graph{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Nodes:       make(map[NodeID]*Node),
		Edges:       make(map[EdgeID]*Edge),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Metadata:    make(map[string]interface{}),
	}
}

func (g *Graph) AddNode(nodeType NodeType, name, description string, handler NodeHandler) *Node {
	g.nodeLock.Lock()
	defer g.nodeLock.Unlock()

	node := &Node{
		ID:          NodeID(uuid.New().String()),
		Type:        nodeType,
		Name:        name,
		Description: description,
		Metadata:    make(map[string]interface{}),
		State:       make(map[string]interface{}),
		Inputs:      make(map[string]interface{}),
		Outputs:     make(map[string]interface{}),
		Handler:     handler,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	g.Nodes[node.ID] = node
	g.UpdatedAt = time.Now().UTC()

	return node
}

func (g *Graph) AddAgentNode(agentType AgentType, name, description string, handler NodeHandler) *Node {
	node := g.AddNode(NodeTypeAgent, name, description, handler)
	node.AgentType = agentType
	node.Metadata["agent_type"] = string(agentType)
	return node
}

func (g *Graph) AddEdge(source, target NodeID, name, description string) (*Edge, error) {
	g.nodeLock.RLock()
	_, sourceExists := g.Nodes[source]
	_, targetExists := g.Nodes[target]
	g.nodeLock.RUnlock()

	if !sourceExists {
		return nil, fmt.Errorf("source node %s does not exist", source)
	}
	if !targetExists {
		return nil, fmt.Errorf("target node %s does not exist", target)
	}

	g.edgeLock.Lock()
	defer g.edgeLock.Unlock()

	edge := &Edge{
		ID:          EdgeID(uuid.New().String()),
		Source:      source,
		Target:      target,
		Name:        name,
		Description: description,
		Metadata:    make(map[string]interface{}),
		Weight:      1.0,
		CreatedAt:   time.Now().UTC(),
	}

	g.Edges[edge.ID] = edge
	g.UpdatedAt = time.Now().UTC()

	return edge, nil
}

func (g *Graph) RemoveNode(id NodeID) error {
	g.nodeLock.Lock()
	defer g.nodeLock.Unlock()

	if _, exists := g.Nodes[id]; !exists {
		return fmt.Errorf("node %s does not exist", id)
	}

	g.edgeLock.Lock()
	defer g.edgeLock.Unlock()

	for edgeID, edge := range g.Edges {
		if edge.Source == id || edge.Target == id {
			delete(g.Edges, edgeID)
		}
	}

	delete(g.Nodes, id)
	g.UpdatedAt = time.Now().UTC()

	return nil
}

func (g *Graph) RemoveEdge(id EdgeID) error {
	g.edgeLock.Lock()
	defer g.edgeLock.Unlock()

	if _, exists := g.Edges[id]; !exists {
		return fmt.Errorf("edge %s does not exist", id)
	}

	delete(g.Edges, id)
	g.UpdatedAt = time.Now().UTC()

	return nil
}

func (g *Graph) GetNode(id NodeID) (*Node, error) {
	g.nodeLock.RLock()
	defer g.nodeLock.RUnlock()

	node, exists := g.Nodes[id]
	if !exists {
		return nil, fmt.Errorf("node %s does not exist", id)
	}

	return node, nil
}

func (g *Graph) GetEdge(id EdgeID) (*Edge, error) {
	g.edgeLock.RLock()
	defer g.edgeLock.RUnlock()

	edge, exists := g.Edges[id]
	if !exists {
		return nil, fmt.Errorf("edge %s does not exist", id)
	}

	return edge, nil
}

func (g *Graph) GetOutgoingEdges(nodeID NodeID) []*Edge {
	g.edgeLock.RLock()
	defer g.edgeLock.RUnlock()

	var edges []*Edge
	for _, edge := range g.Edges {
		if edge.Source == nodeID {
			edges = append(edges, edge)
		}
	}

	return edges
}

func (g *Graph) GetIncomingEdges(nodeID NodeID) []*Edge {
	g.edgeLock.RLock()
	defer g.edgeLock.RUnlock()

	var edges []*Edge
	for _, edge := range g.Edges {
		if edge.Target == nodeID {
			edges = append(edges, edge)
		}
	}

	return edges
}

func (g *Graph) GetNodesByType(nodeType NodeType) []*Node {
	g.nodeLock.RLock()
	defer g.nodeLock.RUnlock()

	var nodes []*Node
	for _, node := range g.Nodes {
		if node.Type == nodeType {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (g *Graph) GetAgentNodes() []*Node {
	return g.GetNodesByType(NodeTypeAgent)
}

func (g *Graph) GetAgentNodesByType(agentType AgentType) []*Node {
	g.nodeLock.RLock()
	defer g.nodeLock.RUnlock()

	var nodes []*Node
	for _, node := range g.Nodes {
		if node.Type == NodeTypeAgent && node.AgentType == agentType {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (g *Graph) SetNodeState(nodeID NodeID, state map[string]interface{}) error {
	node, err := g.GetNode(nodeID)
	if err != nil {
		return err
	}

	node.stateLock.Lock()
	defer node.stateLock.Unlock()

	node.State = state
	node.UpdatedAt = time.Now().UTC()

	return nil
}

func (g *Graph) UpdateNodeState(nodeID NodeID, updates map[string]interface{}) error {
	node, err := g.GetNode(nodeID)
	if err != nil {
		return err
	}

	node.stateLock.Lock()
	defer node.stateLock.Unlock()

	if node.State == nil {
		node.State = make(map[string]interface{})
	}

	for k, v := range updates {
		node.State[k] = v
	}
	node.UpdatedAt = time.Now().UTC()

	return nil
}

func (g *Graph) GetNodeState(nodeID NodeID) (map[string]interface{}, error) {
	node, err := g.GetNode(nodeID)
	if err != nil {
		return nil, err
	}

	node.stateLock.RLock()
	defer node.stateLock.RUnlock()

	stateCopy := make(map[string]interface{})
	for k, v := range node.State {
		stateCopy[k] = v
	}

	return stateCopy, nil
}

func (g *Graph) ProcessNode(ctx context.Context, nodeID NodeID, inputs map[string]interface{}) (map[string]interface{}, error) {
	node, err := g.GetNode(nodeID)
	if err != nil {
		return nil, err
	}

	if node.Handler == nil {
		return nil, fmt.Errorf("node %s has no handler", nodeID)
	}

	node.stateLock.Lock()
	node.Inputs = inputs
	node.stateLock.Unlock()

	outputs, err := node.Handler(ctx, node, inputs)
	if err != nil {
		return nil, fmt.Errorf("error processing node %s: %w", nodeID, err)
	}

	node.stateLock.Lock()
	node.Outputs = outputs
	node.UpdatedAt = time.Now().UTC()
	node.stateLock.Unlock()

	return outputs, nil
}

func (g *Graph) Traverse(ctx context.Context, startNodeID NodeID, initialInputs map[string]interface{}) (map[NodeID]map[string]interface{}, error) {
	results := make(map[NodeID]map[string]interface{})
	visited := make(map[NodeID]bool)

	var traverse func(nodeID NodeID, inputs map[string]interface{}) error
	traverse = func(nodeID NodeID, inputs map[string]interface{}) error {
		if visited[nodeID] {
			return nil
		}
		visited[nodeID] = true

		outputs, err := g.ProcessNode(ctx, nodeID, inputs)
		if err != nil {
			return err
		}

		results[nodeID] = outputs

		outgoingEdges := g.GetOutgoingEdges(nodeID)
		for _, edge := range outgoingEdges {
			if err := traverse(edge.Target, outputs); err != nil {
				return err
			}
		}

		return nil
	}

	if err := traverse(startNodeID, initialInputs); err != nil {
		return nil, err
	}

	return results, nil
}

func (g *Graph) TraverseAsync(ctx context.Context, startNodeID NodeID, initialInputs map[string]interface{}) (map[NodeID]map[string]interface{}, error) {
	results := make(map[NodeID]map[string]interface{})
	resultsMutex := sync.Mutex{}
	visited := make(map[NodeID]bool)
	visitedMutex := sync.Mutex{}
	var wg sync.WaitGroup

	var traverse func(nodeID NodeID, inputs map[string]interface{})
	traverse = func(nodeID NodeID, inputs map[string]interface{}) {
		defer wg.Done()

		visitedMutex.Lock()
		if visited[nodeID] {
			visitedMutex.Unlock()
			return
		}
		visited[nodeID] = true
		visitedMutex.Unlock()

		outputs, err := g.ProcessNode(ctx, nodeID, inputs)
		if err != nil {
			fmt.Printf("Error processing node %s: %v\n", nodeID, err)
			return
		}

		resultsMutex.Lock()
		results[nodeID] = outputs
		resultsMutex.Unlock()

		outgoingEdges := g.GetOutgoingEdges(nodeID)
		for _, edge := range outgoingEdges {
			wg.Add(1)
			go traverse(edge.Target, outputs)
		}
	}

	wg.Add(1)
	traverse(startNodeID, initialInputs)
	wg.Wait()

	return results, nil
}
