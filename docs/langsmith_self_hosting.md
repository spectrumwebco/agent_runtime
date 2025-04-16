# Self-Hosting LangSmith for Agent Runtime

This document provides comprehensive guidance on self-hosting LangSmith for the Agent Runtime system, enabling tracing, monitoring, and evaluation of multi-agent systems without requiring the paid LangSmith platform.

## Overview

LangSmith is a platform for debugging, testing, evaluating, and monitoring LLM applications. By self-hosting LangSmith, we can:

1. Avoid the costs associated with the hosted LangSmith platform
2. Maintain full control over our data and infrastructure
3. Integrate seamlessly with our Kubernetes environment
4. Customize the deployment to our specific needs

This document covers both the client and server components of LangSmith self-hosting.

## Architecture

The self-hosted LangSmith deployment consists of the following components:

1. **LangSmith API Server**: The core backend service that handles trace data, run management, and API requests
2. **LangSmith Frontend**: The web UI for visualizing traces, runs, and evaluations
3. **PostgreSQL Database**: Stores trace data, runs, and other persistent information
4. **Redis**: Used for caching and pub/sub messaging
5. **Object Storage**: For storing large artifacts (using filesystem, S3, GCS, or Azure Blob Storage)

The Agent Runtime system integrates with LangSmith through:

1. **LangSmith Go Client**: A Go client library for interacting with the LangSmith API
2. **LangGraph Integration**: Integration between LangGraph Go and LangSmith for tracing multi-agent systems
3. **Kubernetes Deployment**: Configuration for deploying LangSmith on Kubernetes

## Self-Hosting LangSmith Server

### Prerequisites

- Kubernetes cluster with at least 4 vCPUs and 8GB of memory
- Helm 3.0 or later
- kubectl 1.18 or later
- CrunchyData PostgreSQL Operator installed in the cluster
- Hashicorp Vault for secrets management (optional but recommended)

### Deployment Options

#### Option 1: Using the Deployment Script

The simplest way to deploy LangSmith is using the provided deployment script:

```bash
./scripts/deploy_langsmith.sh --namespace langsmith --license-key YOUR_LICENSE_KEY
```

This script:
1. Creates the necessary Kubernetes namespace
2. Sets up secrets using Hashicorp Vault or Kubernetes secrets
3. Deploys PostgreSQL using the CrunchyData operator
4. Deploys Redis
5. Deploys the LangSmith API server and frontend
6. Configures ingress for external access

#### Option 2: Manual Deployment with Kustomize

For more control over the deployment, you can use Kustomize:

```bash
# Create the namespace
kubectl create namespace langsmith

# Create secrets (using Vault or Kubernetes secrets)
kubectl apply -f kubernetes/langsmith/vault-integration.yaml -n langsmith

# Deploy PostgreSQL
kubectl apply -f kubernetes/langsmith/postgres-operator.yaml -n langsmith

# Deploy LangSmith components
kubectl apply -k kubernetes/langsmith -n langsmith
```

### Monitoring and Observability

The LangSmith deployment includes integration with the following monitoring tools:

1. **Prometheus**: For metrics collection
2. **Grafana**: For metrics visualization
3. **Loki**: For log aggregation
4. **Jaeger**: For distributed tracing
5. **OpenTelemetry**: For standardized observability data collection

To enable monitoring, apply the monitoring configuration:

```bash
kubectl apply -f kubernetes/langsmith/monitoring.yaml -n langsmith
```

## LangSmith Client Integration

### Go Client Library

The Agent Runtime system includes a Go client library for interacting with LangSmith:

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/pkg/langsmith"
)

func main() {
    // Create a LangSmith client
    config := langsmith.DefaultConfig()
    client := langsmith.NewClient(config)

    // Create a run
    run, err := client.CreateRun(context.Background(), &langsmith.CreateRunRequest{
        Name:        "Example Run",
        RunType:     "chain",
        StartTime:   time.Now(),
        Inputs:      map[string]interface{}{"input": "Hello, world!"},
        ProjectName: "my-project",
        Status:      "running",
    })
    if err != nil {
        panic(err)
    }

    // Update the run with outputs
    _, err = client.UpdateRun(context.Background(), run.ID, map[string]interface{}{
        "outputs": map[string]interface{}{"output": "Hello from LangSmith!"},
        "end_time": time.Now(),
        "status": "completed",
    })
    if err != nil {
        panic(err)
    }
}
```

### Configuration

The LangSmith client can be configured using environment variables or a configuration file:

#### Environment Variables

```bash
export LANGCHAIN_TRACING_V2=true
export LANGCHAIN_ENDPOINT="http://langsmith-api.langsmith.svc.cluster.local:8000"
export LANGCHAIN_API_KEY="your-api-key"
export LANGCHAIN_PROJECT="agent-runtime"
```

#### Configuration File

```yaml
# langsmith.yaml
enabled: true
api_key: "your-api-key"
api_url: "http://langsmith-api.langsmith.svc.cluster.local:8000"
project_name: "agent-runtime"
self_hosted: true
self_hosted_config:
  url: "http://langsmith-api.langsmith.svc.cluster.local:8000"
  admin_username: "admin"
  admin_password: "your-admin-password"
  license_key: "your-license-key"
  database_url: "postgresql://langsmith:password@langsmith-postgres:5432/langsmith"
  redis_url: "redis://langsmith-redis:6379/0"
  storage_type: "filesystem"
  storage_config:
    path: "/data/langsmith"
