package socketio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
)

type Server struct {
	clients       map[string]*Client
	clientsMutex  sync.RWMutex
	upgrader      websocket.Upgrader
	eventStream   *eventstream.EventStream
	stateManager  *statemanager.StateManager
	eventHandlers map[string]EventHandler
}

type Client struct {
	ID             string
	Conn           *websocket.Conn
	ConversationID string
	Send           chan []byte
	Server         *Server
}

type Event struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

type EventHandler func(client *Client, data json.RawMessage) error

func NewServer(eventStream *eventstream.EventStream, stateManager *statemanager.StateManager) *Server {
	server := &Server{
		clients: make(map[string]*Client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		eventStream:   eventStream,
		stateManager:  stateManager,
		eventHandlers: make(map[string]EventHandler),
	}

	server.RegisterEventHandler("oh_user_action", server.handleUserAction)
	server.RegisterEventHandler("oh_agent_state_change", server.handleAgentStateChange)
	server.RegisterEventHandler("oh_message", server.handleMessage)

	return server
}

func (s *Server) RegisterEventHandler(eventType string, handler EventHandler) {
	s.eventHandlers[eventType] = handler
}

func (s *Server) HandleConnection(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		conversationID = "default"
	}

	client := &Client{
		ID:             uuid.New().String(),
		Conn:           conn,
		ConversationID: conversationID,
		Send:           make(chan []byte, 256),
		Server:         s,
	}

	s.clientsMutex.Lock()
	s.clients[client.ID] = client
	s.clientsMutex.Unlock()

	go client.readPump()
	go client.writePump()

	connectionEvent := map[string]interface{}{
		"type":            "connection_established",
		"client_id":       client.ID,
		"conversation_id": client.ConversationID,
		"timestamp":       time.Now().UnixMilli(),
	}
	eventJSON, _ := json.Marshal(connectionEvent)
	client.Send <- formatSocketIOMessage("oh_event", eventJSON)

	log.Printf("Client connected: %s (conversation: %s)", client.ID, client.ConversationID)
}

func (s *Server) BroadcastToConversation(conversationID string, eventType string, data interface{}) {
	eventJSON, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	message := formatSocketIOMessage(eventType, eventJSON)

	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	for _, client := range s.clients {
		if client.ConversationID == conversationID {
			select {
			case client.Send <- message:
			default:
				s.removeClient(client)
			}
		}
	}
}

func (s *Server) removeClient(client *Client) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if _, ok := s.clients[client.ID]; ok {
		delete(s.clients, client.ID)
		close(client.Send)
	}
}

func (s *Server) handleUserAction(client *Client, data json.RawMessage) error {
	var action struct {
		Type    string          `json:"type"`
		Message json.RawMessage `json:"message,omitempty"`
		State   string          `json:"state,omitempty"`
	}

	if err := json.Unmarshal(data, &action); err != nil {
		return err
	}

	switch action.Type {
	case "message":
		return s.handleUserMessage(client, action.Message)
	case "agent_state_change":
		return s.handleUserAgentStateChange(client, action.State)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (s *Server) handleUserMessage(client *Client, messageData json.RawMessage) error {
	var message struct {
		Content string `json:"content"`
	}

	if err := json.Unmarshal(messageData, &message); err != nil {
		return err
	}

	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"source":    "user",
		"type":      "message",
		"message":   message.Content,
		"timestamp": time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(client.ConversationID, "oh_event", event)

	s.eventStream.Publish("user_message", map[string]interface{}{
		"conversation_id": client.ConversationID,
		"client_id":       client.ID,
		"message":         message.Content,
	})

	return nil
}

func (s *Server) handleUserAgentStateChange(client *Client, state string) error {
	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"source":    "user",
		"type":      "agent_state_change",
		"state":     state,
		"timestamp": time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(client.ConversationID, "oh_event", event)

	s.eventStream.Publish("agent_state_change", map[string]interface{}{
		"conversation_id": client.ConversationID,
		"client_id":       client.ID,
		"state":           state,
	})

	return nil
}

func (s *Server) handleAgentStateChange(client *Client, data json.RawMessage) error {
	var stateChange struct {
		State string `json:"state"`
	}

	if err := json.Unmarshal(data, &stateChange); err != nil {
		return err
	}

	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"source":    "agent",
		"type":      "agent_state_update",
		"state":     stateChange.State,
		"timestamp": time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(client.ConversationID, "oh_event", event)

	return nil
}

func (s *Server) handleMessage(client *Client, data json.RawMessage) error {
	var message struct {
		Content string `json:"content"`
	}

	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}

	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"source":    "agent",
		"type":      "message",
		"message":   message.Content,
		"timestamp": time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(client.ConversationID, "oh_event", event)

	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.Server.removeClient(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512KB
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		if len(message) > 0 {
			msgStr := string(message)

			if msgStr == "2" {
				c.Send <- []byte("3") // pong
				continue
			}

			if len(msgStr) > 2 && msgStr[:2] == "42" {
				var event []interface{}
				if err := json.Unmarshal([]byte(msgStr[2:]), &event); err != nil {
					log.Printf("Error parsing Socket.IO event: %v", err)
					continue
				}

				if len(event) < 2 {
					log.Printf("Invalid Socket.IO event format")
					continue
				}

				eventName, ok := event[0].(string)
				if !ok {
					log.Printf("Invalid event name type")
					continue
				}

				eventDataJSON, err := json.Marshal(event[1])
				if err != nil {
					log.Printf("Error marshaling event data: %v", err)
					continue
				}

				if handler, ok := c.Server.eventHandlers[eventName]; ok {
					if err := handler(c, eventDataJSON); err != nil {
						log.Printf("Error handling event %s: %v", eventName, err)
					}
				} else {
					log.Printf("Unknown event type: %s", eventName)
				}
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func formatSocketIOMessage(eventName string, data []byte) []byte {
	return []byte(fmt.Sprintf("42[\"%s\",%s]", eventName, string(data)))
}

func (s *Server) SendAgentMessage(conversationID string, message string) {
	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"source":    "agent",
		"type":      "message",
		"message":   message,
		"timestamp": time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(conversationID, "oh_event", event)
}

func (s *Server) SendAgentObservation(conversationID string, observation string, observationType string) {
	event := map[string]interface{}{
		"id":          uuid.New().String(),
		"source":      "agent",
		"type":        "observation",
		"observation": observationType,
		"message":     observation,
		"timestamp":   time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(conversationID, "oh_event", event)
}

func (s *Server) SendAgentStateUpdate(conversationID string, state string) {
	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"source":    "agent",
		"type":      "agent_state_update",
		"state":     state,
		"timestamp": time.Now().UnixMilli(),
	}

	s.BroadcastToConversation(conversationID, "oh_event", event)
}
