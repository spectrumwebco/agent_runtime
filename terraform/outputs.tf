
output "ml_infrastructure_namespace" {
  description = "Namespace for ML infrastructure"
  value       = kubernetes_namespace.ml_infrastructure.metadata[0].name
}

output "kubeflow_namespace" {
  description = "Namespace for KubeFlow"
  value       = module.kubeflow.kubeflow_namespace
}

output "mlflow_namespace" {
  description = "Namespace for MLFlow"
  value       = var.mlflow_namespace
}

output "kserve_namespace" {
  description = "Namespace for KServe"
  value       = var.kserve_namespace
}

output "minio_namespace" {
  description = "Namespace for MinIO"
  value       = var.minio_namespace
}

output "feast_namespace" {
  description = "Namespace for Feast"
  value       = var.feast_namespace
}

output "jupyterhub_namespace" {
  description = "Namespace for JupyterHub"
  value       = var.jupyterhub_namespace
}

output "seldon_namespace" {
  description = "Namespace for Seldon Core"
  value       = var.seldon_namespace
}

output "h2o_namespace" {
  description = "Namespace for h2o.ai"
  value       = var.h2o_namespace
}

output "mlflow_tracking_uri" {
  description = "MLFlow tracking URI"
  value       = var.mlflow_tracking_uri
}
