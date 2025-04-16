package gomicro

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/microservices/gomicro"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config          *config.Config
	serviceRegistry *gomicro.ServiceRegistry
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config:          cfg,
		serviceRegistry: gomicro.NewServiceRegistry(),
	}
}

func (m *Module) Initialize() error {
	return nil
}

func (m *Module) Start(ctx context.Context) error {
	return nil
}

func (m *Module) Stop(ctx context.Context) error {
	return m.serviceRegistry.StopAll()
}

func (m *Module) CreateService(name, version, address string, metadata map[string]string) (*gomicro.Service, error) {
	service := gomicro.NewService(gomicro.ServiceOptions{
		Name:      name,
		Version:   version,
		Address:   address,
		Metadata:  metadata,
		Timeout:   m.config.GetDuration("microservices.timeout"),
		Retries:   m.config.GetInt("microservices.retries"),
		Namespace: m.config.GetString("microservices.namespace"),
	})

	m.serviceRegistry.Register(service)

	return service, nil
}

func (m *Module) GetService(name string) (*gomicro.Service, error) {
	return m.serviceRegistry.Get(name)
}

func (m *Module) ListServices() []*gomicro.Service {
	return m.serviceRegistry.List()
}

func (m *Module) RemoveService(name string) {
	m.serviceRegistry.Remove(name)
}

func (m *Module) Call(ctx context.Context, serviceName, endpoint string, req, rsp interface{}) error {
	service, err := m.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	return service.Call(ctx, serviceName, endpoint, req, rsp)
}

func (m *Module) Publish(ctx context.Context, topic string, msg interface{}) error {
	services := m.ListServices()
	if len(services) == 0 {
		return fmt.Errorf("no services available to publish message")
	}

	return services[0].Publish(ctx, topic, msg)
}

func (m *Module) RegisterHandler(serviceName, handlerName string, handler interface{}) error {
	service, err := m.GetService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	return service.RegisterHandler(handlerName, handler)
}
