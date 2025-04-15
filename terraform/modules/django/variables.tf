
variable "namespace" {
  description = "Kubernetes namespace to deploy the Django backend"
  type        = string
  default     = "agent-runtime"
}

variable "replicas" {
  description = "Number of replicas for the Django backend deployment"
  type        = number
  default     = 2
}

variable "image" {
  description = "Docker image for the Django backend"
  type        = string
  default     = "spectrumwebco/agent-runtime-django:latest"
}

variable "ingress_host" {
  description = "Hostname for the Django backend ingress"
  type        = string
  default     = "api.agent-runtime.spectrumwebco.com"
}

variable "storage_class" {
  description = "Storage class for persistent volumes"
  type        = string
  default     = "standard"
}

variable "workspace_storage" {
  description = "Size of the workspace storage"
  type        = string
  default     = "10Gi"
}

variable "enable_ssl" {
  description = "Enable SSL for the ingress"
  type        = bool
  default     = true
}

variable "cpu_limit" {
  description = "CPU limit for the Django backend"
  type        = string
  default     = "1"
}

variable "memory_limit" {
  description = "Memory limit for the Django backend"
  type        = string
  default     = "2Gi"
}

variable "cpu_request" {
  description = "CPU request for the Django backend"
  type        = string
  default     = "500m"
}

variable "memory_request" {
  description = "Memory request for the Django backend"
  type        = string
  default     = "1Gi"
}
