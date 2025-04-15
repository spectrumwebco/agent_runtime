
output "h2o_namespace" {
  description = "Namespace where h2o.ai is deployed"
  value       = kubernetes_namespace.h2o.metadata[0].name
}

output "h2o_service" {
  description = "Name of h2o.ai service"
  value       = kubernetes_service.h2o.metadata[0].name
}

output "h2o_automl_config_map" {
  description = "Name of ConfigMap for h2o.ai AutoML configuration"
  value       = kubernetes_config_map.h2o_automl_config.metadata[0].name
}
