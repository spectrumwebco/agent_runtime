output "namespace" {
  description = "Namespace where Yap2DB is deployed"
  value       = var.namespace
}

output "service_name" {
  description = "Name of the Yap2DB service"
  value       = kubernetes_service.yap2db.metadata[0].name
}

output "service_port" {
  description = "Port of the Yap2DB service"
  value       = 10824
}

output "api_url" {
  description = "URL for Yap2DB API"
  value       = "http://${kubernetes_service.yap2db.metadata[0].name}.${var.namespace}.svc.cluster.local:10824"
}
