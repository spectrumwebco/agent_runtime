variable "namespace" {
  description = "Kubernetes namespace for Neovim deployment"
  type        = string
  default     = "neovim"
}

variable "replicas" {
  description = "Number of Neovim replicas"
  type        = number
  default     = 1
}

variable "image" {
  description = "Neovim container image"
  type        = string
  default     = "neovim/neovim:latest"
}

variable "supabase_url" {
  description = "Supabase URL for state persistence"
  type        = string
}

variable "supabase_key" {
  description = "Supabase key for state persistence"
  type        = string
  sensitive   = true
}

variable "storage_size" {
  description = "Size of persistent volume for Neovim data"
  type        = string
  default     = "1Gi"
}

variable "kata_runtime_class" {
  description = "Kata Containers runtime class"
  type        = string
  default     = "kata"
}

variable "enable_kata" {
  description = "Whether to enable Kata Containers for Neovim"
  type        = bool
  default     = true
}

variable "resource_limits" {
  description = "Resource limits for Neovim container"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "200m"
    memory = "512Mi"
  }
}

variable "resource_requests" {
  description = "Resource requests for Neovim container"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "100m"
    memory = "256Mi"
  }
}
