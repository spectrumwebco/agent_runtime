
output "minio_namespace" {
  description = "Namespace where MinIO is deployed"
  value       = kubernetes_namespace.minio.metadata[0].name
}

output "minio_service" {
  description = "Name of MinIO service"
  value       = "minio.${kubernetes_namespace.minio.metadata[0].name}.svc.cluster.local"
}

output "minio_endpoint" {
  description = "MinIO endpoint URL"
  value       = "http://minio.${kubernetes_namespace.minio.metadata[0].name}.svc.cluster.local:9000"
}

output "minio_credentials_secret" {
  description = "Name of MinIO credentials secret"
  value       = kubernetes_secret.minio_credentials.metadata[0].name
}
