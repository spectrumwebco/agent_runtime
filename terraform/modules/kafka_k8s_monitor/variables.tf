variable "namespace" {
  description = "Kubernetes namespace for Kafka and K8s monitor"
  type        = string
  default     = "default"
}

variable "kafka_replicas" {
  description = "Number of Kafka replicas"
  type        = number
  default     = 1
}

variable "monitor_namespace" {
  description = "Kubernetes namespace to monitor"
  type        = string
  default     = "default"
}

variable "poll_interval" {
  description = "Poll interval in seconds"
  type        = number
  default     = 30
}

variable "resources_to_monitor" {
  description = "Comma-separated list of resources to monitor"
  type        = string
  default     = "pods,services,deployments,statefulsets,configmaps,secrets"
}
