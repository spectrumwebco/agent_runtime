package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type RocketMQAdapter struct {
	nameServers    []string
	producer       rocketmq.Producer
	consumer       rocketmq.PushConsumer
	stateCache     map[string]*models.State
	stateCacheLock sync.RWMutex
	initialized    bool
	consumerGroup  string
	stateTopic     string
	clientID       string
}

type RocketMQConfig struct {
	NameServers   []string
	ConsumerGroup string
	StateTopic    string
	ClientID      string
}

func NewRocketMQAdapter(config RocketMQConfig) *RocketMQAdapter {
	if config.ClientID == "" {
		config.ClientID = uuid.New().String()
	}
	
	if config.StateTopic == "" {
		config.StateTopic = "shared_state"
	}
	
	if config.ConsumerGroup == "" {
		config.ConsumerGroup = "shared_state_consumers"
	}

	return &RocketMQAdapter{
		nameServers:   config.NameServers,
		stateCache:    make(map[string]*models.State),
		consumerGroup: config.ConsumerGroup,
		stateTopic:    config.StateTopic,
		clientID:      config.ClientID,
	}
}

func (a *RocketMQAdapter) Initialize() error {
	if a.initialized {
		return nil
	}

	if len(a.nameServers) == 0 {
		return fmt.Errorf("RocketMQ name servers are required")
	}

	p, err := rocketmq.NewProducer(
		producer.WithNameServer(a.nameServers),
		producer.WithRetry(2),
		producer.WithGroupName(fmt.Sprintf("%s-producer", a.consumerGroup)),
	)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start producer: %w", err)
	}
	a.producer = p

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(a.nameServers),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(a.consumerGroup),
		consumer.WithInstance(a.clientID),
	)
	if err != nil {
		p.Shutdown()
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	if err := c.Subscribe(a.stateTopic, consumer.MessageSelector{}, a.handleStateMessage); err != nil {
		p.Shutdown()
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	if err := c.Start(); err != nil {
		p.Shutdown()
		return fmt.Errorf("failed to start consumer: %w", err)
	}
	a.consumer = c

	a.initialized = true
	return nil
}

func (a *RocketMQAdapter) handleStateMessage(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		var state models.State
		if err := json.Unmarshal(msg.Body, &state); err != nil {
			fmt.Printf("Failed to unmarshal state message: %v\n", err)
			continue
		}

		a.stateCacheLock.Lock()
		a.stateCache[state.ID] = &state
		a.stateCacheLock.Unlock()
	}

	return consumer.ConsumeSuccess, nil
}

func (a *RocketMQAdapter) GetState(stateType models.StateType, stateID string) (*models.State, error) {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	a.stateCacheLock.RLock()
	state, exists := a.stateCache[stateID]
	a.stateCacheLock.RUnlock()

	if exists && state.Type == stateType {
		return state, nil
	}

	return nil, fmt.Errorf("state not found in cache")
}

func (a *RocketMQAdapter) UpdateState(state *models.State) error {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return err
		}
	}

	state.Version++
	state.UpdatedAt = time.Now().UTC()

	a.stateCacheLock.Lock()
	a.stateCache[state.ID] = state
	a.stateCacheLock.Unlock()

	body, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	msg := &primitive.Message{
		Topic: a.stateTopic,
		Body:  body,
	}

	msg.WithProperty("type", string(state.Type))
	msg.WithProperty("id", state.ID)
	msg.WithProperty("version", fmt.Sprintf("%d", state.Version))

	_, err = a.producer.SendSync(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (a *RocketMQAdapter) CreateState(state *models.State) error {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return err
		}
	}

	a.stateCacheLock.Lock()
	a.stateCache[state.ID] = state
	a.stateCacheLock.Unlock()

	body, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	msg := &primitive.Message{
		Topic: a.stateTopic,
		Body:  body,
	}

	msg.WithProperty("type", string(state.Type))
	msg.WithProperty("id", state.ID)
	msg.WithProperty("version", fmt.Sprintf("%d", state.Version))

	_, err = a.producer.SendSync(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (a *RocketMQAdapter) DeleteState(stateType models.StateType, stateID string) error {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return err
		}
	}

	a.stateCacheLock.Lock()
	delete(a.stateCache, stateID)
	a.stateCacheLock.Unlock()

	tombstone := &models.State{
		ID:        stateID,
		Type:      stateType,
		Data:      map[string]interface{}{"_deleted": true},
		Metadata:  map[string]string{"_tombstone": "true"},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Version:   0,
	}

	body, err := json.Marshal(tombstone)
	if err != nil {
		return fmt.Errorf("failed to marshal tombstone: %w", err)
	}

	msg := &primitive.Message{
		Topic: a.stateTopic,
		Body:  body,
	}

	msg.WithProperty("type", string(stateType))
	msg.WithProperty("id", stateID)
	msg.WithProperty("tombstone", "true")

	_, err = a.producer.SendSync(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to send tombstone message: %w", err)
	}

	return nil
}

func (a *RocketMQAdapter) ListStates(stateType models.StateType) ([]*models.State, error) {
	if !a.initialized {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	var states []*models.State

	a.stateCacheLock.RLock()
	for _, state := range a.stateCache {
		if state.Type == stateType {
			states = append(states, state)
		}
	}
	a.stateCacheLock.RUnlock()

	return states, nil
}

func (a *RocketMQAdapter) Close() error {
	if !a.initialized {
		return nil
	}

	if a.producer != nil {
		a.producer.Shutdown()
	}
	if a.consumer != nil {
		a.consumer.Shutdown()
	}

	a.initialized = false
	return nil
}
