variable "namespace" {
  description = "Kubernetes namespace for OTF deployment"
  type        = string
  default     = "otf"
}

variable "api_key" {
  description = "API key for OTF"
  type        = string
  sensitive   = true
  default     = ""
}
