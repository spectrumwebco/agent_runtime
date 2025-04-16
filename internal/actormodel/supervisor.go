package actormodel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SupervisorStrategy string

const (
	OneForOne SupervisorStrategy = "one_for_one"
	
	OneForAll SupervisorStrategy = "one_for_all"
	
	RestForOne SupervisorStrategy = "rest_for_one"
)

type SupervisorOptions struct {
	ID       string
	Strategy SupervisorStrategy
	MaxRestarts int
	WithinDuration time.Duration
	Parent   *Actor
	State    map[string]interface{}
}

type Supervisor struct {
	*Actor
	strategy      SupervisorStrategy
	maxRestarts   int
	withinDuration time.Duration
	restarts      map[string][]time.Time
	mu            sync.RWMutex
}

func NewSupervisor(opts SupervisorOptions) *Supervisor {
	if opts.Strategy == "" {
		opts.Strategy = OneForOne
	}
	
	if opts.MaxRestarts <= 0 {
		opts.MaxRestarts = 10
	}
	
	if opts.WithinDuration <= 0 {
		opts.WithinDuration = 1 * time.Minute
	}
	
	actor := NewActor(ActorOptions{
		ID:       opts.ID,
		Behavior: supervisorBehavior,
		Parent:   opts.Parent,
		State:    opts.State,
	})
	
	supervisor := &Supervisor{
		Actor:         actor,
		strategy:      opts.Strategy,
		maxRestarts:   opts.MaxRestarts,
		withinDuration: opts.WithinDuration,
		restarts:      make(map[string][]time.Time),
	}
	
	return supervisor
}

func supervisorBehavior(ctx context.Context, msg Message) error {
	switch msg.Type {
	case "child_failed":
		payload, ok := msg.Payload.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid payload for child_failed message")
		}
		
		childID, ok := payload["child_id"].(string)
		if !ok {
			return fmt.Errorf("invalid child_id in child_failed message")
		}
		
		supervisor, ok := msg.Sender.(*Supervisor)
		if !ok {
			return fmt.Errorf("child_failed message not sent by a supervisor")
		}
		
		return supervisor.handleChildFailure(childID)
	default:
		return nil
	}
}

func (s *Supervisor) SpawnChild(id string, behavior Behavior, state map[string]interface{}) (*Actor, error) {
	child, err := s.Spawn(id, behavior, state)
	if err != nil {
		return nil, err
	}
	
	s.mu.Lock()
	s.restarts[id] = []time.Time{}
	s.mu.Unlock()
	
	return child, nil
}

func (s *Supervisor) handleChildFailure(childID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.canRestartChild(childID) {
		return fmt.Errorf("too many restarts for child %s", childID)
	}
	
	now := time.Now()
	s.restarts[childID] = append(s.restarts[childID], now)
	
	switch s.strategy {
	case OneForOne:
		return s.restartChild(childID)
	case OneForAll:
		return s.restartAllChildren()
	case RestForOne:
		return s.restartChildAndRest(childID)
	default:
		return fmt.Errorf("unknown supervisor strategy: %s", s.strategy)
	}
}

func (s *Supervisor) canRestartChild(childID string) bool {
	restarts, exists := s.restarts[childID]
	if !exists {
		return true
	}
	
	now := time.Now()
	recentRestarts := []time.Time{}
	for _, t := range restarts {
		if now.Sub(t) <= s.withinDuration {
			recentRestarts = append(recentRestarts, t)
		}
	}
	
	s.restarts[childID] = recentRestarts
	
	return len(recentRestarts) < s.maxRestarts
}

func (s *Supervisor) restartChild(childID string) error {
	child, exists := s.children[childID]
	if !exists {
		return fmt.Errorf("child %s does not exist", childID)
	}
	
	err := child.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop child %s: %w", childID, err)
	}
	
	err = child.Start()
	if err != nil {
		return fmt.Errorf("failed to restart child %s: %w", childID, err)
	}
	
	return nil
}

func (s *Supervisor) restartAllChildren() error {
	for id, child := range s.children {
		err := child.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop child %s: %w", id, err)
		}
	}
	
	for id, child := range s.children {
		err := child.Start()
		if err != nil {
			return fmt.Errorf("failed to restart child %s: %w", id, err)
		}
	}
	
	return nil
}

func (s *Supervisor) restartChildAndRest(childID string) error {
	childrenIDs := []string{}
	for id := range s.children {
		childrenIDs = append(childrenIDs, id)
	}
	
	childIndex := -1
	for i, id := range childrenIDs {
		if id == childID {
			childIndex = i
			break
		}
	}
	
	if childIndex == -1 {
		return fmt.Errorf("child %s does not exist", childID)
	}
	
	for i := childIndex; i < len(childrenIDs); i++ {
		id := childrenIDs[i]
		err := s.children[id].Stop()
		if err != nil {
			return fmt.Errorf("failed to stop child %s: %w", id, err)
		}
	}
	
	for i := childIndex; i < len(childrenIDs); i++ {
		id := childrenIDs[i]
		err := s.children[id].Start()
		if err != nil {
			return fmt.Errorf("failed to restart child %s: %w", id, err)
		}
	}
	
	return nil
}
