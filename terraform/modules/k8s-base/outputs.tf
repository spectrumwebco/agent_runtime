output "namespace" {
  description = "The namespace where agent runtime resources are deployed"
  value       = var.namespace
}

output "service_account_name" {
  description = "The name of the service account created for agent runtime"
  value       = kubernetes_service_account.agent_runtime_controller.metadata[0].name
}

output "cluster_role_name" {
  description = "The name of the cluster role created for agent runtime"
  value       = kubernetes_cluster_role.agent_runtime_controller.metadata[0].name
}

output "config_map_name" {
  description = "The name of the config map created for agent runtime"
  value       = kubernetes_config_map.agent_runtime_config.metadata[0].name
}
