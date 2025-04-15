
output "service_name" {
  description = "Name of the Django backend service"
  value       = kubernetes_service.django_backend.metadata[0].name
}

output "service_namespace" {
  description = "Namespace of the Django backend service"
  value       = kubernetes_service.django_backend.metadata[0].namespace
}

output "service_cluster_ip" {
  description = "Cluster IP of the Django backend service"
  value       = kubernetes_service.django_backend.spec[0].cluster_ip
}

output "ingress_host" {
  description = "Hostname for the Django backend ingress"
  value       = var.ingress_host
}

output "deployment_name" {
  description = "Name of the Django backend deployment"
  value       = kubernetes_deployment.django_backend.metadata[0].name
}

output "deployment_replicas" {
  description = "Number of replicas for the Django backend deployment"
  value       = kubernetes_deployment.django_backend.spec[0].replicas
}

output "config_map_name" {
  description = "Name of the Django backend config map"
  value       = kubernetes_config_map.django_config.metadata[0].name
}

output "workspace_pvc_name" {
  description = "Name of the workspace persistent volume claim"
  value       = kubernetes_persistent_volume_claim.workspaces_pvc.metadata[0].name
}

output "workspace_pvc_size" {
  description = "Size of the workspace persistent volume claim"
  value       = kubernetes_persistent_volume_claim.workspaces_pvc.spec[0].resources[0].requests.storage
}
