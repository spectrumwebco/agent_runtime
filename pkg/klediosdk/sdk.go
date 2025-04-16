package klediosdk

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/langsmith"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type SDK struct {
	config *Config

	stateManager sharedstate.Manager

	eventStream eventstream.Stream

	multiAgentSystem *langraph.MultiAgentSystem

	langSmithIntegration *langsmith.LangGraphIntegration

	runtimeAdapters map[RuntimeType]RuntimeAdapter

	ctx    context.Context
	cancel context.CancelFunc

	mu sync.RWMutex

	sessionID string
}

type Config struct {
	StateManagerConfig *sharedstate.Config

	EventStreamConfig *eventstream.Config

	MultiAgentSystemConfig *langraph.MultiAgentSystemConfig

	LangSmithConfig *langsmith.LangSmithConfig

	RuntimeAdaptersConfig map[RuntimeType]*RuntimeAdapterConfig

	DefaultRuntime RuntimeType
}

type RuntimeType string

const (
	RuntimeTypeGo RuntimeType = "go"

	RuntimeTypePython RuntimeType = "python"

	RuntimeTypeReact RuntimeType = "react"

	RuntimeTypeNode RuntimeType = "node"
)

type RuntimeAdapterConfig struct {
	Type RuntimeType

	ConnectionConfig map[string]interface{}

	AuthConfig map[string]interface{}

	Timeout time.Duration
}

type RuntimeAdapter interface {
	Initialize(ctx context.Context, config *RuntimeAdapterConfig) error

	Connect(ctx context.Context) error

	Disconnect(ctx context.Context) error

	IsConnected() bool

	ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error)

	GetState(ctx context.Context, key string) (interface{}, error)

	SetState(ctx context.Context, key string, value interface{}) error

	SubscribeToEvents(ctx context.Context, eventType string, handler EventHandler) (string, error)

	UnsubscribeFromEvents(ctx context.Context, subscriptionID string) error

	PublishEvent(ctx context.Context, event *Event) error
}

type Task struct {
	ID string

	Type string

	Description string

	Input map[string]interface{}

	AgentID string

	Runtime RuntimeType

	Timeout time.Duration

	Metadata map[string]interface{}
}

type TaskResult struct {
	TaskID string

	AgentID string

	Status string

	Output map[string]interface{}

	Error string

	Runtime RuntimeType

	ExecutionTime time.Duration

	Metadata map[string]interface{}
}

type Event struct {
	ID string

	Type string

	Source string

	Timestamp time.Time

	Data map[string]interface{}

	Runtime RuntimeType

	Metadata map[string]interface{}
}

type EventHandler func(ctx context.Context, event *Event) error

func NewSDK(config *Config) (*SDK, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	sdk := &SDK{
		config:          config,
		runtimeAdapters: make(map[RuntimeType]RuntimeAdapter),
		ctx:             ctx,
		cancel:          cancel,
		sessionID:       uuid.New().String(),
	}

	stateManager, err := sharedstate.NewManager(config.StateManagerConfig)
	if err != nil {
		cancel()
		return nil, errors.New("failed to initialize state manager: " + err.Error())
	}
	sdk.stateManager = stateManager

	eventStream, err := eventstream.NewStream(config.EventStreamConfig)
	if err != nil {
		cancel()
		return nil, errors.New("failed to initialize event stream: " + err.Error())
	}
	sdk.eventStream = eventStream

	if config.LangSmithConfig != nil {
		langSmithIntegration, err := langsmith.NewLangGraphIntegration(config.LangSmithConfig)
		if err != nil {
			cancel()
			return nil, errors.New("failed to initialize LangSmith integration: " + err.Error())
		}
		sdk.langSmithIntegration = langSmithIntegration
	}

	if config.MultiAgentSystemConfig != nil {
		var executor *langraph.Executor
		if sdk.langSmithIntegration != nil {
			tracer := sdk.langSmithIntegration.CreateTracer()
			executor = langraph.NewExecutor(langraph.WithTracer(tracer))
		} else {
			executor = langraph.NewExecutor()
		}

		multiAgentSystem, err := langraph.NewMultiAgentSystem(config.MultiAgentSystemConfig, executor)
		if err != nil {
			cancel()
			return nil, errors.New("failed to initialize multi-agent system: " + err.Error())
		}
		sdk.multiAgentSystem = multiAgentSystem

		if sdk.langSmithIntegration != nil {
			_, err = sdk.langSmithIntegration.RegisterMultiAgentSystem(ctx, multiAgentSystem)
			if err != nil {
			}
		}
	}

	for runtimeType, adapterConfig := range config.RuntimeAdaptersConfig {
		adapter, err := NewRuntimeAdapter(runtimeType)
		if err != nil {
			cancel()
			return nil, errors.New("failed to create runtime adapter for " + string(runtimeType) + ": " + err.Error())
		}

		err = adapter.Initialize(ctx, adapterConfig)
		if err != nil {
			cancel()
			return nil, errors.New("failed to initialize runtime adapter for " + string(runtimeType) + ": " + err.Error())
		}

		sdk.runtimeAdapters[runtimeType] = adapter
	}

	return sdk, nil
}

