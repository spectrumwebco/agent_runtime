output "namespace" {
  description = "The namespace where MCP resources are deployed"
  value       = local.namespace
}

output "mcp_client_service_name" {
  description = "The name of the MCP client service"
  value       = kubernetes_service.mcp_client.metadata[0].name
}

output "mcp_client_service_port" {
  description = "The port of the MCP client service"
  value       = kubernetes_service.mcp_client.spec[0].port[0].port
}

output "mcp_client_deployment_name" {
  description = "The name of the MCP client deployment"
  value       = kubernetes_deployment.mcp_client.metadata[0].name
}

output "mcp_config_map_name" {
  description = "The name of the MCP config map"
  value       = kubernetes_config_map.mcp_config.metadata[0].name
}

output "mcp_secrets_name" {
  description = "The name of the MCP secrets"
  value       = kubernetes_secret.mcp_secrets.metadata[0].name
}
