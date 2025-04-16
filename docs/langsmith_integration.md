# LangSmith Integration for Agent Runtime

This document describes how to integrate LangSmith with the Agent Runtime system, both using the hosted LangSmith service and self-hosting LangSmith on Kubernetes.

## Overview

LangSmith is a platform for debugging, testing, evaluating, and monitoring LLM applications. It provides tools for:

- Tracing and debugging LLM applications
- Evaluating LLM outputs
- Monitoring LLM applications in production
- Testing LLM applications

The Agent Runtime system integrates with LangSmith to provide these capabilities for our multi-agent system built with LangGraph Go and LangChain Go.

## Client Integration

The Agent Runtime system includes a Go client for LangSmith that can be used to:

- Create and manage runs in LangSmith
- Trace LangGraph executions
- Evaluate agent outputs
- Monitor agent performance

### Using the LangSmith Client

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

### Tracing LangGraph Executions

The Agent Runtime system includes a tracer for LangGraph that sends traces to LangSmith:

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/pkg/langraph"
    "github.com/spectrumwebco/agent_runtime/pkg/langsmith"
)

func main() {
    // Create a LangSmith client
    client := langsmith.NewClient(nil)

    // Create a LangGraph tracer
    tracer := langsmith.NewLangGraphTracer(client, "my-project")

    // Create a LangGraph executor with the tracer
    executor := langraph.NewExecutor(langraph.WithTracer(tracer))

    // Create a graph
    graph := langraph.NewGraph("My Graph")
    // Add nodes and edges to the graph...

    // Execute the graph
    result, err := executor.ExecuteGraph(context.Background(), graph)
    if err != nil {
        panic(err)
    }

    // Use the result...
}
```

## Self-Hosting LangSmith

The Agent Runtime system includes Kubernetes manifests for self-hosting LangSmith. These manifests are located in the `kubernetes/langsmith` directory.

### Prerequisites

- A Kubernetes cluster with at least 16 vCPUs and 64GB of memory
- Helm 3.0 or later
- kubectl 1.18 or later
- A LangSmith Enterprise license key

### Installation

1. Clone the repository:

```bash
git clone https://github.com/spectrumwebco/agent_runtime.git
cd agent_runtime
```

2. Create a namespace for LangSmith:

```bash
kubectl create namespace langsmith
```

3. Create a secret for LangSmith:

```bash
kubectl create secret generic langsmith-secrets \
  --namespace langsmith \
  --from-literal=database-url="postgresql://langsmith:password@langsmith-postgres:5432/langsmith" \
  --from-literal=redis-url="redis://langsmith-redis:6379/0" \
  --from-literal=secret-key="your-secret-key" \
  --from-literal=license-key="your-license-key" \
  --from-literal=postgres-password="password"
```

4. Apply the Kubernetes manifests:

```bash
kubectl apply -f kubernetes/langsmith/deployment.yaml -n langsmith
kubectl apply -f kubernetes/langsmith/service.yaml -n langsmith
kubectl apply -f kubernetes/langsmith/ingress.yaml -n langsmith
```

5. Create a persistent volume claim for Redis:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: langsmith-redis-data
  namespace: langsmith
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
EOF
```

6. Wait for the pods to be ready:

```bash
kubectl wait --for=condition=ready pod -l app=langsmith -n langsmith --timeout=300s
```

7. Access LangSmith at the URL specified in the ingress configuration.

### Configuration

The LangSmith deployment can be configured by editing the `kubernetes/langsmith/values.yaml` file. This file contains configuration for:

- API server
- Frontend
- Database
- Redis
- Blob storage
- Authentication
- Monitoring
- Backup

After editing the values file, you can update the deployment using Helm:

```bash
helm upgrade langsmith ./kubernetes/langsmith -n langsmith -f kubernetes/langsmith/values.yaml
```

## Integration with LangGraph Go

The Agent Runtime system integrates LangSmith with LangGraph Go to provide tracing, debugging, and monitoring capabilities for our multi-agent system.

### Tracing Multi-Agent Communication

The LangSmith integration with LangGraph Go allows tracing of multi-agent communication:

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/pkg/langraph"
    "github.com/spectrumwebco/agent_runtime/pkg/langsmith"
)

func main() {
    // Create a LangSmith client
    client := langsmith.NewClient(nil)

    // Create a LangGraph tracer
    tracer := langsmith.NewLangGraphTracer(client, "multi-agent-system")

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
        Tracer: tracer,
    }

    // Create a multi-agent system
    system, err := langraph.NewMultiAgentSystem(config, langraph.NewExecutor())
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

## Conclusion

The LangSmith integration with the Agent Runtime system provides powerful tools for debugging, testing, evaluating, and monitoring our multi-agent system. By using LangSmith, we can:

- Trace and debug agent communication
- Evaluate agent outputs
- Monitor agent performance in production
- Test agent behavior

Whether using the hosted LangSmith service or self-hosting LangSmith on Kubernetes, the integration provides the same capabilities and benefits.
