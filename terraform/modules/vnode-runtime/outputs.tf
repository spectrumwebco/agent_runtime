output "namespace" {
  description = "Namespace where vNode runtime is deployed"
  value       = kubernetes_namespace.vnode_runtime.metadata[0].name
}

output "vnode_runtime_endpoint" {
  description = "Endpoint for vNode runtime service"
  value       = "http://vnode-runtime.${kubernetes_namespace.vnode_runtime.metadata[0].name}.svc.cluster.local:8080"
}

output "vnode_runtime_version" {
  description = "Version of vNode runtime deployed"
  value       = var.vnode_runtime_version
}

output "kubeflow_integration_configmap" {
  description = "ConfigMap for KubeFlow integration with vNode runtime"
  value       = "vnode-kubeflow-integration"
}

output "mlflow_integration_configmap" {
  description = "ConfigMap for MLflow integration with vNode runtime"
  value       = "vnode-mlflow-integration"
}

output "kserve_integration_configmap" {
  description = "ConfigMap for KServe integration with vNode runtime"
  value       = "vnode-kserve-integration"
}
