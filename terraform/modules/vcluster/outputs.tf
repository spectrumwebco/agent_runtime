output "namespace" {
  description = "The namespace where vCluster resources are deployed"
  value       = var.namespace
}

output "webhook_service_name" {
  description = "The name of the vCluster webhook service"
  value       = kubernetes_service.vcluster_webhook.metadata[0].name
}

output "webhook_service_port" {
  description = "The port of the vCluster webhook service"
  value       = kubernetes_service.vcluster_webhook.spec[0].port[0].port
}

output "controller_deployment_name" {
  description = "The name of the vCluster controller deployment"
  value       = kubernetes_deployment.vcluster_controller.metadata[0].name
}
