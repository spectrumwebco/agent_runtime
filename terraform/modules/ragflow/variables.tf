variable "namespace" {
  description = "Kubernetes namespace for RAGflow deployment"
  type        = string
  default     = "ragflow"
}

variable "create_namespace" {
  description = "Whether to create the RAGflow namespace"
  type        = bool
  default     = true
}

variable "name" {
  description = "Name for RAGflow resources"
  type        = string
  default     = "ragflow"
}

variable "vcluster_enabled" {
  description = "Whether to deploy RAGflow in a vCluster"
  type        = bool
  default     = true
}

variable "vcluster_name" {
  description = "Name of the vCluster for RAGflow"
  type        = string
  default     = "ragflow-vcluster"
}

variable "vcluster_namespace" {
  description = "Namespace for the vCluster"
  type        = string
  default     = "vcluster"
}

variable "ragflow_version" {
  description = "Version of RAGflow to deploy"
  type        = string
  default     = "latest"
}

variable "ragflow_image" {
  description = "Docker image for RAGflow"
  type        = string
  default     = "ragflow/ragflow"
}

variable "ragflow_image_tag" {
  description = "Docker image tag for RAGflow"
  type        = string
  default     = "latest"
}

variable "ragflow_replicas" {
  description = "Number of RAGflow replicas"
  type        = number
  default     = 3
}

variable "ragflow_resources" {
  description = "Resource limits and requests for RAGflow"
  type = object({
    limits = object({
      cpu    = string
      memory = string
    })
    requests = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    limits = {
      cpu    = "2000m"
      memory = "4Gi"
    }
    requests = {
      cpu    = "1000m"
      memory = "2Gi"
    }
  }
}

variable "ragflow_storage_size" {
  description = "Storage size for RAGflow"
  type        = string
  default     = "20Gi"
}

variable "ragflow_storage_class" {
  description = "Storage class for RAGflow"
  type        = string
  default     = null
}

variable "enable_high_availability" {
  description = "Whether to enable high availability for RAGflow"
  type        = bool
  default     = true
}

variable "enable_prometheus_integration" {
  description = "Whether to enable Prometheus integration"
  type        = bool
  default     = true
}

variable "enable_jaeger_integration" {
  description = "Whether to enable Jaeger integration"
  type        = bool
  default     = true
}

variable "enable_opentelemetry_integration" {
  description = "Whether to enable OpenTelemetry integration"
  type        = bool
  default     = true
}

variable "enable_loki_integration" {
  description = "Whether to enable Loki integration"
  type        = bool
  default     = true
}

variable "enable_vector_integration" {
  description = "Whether to enable Vector integration"
  type        = bool
  default     = true
}

variable "enable_kata_container_integration" {
  description = "Whether to enable Kata Containers integration"
  type        = bool
  default     = true
}

variable "enable_ingress" {
  description = "Whether to enable ingress for RAGflow"
  type        = bool
  default     = false
}

variable "ingress_domain" {
  description = "Domain for RAGflow ingress"
  type        = string
  default     = ""
}

variable "ingress_class" {
  description = "Ingress class for RAGflow"
  type        = string
  default     = "nginx"
}

variable "ingress_tls_enabled" {
  description = "Whether to enable TLS for RAGflow ingress"
  type        = bool
  default     = false
}

variable "ingress_tls_secret" {
  description = "TLS secret for RAGflow ingress"
  type        = string
  default     = ""
}

variable "ragflow_config" {
  description = "Configuration for RAGflow"
  type = object({
    embedding_model = string
    llm_model       = string
    vector_db       = string
    document_store  = string
    api_key         = string
  })
  default = {
    embedding_model = "sentence-transformers/all-mpnet-base-v2"
    llm_model       = "gpt-3.5-turbo"
    vector_db       = "qdrant"
    document_store  = "minio"
    api_key         = ""
  }
  sensitive = true
}

variable "ragflow_qdrant_config" {
  description = "Configuration for Qdrant vector database"
  type = object({
    enabled     = bool
    url         = string
    api_key     = string
    collection  = string
    replicas    = number
    resources = object({
      limits = object({
        cpu    = string
        memory = string
      })
      requests = object({
        cpu    = string
        memory = string
      })
    })
  })
  default = {
    enabled    = true
    url        = ""
    api_key    = ""
    collection = "ragflow"
    replicas   = 3
    resources = {
      limits = {
        cpu    = "1000m"
        memory = "2Gi"
      }
      requests = {
        cpu    = "500m"
        memory = "1Gi"
      }
    }
  }
  sensitive = true
}

variable "ragflow_minio_config" {
  description = "Configuration for MinIO document store"
  type = object({
    enabled     = bool
    url         = string
    access_key  = string
    secret_key  = string
    bucket      = string
    replicas    = number
    resources = object({
      limits = object({
        cpu    = string
        memory = string
      })
      requests = object({
        cpu    = string
        memory = string
      })
    })
  })
  default = {
    enabled    = true
    url        = ""
    access_key = ""
    secret_key = ""
    bucket     = "ragflow"
    replicas   = 3
    resources = {
      limits = {
        cpu    = "1000m"
        memory = "2Gi"
      }
      requests = {
        cpu    = "500m"
        memory = "1Gi"
      }
    }
  }
  sensitive = true
}

variable "ragflow_redis_config" {
  description = "Configuration for Redis cache"
  type = object({
    enabled     = bool
    url         = string
    password    = string
    replicas    = number
    resources = object({
      limits = object({
        cpu    = string
        memory = string
      })
      requests = object({
        cpu    = string
        memory = string
      })
    })
  })
  default = {
    enabled  = true
    url      = ""
    password = ""
    replicas = 3
    resources = {
      limits = {
        cpu    = "1000m"
        memory = "2Gi"
      }
      requests = {
        cpu    = "500m"
        memory = "1Gi"
      }
    }
  }
  sensitive = true
}

variable "ragflow_api_config" {
  description = "Configuration for RAGflow API"
  type = object({
    port           = number
    max_tokens     = number
    temperature    = number
    top_p          = number
    top_k          = number
    chunk_size     = number
    chunk_overlap  = number
    max_documents  = number
  })
  default = {
    port          = 8000
    max_tokens    = 1024
    temperature   = 0.7
    top_p         = 0.95
    top_k         = 40
    chunk_size    = 1000
    chunk_overlap = 200
    max_documents = 10
  }
}

variable "ragflow_labels" {
  description = "Additional labels for RAGflow resources"
  type        = map(string)
  default     = {}
}

variable "ragflow_annotations" {
  description = "Additional annotations for RAGflow resources"
  type        = map(string)
  default     = {}
}
