output "postgres_operator_namespace" {
  description = "Namespace where the PostgreSQL Operator is deployed"
  value       = kubernetes_namespace.postgres_operator.metadata[0].name
}

output "postgres_cluster_name" {
  description = "Name of the PostgreSQL cluster"
  value       = kubernetes_manifest.ml_postgres_cluster.manifest.metadata.name
}

output "postgres_service_name" {
  description = "Service name for the PostgreSQL cluster"
  value       = "${kubernetes_manifest.ml_postgres_cluster.manifest.metadata.name}.${kubernetes_namespace.postgres_operator.metadata[0].name}.svc.cluster.local"
}

output "postgres_port" {
  description = "Port for the PostgreSQL cluster"
  value       = 5432
}

output "postgres_users" {
  description = "List of PostgreSQL users created"
  value       = kubernetes_manifest.ml_postgres_cluster.manifest.spec.users[*].name
}

output "postgres_databases" {
  description = "List of PostgreSQL databases created"
  value       = flatten([for user in kubernetes_manifest.ml_postgres_cluster.manifest.spec.users : user.databases])
}

output "vnode_runtime_namespace" {
  description = "Namespace where the vNode runtime is deployed"
  value       = helm_release.vnode_runtime.namespace
}

output "vnode_runtime_version" {
  description = "Version of the vNode runtime"
  value       = var.vnode_runtime_version
}
