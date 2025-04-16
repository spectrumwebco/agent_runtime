# Kubernetes Operations (kOps) Integration

This document describes the integration of [kOps](https://github.com/kubernetes/kops) into the Kled.io Framework.

## Overview

kOps is a production-grade Kubernetes cluster management tool that helps users create, destroy, upgrade, and maintain highly available Kubernetes clusters. This integration enables the Kled.io Framework to leverage kOps for comprehensive Kubernetes cluster management across multiple cloud providers.

## Directory Structure

The integration follows the established pattern for the agent_runtime repository:

```
agent_runtime/
├── cmd/cli/commands/
│   └── kops.go                # CLI commands for kOps
├── internal/kubernetes/kops/
│   ├── client.go              # Client wrapper for kOps
│   └── options.go             # Options for kOps operations
└── pkg/modules/kops/
    └── module.go              # Module integration for the framework
```

## Features

The kOps integration provides the following capabilities:

1. **Cluster Management**
   - Create Kubernetes clusters on multiple cloud providers
   - Update existing cluster configurations
   - Delete clusters when no longer needed
   - Validate cluster health and configuration

2. **Instance Group Management**
   - List and edit instance groups
   - Scale instance groups up or down
   - Modify instance group configurations

3. **Cluster Operations**
   - Perform rolling updates of clusters
   - Export kubeconfig for cluster access
   - Manage cluster secrets
   - Dump cluster state for debugging

4. **Installation Management**
   - Install kOps binary for local use
   - Configure state store for cluster state management

## Integration with Multiple Container Runtimes

This integration is designed to work with Spectrum Web Co's infrastructure that supports multiple container runtimes including LXC, Podman, Docker, and Kata Containers. The kOps integration can deploy Kubernetes clusters that use any of these container runtimes, providing consistent behavior across different deployment scenarios.

## CLI Commands

The kOps integration provides the following CLI commands:

```bash
# Create a new Kubernetes cluster
kled kops create-cluster --name=mycluster.example.com --zones=us-east-1a --node-count=2

# Update a cluster
kled kops update-cluster --name=mycluster.example.com --yes

# Delete a cluster
kled kops delete-cluster --name=mycluster.example.com --yes

# Validate a cluster
kled kops validate-cluster --name=mycluster.example.com

# List all clusters
kled kops get-clusters

# Export kubeconfig
kled kops export-kubecfg --name=mycluster.example.com

# Perform a rolling update
kled kops rolling-update --name=mycluster.example.com --yes

# List instance groups
kled kops get-instance-groups --name=mycluster.example.com

# Edit an instance group
kled kops edit-instance-group --name=mycluster.example.com --instance-group=nodes

# List secrets
kled kops get-secrets

# Dump cluster state
kled kops toolbox-dump --name=mycluster.example.com

# Install kOps binary
kled kops install --version=latest
```

## Integration with Other Framework Components

The kOps integration works seamlessly with other components of the Kled.io Framework:

1. **K9s Integration**: Use K9s for cluster management UI after creating clusters with kOps
2. **Microservices**: Deploy microservices to kOps-managed Kubernetes clusters
3. **Authentication**: Use Casbin for RBAC within kOps-managed clusters
4. **Monitoring**: Integrate with monitoring tools for cluster health tracking

## Use Cases

The kOps integration enables various use cases within the Kled.io Framework:

1. **Production Kubernetes Deployment**: Deploy production-grade Kubernetes clusters
2. **Multi-Cloud Strategy**: Create and manage clusters across different cloud providers
3. **Cluster Lifecycle Management**: Manage the entire lifecycle of Kubernetes clusters
4. **Infrastructure as Code**: Define cluster configurations as code
5. **Disaster Recovery**: Create backup clusters for disaster recovery

## Dependencies

- kubernetes/kops
- AWS CLI (for AWS deployments)
- Google Cloud SDK (for GCP deployments)
- Azure CLI (for Azure deployments)

## Future Enhancements

1. Add support for additional cloud providers
2. Implement cluster templates for quick deployment
3. Create high-level abstractions for common cluster patterns
4. Add support for advanced networking configurations
5. Implement cluster federation capabilities

## Kata Containers Integration

The kOps integration includes specific support for Kata Containers as a container runtime option. When creating clusters with kOps, users can specify Kata Containers as the runtime, which provides enhanced isolation and security for containerized workloads.

The integration includes support for two crucial Rust-based components required for Kata Containers in Kubernetes:
1. runtime-rs - A Rust-based runtime component for Kata Containers
2. mem-agent - A Rust-based memory management agent for Kata Containers

These components are automatically configured when deploying clusters with Kata Containers support.
