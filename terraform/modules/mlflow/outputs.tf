
output "mlflow_namespace" {
  description = "Namespace where MLFlow is deployed"
  value       = kubernetes_namespace.mlflow.metadata[0].name
}

output "mlflow_tracking_uri" {
  description = "MLFlow tracking URI"
  value       = var.mlflow_tracking_uri
}

output "mlflow_server_service" {
  description = "Name of MLFlow server service"
  value       = kubernetes_service.mlflow_server.metadata[0].name
}

output "mlflow_db_service" {
  description = "Name of MLFlow database service"
  value       = kubernetes_service.mlflow_db.metadata[0].name
}

output "llama4_experiment_config_map" {
  description = "Name of ConfigMap for Llama4 experiment configuration"
  value       = kubernetes_config_map.llama4_experiment_config.metadata[0].name
}
