package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/kled/socketio"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
)

func RegisterSocketIORoutes(router *gin.Engine, eventStream *eventstream.EventStream, stateManager *statemanager.StateManager) {
	socketio.RegisterSocketIOServer(router, eventStream, stateManager, "http://localhost:8000")
}
