package rocketmq

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type StateUpdateHandler func(ctx context.Context, msg *primitive.MessageExt) error

type StateConsumer struct {
	consumer rocketmq.PushConsumer
	handler  StateUpdateHandler
}

type ConsumerConfig struct {
	NameServerAddrs []string
	Topic           string
	GroupName       string
	ConsumeFromWhere consumer.ConsumeFromWhere // e.g., consumer.ConsumeFromFirstOffset
	SubscriptionExpression string // e.g., "*" for all tags, or "TagA || TagB"
}

func NewStateConsumer(cfg ConsumerConfig, handler StateUpdateHandler) (*StateConsumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(cfg.NameServerAddrs),
		consumer.WithGroupName(cfg.GroupName),
		consumer.WithConsumeFromWhere(cfg.ConsumeFromWhere),
		consumer.WithConsumerModel(consumer.Clustering), // Clustering is typical for state updates
		consumer.WithConsumeMessageBatchMaxSize(1), // Process one message at a time for simplicity
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create RocketMQ consumer: %w", err)
	}

	if handler == nil {
		return nil, fmt.Errorf("state update handler cannot be nil")
	}

	stateConsumer := &StateConsumer{
		consumer: c,
		handler:  handler,
	}

	err = c.Subscribe(cfg.Topic, consumer.MessageSelector{
		Type:       consumer.TAG, // Or SQL92
		Expression: cfg.SubscriptionExpression,
	}, stateConsumer.handleMessage) // Pass the internal handler method

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", cfg.Topic, err)
	}

	err = c.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start RocketMQ consumer: %w", err)
	}

	log.Printf("RocketMQ Consumer started successfully for topic %s, group %s\n", cfg.Topic, cfg.GroupName)

	return stateConsumer, nil
}

func (sc *StateConsumer) handleMessage(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		msg := msgs[i]
		log.Printf("Received state update message: ID=%s, Topic=%s, Keys=%s\n", msg.MsgId, msg.Topic, msg.GetKeys())

		err := sc.handler(ctx, msg)
		if err != nil {
			log.Printf("Error processing state update message ID %s: %v. Retrying...\n", msg.MsgId, err)
			return consumer.ConsumeRetryLater, err
		}
	}

	return consumer.ConsumeSuccess, nil
}

func (sc *StateConsumer) Close() error {
	log.Println("Shutting down RocketMQ Consumer...")
	err := sc.consumer.Shutdown()
	if err != nil {
		log.Printf("Error shutting down RocketMQ consumer: %v\n", err)
		return err
	}
	log.Println("RocketMQ Consumer shut down successfully.")
	return nil
}
