variable "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "postgres_operator_namespace" {
  description = "Namespace for the PostgreSQL Operator"
  type        = string
  default     = "postgres-operator"
}

variable "postgres_version" {
  description = "PostgreSQL version to deploy"
  type        = number
  default     = 15
}

variable "postgres_replicas" {
  description = "Number of PostgreSQL replicas"
  type        = number
  default     = 3
}

variable "postgres_storage_size" {
  description = "Storage size for PostgreSQL data"
  type        = string
  default     = "10Gi"
}

variable "backup_storage_size" {
  description = "Storage size for PostgreSQL backups"
  type        = string
  default     = "20Gi"
}

variable "storage_class_name" {
  description = "Storage class name for PostgreSQL volumes"
  type        = string
  default     = "standard"
}

variable "vault_integration_enabled" {
  description = "Enable Vault integration for PostgreSQL"
  type        = bool
  default     = true
}

variable "vault_address" {
  description = "Vault server address"
  type        = string
  default     = "http://vault.vault.svc.cluster.local:8200"
}

variable "vnode_runtime_enabled" {
  description = "Enable vNode runtime integration"
  type        = bool
  default     = true
}

variable "vnode_runtime_version" {
  description = "vNode runtime version"
  type        = string
  default     = "0.0.2"
}

variable "postgres_resources_limits_cpu" {
  description = "CPU limits for PostgreSQL"
  type        = string
  default     = "1"
}

variable "postgres_resources_limits_memory" {
  description = "Memory limits for PostgreSQL"
  type        = string
  default     = "1Gi"
}

variable "postgres_resources_requests_cpu" {
  description = "CPU requests for PostgreSQL"
  type        = string
  default     = "500m"
}

variable "postgres_resources_requests_memory" {
  description = "Memory requests for PostgreSQL"
  type        = string
  default     = "512Mi"
}
