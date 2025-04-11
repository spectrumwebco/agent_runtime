# Variables for Kubernetes Module
variable "namespace" {
  description = "Kubernetes namespace for Agent Runtime"
  type        = string
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
}

variable "vcluster_enabled" {
  description = "Enable vCluster deployment"
  type        = bool
}

variable "vcluster_version" {
  description = "vCluster version to deploy"
  type        = string
}

variable "jspolicy_enabled" {
  description = "Enable jsPolicy deployment"
  type        = bool
}
