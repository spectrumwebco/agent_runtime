package rocketmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type StateProducer struct {
	producer rocketmq.Producer
	topic    string
}

type ProducerConfig struct {
	NameServerAddrs []string
	Topic           string
	GroupName       string
	Retry           int
}

func NewStateProducer(cfg ProducerConfig) (*StateProducer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.NameServerAddrs),
		producer.WithRetry(cfg.Retry),
		producer.WithGroupName(cfg.GroupName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create RocketMQ producer: %w", err)
	}

	err = p.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start RocketMQ producer: %w", err)
	}

	log.Printf("RocketMQ Producer started successfully for topic %s\n", cfg.Topic)

	return &StateProducer{
		producer: p,
		topic:    cfg.Topic,
	}, nil
}

func (sp *StateProducer) PublishStateUpdate(ctx context.Context, stateData []byte, keys []string) error {
	msg := primitive.NewMessage(sp.topic, stateData)
	msg.WithKeys(keys) // Use keys for potential filtering/routing

	result, err := sp.producer.SendSync(ctx, msg)

	if err != nil {
		log.Printf("Failed to send state update message: %v\n", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Successfully sent state update message: %s\n", result.String())
	return nil
}

func (sp *StateProducer) Close() error {
	log.Println("Shutting down RocketMQ Producer...")
	err := sp.producer.Shutdown()
	if err != nil {
		log.Printf("Error shutting down RocketMQ producer: %v\n", err)
		return err
	}
	log.Println("RocketMQ Producer shut down successfully.")
	return nil
}
