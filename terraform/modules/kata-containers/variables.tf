variable "namespace" {
  description = "Namespace for kata containers resources"
  type        = string
  default     = "kata-containers-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "node_selector" {
  description = "Node selector for kata containers runtime"
  type        = map(string)
  default     = {}
}

variable "librechat_code_api_key" {
  description = "API key for LibreChat Code Interpreter"
  type        = string
  sensitive   = true
}

variable "rdp_password" {
  description = "Password for RDP access to sandbox"
  type        = string
  sensitive   = true
  default     = ""
}

variable "ssh_key" {
  description = "SSH key for access to sandbox"
  type        = string
  sensitive   = true
  default     = ""
}
