package database

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRocketMQClient struct {
	mock.Mock
}

func (m *MockRocketMQClient) CreateTopic(ctx context.Context, topic string) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockRocketMQClient) DeleteTopic(ctx context.Context, topic string) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockRocketMQClient) ListTopics(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRocketMQClient) SendMessage(ctx context.Context, topic, message string) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *MockRocketMQClient) ConsumeMessages(ctx context.Context, topic, group string, count int) ([]string, error) {
	args := m.Called(ctx, topic, group, count)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRocketMQClient) GetTopicStats(ctx context.Context, topic string) (map[string]interface{}, error) {
	args := m.Called(ctx, topic)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestRocketMQAdapter(t *testing.T) {
	t.Skip("Skipping tests that require actual RocketMQ connection")
	
	t.Run("NewRocketMQAdapter", func(t *testing.T) {
		adapter, err := NewRocketMQAdapter("localhost:9876")
		assert.NoError(t, err)
		assert.NotNil(t, adapter)
	})
	
	t.Run("ExecuteCommand", func(t *testing.T) {
		mockClient := new(MockRocketMQClient)
		
		mockClient.On("CreateTopic", mock.Anything, "test-topic").Return(nil)
		mockClient.On("SendMessage", mock.Anything, "test-topic", "test-message").Return(nil)
		mockClient.On("ListTopics", mock.Anything).Return([]string{"test-topic"}, nil)
		mockClient.On("DeleteTopic", mock.Anything, "test-topic").Return(nil)
		
		adapter := &RocketMQAdapter{
			connStr: "localhost:9876",
			client:  mockClient,
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := adapter.ExecuteCommand(ctx, "CREATE_TOPIC test-topic")
		assert.NoError(t, err)
		assert.Equal(t, "Topic 'test-topic' created successfully", result)
		
		result, err = adapter.ExecuteCommand(ctx, "SEND_MESSAGE test-topic test-message")
		assert.NoError(t, err)
		assert.Equal(t, "Message sent to topic 'test-topic' successfully", result)
		
		result, err = adapter.ExecuteCommand(ctx, "LIST_TOPICS")
		assert.NoError(t, err)
		assert.Equal(t, []string{"test-topic"}, result)
		
		result, err = adapter.ExecuteCommand(ctx, "DELETE_TOPIC test-topic")
		assert.NoError(t, err)
		assert.Equal(t, "Topic 'test-topic' deleted successfully", result)
		
		mockClient.AssertExpectations(t)
	})
	
	t.Run("Query", func(t *testing.T) {
		mockClient := new(MockRocketMQClient)
		
		mockClient.On("GetTopicStats", mock.Anything, "test-topic").Return(
			map[string]interface{}{
				"messageCount": 100,
				"tps": 10.5,
			}, nil)
		mockClient.On("ConsumeMessages", mock.Anything, "test-topic", "test-group", 10).Return(
			[]string{"message1", "message2"}, nil)
		
		adapter := &RocketMQAdapter{
			connStr: "localhost:9876",
			client:  mockClient,
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := adapter.Query(ctx, "GET_TOPIC_STATS test-topic")
		assert.NoError(t, err)
		stats, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, 100, stats["messageCount"])
		assert.Equal(t, 10.5, stats["tps"])
		
		result, err = adapter.Query(ctx, "CONSUME_MESSAGES test-topic test-group 10")
		assert.NoError(t, err)
		messages, ok := result.([]string)
		assert.True(t, ok)
		assert.Equal(t, 2, len(messages))
		assert.Equal(t, "message1", messages[0])
		assert.Equal(t, "message2", messages[1])
		
		mockClient.AssertExpectations(t)
	})
	
	t.Run("ErrorHandling", func(t *testing.T) {
		adapter := &RocketMQAdapter{
			connStr: "localhost:9876",
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		_, err := adapter.ExecuteCommand(ctx, "INVALID_COMMAND")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported command")
		
		_, err = adapter.ExecuteCommand(ctx, "CREATE_TOPIC")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires a topic name")
	})
}
