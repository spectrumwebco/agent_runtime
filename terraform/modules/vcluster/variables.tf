variable "namespace" {
  description = "Namespace for vCluster resources"
  type        = string
  default     = "vcluster-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "vcluster_version" {
  description = "Version of vCluster to deploy"
  type        = string
  default     = "0.15.0"
}

variable "kube_config" {
  description = "Kubernetes config for vCluster"
  type        = string
  sensitive   = true
  default     = ""
}
