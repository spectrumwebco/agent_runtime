variable "namespace" {
  description = "Kubernetes namespace for Apache Doris"
  type        = string
  default     = "default"
}

variable "doris_version" {
  description = "Apache Doris version"
  type        = string
  default     = "2.0.2"
}

variable "fe_replicas" {
  description = "Number of Frontend replicas"
  type        = number
  default     = 1
}

variable "be_replicas" {
  description = "Number of Backend replicas"
  type        = number
  default     = 1
}

variable "fe_memory_request" {
  description = "Memory request for Frontend"
  type        = string
  default     = "2Gi"
}

variable "fe_memory_limit" {
  description = "Memory limit for Frontend"
  type        = string
  default     = "4Gi"
}

variable "fe_cpu_request" {
  description = "CPU request for Frontend"
  type        = string
  default     = "1000m"
}

variable "fe_cpu_limit" {
  description = "CPU limit for Frontend"
  type        = string
  default     = "2000m"
}

variable "be_memory_request" {
  description = "Memory request for Backend"
  type        = string
  default     = "4Gi"
}

variable "be_memory_limit" {
  description = "Memory limit for Backend"
  type        = string
  default     = "8Gi"
}

variable "be_cpu_request" {
  description = "CPU request for Backend"
  type        = string
  default     = "2000m"
}

variable "be_cpu_limit" {
  description = "CPU limit for Backend"
  type        = string
  default     = "4000m"
}

variable "fe_storage_size" {
  description = "Storage size for Frontend"
  type        = string
  default     = "10Gi"
}

variable "be_storage_size" {
  description = "Storage size for Backend"
  type        = string
  default     = "20Gi"
}

variable "admin_password" {
  description = "Admin password for Apache Doris"
  type        = string
  sensitive   = true
}
