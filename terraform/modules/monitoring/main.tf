
resource "kubernetes_namespace" "monitoring" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    
    labels = merge({
      "app.kubernetes.io/name"       = "monitoring"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
    
    annotations = var.annotations
  }
}

locals {
  namespace = var.create_namespace ? kubernetes_namespace.monitoring[0].metadata[0].name : var.namespace
}

resource "helm_release" "prometheus" {
  count = var.prometheus_enabled ? 1 : 0
  
  name       = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "prometheus"
  version    = var.prometheus_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "grafana" {
  count = var.grafana_enabled ? 1 : 0
  
  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  version    = var.grafana_version
  namespace  = local.namespace
  
  
  depends_on = [
    kubernetes_namespace.monitoring,
    helm_release.prometheus,
    helm_release.thanos,
    helm_release.loki,
    helm_release.jaeger,
    helm_release.elasticsearch
  ]
}

resource "helm_release" "thanos" {
  count = var.thanos_enabled ? 1 : 0
  
  name       = "thanos"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "thanos"
  version    = var.thanos_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "kubernetes_secret" "thanos_objstore_config" {
  count = var.thanos_enabled && var.thanos_objstore_config != "" ? 1 : 0
  
  metadata {
    name      = "thanos-objstore-config"
    namespace = local.namespace
  }
  
  data = {
    "objstore.yml" = var.thanos_objstore_config
  }
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "loki" {
  count = var.loki_enabled ? 1 : 0
  
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki-stack"
  version    = var.loki_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "jaeger" {
  count = var.jaeger_enabled ? 1 : 0
  
  name       = "jaeger"
  repository = "https://jaegertracing.github.io/helm-charts"
  chart      = "jaeger"
  version    = var.jaeger_version
  namespace  = local.namespace
  
  
  depends_on = [
    kubernetes_namespace.monitoring,
    helm_release.elasticsearch
  ]
}

resource "helm_release" "elasticsearch" {
  count = var.elasticsearch_enabled ? 1 : 0
  
  name       = "elasticsearch"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "elasticsearch"
  version    = var.elasticsearch_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "kibana" {
  count = var.elasticsearch_enabled && var.kibana_enabled ? 1 : 0
  
  name       = "kibana"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "kibana"
  version    = var.kibana_version
  namespace  = local.namespace
  
  
  depends_on = [
    kubernetes_namespace.monitoring,
    helm_release.elasticsearch
  ]
}

resource "helm_release" "filebeat" {
  count = var.elasticsearch_enabled && var.filebeat_enabled ? 1 : 0
  
  name       = "filebeat"
  repository = "https://charts.elastic.co/helm"
  chart      = "filebeat"
  version    = var.filebeat_version
  namespace  = local.namespace
  
  
  depends_on = [
    kubernetes_namespace.monitoring,
    helm_release.elasticsearch
  ]
}

resource "helm_release" "vector" {
  count = var.vector_enabled ? 1 : 0
  
  name       = "vector"
  repository = "https://helm.vector.dev"
  chart      = "vector"
  version    = var.vector_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "opentelemetry" {
  count = var.opentelemetry_enabled ? 1 : 0
  
  name       = "opentelemetry-collector"
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-collector"
  version    = var.opentelemetry_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "kube_state_metrics" {
  count = var.kube_state_metrics_enabled ? 1 : 0
  
  name       = "kube-state-metrics"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-state-metrics"
  version    = var.kube_state_metrics_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "kubernetes_daemon_set" "cadvisor" {
  count = var.cadvisor_enabled ? 1 : 0
  
  metadata {
    name      = "cadvisor"
    namespace = local.namespace
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
        container {
          name  = "cadvisor"
          image = "gcr.io/cadvisor/cadvisor:${var.cadvisor_version}"
          
          resources {
            limits = {
              cpu    = "300m"
              memory = "500Mi"
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
        }
        
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
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "helm_release" "kubernetes_dashboard" {
  count = var.kubernetes_dashboard_enabled ? 1 : 0
  
  name       = "kubernetes-dashboard"
  repository = "https://kubernetes.github.io/dashboard/"
  chart      = "kubernetes-dashboard"
  version    = var.kubernetes_dashboard_version
  namespace  = local.namespace
  
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "kubernetes_service_account" "dashboard_admin" {
  count = var.kubernetes_dashboard_enabled ? 1 : 0
  
  metadata {
    name      = "dashboard-admin"
    namespace = local.namespace
  }
  
  depends_on = [kubernetes_namespace.monitoring]
}

resource "kubernetes_cluster_role_binding" "dashboard_admin" {
  count = var.kubernetes_dashboard_enabled ? 1 : 0
  
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
    namespace = local.namespace
  }
  
  depends_on = [kubernetes_service_account.dashboard_admin]
}
