output "namespace" {
  description = "The Kubernetes namespace for Lynx"
  value       = kubernetes_namespace.lynx.metadata[0].name
}

output "service_name" {
  description = "The name of the Lynx service"
  value       = kubernetes_service.lynx.metadata[0].name
}

output "service_endpoint" {
  description = "The endpoint for the Lynx service"
  value       = "${kubernetes_service.lynx.metadata[0].name}.${kubernetes_namespace.lynx.metadata[0].name}.svc.cluster.local"
}
