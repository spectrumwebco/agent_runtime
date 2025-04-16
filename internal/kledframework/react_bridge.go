package kledframework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ReactBridge struct {
	config     *ReactBridgeConfig
	clients    map[string]*websocketClient
	clientsMu  sync.RWMutex
	eventChan  chan map[string]interface{}
	stateMu    sync.RWMutex
	state      map[string]interface{}
	httpServer *http.Server
	upgrader   websocket.Upgrader
	isRunning  bool
	stopCh     chan struct{}
}

type ReactBridgeConfig struct {
	Address         string
	Path            string
	AllowedOrigins  []string
	HeartbeatPeriod time.Duration
	StateUpdatePath string
	StateGetPath    string
}

type websocketClient struct {
	id       string
	conn     *websocket.Conn
	send     chan []byte
	reactBridge *ReactBridge
}

func DefaultReactBridgeConfig() *ReactBridgeConfig {
	return &ReactBridgeConfig{
		Address:         ":8081",
		Path:            "/ws",
		AllowedOrigins:  []string{"*"},
		HeartbeatPeriod: 30 * time.Second,
		StateUpdatePath: "/api/state",
		StateGetPath:    "/api/state",
	}
}

func NewReactBridge(config *ReactBridgeConfig) (*ReactBridge, error) {
	if config == nil {
		config = DefaultReactBridgeConfig()
	}

	bridge := &ReactBridge{
		config:    config,
		clients:   make(map[string]*websocketClient),
		eventChan: make(chan map[string]interface{}, 100),
		state:     make(map[string]interface{}),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				if len(config.AllowedOrigins) == 0 || (len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*") {
					return true
				}
				origin := r.Header.Get("Origin")
				for _, allowed := range config.AllowedOrigins {
					if allowed == origin {
						return true
					}
				}
				return false
			},
		},
		stopCh: make(chan struct{}),
	}

	return bridge, nil
}

func (rb *ReactBridge) Start(ctx context.Context) error {
	rb.stateMu.Lock()
	rb.isRunning = true
	rb.stateMu.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc(rb.config.Path, rb.handleWebSocket)
	mux.HandleFunc(rb.config.StateUpdatePath, rb.handleStateUpdate)
	mux.HandleFunc(rb.config.StateGetPath, rb.handleStateGet)

	rb.httpServer = &http.Server{
		Addr:    rb.config.Address,
		Handler: mux,
	}

	go rb.run()

	go func() {
		log.Printf("Starting React bridge on %s", rb.config.Address)
		if err := rb.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("React bridge server error: %v", err)
		}
	}()

	return nil
}

