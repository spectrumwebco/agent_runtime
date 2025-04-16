package actormodel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ActorSystem struct {
	name      string
	rootActor *Supervisor
	registry  *ActorRegistry
	ctx       context.Context
	cancelFn  context.CancelFunc
	mu        sync.RWMutex
	isStarted bool
}

type ActorSystemOptions struct {
	Name            string
	SupervisorOptions SupervisorOptions
}

func NewActorSystem(opts ActorSystemOptions) *ActorSystem {
	ctx, cancelFn := context.WithCancel(context.Background())
	
	if opts.SupervisorOptions.ID == "" {
		opts.SupervisorOptions.ID = "root"
	}
	
	rootSupervisor := NewSupervisor(opts.SupervisorOptions)
	
	registry := NewActorRegistry()
	
	return &ActorSystem{
		name:      opts.Name,
		rootActor: rootSupervisor,
		registry:  registry,
		ctx:       ctx,
		cancelFn:  cancelFn,
	}
}

func (as *ActorSystem) Start() error {
	as.mu.Lock()
	defer as.mu.Unlock()
	
	if as.isStarted {
		return fmt.Errorf("actor system %s is already started", as.name)
	}
	
	err := as.rootActor.Start()
	if err != nil {
		return fmt.Errorf("failed to start root supervisor: %w", err)
	}
	
	as.isStarted = true
	return nil
}

func (as *ActorSystem) Stop() error {
	as.mu.Lock()
	defer as.mu.Unlock()
	
	if !as.isStarted {
		return fmt.Errorf("actor system %s is not started", as.name)
	}
	
	err := as.rootActor.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop root supervisor: %w", err)
	}
	
	as.cancelFn()
	
	as.isStarted = false
	return nil
}

func (as *ActorSystem) SpawnActor(id string, behavior Behavior, state map[string]interface{}) (*Actor, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	
	if !as.isStarted {
		return nil, fmt.Errorf("actor system %s is not started", as.name)
	}
	
	if as.registry.Exists(id) {
		return nil, fmt.Errorf("actor with id %s already exists", id)
	}
	
	actor, err := as.rootActor.SpawnChild(id, behavior, state)
	if err != nil {
		return nil, err
	}
	
	as.registry.Register(id, actor)
	
	err = actor.Start()
	if err != nil {
		return nil, err
	}
	
	return actor, nil
}

func (as *ActorSystem) SpawnSupervisor(opts SupervisorOptions) (*Supervisor, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	
	if !as.isStarted {
		return nil, fmt.Errorf("actor system %s is not started", as.name)
	}
	
	if as.registry.Exists(opts.ID) {
		return nil, fmt.Errorf("actor with id %s already exists", opts.ID)
	}
	
	opts.Parent = as.rootActor.Actor
	
	supervisor := NewSupervisor(opts)
	
	as.registry.Register(opts.ID, supervisor.Actor)
	
	err := supervisor.Start()
	if err != nil {
		return nil, err
	}
	
	return supervisor, nil
}

func (as *ActorSystem) GetActor(id string) (*Actor, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	
	actor, err := as.registry.Get(id)
	if err != nil {
		return nil, err
	}
	
	return actor, nil
}

type ActorRegistry struct {
	actors map[string]*Actor
	mu     sync.RWMutex
}

func NewActorRegistry() *ActorRegistry {
	return &ActorRegistry{
		actors: make(map[string]*Actor),
	}
}

func (ar *ActorRegistry) Register(id string, actor *Actor) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	
	ar.actors[id] = actor
}

func (ar *ActorRegistry) Unregister(id string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	
	delete(ar.actors, id)
}

func (ar *ActorRegistry) Get(id string) (*Actor, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	
	actor, exists := ar.actors[id]
	if !exists {
		return nil, fmt.Errorf("actor with id %s does not exist", id)
	}
	
	return actor, nil
}

func (ar *ActorRegistry) Exists(id string) bool {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	
	_, exists := ar.actors[id]
	return exists
}

func (ar *ActorRegistry) List() []*Actor {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	
	actors := make([]*Actor, 0, len(ar.actors))
	for _, actor := range ar.actors {
		actors = append(actors, actor)
	}
	
	return actors
}

type ActorSystemManager struct {
	systems map[string]*ActorSystem
	mu      sync.RWMutex
}

func NewActorSystemManager() *ActorSystemManager {
	return &ActorSystemManager{
		systems: make(map[string]*ActorSystem),
	}
}

func (asm *ActorSystemManager) CreateSystem(opts ActorSystemOptions) (*ActorSystem, error) {
	asm.mu.Lock()
	defer asm.mu.Unlock()
	
	if _, exists := asm.systems[opts.Name]; exists {
		return nil, fmt.Errorf("actor system with name %s already exists", opts.Name)
	}
	
	system := NewActorSystem(opts)
	asm.systems[opts.Name] = system
	
	return system, nil
}

func (asm *ActorSystemManager) GetSystem(name string) (*ActorSystem, error) {
	asm.mu.RLock()
	defer asm.mu.RUnlock()
	
	system, exists := asm.systems[name]
	if !exists {
		return nil, fmt.Errorf("actor system with name %s does not exist", name)
	}
	
	return system, nil
}

func (asm *ActorSystemManager) StartSystem(name string) error {
	system, err := asm.GetSystem(name)
	if err != nil {
		return err
	}
	
	return system.Start()
}

func (asm *ActorSystemManager) StopSystem(name string) error {
	system, err := asm.GetSystem(name)
	if err != nil {
		return err
	}
	
	return system.Stop()
}

func (asm *ActorSystemManager) RemoveSystem(name string) error {
	asm.mu.Lock()
	defer asm.mu.Unlock()
	
	system, exists := asm.systems[name]
	if !exists {
		return fmt.Errorf("actor system with name %s does not exist", name)
	}
	
	if system.isStarted {
		err := system.Stop()
		if err != nil {
			return err
		}
	}
	
	delete(asm.systems, name)
	return nil
}

func (asm *ActorSystemManager) ListSystems() []string {
	asm.mu.RLock()
	defer asm.mu.RUnlock()
	
	systems := make([]string, 0, len(asm.systems))
	for name := range asm.systems {
		systems = append(systems, name)
	}
	
	return systems
}
