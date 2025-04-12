output "namespace" {
  description = "The namespace where jsPolicy resources are deployed"
  value       = var.namespace
}

output "webhook_service_name" {
  description = "The name of the jsPolicy webhook service"
  value       = kubernetes_service.jspolicy_webhook.metadata[0].name
}

output "webhook_service_port" {
  description = "The port of the jsPolicy webhook service"
  value       = kubernetes_service.jspolicy_webhook.spec[0].port[0].port
}

output "controller_deployment_name" {
  description = "The name of the jsPolicy controller deployment"
  value       = kubernetes_deployment.jspolicy_controller.metadata[0].name
}

output "config_map_name" {
  description = "The name of the jsPolicy config map"
  value       = kubernetes_config_map.jspolicy_config.metadata[0].name
}
