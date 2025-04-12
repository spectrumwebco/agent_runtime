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
  default     = "v1.21.4-k3s1"
}
