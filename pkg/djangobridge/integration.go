package djangobridge

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/klediosdk"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/langsmith"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type DjangoIntegration struct {
	config *IntegrationConfig

	bridge *Bridge

	sdk *klediosdk.SDK

	multiAgentSystem *langraph.MultiAgentSystem

	langSmithIntegration *langsmith.LangGraphIntegration

	stateManager sharedstate.Manager

	ctx    context.Context
	cancel context.CancelFunc

	mu sync.RWMutex

	connected bool
}

type IntegrationConfig struct {
	BridgeConfig *BridgeConfig

	SDKConfig *klediosdk.Config

	MultiAgentSystemConfig *langraph.MultiAgentSystemConfig

	LangSmithConfig *langsmith.LangSmithConfig

	StateManagerConfig *sharedstate.Config

	ConnectionTimeout time.Duration

	ReconnectionAttempts int

	ReconnectionDelay time.Duration
}

func DefaultIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{
		BridgeConfig:         DefaultBridgeConfig(),
		SDKConfig:            &klediosdk.Config{},
		MultiAgentSystemConfig: &langraph.MultiAgentSystemConfig{
			Name: "Django Integration Multi-Agent System",
		},
		LangSmithConfig:      &langsmith.LangSmithConfig{},
		StateManagerConfig:   &sharedstate.Config{},
		ConnectionTimeout:    30 * time.Second,
		ReconnectionAttempts: 5,
		ReconnectionDelay:    5 * time.Second,
	}
}

func NewDjangoIntegration(config *IntegrationConfig) (*DjangoIntegration, error) {
	if config == nil {
		config = DefaultIntegrationConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	integration := &DjangoIntegration{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}

	bridge, err := NewBridge(config.BridgeConfig)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize bridge: %w", err)
	}
	integration.bridge = bridge

	stateManager, err := sharedstate.NewManager(config.StateManagerConfig)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize state manager: %w", err)
	}
	integration.stateManager = stateManager

	sdkConfig := config.SDKConfig
	if sdkConfig == nil {
		sdkConfig = &klediosdk.Config{}
	}
	sdkConfig.StateManagerConfig = config.StateManagerConfig

	sdk, err := klediosdk.NewSDK(sdkConfig)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize SDK: %w", err)
	}
	integration.sdk = sdk

	if config.LangSmithConfig != nil && config.LangSmithConfig.Enabled {
		langSmithIntegration, err := langsmith.NewLangGraphIntegration(config.LangSmithConfig)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to initialize LangSmith integration: %w", err)
		}
		integration.langSmithIntegration = langSmithIntegration
	}

	if config.MultiAgentSystemConfig != nil {
		var executor *langraph.Executor
		if integration.langSmithIntegration != nil {
			tracer := integration.langSmithIntegration.CreateTracer()
			executor = langraph.NewExecutor(langraph.WithTracer(tracer))
		} else {
			executor = langraph.NewExecutor()
		}

		multiAgentSystem, err := langraph.NewMultiAgentSystem(config.MultiAgentSystemConfig, executor)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to initialize multi-agent system: %w", err)
		}
		integration.multiAgentSystem = multiAgentSystem

		if integration.langSmithIntegration != nil {
			_, err = integration.langSmithIntegration.RegisterMultiAgentSystem(ctx, multiAgentSystem)
			if err != nil {
				log.Printf("Failed to register multi-agent system with LangSmith: %v", err)
			}
		}
	}

	return integration, nil
}

func (i *DjangoIntegration) Connect(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.connected {
		return nil
	}

	connectCtx, cancel := context.WithTimeout(ctx, i.config.ConnectionTimeout)
	defer cancel()

	err := i.bridge.Connect(connectCtx)
	if err != nil {
		return fmt.Errorf("failed to connect to bridge: %w", err)
	}

	err = i.sdk.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize SDK: %w", err)
	}

	i.connected = true
	return nil
}

func (i *DjangoIntegration) Disconnect(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.connected {
		return nil
	}

	err := i.bridge.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect from bridge: %w", err)
	}

	err = i.sdk.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown SDK: %w", err)
	}

	i.connected = false
	return nil
}

