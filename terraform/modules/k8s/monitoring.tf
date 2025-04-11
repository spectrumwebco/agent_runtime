
locals {
  monitoring_namespace = "monitoring"
}

resource "kubernetes_namespace" "monitoring" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name = local.monitoring_namespace
    
    labels = {
      "app.kubernetes.io/name" = "monitoring"
      "app.kubernetes.io/part-of" = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "helm_release" "prometheus" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "prometheus"
  version    = "19.6.1"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/prometheus/prometheus.yaml")
  ]
  
  set {
    name  = "server.persistentVolume.enabled"
    value = "true"
  }
  
  set {
    name  = "server.persistentVolume.size"
    value = "50Gi"
  }
  
  set {
    name  = "server.retention"
    value = "15d"
  }
  
  set {
    name  = "alertmanager.persistentVolume.enabled"
    value = "true"
  }
  
  set {
    name  = "alertmanager.persistentVolume.size"
    value = "10Gi"
  }
  
  set {
    name  = "server.resources.limits.cpu"
    value = "1000m"
  }
  
  set {
    name  = "server.resources.limits.memory"
    value = "2Gi"
  }
  
  set {
    name  = "server.resources.requests.cpu"
    value = "500m"
  }
  
  set {
    name  = "server.resources.requests.memory"
    value = "1Gi"
  }
}

resource "helm_release" "thanos" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "thanos"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "thanos"
  version    = "12.5.1"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/thanos/thanos.yaml")
  ]
  
  set {
    name  = "query.replicaCount"
    value = "2"
  }
  
  set {
    name  = "storegateway.replicaCount"
    value = "2"
  }
  
  set {
    name  = "compactor.retentionResolutionRaw"
    value = "30d"
  }
  
  set {
    name  = "compactor.retentionResolution5m"
    value = "90d"
  }
  
  set {
    name  = "compactor.retentionResolution1h"
    value = "1y"
  }
  
  depends_on = [
    helm_release.prometheus
  ]
}

resource "helm_release" "grafana" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  version    = "6.56.5"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/grafana/grafana.yaml")
  ]
  
  set {
    name  = "persistence.enabled"
    value = "true"
  }
  
  set {
    name  = "persistence.size"
    value = "10Gi"
  }
  
  set {
    name  = "adminPassword"
    value = "admin"  # In production, use a secret
  }
  
  set {
    name  = "datasources.datasources\\.yaml.apiVersion"
    value = "1"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[0].name"
    value = "Prometheus"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[0].type"
    value = "prometheus"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[0].url"
    value = "http://prometheus-server.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[0].access"
    value = "proxy"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[0].isDefault"
    value = "true"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[1].name"
    value = "Thanos"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[1].type"
    value = "prometheus"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[1].url"
    value = "http://thanos-query.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:9090"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[1].access"
    value = "proxy"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[2].name"
    value = "Loki"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[2].type"
    value = "loki"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[2].url"
    value = "http://loki-gateway.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80"
  }
  
  set {
    name  = "datasources.datasources\\.yaml.datasources[2].access"
    value = "proxy"
  }
  
  depends_on = [
    helm_release.prometheus,
    helm_release.thanos,
    helm_release.loki
  ]
}

resource "helm_release" "loki" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki-stack"
  version    = "2.9.10"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/loki/loki.yaml")
  ]
  
  set {
    name  = "loki.persistence.enabled"
    value = "true"
  }
  
  set {
    name  = "loki.persistence.size"
    value = "50Gi"
  }
  
  set {
    name  = "loki.auth_enabled"
    value = "false"
  }
  
  set {
    name  = "promtail.enabled"
    value = "true"
  }
  
  set {
    name  = "promtail.config.lokiAddress"
    value = "http://loki-gateway.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80/loki/api/v1/push"
  }
}

resource "helm_release" "jaeger" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "jaeger"
  repository = "https://jaegertracing.github.io/helm-charts"
  chart      = "jaeger"
  version    = "0.71.5"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/jaeger/jaeger.yaml")
  ]
  
  set {
    name  = "persistence.enabled"
    value = "true"
  }
  
  set {
    name  = "persistence.size"
    value = "10Gi"
  }
  
  set {
    name  = "collector.replicaCount"
    value = "2"
  }
  
  set {
    name  = "query.replicaCount"
    value = "2"
  }
}

