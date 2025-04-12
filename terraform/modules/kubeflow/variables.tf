
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
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

variable "storage_class_name" {
  description = "Storage class name for persistent volume claims"
  type        = string
  default     = "standard"
}
