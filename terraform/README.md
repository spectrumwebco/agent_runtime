# Agent Runtime Terraform Configuration

This directory contains Terraform configurations for deploying the Agent Runtime infrastructure to Kubernetes.

## Overview

The Agent Runtime infrastructure consists of several components:

- **Kata Containers**: Secure container runtime for sandboxed execution
- **vCluster**: Virtual Kubernetes cluster for runtime isolation
- **jsPolicy**: Kubernetes policy enforcement system
- **vNode**: Virtual node runtime for Kubernetes
- **DragonflyDB**: In-memory data store for caching
- **Supabase**: Database and authentication service
- **MCP (Model Control Plane)**: Communication interface between agent and tools
- **ArgoCD**: GitOps continuous delivery tool
- **Flux**: GitOps toolkit for Kubernetes

## Terraform and Kubernetes Alignment

The Terraform configurations in this directory are designed to mirror the Kubernetes configurations in the `k8s` directory. This ensures that:

1. All services defined in Kubernetes are also defined in Terraform
2. Terraform can be used to create all Kubernetes services
3. Infrastructure can be redeployed according to the initial configuration

This alignment enables:
- Drift detection between the intended infrastructure state (Terraform) and the actual deployed state (Kubernetes)
- Consistent infrastructure management across environments
- Reliable disaster recovery

## Module Structure

Each component has its own module in the `modules` directory:

- `modules/kata-containers`: Kata Containers runtime configuration
- `modules/vcluster`: Virtual Kubernetes cluster configuration
- `modules/jspolicy`: Kubernetes policy enforcement configuration
- `modules/vnode`: Virtual node runtime configuration
- `modules/dragonfly`: DragonflyDB configuration
- `modules/supabase`: Supabase configuration
- `modules/mcp`: Model Control Plane configuration
- `modules/argocd`: ArgoCD configuration
- `modules/flux-system`: Flux configuration
- `modules/k8s-base`: Base Kubernetes configuration

## Usage

To apply the Terraform configuration:

```bash
terraform init
terraform plan
terraform apply
```

To detect drift between Terraform and Kubernetes:

```bash
../scripts/enhanced-drift-detection.sh
```
