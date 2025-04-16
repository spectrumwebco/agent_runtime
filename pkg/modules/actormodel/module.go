package actormodel

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/actormodel"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config        *config.Config
	systemManager *actormodel.ActorSystemManager
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config:        cfg,
		systemManager: actormodel.NewActorSystemManager(),
	}
}

func (m *Module) Initialize() error {
	_, err := m.systemManager.CreateSystem(actormodel.ActorSystemOptions{
		Name: "default",
		SupervisorOptions: actormodel.SupervisorOptions{
			ID:              "root",
			Strategy:        actormodel.OneForOne,
			MaxRestarts:     10,
			WithinDuration:  60,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create default actor system: %w", err)
	}

	return nil
}

func (m *Module) Start(ctx context.Context) error {
	err := m.systemManager.StartSystem("default")
	if err != nil {
		return fmt.Errorf("failed to start default actor system: %w", err)
	}

	return nil
}

func (m *Module) Stop(ctx context.Context) error {
	err := m.systemManager.StopSystem("default")
	if err != nil {
		return fmt.Errorf("failed to stop default actor system: %w", err)
	}

	return nil
}

func (m *Module) CreateActorSystem(name string, opts actormodel.SupervisorOptions) (*actormodel.ActorSystem, error) {
	return m.systemManager.CreateSystem(actormodel.ActorSystemOptions{
		Name:              name,
		SupervisorOptions: opts,
	})
}

func (m *Module) GetActorSystem(name string) (*actormodel.ActorSystem, error) {
	return m.systemManager.GetSystem(name)
}

func (m *Module) ListActorSystems() []string {
	return m.systemManager.ListSystems()
}

func (m *Module) RemoveActorSystem(name string) error {
	return m.systemManager.RemoveSystem(name)
}

func (m *Module) CreateActor(id string, behavior actormodel.Behavior, state map[string]interface{}) (*actormodel.Actor, error) {
	system, err := m.systemManager.GetSystem("default")
	if err != nil {
		return nil, err
	}

	return system.SpawnActor(id, behavior, state)
}

func (m *Module) CreateSupervisor(opts actormodel.SupervisorOptions) (*actormodel.Supervisor, error) {
	system, err := m.systemManager.GetSystem("default")
	if err != nil {
		return nil, err
	}

	return system.SpawnSupervisor(opts)
}

func (m *Module) GetActor(id string) (*actormodel.Actor, error) {
	system, err := m.systemManager.GetSystem("default")
	if err != nil {
		return nil, err
	}

	return system.GetActor(id)
}

func (m *Module) SendMessage(actorID string, msgType string, payload interface{}) error {
	actor, err := m.GetActor(actorID)
	if err != nil {
		return err
	}

	msg := actormodel.Message{
		Type:    msgType,
		Payload: payload,
	}

	return actor.Send(msg)
}

func (m *Module) SendMessageAndWait(actorID string, msgType string, payload interface{}) (interface{}, error) {
	actor, err := m.GetActor(actorID)
	if err != nil {
		return nil, err
	}

	reply, err := actor.SendAndWait(msgType, payload)
	if err != nil {
		return nil, err
	}

	return reply.Payload, nil
}
