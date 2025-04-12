output "namespace" {
  description = "The namespace where Supabase resources are deployed"
  value       = var.namespace
}

output "http_service_name" {
  description = "The name of the Supabase HTTP service"
  value       = kubernetes_service.supabase_http.metadata[0].name
}

output "http_service_port" {
  description = "The port of the Supabase HTTP service"
  value       = kubernetes_service.supabase_http.spec[0].port[0].port
}

output "postgres_service_name" {
  description = "The name of the Supabase Postgres service"
  value       = kubernetes_service.supabase_postgres.metadata[0].name
}

output "postgres_service_port" {
  description = "The port of the Supabase Postgres service"
  value       = kubernetes_service.supabase_postgres.spec[0].port[0].port
}

output "deployment_name" {
  description = "The name of the Supabase deployment"
  value       = kubernetes_deployment.supabase.metadata[0].name
}
