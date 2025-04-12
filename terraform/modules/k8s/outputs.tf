output "cluster_id" {
  description = "ID of the created Kubernetes cluster"
  value = var.cloud_provider == "aws" ? aws_eks_cluster.this[0].id : (
    var.cloud_provider == "azure" ? azurerm_kubernetes_cluster.this[0].id : (
      var.cloud_provider == "ovh" ? ovh_cloud_project_kube.this[0].id : (
        var.cloud_provider == "fly" ? "${fly_app.k3s_server[0].name}" : null
      )
    )
  )
}

output "kubeconfig_path" {
  description = "Path to the kubeconfig file for the created Kubernetes cluster"
  value = var.cloud_provider == "aws" ? aws_eks_cluster.this[0].kubeconfig : (
    var.cloud_provider == "azure" ? azurerm_kubernetes_cluster.this[0].kube_config_raw : (
      var.cloud_provider == "ovh" ? ovh_cloud_project_kube.this[0].kubeconfig : null
    )
  )
  sensitive = true
}

output "api_endpoint" {
  description = "Endpoint for the Kubernetes API server"
  value = var.cloud_provider == "aws" ? aws_eks_cluster.this[0].endpoint : (
    var.cloud_provider == "azure" ? azurerm_kubernetes_cluster.this[0].kube_config.0.host : (
      var.cloud_provider == "ovh" ? ovh_cloud_project_kube.this[0].endpoint : (
        var.cloud_provider == "fly" ? "https://${fly_app.k3s_server[0].name}.fly.dev:6443" : null
      )
    )
  )
}

output "k8s_name" {
  description = "Name of the Kubernetes cluster"
  value = var.cluster_name
}

output "k8s_region" {
  description = "Region of the Kubernetes cluster"
  value = var.region
}

output "k8s_version" {
  description = "Kubernetes version of the cluster"
  value = var.kubernetes_version
}

output "namespace" {
  description = "Namespace for agent runtime components"
  value = kubernetes_namespace.agent_runtime.metadata[0].name
}

output "vcluster_enabled" {
  description = "Whether vCluster is enabled"
  value = var.vcluster_enabled
}

output "vnode_enabled" {
  description = "Whether vNode runtime is enabled"
  value = var.vnode_enabled
}

output "jspolicy_enabled" {
  description = "Whether jsPolicy is enabled"
  value = var.jspolicy_enabled
}
