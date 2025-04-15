output "vault_namespace" {
  description = "Namespace where Vault is deployed"
  value       = kubernetes_namespace.vault.metadata[0].name
}

output "vault_service_name" {
  description = "Service name for Vault"
  value       = "${kubernetes_service.vault.metadata[0].name}.${kubernetes_namespace.vault.metadata[0].name}.svc.cluster.local"
}

output "vault_port" {
  description = "Port for Vault service"
  value       = kubernetes_service.vault.spec[0].port[0].port
}

output "vault_agent_injector_service_name" {
  description = "Service name for Vault Agent Injector"
  value       = "${kubernetes_service.vault_agent_injector.metadata[0].name}.${kubernetes_namespace.vault.metadata[0].name}.svc.cluster.local"
}

output "vault_agent_injector_port" {
  description = "Port for Vault Agent Injector service"
  value       = kubernetes_service.vault_agent_injector.spec[0].port[0].port
}

output "vault_version" {
  description = "Version of Vault deployed"
  value       = var.vault_version
}
