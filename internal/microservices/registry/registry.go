package registry

import (
	"time"

	"github.com/micro/go-micro/registry"
)

type Registry struct {
	registry.Registry
}

type RegistryOption func(*RegistryOptions)

type RegistryOptions struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig map[string]interface{}
}

func WithAddrs(addrs ...string) RegistryOption {
	return func(o *RegistryOptions) {
		o.Addrs = addrs
	}
}

func WithTimeout(timeout time.Duration) RegistryOption {
	return func(o *RegistryOptions) {
		o.Timeout = timeout
	}
}

func WithSecure(secure bool) RegistryOption {
	return func(o *RegistryOptions) {
		o.Secure = secure
	}
}

func WithTLSConfig(tlsConfig map[string]interface{}) RegistryOption {
	return func(o *RegistryOptions) {
		o.TLSConfig = tlsConfig
	}
}

func NewRegistry(opts ...RegistryOption) *Registry {
	options := &RegistryOptions{
		Addrs:   []string{},
		Timeout: time.Second * 5,
		Secure:  false,
	}

	for _, o := range opts {
		o(options)
	}

	r := registry.NewRegistry(
		registry.Addrs(options.Addrs...),
		registry.Timeout(options.Timeout),
	)

	if options.Secure {
		r = registry.NewRegistry(
			registry.Addrs(options.Addrs...),
			registry.Timeout(options.Timeout),
			registry.Secure(true),
		)
	}

	return &Registry{
		Registry: r,
	}
}

func (r *Registry) GetService(name string) ([]*registry.Service, error) {
	return r.Registry.GetService(name)
}

func (r *Registry) ListServices() ([]*registry.Service, error) {
	return r.Registry.ListServices()
}

func (r *Registry) Register(service *registry.Service, opts ...registry.RegisterOption) error {
	return r.Registry.Register(service, opts...)
}

func (r *Registry) Deregister(service *registry.Service, opts ...registry.DeregisterOption) error {
	return r.Registry.Deregister(service, opts...)
}