resource "helm_release" "vector" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "vector"
  repository = "https://helm.vector.dev"
  chart      = "vector"
  version    = "0.20.0"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/vector/vector.yaml")
  ]
  
  set {
    name  = "role"
    value = "Agent"
  }
  
  set {
    name  = "customConfig.data_dir"
    value = "/vector-data-dir"
  }
  
  set {
    name  = "customConfig.sinks.loki.type"
    value = "loki"
  }
  
  set {
    name  = "customConfig.sinks.loki.endpoint"
    value = "http://loki-gateway.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80"
  }
  
  set {
    name  = "customConfig.sinks.loki.encoding.codec"
    value = "json"
  }
  
  depends_on = [
    helm_release.loki
  ]
}

resource "helm_release" "opentelemetry" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "opentelemetry"
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-collector"
  version    = "0.55.0"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/opentelemetry/opentelemetry.yaml")
  ]
  
  set {
    name  = "mode"
    value = "daemonset"
  }
  
  set {
    name  = "config.exporters.prometheus.endpoint"
    value = "0.0.0.0:8889"
  }
  
  set {
    name  = "config.exporters.jaeger.endpoint"
    value = "jaeger-collector.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:14250"
  }
  
  set {
    name  = "config.exporters.loki.endpoint"
    value = "http://loki-gateway.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80/loki/api/v1/push"
  }
  
  depends_on = [
    helm_release.jaeger,
    helm_release.loki
  ]
}

resource "helm_release" "kube_state_metrics" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "kube-state-metrics"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-state-metrics"
  version    = "5.6.2"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/monitoring/kube-state-metrics/kube-state-metrics.yaml")
  ]
  
  set {
    name  = "replicas"
    value = "2"
  }
  
  set {
    name  = "collectors.certificatesigningrequests"
    value = "true"
  }
  
  set {
    name  = "collectors.configmaps"
    value = "true"
  }
  
  set {
    name  = "collectors.cronjobs"
    value = "true"
  }
  
  set {
    name  = "collectors.daemonsets"
    value = "true"
  }
  
  set {
    name  = "collectors.deployments"
    value = "true"
  }
  
  set {
    name  = "collectors.endpoints"
    value = "true"
  }
  
  set {
    name  = "collectors.horizontalpodautoscalers"
    value = "true"
  }
  
  set {
    name  = "collectors.ingresses"
    value = "true"
  }
  
  set {
    name  = "collectors.jobs"
    value = "true"
  }
  
  set {
    name  = "collectors.limitranges"
    value = "true"
  }
  
  set {
    name  = "collectors.namespaces"
    value = "true"
  }
  
  set {
    name  = "collectors.nodes"
    value = "true"
  }
  
  set {
    name  = "collectors.persistentvolumeclaims"
    value = "true"
  }
  
  set {
    name  = "collectors.persistentvolumes"
    value = "true"
  }
  
  set {
    name  = "collectors.poddisruptionbudgets"
    value = "true"
  }
  
  set {
    name  = "collectors.pods"
    value = "true"
  }
  
  set {
    name  = "collectors.replicasets"
    value = "true"
  }
  
  set {
    name  = "collectors.replicationcontrollers"
    value = "true"
  }
  
  set {
    name  = "collectors.resourcequotas"
    value = "true"
  }
  
  set {
    name  = "collectors.secrets"
    value = "true"
  }
  
  set {
    name  = "collectors.services"
    value = "true"
  }
  
  set {
    name  = "collectors.statefulsets"
    value = "true"
  }
  
  set {
    name  = "collectors.storageclasses"
    value = "true"
  }
  
  set {
    name  = "collectors.verticalpodautoscalers"
    value = "true"
  }
}

