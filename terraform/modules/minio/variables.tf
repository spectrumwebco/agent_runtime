variable "namespace" {
  description = "Kubernetes namespace for MinIO deployment"
  type        = string
  default     = "ml-infrastructure"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "minio_image" {
  description = "MinIO container image"
  type        = string
  default     = "minio/minio:RELEASE.2023-06-29T05-12-28Z"
}

variable "minio_access_key" {
  description = "MinIO access key"
  type        = string
  sensitive   = true
}

variable "minio_secret_key" {
  description = "MinIO secret key"
  type        = string
  sensitive   = true
}

variable "storage_size" {
  description = "Size of the persistent volume claim"
  type        = string
  default     = "100Gi"
}

variable "storage_class" {
  description = "Storage class for the persistent volume claim"
  type        = string
  default     = "standard"
}

variable "memory_request" {
  description = "Memory request for MinIO container"
  type        = string
  default     = "512Mi"
}

variable "cpu_request" {
  description = "CPU request for MinIO container"
  type        = string
  default     = "250m"
}

variable "memory_limit" {
  description = "Memory limit for MinIO container"
  type        = string
  default     = "2Gi"
}

variable "cpu_limit" {
  description = "CPU limit for MinIO container"
  type        = string
  default     = "1"
}

variable "create_ingress" {
  description = "Whether to create an ingress for MinIO"
  type        = bool
  default     = true
}

variable "minio_domain" {
  description = "Domain for MinIO API"
  type        = string
  default     = "minio.example.com"
}

variable "minio_console_domain" {
  description = "Domain for MinIO console"
  type        = string
  default     = "minio-console.example.com"
}

variable "minio_console_url" {
  description = "URL for MinIO console"
  type        = string
  default     = "https://minio-console.example.com"
}

variable "tls_secret_name" {
  description = "Name of the TLS secret for MinIO ingress"
  type        = string
  default     = "minio-tls"
}

variable "prometheus_url" {
  description = "URL for Prometheus server"
  type        = string
  default     = "http://prometheus-server.monitoring.svc.cluster.local:9090"
}

variable "region" {
  description = "MinIO region"
  type        = string
  default     = "us-east-1"
}

variable "create_buckets" {
  description = "Whether to create buckets"
  type        = bool
  default     = true
}

variable "buckets" {
  description = "List of buckets to create"
  type        = list(string)
  default     = [
    "mlflow-artifacts",
    "model-registry",
    "training-data",
    "model-serving",
    "datasets",
    "checkpoints",
    "logs",
    "metrics",
    "backups"
  ]
}

variable "configure_lifecycle" {
  description = "Whether to configure lifecycle management"
  type        = bool
  default     = true
}

variable "minio_endpoint" {
  description = "MinIO endpoint URL"
  type        = string
  default     = "http://minio.ml-infrastructure.svc.cluster.local:9000"
}
