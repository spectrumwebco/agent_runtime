package client

import (
	"context"
	"time"

)

type Client struct {
	client.Client
}

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	Registry registry.Registry
	Timeout  time.Duration
	Retries  int
}

func WithRegistry(registry registry.Registry) ClientOption {
	return func(o *ClientOptions) {
		o.Registry = registry
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

func WithRetries(retries int) ClientOption {
	return func(o *ClientOptions) {
		o.Retries = retries
	}
}

func NewClient(opts ...ClientOption) *Client {
	options := &ClientOptions{
		Timeout: time.Second * 5,
		Retries: 3,
	}

	for _, o := range opts {
		o(options)
	}

	c := client.NewClient(
		client.Retries(options.Retries),
		client.RequestTimeout(options.Timeout),
	)

	if options.Registry != nil {
		c = client.NewClient(
			client.Registry(options.Registry),
			client.Retries(options.Retries),
			client.RequestTimeout(options.Timeout),
		)
	}

	return &Client{
		Client: c,
	}
}

func (c *Client) Call(ctx context.Context, service, endpoint string, req, rsp interface{}, opts ...client.CallOption) error {
	return c.Client.Call(ctx, &client.Request{
		Service: service,
		Method:  endpoint,
		Body:    req,
	}, rsp, opts...)
}

func (c *Client) Stream(ctx context.Context, service, endpoint string, req interface{}, opts ...client.CallOption) (client.Stream, error) {
	return c.Client.Stream(ctx, &client.Request{
		Service: service,
		Method:  endpoint,
		Body:    req,
	}, opts...)
}

func (c *Client) Publish(ctx context.Context, topic string, msg interface{}, opts ...client.PublishOption) error {
	return c.Client.Publish(ctx, &client.Message{
		Topic: topic,
		Body:  msg,
	}, opts...)
}
