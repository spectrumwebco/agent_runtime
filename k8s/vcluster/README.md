# vCluster Configuration

This directory contains the configuration for vCluster, a virtual Kubernetes cluster technology from Loft.

## Features

- **High Availability**: Configured with 3 replicas for control plane redundancy
- **Disaster Recovery**: Automated backups every 6 hours with retention policy
- **Resource Isolation**: Dedicated resources for virtual clusters
- **Node Syncing**: All nodes from host cluster are synced to virtual cluster

## Usage

To deploy a vCluster:

```bash
vcluster create agent-runtime -n agent-runtime-system -f values.yaml
```

To connect to a vCluster:

```bash
vcluster connect agent-runtime -n agent-runtime-system
```

## Integration with Agent Runtime

The vCluster provides isolated Kubernetes environments for the Agent Runtime, allowing:

1. Separate development and production environments
2. Isolated testing environments
3. Multi-tenant deployments
