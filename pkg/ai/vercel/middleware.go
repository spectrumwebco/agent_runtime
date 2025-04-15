package vercel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StreamingMiddleware(aiClient *VercelAIClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Accept") == "text/event-stream" {
			handleSSE(c, aiClient)
			return
		}
		c.Next()
	}
}

func ActionMiddleware(adapter *ServerActionsAdapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.URL.Path == "/api/actions" {
			handleAction(c, adapter)
			return
		}
		c.Next()
	}
}

func handleSSE(c *gin.Context, aiClient *VercelAIClient) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	events := make(chan interface{}, 10)
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = fmt.Sprintf("client-%d", time.Now().UnixNano())
	}

	aiClient.RegisterStreamingHandler(clientID, events)
	defer aiClient.UnregisterStreamingHandler(clientID)

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	go func() {
		<-ctx.Done()
		aiClient.UnregisterStreamingHandler(clientID)
		close(events)
	}()

	sendSSEEvent(c.Writer, "connected", map[string]interface{}{
		"client_id": clientID,
		"timestamp": time.Now().Unix(),
	})

	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			sendSSEEvent(c.Writer, "update", event)
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			sendSSEEvent(c.Writer, "keepalive", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			})
		}
	}
}

func handleAction(c *gin.Context, adapter *ServerActionsAdapter) {
	var request struct {
		Action string                 `json:"action"`
		Params map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if request.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action name is required"})
		return
	}

	result, err := adapter.ExecuteAction(c.Request.Context(), request.Action, request.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"action": request.Action,
		"result": result,
	})
}

func sendSSEEvent(w http.ResponseWriter, eventType string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", eventType)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.(http.Flusher).Flush()
}
