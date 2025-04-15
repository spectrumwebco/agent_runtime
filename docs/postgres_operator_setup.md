# CrunchyData PostgreSQL Operator Setup

This document describes how to set up and use the CrunchyData PostgreSQL Operator in the agent_runtime project.

## Overview

The CrunchyData PostgreSQL Operator provides a Kubernetes-native PostgreSQL-as-a-Service solution for the agent_runtime project. It automates and simplifies deploying, managing, and maintaining PostgreSQL clusters on Kubernetes.

## Installation

### Prerequisites

- Kubernetes cluster
- kubectl command-line tool
- Helm (optional)

### Installing the Operator

The PostgreSQL Operator can be installed using the provided Kubernetes manifests:

```bash
kubectl apply -f kubernetes/postgres-operator-deployment.yaml
```

This will create:
- A namespace called `pgo`
- Service accounts, roles, and role bindings
- The PostgreSQL Operator deployment
- Configuration for the operator

### Creating a PostgreSQL Cluster

After the operator is installed, you can create a PostgreSQL cluster:

```bash
kubectl apply -f kubernetes/postgres-cluster-crunchy.yaml
```

This will create a PostgreSQL cluster with:
- 2 replicas
- 10Gi of storage
- Backup configuration
- Users and databases for the agent_runtime application

## Integration with Django

The agent_runtime project integrates with the PostgreSQL cluster using the Django ORM. The integration is configured in:

- `backend/agent_api/database_config_postgres.py`
- `backend/apps/python_agent/integrations/crunchydata.py`

### Local Development

For local development, you can set up PostgreSQL using the provided Django management command:

```bash
python manage.py setup_postgres
```

This will:
- Check if PostgreSQL is running
- Create the necessary users
- Create the required databases

## Terraform Integration

For production deployments, the PostgreSQL Operator and cluster are managed using Terraform. The Terraform configuration is in:

- `terraform/modules/postgres/main.tf`
- `terraform/modules/postgres/variables.tf`
- `terraform/modules/postgres/outputs.tf`

## Database Structure

The PostgreSQL cluster includes the following databases:

1. `agent_runtime` - Main application database
2. `agent_db` - Agent-specific data
3. `trajectory_db` - Trajectory data
4. `ml_db` - Machine learning data

## Users

The PostgreSQL cluster includes the following users:

1. `agent_user` - Superuser with access to all databases
2. `app_user` - Regular user with access to the agent_runtime database

## Integration with Apache Doris

The PostgreSQL databases are configured to work with Apache Doris for enterprise-grade data analytics. The integration allows:

- Data synchronization between PostgreSQL and Doris
- Complex analytics queries using Doris
- High-performance data processing

## Monitoring

The PostgreSQL cluster is monitored using:

- Prometheus for metrics collection
- Grafana for visualization
- Kubernetes events fed into Kafka

## Backup and Recovery

Backups are managed by pgBackRest, which is included in the PostgreSQL cluster configuration. The backup configuration includes:

- Volume-based backups
- Point-in-time recovery
- Scheduled backups

## Troubleshooting

If you encounter issues with the PostgreSQL cluster, you can:

1. Check the operator logs:
   ```bash
   kubectl logs -n pgo deployment/postgres-operator
   ```

2. Check the cluster status:
   ```bash
   kubectl get postgrescluster -n default
   ```

3. Check the database pods:
   ```bash
   kubectl get pods -n default -l postgres-operator.crunchydata.com/cluster=agent-postgres-cluster
   ```

4. Use the CrunchyData client in Django:
   ```python
   from apps.python_agent.integrations.crunchydata import CrunchyDataClient
   client = CrunchyDataClient()
   status = client.get_cluster_status()
   print(status)
   ```
