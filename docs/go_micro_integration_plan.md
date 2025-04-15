# Go-Micro Integration Plan for Kled.io Framework

This document outlines the plan for deeply integrating Go-Micro into the agent_runtime codebase to enhance the Kled.io Framework's distributed systems capabilities.

## 1. Current Implementation Assessment

The agent_runtime repository currently has limited microservices architecture:

- **Basic Service Communication**: Simple RPC implementation
- **Limited Service Discovery**: No automatic service registration
- **No Load Balancing**: No client-side load balancing
- **Limited Resilience**: No built-in retries or circuit breaking

### Existing Files:
- `/internal/microservices/`: Contains basic microservices implementation
- `/internal/server/grpc.go`: Basic gRPC server implementation
- `/cmd/agent_runtime/grpc_server.go`: gRPC server initialization

## 2. Enhancement Goals

The enhanced Go-Micro integration will provide:

- **Service Discovery**: Automatic service registration and discovery
- **Load Balancing**: Client-side load balancing
- **Message Encoding**: Dynamic message encoding
- **Request/Response**: Enhanced RPC communication
- **Async Messaging**: PubSub for event-driven architecture
- **Resilience**: Retries, circuit breaking, and rate limiting

## 3. Implementation Plan

### 3.1 Core Components

#### Service Registry
Create a service registry for service discovery:

```go
// internal/microservices/registry/registry.go
package registry

import (
	"github.com/go-micro/go-micro/v4/registry"
	"github.com/go-micro/go-micro/v4/registry/mdns"
)

var (
	// DefaultRegistry is the default registry
	DefaultRegistry registry.Registry
)

// Initialize initializes the registry
func Initialize() error {
	DefaultRegistry = mdns.NewRegistry()
	return DefaultRegistry.Init()
}

// GetRegistry returns the registry
func GetRegistry() registry.Registry {
	return DefaultRegistry
}
```

#### Service Client
Create a service client for service communication:

```go
// internal/microservices/client/client.go
package client

import (
	"github.com/go-micro/go-micro/v4/client"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/registry"
)

var (
	// DefaultClient is the default client
	DefaultClient client.Client
)

// Initialize initializes the client
func Initialize() error {
	DefaultClient = client.NewClient(
		client.Registry(registry.GetRegistry()),
		client.Retries(3),
	)
	return nil
}

// GetClient returns the client
func GetClient() client.Client {
	return DefaultClient
}
```

#### Service Server
Create a service server for handling requests:

```go
// internal/microservices/server/server.go
package server

import (
	"github.com/go-micro/go-micro/v4/server"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/registry"
)

var (
	// DefaultServer is the default server
	DefaultServer server.Server
)

// Initialize initializes the server
func Initialize() error {
	DefaultServer = server.NewServer(
		server.Registry(registry.GetRegistry()),
	)
	return DefaultServer.Init()
}

// GetServer returns the server
func GetServer() server.Server {
	return DefaultServer
}
```

#### Event Broker
Create an event broker for event-driven communication:

```go
// internal/microservices/broker/broker.go
package broker

import (
	"github.com/go-micro/go-micro/v4/broker"
	"github.com/go-micro/go-micro/v4/broker/memory"
)

var (
	// DefaultBroker is the default broker
	DefaultBroker broker.Broker
)

// Initialize initializes the broker
func Initialize() error {
	DefaultBroker = memory.NewBroker()
	return DefaultBroker.Init()
}

// GetBroker returns the broker
func GetBroker() broker.Broker {
	return DefaultBroker
}
```

#### Service Framework
Create a service framework for creating microservices:

```go
// internal/microservices/service/service.go
package service

import (
	"github.com/go-micro/go-micro/v4"
	"github.com/go-micro/go-micro/v4/server"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/broker"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/client"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/registry"
)

// Service represents a microservice
type Service struct {
	micro.Service
}

// NewService creates a new service
func NewService(name string, opts ...micro.Option) (*Service, error) {
	// Initialize components
	if err := registry.Initialize(); err != nil {
		return nil, err
	}

	if err := client.Initialize(); err != nil {
		return nil, err
	}

	if err := broker.Initialize(); err != nil {
		return nil, err
	}

	// Create service
	service := micro.NewService(
		micro.Name(name),
		micro.Registry(registry.GetRegistry()),
		micro.Broker(broker.GetBroker()),
		micro.Client(client.GetClient()),
		micro.Server(server.NewServer()),
	)

	// Apply options
	for _, opt := range opts {
		opt(&service.Options())
	}

	return &Service{
		Service: service,
	}, nil
}

// Run starts the service
func (s *Service) Run() error {
	return s.Service.Run()
}
```

