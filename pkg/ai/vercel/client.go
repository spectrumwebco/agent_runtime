package vercel

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type VercelAIClient struct {
	sharedStateManager *sharedstate.SharedStateManager
	eventStream        *eventstream.Stream
	mu                 sync.RWMutex
	streamingHandlers  map[string]chan<- interface{}
}

func NewVercelAIClient(sharedStateManager *sharedstate.SharedStateManager, eventStream *eventstream.Stream) *VercelAIClient {
	return &VercelAIClient{
		sharedStateManager: sharedStateManager,
		eventStream:        eventStream,
		streamingHandlers:  make(map[string]chan<- interface{}),
	}
}

func (c *VercelAIClient) RegisterStreamingHandler(id string, handler chan<- interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.streamingHandlers[id] = handler
}

func (c *VercelAIClient) UnregisterStreamingHandler(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.streamingHandlers, id)
}

func (c *VercelAIClient) SendStreamingEvent(event interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for id, handler := range c.streamingHandlers {
		select {
		case handler <- event:
		default:
			log.Printf("Streaming handler %s is full, dropping event", id)
		}
	}
}

func (c *VercelAIClient) HandleStateUpdate(stateType sharedstate.StateType, stateID string, data map[string]interface{}) {
	c.SendStreamingEvent(map[string]interface{}{
		"type":       "state_update",
		"state_type": string(stateType),
		"state_id":   stateID,
		"data":       data,
	})
}

func (c *VercelAIClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for id, handler := range c.streamingHandlers {
		close(handler)
		delete(c.streamingHandlers, id)
	}
}
