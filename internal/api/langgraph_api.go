package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spectrumwebco/agent_runtime/internal/langchain"
	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}
	
	activeGraphs     = make(map[string]*langgraph.Graph)
	activeGraphsLock sync.RWMutex
)

type GraphState struct {
	GraphID    string                   `json:"graphId"`
	Nodes      []langgraph.NodeState    `json:"nodes"`
	Edges      []langgraph.EdgeState    `json:"edges"`
	CurrentNode string                  `json:"currentNode"`
	State      map[string]interface{}   `json:"state"`
	Status     string                   `json:"status"` // "running", "completed", "error"
	Error      string                   `json:"error,omitempty"`
}

type NodeState struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Status   string `json:"status"` // "idle", "active", "completed", "error"
}

type EdgeState struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

type AgentConfig struct {
	AgentType  string                 `json:"agentType"` // "swe", "ui", "scaffolding", "codegen"
	ModelID    string                 `json:"modelId"`   // "gemini-2.5-pro", "llama-4", etc.
	InitialState map[string]interface{} `json:"initialState,omitempty"`
}

func RegisterLangGraphRoutes(router *gin.Engine) {
	router.POST("/api/langgraph/create", createGraph)
	router.GET("/api/langgraph/state/:graphId", getGraphState)
	router.POST("/api/langgraph/execute/:graphId", executeGraphStep)
	router.GET("/api/ws/langgraph/:graphId", websocketGraphUpdates)
}

func createGraph(c *gin.Context) {
	var config AgentConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}
	
	graph, err := createAgentGraph(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create agent graph: %v", err)})
		return
	}
	
	graphID := fmt.Sprintf("%s-%s-%d", config.AgentType, config.ModelID, len(activeGraphs))
	
	activeGraphsLock.Lock()
	activeGraphs[graphID] = graph
	activeGraphsLock.Unlock()
	
	c.JSON(http.StatusOK, gin.H{
		"graphId": graphID,
		"state": getGraphStateObject(graphID, graph, nil, "created"),
	})
}

func getGraphState(c *gin.Context) {
	graphID := c.Param("graphId")
	
	activeGraphsLock.RLock()
	graph, exists := activeGraphs[graphID]
	activeGraphsLock.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Graph not found"})
		return
	}
	
	c.JSON(http.StatusOK, getGraphStateObject(graphID, graph, nil, "idle"))
}

func executeGraphStep(c *gin.Context) {
	graphID := c.Param("graphId")
	
	activeGraphsLock.RLock()
	graph, exists := activeGraphs[graphID]
	activeGraphsLock.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Graph not found"})
		return
	}
	
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}
	
	ctx := context.Background()
	result, err := graph.Execute(ctx, input)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to execute graph: %v", err),
			"state": getGraphStateObject(graphID, graph, nil, "error"),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"state": getGraphStateObject(graphID, graph, result, "completed"),
	})
}

func websocketGraphUpdates(c *gin.Context) {
	graphID := c.Param("graphId")
	
	activeGraphsLock.RLock()
	_, exists := activeGraphs[graphID]
	activeGraphsLock.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Graph not found"})
		return
	}
	
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upgrade connection: %v", err)})
		return
	}
	defer conn.Close()
	
	updateChan := make(chan GraphState)
	
	
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		
		if err := conn.WriteMessage(messageType, message); err != nil {
			break
		}
		
		select {
		case update := <-updateChan:
			updateJSON, _ := json.Marshal(update)
			if err := conn.WriteMessage(websocket.TextMessage, updateJSON); err != nil {
				break
			}
		default:
		}
	}
}

