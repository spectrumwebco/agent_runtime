output "namespace" {
  description = "The Kubernetes namespace for terraform-operator"
  value       = kubernetes_namespace.terraform_operator.metadata[0].name
}

output "deployment_name" {
  description = "The name of the terraform-operator deployment"
  value       = kubernetes_deployment.terraform_operator.metadata[0].name
}

output "service_account_name" {
  description = "The name of the terraform-operator service account"
  value       = kubernetes_service_account.terraform_operator_sa.metadata[0].name
}
