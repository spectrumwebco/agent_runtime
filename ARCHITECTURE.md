# Agent Runtime Architecture

## Overview

This document outlines the architecture for the Agent Runtime, a Go-based implementation of an autonomous software engineering agent with MCP (Model Context Protocol) support. The system is designed to run in containerized environments (Kubernetes with Kata Containers) and includes robust infrastructure management, state tracking, and event distribution capabilities.

## File Structure

```
agent_runtime/
├── k8s/                           # Kubernetes deployment and MCP server
│   ├── main.go                    # Entry point for MCP server
│   ├── k8s_handlers.go            # Kubernetes resource management handlers
│   ├── vcluster_handlers.go       # vCluster management handlers
│   ├── policy_handlers.go         # jsPolicy management handlers
│   ├── ragflow_handlers.go        # RAGflow management handlers
│   ├── manifests/                 # Kubernetes manifests
│   │   ├── base/                  # Base configurations
│   │   │   ├── namespace.yaml     # Namespace definition
│   │   │   ├── rbac.yaml          # RBAC configurations
│   │   │   └── crds/              # Custom Resource Definitions
│   │   ├── dragonfly/             # DragonflyDB deployment
│   │   ├── rocketmq/              # RocketMQ deployment
│   │   ├── supabase/              # Supabase PostgreSQL deployment
│   │   ├── vcluster/              # vCluster deployment
│   │   ├── jspolicy/              # jsPolicy deployment
│   │   ├── ragflow/               # RAGflow deployment
│   │   └── kata/                  # Kata container configuration
│   └── Dockerfile                 # Container definition for MCP server
├── terraform/                     # Infrastructure as Code
│   ├── main.tf                    # Main Terraform configuration
│   ├── variables.tf               # Input variables
│   ├── outputs.tf                 # Output values
│   ├── providers.tf               # Provider configurations
│   ├── modules/                   # Reusable Terraform modules
│   │   ├── kubernetes/            # Kubernetes cluster module
│   │   ├── vcluster/              # vCluster module
│   │   ├── dragonfly/             # DragonflyDB module
│   │   ├── rocketmq/              # RocketMQ module
│   │   ├── supabase/              # Supabase PostgreSQL module
│   │   └── kata/                  # Kata container module
│   └── environments/              # Environment-specific configurations
│       ├── dev/                   # Development environment
│       ├── staging/               # Staging environment
│       └── prod/                  # Production environment
├── internal/                      # Internal Go packages
│   ├── agent/                     # Agent implementation
│   │   ├── agent.go               # Main agent implementation
│   │   ├── loop.go                # Agent loop implementation
│   │   ├── tools/                 # Tool implementations
│   │   │   ├── tools.go           # Tool registry
│   │   │   ├── filesystem.go      # Filesystem tools
│   │   │   ├── shell.go           # Shell command tools
│   │   │   ├── browser.go         # Browser tools
│   │   │   └── mcp/               # MCP tool integrations
│   │   ├── models/                # Data models
│   │   │   ├── state.go           # Agent state model
│   │   │   ├── action.go          # Action model
│   │   │   └── observation.go     # Observation model
│   │   └── runtime/               # Runtime implementations
│   │       ├── runtime.go         # Runtime interface
│   │       ├── local.go           # Local runtime
│   │       ├── docker.go          # Docker runtime
│   │       ├── kata.go            # Kata container runtime
│   │       └── e2b.go             # E2B runtime
│   ├── ffi/                       # Foreign Function Interface
│   │   ├── python/                # Python integration
│   │   │   ├── script_runner.go   # Python script runner
│   │   │   └── module_loader.go   # Python module loader
│   │   ├── js/                    # JavaScript integration
│   │   │   ├── script_runner.go   # JavaScript script runner
│   │   │   └── module_loader.go   # JavaScript module loader
│   │   └── cpp/                   # C++ integration
│   │       ├── script_runner.go   # C++ script runner
│   │       └── module_loader.go   # C++ module loader
│   └── eventstream/               # Event distribution system
│       ├── stream.go              # Event stream implementation
│       ├── dragonfly.go           # DragonflyDB integration
│       ├── rocketmq.go            # RocketMQ integration
│       ├── models/                # Event models
│       │   ├── event.go           # Base event model
│       │   ├── agent_event.go     # Agent-specific events
│       │   └── system_event.go    # System-specific events
│       └── handlers/              # Event handlers
│           ├── handler.go         # Handler interface
│           ├── agent_handler.go   # Agent event handler
│           └── system_handler.go  # System event handler
├── cmd/                           # Command-line applications
│   ├── agent/                     # Agent CLI
│   │   └── main.go                # Agent entry point
│   └── server/                    # Server CLI
│       └── main.go                # Server entry point
├── pkg/                           # Public Go packages
│   ├── api/                       # API definitions
│   │   ├── server.go              # API server
│   │   ├── routes.go              # API routes
│   │   └── handlers/              # API handlers
│   ├── config/                    # Configuration
│   │   ├── config.go              # Configuration loader
│   │   └── defaults.go            # Default configurations
│   └── utils/                     # Utility functions
│       ├── logger.go              # Logging utilities
│       └── errors.go              # Error handling utilities
├── go.mod                         # Go module definition
└── go.sum                         # Go module checksums
```

