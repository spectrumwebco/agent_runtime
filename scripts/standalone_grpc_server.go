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

// Import the generated protobuf code
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative agent.proto

// Task represents a task being executed by the agent
type Task struct {
	ID      string
	Status  string
	Prompt  string
	Context map[string]string
	Tools   []string
	Result  string
	Events  []string
}

// AgentService implements the gRPC service defined in agent.proto
type AgentService struct {
	tasks   map[string]*Task
	tasksMu sync.RWMutex
}

// ExecuteTask implements the ExecuteTask RPC method
func (s *AgentService) ExecuteTask(ctx context.Context, req *ExecuteTaskRequest) (*ExecuteTaskResponse, error) {
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

// GetTaskStatus implements the GetTaskStatus RPC method
func (s *AgentService) GetTaskStatus(ctx context.Context, req *GetTaskStatusRequest) (*GetTaskStatusResponse, error) {
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

// CancelTask implements the CancelTask RPC method
func (s *AgentService) CancelTask(ctx context.Context, req *CancelTaskRequest) (*CancelTaskResponse, error) {
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

func (s *AgentService) executeTaskAsync(task *Task) {
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

func main() {
	server := grpc.NewServer()

	agentService := &AgentService{
		tasks: make(map[string]*Task),
	}
	RegisterAgentServiceServer(server, agentService)

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
