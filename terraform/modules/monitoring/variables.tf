variable "namespace" {
  description = "Namespace for monitoring components"
  type        = string
  default     = "monitoring"
}

variable "create_namespace" {
  description = "Whether to create the monitoring namespace"
  type        = bool
  default     = true
}

variable "prometheus_enabled" {
  description = "Whether to enable Prometheus"
  type        = bool
  default     = true
}

variable "prometheus_version" {
  description = "Version of Prometheus Helm chart"
  type        = string
  default     = "19.6.1"
}

variable "prometheus_retention" {
  description = "Retention period for Prometheus metrics"
  type        = string
  default     = "15d"
}

variable "prometheus_storage_size" {
  description = "Storage size for Prometheus"
  type        = string
  default     = "50Gi"
}

variable "prometheus_storage_class" {
  description = "Storage class for Prometheus"
  type        = string
  default     = null
}

variable "prometheus_replicas" {
  description = "Number of Prometheus replicas"
  type        = number
  default     = 2
}

variable "grafana_enabled" {
  description = "Whether to enable Grafana"
  type        = bool
  default     = true
}

variable "grafana_version" {
  description = "Version of Grafana Helm chart"
  type        = string
  default     = "6.56.6"
}

variable "grafana_admin_password" {
  description = "Admin password for Grafana"
  type        = string
  default     = "admin"
  sensitive   = true
}

variable "grafana_storage_size" {
  description = "Storage size for Grafana"
  type        = string
  default     = "10Gi"
}

variable "grafana_storage_class" {
  description = "Storage class for Grafana"
  type        = string
  default     = null
}

variable "thanos_enabled" {
  description = "Whether to enable Thanos"
  type        = bool
  default     = true
}

variable "thanos_version" {
  description = "Version of Thanos Helm chart"
  type        = string
  default     = "12.5.1"
}

variable "thanos_storage_size" {
  description = "Storage size for Thanos"
  type        = string
  default     = "100Gi"
}

variable "thanos_storage_class" {
  description = "Storage class for Thanos"
  type        = string
  default     = null
}

variable "thanos_objstore_config" {
  description = "Object store configuration for Thanos"
  type        = string
  default     = ""
  sensitive   = true
}

variable "loki_enabled" {
  description = "Whether to enable Loki"
  type        = bool
  default     = true
}

variable "loki_version" {
  description = "Version of Loki Helm chart"
  type        = string
  default     = "5.8.9"
}

variable "loki_storage_size" {
  description = "Storage size for Loki"
  type        = string
  default     = "50Gi"
}

variable "loki_storage_class" {
  description = "Storage class for Loki"
  type        = string
  default     = null
}

variable "loki_retention" {
  description = "Retention period for Loki logs"
  type        = string
  default     = "168h"
}

variable "jaeger_enabled" {
  description = "Whether to enable Jaeger"
  type        = bool
  default     = true
}

variable "jaeger_version" {
  description = "Version of Jaeger Helm chart"
  type        = string
  default     = "0.71.6"
}

variable "jaeger_storage_type" {
  description = "Storage type for Jaeger (memory, elasticsearch, cassandra)"
  type        = string
  default     = "elasticsearch"
}

variable "jaeger_elasticsearch_host" {
  description = "Elasticsearch host for Jaeger"
  type        = string
  default     = "elasticsearch-master.monitoring.svc.cluster.local"
}

variable "elasticsearch_enabled" {
  description = "Whether to enable Elasticsearch for ELK stack"
  type        = bool
  default     = true
}

variable "elasticsearch_version" {
  description = "Version of Elasticsearch Helm chart"
  type        = string
  default     = "19.5.7"
}

variable "elasticsearch_storage_size" {
  description = "Storage size for Elasticsearch"
  type        = string
  default     = "100Gi"
}

variable "elasticsearch_storage_class" {
  description = "Storage class for Elasticsearch"
  type        = string
  default     = null
}

variable "elasticsearch_replicas" {
  description = "Number of Elasticsearch replicas"
  type        = number
  default     = 3
}

variable "kibana_enabled" {
  description = "Whether to enable Kibana for ELK stack"
  type        = bool
  default     = true
}

variable "kibana_version" {
  description = "Version of Kibana Helm chart"
  type        = string
  default     = "10.4.1"
}

variable "filebeat_enabled" {
  description = "Whether to enable Filebeat for ELK stack"
  type        = bool
  default     = true
}

variable "filebeat_version" {
  description = "Version of Filebeat Helm chart"
  type        = string
  default     = "7.17.3"
}

variable "vector_enabled" {
  description = "Whether to enable Vector"
  type        = bool
  default     = true
}

