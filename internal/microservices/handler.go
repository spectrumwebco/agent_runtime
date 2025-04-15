package microservices

import (
	"context"
	"fmt"
	"reflect"

	"go-micro.dev/v4/server"
)

// Handler represents a handler for a microservice
type Handler struct {
	Name    string
	Handler interface{}
	Options []server.HandlerOption
}

// NewHandler creates a new handler
func NewHandler(name string, handler interface{}, opts ...server.HandlerOption) *Handler {
	return &Handler{
		Name:    name,
		Handler: handler,
		Options: opts,
	}
}

// Register registers the handler with a server
func (h *Handler) Register(s server.Server) error {
	// Check if the handler is valid
	if h.Handler == nil {
		return fmt.Errorf("handler is nil")
	}

	// Check if the handler is a struct
	handlerType := reflect.TypeOf(h.Handler)
	if handlerType.Kind() != reflect.Ptr || handlerType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("handler must be a pointer to a struct")
	}

	// Register the handler
	return s.Handle(server.NewHandler(h.Handler, h.Options...))
}

// RegisterHandlers registers multiple handlers with a server
func RegisterHandlers(s server.Server, handlers ...*Handler) error {
	for _, handler := range handlers {
		if err := handler.Register(s); err != nil {
			return fmt.Errorf("failed to register handler %s: %w", handler.Name, err)
		}
	}
	return nil
}
