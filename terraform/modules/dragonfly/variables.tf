variable "namespace" {
  description = "Namespace for DragonflyDB resources"
  type        = string
  default     = "dragonfly-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "container_registry" {
  description = "Container registry for DragonflyDB images"
  type        = string
  default     = "ghcr.io/spectrumwebco"
}

variable "dragonfly_version" {
  description = "Version of DragonflyDB to deploy"
  type        = string
  default     = "latest"
}

variable "replicas" {
  description = "Number of replicas for DragonflyDB deployment"
  type        = number
  default     = 1
}

variable "password" {
  description = "Password for DragonflyDB"
  type        = string
  sensitive   = true
  default     = ""
}

variable "max_memory" {
  description = "Maximum memory for DragonflyDB"
  type        = string
  default     = "512mb"
}

variable "memory_policy" {
  description = "Memory policy for DragonflyDB"
  type        = string
  default     = "allkeys-lru"
}

variable "io_threads" {
  description = "Number of IO threads for DragonflyDB"
  type        = number
  default     = 4
}
