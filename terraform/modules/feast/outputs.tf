
output "feast_namespace" {
  description = "Namespace where Feast is deployed"
  value       = kubernetes_namespace.feast.metadata[0].name
}

output "feast_server_service" {
  description = "Name of Feast server service"
  value       = kubernetes_service.feast_server.metadata[0].name
}

output "feast_redis_service" {
  description = "Name of Feast Redis service"
  value       = kubernetes_service.feast_redis.metadata[0].name
}

output "feast_features_config_map" {
  description = "Name of ConfigMap for Feast feature definitions"
  value       = kubernetes_config_map.feast_features.metadata[0].name
}
