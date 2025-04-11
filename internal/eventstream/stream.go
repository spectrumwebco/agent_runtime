package eventstream

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"encoding/json"
	"strings"
	"path/filepath"

	"github.com/go-redis/redis/v8" // Using go-redis for DragonflyDB compatibility
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
)

type Stream struct {
	redisClient     *redis.Client
	rocketProducer  rocketmq.Producer
	rocketConsumer  rocketmq.PushConsumer
	subscribers     map[eventstream.EventType][]chan<- *eventstream.Event
	rebuildRegistry ContextRebuildRegistry
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
}

type Config struct {
	DragonflyAddr  string
	Password       string
	DB             int
	RocketMQAddr   []string
	RocketMQTopic  string
	RocketMQGroup  string
	RocketMQRetry  int
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

	p, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.RocketMQAddr),
		producer.WithRetry(cfg.RocketMQRetry),
		producer.WithGroupName(cfg.RocketMQGroup),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create RocketMQ producer: %w", err)
	}

	if err := p.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start RocketMQ producer: %w", err)
	}

	log.Println("Successfully started RocketMQ producer")

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(cfg.RocketMQAddr),
		consumer.WithGroupName(cfg.RocketMQGroup),
	)
	if err != nil {
		p.Shutdown()
		cancel()
		return nil, fmt.Errorf("failed to create RocketMQ consumer: %w", err)
	}

	rebuildRegistry := NewContextRebuildRegistry()

	stream := &Stream{
		redisClient:     rdb,
		rocketProducer:  p,
		rocketConsumer:  c,
		subscribers:     make(map[eventstream.EventType][]chan<- *eventstream.Event),
		rebuildRegistry: rebuildRegistry,
		ctx:             ctx,
		cancel:          cancel,
	}

	err = c.Subscribe(cfg.RocketMQTopic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var event eventstream.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Error unmarshalling RocketMQ message: %v\n", err)
				continue
			}
			
			stream.processRemoteEvent(&event)
		}
		return consumer.ConsumeSuccess, nil
	})
	
	if err != nil {
		p.Shutdown()
		cancel()
		return nil, fmt.Errorf("failed to subscribe to RocketMQ topic: %w", err)
	}

	if err := c.Start(); err != nil {
		p.Shutdown()
		cancel()
		return nil, fmt.Errorf("failed to start RocketMQ consumer: %w", err)
	}

	log.Println("Successfully started RocketMQ consumer")

	return stream, nil
}

