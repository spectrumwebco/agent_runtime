# Kubernetes and Terraform Compatibility

This document explains how the Django database integration is compatible with existing Kubernetes and Terraform configurations.

## Overview

The Agent Runtime system uses four main database services:

1. **Supabase (PostgreSQL)** - Primary database for Django models
2. **DragonflyDB** - Redis replacement for caching and channels
3. **RAGflow** - Vector database for semantic search
4. **RocketMQ** - Messaging system for state communication

Each of these services is defined in both Kubernetes YAML configurations and Terraform modules to ensure consistency between infrastructure definitions.

## Kubernetes Configuration

The Kubernetes configurations are located in the `kubernetes/` directory and include:

- `supabase-deployment.yaml` - PostgreSQL database deployment
- `dragonfly-deployment.yaml` - DragonflyDB (Redis) deployment
- `ragflow-deployment.yaml` - RAGflow vector database deployment
- `rocketmq-deployment.yaml` - RocketMQ messaging system deployment
- `vault-deployment.yaml` - Hashicorp Vault for secret management
- `database-config.yaml` - ConfigMap with database connection information

These configurations define the services, deployments, and persistent storage needed for each database component.

## Terraform Configuration

The Terraform configurations are located in the `terraform/modules/` directory and include:

- `supabase/` - Terraform module for Supabase deployment
- `dragonfly/` - Terraform module for DragonflyDB deployment
- `ragflow/` - Terraform module for RAGflow deployment
- `rocketmq/` - Terraform module for RocketMQ deployment

Each module defines the same resources as the corresponding Kubernetes YAML files, ensuring that the infrastructure can be managed through either Kubernetes or Terraform.

## Django Integration

The Django database integration is designed to be compatible with both the Kubernetes and Terraform configurations:

1. **Service Discovery** - Django uses Kubernetes service discovery to locate database services:
   ```python
   # Example service discovery in database_config.py
   if is_running_in_kubernetes():
       db_host = "supabase-db.default.svc.cluster.local"
   else:
       db_host = "localhost"
   ```

2. **Secret Management** - Database credentials are retrieved from Hashicorp Vault:
   ```python
   # Example Vault integration in vault.py
   if is_running_in_kubernetes():
       vault_client.authenticate_kubernetes()
   else:
       vault_client.authenticate_token()
   ```

3. **Multiple Database Support** - Django is configured to use multiple databases:
   ```python
   # Example database router in database_routers.py
   class AgentDatabaseRouter:
       def db_for_read(self, model, **hints):
           if model._meta.app_label == 'python_agent':
               return 'agent_db'
           return None
   ```

4. **Fallback Configuration** - Local development uses MariaDB as a fallback:
   ```python
   # Example fallback in settings.py
   try:
       from .vault import database_secrets
       DATABASES = database_secrets.configure_django_databases()
   except Exception:
       from .database_config import DATABASES
   ```

## Ensuring Compatibility

To ensure compatibility between Django, Kubernetes, and Terraform:

1. **Service Names** - Use consistent service names across all configurations
2. **Port Numbers** - Use the same port numbers for each service
3. **Environment Variables** - Use environment variables for configuration
4. **Secret Management** - Use Vault for credential management
5. **Health Checks** - Implement health checks for all services

## Testing Compatibility

To test compatibility between Django and Kubernetes:

1. Run the database verification script:
   ```bash
   python verify_db_connections.py
   ```

2. Run the Django database verification command:
   ```bash
   python manage.py verify_db_connections
   ```

3. Deploy to Kubernetes and verify connections:
   ```bash
   kubectl apply -f kubernetes/
   kubectl exec -it <pod-name> -- python manage.py verify_db_connections
   ```

## Conclusion

The Django database integration is fully compatible with both Kubernetes and Terraform configurations. It uses service discovery to locate database services, Vault for credential management, and supports multiple databases with appropriate routing.
