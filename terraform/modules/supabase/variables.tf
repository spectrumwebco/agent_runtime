variable "namespace" {
  description = "Namespace for Supabase resources"
  type        = string
  default     = "supabase-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "container_registry" {
  description = "Container registry for Supabase images"
  type        = string
  default     = "ghcr.io/spectrumwebco"
}

variable "supabase_version" {
  description = "Version of Supabase to deploy"
  type        = string
  default     = "latest"
}

variable "replicas" {
  description = "Number of replicas for Supabase deployment"
  type        = number
  default     = 1
}

variable "postgres_password" {
  description = "Password for Postgres"
  type        = string
  sensitive   = true
  default     = ""
}

variable "jwt_secret" {
  description = "JWT secret for Supabase"
  type        = string
  sensitive   = true
  default     = ""
}

variable "anon_key" {
  description = "Anonymous key for Supabase"
  type        = string
  sensitive   = true
  default     = ""
}

variable "service_role_key" {
  description = "Service role key for Supabase"
  type        = string
  sensitive   = true
  default     = ""
}
