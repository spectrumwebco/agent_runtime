package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Task struct {
	ID      string
	Status  string
	Prompt  string
	Context map[string]string
	Tools   []string
	Result  string
	Events  []string
}

type AgentServiceServer interface {
	ExecuteTask(context.Context, *ExecuteTaskRequest) (*ExecuteTaskResponse, error)
	GetTaskStatus(context.Context, *GetTaskStatusRequest) (*GetTaskStatusResponse, error)
	CancelTask(context.Context, *CancelTaskRequest) (*CancelTaskResponse, error)
}

type ExecuteTaskRequest struct {
	Prompt  string
	Context map[string]string
	Tools   []string
}

type ExecuteTaskResponse struct {
	TaskId  string
	Status  string
	Message string
}

type GetTaskStatusRequest struct {
	TaskId string
}

type GetTaskStatusResponse struct {
	TaskId  string
	Status  string
	Result  string
	Events  []string
}

type CancelTaskRequest struct {
	TaskId string
}

type CancelTaskResponse struct {
	TaskId  string
	Status  string
	Message string
}

type GRPCServer struct {
	tasks   map[string]*Task
	tasksMu sync.RWMutex
}

func (s *GRPCServer) ExecuteTask(ctx context.Context, req *ExecuteTaskRequest) (*ExecuteTaskResponse, error) {
	log.Printf("Received ExecuteTask request with prompt: %s", req.Prompt)

	taskID := uuid.New().String()

	task := &Task{
		ID:      taskID,
		Status:  "running",
		Prompt:  req.Prompt,
		Context: req.Context,
		Tools:   req.Tools,
		Events:  []string{"Task created"},
	}

	s.tasksMu.Lock()
	s.tasks[taskID] = task
	s.tasksMu.Unlock()

	go s.executeTaskAsync(task)

	return &ExecuteTaskResponse{
		TaskId:  taskID,
		Status:  "accepted",
		Message: "Task submitted for execution",
	}, nil
}

func (s *GRPCServer) GetTaskStatus(ctx context.Context, req *GetTaskStatusRequest) (*GetTaskStatusResponse, error) {
	log.Printf("Received GetTaskStatus request for task: %s", req.TaskId)

	s.tasksMu.RLock()
	task, exists := s.tasks[req.TaskId]
	s.tasksMu.RUnlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "task %s not found", req.TaskId)
	}

	return &GetTaskStatusResponse{
		TaskId:  task.ID,
		Status:  task.Status,
		Result:  task.Result,
		Events:  task.Events,
	}, nil
}

func (s *GRPCServer) CancelTask(ctx context.Context, req *CancelTaskRequest) (*CancelTaskResponse, error) {
	log.Printf("Received CancelTask request for task: %s", req.TaskId)

	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	task, exists := s.tasks[req.TaskId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "task %s not found", req.TaskId)
	}

	if task.Status == "running" {
		task.Status = "cancelled"
		task.Events = append(task.Events, "Task cancelled")

		return &CancelTaskResponse{
			TaskId:  task.ID,
			Status:  "cancelled",
			Message: "Task cancelled successfully",
		}, nil
	}

	return &CancelTaskResponse{
		TaskId:  task.ID,
		Status:  task.Status,
		Message: "Cannot cancel task with status: " + task.Status,
	}, nil
}

func (s *GRPCServer) executeTaskAsync(task *Task) {
	time.Sleep(2 * time.Second)

	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	if task.Status == "cancelled" {
		return
	}

	task.Status = "completed"
	task.Result = "Task completed successfully: " + task.Prompt
	task.Events = append(task.Events, "Task execution started", "Task execution completed")
}

func RegisterAgentServiceServer(s *grpc.Server, srv AgentServiceServer) {
	s.RegisterService(&_AgentService_serviceDesc, srv)
}

var _AgentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "agent.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteTask",
			Handler:    _AgentService_ExecuteTask_Handler,
		},
		{
			MethodName: "GetTaskStatus",
			Handler:    _AgentService_GetTaskStatus_Handler,
		},
		{
			MethodName: "CancelTask",
			Handler:    _AgentService_CancelTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent.proto",
}

func _AgentService_ExecuteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).ExecuteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/agent.AgentService/ExecuteTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).ExecuteTask(ctx, req.(*ExecuteTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_GetTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).GetTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/agent.AgentService/GetTaskStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).GetTaskStatus(ctx, req.(*GetTaskStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_CancelTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).CancelTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/agent.AgentService/CancelTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).CancelTask(ctx, req.(*CancelTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func main() {
	server := grpc.NewServer()

	grpcServer := &GRPCServer{
		tasks: make(map[string]*Task),
	}
	RegisterAgentServiceServer(server, grpcServer)

	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal, stopping server...")
		server.GracefulStop()
	}()

	log.Printf("Starting gRPC server on port %d...", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
