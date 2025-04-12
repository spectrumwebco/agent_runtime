
output "kserve_namespace" {
  description = "Namespace where KServe is deployed"
  value       = kubernetes_namespace.kserve.metadata[0].name
}

output "llama4_model_config_map" {
  description = "Name of ConfigMap for Llama4 model configuration"
  value       = kubernetes_config_map.llama4_model_config.metadata[0].name
}
