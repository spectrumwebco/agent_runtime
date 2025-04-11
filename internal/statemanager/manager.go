package statemanager

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager/rocketmq"
)

type StateManager struct {
	producer       *rocketmq.StateProducer
	consumer       *rocketmq.StateConsumer
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	stateCache     map[string][]byte
	stateCacheLock sync.RWMutex
}

func NewStateManager(cfg *config.Config) (*StateManager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	producerCfg := rocketmq.ProducerConfig{
		NameServerAddrs: cfg.RocketMQ.NameServerAddrs,
		Topic:           cfg.RocketMQ.StateTopic,
		GroupName:       cfg.RocketMQ.ProducerGroupName,
		Retry:           cfg.RocketMQ.ProducerRetry,
	}

	producer, err := rocketmq.NewStateProducer(producerCfg)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create state producer: %w", err)
	}

	manager := &StateManager{
		producer:   producer,
		ctx:        ctx,
		cancel:     cancel,
		stateCache: make(map[string][]byte),
	}

	consumerCfg := rocketmq.ConsumerConfig{
		NameServerAddrs:        cfg.RocketMQ.NameServerAddrs,
		Topic:                  cfg.RocketMQ.StateTopic,
		GroupName:              cfg.RocketMQ.ConsumerGroupName,
		SubscriptionExpression: cfg.RocketMQ.ConsumerSubscription,
		ConsumeFromWhere:       rocketmq.ConsumeFromLastOffset,
	}

	consumer, err := rocketmq.NewStateConsumer(consumerCfg, manager.handleStateUpdate)
	if err != nil {
		producer.Close()
		cancel()
		return nil, fmt.Errorf("failed to create state consumer: %w", err)
	}

	manager.consumer = consumer

	log.Println("State Manager initialized successfully")
	return manager, nil
}

func (sm *StateManager) UpdateState(stateID string, stateType string, stateData []byte) error {
	sm.stateCacheLock.Lock()
	sm.stateCache[stateID] = stateData
	sm.stateCacheLock.Unlock()

	keys := []string{stateID, stateType}
	err := sm.producer.PublishStateUpdate(sm.ctx, stateData, keys)
	if err != nil {
		return fmt.Errorf("failed to publish state update: %w", err)
	}

	log.Printf("State update published for ID: %s, Type: %s\n", stateID, stateType)
	return nil
}

func (sm *StateManager) GetState(stateID string) ([]byte, bool) {
	sm.stateCacheLock.RLock()
	defer sm.stateCacheLock.RUnlock()

	data, exists := sm.stateCache[stateID]
	return data, exists
}

func (sm *StateManager) handleStateUpdate(ctx context.Context, msg *rocketmq.MessageExt) error {
	keys := msg.GetKeys()
	if len(keys) < 2 {
		return fmt.Errorf("invalid state update message: missing keys")
	}

	stateID := keys[0]
	stateType := keys[1]

	log.Printf("Received state update for ID: %s, Type: %s\n", stateID, stateType)

	sm.stateCacheLock.Lock()
	sm.stateCache[stateID] = msg.Body
	sm.stateCacheLock.Unlock()


	return nil
}

func (sm *StateManager) StartLifoTaskProcessor() {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
			case <-sm.ctx.Done():
				log.Println("LIFO task processor shutting down...")
				return
			}
		}
	}()

	log.Println("LIFO task processor started")
}

func (sm *StateManager) Close() error {
	log.Println("Shutting down State Manager...")
	sm.cancel()

	sm.wg.Wait()

	if err := sm.producer.Close(); err != nil {
		log.Printf("Error closing state producer: %v\n", err)
	}

	if err := sm.consumer.Close(); err != nil {
		log.Printf("Error closing state consumer: %v\n", err)
	}

	log.Println("State Manager shut down successfully")
	return nil
}