## Infrastructure Components

### 1. Kubernetes Cluster with vCluster

The agent runtime will be deployed on a Kubernetes cluster with vCluster for virtual cluster management. This allows for isolated environments for different agent instances.

#### Key Components:
- **Base Kubernetes Cluster**: Production-grade Kubernetes cluster
- **vCluster**: Virtual Kubernetes clusters for isolation
- **jsPolicy**: Policy management for Kubernetes resources
- **Flux CD**: GitOps for Kubernetes components
- **Argo CD**: GitOps for software components

#### Implementation:
- Kubernetes manifests in `k8s/manifests/`
- Terraform modules in `terraform/modules/kubernetes/` and `terraform/modules/vcluster/`
- MCP server handlers in `k8s/k8s_handlers.go` and `k8s/vcluster_handlers.go`

### 2. jsPolicy Integration

jsPolicy will be used to define robust Kubernetes lifecycle policies to ensure the system works without failure for production and enables recovery in case of failures.

#### Key Components:
- **Policy Definitions**: JavaScript-based policies
- **Policy Enforcement**: Kubernetes admission controllers
- **Policy Management**: API for creating, updating, and deleting policies

#### Implementation:
- Policy definitions in `k8s/manifests/jspolicy/`
- Terraform module in `terraform/modules/jspolicy/`
- MCP server handlers in `k8s/policy_handlers.go`

### 3. RAGflow in Kubernetes

RAGflow will be deployed in the Kubernetes cluster to provide retrieval-augmented generation capabilities for the agent.

#### Key Components:
- **Vector Database**: For storing embeddings
- **Document Processor**: For processing documents
- **Query Engine**: For retrieving relevant information
- **Integration API**: For agent interaction

#### Implementation:
- RAGflow manifests in `k8s/manifests/ragflow/`
- MCP server handlers in `k8s/ragflow_handlers.go`

### 4. Supabase PostgreSQL for Task State

Supabase PostgreSQL will be used for storing task state with high availability and read replicas.

#### Key Components:
- **Primary Instance**: For write operations
- **3 Read Replicas**: For read operations and high availability
- **Backup System**: For data recovery
- **Connection Pooling**: For efficient database connections

#### Implementation:
- Supabase manifests in `k8s/manifests/supabase/`
- Terraform module in `terraform/modules/supabase/`
- Database models in `internal/agent/models/`

### 5. DragonflyDB for Rollback and Persistent Context

Two DragonflyDB instances will be deployed for rollback mechanism and persistent context.

