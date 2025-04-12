output "namespace" {
  description = "The namespace where vNode resources are deployed"
  value       = var.namespace
}

output "service_name" {
  description = "The name of the vNode service"
  value       = kubernetes_service.vnode_runtime.metadata[0].name
}

output "service_port" {
  description = "The port of the vNode service"
  value       = kubernetes_service.vnode_runtime.spec[0].port[0].port
}

output "daemon_set_name" {
  description = "The name of the vNode daemon set"
  value       = kubernetes_daemon_set.vnode_runtime.metadata[0].name
}
