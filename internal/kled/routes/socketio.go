package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/kled/socketio"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
)

func RegisterSocketIORoutes(router *gin.Engine, eventStream *eventstream.EventStream, stateManager *statemanager.StateManager) {
	socketio.RegisterSocketIOServer(router, eventStream, stateManager, "http://localhost:8000")
}

func AddSocketIORoutesToRouter(router *gin.Engine, eventStream *eventstream.EventStream, stateManager *statemanager.StateManager) {
	socketIOServer := socketio.NewServer(eventStream, stateManager)

	djangoMiddleware := socketio.NewDjangoMiddleware(socketIOServer, "http://localhost:8000", eventStream)

	djangoMiddleware.RegisterRoutes(router)
}