variable "vector_version" {
  description = "Version of Vector Helm chart"
  type        = string
  default     = "0.28.0"
}

variable "opentelemetry_enabled" {
  description = "Whether to enable OpenTelemetry"
  type        = bool
  default     = true
}

variable "opentelemetry_version" {
  description = "Version of OpenTelemetry Helm chart"
  type        = string
  default     = "0.41.0"
}

variable "kube_state_metrics_enabled" {
  description = "Whether to enable kube-state-metrics"
  type        = bool
  default     = true
}

variable "kube_state_metrics_version" {
  description = "Version of kube-state-metrics Helm chart"
  type        = string
  default     = "5.8.0"
}

variable "cadvisor_enabled" {
  description = "Whether to enable cAdvisor"
  type        = bool
  default     = true
}

variable "cadvisor_version" {
  description = "Version of cAdvisor image"
  type        = string
  default     = "v0.47.2"
}

variable "kubernetes_dashboard_enabled" {
  description = "Whether to enable Kubernetes Dashboard"
  type        = bool
  default     = true
}

variable "kubernetes_dashboard_version" {
  description = "Version of Kubernetes Dashboard Helm chart"
  type        = string
  default     = "6.0.8"
}

variable "high_availability" {
  description = "Whether to enable high availability for monitoring components"
  type        = bool
  default     = true
}

variable "resource_limits" {
  description = "Resource limits for monitoring components"
  type = object({
    prometheus = object({
      cpu    = string
      memory = string
    })
    grafana = object({
      cpu    = string
      memory = string
    })
    thanos = object({
      cpu    = string
      memory = string
    })
    loki = object({
      cpu    = string
      memory = string
    })
    jaeger = object({
      cpu    = string
      memory = string
    })
    elasticsearch = object({
      cpu    = string
      memory = string
    })
    vector = object({
      cpu    = string
      memory = string
    })
    opentelemetry = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    prometheus = {
      cpu    = "2000m"
      memory = "4Gi"
    }
    grafana = {
      cpu    = "500m"
      memory = "1Gi"
    }
    thanos = {
      cpu    = "1000m"
      memory = "2Gi"
    }
    loki = {
      cpu    = "1000m"
      memory = "2Gi"
    }
    jaeger = {
      cpu    = "1000m"
      memory = "2Gi"
    }
    elasticsearch = {
      cpu    = "2000m"
      memory = "4Gi"
    }
    vector = {
      cpu    = "500m"
      memory = "1Gi"
    }
    opentelemetry = {
      cpu    = "1000m"
      memory = "2Gi"
    }
  }
}

variable "resource_requests" {
  description = "Resource requests for monitoring components"
  type = object({
    prometheus = object({
      cpu    = string
      memory = string
    })
    grafana = object({
      cpu    = string
      memory = string
    })
    thanos = object({
      cpu    = string
      memory = string
    })
    loki = object({
      cpu    = string
      memory = string
    })
    jaeger = object({
      cpu    = string
      memory = string
    })
    elasticsearch = object({
      cpu    = string
      memory = string
    })
    vector = object({
      cpu    = string
      memory = string
    })
    opentelemetry = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    prometheus = {
      cpu    = "500m"
      memory = "2Gi"
    }
    grafana = {
      cpu    = "200m"
      memory = "512Mi"
    }
    thanos = {
      cpu    = "500m"
      memory = "1Gi"
    }
    loki = {
      cpu    = "500m"
      memory = "1Gi"
    }
    jaeger = {
      cpu    = "500m"
      memory = "1Gi"
    }
    elasticsearch = {
      cpu    = "1000m"
      memory = "2Gi"
    }
    vector = {
      cpu    = "200m"
      memory = "512Mi"
    }
    opentelemetry = {
      cpu    = "500m"
      memory = "1Gi"
    }
  }
}

variable "ingress_enabled" {
  description = "Whether to enable ingress for monitoring components"
  type        = bool
  default     = false
}

variable "ingress_domain" {
  description = "Domain for ingress"
  type        = string
  default     = ""
}

variable "ingress_class" {
  description = "Ingress class"
  type        = string
  default     = "nginx"
}

variable "ingress_tls_enabled" {
  description = "Whether to enable TLS for ingress"
  type        = bool
  default     = false
}

variable "ingress_tls_secret" {
  description = "TLS secret for ingress"
  type        = string
  default     = ""
}

variable "labels" {
  description = "Additional labels for monitoring resources"
  type        = map(string)
  default     = {}
}

variable "annotations" {
  description = "Additional annotations for monitoring resources"
  type        = map(string)
  default     = {}
}
