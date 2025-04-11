output "namespace" {
  description = "Namespace where monitoring components are deployed"
  value       = local.namespace
}

output "prometheus_server_endpoint" {
  description = "Endpoint for Prometheus server"
  value       = var.prometheus_enabled ? "http://prometheus-server.${local.namespace}.svc.cluster.local:80" : null
}

output "grafana_endpoint" {
  description = "Endpoint for Grafana"
  value       = var.grafana_enabled ? "http://grafana.${local.namespace}.svc.cluster.local:80" : null
}

output "thanos_query_endpoint" {
  description = "Endpoint for Thanos Query"
  value       = var.thanos_enabled ? "http://thanos-query.${local.namespace}.svc.cluster.local:9090" : null
}

output "loki_endpoint" {
  description = "Endpoint for Loki"
  value       = var.loki_enabled ? "http://loki-gateway.${local.namespace}.svc.cluster.local:80" : null
}

output "jaeger_query_endpoint" {
  description = "Endpoint for Jaeger Query"
  value       = var.jaeger_enabled ? "http://jaeger-query.${local.namespace}.svc.cluster.local:16686" : null
}

output "elasticsearch_endpoint" {
  description = "Endpoint for Elasticsearch"
  value       = var.elasticsearch_enabled ? "http://elasticsearch-master.${local.namespace}.svc.cluster.local:9200" : null
}

output "kibana_endpoint" {
  description = "Endpoint for Kibana"
  value       = var.elasticsearch_enabled && var.kibana_enabled ? "http://kibana.${local.namespace}.svc.cluster.local:5601" : null
}

output "opentelemetry_collector_endpoint" {
  description = "Endpoint for OpenTelemetry Collector"
  value       = var.opentelemetry_enabled ? "http://opentelemetry-collector.${local.namespace}.svc.cluster.local:4317" : null
}

output "kubernetes_dashboard_endpoint" {
  description = "Endpoint for Kubernetes Dashboard"
  value       = var.kubernetes_dashboard_enabled ? "http://kubernetes-dashboard.${local.namespace}.svc.cluster.local:443" : null
}

output "dashboard_admin_token_command" {
  description = "Command to get the Kubernetes Dashboard admin token"
  value       = var.kubernetes_dashboard_enabled ? "kubectl -n ${local.namespace} create token dashboard-admin" : null
}

output "grafana_admin_password" {
  description = "Admin password for Grafana"
  value       = var.grafana_enabled ? var.grafana_admin_password : null
  sensitive   = true
}

output "monitoring_components" {
  description = "List of enabled monitoring components"
  value = compact([
    var.prometheus_enabled ? "prometheus" : "",
    var.grafana_enabled ? "grafana" : "",
    var.thanos_enabled ? "thanos" : "",
    var.loki_enabled ? "loki" : "",
    var.jaeger_enabled ? "jaeger" : "",
    var.elasticsearch_enabled ? "elasticsearch" : "",
    var.elasticsearch_enabled && var.kibana_enabled ? "kibana" : "",
    var.elasticsearch_enabled && var.filebeat_enabled ? "filebeat" : "",
    var.vector_enabled ? "vector" : "",
    var.opentelemetry_enabled ? "opentelemetry" : "",
    var.kube_state_metrics_enabled ? "kube-state-metrics" : "",
    var.cadvisor_enabled ? "cadvisor" : "",
    var.kubernetes_dashboard_enabled ? "kubernetes-dashboard" : ""
  ])
}

output "ingress_urls" {
  description = "URLs for accessing monitoring components via ingress (if enabled)"
  value = var.ingress_enabled ? {
    grafana = var.grafana_enabled ? "https://grafana.${var.ingress_domain}" : null
    jaeger = var.jaeger_enabled ? "https://jaeger.${var.ingress_domain}" : null
    kibana = var.elasticsearch_enabled && var.kibana_enabled ? "https://kibana.${var.ingress_domain}" : null
    prometheus = var.prometheus_enabled ? "https://prometheus.${var.ingress_domain}" : null
    thanos = var.thanos_enabled ? "https://thanos.${var.ingress_domain}" : null
    kubernetes_dashboard = var.kubernetes_dashboard_enabled ? "https://kubernetes-dashboard.${var.ingress_domain}" : null
  } : null
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled for monitoring components"
  value       = var.high_availability
}

output "prometheus_retention" {
  description = "Retention period for Prometheus metrics"
  value       = var.prometheus_enabled ? var.prometheus_retention : null
}

output "loki_retention" {
  description = "Retention period for Loki logs"
  value       = var.loki_enabled ? var.loki_retention : null
}
