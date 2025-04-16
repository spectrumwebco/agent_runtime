package kledframework

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/agent"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/models"
)

type KledFramework struct {
	config *Config

	registry       *ServiceRegistry
	agentFactory   *agent.AgentFactory
	eventStream    EventStream
	stateManager   StateManager
	graphExecutor  *langraph.Executor
	djangoBridge   DjangoBridge
	reactBridge    ReactBridge
	serviceManager *ServiceManager

	isRunning bool
	mutex     sync.RWMutex
	stopCh    chan struct{}
}

type Config struct {
	Name string `json:"name" yaml:"name"`

	Version string `json:"version" yaml:"version"`

	LogLevel string `json:"log_level" yaml:"log_level"`

	EventStreamConfig map[string]interface{} `json:"event_stream_config" yaml:"event_stream_config"`

	StateManagerConfig map[string]interface{} `json:"state_manager_config" yaml:"state_manager_config"`

	DjangoBridgeConfig map[string]interface{} `json:"django_bridge_config" yaml:"django_bridge_config"`

	ReactBridgeConfig map[string]interface{} `json:"react_bridge_config" yaml:"react_bridge_config"`

	AgentConfigs map[string]map[string]interface{} `json:"agent_configs" yaml:"agent_configs"`

	ServiceConfigs map[string]map[string]interface{} `json:"service_configs" yaml:"service_configs"`
}

type EventStream interface {
	AddEvent(event *models.Event) error
}

type StateManager interface {
	GetState(id string) (*statemodels.State, error)
	UpdateState(id string, state *statemodels.State) error
}

type DjangoBridge interface {
	ExecutePythonCode(ctx context.Context, code string, timeout time.Duration) (map[string]interface{}, error)
	QueryDjangoModel(ctx context.Context, modelName string, filters map[string]interface{}) ([]map[string]interface{}, error)
	CreateDjangoModel(ctx context.Context, modelName string, data map[string]interface{}) (map[string]interface{}, error)
	UpdateDjangoModel(ctx context.Context, modelName string, id int, data map[string]interface{}) error
	DeleteDjangoModel(ctx context.Context, modelName string, id int) error
	ExecuteDjangoManagementCommand(ctx context.Context, command string, args []string) (string, error)
	GetDjangoSettings(ctx context.Context) (map[string]interface{}, error)
	CreateDjangoAgent(ctx context.Context, id, name, role string) (map[string]interface{}, error)
	Close() error
}

type ReactBridge interface {
	SendEvent(ctx context.Context, event map[string]interface{}) error
	SubscribeToEvents(ctx context.Context) (<-chan map[string]interface{}, error)
	UpdateState(ctx context.Context, state map[string]interface{}) error
	GetState(ctx context.Context) (map[string]interface{}, error)
	Close() error
}

func NewKledFramework(config *Config) (*KledFramework, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	framework := &KledFramework{
		config:    config,
		isRunning: false,
		mutex:     sync.RWMutex{},
		stopCh:    make(chan struct{}),
	}

	if err := framework.initComponents(); err != nil {
		return nil, fmt.Errorf("failed to initialize components: %v", err)
	}

	return framework, nil
}

func (f *KledFramework) initComponents() error {
	f.registry = NewServiceRegistry()

	f.serviceManager = NewServiceManager(f.registry)



	f.agentFactory = agent.NewAgentFactory(f.eventStream, f.stateManager)

	f.graphExecutor = langraph.NewExecutor()



	return nil
}

func (f *KledFramework) Start(ctx context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.isRunning {
		return fmt.Errorf("framework is already running")
	}

	if err := f.serviceManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start service manager: %v", err)
	}

	if err := f.agentFactory.StartAllAgents(ctx); err != nil {
		return fmt.Errorf("failed to start agents: %v", err)
	}

	if err := f.graphExecutor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start graph executor: %v", err)
	}

	f.isRunning = true
	go f.run(ctx)

	log.Printf("Kled.io framework started: %s (version %s)", f.config.Name, f.config.Version)

	return nil
}

func (f *KledFramework) Stop() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.isRunning {
		return fmt.Errorf("framework is not running")
	}

	close(f.stopCh)

	if err := f.graphExecutor.Stop(); err != nil {
		log.Printf("Failed to stop graph executor: %v", err)
	}

	if err := f.agentFactory.StopAllAgents(); err != nil {
		log.Printf("Failed to stop agents: %v", err)
	}

	if err := f.serviceManager.Stop(); err != nil {
		log.Printf("Failed to stop service manager: %v", err)
	}

	f.isRunning = false

	log.Printf("Kled.io framework stopped: %s", f.config.Name)

	return nil
}

