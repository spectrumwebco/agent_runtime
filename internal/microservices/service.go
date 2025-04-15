package microservices

import (
	"context"
	"fmt"
	"time"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

// Service represents a microservice
type Service struct {
	Name    string
	Version string
	Service micro.Service
}

// ServiceConfig contains configuration for a microservice
type ServiceConfig struct {
	Name      string
	Version   string
	Address   string
	Registry  registry.Registry
	Timeout   time.Duration
	Retries   int
	Namespace string
}

// NewService creates a new microservice
func NewService(config ServiceConfig) (*Service, error) {
	// Set default values
	if config.Timeout == 0 {
		config.Timeout = time.Second * 5
	}
	if config.Retries == 0 {
		config.Retries = 3
	}
	if config.Namespace == "" {
		config.Namespace = "kled"
	}

	// Create service options
	options := []micro.Option{
		micro.Name(config.Name),
		micro.Version(config.Version),
		micro.Address(config.Address),
		micro.Client(client.NewClient(
			client.RequestTimeout(config.Timeout),
			client.Retries(config.Retries),
		)),
		micro.Server(server.NewServer(
			server.Name(config.Name),
			server.Version(config.Version),
			server.Address(config.Address),
		)),
	}

	// Add registry if provided
	if config.Registry != nil {
		options = append(options, micro.Registry(config.Registry))
	}

	// Create the service
	service := micro.NewService(options...)

	return &Service{
		Name:    config.Name,
		Version: config.Version,
		Service: service,
	}, nil
}

// Start starts the microservice
func (s *Service) Start() error {
	return s.Service.Run()
}

// Stop stops the microservice
func (s *Service) Stop() error {
	return nil
}

// Client returns the service client
func (s *Service) Client() client.Client {
	return s.Service.Client()
}

// Server returns the service server
func (s *Service) Server() server.Server {
	return s.Service.Server()
}

// Call calls a method on a service
func (s *Service) Call(ctx context.Context, service, endpoint string, req, rsp interface{}, opts ...client.CallOption) error {
	return s.Client().Call(ctx, s.Client().NewRequest(service, endpoint, req), rsp, opts...)
}

// Publish publishes a message to a topic
func (s *Service) Publish(ctx context.Context, topic string, msg interface{}, opts ...client.PublishOption) error {
	return s.Client().Publish(ctx, s.Client().NewMessage(topic, msg), opts...)
}

// Subscribe subscribes to a topic
func (s *Service) Subscribe(topic string, handler interface{}, opts ...client.SubscribeOption) error {
	_, err := s.Client().Subscribe(context.Background(), topic, handler, opts...)
	return err
}
