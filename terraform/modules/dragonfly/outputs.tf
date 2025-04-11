output "service_name" {
  description = "Name of the DragonflyDB service"
  value       = kubernetes_service.dragonfly.metadata[0].name
}

output "service_namespace" {
  description = "Namespace of the DragonflyDB service"
  value       = kubernetes_service.dragonfly.metadata[0].namespace
}

output "service_port" {
  description = "Port of the DragonflyDB service"
  value       = kubernetes_service.dragonfly.spec[0].port[0].port
}

output "connection_string" {
  description = "Connection string for DragonflyDB"
  value       = "${kubernetes_service.dragonfly.metadata[0].name}.${kubernetes_service.dragonfly.metadata[0].namespace}.svc.cluster.local:${kubernetes_service.dragonfly.spec[0].port[0].port}"
}
