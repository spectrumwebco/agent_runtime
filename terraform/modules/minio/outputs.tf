output "minio_service_name" {
  description = "Name of the MinIO service"
  value       = kubernetes_service.minio.metadata[0].name
}

output "minio_service_namespace" {
  description = "Namespace of the MinIO service"
  value       = kubernetes_service.minio.metadata[0].namespace
}

output "minio_endpoint" {
  description = "MinIO endpoint URL"
  value       = var.minio_endpoint
}

output "minio_buckets" {
  description = "List of created MinIO buckets"
  value       = var.buckets
}

output "minio_console_url" {
  description = "URL for MinIO console"
  value       = var.minio_console_url
}

output "minio_ingress_hosts" {
  description = "Hosts for MinIO ingress"
  value       = var.create_ingress ? [var.minio_domain, var.minio_console_domain] : []
}

output "minio_credentials_secret" {
  description = "Name of the MinIO credentials secret"
  value       = kubernetes_secret.minio_credentials.metadata[0].name
  sensitive   = true
}
