package vercel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type RSCMiddleware struct {
	rscAdapter        *RSCAdapter
	sharedStateManager *sharedstate.SharedStateManager
	eventStream       *eventstream.Stream
}

func NewRSCMiddleware(rscAdapter *RSCAdapter, sharedStateManager *sharedstate.SharedStateManager, eventStream *eventstream.Stream) *RSCMiddleware {
	return &RSCMiddleware{
		rscAdapter:        rscAdapter,
		sharedStateManager: sharedStateManager,
		eventStream:       eventStream,
	}
}

func (m *RSCMiddleware) HandleRSCRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Accept") == "application/rsc+json" {
			m.handleRSCResponse(c)
			return
		}

		c.Next()
	}
}

func (m *RSCMiddleware) handleRSCResponse(c *gin.Context) {
	var request struct {
		ComponentType string                 `json:"componentType"`
		Props         map[string]interface{} `json:"props"`
		AgentID       string                 `json:"agentID"`
		ActionID      string                 `json:"actionID"`
		ActionType    string                 `json:"actionType"`
		ActionData    map[string]interface{} `json:"actionData"`
		ToolID        string                 `json:"toolID"`
		ToolName      string                 `json:"toolName"`
		ToolInput     map[string]interface{} `json:"toolInput"`
		ToolOutput    map[string]interface{} `json:"toolOutput"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	ctx := context.Background()
	var component GeneratedComponent
	var err error

	if request.ComponentType != "" {
		componentType := ComponentType(request.ComponentType)
		component, err = m.rscAdapter.GenerateComponent(ctx, componentType, request.Props)
	} else if request.AgentID != "" && request.ActionType != "" {
		component, err = m.rscAdapter.GenerateComponentFromAgentAction(
			ctx,
			request.AgentID,
			request.ActionID,
			request.ActionType,
			request.ActionData,
		)
	} else if request.AgentID != "" && request.ToolName != "" {
		component, err = m.rscAdapter.GenerateComponentFromToolUsage(
			ctx,
			request.AgentID,
			request.ToolID,
			request.ToolName,
			request.ToolInput,
			request.ToolOutput,
		)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: missing required fields"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error generating component: %v", err)})
		return
	}

	if component.ID == "" {
		component.ID = uuid.New().String()
	}
	
	now := time.Now().Unix()
	if component.CreatedAt == 0 {
		component.CreatedAt = now
	}
	
	component.UpdatedAt = now

	if request.AgentID != "" {
		stateID := fmt.Sprintf("component:%s", component.ID)
		componentJSON, err := json.Marshal(component)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error marshaling component: %v", err)})
			return
		}

		var componentData map[string]interface{}
		err = json.Unmarshal(componentJSON, &componentData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error unmarshaling component: %v", err)})
			return
		}

		err = m.sharedStateManager.UpdateState(sharedstate.ComponentState, stateID, componentData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating state: %v", err)})
			return
		}

		component.StateID = stateID
	}

	c.Header("Content-Type", "application/rsc+json")
	c.JSON(http.StatusOK, component)
}

func (m *RSCMiddleware) StreamRSCComponents() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Accept") == "text/event-stream" {
			m.handleRSCStream(c)
			return
		}

		c.Next()
	}
}

func (m *RSCMiddleware) handleRSCStream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	events := make(chan interface{}, 10)
	defer close(events)

	handlerID := uuid.New().String()
	if m.rscAdapter.aiClient != nil {
		m.rscAdapter.aiClient.RegisterStreamingHandler(handlerID, events)
		defer m.rscAdapter.aiClient.UnregisterStreamingHandler(handlerID)
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	components := m.rscAdapter.GetComponentHistory()
	for _, component := range components {
		componentJSON, err := json.Marshal(component)
		if err != nil {
			continue
		}

		c.SSEvent("component", string(componentJSON))
		c.Writer.Flush()
	}

	c.SSEvent("ping", time.Now().Unix())
	c.Writer.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.SSEvent("ping", time.Now().Unix())
			c.Writer.Flush()
		case event := <-events:
			eventJSON, err := json.Marshal(event)
			if err != nil {
				continue
			}

			c.SSEvent("event", string(eventJSON))
			c.Writer.Flush()
		}
	}
}
