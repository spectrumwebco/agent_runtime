package gomicro

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
)

type Service struct {
	service micro.Service
	options ServiceOptions
}

type ServiceOptions struct {
	Name      string
	Version   string
	Address   string
	Registry  registry.Registry
	Metadata  map[string]string
	Timeout   time.Duration
	Retries   int
	Namespace string
}

func NewService(opts ServiceOptions) *Service {
	if opts.Timeout == 0 {
		opts.Timeout = time.Second * 5
	}

	if opts.Retries == 0 {
		opts.Retries = 3
	}

	if opts.Namespace == "" {
		opts.Namespace = "kled"
	}

	serviceOpts := []micro.Option{
		micro.Name(opts.Name),
		micro.Version(opts.Version),
		micro.Metadata(opts.Metadata),
	}

	if opts.Registry != nil {
		serviceOpts = append(serviceOpts, micro.Registry(opts.Registry))
	}

	if opts.Address != "" {
		serviceOpts = append(serviceOpts, micro.Address(opts.Address))
	}

	serviceOpts = append(serviceOpts, micro.Client(
		client.NewClient(
			client.RequestTimeout(opts.Timeout),
			client.Retries(opts.Retries),
		),
	))

	service := micro.NewService(serviceOpts...)

	return &Service{
		service: service,
		options: opts,
	}
}

func (s *Service) Init(opts ...micro.Option) error {
	return s.service.Init(opts...)
}

func (s *Service) Options() micro.Options {
	return s.service.Options()
}

func (s *Service) Client() client.Client {
	return s.service.Client()
}

func (s *Service) Server() server.Server {
	return s.service.Server()
}

func (s *Service) Run() error {
	return s.service.Run()
}

func (s *Service) String() string {
	return s.service.String()
}

func (s *Service) Name() string {
	return s.options.Name
}

func (s *Service) Version() string {
	return s.options.Version
}

func (s *Service) RegisterHandler(name string, handler interface{}, opts ...server.HandlerOption) error {
	return s.service.Server().Handle(
		s.service.Server().NewHandler(handler, opts...),
	)
}

func (s *Service) Call(ctx context.Context, service, endpoint string, req, rsp interface{}, opts ...client.CallOption) error {
	return s.service.Client().Call(ctx, s.createRequest(service, endpoint), req, rsp, opts...)
}

func (s *Service) Publish(ctx context.Context, topic string, msg interface{}, opts ...client.PublishOption) error {
	return s.service.Client().Publish(ctx, s.createRequest("", topic), msg, opts...)
}

func (s *Service) createRequest(service, endpoint string) client.Request {
	if service != "" {
		service = fmt.Sprintf("%s.%s", s.options.Namespace, service)
	}

	return s.service.Client().NewRequest(service, endpoint, nil)
}

func (s *Service) Stop() error {
	return s.service.Server().Stop()
}

type ServiceRegistry struct {
	services map[string]*Service
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*Service),
	}
}

func (sr *ServiceRegistry) Register(service *Service) {
	sr.services[service.Name()] = service
}

func (sr *ServiceRegistry) Get(name string) (*Service, error) {
	service, ok := sr.services[name]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

func (sr *ServiceRegistry) List() []*Service {
	services := make([]*Service, 0, len(sr.services))
	for _, service := range sr.services {
		services = append(services, service)
	}

	return services
}

func (sr *ServiceRegistry) Remove(name string) {
	delete(sr.services, name)
}

func (sr *ServiceRegistry) StopAll() error {
	for _, service := range sr.services {
		err := service.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}
