variable "namespace" {
  description = "Namespace for agent runtime resources"
  type        = string
  default     = "agent-runtime-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}
