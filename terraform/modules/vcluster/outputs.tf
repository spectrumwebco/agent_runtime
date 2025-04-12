output "namespace" {
  description = "The namespace where vCluster resources are deployed"
  value       = var.namespace
}

output "service_name" {
  description = "The name of the vCluster service"
  value       = kubernetes_service.vcluster.metadata[0].name
}

output "service_port" {
  description = "The port of the vCluster service"
  value       = kubernetes_service.vcluster.spec[0].port[0].port
}

output "controller_deployment_name" {
  description = "The name of the vCluster controller deployment"
  value       = kubernetes_deployment.vcluster_controller.metadata[0].name
}

output "config_map_name" {
  description = "The name of the vCluster config map"
  value       = kubernetes_config_map.vcluster_config.metadata[0].name
}
