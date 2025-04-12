
output "seldon_namespace" {
  description = "Namespace where Seldon Core is deployed"
  value       = kubernetes_namespace.seldon.metadata[0].name
}

output "seldon_model_config_map" {
  description = "Name of ConfigMap for Seldon Core model deployment"
  value       = kubernetes_config_map.seldon_model_config.metadata[0].name
}
