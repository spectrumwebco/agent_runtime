
output "kubeflow_namespace" {
  description = "Namespace where KubeFlow is deployed"
  value       = kubernetes_namespace.kubeflow.metadata[0].name
}

output "kubeflow_data_pvc" {
  description = "Name of persistent volume claim for KubeFlow data"
  value       = kubernetes_persistent_volume_claim.kubeflow_data.metadata[0].name
}

output "llama4_training_config_map" {
  description = "Name of ConfigMap for Llama4 training configuration"
  value       = kubernetes_config_map.llama4_training_config.metadata[0].name
}
