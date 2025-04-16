variable "namespace" {
  description = "Kubernetes namespace for Kafka"
  type        = string
  default     = "default"
}

variable "kafka_version" {
  description = "Kafka version"
  type        = string
  default     = "3.4.0"
}

variable "zookeeper_version" {
  description = "ZooKeeper version"
  type        = string
  default     = "3.8.1"
}

variable "kafka_replicas" {
  description = "Number of Kafka replicas"
  type        = number
  default     = 1
}

variable "kafka_heap_size" {
  description = "Kafka heap size"
  type        = string
  default     = "1G"
}

variable "kafka_memory_request" {
  description = "Memory request for Kafka"
  type        = string
  default     = "2Gi"
}

variable "kafka_memory_limit" {
  description = "Memory limit for Kafka"
  type        = string
  default     = "4Gi"
}

variable "kafka_cpu_request" {
  description = "CPU request for Kafka"
  type        = string
  default     = "500m"
}

variable "kafka_cpu_limit" {
  description = "CPU limit for Kafka"
  type        = string
  default     = "1000m"
}

variable "kafka_storage_size" {
  description = "Storage size for Kafka"
  type        = string
  default     = "10Gi"
}

variable "zookeeper_memory_request" {
  description = "Memory request for ZooKeeper"
  type        = string
  default     = "512Mi"
}

variable "zookeeper_memory_limit" {
  description = "Memory limit for ZooKeeper"
  type        = string
  default     = "1Gi"
}

variable "zookeeper_cpu_request" {
  description = "CPU request for ZooKeeper"
  type        = string
  default     = "250m"
}

variable "zookeeper_cpu_limit" {
  description = "CPU limit for ZooKeeper"
  type        = string
  default     = "500m"
}

variable "zookeeper_storage_size" {
  description = "Storage size for ZooKeeper"
  type        = string
  default     = "5Gi"
}
