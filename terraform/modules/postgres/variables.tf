variable "namespace" {
  description = "Kubernetes namespace for PostgreSQL cluster"
  type        = string
  default     = "default"
}

variable "cluster_name" {
  description = "Name of the PostgreSQL cluster"
  type        = string
  default     = "agent-postgres-cluster"
}

variable "replicas" {
  description = "Number of PostgreSQL replicas"
  type        = number
  default     = 2
}

variable "storage_size" {
  description = "Size of the PostgreSQL storage"
  type        = string
  default     = "10Gi"
}
