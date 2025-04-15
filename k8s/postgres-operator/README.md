# PostgreSQL Operator Integration

This directory contains the Kubernetes manifests for the CrunchyData PostgreSQL Operator integration, which provides a production-ready PostgreSQL database for the ML infrastructure.

## Directory Structure

- `manifests/`: Contains the core Kubernetes manifests for deploying the PostgreSQL Operator
- `config/`: Contains configuration files for PostgreSQL
- `examples/`: Contains example PostgreSQL cluster configurations and integrations

## Features

- Production-ready PostgreSQL deployment using the CrunchyData PostgreSQL Operator
- High availability with multiple replicas
- Automated backups using pgBackRest
- Integration with Hashicorp Vault for secure credential management
- Monitoring with Prometheus and Grafana
- Integration with MLflow, JupyterLab, and other ML components
- vNode runtime integration for enhanced performance

## Usage

1. Deploy the PostgreSQL Operator:
   ```bash
   kubectl apply -f k8s/postgres-operator/manifests/
   ```

2. Deploy a PostgreSQL cluster:
   ```bash
   kubectl apply -f k8s/postgres-operator/examples/ml-postgres-cluster.yaml
   ```

3. Configure Vault integration:
   ```bash
   kubectl apply -f k8s/vault-integration/vault-postgres-integration.yaml
   ```

4. Integrate with MLflow:
   ```bash
   kubectl apply -f k8s/postgres-operator/examples/mlflow-postgres-integration.yaml
   ```

5. Integrate with JupyterLab:
   ```bash
   kubectl apply -f k8s/postgres-operator/examples/jupyterlab-postgres-integration.yaml
   ```

## Security

All sensitive information such as database credentials is managed by Hashicorp Vault. The PostgreSQL Operator is configured to use Vault for credential management, ensuring that no sensitive information is stored directly in Kubernetes manifests or configuration files.

## Monitoring

The PostgreSQL Operator includes built-in monitoring capabilities using the pgMonitor stack, which provides Prometheus exporters for collecting metrics from PostgreSQL instances. These metrics can be visualized using Grafana dashboards.

## Backup and Recovery

Automated backups are configured using pgBackRest, which provides efficient and reliable backup and recovery capabilities for PostgreSQL. Backups are stored in a dedicated persistent volume and can be scheduled to run at regular intervals.
