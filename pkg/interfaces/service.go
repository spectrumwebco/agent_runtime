package interfaces

import (
	"context"
	"time"
)

type Service interface {
	Start(ctx context.Context) error

	Stop() error

	Status() ServiceStatus

	Name() string

	Version() string

	Health() ServiceHealth
}

type ServiceStatus string

const (
	ServiceStatusUnknown ServiceStatus = "unknown"

	ServiceStatusStarting ServiceStatus = "starting"

	ServiceStatusRunning ServiceStatus = "running"

	ServiceStatusStopping ServiceStatus = "stopping"

	ServiceStatusStopped ServiceStatus = "stopped"

	ServiceStatusFailed ServiceStatus = "failed"
)

type ServiceHealth struct {
	Status ServiceHealthStatus

	Message string

	LastChecked time.Time

	Details map[string]interface{}
}

type ServiceHealthStatus string

const (
	ServiceHealthStatusUnknown ServiceHealthStatus = "unknown"

	ServiceHealthStatusHealthy ServiceHealthStatus = "healthy"

	ServiceHealthStatusDegraded ServiceHealthStatus = "degraded"

	ServiceHealthStatusUnhealthy ServiceHealthStatus = "unhealthy"
)

type ServiceRegistry interface {
	Register(service Service) error

	Unregister(serviceName string) error

	Get(serviceName string) (Service, error)

	List() []Service

	Start(ctx context.Context) error

	Stop() error
}

type ServiceFactory interface {
	Create(serviceType string, config map[string]interface{}) (Service, error)
}

type ServiceManager interface {
	ServiceRegistry

	Start(ctx context.Context) error

	Stop() error

	AddDependency(serviceName, dependsOn string) error

	RemoveDependency(serviceName, dependsOn string) error

	GetDependencies(serviceName string) ([]string, error)

	GetDependents(serviceName string) ([]string, error)
}

type ServiceConfig struct {
	Name string `json:"name" yaml:"name"`

	Type string `json:"type" yaml:"type"`

	Version string `json:"version" yaml:"version"`

	Enabled bool `json:"enabled" yaml:"enabled"`

	Dependencies []string `json:"dependencies" yaml:"dependencies"`

	Config map[string]interface{} `json:"config" yaml:"config"`
}

type ServiceEvent struct {
	Type ServiceEventType `json:"type" yaml:"type"`

	ServiceName string `json:"service_name" yaml:"service_name"`

	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`

	Data map[string]interface{} `json:"data" yaml:"data"`
}

type ServiceEventType string

const (
	ServiceEventTypeRegistered ServiceEventType = "registered"

	ServiceEventTypeUnregistered ServiceEventType = "unregistered"

	ServiceEventTypeStarted ServiceEventType = "started"

	ServiceEventTypeStopped ServiceEventType = "stopped"

	ServiceEventTypeFailed ServiceEventType = "failed"

	ServiceEventTypeHealthChanged ServiceEventType = "health_changed"
)

type ServiceEventListener interface {
	OnEvent(event ServiceEvent)
}

type ServiceEventEmitter interface {
	AddListener(listener ServiceEventListener)

	RemoveListener(listener ServiceEventListener)

	EmitEvent(event ServiceEvent)
}