func NewRuntimeAdapter(runtimeType RuntimeType) (RuntimeAdapter, error) {
	switch runtimeType {
	case RuntimeTypeGo:
		return NewGoRuntimeAdapter(), nil
	case RuntimeTypePython:
		return NewPythonRuntimeAdapter(), nil
	case RuntimeTypeReact:
		return NewReactRuntimeAdapter(), nil
	case RuntimeTypeNode:
		return NewNodeRuntimeAdapter(), nil
	default:
		return nil, errors.New("unsupported runtime type: " + string(runtimeType))
	}
}

func (s *SDK) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for runtimeType, adapter := range s.runtimeAdapters {
		err := adapter.Connect(s.ctx)
		if err != nil {
			return errors.New("failed to connect to runtime " + string(runtimeType) + ": " + err.Error())
		}
	}

	return nil
}

func (s *SDK) Shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cancel()

	for runtimeType, adapter := range s.runtimeAdapters {
		err := adapter.Disconnect(s.ctx)
		if err != nil {
		}
	}

	return nil
}

func (s *SDK) ExecuteTask(ctx context.Context, task *Task) (*TaskResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	runtimeType := task.Runtime
	if runtimeType == "" {
		runtimeType = s.config.DefaultRuntime
	}

	adapter, ok := s.runtimeAdapters[runtimeType]
	if !ok {
		return nil, errors.New("no adapter available for runtime type: " + string(runtimeType))
	}

	if !adapter.IsConnected() {
		return nil, errors.New("adapter for runtime type " + string(runtimeType) + " is not connected")
	}

	var execCtx context.Context
	var cancel context.CancelFunc
	if task.Timeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	} else {
		execCtx = ctx
	}

	var taskCtx context.Context
	if s.langSmithIntegration != nil && task.AgentID != "" {
		var taskRunID string
		var err error
		taskCtx, taskRunID, err = s.langSmithIntegration.TraceAgentTask(
			execCtx,
			task.AgentID,
			task.ID,
			task.Type,
			task.Description,
			task.Input,
		)
		if err != nil {
			taskCtx = execCtx
		}

		defer func() {
			if taskRunID != "" {
				_ = s.langSmithIntegration.EndAgentTask(
					taskCtx,
					task.ID,
					nil,
					nil,
				)
			}
		}()
	} else {
		taskCtx = execCtx
	}

	startTime := time.Now()
	result, err := adapter.ExecuteTask(taskCtx, task)
	executionTime := time.Since(startTime)

	if result != nil {
		result.ExecutionTime = executionTime
		result.Runtime = runtimeType
	}

	if s.langSmithIntegration != nil && task.AgentID != "" && result != nil {
		var traceErr error
		if err != nil {
			traceErr = err
		}
		_ = s.langSmithIntegration.EndAgentTask(
			taskCtx,
			task.ID,
			result.Output,
			traceErr,
		)
	}

	return result, err
}

func (s *SDK) GetState(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stateManager.Get(ctx, key)
}

func (s *SDK) SetState(ctx context.Context, key string, value interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stateManager.Set(ctx, key, value)
}

func (s *SDK) DeleteState(ctx context.Context, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stateManager.Delete(ctx, key)
}

func (s *SDK) SubscribeToEvents(ctx context.Context, eventType string, handler EventHandler) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wrapperHandler := func(ctx context.Context, e *eventstream.Event) error {
		event := &Event{
			ID:        e.ID,
			Type:      e.Type,
			Source:    e.Source,
			Timestamp: e.Timestamp,
			Data:      e.Data,
			Metadata:  e.Metadata,
		}

		return handler(ctx, event)
	}

	return s.eventStream.Subscribe(ctx, eventType, wrapperHandler)
}

func (s *SDK) UnsubscribeFromEvents(ctx context.Context, subscriptionID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.eventStream.Unsubscribe(ctx, subscriptionID)
}

func (s *SDK) PublishEvent(ctx context.Context, event *Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	e := &eventstream.Event{
		ID:        event.ID,
		Type:      event.Type,
		Source:    event.Source,
		Timestamp: event.Timestamp,
		Data:      event.Data,
		Metadata:  event.Metadata,
	}

	return s.eventStream.Publish(ctx, e)
}

func (s *SDK) GetMultiAgentSystem() *langraph.MultiAgentSystem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.multiAgentSystem
}

func (s *SDK) GetLangSmithIntegration() *langsmith.LangGraphIntegration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.langSmithIntegration
}

func (s *SDK) GetRuntimeAdapter(runtimeType RuntimeType) (RuntimeAdapter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	adapter, ok := s.runtimeAdapters[runtimeType]
	if !ok {
		return nil, errors.New("no adapter available for runtime type: " + string(runtimeType))
	}

	return adapter, nil
}

func (s *SDK) GetSessionID() string {
	return s.sessionID
}