```

## Integration with LangGraph Go

The Agent Runtime system integrates LangSmith with LangGraph Go to provide tracing, debugging, and monitoring capabilities for multi-agent systems.

### LangGraph Tracer

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/pkg/langraph"
    "github.com/spectrumwebco/agent_runtime/pkg/langsmith"
)

func main() {
    // Create a LangSmith integration
    integration, err := langsmith.NewLangGraphIntegration(nil)
    if err != nil {
        panic(err)
    }

    // Create a tracer
    tracer := integration.CreateTracer()

    // Create an executor with the tracer
    executor := langraph.NewExecutor(langraph.WithTracer(tracer))

    // Create a graph
    graph := langraph.NewGraph("My Graph")
    // Add nodes and edges to the graph...

    // Register the graph with LangSmith
    graphID, err := integration.RegisterGraph(context.Background(), graph)
    if err != nil {
        panic(err)
    }

    // Execute the graph
    result, err := executor.ExecuteGraph(context.Background(), graph)
    if err != nil {
        panic(err)
    }

    // Use the result...
}
```

### Multi-Agent System Tracing

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/pkg/langraph"
    "github.com/spectrumwebco/agent_runtime/pkg/langsmith"
)

func main() {
    // Create a LangSmith integration
    integration, err := langsmith.NewLangGraphIntegration(nil)
    if err != nil {
        panic(err)
    }

    // Create a tracer
    tracer := integration.CreateTracer()

    // Create a multi-agent system config
    config := &langraph.MultiAgentSystemConfig{
        Name: "My Multi-Agent System",
        Agents: []*langraph.AgentConfig{
            {
                ID:   "frontend",
                Name: "Frontend Agent",
                Type: "frontend",
            },
            {
                ID:   "codegen",
                Name: "Codegen Agent",
                Type: "codegen",
            },
            {
                ID:   "engineering",
                Name: "Engineering Agent",
                Type: "engineering",
            },
        },
    }

    // Create a multi-agent system
    system, err := langraph.NewMultiAgentSystem(config, langraph.NewExecutor(langraph.WithTracer(tracer)))
    if err != nil {
        panic(err)
    }

    // Register the multi-agent system with LangSmith
    systemID, err := integration.RegisterMultiAgentSystem(context.Background(), system)
    if err != nil {
        panic(err)
    }

    // Execute a task
    result, err := system.ExecuteTask(context.Background(), &langraph.Task{
        ID:          "task-1",
        Description: "Build a web application",
        AssignedTo:  "frontend",
    })
    if err != nil {
        panic(err)
    }

    // Use the result...
}
```

## Integration with LangChain Go

The Agent Runtime system also integrates LangSmith with LangChain Go (langchaingo) to provide tracing and monitoring capabilities for LangChain components.

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/pkg/langsmith"
    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/chains"
)

func main() {
    // Create a LangSmith client
    client := langsmith.NewClient(nil)

    // Create a LangChain tracer
    tracer := langsmith.NewLangChainTracer(client, "my-project")

    // Create a LangChain LLM with tracing
    llm := llms.NewChatOpenAI(llms.WithTracer(tracer))

    // Create a LangChain chain with tracing
    chain := chains.NewLLMChain(llm, chains.WithTracer(tracer))

    // Execute the chain
    result, err := chain.Call(context.Background(), map[string]interface{}{
        "input": "Hello, world!",
    })
    if err != nil {
        panic(err)
    }

    // Use the result...
}
```

## Verification and Testing

The Agent Runtime system includes scripts for verifying and testing the LangSmith deployment and integration:

```bash
# Verify the LangSmith deployment
./scripts/verify_langsmith_deployment.sh --namespace langsmith

# Test the LangSmith integration
./scripts/test_langsmith_integration.sh --api-key YOUR_API_KEY --self-hosted
```

## Troubleshooting

### Common Issues

1. **Connection Issues**:
   - Check that the LangSmith API server is running: `kubectl get pods -n langsmith`
   - Verify the API URL is correct in your configuration
   - Check that the API key is valid

2. **Database Issues**:
   - Check the PostgreSQL deployment: `kubectl get pods -n langsmith -l app=langsmith-postgres`
   - Verify the database connection string in the LangSmith API server configuration

3. **Redis Issues**:
   - Check the Redis deployment: `kubectl get pods -n langsmith -l app=langsmith-redis`
   - Verify the Redis connection string in the LangSmith API server configuration

4. **Ingress Issues**:
   - Check the ingress configuration: `kubectl get ingress -n langsmith`
   - Verify that the ingress controller is properly configured

### Logs

To view logs for the LangSmith components:

```bash
# API server logs
kubectl logs -n langsmith -l app=langsmith,component=api

# Frontend logs
kubectl logs -n langsmith -l app=langsmith,component=frontend

# PostgreSQL logs
kubectl logs -n langsmith -l app=langsmith,component=postgres

# Redis logs
kubectl logs -n langsmith -l app=langsmith,component=redis
```

## Conclusion

By self-hosting LangSmith, the Agent Runtime system can leverage the powerful tracing, debugging, and monitoring capabilities of LangSmith without the costs associated with the hosted platform. The integration with LangGraph Go and LangChain Go provides a comprehensive solution for developing, debugging, and monitoring multi-agent systems.

The Kubernetes deployment configuration, along with the Go client library and integration components, makes it easy to set up and use LangSmith in any environment, from development to production.
