output "namespace" {
  description = "The namespace where Flux System is deployed"
  value       = local.namespace
}

output "git_repository_name" {
  description = "The name of the Git repository"
  value       = var.git_repository_name
}

output "kustomization_name" {
  description = "The name of the Kustomization"
  value       = var.kustomization_name
}
