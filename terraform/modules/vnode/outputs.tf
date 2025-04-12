output "namespace" {
  description = "The namespace where vNode resources are deployed"
  value       = var.namespace
}

output "service_name" {
  description = "The name of the vNode runtime service"
  value       = kubernetes_service.vnode_runtime.metadata[0].name
}

output "service_port" {
  description = "The port of the vNode runtime service"
  value       = kubernetes_service.vnode_runtime.spec[0].port[0].port
}

output "config_map_name" {
  description = "The name of the vNode config map"
  value       = kubernetes_config_map.vnode_config.metadata[0].name
}
