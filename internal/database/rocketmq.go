package database

import (
	"context"
	"fmt"
	"strings"
	"time"
	
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketMQClient interface {
	CreateTopic(ctx context.Context, topic string) error
	DeleteTopic(ctx context.Context, topic string) error
	ListTopics(ctx context.Context) ([]string, error)
	SendMessage(ctx context.Context, topic, message string) error
	ConsumeMessages(ctx context.Context, topic, group string, count int) ([]string, error)
	GetTopicStats(ctx context.Context, topic string) (map[string]interface{}, error)
}

type RocketMQAdapter struct {
	connStr     string
	producer    rocketmq.Producer
	consumer    rocketmq.PushConsumer
	admin       admin.Admin
	nameServers []string
	client      RocketMQClient
}

func NewRocketMQAdapter(connStr string) (*RocketMQAdapter, error) {
	parts := strings.Split(connStr, ";")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid connection string format")
	}
	
	nameServers := []string{"localhost:9876"} // Default
	for _, part := range parts {
		if strings.HasPrefix(part, "nameServer=") {
			servers := strings.TrimPrefix(part, "nameServer=")
			nameServers = strings.Split(servers, ",")
		}
	}
	
	p, err := producer.NewDefaultProducer(
		producer.WithNameServer(nameServers),
		producer.WithRetry(2),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	
	if err := p.Start(); err != nil {
		return nil, fmt.Errorf("failed to start producer: %w", err)
	}
	
	c, err := consumer.NewPushConsumer(
		consumer.WithNameServer(nameServers),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		p.Shutdown()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	
	adm, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameServers)),
	)
	if err != nil {
		p.Shutdown()
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}
	
	adapter := &RocketMQAdapter{
		connStr:     connStr,
		producer:    p,
		consumer:    c,
		admin:       adm,
		nameServers: nameServers,
	}
	
	adapter.client = &defaultRocketMQClient{
		producer: p,
		consumer: c,
		admin:    adm,
	}
	
	return adapter, nil
}

func (a *RocketMQAdapter) Query(ctx context.Context, query string) (interface{}, error) {
	if strings.HasPrefix(query, "GET_TOPIC_STATS") {
		topicName := strings.TrimPrefix(query, "GET_TOPIC_STATS ")
		stats, err := a.client.GetTopicStats(ctx, topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to get topic stats: %w", err)
		}
		return stats, nil
	} else if strings.HasPrefix(query, "GET_CONSUMER_GROUP_INFO") {
		groupName := strings.TrimPrefix(query, "GET_CONSUMER_GROUP_INFO ")
		return map[string]interface{}{
			"group": groupName,
			"connections": 0,
		}, nil
	} else if strings.HasPrefix(query, "GET_TOPICS") {
		topics, err := a.client.ListTopics(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch topics: %w", err)
		}
		return topics, nil
	} else if strings.HasPrefix(query, "GET_CONSUMER_GROUPS") {
		return []string{"default-consumer-group"}, nil
	} else if strings.HasPrefix(query, "CONSUME_MESSAGES") {
		parts := strings.Split(strings.TrimPrefix(query, "CONSUME_MESSAGES "), " ")
		if len(parts) < 3 {
			return nil, fmt.Errorf("topic, group, and count are required")
		}
		
		topic := parts[0]
		group := parts[1]
		count := 10 // Default
		
		messages, err := a.client.ConsumeMessages(ctx, topic, group, count)
		if err != nil {
			return nil, fmt.Errorf("failed to consume messages: %w", err)
		}
		
		return messages, nil
	}
	
	return nil, fmt.Errorf("unsupported query operation: %s", query)
}

func (a *RocketMQAdapter) Execute(ctx context.Context, command string) (interface{}, error) {
	if strings.HasPrefix(command, "CREATE_TOPIC") {
		parts := strings.Split(strings.TrimPrefix(command, "CREATE_TOPIC "), " ")
		if len(parts) < 1 {
			return nil, fmt.Errorf("topic name is required")
		}
		
		topicName := parts[0]
		err := a.client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic: %w", err)
		}
		
		return "Topic '" + topicName + "' created successfully", nil
	} else if strings.HasPrefix(command, "SEND_MESSAGE") {
		parts := strings.Split(strings.TrimPrefix(command, "SEND_MESSAGE "), " ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("topic and message body are required")
		}
		
		topicName := parts[0]
		messageBody := strings.Join(parts[1:], " ")
		
		err := a.client.SendMessage(ctx, topicName, messageBody)
		if err != nil {
			return nil, fmt.Errorf("failed to send message: %w", err)
		}
		
		return "Message sent to topic '" + topicName + "' successfully", nil
	} else if strings.HasPrefix(command, "DELETE_TOPIC") {
		topicName := strings.TrimPrefix(command, "DELETE_TOPIC ")
		err := a.client.DeleteTopic(ctx, topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to delete topic: %w", err)
		}
		
		return "Topic '" + topicName + "' deleted successfully", nil
	} else if strings.HasPrefix(command, "LIST_TOPICS") {
		topics, err := a.client.ListTopics(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list topics: %w", err)
		}
		
		return topics, nil
	}
	
	return nil, fmt.Errorf("unsupported command: %s", command)
}

func (a *RocketMQAdapter) GetSchema(ctx context.Context) (interface{}, error) {
	schema := make(map[string]interface{})
	
	topics, err := a.client.ListTopics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch topics: %w", err)
	}
	
	topicDetails := make([]map[string]interface{}, 0, len(topics))
	for _, topic := range topics {
		stats, err := a.client.GetTopicStats(ctx, topic)
		if err != nil {
			continue
		}
		
		topicDetail := map[string]interface{}{
			"name":          topic,
			"message_count": stats["totalCount"],
			"creation_time": time.Now().String(), // Not provided by API
		}
		
		topicDetails = append(topicDetails, topicDetail)
	}
	
	schema["topics"] = topicDetails
	schema["brokers"] = []map[string]interface{}{
		{
			"name":    "broker-1",
			"address": a.nameServers[0],
			"version": "4.9.3",
		},
	}
	schema["clusters"] = []map[string]interface{}{
		{
			"name":          "DefaultCluster",
			"broker_count":  1,
			"topic_count":   len(topics),
			"message_count": 0,
		},
	}
	
	return schema, nil
}

func (a *RocketMQAdapter) Close() error {
	var errs []string
	
	if err := a.producer.Shutdown(); err != nil {
		errs = append(errs, fmt.Sprintf("producer shutdown error: %v", err))
	}
	
	if err := a.consumer.Shutdown(); err != nil {
		errs = append(errs, fmt.Sprintf("consumer shutdown error: %v", err))
	}
	
	if err := a.admin.Close(); err != nil {
		errs = append(errs, fmt.Sprintf("admin close error: %v", err))
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %s", strings.Join(errs, "; "))
	}
	
	return nil
}
