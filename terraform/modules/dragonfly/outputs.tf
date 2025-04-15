output "namespace" {
  description = "The namespace where DragonflyDB resources are deployed"
  value       = var.namespace
}

output "redis_service_name" {
  description = "The name of the DragonflyDB Redis service"
  value       = kubernetes_service.dragonfly_redis.metadata[0].name
}

output "redis_service_port" {
  description = "The port of the DragonflyDB Redis service"
  value       = kubernetes_service.dragonfly_redis.spec[0].port[0].port
}

output "http_service_name" {
  description = "The name of the DragonflyDB HTTP service"
  value       = kubernetes_service.dragonfly_http.metadata[0].name
}

output "http_service_port" {
  description = "The port of the DragonflyDB HTTP service"
  value       = kubernetes_service.dragonfly_http.spec[0].port[0].port
}

output "deployment_name" {
  description = "The name of the DragonflyDB deployment"
  value       = kubernetes_deployment.dragonfly.metadata[0].name
}
