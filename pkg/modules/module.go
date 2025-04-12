package modules

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type BaseModule struct {
	name        string
	description string
	tools       []tools.Tool
	initialized bool
	mutex       sync.RWMutex
}

func NewBaseModule(name, description string) *BaseModule {
	return &BaseModule{
		name:        name,
		description: description,
		tools:       make([]tools.Tool, 0),
	}
}

func (m *BaseModule) Name() string {
	return m.name
}

func (m *BaseModule) Description() string {
	return m.description
}

func (m *BaseModule) Tools() []tools.Tool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	return m.tools
}

func (m *BaseModule) AddTool(tool tools.Tool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.tools = append(m.tools, tool)
}

func (m *BaseModule) Initialize(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.initialized {
		return fmt.Errorf("module already initialized: %s", m.name)
	}
	
	m.initialized = true
	return nil
}

func (m *BaseModule) Cleanup() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if !m.initialized {
		return fmt.Errorf("module not initialized: %s", m.name)
	}
	
	m.initialized = false
	return nil
}

type Registry struct {
	modules map[string]Module
	mutex   sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

func (r *Registry) Register(module Module) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.modules[module.Name()]; exists {
		return fmt.Errorf("module already registered: %s", module.Name())
	}
	
	r.modules[module.Name()] = module
	
	return nil
}

func (r *Registry) Get(name string) (Module, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	module, exists := r.modules[name]
	if !exists {
		return nil, fmt.Errorf("module not found: %s", name)
	}
	
	return module, nil
}

func (r *Registry) List() []Module {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	modules := make([]Module, 0, len(r.modules))
	for _, module := range r.modules {
		modules = append(modules, module)
	}
	
	return modules
}

func (r *Registry) ListNames() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	
	return names
}
