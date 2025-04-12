variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "vault_namespace" {
  description = "Namespace for Vault deployment"
  type        = string
  default     = "vault"
}

variable "vault_version" {
  description = "Version of Vault to deploy"
  type        = string
  default     = "1.13.0"
}

variable "vault_k8s_version" {
  description = "Version of Vault K8s to deploy"
  type        = string
  default     = "1.1.0"
}

variable "vault_token" {
  description = "Root token for Vault"
  type        = string
  default     = "vault-token-secret-key"
  sensitive   = true
}

variable "vault_resources_limits_cpu" {
  description = "CPU limits for Vault"
  type        = string
  default     = "500m"
}

variable "vault_resources_limits_memory" {
  description = "Memory limits for Vault"
  type        = string
  default     = "512Mi"
}

variable "vault_resources_requests_cpu" {
  description = "CPU requests for Vault"
  type        = string
  default     = "250m"
}

variable "vault_resources_requests_memory" {
  description = "Memory requests for Vault"
  type        = string
  default     = "256Mi"
}
