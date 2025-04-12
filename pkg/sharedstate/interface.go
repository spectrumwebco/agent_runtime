package sharedstate

import (
	"github.com/gin-gonic/gin"
)

type StateManager interface {
	GetState(stateType StateType, stateID string) (map[string]interface{}, error)
	
	UpdateState(stateType StateType, stateID string, data map[string]interface{}) error
	
	RegisterWebSocketHandler(router gin.IRouter)
	
	Close() error
}

type EventHandler interface {
	HandleEvent(eventType string, data map[string]interface{}) error
}

type WebSocketServer interface {
	Start() error
	
	Stop() error
	
	PublishToTopic(topic string, data interface{})
	
	SubscribeToTopic(clientID, topic string)
	
	UnsubscribeFromTopic(clientID, topic string)
	
	HandleWebSocket(c *gin.Context)
}
