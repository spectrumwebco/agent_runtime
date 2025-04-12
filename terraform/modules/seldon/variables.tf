
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "seldon_version" {
  description = "Version of Seldon Core to deploy"
  type        = string
  default     = "1.15.0"
}

variable "seldon_namespace" {
  description = "Namespace for Seldon Core deployment"
  type        = string
  default     = "seldon"
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
