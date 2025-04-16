package langgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type NodeID string

type EdgeCondition func(ctx context.Context, state map[string]interface{}) (NodeID, error)

type NodeFn func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error)

type Node struct {
	ID       NodeID
	Process  NodeFn
	Metadata map[string]interface{}
}

type Edge struct {
	Source    NodeID
	Target    NodeID
	Condition EdgeCondition
}

type Graph struct {
	Nodes       map[NodeID]*Node
	Edges       map[NodeID][]Edge
	EntryPoint  NodeID
	ExitPoints  map[NodeID]bool
	StateSchema map[string]string
	mu          sync.RWMutex
}

func NewGraph(entryPoint NodeID) *Graph {
	return &Graph{
		Nodes:       make(map[NodeID]*Node),
		Edges:       make(map[NodeID][]Edge),
		EntryPoint:  entryPoint,
		ExitPoints:  make(map[NodeID]bool),
		StateSchema: make(map[string]string),
	}
}

func (g *Graph) AddNode(id NodeID, processFn NodeFn, metadata map[string]interface{}) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Nodes[id]; exists {
		return fmt.Errorf("node with ID %s already exists", id)
	}

	g.Nodes[id] = &Node{
		ID:       id,
		Process:  processFn,
		Metadata: metadata,
	}

	return nil
}

func (g *Graph) AddEdge(source NodeID, target NodeID, condition EdgeCondition) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Nodes[source]; !exists {
		return fmt.Errorf("source node %s does not exist", source)
	}

	if _, exists := g.Nodes[target]; !exists {
		return fmt.Errorf("target node %s does not exist", target)
	}

	edge := Edge{
		Source:    source,
		Target:    target,
		Condition: condition,
	}

	g.Edges[source] = append(g.Edges[source], edge)
	return nil
}

func (g *Graph) SetExitPoint(id NodeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Nodes[id]; !exists {
		return fmt.Errorf("node %s does not exist", id)
	}

	g.ExitPoints[id] = true
	return nil
}

func (g *Graph) DefineStateSchema(schema map[string]string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.StateSchema = schema
}

func (g *Graph) ValidateState(state map[string]interface{}) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for key, expectedType := range g.StateSchema {
		val, exists := state[key]
		if !exists {
			return fmt.Errorf("required state key %s is missing", key)
		}

		switch expectedType {
		case "string":
			if _, ok := val.(string); !ok {
				return fmt.Errorf("state key %s should be a string", key)
			}
		case "int":
			if _, ok := val.(int); !ok {
				return fmt.Errorf("state key %s should be an int", key)
			}
		case "bool":
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("state key %s should be a bool", key)
			}
		case "map":
			if _, ok := val.(map[string]interface{}); !ok {
				return fmt.Errorf("state key %s should be a map", key)
			}
		case "array":
			if _, ok := val.([]interface{}); !ok {
				return fmt.Errorf("state key %s should be an array", key)
			}
		}
	}

	return nil
}

func (g *Graph) Execute(ctx context.Context, initialState map[string]interface{}) (map[string]interface{}, error) {
	g.mu.RLock()
	if _, exists := g.Nodes[g.EntryPoint]; !exists {
		g.mu.RUnlock()
		return nil, fmt.Errorf("entry point node %s does not exist", g.EntryPoint)
	}
	g.mu.RUnlock()

	if err := g.ValidateState(initialState); err != nil {
		return nil, err
	}

	state := initialState
	currentNode := g.EntryPoint
	visited := make(map[NodeID]bool)
	trajectory := []map[string]interface{}{}

	stateCopy := make(map[string]interface{})
	for k, v := range state {
		stateCopy[k] = v
	}
	stateCopy["node"] = string(currentNode)
	trajectory = append(trajectory, stateCopy)

	for {
		g.mu.RLock()
		node, exists := g.Nodes[currentNode]
		if !exists {
			g.mu.RUnlock()
			return nil, fmt.Errorf("node %s does not exist", currentNode)
		}

		if visited[currentNode] {
			g.mu.RUnlock()
			return nil, fmt.Errorf("cycle detected at node %s", currentNode)
		}
		visited[currentNode] = true

		g.mu.RUnlock()
		var err error
		state, err = node.Process(ctx, state)
		if err != nil {
			return nil, fmt.Errorf("error processing node %s: %w", currentNode, err)
		}

		stateCopy = make(map[string]interface{})
		for k, v := range state {
			stateCopy[k] = v
		}
		stateCopy["node"] = string(currentNode)
		trajectory = append(trajectory, stateCopy)

		g.mu.RLock()
		if g.ExitPoints[currentNode] {
			g.mu.RUnlock()
			trajectoryJSON, _ := json.Marshal(trajectory)
			state["trajectory"] = string(trajectoryJSON)
			return state, nil
		}

		edges := g.Edges[currentNode]
		g.mu.RUnlock()

		if len(edges) == 0 {
			return nil, fmt.Errorf("node %s has no outgoing edges", currentNode)
		}

		nextNode := NodeID("")
		for _, edge := range edges {
			var conditionErr error
			nextNode, conditionErr = edge.Condition(ctx, state)
			if conditionErr == nil && nextNode != "" {
				break
			}
		}

		if nextNode == "" {
			return nil, errors.New("no valid edge condition found")
		}

		currentNode = nextNode
	}
}

func (g *Graph) Compile() error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.Nodes[g.EntryPoint]; !exists {
		return fmt.Errorf("entry point node %s does not exist", g.EntryPoint)
	}

	for exitPoint := range g.ExitPoints {
		if _, exists := g.Nodes[exitPoint]; !exists {
			return fmt.Errorf("exit point node %s does not exist", exitPoint)
		}
	}

	for source, edges := range g.Edges {
		if _, exists := g.Nodes[source]; !exists {
			return fmt.Errorf("edge source node %s does not exist", source)
		}

		for _, edge := range edges {
			if _, exists := g.Nodes[edge.Target]; !exists {
				return fmt.Errorf("edge target node %s does not exist", edge.Target)
			}
		}
	}

	reachable := make(map[NodeID]bool)
	var dfs func(NodeID)
	dfs = func(node NodeID) {
		if reachable[node] {
			return
		}
		reachable[node] = true
		for _, edge := range g.Edges[node] {
			dfs(edge.Target)
		}
	}
	dfs(g.EntryPoint)

	for id := range g.Nodes {
		if !reachable[id] {
			return fmt.Errorf("node %s is unreachable", id)
		}
	}

	return nil
}
