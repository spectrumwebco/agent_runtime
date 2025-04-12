package ai

import "context"

type Option func(*options)

type options struct {
	Model   string
	Context context.Context
}

func WithModel(model string) Option {
	return func(o *options) {
		o.Model = model
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.Context = ctx
	}
}

func newOptions(opts ...Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
