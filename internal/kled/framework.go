package kled

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Framework struct {
	Config FrameworkConfig

	integrations map[string]Integration

	tools map[string]Tool

	context map[string]interface{}

	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

type FrameworkConfig struct {
	Name string `json:"name"`

	Version string `json:"version"`

	Debug bool `json:"debug"`
}

type Integration interface {
	Start() error

	Stop() error

	Name() string
}

type Tool func(ctx context.Context, params map[string]interface{}) (interface{}, error)

func NewFramework(config FrameworkConfig) (*Framework, error) {
	ctx, cancel := context.WithCancel(context.Background())

	framework := &Framework{
		Config:       config,
		integrations: make(map[string]Integration),
		tools:        make(map[string]Tool),
		context:      make(map[string]interface{}),
		ctx:          ctx,
		cancel:       cancel,
	}


	if config.Debug {
		log.Printf("Kled.io Framework %s (v%s) initialized in debug mode\n", 
			config.Name, config.Version)
	}

	return framework, nil
}

func (f *Framework) RegisterIntegration(integration Integration) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	name := integration.Name()
	if _, exists := f.integrations[name]; exists {
		return fmt.Errorf("integration with name %s already registered", name)
	}

	f.integrations[name] = integration
	
	if f.Config.Debug {
		log.Printf("Registered integration: %s\n", name)
	}

	return nil
}

func (f *Framework) GetIntegration(name string) (Integration, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	integration, exists := f.integrations[name]
	if !exists {
		return nil, fmt.Errorf("integration with name %s not found", name)
	}

	return integration, nil
}

func (f *Framework) RegisterTool(name string, tool Tool) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.tools[name]; exists {
		return fmt.Errorf("tool with name %s already registered", name)
	}

	f.tools[name] = tool
	
	if f.Config.Debug {
		log.Printf("Registered tool: %s\n", name)
	}

	return nil
}

func (f *Framework) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	f.mu.RLock()
	tool, exists := f.tools[name]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tool with name %s not found", name)
	}

	if f.Config.Debug {
		log.Printf("Executing tool: %s\n", name)
	}

	result, err := tool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error executing tool %s: %w", name, err)
	}

	return result, nil
}

func (f *Framework) Start() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for name, integration := range f.integrations {
		if err := integration.Start(); err != nil {
			return fmt.Errorf("failed to start integration %s: %w", name, err)
		}
		
		if f.Config.Debug {
			log.Printf("Started integration: %s\n", name)
		}
	}

	return nil
}

func (f *Framework) Stop() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var errs []error

	for name, integration := range f.integrations {
		if err := integration.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop integration %s: %w", name, err))
		} else if f.Config.Debug {
			log.Printf("Stopped integration: %s\n", name)
		}
	}


	f.cancel()

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping framework: %v", errs)
	}

	return nil
}

func (f *Framework) SetContext(key string, value interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.context[key] = value
}

func (f *Framework) GetContext(key string) (interface{}, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	value, exists := f.context[key]
	return value, exists
}

func (f *Framework) PublishEvent(eventType string, eventSource string, data interface{}) error {
	log.Printf("Event published: type=%s, source=%s", eventType, eventSource)
	return nil
}

func (f *Framework) SubscribeToEvents(eventType string, channel chan<- interface{}) {
	log.Printf("Event subscription registered: type=%s", eventType)
}

func (f *Framework) SaveState() ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	state := map[string]interface{}{
		"context": f.context,
	}

	return json.Marshal(state)
}

func (f *Framework) LoadState(data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if context, ok := state["context"].(map[string]interface{}); ok {
		f.context = context
	}

	return nil
}

func (f *Framework) Version() string {
	return f.Config.Version
}

func (f *Framework) Name() string {
	return f.Config.Name
}
