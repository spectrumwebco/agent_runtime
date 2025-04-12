variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "namespace" {
  description = "Namespace for vNode runtime"
  type        = string
  default     = "vnode-runtime"
}

variable "vnode_runtime_version" {
  description = "Version of vNode runtime to deploy"
  type        = string
  default     = "0.0.2"
}

variable "replica_count" {
  description = "Number of vNode runtime replicas"
  type        = number
  default     = 1
}

variable "resources_limits_cpu" {
  description = "CPU limits for vNode runtime"
  type        = string
  default     = "500m"
}

variable "resources_limits_memory" {
  description = "Memory limits for vNode runtime"
  type        = string
  default     = "512Mi"
}

variable "resources_requests_cpu" {
  description = "CPU requests for vNode runtime"
  type        = string
  default     = "250m"
}

variable "resources_requests_memory" {
  description = "Memory requests for vNode runtime"
  type        = string
  default     = "256Mi"
}
