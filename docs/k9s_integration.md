# K9s Kubernetes UI Integration

This document describes the integration of [K9s](https://github.com/derailed/k9s) into the Kled.io Framework.

## Overview

K9s is a powerful terminal-based UI for managing Kubernetes clusters. This integration provides a seamless way to access K9s functionality within the Kled.io Framework, allowing users to interact with their Kubernetes clusters through a rich terminal interface.

## Directory Structure

The integration follows the established pattern for the agent_runtime repository:

```
agent_runtime/
├── cmd/cli/commands/
│   └── k9s.go                # CLI commands for K9s
├── internal/kubernetes/k9s/
│   └── client.go             # Client wrapper for K9s
└── pkg/modules/k9s/
    └── module.go             # Module integration for the framework
```

## Features

The K9s integration provides the following capabilities:

1. **Cluster Management**
   - View and manage all Kubernetes resources
   - Real-time monitoring of cluster state
   - Resource filtering and navigation

2. **Resource Operations**
   - Create, edit, and delete resources
   - Scale deployments
   - View logs and events

3. **Context Management**
   - Switch between clusters and namespaces
   - Manage multiple Kubernetes contexts

4. **CLI Commands**
   - `kled k9s run`: Launch the K9s UI
   - `kled k9s resource [resource]`: Navigate directly to a specific resource
   - `kled k9s install`: Install K9s
   - `kled k9s version`: Display the installed K9s version

## Integration with On-Premises Infrastructure

This integration is specifically designed to work with Spectrum Web Co's on-premises Kubernetes infrastructure running on dedicated servers rather than cloud providers. It avoids cloud provider-specific dependencies and focuses on solutions compatible with on-premises deployments.

## Kata Container Support

The integration includes support for Kata containers in Kubernetes deployments, with specific attention to the required Rust-based components:

1. runtime-rs - A Rust-based runtime component for kata containers
2. mem-agent - A Rust-based memory management agent for kata containers

These components are essential for proper kata container functionality in Kubernetes environments.

## vNode Runtime Integration

The integration supports vNode runtime, which can be implemented using the Helm chart available at https://artifacthub.io/packages/helm/loft/vnode-runtime/0.0.2. This component is critical for the Kubernetes cluster and is properly configured alongside other infrastructure components.

## Usage

### Basic Example

```bash
# Launch K9s UI
kled k9s run

# Navigate directly to pods in the default namespace
kled k9s resource pods -n default

# Install K9s
kled k9s install

# Check K9s version
kled k9s version
```

### Configuration Options

The K9s integration supports various configuration options:

- `--kubeconfig`: Path to the kubeconfig file
- `-n, --namespace`: Namespace to use
- `-c, --context`: Kubernetes context to use
- `--readonly`: Run in read-only mode
- `--headless`: Run in headless mode

## Dependencies

- K9s binary (installed via `kled k9s install`)
- Kubernetes cluster access configured via kubeconfig

## Future Enhancements

1. Add support for custom K9s skins and plugins
2. Implement integration with the framework's authentication system
3. Add support for Gitee-specific Kubernetes deployments
4. Enhance integration with other framework components
