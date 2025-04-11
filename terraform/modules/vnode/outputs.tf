output "service_name" {
  description = "Name of the vNode runtime service"
  value       = kubernetes_service.vnode_runtime.metadata[0].name
}

output "service_namespace" {
  description = "Namespace of the vNode runtime service"
  value       = kubernetes_service.vnode_runtime.metadata[0].namespace
}

output "deployment_name" {
  description = "Name of the vNode runtime deployment"
  value       = kubernetes_deployment.vnode_runtime.metadata[0].name
}

output "runtime_class_name" {
  description = "Name of the vNode runtime class"
  value       = kubernetes_runtime_class.vnode.metadata[0].name
}

output "kata_runtime_class_name" {
  description = "Name of the Kata Containers runtime class"
  value       = var.enable_kata_integration ? kubernetes_runtime_class.kata_containers[0].metadata[0].name : null
}
