package microservices

import (
	"fmt"
	"time"

	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/registry/mdns"
	"go-micro.dev/v4/registry/memory"
)

// RegistryType represents the type of registry
type RegistryType string

const (
	// RegistryTypeMDNS is the MDNS registry type
	RegistryTypeMDNS RegistryType = "mdns"
	// RegistryTypeMemory is the memory registry type
	RegistryTypeMemory RegistryType = "memory"
)

// RegistryConfig contains configuration for a registry
type RegistryConfig struct {
	Type      RegistryType
	Namespace string
	TTL       time.Duration
}

// NewRegistry creates a new registry
func NewRegistry(config RegistryConfig) (registry.Registry, error) {
	// Set default values
	if config.TTL == 0 {
		config.TTL = time.Minute
	}
	if config.Namespace == "" {
		config.Namespace = "kled"
	}

	// Create registry options
	options := []registry.Option{
		registry.Timeout(time.Second * 5),
		registry.Namespace(config.Namespace),
	}

	// Create the registry based on the type
	switch config.Type {
	case RegistryTypeMDNS:
		return mdns.NewRegistry(options...), nil
	case RegistryTypeMemory:
		return memory.NewRegistry(options...), nil
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", config.Type)
	}
}