func (i *DjangoIntegration) IsConnected() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.connected
}

func (i *DjangoIntegration) ExecuteTask(ctx context.Context, task *klediosdk.Task) (*klediosdk.TaskResult, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return nil, errors.New("integration is not connected")
	}

	task.Runtime = klediosdk.RuntimeTypePython

	return i.sdk.ExecuteTask(ctx, task)
}

func (i *DjangoIntegration) ExecuteAgentTask(ctx context.Context, task *langraph.Task) (*langraph.TaskResult, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return nil, errors.New("integration is not connected")
	}

	if i.multiAgentSystem == nil {
		return nil, errors.New("multi-agent system is not initialized")
	}

	var taskCtx context.Context
	if i.langSmithIntegration != nil && task.AssignedTo != "" {
		var taskRunID string
		var err error
		taskCtx, taskRunID, err = i.langSmithIntegration.TraceAgentTask(
			ctx,
			task.AssignedTo,
			task.ID,
			"task",
			task.Description,
			task.Input,
		)
		if err != nil {
			log.Printf("Failed to trace agent task with LangSmith: %v", err)
			taskCtx = ctx
		}

		defer func() {
			if taskRunID != "" {
				_ = i.langSmithIntegration.EndAgentTask(
					taskCtx,
					task.ID,
					nil,
					nil,
				)
			}
		}()
	} else {
		taskCtx = ctx
	}

	result, err := i.multiAgentSystem.ExecuteTask(taskCtx, task)

	if i.langSmithIntegration != nil && task.AssignedTo != "" && result != nil {
		var traceErr error
		if err != nil {
			traceErr = err
		}
		_ = i.langSmithIntegration.EndAgentTask(
			taskCtx,
			task.ID,
			result.Output,
			traceErr,
		)
	}

	return result, err
}

func (i *DjangoIntegration) GetState(ctx context.Context, key string) (interface{}, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return nil, errors.New("integration is not connected")
	}

	return i.stateManager.Get(ctx, key)
}

func (i *DjangoIntegration) SetState(ctx context.Context, key string, value interface{}) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return errors.New("integration is not connected")
	}

	return i.stateManager.Set(ctx, key, value)
}

func (i *DjangoIntegration) DeleteState(ctx context.Context, key string) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return errors.New("integration is not connected")
	}

	return i.stateManager.Delete(ctx, key)
}

func (i *DjangoIntegration) SubscribeToEvents(ctx context.Context, eventType string, handler klediosdk.EventHandler) (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return "", errors.New("integration is not connected")
	}

	return i.sdk.SubscribeToEvents(ctx, eventType, handler)
}

func (i *DjangoIntegration) UnsubscribeFromEvents(ctx context.Context, subscriptionID string) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return errors.New("integration is not connected")
	}

	return i.sdk.UnsubscribeFromEvents(ctx, subscriptionID)
}

func (i *DjangoIntegration) PublishEvent(ctx context.Context, event *klediosdk.Event) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.connected {
		return errors.New("integration is not connected")
	}

	return i.sdk.PublishEvent(ctx, event)
}

func (i *DjangoIntegration) GetSDK() *klediosdk.SDK {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.sdk
}

func (i *DjangoIntegration) GetMultiAgentSystem() *langraph.MultiAgentSystem {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.multiAgentSystem
}

func (i *DjangoIntegration) GetLangSmithIntegration() *langsmith.LangGraphIntegration {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.langSmithIntegration
}

func (i *DjangoIntegration) GetBridge() *Bridge {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.bridge
}

func (i *DjangoIntegration) Shutdown() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.cancel()

	if i.connected {
		err := i.bridge.Disconnect(context.Background())
		if err != nil {
			return fmt.Errorf("failed to disconnect from bridge: %w", err)
		}

		err = i.sdk.Shutdown()
		if err != nil {
			return fmt.Errorf("failed to shutdown SDK: %w", err)
		}

		i.connected = false
	}

	return nil
}
