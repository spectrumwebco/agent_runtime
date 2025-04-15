package sharedstate

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
	"github.com/spectrumwebco/agent_runtime/pkg/ai/vercel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*Server, *eventstream.Stream, *statemanager.StateManager) {
	gin.SetMode(gin.TestMode)
	
	eventStream := eventstream.NewStream()
	
	stateManager := statemanager.NewInMemoryStateManager()
	
	sharedStateCfg := SharedStateConfig{
		EventStream:  eventStream,
		StateManager: stateManager,
	}
	sharedStateManager, err := NewSharedStateManager(sharedStateCfg)
	require.NoError(t, err)
	
	aiClient := vercel.NewVercelAIClient(sharedStateManager, eventStream)
	actionAdapter := vercel.NewServerActionsAdapter(sharedStateManager, aiClient)
	
	actionAdapter.RegisterAction("test_action", func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{
			"result": "success",
			"params": params,
		}, nil
	})
	
	serverCfg := ServerConfig{
		Port:               8080,
		EventStream:        eventStream,
		StateManager:       stateManager,
		SharedStateManager: sharedStateManager,
		AIClient:           aiClient,
		ActionAdapter:      actionAdapter,
	}
	server, err := NewServer(serverCfg)
	require.NoError(t, err)
	
	return server, eventStream, stateManager
}

func TestGetState(t *testing.T) {
	server, _, stateManager := setupTestServer(t)
	
	err := stateManager.SetState("test", "123", map[string]interface{}{
		"key": "value",
	})
	require.NoError(t, err)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/state/test/123", nil)
	server.router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "value", response["key"])
}

func TestUpdateState(t *testing.T) {
	server, _, _ := setupTestServer(t)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/state/test/123", strings.NewReader(`{"key": "updated"}`))
	req.Header.Set("Content-Type", "application/json")
	server.router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/state/test/123", nil)
	server.router.ServeHTTP(w, req)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "updated", response["key"])
}

func TestListActions(t *testing.T) {
	server, _, _ := setupTestServer(t)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/actions", nil)
	server.router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	actions, ok := response["actions"].([]interface{})
	assert.True(t, ok)
	assert.Contains(t, actions, "test_action")
}

func TestExecuteAction(t *testing.T) {
	server, _, _ := setupTestServer(t)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/actions", strings.NewReader(`{"action": "test_action", "params": {"test": "value"}}`))
	req.Header.Set("Content-Type", "application/json")
	server.router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "test_action", response["action"])
	
	result, ok := response["result"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "success", result["result"])
	
	params, ok := result["params"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", params["test"])
}

func TestWebSocketConnection(t *testing.T) {
	server, eventStream, _ := setupTestServer(t)
	
	ts := httptest.NewServer(server.router)
	defer ts.Close()
	
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()
	
	err = ws.WriteJSON(map[string]interface{}{
		"action": "subscribe",
		"type":   "test",
		"id":     "123",
	})
	require.NoError(t, err)
	
	eventStream.Emit(eventstream.Event{
		Type: "state_update",
		Data: map[string]interface{}{
			"type": "test",
			"id":   "123",
			"data": map[string]interface{}{
				"key": "websocket_value",
			},
		},
	})
	
	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)
	
	var response map[string]interface{}
	err = json.Unmarshal(message, &response)
	require.NoError(t, err)
	
	assert.Equal(t, "state_update", response["type"])
	
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", data["type"])
	assert.Equal(t, "123", data["id"])
	
	stateData, ok := data["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "websocket_value", stateData["key"])
}
