# Variables for DragonflyDB Module
variable "namespace" {
  description = "Kubernetes namespace for Agent Runtime"
  type        = string
}

variable "replicas" {
  description = "Number of DragonflyDB replicas"
  type        = number
  default     = 3
}

variable "dragonfly_password" {
  description = "Password for DragonflyDB authentication"
  type        = string
  sensitive   = true
  # Default value should likely be fetched from Vault or passed securely
  default     = "changeme" 
}