func (rb *ReactBridge) Stop() error {
	rb.stateMu.Lock()
	if !rb.isRunning {
		rb.stateMu.Unlock()
		return nil
	}
	rb.isRunning = false
	rb.stateMu.Unlock()

	close(rb.stopCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rb.clientsMu.Lock()
	for _, client := range rb.clients {
		client.conn.Close()
	}
	rb.clients = make(map[string]*websocketClient)
	rb.clientsMu.Unlock()

	return rb.httpServer.Shutdown(ctx)
}

func (rb *ReactBridge) SendEvent(ctx context.Context, event map[string]interface{}) error {
	select {
	case rb.eventChan <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("event channel full")
	}
}

func (rb *ReactBridge) SubscribeToEvents(ctx context.Context) (<-chan map[string]interface{}, error) {
	events := make(chan map[string]interface{}, 100)

	go func() {
		defer close(events)

		for {
			select {
			case <-ctx.Done():
				return
			case <-rb.stopCh:
				return
			}
		}
	}()

	return events, nil
}

func (rb *ReactBridge) UpdateState(ctx context.Context, state map[string]interface{}) error {
	rb.stateMu.Lock()
	defer rb.stateMu.Unlock()

	for k, v := range state {
		rb.state[k] = v
	}

	data, err := json.Marshal(map[string]interface{}{
		"type":  "state_update",
		"state": rb.state,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal state update: %w", err)
	}

	rb.clientsMu.RLock()
	defer rb.clientsMu.RUnlock()

	for _, client := range rb.clients {
		select {
		case client.send <- data:
		default:
			log.Printf("Failed to send state update to client %s: channel full", client.id)
		}
	}

	return nil
}

func (rb *ReactBridge) GetState(ctx context.Context) (map[string]interface{}, error) {
	rb.stateMu.RLock()
	defer rb.stateMu.RUnlock()

	stateCopy := make(map[string]interface{}, len(rb.state))
	for k, v := range rb.state {
		stateCopy[k] = v
	}

	return stateCopy, nil
}

func (rb *ReactBridge) Close() error {
	return rb.Stop()
}

func (rb *ReactBridge) run() {
	ticker := time.NewTicker(rb.config.HeartbeatPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-rb.stopCh:
			return
		case event := <-rb.eventChan:
			rb.broadcastEvent(event)
		case <-ticker.C:
			rb.sendHeartbeat()
		}
	}
}

func (rb *ReactBridge) broadcastEvent(event map[string]interface{}) {
	data, err := json.Marshal(map[string]interface{}{
		"type":  "event",
		"event": event,
	})
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	rb.clientsMu.RLock()
	defer rb.clientsMu.RUnlock()

	for _, client := range rb.clients {
		select {
		case client.send <- data:
		default:
			log.Printf("Failed to send event to client %s: channel full", client.id)
		}
	}
}

func (rb *ReactBridge) sendHeartbeat() {
	data, err := json.Marshal(map[string]interface{}{
		"type": "heartbeat",
		"time": time.Now().Unix(),
	})
	if err != nil {
		log.Printf("Failed to marshal heartbeat: %v", err)
		return
	}

	rb.clientsMu.RLock()
	defer rb.clientsMu.RUnlock()

	for _, client := range rb.clients {
		select {
		case client.send <- data:
		default:
			go rb.removeClient(client)
		}
	}
}

func (rb *ReactBridge) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := rb.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &websocketClient{
		id:       fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().UnixNano()),
		conn:     conn,
		send:     make(chan []byte, 256),
		reactBridge: rb,
	}

	rb.clientsMu.Lock()
	rb.clients[client.id] = client
	rb.clientsMu.Unlock()

	go client.writePump()
	go client.readPump()

	rb.stateMu.RLock()
	stateCopy := make(map[string]interface{}, len(rb.state))
	for k, v := range rb.state {
		stateCopy[k] = v
	}
	rb.stateMu.RUnlock()

	data, err := json.Marshal(map[string]interface{}{
		"type":  "state_init",
		"state": stateCopy,
	})
	if err != nil {
		log.Printf("Failed to marshal initial state: %v", err)
	} else {
		client.send <- data
	}
}

func (rb *ReactBridge) handleStateUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var state map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode state: %v", err), http.StatusBadRequest)
		return
	}

	if err := rb.UpdateState(r.Context(), state); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update state: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rb *ReactBridge) handleStateGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state, err := rb.GetState(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get state: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		log.Printf("Failed to encode state: %v", err)
	}
}

func (rb *ReactBridge) removeClient(client *websocketClient) {
	rb.clientsMu.Lock()
	defer rb.clientsMu.Unlock()

	if _, ok := rb.clients[client.id]; ok {
		delete(rb.clients, client.id)
		close(client.send)
	}
}

func (c *websocketClient) writePump() {
	ticker := time.NewTicker(time.Second * 54)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
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

func (c *websocketClient) readPump() {
	defer func() {
		c.reactBridge.removeClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var event map[string]interface{}
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Failed to unmarshal client message: %v", err)
			continue
		}

		if eventType, ok := event["type"].(string); ok {
			switch eventType {
			case "state_update":
				if state, ok := event["state"].(map[string]interface{}); ok {
					c.reactBridge.UpdateState(context.Background(), state)
				}
			case "event":
				if eventData, ok := event["event"].(map[string]interface{}); ok {
					c.reactBridge.eventChan <- eventData
				}
			}
		}
	}
}