#### Key Components:
- **Context Cache**: For storing application context
- **Rollback Cache**: For storing previous states
- **Invalidation Mechanism**: For updating cache when data changes

#### Implementation:
- DragonflyDB manifests in `k8s/manifests/dragonfly/`
- Terraform module in `terraform/modules/dragonfly/`
- Integration in `internal/eventstream/dragonfly.go`

### 6. RocketMQ for Communication

RocketMQ will be used for communication between components using the Go client.

#### Key Components:
- **Message Producers**: For sending messages
- **Message Consumers**: For receiving messages
- **Topics**: For organizing messages
- **Message Persistence**: For reliable message delivery

#### Implementation:
- RocketMQ manifests in `k8s/manifests/rocketmq/`
- Terraform module in `terraform/modules/rocketmq/`
- Integration in `internal/eventstream/rocketmq.go`

### 7. Kata Container with Ubuntu 22.04 LTS Desktop

Kata containers will be used to provide a sandboxed environment for the agent with Ubuntu 22.04 LTS Desktop.

#### Key Components:
- **Kata Runtime**: For running containers
- **Ubuntu 22.04 LTS Desktop**: Base operating system
- **JetBrains Toolbox**: For JetBrains IDEs
- **VSCode**: Default IDE
- **Windsurf IDE**: Alternative IDE

#### Implementation:
- Kata container manifests in `k8s/manifests/kata/`
- Terraform module in `terraform/modules/kata/`
- Runtime implementation in `internal/agent/runtime/kata.go`

### 8. E2B Surf Agent Integration

E2B Surf Agent will be adapted to work with the Kata container for browser automation.

#### Key Components:
- **E2B Surf Agent**: For browser automation
- **Browser Integration**: For web interaction
- **API Integration**: For agent communication

#### Implementation:
- E2B integration in `internal/agent/tools/browser.go`
- Runtime implementation in `internal/agent/runtime/e2b.go`

## Multi-Cloud Support

The infrastructure will support deployment on multiple cloud providers:

1. **Microsoft Azure**
   - AKS for Kubernetes
   - Azure Database for PostgreSQL
   - Azure Cache for Redis (alternative to DragonflyDB)
   - Azure Service Bus (alternative to RocketMQ)

2. **Amazon Web Services (AWS)**
   - EKS for Kubernetes
   - RDS for PostgreSQL
   - ElastiCache for Redis (alternative to DragonflyDB)
   - Amazon MQ for RocketMQ

3. **OVHcloud US**
   - Managed Kubernetes Service
   - Managed PostgreSQL
   - Self-hosted DragonflyDB
   - Self-hosted RocketMQ

4. **Fly.io**
   - Fly Machines for Kubernetes
   - Fly Postgres for PostgreSQL
   - Self-hosted DragonflyDB
   - Self-hosted RocketMQ

## CI/CD Integration

The CI/CD system will integrate with both GitHub and Gitea Enterprise Server:

1. **GitHub Integration**
   - GitHub Actions for CI/CD
   - CodeRabbit AI for PR reviews
   - 100% code coverage with comprehensive tests
   - Public test repository structure

2. **Gitea Enterprise Server Integration**
   - Gitea Actions for CI/CD
   - Tea CLI for automation
   - Private repository management
   - Synchronization with GitHub repositories

## Event Stream Architecture

The Event Stream is a central component that provides the application context cached in memory. It uses DragonflyDB for caching and invalidates the cache only when something changes in the application.

### Key Components:
- **Event Distribution**: Distributes events between components
- **Context Caching**: Caches application context in DragonflyDB
- **Invalidation Logic**: Invalidates cache when data changes
- **Event Persistence**: Persists events for recovery

### Implementation:
- Event stream implementation in `internal/eventstream/stream.go`
- DragonflyDB integration in `internal/eventstream/dragonfly.go`
- RocketMQ integration in `internal/eventstream/rocketmq.go`
- Event models in `internal/eventstream/models/`
