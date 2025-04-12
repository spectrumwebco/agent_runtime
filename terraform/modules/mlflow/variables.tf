
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "mlflow_version" {
  description = "Version of MLFlow to deploy"
  type        = string
  default     = "2.3.0"
}

variable "mlflow_namespace" {
  description = "Namespace for MLFlow deployment"
  type        = string
  default     = "mlflow"
}

variable "mlflow_storage_size" {
  description = "Size of persistent volume claim for MLFlow"
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

variable "mlflow_tracking_uri" {
  description = "MLFlow tracking URI"
  type        = string
  default     = "http://mlflow-server.mlflow.svc.cluster.local:5000"
}
