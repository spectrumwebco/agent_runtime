output "namespace" {
  description = "Namespace where Kata Containers are deployed"
  value       = local.namespace
}

output "runtime_class_name" {
  description = "Name of the Kata runtime class"
  value       = kubernetes_runtime_class.kata.metadata[0].name
}

output "sandbox_service_name" {
  description = "Name of the Kata sandbox service"
  value       = kubernetes_service.kata_sandbox.metadata[0].name
}

output "sandbox_service_namespace" {
  description = "Namespace of the Kata sandbox service"
  value       = kubernetes_service.kata_sandbox.metadata[0].namespace
}

output "rdp_port" {
  description = "RDP port for connecting to Ubuntu Desktop"
  value       = 3389
}

output "e2b_surf_port" {
  description = "Port for E2B Surf browser agent"
  value       = 3000
}

output "rdp_connection_string" {
  description = "Connection string for RDP"
  value       = "${kubernetes_service.kata_sandbox.metadata[0].name}.${kubernetes_service.kata_sandbox.metadata[0].namespace}.svc.cluster.local:3389"
}

output "e2b_surf_url" {
  description = "URL for E2B Surf browser agent"
  value       = "http://${kubernetes_service.kata_sandbox.metadata[0].name}.${kubernetes_service.kata_sandbox.metadata[0].namespace}.svc.cluster.local:3000"
}

output "ubuntu_version" {
  description = "Ubuntu version used in Kata Containers"
  value       = var.ubuntu_version
}

output "desktop_enabled" {
  description = "Whether Ubuntu Desktop is enabled"
  value       = var.enable_desktop
}

output "browser_agent_enabled" {
  description = "Whether E2B Surf browser agent is enabled"
  value       = var.enable_browser_agent
}

output "jetbrains_toolbox_enabled" {
  description = "Whether JetBrains Toolbox is enabled"
  value       = var.enable_jetbrains_toolbox
}

output "vscode_enabled" {
  description = "Whether VSCode is enabled"
  value       = var.enable_vscode
}

output "windsurf_ide_enabled" {
  description = "Whether Windsurf IDE is enabled"
  value       = var.enable_windsurf_ide
}

output "sandbox_network_policy" {
  description = "Network policy for the sandbox"
  value       = var.sandbox_network_policy
}

output "gpu_support_enabled" {
  description = "Whether GPU support is enabled"
  value       = var.enable_gpu_support
}

output "kata_runtime_version" {
  description = "Version of Kata Containers runtime"
  value       = var.kata_runtime_version
}

output "kata_agent_version" {
  description = "Version of Kata Containers agent"
  value       = var.kata_agent_version
}

output "kata_config_path" {
  description = "Path to the Kata configuration directory"
  value       = var.kata_config_path
}

output "kata_runtime_path" {
  description = "Path to the Kata runtime binary"
  value       = var.kata_runtime_path
}

output "prometheus_integration_enabled" {
  description = "Whether Prometheus integration is enabled"
  value       = var.enable_prometheus_integration
}

output "jaeger_integration_enabled" {
  description = "Whether Jaeger integration is enabled"
  value       = var.enable_jaeger_integration
}

output "opentelemetry_integration_enabled" {
  description = "Whether OpenTelemetry integration is enabled"
  value       = var.enable_opentelemetry_integration
}
