variable "namespace" {
  description = "Kubernetes namespace for RocketMQ deployment"
  type        = string
  default     = "rocketmq"
}

variable "create_namespace" {
  description = "Whether to create the RocketMQ namespace"
  type        = bool
  default     = true
}

variable "name" {
  description = "Name for RocketMQ resources"
  type        = string
  default     = "rocketmq"
}

variable "chart_version" {
  description = "Version of the RocketMQ Helm chart"
  type        = string
  default     = "0.3.0"
}

variable "image_repository" {
  description = "RocketMQ image repository"
  type        = string
  default     = "apache/rocketmq"
}

variable "image_tag" {
  description = "RocketMQ image tag"
  type        = string
  default     = "4.9.4"
}

variable "high_availability" {
  description = "Whether to enable high availability for RocketMQ"
  type        = bool
  default     = true
}

variable "name_server_replicas" {
  description = "Number of RocketMQ name server replicas"
  type        = number
  default     = 3
}

variable "broker_replicas" {
  description = "Number of RocketMQ broker replicas"
  type        = number
  default     = 3
}

variable "dashboard_enabled" {
  description = "Whether to enable RocketMQ dashboard"
  type        = bool
  default     = true
}

variable "storage_size" {
  description = "Storage size for RocketMQ brokers"
  type        = string
  default     = "20Gi"
}

variable "storage_class" {
  description = "Storage class for RocketMQ persistent volumes"
  type        = string
  default     = null
}

variable "resource_limits" {
  description = "Resource limits for RocketMQ components"
  type = object({
    name_server = object({
      cpu    = string
      memory = string
    })
    broker = object({
      cpu    = string
      memory = string
    })
    dashboard = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    name_server = {
      cpu    = "1000m"
      memory = "2Gi"
    }
    broker = {
      cpu    = "2000m"
      memory = "4Gi"
    }
    dashboard = {
      cpu    = "500m"
      memory = "1Gi"
    }
  }
}

variable "resource_requests" {
  description = "Resource requests for RocketMQ components"
  type = object({
    name_server = object({
      cpu    = string
      memory = string
    })
    broker = object({
      cpu    = string
      memory = string
    })
    dashboard = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    name_server = {
      cpu    = "500m"
      memory = "1Gi"
    }
    broker = {
      cpu    = "1000m"
      memory = "2Gi"
    }
    dashboard = {
      cpu    = "200m"
      memory = "512Mi"
    }
  }
}

variable "broker_config" {
  description = "Additional broker configuration"
  type        = map(string)
  default     = {}
}

variable "name_server_config" {
  description = "Additional name server configuration"
  type        = map(string)
  default     = {}
}

variable "ingress_enabled" {
  description = "Whether to enable ingress for RocketMQ dashboard"
  type        = bool
  default     = false
}

variable "ingress_domain" {
  description = "Domain for RocketMQ dashboard ingress"
  type        = string
  default     = ""
}

variable "ingress_class" {
  description = "Ingress class for RocketMQ dashboard"
  type        = string
  default     = "nginx"
}

variable "ingress_tls_enabled" {
  description = "Whether to enable TLS for RocketMQ dashboard ingress"
  type        = bool
  default     = false
}

variable "ingress_tls_secret" {
  description = "TLS secret for RocketMQ dashboard ingress"
  type        = string
  default     = ""
}

variable "prometheus_integration" {
  description = "Whether to enable Prometheus integration"
  type        = bool
  default     = true
}

variable "kata_container_integration" {
  description = "Whether to enable Kata Containers integration"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Additional labels for RocketMQ resources"
  type        = map(string)
  default     = {}
}

variable "annotations" {
  description = "Additional annotations for RocketMQ resources"
  type        = map(string)
  default     = {}
}

variable "acl_enabled" {
  description = "Whether to enable ACL for RocketMQ"
  type        = bool
  default     = true
}

variable "acl_access_key" {
  description = "Access key for RocketMQ ACL"
  type        = string
  default     = "rocketmq"
  sensitive   = true
}

variable "acl_secret_key" {
  description = "Secret key for RocketMQ ACL"
  type        = string
  default     = ""
  sensitive   = true
}

variable "topic_configs" {
  description = "Configuration for RocketMQ topics"
  type = list(object({
    name            = string
    read_queue_nums = number
    write_queue_nums = number
    perm            = string
  }))
  default = [
    {
      name            = "agent-events"
      read_queue_nums = 8
      write_queue_nums = 8
      perm            = "READ | WRITE"
    },
    {
      name            = "k8s-lifecycle"
      read_queue_nums = 8
      write_queue_nums = 8
      perm            = "READ | WRITE"
    },
    {
      name            = "kata-lifecycle"
      read_queue_nums = 8
      write_queue_nums = 8
      perm            = "READ | WRITE"
    },
    {
      name            = "state-updates"
      read_queue_nums = 16
      write_queue_nums = 16
      perm            = "READ | WRITE"
    }
  ]
}