func (s *Stream) Publish(event *eventstream.Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("Publishing event: Type=%s, Source=%s, ID=%s\n", event.Type, event.Source, event.ID)

	eventKey := fmt.Sprintf("event:%s", event.ID)
	eventDataJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event data: %v\n", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	
	err = s.redisClient.Set(s.ctx, eventKey, eventDataJSON, 24*time.Hour).Err() // Store for 24 hours
	if err != nil {
		log.Printf("Error storing event in DragonflyDB: %v\n", err)
	}

	msg := primitive.NewMessage(
		"TOPIC_EVENT_STREAM", // This should be configurable
		eventDataJSON,
	)
	msg.WithKeys([]string{string(event.Type), string(event.Source)})
	
	res, err := s.rocketProducer.SendSync(context.Background(), msg)
	if err != nil {
		log.Printf("Error publishing event to RocketMQ: %v\n", err)
	} else {
		log.Printf("Event published to RocketMQ: MessageID=%s, Queue=%d\n", 
			res.MsgID, res.QueueOffset)
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

	return nil
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
	log.Printf("Handling cache invalidation for event: %s (Type: %s, Source: %s)\n", 
		event.ID, event.Type, event.Source)
	
	var keysToInvalidate []string
	
	switch event.Source {
	case eventstream.EventSourceSandbox:
		if event.Type == eventstream.EventTypeObservation {
			if dataMap, ok := event.Data.(map[string]interface{}); ok {
				if dataType, ok := dataMap["type"].(string); ok && dataType == "file_change" {
					if filePath, ok := dataMap["path"].(string); ok && filePath != "" {
						cacheKey := fmt.Sprintf("context:file:%s", filePath)
						keysToInvalidate = append(keysToInvalidate, cacheKey)
						
						dirPath := filepath.Dir(filePath)
						keysToInvalidate = append(keysToInvalidate, fmt.Sprintf("context:dir:%s", dirPath))
					}
				}
			}
		}
		
	case eventstream.EventSourceCICD:
		if event.Type == eventstream.EventTypeStateUpdate {
			if dataMap, ok := event.Data.(map[string]interface{}); ok {
				if service, ok := dataMap["service"].(string); ok && service != "" {
					keysToInvalidate = append(keysToInvalidate, fmt.Sprintf("context:ci:pipeline:%s", service))
					keysToInvalidate = append(keysToInvalidate, fmt.Sprintf("context:k8s:deployment:%s", service))
				}
			}
		}
		
	case eventstream.EventSourceK8s:
		if event.Type == eventstream.EventTypeStateUpdate {
			if dataMap, ok := event.Data.(map[string]interface{}); ok {
				resourceType, hasType := dataMap["resource_type"].(string)
				resourceName, hasName := dataMap["resource_name"].(string)
				
				if hasType && hasName {
					keysToInvalidate = append(keysToInvalidate, 
						fmt.Sprintf("context:k8s:%s:%s", resourceType, resourceName))
				}
			}
		}
		
	case eventstream.EventSourceAgent:
		if event.Type == eventstream.EventTypeAction {
			if dataMap, ok := event.Data.(map[string]interface{}); ok {
				if action, ok := dataMap["action"].(string); ok {
					if strings.HasPrefix(action, "git") {
						keysToInvalidate = append(keysToInvalidate, "context:git:status")
						keysToInvalidate = append(keysToInvalidate, "context:git:log")
					} else if strings.Contains(action, "file") || strings.Contains(action, "edit") {
						if filePath, ok := dataMap["path"].(string); ok {
							keysToInvalidate = append(keysToInvalidate, 
								fmt.Sprintf("context:file:%s", filePath))
						}
					}
				}
			}
		}
	}
	
	
	if len(keysToInvalidate) > 0 {
		log.Printf("Invalidating %d cache keys: %v\n", len(keysToInvalidate), keysToInvalidate)
		
		deletedCount, err := s.redisClient.Unlink(s.ctx, keysToInvalidate...).Result()
		if err != nil {
			log.Printf("Error invalidating cache keys with UNLINK: %v. Falling back to DEL.\n", err)
			deletedCount, err = s.redisClient.Del(s.ctx, keysToInvalidate...).Result()
			if err != nil {
				log.Printf("Error invalidating cache keys with DEL: %v\n", err)
			}
		}
		
		log.Printf("Successfully invalidated %d keys.\n", deletedCount)
		
		cacheUpdateData := map[string]interface{}{
			"invalidated_keys": keysToInvalidate,
			"timestamp":        time.Now().Unix(),
		}
		
		go func() {
			cacheUpdateEvent := eventstream.NewEvent(
				eventstream.EventTypeCacheUpdate,
				eventstream.EventSourceSystem,
				cacheUpdateData,
				map[string]string{"origin_event_id": event.ID},
			)
			s.publishWithoutInvalidation(cacheUpdateEvent)
		}()
	}
}

func (s *Stream) GetAppContext(key string) (string, error) {
	cacheKey := fmt.Sprintf("context:%s", key)
	
	val, err := s.redisClient.Get(s.ctx, cacheKey).Result()
	if err == nil {
		log.Printf("Cache hit for context key: %s\n", cacheKey)
		return val, nil
	} else if err != redis.Nil {
		log.Printf("Error retrieving context from cache for key %s: %v\n", cacheKey, err)
		return "", fmt.Errorf("failed to retrieve context from cache: %w", err)
	}
	
	log.Printf("Cache miss for context key: %s. Attempting to rebuild.\n", cacheKey)
	
	parts := strings.SplitN(key, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid context key format: %s", key)
	}
	
	contextType := parts[0]
	contextIdentifier := parts[1]
	
	rebuildFunc, ok := s.rebuildRegistry[contextType]
	if !ok {
		return "", fmt.Errorf("no rebuild function registered for context type: %s", contextType)
	}
	
	rebuiltContext, err := rebuildFunc(contextIdentifier)
	if err != nil {
		log.Printf("Failed to rebuild context for key %s: %v\n", key, err)
		return "", fmt.Errorf("failed to rebuild context: %w", err)
	}
	
	expiration := GetCacheExpiration(contextType)
	err = s.SetAppContext(key, rebuiltContext, expiration)
	if err != nil {
		log.Printf("Warning: Failed to cache rebuilt context for key %s: %v\n", key, err)
	} else {
		log.Printf("Successfully rebuilt and cached context for key %s with TTL %v\n", key, expiration)
	}
	
	return rebuiltContext, nil
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


func (s *Stream) processRemoteEvent(event *eventstream.Event) {
	log.Printf("Processing remote event: Type=%s, Source=%s, ID=%s\n", 
		event.Type, event.Source, event.ID)
	
	eventKey := fmt.Sprintf("event:%s", event.ID)
	eventDataJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling remote event data: %v\n", err)
		return
	}
	
	err = s.redisClient.Set(s.ctx, eventKey, eventDataJSON, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Error storing remote event in DragonflyDB: %v\n", err)
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if subs, ok := s.subscribers[event.Type]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
				log.Printf("Subscriber channel full for remote event type %s\n", event.Type)
			}
		}
	}
	
	s.handleCacheInvalidation(event)
}

func (s *Stream) Close() error {
	s.cancel() // Signal background tasks to stop
	
	if s.rocketProducer != nil {
		if err := s.rocketProducer.Shutdown(); err != nil {
			log.Printf("Error shutting down RocketMQ producer: %v\n", err)
		} else {
			log.Println("RocketMQ producer shut down successfully")
		}
	}
	
	if s.rocketConsumer != nil {
		if err := s.rocketConsumer.Shutdown(); err != nil {
			log.Printf("Error shutting down RocketMQ consumer: %v\n", err)
		} else {
			log.Println("RocketMQ consumer shut down successfully")
		}
	}
	
	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			log.Printf("Error closing DragonflyDB connection: %v\n", err)
			return err
		}
		log.Println("DragonflyDB connection closed successfully")
	}
	
	log.Println("Event Stream closed.")
	return nil
}


func (s *Stream) publishWithoutInvalidation(event *eventstream.Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("Publishing event without invalidation: Type=%s, Source=%s, ID=%s\n", 
		event.Type, event.Source, event.ID)

	eventKey := fmt.Sprintf("event:%s", event.ID)
	eventDataJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event data: %v\n", err)
		return err
	}
	
	err = s.redisClient.Set(s.ctx, eventKey, eventDataJSON, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Error storing event in DragonflyDB: %v\n", err)
		return err
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

	return nil
}
