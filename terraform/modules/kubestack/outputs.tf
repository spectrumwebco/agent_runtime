output "namespace" {
  description = "The Kubernetes namespace for Kubestack"
  value       = kubernetes_namespace.kubestack.metadata[0].name
}

output "service_name" {
  description = "The name of the Kubestack service"
  value       = kubernetes_service.kubestack.metadata[0].name
}

output "service_endpoint" {
  description = "The endpoint for the Kubestack service"
  value       = "${kubernetes_service.kubestack.metadata[0].name}.${kubernetes_namespace.kubestack.metadata[0].name}.svc.cluster.local"
}