### 3.2 MCP Integration

Integrate Go-Micro with MCP Host:

```go
// internal/mcp/service/service.go
package service

import (
	"context"

	"github.com/go-micro/go-micro/v4"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/service"
	mcpproto "github.com/spectrumwebco/agent_runtime/internal/mcp/proto"
)

// MCPService represents the MCP service
type MCPService struct {
	service *service.Service
	handler mcpproto.MCPHandler
}

// NewMCPService creates a new MCP service
func NewMCPService(name string, handler mcpproto.MCPHandler) (*MCPService, error) {
	// Create service
	svc, err := service.NewService(
		name,
		micro.WrapHandler(logWrapper),
	)
	if err != nil {
		return nil, err
	}

	// Register handler
	if err := mcpproto.RegisterMCPHandler(svc.Server(), handler); err != nil {
		return nil, err
	}

	return &MCPService{
		service: svc,
		handler: handler,
	}, nil
}

// Run starts the MCP service
func (s *MCPService) Run() error {
	return s.service.Run()
}

// logWrapper is a handler wrapper that logs requests
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[MCP] Received request: %s", req.Method())
		return fn(ctx, req, rsp)
	}
}
```

### 3.3 Agent Integration

Integrate Go-Micro with Agent:

```go
// internal/agent/service/service.go
package service

import (
	"context"

	"github.com/go-micro/go-micro/v4"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/service"
	agentproto "github.com/spectrumwebco/agent_runtime/internal/agent/proto"
)

// AgentService represents the Agent service
type AgentService struct {
	service *service.Service
	handler agentproto.AgentHandler
}

// NewAgentService creates a new Agent service
func NewAgentService(name string, handler agentproto.AgentHandler) (*AgentService, error) {
	// Create service
	svc, err := service.NewService(
		name,
		micro.WrapHandler(logWrapper),
	)
	if err != nil {
		return nil, err
	}

	// Register handler
	if err := agentproto.RegisterAgentHandler(svc.Server(), handler); err != nil {
		return nil, err
	}

	return &AgentService{
		service: svc,
		handler: handler,
	}, nil
}

// Run starts the Agent service
func (s *AgentService) Run() error {
	return s.service.Run()
}

// logWrapper is a handler wrapper that logs requests
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[Agent] Received request: %s", req.Method())
		return fn(ctx, req, rsp)
	}
}
```

### 3.4 CLI Commands

Create CLI commands for service management:

```go
// cmd/cli/commands/service.go
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/registry"
)

// NewServiceCommand creates a new service command
func NewServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Service management commands",
		Long:  `Commands for managing microservices.`,
	}

	cmd.AddCommand(newServiceListCommand())
	cmd.AddCommand(newServiceStartCommand())
	cmd.AddCommand(newServiceStopCommand())

	return cmd
}

// newServiceListCommand creates a command to list services
func newServiceListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List services",
		Long:  `List all registered services.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := registry.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize registry: %w", err)
			}

			services, err := registry.GetRegistry().ListServices()
			if err != nil {
				return fmt.Errorf("failed to list services: %w", err)
			}

			fmt.Println("Services:")
			for _, service := range services {
				fmt.Printf("  %s\n", service.Name)
				nodes, err := registry.GetRegistry().GetService(service.Name)
				if err != nil {
					fmt.Printf("    Error: %v\n", err)
					continue
				}
				for _, node := range nodes {
					fmt.Printf("    %s (%s)\n", node.Id, node.Address)
				}
			}

			return nil
		},
	}

	return cmd
}

// newServiceStartCommand creates a command to start a service
func newServiceStartCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a service",
		Long:  `Start a microservice.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Implementation depends on the service type
			fmt.Printf("Starting service: %s\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Service name")
	cmd.MarkFlagRequired("name")

	return cmd
}

// newServiceStopCommand creates a command to stop a service
func newServiceStopCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a service",
		Long:  `Stop a microservice.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Implementation depends on the service type
			fmt.Printf("Stopping service: %s\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Service name")
	cmd.MarkFlagRequired("name")

	return cmd
}
```

### 3.5 Server Integration

Update the server initialization to use Go-Micro:

```go
// cmd/agent_runtime/main.go
package main

