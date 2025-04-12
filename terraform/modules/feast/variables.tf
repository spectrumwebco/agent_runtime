
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "feast_version" {
  description = "Version of Feast to deploy"
  type        = string
  default     = "0.30.0"
}

variable "feast_namespace" {
  description = "Namespace for Feast deployment"
  type        = string
  default     = "feast"
}

variable "storage_class_name" {
  description = "Storage class name for persistent volume claims"
  type        = string
  default     = "standard"
}

variable "minio_access_key" {
  description = "MinIO access key"
  type        = string
  default     = "minioadmin"
  sensitive   = true
}

variable "minio_secret_key" {
  description = "MinIO secret key"
  type        = string
  default     = "minioadmin"
  sensitive   = true
}
