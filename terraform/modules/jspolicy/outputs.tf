output "namespace" {
  description = "Namespace where jsPolicy is deployed"
  value       = var.namespace
}

output "name" {
  description = "Name of the jsPolicy deployment"
  value       = var.name
}

output "chart_version" {
  description = "Version of the jsPolicy Helm chart used"
  value       = var.chart_version
}

output "image_tag" {
  description = "Tag of the jsPolicy image used"
  value       = var.image_tag
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled for jsPolicy"
  value       = var.high_availability
}

output "strict_mode_enabled" {
  description = "Whether strict mode is enabled for policy enforcement"
  value       = var.strict_mode
}

output "prometheus_integration_enabled" {
  description = "Whether Prometheus integration is enabled"
  value       = var.prometheus_enabled
}

output "reporting_enabled" {
  description = "Whether policy reporting is enabled"
  value       = var.enable_reporting
}

output "recovery_controller_enabled" {
  description = "Whether the recovery controller is enabled"
  value       = var.enable_recovery_controller
}

output "policies" {
  description = "List of policies created by this module"
  value = [
    "resource-quota-policy",
    "high-availability-policy",
    "security-policy",
    "network-policy-enforcer",
    "disaster-recovery-policy",
    "rollback-policy",
    "resource-cleanup-policy"
  ]
}

output "policy_bundle_name" {
  description = "Name of the policy bundle"
  value       = "agent-runtime-policies"
}

output "recovery_controller_service_account" {
  description = "Name of the recovery controller service account"
  value       = var.enable_recovery_controller ? kubernetes_service_account.recovery_controller[0].metadata[0].name : null
}

output "recovery_controller_cluster_role" {
  description = "Name of the recovery controller cluster role"
  value       = var.enable_recovery_controller ? kubernetes_cluster_role.recovery_controller[0].metadata[0].name : null
}

output "recovery_controller_cluster_role_binding" {
  description = "Name of the recovery controller cluster role binding"
  value       = var.enable_recovery_controller ? kubernetes_cluster_role_binding.recovery_controller[0].metadata[0].name : null
}
