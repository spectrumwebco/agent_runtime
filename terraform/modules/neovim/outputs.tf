output "namespace" {
  description = "Kubernetes namespace for Neovim deployment"
  value       = kubernetes_namespace.neovim.metadata[0].name
}

output "service_name" {
  description = "Kubernetes service name for Neovim"
  value       = kubernetes_service.neovim.metadata[0].name
}

output "service_cluster_ip" {
  description = "Cluster IP of the Neovim service"
  value       = kubernetes_service.neovim.spec[0].cluster_ip
}

output "config_map_name" {
  description = "Name of the Neovim config map"
  value       = kubernetes_config_map.neovim_config.metadata[0].name
}

output "secret_name" {
  description = "Name of the Neovim secrets"
  value       = kubernetes_secret.neovim_secrets.metadata[0].name
}

output "pvc_name" {
  description = "Name of the Neovim persistent volume claim"
  value       = kubernetes_persistent_volume_claim.neovim_data.metadata[0].name
}

output "deployment_name" {
  description = "Name of the Neovim deployment"
  value       = kubernetes_deployment.neovim.metadata[0].name
}

output "kata_enabled" {
  description = "Whether Kata Containers are enabled for Neovim"
  value       = var.enable_kata
}
