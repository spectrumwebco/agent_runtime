variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "minio_version" {
  description = "Version of MinIO to deploy"
  type        = string
  default     = "12.1.3"
}

variable "minio_namespace" {
  description = "Namespace for MinIO deployment"
  type        = string
  default     = "minio"
}

variable "minio_storage_size" {
  description = "Size of persistent volume claim for MinIO"
  type        = string
  default     = "100Gi"
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
