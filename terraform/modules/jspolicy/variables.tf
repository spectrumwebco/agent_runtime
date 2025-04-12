variable "namespace" {
  description = "Namespace for JSPolicy resources"
  type        = string
  default     = "jspolicy-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "jspolicy_version" {
  description = "Version of JSPolicy to deploy"
  type        = string
  default     = "0.3.0-beta.5"
}
