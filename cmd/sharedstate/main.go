package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	mockStream := sharedstate.NewMockStream()
	mockStateManager := sharedstate.NewMockStateManager()
	mockSupabaseClient := sharedstate.NewMockSupabaseClient()

	wsManager := sharedstate.NewWebSocketManager()
	ctx, cancel := context.WithCancel(context.Background())
	go wsManager.Start(ctx)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			clientID = "client-" + fmt.Sprintf("%d", time.Now().UnixNano())
		}
		
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		}
		
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}
		
		log.Printf("Client connected: %s", clientID)
		
		go handleWebSocketConnection(conn, clientID, wsManager)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
		os.Exit(0)
	}()

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting shared state test server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleWebSocketConnection(conn *websocket.Conn, clientID string, wsManager *sharedstate.WebSocketManager) {
	defer conn.Close()
	
	welcomeMsg := map[string]interface{}{
		"type": "welcome",
		"data": map[string]interface{}{
			"client_id": clientID,
			"message":   "Welcome to the shared state system",
		},
		"timestamp": time.Now().UTC(),
	}
	
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}
	
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		log.Printf("Received message from client %s: %v", clientID, msg)
		
		if msgType, ok := msg["type"].(string); ok {
			switch msgType {
			case "subscribe":
				if topic, ok := msg["data"].(string); ok {
					wsManager.SubscribeToTopic(clientID, topic)
					log.Printf("Client %s subscribed to topic: %s", clientID, topic)
				}
			case "unsubscribe":
				if topic, ok := msg["data"].(string); ok {
					wsManager.UnsubscribeFromTopic(clientID, topic)
					log.Printf("Client %s unsubscribed from topic: %s", clientID, topic)
				}
			case "event":
				response := map[string]interface{}{
					"type": "event_received",
					"data": msg["data"],
					"timestamp": time.Now().UTC(),
				}
				if err := conn.WriteJSON(response); err != nil {
					log.Printf("Error sending event response: %v", err)
				}
			}
		}
	}
	
	log.Printf("Client disconnected: %s", clientID)
}
