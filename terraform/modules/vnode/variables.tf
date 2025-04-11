variable "name" {
  description = "Name of the vNode runtime deployment"
  type        = string
  default     = "vnode-runtime"
}

variable "namespace" {
  description = "Namespace to deploy the vNode runtime"
  type        = string
  default     = "vcluster"
}

variable "replicas" {
  description = "Number of vNode runtime replicas"
  type        = number
  default     = 2
}

variable "vcluster_name" {
  description = "Name of the vCluster to integrate with"
  type        = string
  default     = "agent-runtime-vcluster"
}

variable "vcluster_namespace" {
  description = "Namespace of the vCluster to integrate with"
  type        = string
  default     = "vcluster"
}

variable "resources" {
  description = "Resource limits and requests for the vNode runtime"
  type = object({
    limits = object({
      cpu    = string
      memory = string
    })
    requests = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    limits = {
      cpu    = "500m"
      memory = "512Mi"
    }
    requests = {
      cpu    = "250m"
      memory = "256Mi"
    }
  }
}

variable "enable_kata_integration" {
  description = "Whether to enable Kata Containers integration"
  type        = bool
  default     = true
}

variable "node_names" {
  description = "List of node names to label for vNode integration"
  type        = list(string)
  default     = []
}
