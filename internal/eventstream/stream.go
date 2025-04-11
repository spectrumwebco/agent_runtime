package eventstream

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8" // Using go-redis for DragonflyDB compatibility
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
)

type Stream struct {
	redisClient *redis.Client
	subscribers map[eventstream.EventType][]chan<- *eventstream.Event
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

type Config struct {
	DragonflyAddr string
	Password      string
	DB            int
}

func NewStream(cfg Config) (*Stream, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.DragonflyAddr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithCancel(context.Background())

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to DragonflyDB: %w", err)
	}

	log.Println("Successfully connected to DragonflyDB")

	stream := &Stream{
		redisClient: rdb,
		subscribers: make(map[eventstream.EventType][]chan<- *eventstream.Event),
		ctx:         ctx,
		cancel:      cancel,
	}


	return stream, nil
}

func (s *Stream) Publish(event *eventstream.Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("Publishing event: Type=%s, Source=%s, ID=%s\n", event.Type, event.Source, event.ID)

	eventKey := fmt.Sprintf("event:%s", event.ID)
	eventDataJSON, err := json.Marshal(event) // Requires importing "encoding/json"
	if err != nil {
		log.Printf("Error marshalling event data: %v\n", err)
	} else {
		err = s.redisClient.Set(s.ctx, eventKey, eventDataJSON, 24*time.Hour).Err() // Store for 24 hours
		if err != nil {
			log.Printf("Error storing event in DragonflyDB: %v\n", err)
		}
	}


	if subs, ok := s.subscribers[event.Type]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
				log.Printf("Subscriber channel full for event type %s\n", event.Type)
			}
		}
	}

	s.handleCacheInvalidation(event)

	return nil // Or return error if publishing fails critically
}

func (s *Stream) Subscribe(eventType eventstream.EventType, channel chan<- *eventstream.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribers[eventType] = append(s.subscribers[eventType], channel)
	log.Printf("New subscriber for event type: %s\n", eventType)
}

func (s *Stream) Unsubscribe(eventType eventstream.EventType, channel chan<- *eventstream.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if subs, ok := s.subscribers[eventType]; ok {
		newSubs := []chan<- *eventstream.Event{}
		for _, sub := range subs {
			if sub != channel {
				newSubs = append(newSubs, sub)
			}
		}
		s.subscribers[eventType] = newSubs
		log.Printf("Unsubscribed channel for event type: %s\n", eventType)
	}
}

func (s *Stream) handleCacheInvalidation(event *eventstream.Event) {

	log.Printf("Handling cache invalidation for event: %s\n", event.ID)

	if event.Type == eventstream.EventTypeObservation && event.Source == eventstream.EventSourceSandbox {
	}

	if event.Type == eventstream.EventTypeStateUpdate {
	}

}

func (s *Stream) GetAppContext(key string) (string, error) {
	cacheKey := fmt.Sprintf("context:%s", key)
	val, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	if err == redis.Nil {
		log.Printf("Cache miss for context key: %s\n", cacheKey)
		return "", fmt.Errorf("context key not found and rebuild logic not implemented")
	} else if err != nil {
		log.Printf("Error retrieving context from cache for key %s: %v\n", cacheKey, err)
		return "", fmt.Errorf("failed to retrieve context from cache: %w", err)
	}

	log.Printf("Cache hit for context key: %s\n", cacheKey)
	return val, nil
}

func (s *Stream) SetAppContext(key string, value string, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("context:%s", key)
	err := s.redisClient.Set(s.ctx, cacheKey, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting context cache for key %s: %v\n", cacheKey, err)
		return fmt.Errorf("failed to set context cache: %w", err)
	}
	log.Printf("Set context cache for key: %s\n", cacheKey)
	return nil
}


func (s *Stream) Close() error {
	s.cancel() // Signal background tasks to stop
	if err := s.redisClient.Close(); err != nil {
		log.Printf("Error closing DragonflyDB connection: %v\n", err)
		return err
	}
	log.Println("Event Stream closed.")
	return nil
}


import "encoding/json"
