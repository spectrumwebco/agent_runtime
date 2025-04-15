package socketio

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
)

type DjangoMiddleware struct {
	socketIOServer *Server
	djangoBaseURL  string
	eventStream    *eventstream.EventStream
}

func NewDjangoMiddleware(socketIOServer *Server, djangoBaseURL string, eventStream *eventstream.EventStream) *DjangoMiddleware {
	return &DjangoMiddleware{
		socketIOServer: socketIOServer,
		djangoBaseURL:  djangoBaseURL,
		eventStream:    eventStream,
	}
}

func (m *DjangoMiddleware) RegisterRoutes(router *gin.Engine) {
	router.GET("/socket.io/", func(c *gin.Context) {
		m.socketIOServer.HandleConnection(c)
	})

	router.Any("/api/*path", m.proxyToDjango)

	m.eventStream.Subscribe("django_event", m.handleDjangoEvent)
}

func (m *DjangoMiddleware) proxyToDjango(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	djangoURL := m.djangoBaseURL + "/api" + path

	req, err := http.NewRequest(c.Request.Method, djangoURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.URL.RawQuery = c.Request.URL.RawQuery

	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
			return
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request to Django"})
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response body"})
		return
	}

	c.Status(resp.StatusCode)

	c.Writer.Write(bodyBytes)
}

func (m *DjangoMiddleware) handleDjangoEvent(eventType string, data map[string]interface{}) {
	eventData, ok := data["event"].(map[string]interface{})
	if !ok {
		log.Printf("Invalid event data format")
		return
	}

	conversationID, ok := data["conversation_id"].(string)
	if !ok {
		log.Printf("Missing conversation ID")
		return
	}

	switch eventType {
	case "agent_message":
		message, ok := eventData["message"].(string)
		if !ok {
			log.Printf("Missing message")
			return
		}
		m.socketIOServer.SendAgentMessage(conversationID, message)

	case "agent_observation":
		observation, ok := eventData["observation"].(string)
		if !ok {
			log.Printf("Missing observation")
			return
		}
		observationType, ok := eventData["observation_type"].(string)
		if !ok {
			observationType = "info"
		}
		m.socketIOServer.SendAgentObservation(conversationID, observation, observationType)

	case "agent_state_update":
		state, ok := eventData["state"].(string)
		if !ok {
			log.Printf("Missing state")
			return
		}
		m.socketIOServer.SendAgentStateUpdate(conversationID, state)

	default:
		log.Printf("Unknown event type: %s", eventType)
	}
}

func (m *DjangoMiddleware) SendEventToDjango(conversationID string, eventType string, data map[string]interface{}) error {
	djangoURL := m.djangoBaseURL + "/api/events/"

	requestBody := map[string]interface{}{
		"conversation_id": conversationID,
		"event_type":      eventType,
		"data":            data,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", djangoURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func RegisterSocketIOServer(router *gin.Engine, eventStream *eventstream.EventStream, stateManager *statemanager.StateManager, djangoBaseURL string) {
	socketIOServer := NewServer(eventStream, stateManager)

	djangoMiddleware := NewDjangoMiddleware(socketIOServer, djangoBaseURL, eventStream)

	djangoMiddleware.RegisterRoutes(router)
}
