output "namespace" {
  description = "The namespace where JSPolicy resources are deployed"
  value       = var.namespace
}

output "controller_service_name" {
  description = "The name of the JSPolicy controller service"
  value       = kubernetes_service.jspolicy_controller.metadata[0].name
}

output "controller_service_port" {
  description = "The port of the JSPolicy controller service"
  value       = kubernetes_service.jspolicy_controller.spec[0].port[0].port
}

output "validator_service_name" {
  description = "The name of the JSPolicy validator service"
  value       = kubernetes_service.jspolicy_validator.metadata[0].name
}

output "validator_service_port" {
  description = "The port of the JSPolicy validator service"
  value       = kubernetes_service.jspolicy_validator.spec[0].port[0].port
}
