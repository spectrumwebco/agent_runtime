package broker

import (
	"context"

	"github.com/micro/go-micro/broker"
)

type Broker struct {
	broker.Broker
}

type BrokerOption func(*BrokerOptions)

type BrokerOptions struct {
	Addrs    []string
	Secure   bool
	Username string
	Password string
}

func WithAddrs(addrs ...string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Addrs = addrs
	}
}

func WithSecure(secure bool) BrokerOption {
	return func(o *BrokerOptions) {
		o.Secure = secure
	}
}

func WithAuth(username, password string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Username = username
		o.Password = password
	}
}

func NewBroker(opts ...BrokerOption) *Broker {
	options := &BrokerOptions{
		Addrs:  []string{},
		Secure: false,
	}

	for _, o := range opts {
		o(options)
	}

	b := broker.NewBroker(
		broker.Addrs(options.Addrs...),
	)

	if options.Secure {
		b = broker.NewBroker(
			broker.Addrs(options.Addrs...),
			broker.Secure(true),
		)
	}

	if options.Username != "" && options.Password != "" {
		b = broker.NewBroker(
			broker.Addrs(options.Addrs...),
			broker.Secure(options.Secure),
			broker.Auth(options.Username, options.Password),
		)
	}

	return &Broker{
		Broker: b,
	}
}

func (b *Broker) Publish(ctx context.Context, topic string, msg interface{}) error {
	data, err := broker.Marshal(msg)
	if err != nil {
		return err
	}

	return b.Broker.Publish(topic, &broker.Message{
		Header: map[string]string{},
		Body:   data,
	})
}

func (b *Broker) Subscribe(topic string, handler interface{}) (broker.Subscriber, error) {
	return b.Broker.Subscribe(topic, func(event broker.Event) error {
		return broker.Unmarshal(event.Message().Body, handler)
	})
}
