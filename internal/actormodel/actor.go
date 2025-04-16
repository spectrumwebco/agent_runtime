package actormodel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Actor struct {
	id        string
	mailbox   chan Message
	behavior  Behavior
	children  map[string]*Actor
	parent    *Actor
	ctx       context.Context
	cancelFn  context.CancelFunc
	mu        sync.RWMutex
	state     map[string]interface{}
	isStarted bool
}

type Message struct {
	Type    string
	Payload interface{}
	Sender  *Actor
	ReplyTo chan Message
}

type Behavior func(ctx context.Context, msg Message) error

type ActorOptions struct {
	ID       string
	Behavior Behavior
	Parent   *Actor
	State    map[string]interface{}
}

func NewActor(opts ActorOptions) *Actor {
	ctx, cancelFn := context.WithCancel(context.Background())
	if opts.Parent != nil {
		parentCtx := opts.Parent.ctx
		ctx, cancelFn = context.WithCancel(parentCtx)
	}

	if opts.State == nil {
		opts.State = make(map[string]interface{})
	}

	return &Actor{
		id:       opts.ID,
		mailbox:  make(chan Message, 100),
		behavior: opts.Behavior,
		children: make(map[string]*Actor),
		parent:   opts.Parent,
		ctx:      ctx,
		cancelFn: cancelFn,
		state:    opts.State,
	}
}

func (a *Actor) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isStarted {
		return fmt.Errorf("actor %s is already started", a.id)
	}

	go a.processMessages()
	a.isStarted = true
	return nil
}

func (a *Actor) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isStarted {
		return fmt.Errorf("actor %s is not started", a.id)
	}

	for _, child := range a.children {
		_ = child.Stop()
	}

	a.cancelFn()
	a.isStarted = false
	return nil
}

func (a *Actor) Send(msg Message) error {
	select {
	case a.mailbox <- msg:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending message to actor %s", a.id)
	case <-a.ctx.Done():
		return fmt.Errorf("actor %s is stopped", a.id)
	}
}

func (a *Actor) SendAndWait(msgType string, payload interface{}) (Message, error) {
	replyTo := make(chan Message, 1)
	msg := Message{
		Type:    msgType,
		Payload: payload,
		ReplyTo: replyTo,
	}

	err := a.Send(msg)
	if err != nil {
		return Message{}, err
	}

	select {
	case reply := <-replyTo:
		return reply, nil
	case <-time.After(10 * time.Second):
		return Message{}, fmt.Errorf("timeout waiting for reply from actor %s", a.id)
	case <-a.ctx.Done():
		return Message{}, fmt.Errorf("actor %s is stopped", a.id)
	}
}

func (a *Actor) Reply(original Message, msgType string, payload interface{}) error {
	if original.ReplyTo == nil {
		return fmt.Errorf("cannot reply to message without ReplyTo channel")
	}

	reply := Message{
		Type:    msgType,
		Payload: payload,
		Sender:  a,
	}

	select {
	case original.ReplyTo <- reply:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending reply to actor")
	case <-a.ctx.Done():
		return fmt.Errorf("actor %s is stopped", a.id)
	}
}

func (a *Actor) Spawn(id string, behavior Behavior, state map[string]interface{}) (*Actor, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.children[id]; exists {
		return nil, fmt.Errorf("child actor with id %s already exists", id)
	}

	child := NewActor(ActorOptions{
		ID:       id,
		Behavior: behavior,
		Parent:   a,
		State:    state,
	})

	a.children[id] = child
	return child, nil
}

func (a *Actor) GetChild(id string) (*Actor, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	child, exists := a.children[id]
	if !exists {
		return nil, fmt.Errorf("child actor with id %s does not exist", id)
	}

	return child, nil
}

func (a *Actor) GetState() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	stateCopy := make(map[string]interface{})
	for k, v := range a.state {
		stateCopy[k] = v
	}

	return stateCopy
}

func (a *Actor) SetState(key string, value interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.state[key] = value
}

func (a *Actor) GetID() string {
	return a.id
}

func (a *Actor) processMessages() {
	for {
		select {
		case msg := <-a.mailbox:
			err := a.behavior(a.ctx, msg)
			if err != nil {
				fmt.Printf("Error processing message in actor %s: %v\n", a.id, err)
			}
		case <-a.ctx.Done():
			return
		}
	}
}