func (f *KledFramework) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			f.Stop()
			return
		case <-f.stopCh:
			return
		case <-ticker.C:
			f.checkHealth()
		}
	}
}

func (f *KledFramework) checkHealth() {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	if f.serviceManager != nil {
		health := f.serviceManager.Health()
		if health.Status != ServiceHealthStatusHealthy {
			log.Printf("Service manager health: %s - %s", health.Status, health.Message)
		}
	}




}

func (f *KledFramework) GetAgentFactory() *agent.AgentFactory {
	return f.agentFactory
}

func (f *KledFramework) GetServiceRegistry() *ServiceRegistry {
	return f.registry
}

func (f *KledFramework) GetServiceManager() *ServiceManager {
	return f.serviceManager
}

func (f *KledFramework) GetGraphExecutor() *langraph.Executor {
	return f.graphExecutor
}

func (f *KledFramework) GetDjangoBridge() DjangoBridge {
	return f.djangoBridge
}

func (f *KledFramework) GetReactBridge() ReactBridge {
	return f.reactBridge
}

func (f *KledFramework) CreateAgent(id, name, role string, capabilities []string) (*agent.Agent, error) {
	return f.agentFactory.CreateAgent(id, name, role, capabilities)
}

func (f *KledFramework) GetAgent(id string) (*agent.Agent, error) {
	return f.agentFactory.GetAgent(id)
}

func (f *KledFramework) CreateService(serviceType string, config map[string]interface{}) (Service, error) {
	return f.serviceManager.CreateService(serviceType, config)
}

func (f *KledFramework) GetService(serviceName string) (Service, error) {
	return f.serviceManager.GetService(serviceName)
}

func (f *KledFramework) ExecutePythonCode(ctx context.Context, code string, timeout time.Duration) (map[string]interface{}, error) {
	if f.djangoBridge == nil {
		return nil, fmt.Errorf("Django bridge is not initialized")
	}

	return f.djangoBridge.ExecutePythonCode(ctx, code, timeout)
}

func (f *KledFramework) SendEventToReact(ctx context.Context, event map[string]interface{}) error {
	if f.reactBridge == nil {
		return fmt.Errorf("React bridge is not initialized")
	}

	return f.reactBridge.SendEvent(ctx, event)
}

func (f *KledFramework) UpdateReactState(ctx context.Context, state map[string]interface{}) error {
	if f.reactBridge == nil {
		return fmt.Errorf("React bridge is not initialized")
	}

	return f.reactBridge.UpdateState(ctx, state)
}

func (f *KledFramework) GetReactState(ctx context.Context) (map[string]interface{}, error) {
	if f.reactBridge == nil {
		return nil, fmt.Errorf("React bridge is not initialized")
	}

	return f.reactBridge.GetState(ctx)
}

func (f *KledFramework) CreateMultiAgentSystem(ctx context.Context, config *langraph.MultiAgentSystemConfig) (*langraph.MultiAgentSystem, error) {
	return langraph.NewMultiAgentSystem(config, f.graphExecutor)
}

func (f *KledFramework) ExecuteGraph(ctx context.Context, graph *langraph.Graph) (map[string]interface{}, error) {
	return f.graphExecutor.ExecuteGraph(ctx, graph)
}

func (f *KledFramework) EmitEvent(eventType models.EventType, eventSource models.EventSource, action string, metadata map[string]string) error {
	if f.eventStream == nil {
		return fmt.Errorf("event stream is not initialized")
	}

	event := &models.Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    eventSource,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": action,
		},
		Metadata: metadata,
	}

	return f.eventStream.AddEvent(event)
}

func (f *KledFramework) UpdateState(id string, state *statemodels.State) error {
	if f.stateManager == nil {
		return fmt.Errorf("state manager is not initialized")
	}

	return f.stateManager.UpdateState(id, state)
}

func (f *KledFramework) GetState(id string) (*statemodels.State, error) {
	if f.stateManager == nil {
		return nil, fmt.Errorf("state manager is not initialized")
	}

	return f.stateManager.GetState(id)
}
