output "vcluster_name" {
  description = "Name of the deployed vCluster"
  value       = var.vcluster_name
}

output "vcluster_namespace" {
  description = "Namespace where vCluster is deployed"
  value       = var.vcluster_namespace
}

output "kubeconfig" {
  description = "Kubeconfig for accessing the vCluster"
  value       = data.kubernetes_secret.vcluster_kubeconfig.data.config
  sensitive   = true
}

output "vcluster_service_name" {
  description = "Name of the vCluster service"
  value       = "vc-${var.vcluster_name}"
}

output "vcluster_service_port" {
  description = "Port of the vCluster service"
  value       = var.service_type == "NodePort" ? var.node_port : 443
}

output "vcluster_endpoint" {
  description = "Endpoint for accessing the vCluster API server"
  value       = "${var.vcluster_name}.${var.vcluster_namespace}.svc.cluster.local"
}

output "vcluster_version" {
  description = "Version of the deployed vCluster"
  value       = var.kubernetes_version
}

output "vcluster_distro" {
  description = "Kubernetes distribution used by vCluster"
  value       = var.distro
}

output "vnode_integration_enabled" {
  description = "Whether vNode integration is enabled"
  value       = var.enable_vnode_integration
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled"
  value       = var.enable_high_availability
}

output "replicas" {
  description = "Number of vCluster replicas"
  value       = var.replicas
}

output "service_account_name" {
  description = "Name of the vCluster service account"
  value       = kubernetes_service_account.vcluster.metadata[0].name
}

output "cluster_role_name" {
  description = "Name of the vCluster cluster role"
  value       = kubernetes_cluster_role.vcluster.metadata[0].name
}

output "cluster_role_binding_name" {
  description = "Name of the vCluster cluster role binding"
  value       = kubernetes_cluster_role_binding.vcluster.metadata[0].name
}
