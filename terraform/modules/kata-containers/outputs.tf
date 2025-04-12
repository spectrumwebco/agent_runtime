output "namespace" {
  description = "The namespace where kata containers resources are deployed"
  value       = var.namespace
}

output "runtime_config_name" {
  description = "The name of the kata runtime config map"
  value       = kubernetes_config_map.kata_runtime_config.metadata[0].name
}

output "runtime_class_name" {
  description = "The name of the kata runtime class"
  value       = kubernetes_runtime_class.kata_containers.metadata[0].name
}

output "sandbox_service_name" {
  description = "The name of the agent runtime sandbox service"
  value       = kubernetes_service.agent_runtime_sandbox.metadata[0].name
}

output "sandbox_service_port" {
  description = "The port of the agent runtime sandbox service"
  value       = kubernetes_service.agent_runtime_sandbox.spec[0].port[0].port
}

output "sandbox_deployment_name" {
  description = "The name of the agent runtime sandbox deployment"
  value       = kubernetes_deployment.agent_runtime_sandbox.metadata[0].name
}
