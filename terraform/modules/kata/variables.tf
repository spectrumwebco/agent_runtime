# Variables for Kata Containers Module
variable "namespace" {
  description = "Kubernetes namespace for Agent Runtime"
  type        = string
}

variable "node_selector" {
  description = "Node selector for Kata Containers"
  type        = map(string)
  default     = {
    "kata-containers" = "true"
  }
}

variable "kata_ca_cert" {
  description = "Kata Containers CA certificate"
  type        = string
  default     = ""
  sensitive   = true
}

variable "kata_server_cert" {
  description = "Kata Containers server certificate"
  type        = string
  default     = ""
  sensitive   = true
}

variable "kata_server_key" {
  description = "Kata Containers server key"
  type        = string
  default     = ""
  sensitive   = true
}

variable "enable_desktop" {
  description = "Whether to enable Ubuntu Desktop in Kata Containers"
  type        = bool
  default     = true
}

variable "enable_browser_agent" {
  description = "Whether to enable E2B Surf browser agent"
  type        = bool
  default     = true
}

variable "enable_jetbrains_toolbox" {
  description = "Whether to enable JetBrains Toolbox"
  type        = bool
  default     = true
}

variable "enable_vscode" {
  description = "Whether to enable VSCode"
  type        = bool
  default     = true
}

variable "enable_windsurf_ide" {
  description = "Whether to enable Windsurf IDE"
  type        = bool
  default     = true
}

variable "kata_memory_limit" {
  description = "Memory limit for Kata Containers"
  type        = string
  default     = "8Gi"
}

variable "kata_cpu_limit" {
  description = "CPU limit for Kata Containers"
  type        = string
  default     = "4"
}

variable "kata_storage_size" {
  description = "Storage size for Kata Containers"
  type        = string
  default     = "50Gi"
}
