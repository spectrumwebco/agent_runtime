package veigar

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

	securityTools map[string]SecurityTool

	context map[string]interface{}

	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

type FrameworkConfig struct {
	Name string `json:"name"`

	Version string `json:"version"`

	Debug bool `json:"debug"`

	SecurityConfig SecurityConfig `json:"security_config"`
}

type SecurityConfig struct {
	BlockMergeSeverities []string `json:"block_merge_severities"`
	WarnSeverities       []string `json:"warn_severities"`
	InfoSeverities       []string `json:"info_severities"`
	ComplianceFrameworks []string `json:"compliance_frameworks"`
}

type Integration interface {
	Start() error

	Stop() error

	Name() string
}

type SecurityTool func(ctx context.Context, params map[string]interface{}) (interface{}, error)

func NewFramework(config FrameworkConfig) (*Framework, error) {
	ctx, cancel := context.WithCancel(context.Background())

	framework := &Framework{
		Config:        config,
		integrations:  make(map[string]Integration),
		securityTools: make(map[string]SecurityTool),
		context:       make(map[string]interface{}),
		ctx:           ctx,
		cancel:        cancel,
	}

	if config.Debug {
		log.Printf("Veigar Cybersecurity Framework %s (v%s) initialized in debug mode\n",
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

func (f *Framework) RegisterSecurityTool(name string, tool SecurityTool) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.securityTools[name]; exists {
		return fmt.Errorf("security tool with name %s already registered", name)
	}

	f.securityTools[name] = tool

	if f.Config.Debug {
		log.Printf("Registered security tool: %s\n", name)
	}

	return nil
}

func (f *Framework) ExecuteSecurityTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	f.mu.RLock()
	tool, exists := f.securityTools[name]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("security tool with name %s not found", name)
	}

	if f.Config.Debug {
		log.Printf("Executing security tool: %s\n", name)
	}

	result, err := tool(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error executing security tool %s: %w", name, err)
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

func (f *Framework) PublishSecurityEvent(eventType string, eventSource string, data interface{}, metadata map[string]interface{}) error {
	log.Printf("Security event published: type=%s, source=%s", eventType, eventSource)
	
	
	return nil
}

func (f *Framework) SubscribeToSecurityEvents(eventType string, channel chan<- interface{}) {
	log.Printf("Security event subscription registered: type=%s", eventType)
	
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

func (f *Framework) IsMergeSeverityBlocked(severity string) bool {
	for _, s := range f.Config.SecurityConfig.BlockMergeSeverities {
		if s == severity {
			return true
		}
	}
	return false
}

func (f *Framework) IsComplianceFrameworkEnabled(framework string) bool {
	for _, f := range f.Config.SecurityConfig.ComplianceFrameworks {
		if f == framework {
			return true
		}
	}
	return false
}
