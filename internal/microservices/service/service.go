package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/kledframework"
)

// Service represents a microservice
type Service struct {
	kledService *kledframework.Service
}

// ServiceOption represents an option for a microservice
type ServiceOption func(*ServiceOptions)

// ServiceOptions represents options for a microservice
type ServiceOptions struct {
	Name     string
	Version  string
	Address  string
	Metadata map[string]string
	Timeout  time.Duration
}

// WithName sets the name of the service
func WithName(name string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Name = name
	}
}

// WithVersion sets the version of the service
func WithVersion(version string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Version = version
	}
}

// WithAddress sets the address of the service
func WithAddress(address string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Address = address
	}
}

// WithMetadata sets the metadata of the service
func WithMetadata(metadata map[string]string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Metadata = metadata
	}
}

// WithTimeout sets the timeout of the service
func WithTimeout(timeout time.Duration) ServiceOption {
	return func(o *ServiceOptions) {
		o.Timeout = timeout
	}
}

// NewService creates a new microservice
func NewService(opts ...ServiceOption) (*Service, error) {
	options := &ServiceOptions{
		Name:     "service",
		Version:  "latest",
		Address:  ":8080",
		Metadata: map[string]string{},
		Timeout:  time.Second * 5,
	}

	for _, o := range opts {
		o(options)
	}

	kledService := kledframework.NewService(kledframework.ServiceConfig{
		Name:    options.Name,
		Version: options.Version,
		Address: options.Address,
		Timeout: options.Timeout,
	})

	return &Service{
		kledService: kledService,
	}, nil
}

// Run runs the microservice
func (s *Service) Run() error {
	return s.kledService.Start()
}

// Server returns the server of the microservice
func (s *Service) Server() interface{} {
	return s.kledService.Router
}
