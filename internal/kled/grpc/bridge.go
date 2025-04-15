package grpc

import (
	"context"
	"log"
	"net"

	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
	pb "github.com/spectrumwebco/agent_runtime/protos/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BridgeServer struct {
	pb.UnimplementedAgentBridgeServer
	eventStream  *eventstream.EventStream
	stateManager *statemanager.StateManager
}

func NewBridgeServer(eventStream *eventstream.EventStream, stateManager *statemanager.StateManager) *BridgeServer {
	return &BridgeServer{
		eventStream:  eventStream,
		stateManager: stateManager,
	}
}

func (s *BridgeServer) SendEvent(ctx context.Context, req *pb.SendEventRequest) (*pb.SendEventResponse, error) {
	data := make(map[string]interface{})
	for k, v := range req.Data {
		data[k] = v
	}

	s.eventStream.Publish(req.EventType, data)

	return &pb.SendEventResponse{
		Success: true,
		Message: "Event sent successfully",
	}, nil
}

func (s *BridgeServer) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	state, err := s.stateManager.GetState(req.StateType, req.StateId)
	if err != nil {
		return &pb.GetStateResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	stateData := make(map[string]string)
	for k, v := range state {
		if strVal, ok := v.(string); ok {
			stateData[k] = strVal
		}
	}

	return &pb.GetStateResponse{
		Success: true,
		Message: "State retrieved successfully",
		State:   stateData,
	}, nil
}

func (s *BridgeServer) SetState(ctx context.Context, req *pb.SetStateRequest) (*pb.SetStateResponse, error) {
	state := make(map[string]interface{})
	for k, v := range req.State {
		state[k] = v
	}

	err := s.stateManager.SetState(req.StateType, req.StateId, state)
	if err != nil {
		return &pb.SetStateResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SetStateResponse{
		Success: true,
		Message: "State set successfully",
	}, nil
}

func StartGRPCServer(eventStream *eventstream.EventStream, stateManager *statemanager.StateManager, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	bridgeServer := NewBridgeServer(eventStream, stateManager)
	pb.RegisterAgentBridgeServer(server, bridgeServer)

	log.Printf("Starting gRPC server on port %s", port)
	return server.Serve(lis)
}

func GetGRPCClient(address string) (pb.AgentBridgeClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewAgentBridgeClient(conn)
	return client, conn, nil
}
