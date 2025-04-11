package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketMessage represents a message sent over WebSocket connection
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Topic     string      `json:"topic,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}
// Connection represents a WebSocket client connection

// Connection represents a WebSocket client connection
type Connection struct {
	conn      *websocket.Conn
	send      chan WebSocketMessage
	topics    map[string]bool
	clientID  string
	closed    bool
	closeLock sync.RWMutex
// WebSocketManager manages WebSocket connections and topic subscriptions
}

// WebSocketManager manages WebSocket connections and topic subscriptions
type WebSocketManager struct {
	upgrader    websocket.Upgrader
	connections map[string]*Connection
	topicConns  map[string]map[string]*Connection
// NewWebSocketManager creates a new WebSocketManager instance
	mutex       sync.RWMutex
}

// NewWebSocketManager creates a new WebSocketManager instance
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		connections: make(map[string]*Connection),
// SubscribeToTopic subscribes a client to a specific topic
		topicConns:  make(map[string]map[string]*Connection),
	}
}

// SubscribeToTopic subscribes a client to a specific topic
func (m *WebSocketManager) SubscribeToTopic(clientID, topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	conn, exists := m.connections[clientID]
	if !exists {
		return
	}

	conn.topics[topic] = true

	if _, exists = m.topicConns[topic]; !exists {
		m.topicConns[topic] = make(map[string]*Connection)
	}
// UnsubscribeFromTopic unsubscribes a client from a specific topic
	m.topicConns[topic][clientID] = conn

	log.Printf("Client %s subscribed to topic: %s\n", clientID, topic)
}

// UnsubscribeFromTopic unsubscribes a client from a specific topic
func (m *WebSocketManager) UnsubscribeFromTopic(clientID, topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	conn, exists := m.connections[clientID]
	if !exists {
		return
	}

	delete(conn.topics, topic)

	if conns, exists := m.topicConns[topic]; exists {
		delete(conns, clientID)
		if len(conns) == 0 {
			delete(m.topicConns, topic)
// PublishToTopic publishes a message to all clients subscribed to a topic
		}
	}

	log.Printf("Client %s unsubscribed from topic: %s\n", clientID, topic)
}

// PublishToTopic publishes a message to all clients subscribed to a topic
func (m *WebSocketManager) PublishToTopic(topic string, data interface{}) {
	message := WebSocketMessage{
		Type:      "message",
		Topic:     topic,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	m.mutex.RLock()
	if conns, exists := m.topicConns[topic]; exists {
		for _, conn := range conns {
			select {
			case conn.send <- message:
			default:
				conn.closeLock.Lock()
				if !conn.closed {
					close(conn.send)
					conn.closed = true
				}
// HandleConnection handles new WebSocket connection requests
				conn.closeLock.Unlock()
			}
		}
	}
	m.mutex.RUnlock()
}

// HandleConnection handles new WebSocket connection requests
func (m *WebSocketManager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		clientID = fmt.Sprintf("client-%d", time.Now().UnixNano())
	}

	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Connection{
		conn:     conn,
		send:     make(chan WebSocketMessage, 256),
		topics:   make(map[string]bool),
		clientID: clientID,
		closed:   false,
	}

	m.mutex.Lock()
	m.connections[clientID] = client
	m.mutex.Unlock()

	log.Printf("Client connected: %s", clientID)

	go m.readPump(client)
	go m.writePump(client)
}

func (m *WebSocketManager) readPump(c *Connection) {
	defer func() {
		m.mutex.Lock()
		if _, ok := m.connections[c.clientID]; ok {
			delete(m.connections, c.clientID)
			for topic := range c.topics {
				if conns, exists := m.topicConns[topic]; exists {
					delete(conns, c.clientID)
					if len(conns) == 0 {
						delete(m.topicConns, topic)
					}
				}
			}
		}
		m.mutex.Unlock()

		c.closeLock.Lock()
		if !c.closed {
			c.conn.Close()
			c.closed = true
		}
		c.closeLock.Unlock()

		log.Printf("Client disconnected: %s", c.clientID)
	}()

	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	welcomeMsg := WebSocketMessage{
		Type: "welcome",
		Data: map[string]interface{}{
			"client_id": c.clientID,
			"message":   "Welcome to the shared state system",
		},
		Timestamp: time.Now().UTC(),
	}
	c.send <- welcomeMsg

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		log.Printf("Received message from client %s: %v", c.clientID, msg)

		switch msg.Type {
		case "subscribe":
			if topic, ok := msg.Data.(string); ok {
				m.SubscribeToTopic(c.clientID, topic)

				c.send <- WebSocketMessage{
					Type:      "subscribed",
					Data:      topic,
					Timestamp: time.Now().UTC(),
				}
			}
		case "unsubscribe":
			if topic, ok := msg.Data.(string); ok {
				m.UnsubscribeFromTopic(c.clientID, topic)

				c.send <- WebSocketMessage{
					Type:      "unsubscribed",
					Data:      topic,
					Timestamp: time.Now().UTC(),
				}
			}
		case "event":
			c.send <- WebSocketMessage{
				Type:      "event_received",
				Data:      msg.Data,
				Timestamp: time.Now().UTC(),
			}

			if eventData, ok := msg.Data.(map[string]interface{}); ok {
				if eventType, ok := eventData["type"].(string); ok && eventType == "state_update" {
					if stateType, ok := eventData["state_type"].(string); ok {
						if stateID, ok := eventData["state_id"].(string); ok {
							if data, ok := eventData["data"].(map[string]interface{}); ok {
								topic := fmt.Sprintf("state:%s:%s", stateType, stateID)
								m.PublishToTopic(topic, data)
							}
						}
					}
				}
			}
		}
	}
}

func (m *WebSocketManager) writePump(c *Connection) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.closeLock.Lock()
		if !c.closed {
			c.conn.Close()
			c.closed = true
		}
		c.closeLock.Unlock()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	wsManager := NewWebSocketManager()

	http.HandleFunc("/ws", wsManager.HandleConnection)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.HandleFunc("/api/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Topic string      `json:"topic"`
			Data  interface{} `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		wsManager.PublishToTopic(req.Topic, req.Data)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		stateType := r.URL.Query().Get("type")
		stateID := r.URL.Query().Get("id")

		if stateType == "" || stateID == "" {
			http.Error(w, "Missing type or id parameter", http.StatusBadRequest)
			return
		}

		state := map[string]interface{}{
			"state_type": stateType,
			"state_id":   stateID,
			"status":     "mock",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"data": map[string]interface{}{
				"value": fmt.Sprintf("Mock state for %s:%s", stateType, stateID),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting shared state test server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
