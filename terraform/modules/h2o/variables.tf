
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "h2o_version" {
  description = "Version of h2o.ai to deploy"
  type        = string
  default     = "3.38.0.1"
}

variable "h2o_namespace" {
  description = "Namespace for h2o.ai deployment"
  type        = string
  default     = "h2o"
}

variable "h2o_storage_size" {
  description = "Size of persistent volume claim for h2o.ai"
  type        = string
  default     = "50Gi"
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
