output "postgres_cluster_name" {
  description = "Name of the PostgreSQL cluster"
  value       = var.cluster_name
}

output "postgres_namespace" {
  description = "Namespace of the PostgreSQL cluster"
  value       = var.namespace
}

output "postgres_service_name" {
  description = "Service name for the PostgreSQL cluster"
  value       = "${var.cluster_name}-primary"
}

output "postgres_connection_string" {
  description = "Connection string for the PostgreSQL cluster"
  value       = "postgresql://agent_user:${random_password.postgres_password.result}@${var.cluster_name}-primary.${var.namespace}.svc.cluster.local:5432/agent_runtime"
  sensitive   = true
}

resource "random_password" "postgres_password" {
  length  = 16
  special = false
}
