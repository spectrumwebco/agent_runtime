variable "namespace" {
  description = "Kubernetes namespace to deploy ArgoCD"
  type        = string
  default     = "argocd"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Additional labels to add to resources"
  type        = map(string)
  default     = {}
}

variable "annotations" {
  description = "Additional annotations to add to resources"
  type        = map(string)
  default     = {}
}

variable "chart_version" {
  description = "Version of the ArgoCD Helm chart"
  type        = string
  default     = "5.51.4"
}

variable "values_yaml" {
  description = "Values YAML for the ArgoCD Helm chart"
  type        = string
  default     = ""
}

variable "server_service_type" {
  description = "Service type for ArgoCD server"
  type        = string
  default     = "ClusterIP"
}

variable "controller_cpu_limit" {
  description = "CPU limit for ArgoCD controller"
  type        = string
  default     = "500m"
}

variable "controller_memory_limit" {
  description = "Memory limit for ArgoCD controller"
  type        = string
  default     = "512Mi"
}

variable "controller_cpu_request" {
  description = "CPU request for ArgoCD controller"
  type        = string
  default     = "100m"
}

variable "controller_memory_request" {
  description = "Memory request for ArgoCD controller"
  type        = string
  default     = "128Mi"
}

variable "server_cpu_limit" {
  description = "CPU limit for ArgoCD server"
  type        = string
  default     = "500m"
}

variable "server_memory_limit" {
  description = "Memory limit for ArgoCD server"
  type        = string
  default     = "512Mi"
}

variable "server_cpu_request" {
  description = "CPU request for ArgoCD server"
  type        = string
  default     = "100m"
}

variable "server_memory_request" {
  description = "Memory request for ArgoCD server"
  type        = string
  default     = "128Mi"
}

variable "repo_server_cpu_limit" {
  description = "CPU limit for ArgoCD repo server"
  type        = string
  default     = "500m"
}

variable "repo_server_memory_limit" {
  description = "Memory limit for ArgoCD repo server"
  type        = string
  default     = "512Mi"
}

variable "repo_server_cpu_request" {
  description = "CPU request for ArgoCD repo server"
  type        = string
  default     = "100m"
}

variable "repo_server_memory_request" {
  description = "Memory request for ArgoCD repo server"
  type        = string
  default     = "128Mi"
}

variable "applications" {
  description = "List of ArgoCD applications to deploy"
  type = list(object({
    name                 = string
    project              = string
    repo_url             = string
    target_revision      = string
    path                 = string
    destination_server   = string
    destination_namespace = string
    prune                = bool
    self_heal            = bool
    allow_empty          = bool
    sync_options         = list(string)
  }))
  default = []
}
