output "namespace" {
  description = "The Kubernetes namespace where Kubefirst is deployed"
  value       = kubernetes_namespace.kubefirst.metadata[0].name
}

output "kubefirst_service" {
  description = "The Kubernetes service for Kubefirst"
  value       = kubernetes_service.kubefirst.metadata[0].name
}

output "kubefirst_endpoint" {
  description = "The endpoint for Kubefirst"
  value       = "${kubernetes_service.kubefirst.metadata[0].name}.${kubernetes_namespace.kubefirst.metadata[0].name}.svc.cluster.local"
}

output "gitea_service" {
  description = "The Kubernetes service for Gitea"
  value       = kubernetes_service.gitea.metadata[0].name
}

output "gitea_endpoint" {
  description = "The endpoint for Gitea"
  value       = "${kubernetes_service.gitea.metadata[0].name}.${kubernetes_namespace.kubefirst.metadata[0].name}.svc.cluster.local"
}

output "vault_service" {
  description = "The Kubernetes service for Vault"
  value       = kubernetes_service.vault.metadata[0].name
}

output "vault_endpoint" {
  description = "The endpoint for Vault"
  value       = "${kubernetes_service.vault.metadata[0].name}.${kubernetes_namespace.kubefirst.metadata[0].name}.svc.cluster.local"
}
