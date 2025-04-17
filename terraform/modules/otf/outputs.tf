output "namespace" {
  description = "The Kubernetes namespace for OTF"
  value       = kubernetes_namespace.otf.metadata[0].name
}

output "service_name" {
  description = "The name of the OTF service"
  value       = kubernetes_service.otf.metadata[0].name
}

output "service_endpoint" {
  description = "The endpoint for the OTF service"
  value       = "${kubernetes_service.otf.metadata[0].name}.${kubernetes_namespace.otf.metadata[0].name}.svc.cluster.local"
}
