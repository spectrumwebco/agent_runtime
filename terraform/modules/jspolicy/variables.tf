variable "namespace" {
  description = "Namespace for jsPolicy resources"
  type        = string
  default     = "jspolicy-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "jspolicy_version" {
  description = "Version of jsPolicy to deploy"
  type        = string
  default     = "0.3.0-beta.5"
}

variable "webhook_cert" {
  description = "TLS certificate for jsPolicy webhook"
  type        = string
  sensitive   = true
  default     = ""
}

variable "webhook_key" {
  description = "TLS key for jsPolicy webhook"
  type        = string
  sensitive   = true
  default     = ""
}
