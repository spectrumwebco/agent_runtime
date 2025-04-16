package langgraph

import (
	"context"
	"fmt"
)

type NodeType string

const (
	ProcessorNode NodeType = "processor"
	
	ToolNode NodeType = "tool"
	
	ConditionNode NodeType = "condition"
	
	AgentNode NodeType = "agent"
)

type NodeBuilder struct {
	id       NodeID
	nodeType NodeType
	process  NodeFn
	metadata map[string]interface{}
}

func NewNodeBuilder(id NodeID, nodeType NodeType) *NodeBuilder {
	return &NodeBuilder{
		id:       id,
		nodeType: nodeType,
		metadata: make(map[string]interface{}),
	}
}

func (b *NodeBuilder) WithProcess(processFn NodeFn) *NodeBuilder {
	b.process = processFn
	return b
}

func (b *NodeBuilder) WithMetadata(key string, value interface{}) *NodeBuilder {
	b.metadata[key] = value
	return b
}

func (b *NodeBuilder) Build() *Node {
	if b.process == nil {
		b.process = func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
			return state, nil
		}
	}
	
	return &Node{
		ID:       b.id,
		Process:  b.process,
		Metadata: b.metadata,
	}
}

func ProcessorNodeBuilder(id NodeID) *NodeBuilder {
	return NewNodeBuilder(id, ProcessorNode)
}

func ToolNodeBuilder(id NodeID) *NodeBuilder {
	return NewNodeBuilder(id, ToolNode)
}

func ConditionNodeBuilder(id NodeID) *NodeBuilder {
	return NewNodeBuilder(id, ConditionNode)
}

func AgentNodeBuilder(id NodeID) *NodeBuilder {
	return NewNodeBuilder(id, AgentNode)
}

func ConditionalEdge(condition func(ctx context.Context, state map[string]interface{}) bool, target NodeID) EdgeCondition {
	return func(ctx context.Context, state map[string]interface{}) (NodeID, error) {
		if condition(ctx, state) {
			return target, nil
		}
		return "", fmt.Errorf("condition not met")
	}
}

func DefaultEdge(target NodeID) EdgeCondition {
	return func(ctx context.Context, state map[string]interface{}) (NodeID, error) {
		return target, nil
	}
}

func BranchEdge(key string, branches map[string]NodeID, defaultTarget NodeID) EdgeCondition {
	return func(ctx context.Context, state map[string]interface{}) (NodeID, error) {
		val, exists := state[key]
		if !exists {
			if defaultTarget != "" {
				return defaultTarget, nil
			}
			return "", fmt.Errorf("key %s not found in state", key)
		}
		
		strVal, ok := val.(string)
		if !ok {
			if defaultTarget != "" {
				return defaultTarget, nil
			}
			return "", fmt.Errorf("key %s is not a string", key)
		}
		
		target, exists := branches[strVal]
		if !exists {
			if defaultTarget != "" {
				return defaultTarget, nil
			}
			return "", fmt.Errorf("no branch for value %s", strVal)
		}
		
		return target, nil
	}
}
