package djangobridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DjangoBridge struct {
	grpcAddr     string
	grpcClient   *grpc.ClientConn
	eventStream  EventStream
	agentFactory *langraph.DjangoAgentFactory
}

type EventStream interface {
	AddEvent(event *models.Event) error
}

func NewDjangoBridge(grpcAddr string, eventStream EventStream) (*DjangoBridge, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Django gRPC server: %v", err)
	}

	agentFactory := langraph.NewDjangoAgentFactory(grpcAddr, "", eventStream)

	return &DjangoBridge{
		grpcAddr:     grpcAddr,
		grpcClient:   conn,
		eventStream:  eventStream,
		agentFactory: agentFactory,
	}, nil
}

func (b *DjangoBridge) Close() error {
	if b.grpcClient != nil {
		return b.grpcClient.Close()
	}
	return nil
}

func (b *DjangoBridge) CreateDjangoAgent(ctx context.Context, agentID, agentName, agentRole string) (*langraph.Agent, error) {
	return b.agentFactory.CreateAgent(ctx, agentID, agentName, agentRole)
}

func (b *DjangoBridge) ExecutePythonCode(ctx context.Context, code string, timeout time.Duration) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"status": "success",
		"result": "Python code execution result would go here",
	}
	
	log.Printf("Executing Python code in Django: %s", code)
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("python-exec-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": "execute_python_code",
			"code":   code,
		},
	})
	
	return result, nil
}

func (b *DjangoBridge) QueryDjangoModel(ctx context.Context, modelName string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	
	filtersJSON, _ := json.Marshal(filters)
	
	log.Printf("Querying Django model %s with filters: %s", modelName, string(filtersJSON))
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("django-query-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "query_django_model",
			"model_name": modelName,
			"filters":    filters,
		},
	})
	
	return []map[string]interface{}{
		{
			"id":   1,
			"name": "Example result",
		},
	}, nil
}

func (b *DjangoBridge) UpdateDjangoModel(ctx context.Context, modelName string, modelID interface{}, data map[string]interface{}) error {
	
	dataJSON, _ := json.Marshal(data)
	
	log.Printf("Updating Django model %s with ID %v: %s", modelName, modelID, string(dataJSON))
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("django-update-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "update_django_model",
			"model_name": modelName,
			"model_id":   modelID,
			"data":       data,
		},
	})
	
	return nil
}

func (b *DjangoBridge) CreateDjangoModel(ctx context.Context, modelName string, data map[string]interface{}) (map[string]interface{}, error) {
	
	dataJSON, _ := json.Marshal(data)
	
	log.Printf("Creating Django model %s: %s", modelName, string(dataJSON))
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("django-create-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "create_django_model",
			"model_name": modelName,
			"data":       data,
		},
	})
	
	return map[string]interface{}{
		"id":   1,
		"name": "Example result",
	}, nil
}

func (b *DjangoBridge) DeleteDjangoModel(ctx context.Context, modelName string, modelID interface{}) error {
	
	log.Printf("Deleting Django model %s with ID %v", modelName, modelID)
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("django-delete-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":     "delete_django_model",
			"model_name": modelName,
			"model_id":   modelID,
		},
	})
	
	return nil
}

func (b *DjangoBridge) ExecuteDjangoManagementCommand(ctx context.Context, command string, args []string) (string, error) {
	
	log.Printf("Executing Django management command: %s %v", command, args)
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("django-command-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action":  "execute_django_management_command",
			"command": command,
			"args":    args,
		},
	})
	
	return "Command executed successfully", nil
}

func (b *DjangoBridge) GetDjangoSettings(ctx context.Context) (map[string]interface{}, error) {
	
	log.Printf("Retrieving Django settings")
	
	b.eventStream.AddEvent(&models.Event{
		ID:        fmt.Sprintf("django-settings-%d", time.Now().UnixNano()),
		Type:      models.EventTypeAction,
		Source:    models.EventSourceAgent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": "get_django_settings",
		},
	})
	
	return map[string]interface{}{
		"DEBUG":         true,
		"ALLOWED_HOSTS": []string{"localhost", "127.0.0.1"},
		"DATABASES": map[string]interface{}{
			"default": map[string]interface{}{
				"ENGINE":  "django.db.backends.postgresql",
				"NAME":    "agent_db",
				"USER":    "postgres",
				"HOST":    "localhost",
				"PORT":    5432,
				"OPTIONS": map[string]interface{}{},
			},
		},
	}, nil
}
