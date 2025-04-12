
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "kubeflow_version" {
  description = "Version of KubeFlow to deploy"
  type        = string
  default     = "1.7.0"
}

variable "training_operator_version" {
  description = "Version of KubeFlow Training Operator to deploy"
  type        = string
  default     = "v1.7.0"
}

variable "katib_version" {
  description = "Version of KubeFlow Katib to deploy"
  type        = string
  default     = "v0.15.0"
}

variable "kubeflow_data_storage_size" {
  description = "Size of persistent volume claim for KubeFlow data"
  type        = string
  default     = "100Gi"
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

variable "mlflow_tracking_uri" {
  description = "MLFlow tracking URI"
  type        = string
  default     = "http://mlflow-server.mlflow.svc.cluster.local:5000"
}

variable "kserve_version" {
  description = "Version of KServe to deploy"
  type        = string
  default     = "0.10.0"
}

variable "kserve_namespace" {
  description = "Namespace for KServe deployment"
  type        = string
  default     = "kserve"
}

variable "minio_version" {
  description = "Version of MinIO to deploy"
  type        = string
  default     = "RELEASE.2023-05-04T21-44-30Z"
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
