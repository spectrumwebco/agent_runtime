package microservices

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
)

// Client represents a client for a microservice
type Client struct {
	Name    string
	Client  client.Client
	Timeout time.Duration
	Retries int
}

// ClientConfig contains configuration for a client
type ClientConfig struct {
	Name     string
	Timeout  time.Duration
	Retries  int
	Registry registry.Registry
}

// NewClient creates a new client
func NewClient(config ClientConfig) (*Client, error) {
	// Set default values
	if config.Timeout == 0 {
		config.Timeout = time.Second * 5
	}
	if config.Retries == 0 {
		config.Retries = 3
	}

	// Create client options
	options := []client.Option{
		client.RequestTimeout(config.Timeout),
		client.Retries(config.Retries),
	}

	// Add registry if provided
	if config.Registry != nil {
		options = append(options, client.Registry(config.Registry))
	}

	// Create the client
	c := client.NewClient(options...)

	return &Client{
		Name:    config.Name,
		Client:  c,
		Timeout: config.Timeout,
		Retries: config.Retries,
	}, nil
}

// Call calls a method on a service
func (c *Client) Call(ctx context.Context, service, endpoint string, req, rsp interface{}, opts ...client.CallOption) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	// Create the request
	request := c.Client.NewRequest(service, endpoint, req)

	// Call the service
	return c.Client.Call(ctx, request, rsp, opts...)
}

// Publish publishes a message to a topic
func (c *Client) Publish(ctx context.Context, topic string, msg interface{}, opts ...client.PublishOption) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	// Create the message
	message := c.Client.NewMessage(topic, msg)

	// Publish the message
	return c.Client.Publish(ctx, message, opts...)
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, handler interface{}, opts ...client.SubscribeOption) (client.Subscriber, error) {
	return c.Client.Subscribe(context.Background(), topic, handler, opts...)
}
