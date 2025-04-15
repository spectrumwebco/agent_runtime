output "doris_fe_service_name" {
  description = "Service name for Apache Doris Frontend"
  value       = kubernetes_service.doris_fe.metadata[0].name
}

output "doris_be_service_name" {
  description = "Service name for Apache Doris Backend"
  value       = kubernetes_service.doris_be.metadata[0].name
}

output "doris_fe_service_namespace" {
  description = "Namespace of the Apache Doris Frontend service"
  value       = kubernetes_service.doris_fe.metadata[0].namespace
}

output "doris_be_service_namespace" {
  description = "Namespace of the Apache Doris Backend service"
  value       = kubernetes_service.doris_be.metadata[0].namespace
}

output "doris_query_endpoint" {
  description = "Query endpoint for Apache Doris"
  value       = "${kubernetes_service.doris_fe.metadata[0].name}.${kubernetes_service.doris_fe.metadata[0].namespace}.svc.cluster.local:9030"
}

output "doris_admin_user" {
  description = "Admin username for Apache Doris"
  value       = "root"
}
