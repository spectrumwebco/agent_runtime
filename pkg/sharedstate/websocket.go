package sharedstate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
)

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Topic     string      `json:"topic"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type Connection struct {
	conn      *websocket.Conn
	send      chan WebSocketMessage
	topics    map[string]bool
	manager   *WebSocketManager
	clientID  string
	closed    bool
	closeLock sync.RWMutex
}

type WebSocketManager struct {
	upgrader    websocket.Upgrader
	connections map[string]*Connection
	register    chan *Connection
	unregister  chan *Connection
	broadcast   chan WebSocketMessage
	topicConns  map[string]map[string]*Connection
	mutex       sync.RWMutex
}

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
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		broadcast:   make(chan WebSocketMessage),
		topicConns:  make(map[string]map[string]*Connection),
	}
}

func (m *WebSocketManager) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-m.register:
			m.mutex.Lock()
			m.connections[conn.clientID] = conn
			m.mutex.Unlock()
			log.Printf("WebSocket client registered: %s\n", conn.clientID)

		case conn := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.connections[conn.clientID]; ok {
				delete(m.connections, conn.clientID)
				for topic := range conn.topics {
					if conns, exists := m.topicConns[topic]; exists {
						delete(conns, conn.clientID)
						if len(conns) == 0 {
							delete(m.topicConns, topic)
						}
					}
				}
			}
			m.mutex.Unlock()
			log.Printf("WebSocket client unregistered: %s\n", conn.clientID)

		case message := <-m.broadcast:
			topic := message.Topic
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
						conn.closeLock.Unlock()
					}
				}
			}
			m.mutex.RUnlock()
		}
	}
}

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
	m.topicConns[topic][clientID] = conn
	
	log.Printf("Client %s subscribed to topic: %s\n", clientID, topic)
}

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
		}
	}
	
	log.Printf("Client %s unsubscribed from topic: %s\n", clientID, topic)
}

func (m *WebSocketManager) PublishToTopic(topic string, data interface{}) {
	message := WebSocketMessage{
		Type:      "message",
		Topic:     topic,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	m.broadcast <- message
}

func (m *WebSocketManager) HandleWebSocket(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = fmt.Sprintf("client-%d", time.Now().UnixNano())
	}

	conn, err := m.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	client := &Connection{
		conn:     conn,
		send:     make(chan WebSocketMessage, 256),
		topics:   make(map[string]bool),
		manager:  m,
		clientID: clientID,
		closed:   false,
	}

	m.register <- client

	go client.readPump()
	go client.writePump()
}

func (c *Connection) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.closeLock.Lock()
		if !c.closed {
			c.conn.Close()
			c.closed = true
		}
		c.closeLock.Unlock()
	}()

	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

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

		switch msg.Type {
		case "subscribe":
			if topic, ok := msg.Data.(string); ok {
				c.manager.SubscribeToTopic(c.clientID, topic)
			}
		case "unsubscribe":
			if topic, ok := msg.Data.(string); ok {
				c.manager.UnsubscribeFromTopic(c.clientID, topic)
			}
		case "event":
			if eventData, ok := msg.Data.(map[string]interface{}); ok {
				handleClientEvent(eventData)
			}
		}
	}
}

func (c *Connection) writePump() {
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshalling message: %v", err)
				return
			}
			
			w.Write(messageJSON)

			if err := w.Close(); err != nil {
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

var handleClientEvent = func(eventData map[string]interface{}) {
	log.Printf("Received client event: %v\n", eventData)
}