import (
	"log"

	"github.com/spectrumwebco/agent_runtime/internal/microservices/registry"
	"github.com/spectrumwebco/agent_runtime/internal/mcp/service"
)

func main() {
	// Initialize registry
	if err := registry.Initialize(); err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}

	// Create MCP service
	mcpService, err := service.NewMCPService("mcp", &MCPHandler{})
	if err != nil {
		log.Fatalf("Failed to create MCP service: %v", err)
	}

	// Run service
	if err := mcpService.Run(); err != nil {
		log.Fatalf("Failed to run MCP service: %v", err)
	}
}

// MCPHandler implements the MCP handler
type MCPHandler struct{}

// Implement MCP handler methods...
```

## 4. Integration with Casbin

Integrate Go-Micro with Casbin for service-level authorization:

```go
// internal/microservices/auth/auth.go
package auth

import (
	"context"
	"fmt"

	"github.com/go-micro/go-micro/v4/server"
	"github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

// AuthWrapper creates a server wrapper for authorization
func AuthWrapper(enforcer *auth.DomainEnforcer) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// Extract user from context
			user, ok := ctx.Value("user").(string)
			if !ok {
				user = "anonymous"
			}

			// Extract domain from context
			domain, ok := ctx.Value("domain").(string)
			if !ok {
				domain = "system"
			}

			// Check if the user is authorized
			allowed, err := enforcer.Enforce(user, domain, req.Service(), req.Method())
			if err != nil {
				return fmt.Errorf("authorization error: %w", err)
			}

			if !allowed {
				return fmt.Errorf("forbidden: user %s is not allowed to access %s.%s in domain %s", user, req.Service(), req.Method(), domain)
			}

			return fn(ctx, req, rsp)
		}
	}
}
```

## 5. Integration with GORM

Integrate Go-Micro with GORM for data access:

```go
// internal/microservices/data/data.go
package data

import (
	"context"
	"fmt"

	"github.com/go-micro/go-micro/v4/server"
	"github.com/spectrumwebco/agent_runtime/internal/database"
	"gorm.io/gorm"
)

// DBWrapper creates a server wrapper for database transactions
func DBWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// Start transaction
			tx := database.GetDB().Begin()
			if tx.Error != nil {
				return fmt.Errorf("failed to begin transaction: %w", tx.Error)
			}

			// Add transaction to context
			ctx = context.WithValue(ctx, "tx", tx)

			// Call handler
			err := fn(ctx, req, rsp)

			// Commit or rollback transaction
			if err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("failed to commit transaction: %w", err)
			}

			return nil
		}
	}
}

// GetTx gets the transaction from context
func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return database.GetDB()
	}
	return tx
}
```

## 6. Implementation Timeline

1. **Phase 1: Core Components**
   - Implement service registry
   - Create service client
   - Implement service server
   - Create event broker
   - Implement service framework

2. **Phase 2: MCP Integration**
   - Implement MCP service
   - Create MCP handler
   - Update MCP server initialization

3. **Phase 3: Agent Integration**
   - Implement Agent service
   - Create Agent handler
   - Update Agent server initialization

4. **Phase 4: CLI Integration**
   - Implement service management commands
   - Create service list command
   - Implement service start/stop commands

5. **Phase 5: Cross-Cutting Concerns**
   - Implement authentication integration
   - Create authorization wrapper
   - Implement database integration
   - Create transaction wrapper

## 7. Dependencies

- github.com/go-micro/go-micro/v4
- github.com/go-micro/plugins/v4/registry/etcd
- github.com/go-micro/plugins/v4/broker/nats
- github.com/go-micro/plugins/v4/client/grpc
- github.com/go-micro/plugins/v4/server/grpc

## 8. Conclusion

This integration plan provides a comprehensive approach to enhancing the Kled.io Framework with Go-Micro's distributed systems capabilities. The implementation will provide a robust, flexible, and scalable microservices architecture that can handle the framework's growing needs.

The deep integration approach ensures that the framework leverages the full capabilities of Go-Micro while maintaining a cohesive architecture that users can easily interact with. The microservices architecture will provide proper service isolation, discovery, and communication, enhancing the scalability and reliability of the framework.
