
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "jupyterhub_version" {
  description = "Version of JupyterHub to deploy"
  type        = string
  default     = "2.0.0"
}

variable "jupyterhub_namespace" {
  description = "Namespace for JupyterHub deployment"
  type        = string
  default     = "jupyter"
}

variable "jupyterhub_storage_size" {
  description = "Size of persistent volume claim for JupyterHub"
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
