package service

import (
	"context"
	"fmt"
	"time"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

type Service struct {
	micro.Service
	Name    string
	Version string
}

type ServiceOption func(*ServiceOptions)

type ServiceOptions struct {
	Name        string
	Version     string
	Description string
	Metadata    map[string]string
	Registry    registry.Registry
	Address     string
	Timeout     time.Duration
}

func WithName(name string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Name = name
	}
}

func WithVersion(version string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Version = version
	}
}

func WithDescription(description string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Description = description
	}
}

func WithMetadata(metadata map[string]string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Metadata = metadata
	}
}

func WithRegistry(registry registry.Registry) ServiceOption {
	return func(o *ServiceOptions) {
		o.Registry = registry
	}
}

func WithAddress(address string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Address = address
	}
}

func WithTimeout(timeout time.Duration) ServiceOption {
	return func(o *ServiceOptions) {
		o.Timeout = timeout
	}
}

func NewService(opts ...ServiceOption) (*Service, error) {
	options := &ServiceOptions{
		Name:        "kled.service",
		Version:     "latest",
		Description: "Kled.io Framework Service",
		Metadata:    make(map[string]string),
		Timeout:     time.Second * 5,
	}

	for _, o := range opts {
		o(options)
	}

	service := micro.NewService(
		micro.Name(options.Name),
		micro.Version(options.Version),
		micro.Metadata(options.Metadata),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)

	if options.Registry != nil {
		service.Init(
			micro.Registry(options.Registry),
		)
	}

	if options.Address != "" {
		service.Init(
			micro.Address(options.Address),
		)
	}

	return &Service{
		Service: service,
		Name:    options.Name,
		Version: options.Version,
	}, nil
}

func (s *Service) Run() error {
	return s.Service.Run()
}

func (s *Service) Client() client.Client {
	return s.Service.Client()
}

func (s *Service) Server() server.Server {
	return s.Service.Server()
}

func (s *Service) String() string {
	return fmt.Sprintf("%s@%s", s.Name, s.Version)
}

func (s *Service) Context() context.Context {
	return s.Service.Options().Context
}
