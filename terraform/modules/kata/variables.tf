# Variables for Kata Containers Module
variable "namespace" {
  description = "Kubernetes namespace for Agent Runtime"
  type        = string
}

variable "node_selector" {
  description = "Node selector for Kata Containers"
  type        = map(string)
  default     = {
    "kata-containers" = "true"
  }
}

variable "kata_ca_cert" {
  description = "Kata Containers CA certificate"
  type        = string
  default     = ""
  sensitive   = true
}

variable "kata_server_cert" {
  description = "Kata Containers server certificate"
  type        = string
  default     = ""
  sensitive   = true
}

variable "kata_server_key" {
  description = "Kata Containers server key"
  type        = string
  default     = ""
  sensitive   = true
}
