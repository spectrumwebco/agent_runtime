output "namespace" {
  description = "The namespace where ArgoCD is deployed"
  value       = local.namespace
}

output "argocd_server_service_name" {
  description = "The name of the ArgoCD server service"
  value       = "argocd-server"
}

output "argocd_application_names" {
  description = "The names of the ArgoCD applications"
  value       = [for app in kubernetes_manifest.argocd_application : app.manifest.metadata.name]
}
