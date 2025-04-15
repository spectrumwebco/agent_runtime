package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/internal/agent"
	pb "github.com/spectrumwebco/agent_runtime/internal/server/proto"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
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

type GRPCServer struct {
	pb.UnimplementedAgentServiceServer
	server   *grpc.Server
	config   *config.Config
	agent    *agent.Agent
	tasks    map[string]*Task
	tasksMu  sync.RWMutex
	listener net.Listener
}

func NewGRPCServer(cfg *config.Config, agentInstance *agent.Agent) (*GRPCServer, error) {
	server := grpc.NewServer()
	
	grpcServer := &GRPCServer{
		server: server,
		config: cfg,
		agent:  agentInstance,
		tasks:  make(map[string]*Task),
	}
	
	pb.RegisterAgentServiceServer(server, grpcServer)
	
	return grpcServer, nil
}

func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.GRPCHost, s.config.Server.GRPCPort)
	
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	
	log.Printf("gRPC server listening on %s", addr)
	
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	
	return nil
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
	
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *GRPCServer) ExecuteTask(ctx context.Context, req *pb.ExecuteTaskRequest) (*pb.ExecuteTaskResponse, error) {
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

	return &pb.ExecuteTaskResponse{
		TaskId:  taskID,
		Status:  "accepted",
		Message: "Task submitted for execution",
	}, nil
}

func (s *GRPCServer) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusRequest) (*pb.GetTaskStatusResponse, error) {
	log.Printf("Received GetTaskStatus request for task: %s", req.TaskId)

	s.tasksMu.RLock()
	task, exists := s.tasks[req.TaskId]
	s.tasksMu.RUnlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "task %s not found", req.TaskId)
	}

	return &pb.GetTaskStatusResponse{
		TaskId:  task.ID,
		Status:  task.Status,
		Result:  task.Result,
		Events:  task.Events,
	}, nil
}

func (s *GRPCServer) CancelTask(ctx context.Context, req *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
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

		return &pb.CancelTaskResponse{
			TaskId:  task.ID,
			Status:  "cancelled",
			Message: "Task cancelled successfully",
		}, nil
	}

	return &pb.CancelTaskResponse{
		TaskId:  task.ID,
		Status:  task.Status,
		Message: "Cannot cancel task with status: " + task.Status,
	}, nil
}

func (s *GRPCServer) executeTaskAsync(task *Task) {
	result, err := s.agent.Execute(context.Background(), task.Prompt)
	
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	if task.Status == "cancelled" {
		return
	}

	if err != nil {
		task.Status = "failed"
		task.Result = fmt.Sprintf("Task execution failed: %v", err)
		task.Events = append(task.Events, "Task execution failed")
	} else {
		task.Status = "completed"
		task.Result = fmt.Sprintf("%v", result)
		task.Events = append(task.Events, "Task execution started", "Task execution completed")
	}
}