func createAgentGraph(config AgentConfig) (*langgraph.Graph, error) {
	var startNode langgraph.NodeID = "start"
	var processNode langgraph.NodeID = "process"
	var decisionNode langgraph.NodeID = "decision"
	var toolNode langgraph.NodeID = "tool"
	var finishNode langgraph.NodeID = "finish"
	
	graph := langgraph.NewGraph(startNode)
	
	graph.DefineStateSchema(map[string]string{
		"input":       "string",
		"processed":   "string",
		"decision":    "string",
		"tool_output": "string",
		"output":      "string",
		"agent_type":  "string",
		"model_id":    "string",
	})
	
	graph.AddNode(startNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		state["agent_type"] = config.AgentType
		state["model_id"] = config.ModelID
		return state, nil
	}, map[string]interface{}{
		"description": "Initialize the agent workflow",
	})
	
	graph.AddNode(processNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		input, _ := state["input"].(string)
		processed := fmt.Sprintf("Processed: %s", input)
		state["processed"] = processed
		return state, nil
	}, map[string]interface{}{
		"description": "Process the input",
	})
	
	graph.AddNode(decisionNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		processed, _ := state["processed"].(string)
		
		decision := "use_tool"
		if len(processed) > 20 {
			decision = "finish"
		}
		
		state["decision"] = decision
		return state, nil
	}, map[string]interface{}{
		"description": "Make a decision based on the processed input",
	})
	
	graph.AddNode(toolNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		processed, _ := state["processed"].(string)
		
		toolOutput := fmt.Sprintf("Tool output for: %s", processed)
		state["tool_output"] = toolOutput
		return state, nil
	}, map[string]interface{}{
		"description": "Execute a tool",
	})
	
	graph.AddNode(finishNode, func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
		processed, _ := state["processed"].(string)
		toolOutput, _ := state["tool_output"].(string)
		
		var output string
		if toolOutput != "" {
			output = fmt.Sprintf("Final output: %s with tool: %s", processed, toolOutput)
		} else {
			output = fmt.Sprintf("Final output: %s", processed)
		}
		
		state["output"] = output
		return state, nil
	}, map[string]interface{}{
		"description": "Finalize the agent workflow",
	})
	
	graph.AddEdge(startNode, processNode, langgraph.DefaultEdge(processNode))
	graph.AddEdge(processNode, decisionNode, langgraph.DefaultEdge(decisionNode))
	
	graph.AddEdge(decisionNode, toolNode, langgraph.ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			decision, _ := state["decision"].(string)
			return decision == "use_tool"
		},
		toolNode,
	))
	
	graph.AddEdge(decisionNode, finishNode, langgraph.ConditionalEdge(
		func(ctx context.Context, state map[string]interface{}) bool {
			decision, _ := state["decision"].(string)
			return decision == "finish"
		},
		finishNode,
	))
	
	graph.AddEdge(toolNode, finishNode, langgraph.DefaultEdge(finishNode))
	
	graph.SetExitPoint(finishNode)
	
	err := graph.Compile()
	if err != nil {
		return nil, fmt.Errorf("failed to compile graph: %w", err)
	}
	
	return graph, nil
}

func getGraphStateObject(graphID string, graph *langgraph.Graph, state map[string]interface{}, status string) GraphState {
	nodes := []NodeState{
		{ID: "start", Type: "agent", Label: "Start", Status: "completed"},
		{ID: "process", Type: "agent", Label: "Process", Status: "active"},
		{ID: "decision", Type: "agent", Label: "Decision", Status: "idle"},
		{ID: "tool", Type: "tool", Label: "Tool", Status: "idle"},
		{ID: "finish", Type: "agent", Label: "Finish", Status: "idle"},
	}
	
	edges := []EdgeState{
		{Source: "start", Target: "process"},
		{Source: "process", Target: "decision"},
		{Source: "decision", Target: "tool", Label: "use_tool"},
		{Source: "decision", Target: "finish", Label: "finish"},
		{Source: "tool", Target: "finish"},
	}
	
	return GraphState{
		GraphID:     graphID,
		Nodes:       nodes,
		Edges:       edges,
		CurrentNode: "process", // This would be determined from the actual graph state
		State:       state,
		Status:      status,
	}
}
