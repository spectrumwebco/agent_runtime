package kledframework

import (
	"fmt"
	"sync"
	"time"
)

// ServiceInfo contains information about a service
type ServiceInfo struct {
	Name      string
	Version   string
	Address   string
	Metadata  map[string]string
	Timestamp time.Time
}

// Registry represents a service registry
type Registry struct {
	services map[string]*ServiceInfo
	mutex    sync.RWMutex
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*ServiceInfo),
	}
}

// Register registers a service with the registry
func (r *Registry) Register(info *ServiceInfo) error {
	if info.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if info.Address == "" {
		return fmt.Errorf("service address cannot be empty")
	}

	info.Timestamp = time.Now()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.services[info.Name] = info
	return nil
}

// Deregister removes a service from the registry
func (r *Registry) Deregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.services[name]; !exists {
		return fmt.Errorf("service %s not found", name)
	}

	delete(r.services, name)
	return nil
}

// GetService gets a service from the registry
func (r *Registry) GetService(name string) (*ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	service, exists := r.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// ListServices lists all services in the registry
func (r *Registry) ListServices() []*ServiceInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	services := make([]*ServiceInfo, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}

	return services
}
