package djangobridge

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client AgentServiceClient
}

func NewGRPCClient(grpcAddr string) (*GRPCClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Django gRPC server: %v", err)
	}

	client := NewAgentServiceClient(conn)

	return &GRPCClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *GRPCClient) ExecuteAgentTask(ctx context.Context, taskType string, taskData map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	taskDataProto, err := structpb.NewStruct(taskData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert task data to protobuf struct: %v", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &ExecuteTaskRequest{
		TaskType: taskType,
		TaskData: taskDataProto,
	}

	resp, err := c.client.ExecuteTask(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute task: %v", err)
	}

	result := resp.GetResult().AsMap()
	return result, nil
}

func (c *GRPCClient) QueryModel(ctx context.Context, modelName string, filters map[string]interface{}, timeout time.Duration) ([]map[string]interface{}, error) {
	filtersProto, err := structpb.NewStruct(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filters to protobuf struct: %v", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &QueryModelRequest{
		ModelName: modelName,
		Filters:   filtersProto,
	}

	resp, err := c.client.QueryModel(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query model: %v", err)
	}

	var results []map[string]interface{}
	for _, item := range resp.GetResults() {
		results = append(results, item.AsMap())
	}

	return results, nil
}

func (c *GRPCClient) UpdateModel(ctx context.Context, modelName string, modelID string, data map[string]interface{}, timeout time.Duration) error {
	dataProto, err := structpb.NewStruct(data)
	if err != nil {
		return fmt.Errorf("failed to convert data to protobuf struct: %v", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &UpdateModelRequest{
		ModelName: modelName,
		ModelId:   modelID,
		Data:      dataProto,
	}

	_, err = c.client.UpdateModel(ctxWithTimeout, req)
	if err != nil {
		return fmt.Errorf("failed to update model: %v", err)
	}

	return nil
}

func (c *GRPCClient) CreateModel(ctx context.Context, modelName string, data map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	dataProto, err := structpb.NewStruct(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert data to protobuf struct: %v", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &CreateModelRequest{
		ModelName: modelName,
		Data:      dataProto,
	}

	resp, err := c.client.CreateModel(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}

	result := resp.GetResult().AsMap()
	return result, nil
}

func (c *GRPCClient) DeleteModel(ctx context.Context, modelName string, modelID string, timeout time.Duration) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &DeleteModelRequest{
		ModelName: modelName,
		ModelId:   modelID,
	}

	_, err := c.client.DeleteModel(ctxWithTimeout, req)
	if err != nil {
		return fmt.Errorf("failed to delete model: %v", err)
	}

	return nil
}

func (c *GRPCClient) ExecuteManagementCommand(ctx context.Context, command string, args []string, timeout time.Duration) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &ExecuteCommandRequest{
		Command: command,
		Args:    args,
	}

	resp, err := c.client.ExecuteCommand(ctxWithTimeout, req)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return resp.GetOutput(), nil
}

func (c *GRPCClient) GetSettings(ctx context.Context, timeout time.Duration) (map[string]interface{}, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &GetSettingsRequest{}

	resp, err := c.client.GetSettings(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %v", err)
	}

	settings := resp.GetSettings().AsMap()
	return settings, nil
}

func (c *GRPCClient) ExecutePythonCode(ctx context.Context, code string, timeout time.Duration) (map[string]interface{}, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &ExecutePythonRequest{
		Code: code,
	}

	resp, err := c.client.ExecutePython(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Python code: %v", err)
	}

	result := resp.GetResult().AsMap()
	return result, nil
}

func (c *GRPCClient) Log(ctx context.Context, level string, message string, data map[string]interface{}) error {
	var dataProto *structpb.Struct
	var err error
	if data != nil {
		dataProto, err = structpb.NewStruct(data)
		if err != nil {
			return fmt.Errorf("failed to convert data to protobuf struct: %v", err)
		}
	}

	req := &LogRequest{
		Level:   level,
		Message: message,
		Data:    dataProto,
	}

	_, err = c.client.Log(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to log message: %v", err)
	}

	return nil
}

func (c *GRPCClient) Debug(ctx context.Context, message string, data map[string]interface{}) error {
	return c.Log(ctx, "DEBUG", message, data)
}

func (c *GRPCClient) Info(ctx context.Context, message string, data map[string]interface{}) error {
	return c.Log(ctx, "INFO", message, data)
}

func (c *GRPCClient) Warning(ctx context.Context, message string, data map[string]interface{}) error {
	return c.Log(ctx, "WARNING", message, data)
}

func (c *GRPCClient) Error(ctx context.Context, message string, data map[string]interface{}) error {
	return c.Log(ctx, "ERROR", message, data)
}

func (c *GRPCClient) Critical(ctx context.Context, message string, data map[string]interface{}) error {
	return c.Log(ctx, "CRITICAL", message, data)
}
