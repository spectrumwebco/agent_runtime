variable "namespace" {
  description = "Kubernetes namespace to deploy Flux System"
  type        = string
  default     = "flux-system"
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

variable "identity" {
  description = "SSH private key for Git authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "identity_pub" {
  description = "SSH public key for Git authentication"
  type        = string
  default     = ""
}

variable "known_hosts" {
  description = "Known hosts for SSH authentication"
  type        = string
  default     = ""
}

variable "git_repository_name" {
  description = "Name of the Git repository"
  type        = string
  default     = "flux-system"
}

variable "git_repository_url" {
  description = "URL of the Git repository"
  type        = string
}

variable "git_branch" {
  description = "Branch of the Git repository"
  type        = string
  default     = "main"
}

variable "sync_interval" {
  description = "Interval for synchronization"
  type        = string
  default     = "1m0s"
}

variable "kustomization_name" {
  description = "Name of the Kustomization"
  type        = string
  default     = "flux-system"
}

variable "kustomization_path" {
  description = "Path to the Kustomization"
  type        = string
  default     = "./"
}

variable "prune" {
  description = "Whether to prune resources"
  type        = bool
  default     = true
}
