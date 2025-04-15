package kledframework

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Service represents a simplified microservice
type Service struct {
	Name    string
	Version string
	Router  *gin.Engine
	Server  *http.Server
}

// ServiceConfig contains configuration for a microservice
type ServiceConfig struct {
	Name    string
	Version string
	Address string
	Timeout time.Duration
}

// NewService creates a new microservice
func NewService(config ServiceConfig) *Service {
	// Set default values
	if config.Timeout == 0 {
		config.Timeout = time.Second * 5
	}
	if config.Address == "" {
		config.Address = ":8080"
	}

	// Create the router
	router := gin.Default()

	// Create the server
	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
	}

	return &Service{
		Name:    config.Name,
		Version: config.Version,
		Router:  router,
		Server:  server,
	}
}

// Start starts the microservice
func (s *Service) Start() error {
	fmt.Printf("Starting service %s@%s on %s\n", s.Name, s.Version, s.Server.Addr)
	return s.Server.ListenAndServe()
}

// Stop stops the microservice
func (s *Service) Stop() error {
	return s.Server.Close()
}

// RegisterHandler registers a handler for a route
func (s *Service) RegisterHandler(method, path string, handler gin.HandlerFunc) {
	s.Router.Handle(method, path, handler)
}

// Use adds middleware to the service
func (s *Service) Use(middleware ...gin.HandlerFunc) {
	s.Router.Use(middleware...)
}