resource "kubernetes_daemon_set" "cadvisor" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name      = "cadvisor"
    namespace = kubernetes_namespace.monitoring[0].metadata[0].name
    
    labels = {
      app = "cadvisor"
    }
  }
  
  spec {
    selector {
      match_labels = {
        app = "cadvisor"
      }
    }
    
    template {
      metadata {
        labels = {
          app = "cadvisor"
        }
        
        annotations = {
          "prometheus.io/scrape" = "true"
          "prometheus.io/port"   = "8080"
        }
      }
      
      spec {
        service_account_name = kubernetes_service_account.cadvisor[0].metadata[0].name
        
        container {
          name  = "cadvisor"
          image = "gcr.io/cadvisor/cadvisor:v0.47.2"
          
          args = [
            "--housekeeping_interval=10s",
            "--max_housekeeping_interval=15s",
            "--event_storage_event_limit=default=0",
            "--event_storage_age_limit=default=0",
            "--disable_metrics=disk,diskIO,network,tcp,udp,percpu,sched,process",
            "--docker_only=true",
            "--store_container_labels=false",
            "--whitelisted_container_labels=io.kubernetes.pod.name,io.kubernetes.pod.namespace,io.kubernetes.container.name"
          ]
          
          resources {
            limits = {
              cpu    = "400m"
              memory = "400Mi"
            }
            
            requests = {
              cpu    = "150m"
              memory = "200Mi"
            }
          }
          
          volume_mount {
            name       = "rootfs"
            mount_path = "/rootfs"
            read_only  = true
          }
          
          volume_mount {
            name       = "var-run"
            mount_path = "/var/run"
            read_only  = true
          }
          
          volume_mount {
            name       = "sys"
            mount_path = "/sys"
            read_only  = true
          }
          
          volume_mount {
            name       = "docker"
            mount_path = "/var/lib/docker"
            read_only  = true
          }
          
          volume_mount {
            name       = "disk"
            mount_path = "/dev/disk"
            read_only  = true
          }
          
          port {
            name           = "http"
            container_port = 8080
            protocol       = "TCP"
          }
          
          security_context {
            privileged = true
          }
        }
        
        termination_grace_period_seconds = 30
        
        volume {
          name = "rootfs"
          
          host_path {
            path = "/"
          }
        }
        
        volume {
          name = "var-run"
          
          host_path {
            path = "/var/run"
          }
        }
        
        volume {
          name = "sys"
          
          host_path {
            path = "/sys"
          }
        }
        
        volume {
          name = "docker"
          
          host_path {
            path = "/var/lib/docker"
          }
        }
        
        volume {
          name = "disk"
          
          host_path {
            path = "/dev/disk"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "cadvisor" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name      = "cadvisor"
    namespace = kubernetes_namespace.monitoring[0].metadata[0].name
    
    labels = {
      app = "cadvisor"
    }
  }
  
  spec {
    selector = {
      app = "cadvisor"
    }
    
    port {
      port        = 8080
      target_port = 8080
      name        = "http"
    }
  }
}

resource "kubernetes_service_account" "cadvisor" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name      = "cadvisor"
    namespace = kubernetes_namespace.monitoring[0].metadata[0].name
  }
}

resource "kubernetes_cluster_role" "cadvisor" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name = "cadvisor"
  }
  
  rule {
    api_groups = [""]
    resources  = ["nodes", "nodes/proxy", "nodes/metrics", "services", "endpoints", "pods"]
    verbs      = ["get", "list", "watch"]
  }
  
  rule {
    api_groups = ["extensions", "apps"]
    resources  = ["deployments"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_cluster_role_binding" "cadvisor" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name = "cadvisor"
  }
  
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.cadvisor[0].metadata[0].name
  }
  
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.cadvisor[0].metadata[0].name
    namespace = kubernetes_service_account.cadvisor[0].metadata[0].namespace
  }
}

resource "helm_release" "kubernetes_dashboard" {
  count = var.monitoring_enabled ? 1 : 0
  
  name       = "kubernetes-dashboard"
  repository = "https://kubernetes.github.io/dashboard/"
  chart      = "kubernetes-dashboard"
  version    = "6.0.8"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  
  set {
    name  = "protocolHttp"
    value = "true"
  }
  
  set {
    name  = "service.externalPort"
    value = "80"
  }
  
  set {
    name  = "replicaCount"
    value = "2"
  }
  
  set {
    name  = "metricsScraper.enabled"
    value = "true"
  }
  
  set {
    name  = "metrics-server.enabled"
    value = "true"
  }
}

resource "kubernetes_service_account" "dashboard_admin" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name      = "dashboard-admin"
    namespace = kubernetes_namespace.monitoring[0].metadata[0].name
  }
}

resource "kubernetes_cluster_role_binding" "dashboard_admin" {
  count = var.monitoring_enabled ? 1 : 0
  
  metadata {
    name = "dashboard-admin"
  }
  
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "cluster-admin"
  }
  
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.dashboard_admin[0].metadata[0].name
    namespace = kubernetes_service_account.dashboard_admin[0].metadata[0].namespace
  }
}

output "grafana_url" {
  description = "URL for Grafana dashboard"
  value       = var.monitoring_enabled ? "http://grafana.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80" : null
}

output "prometheus_url" {
  description = "URL for Prometheus dashboard"
  value       = var.monitoring_enabled ? "http://prometheus-server.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80" : null
}

output "jaeger_url" {
  description = "URL for Jaeger UI"
  value       = var.monitoring_enabled ? "http://jaeger-query.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:16686" : null
}

output "kubernetes_dashboard_url" {
  description = "URL for Kubernetes Dashboard"
  value       = var.monitoring_enabled ? "http://kubernetes-dashboard.${kubernetes_namespace.monitoring[0].metadata[0].name}.svc.cluster.local:80" : null
}
