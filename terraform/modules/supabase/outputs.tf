output "namespace" {
  description = "Namespace where Supabase is deployed"
  value       = local.namespace
}

output "service_name" {
  description = "Name of the Supabase service"
  value       = kubernetes_service.supabase.metadata[0].name
}

output "postgres_primary_service" {
  description = "Name of the PostgreSQL primary service"
  value       = kubernetes_service.postgres_primary.metadata[0].name
}

output "postgres_replica_service" {
  description = "Name of the PostgreSQL replica service"
  value       = kubernetes_service.postgres_replica.metadata[0].name
}

output "postgres_connection_string_primary" {
  description = "Connection string for PostgreSQL primary"
  value       = "postgres://postgres:${local.postgres_password}@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/postgres"
  sensitive   = true
}

output "postgres_connection_string_replica" {
  description = "Connection string for PostgreSQL replica"
  value       = "postgres://postgres:${local.postgres_password}@${kubernetes_service.postgres_replica.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/postgres"
  sensitive   = true
}

output "agent_connection_string" {
  description = "Connection string for agent user"
  value       = "postgres://agent:agent_password@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/agent_state"
  sensitive   = true
}

output "task_state_connection_string" {
  description = "Connection string for task state database"
  value       = "postgres://agent:agent_password@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/task_state"
  sensitive   = true
}

output "tool_state_connection_string" {
  description = "Connection string for tool state database"
  value       = "postgres://agent:agent_password@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/tool_state"
  sensitive   = true
}

output "mcp_state_connection_string" {
  description = "Connection string for MCP state database"
  value       = "postgres://agent:agent_password@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/mcp_state"
  sensitive   = true
}

output "prompts_state_connection_string" {
  description = "Connection string for prompts state database"
  value       = "postgres://agent:agent_password@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/prompts_state"
  sensitive   = true
}

output "modules_state_connection_string" {
  description = "Connection string for modules state database"
  value       = "postgres://agent:agent_password@${kubernetes_service.postgres_primary.metadata[0].name}.${local.namespace}.svc.cluster.local:5432/modules_state"
  sensitive   = true
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled for Supabase"
  value       = var.high_availability
}

output "instance_count" {
  description = "Number of Supabase instances"
  value       = local.instance_count
}

output "postgres_replicas" {
  description = "Number of PostgreSQL read replicas"
  value       = local.postgres_replicas
}

output "rollback_enabled" {
  description = "Whether rollback capability is enabled"
  value       = true
}

output "state_databases" {
  description = "List of state databases"
  value       = ["agent_state", "task_state", "tool_state", "mcp_state", "prompts_state", "modules_state"]
}

output "supabase_url" {
  description = "URL for Supabase REST API"
  value       = "http://${kubernetes_service.supabase.metadata[0].name}.${local.namespace}.svc.cluster.local"
}

output "auth_url" {
  description = "URL for Supabase Auth API"
  value       = "http://${kubernetes_service.supabase.metadata[0].name}.${local.namespace}.svc.cluster.local:9999"
}
