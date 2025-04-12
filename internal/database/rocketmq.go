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

type RocketMQAdapter struct {
	connStr     string
	producer    rocketmq.Producer
	consumer    rocketmq.PushConsumer
	admin       admin.Admin
	nameServers []string
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
	
	return &RocketMQAdapter{
		connStr:     connStr,
		producer:    p,
		consumer:    c,
		admin:       adm,
		nameServers: nameServers,
	}, nil
}

func (a *RocketMQAdapter) Query(ctx context.Context, query string) (interface{}, error) {
	if strings.HasPrefix(query, "GET_TOPIC_STATS") {
		topicName := strings.TrimPrefix(query, "GET_TOPIC_STATS ")
		stats, err := a.admin.ExamineTopicStats(ctx, topicName)
		if err != nil {
			return nil, fmt.Errorf("failed to get topic stats: %w", err)
		}
		return stats, nil
	} else if strings.HasPrefix(query, "GET_CONSUMER_GROUP_INFO") {
		groupName := strings.TrimPrefix(query, "GET_CONSUMER_GROUP_INFO ")
		info, err := a.admin.ExamineConsumerGroupInfo(ctx, groupName)
		if err != nil {
			return nil, fmt.Errorf("failed to get consumer group info: %w", err)
		}
		return info, nil
	} else if strings.HasPrefix(query, "GET_TOPICS") {
		topics, err := a.admin.FetchAllTopicList(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch topics: %w", err)
		}
		return topics, nil
	} else if strings.HasPrefix(query, "GET_CONSUMER_GROUPS") {
		groups, err := a.admin.ExamineSubscriptionGroupConfig(ctx, "")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch consumer groups: %w", err)
		}
		return groups, nil
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
		err := a.admin.CreateTopic(ctx, admin.NewTopicConfig(topicName, 8, 3))
		if err != nil {
			return nil, fmt.Errorf("failed to create topic: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"topic":  topicName,
		}, nil
	} else if strings.HasPrefix(command, "SEND_MESSAGE") {
		parts := strings.Split(strings.TrimPrefix(command, "SEND_MESSAGE "), " ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("topic and message body are required")
		}
		
		topicName := parts[0]
		messageBody := strings.Join(parts[1:], " ")
		
		msg := &primitive.Message{
			Topic: topicName,
			Body:  []byte(messageBody),
		}
		
		result, err := a.producer.SendSync(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to send message: %w", err)
		}
		
		return map[string]interface{}{
			"message_id": result.MsgID,
			"status":     "success",
			"topic":      topicName,
		}, nil
	} else if strings.HasPrefix(command, "DELETE_TOPIC") {
		topicName := strings.TrimPrefix(command, "DELETE_TOPIC ")
		err := a.admin.DeleteTopic(ctx, admin.NewTopicConfig(topicName, 0, 0))
		if err != nil {
			return nil, fmt.Errorf("failed to delete topic: %w", err)
		}
		
		return map[string]interface{}{
			"status": "success",
			"topic":  topicName,
		}, nil
	}
	
	return nil, fmt.Errorf("unsupported command: %s", command)
}

func (a *RocketMQAdapter) GetSchema(ctx context.Context) (interface{}, error) {
	schema := make(map[string]interface{})
	
	topics, err := a.admin.FetchAllTopicList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch topics: %w", err)
	}
	
	topicDetails := make([]map[string]interface{}, 0, len(topics.TopicList))
	for _, topic := range topics.TopicList {
		stats, err := a.admin.ExamineTopicStats(ctx, topic)
		if err != nil {
			continue
		}
		
		topicDetail := map[string]interface{}{
			"name":           topic,
			"message_count":  stats.OffsetTable,
			"creation_time":  time.Now().String(), // Not provided by API
		}
		
		topicDetails = append(topicDetails, topicDetail)
	}
	
	schema["topics"] = topicDetails
	
	brokerData, err := a.admin.FetchBrokerNameList(ctx)
	if err == nil {
		schema["brokers"] = brokerData
	}
	
	clusterInfo, err := a.admin.FetchClusterInfo(ctx)
	if err == nil {
		schema["clusters"] = clusterInfo
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
