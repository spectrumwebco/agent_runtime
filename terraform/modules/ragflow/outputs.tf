output "namespace" {
  description = "Namespace where RAGflow is deployed"
  value       = local.namespace
}

output "service_name" {
  description = "Name of the RAGflow service"
  value       = kubernetes_service.ragflow.metadata[0].name
}

output "service_port" {
  description = "Port of the RAGflow service"
  value       = var.ragflow_api_config.port
}

output "api_url" {
  description = "URL for RAGflow API"
  value       = "http://${kubernetes_service.ragflow.metadata[0].name}.${local.namespace}.svc.cluster.local:${var.ragflow_api_config.port}"
}

output "vcluster_enabled" {
  description = "Whether RAGflow is deployed in a vCluster"
  value       = var.vcluster_enabled
}

output "vcluster_name" {
  description = "Name of the vCluster for RAGflow"
  value       = var.vcluster_enabled ? var.vcluster_name : null
}

output "vcluster_namespace" {
  description = "Namespace of the vCluster for RAGflow"
  value       = var.vcluster_enabled ? var.vcluster_namespace : null
}

output "qdrant_service_name" {
  description = "Name of the Qdrant service"
  value       = var.ragflow_qdrant_config.enabled && var.ragflow_qdrant_config.url == "" ? kubernetes_service.qdrant[0].metadata[0].name : null
}

output "qdrant_url" {
  description = "URL for Qdrant"
  value       = var.ragflow_qdrant_config.url != "" ? var.ragflow_qdrant_config.url : var.ragflow_qdrant_config.enabled && var.ragflow_qdrant_config.url == "" ? "http://${kubernetes_service.qdrant[0].metadata[0].name}.${local.namespace}.svc.cluster.local:6333" : null
}

output "minio_service_name" {
  description = "Name of the MinIO service"
  value       = var.ragflow_minio_config.enabled && var.ragflow_minio_config.url == "" ? kubernetes_service.minio[0].metadata[0].name : null
}

output "minio_url" {
  description = "URL for MinIO"
  value       = var.ragflow_minio_config.url != "" ? var.ragflow_minio_config.url : var.ragflow_minio_config.enabled && var.ragflow_minio_config.url == "" ? "http://${kubernetes_service.minio[0].metadata[0].name}.${local.namespace}.svc.cluster.local:9000" : null
}

output "minio_console_url" {
  description = "URL for MinIO Console"
  value       = var.ragflow_minio_config.url != "" ? var.ragflow_minio_config.url : var.ragflow_minio_config.enabled && var.ragflow_minio_config.url == "" ? "http://${kubernetes_service.minio[0].metadata[0].name}.${local.namespace}.svc.cluster.local:9001" : null
}

output "redis_service_name" {
  description = "Name of the Redis service"
  value       = var.ragflow_redis_config.enabled && var.ragflow_redis_config.url == "" ? kubernetes_service.redis[0].metadata[0].name : null
}

output "redis_url" {
  description = "URL for Redis"
  value       = var.ragflow_redis_config.url != "" ? var.ragflow_redis_config.url : var.ragflow_redis_config.enabled && var.ragflow_redis_config.url == "" ? "redis://${kubernetes_service.redis[0].metadata[0].name}.${local.namespace}.svc.cluster.local:6379" : null
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled for RAGflow"
  value       = var.enable_high_availability
}

output "prometheus_integration_enabled" {
  description = "Whether Prometheus integration is enabled"
  value       = var.enable_prometheus_integration
}

output "jaeger_integration_enabled" {
  description = "Whether Jaeger integration is enabled"
  value       = var.enable_jaeger_integration
}

output "opentelemetry_integration_enabled" {
  description = "Whether OpenTelemetry integration is enabled"
  value       = var.enable_opentelemetry_integration
}

output "loki_integration_enabled" {
  description = "Whether Loki integration is enabled"
  value       = var.enable_loki_integration
}

output "vector_integration_enabled" {
  description = "Whether Vector integration is enabled"
  value       = var.enable_vector_integration
}

output "kata_container_integration_enabled" {
  description = "Whether Kata Containers integration is enabled"
  value       = var.enable_kata_container_integration
}

output "ingress_enabled" {
  description = "Whether ingress is enabled for RAGflow"
  value       = var.enable_ingress
}

output "ingress_url" {
  description = "URL for RAGflow ingress"
  value       = var.enable_ingress ? "http${var.ingress_tls_enabled ? "s" : ""}://${var.ingress_domain}" : null
}

output "config_map_name" {
  description = "Name of the RAGflow config map"
  value       = kubernetes_config_map.ragflow_config.metadata[0].name
}

output "api_key_secret_name" {
  description = "Name of the RAGflow API key secret"
  value       = kubernetes_secret.ragflow_api_key.metadata[0].name
}

output "embedding_model" {
  description = "Embedding model used by RAGflow"
  value       = var.ragflow_config.embedding_model
}

output "llm_model" {
  description = "LLM model used by RAGflow"
  value       = var.ragflow_config.llm_model
}

output "vector_db" {
  description = "Vector database used by RAGflow"
  value       = var.ragflow_config.vector_db
}

output "document_store" {
  description = "Document store used by RAGflow"
  value       = var.ragflow_config.document_store
}
