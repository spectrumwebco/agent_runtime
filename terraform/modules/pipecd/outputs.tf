output "namespace" {
  description = "The Kubernetes namespace where PipeCD is deployed"
  value       = kubernetes_namespace.pipecd.metadata[0].name
}

output "control_plane_service" {
  description = "The Kubernetes service for the PipeCD control plane"
  value       = kubernetes_service.pipecd_control_plane.metadata[0].name
}

output "control_plane_endpoint" {
  description = "The endpoint for the PipeCD control plane"
  value       = "${kubernetes_service.pipecd_control_plane.metadata[0].name}.${kubernetes_namespace.pipecd.metadata[0].name}.svc.cluster.local"
}
