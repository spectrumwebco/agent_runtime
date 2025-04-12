
resource "kubernetes_deployment" "vnode_runtime" {
  metadata {
    name      = var.name
    namespace = var.namespace
    labels = {
      app = "vnode-runtime"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "vnode-runtime"
      }
    }

    template {
      metadata {
        labels = {
          app = "vnode-runtime"
        }
        annotations = {
          "prometheus.io/scrape" = "true"
          "prometheus.io/port"   = "8080"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.vnode_runtime.metadata[0].name
        
        container {
          image = "ghcr.io/loft-sh/vnode-runtime:0.0.1-alpha.1"
          name  = "vnode-runtime"

          env {
            name  = "VCLUSTER_NAME"
            value = var.vcluster_name
          }

          env {
            name  = "VCLUSTER_NAMESPACE"
            value = var.vcluster_namespace
          }

          env {
            name  = "LOG_LEVEL"
            value = "info"
          }

          env {
            name = "NODE_NAME"
            value_from {
              field_ref {
                field_path = "spec.nodeName"
              }
            }
          }

          resources {
            limits = {
              cpu    = var.resources.limits.cpu
              memory = var.resources.limits.memory
            }
            requests = {
              cpu    = var.resources.requests.cpu
              memory = var.resources.requests.memory
            }
          }

          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080
            }
            initial_delay_seconds = 30
            period_seconds        = 10
            timeout_seconds       = 5
            failure_threshold     = 3
          }

          readiness_probe {
            http_get {
              path = "/readyz"
              port = 8080
            }
            initial_delay_seconds = 10
            period_seconds        = 10
            timeout_seconds       = 5
            failure_threshold     = 3
          }

          volume_mount {
            name       = "vnode-data"
            mount_path = "/var/lib/vnode"
          }
        }

        volume {
          name = "vnode-data"
          empty_dir {}
        }
      }
    }
  }
}

resource "kubernetes_service" "vnode_runtime" {
  metadata {
    name      = var.name
    namespace = var.namespace
    labels = {
      app = "vnode-runtime"
    }
  }

  spec {
    selector = {
      app = "vnode-runtime"
    }

    port {
      name        = "http"
      port        = 8080
      target_port = 8080
    }

    port {
      name        = "metrics"
      port        = 9090
      target_port = 9090
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_service_account" "vnode_runtime" {
  metadata {
    name      = "${var.name}-sa"
    namespace = var.namespace
  }
}

resource "kubernetes_cluster_role" "vnode_runtime" {
  metadata {
    name = "${var.name}-role"
  }

  rule {
    api_groups = [""]
    resources  = ["nodes"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = [""]
    resources  = ["pods"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = [""]
    resources  = ["events"]
    verbs      = ["create", "patch", "update"]
  }

  rule {
    api_groups = ["node.k8s.io"]
    resources  = ["runtimeclasses"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = [""]
    resources  = ["configmaps"]
    verbs      = ["get", "list", "watch", "create", "update", "patch"]
  }
}

resource "kubernetes_cluster_role_binding" "vnode_runtime" {
  metadata {
    name = "${var.name}-rolebinding"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.vnode_runtime.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.vnode_runtime.metadata[0].name
    namespace = kubernetes_service_account.vnode_runtime.metadata[0].namespace
  }
}

resource "kubernetes_config_map" "vnode_integration" {
  metadata {
    name      = "vnode-integration-config"
    namespace = var.namespace
  }

  data = {
    "config.yaml" = <<-EOT
      vcluster:
        name: ${var.vcluster_name}
        namespace: ${var.vcluster_namespace}
      
      vnode:
        enabled: true
        image: ghcr.io/loft-sh/vnode-runtime:0.0.1-alpha.1
        resources:
          limits:
            cpu: ${var.resources.limits.cpu}
            memory: ${var.resources.limits.memory}
          requests:
            cpu: ${var.resources.requests.cpu}
            memory: ${var.resources.requests.memory}
      
      integration:
        syncNodes: true
        syncPods: true
        syncEvents: true
    EOT
  }
}

resource "kubernetes_runtime_class" "kata_containers" {
  count = var.enable_kata_integration ? 1 : 0
  
  metadata {
    name = "kata-containers"
  }

  handler = "kata"
}

resource "kubernetes_runtime_class" "vnode" {
  metadata {
    name = "vnode"
  }

  handler = "vnode"
  
  scheduling {
    node_selector = {
      "vnode.loft.sh/enabled" = "true"
    }
    
    tolerance {
      key      = "vnode.loft.sh/enabled"
      operator = "Equal"
      value    = "true"
      effect   = "NoSchedule"
    }
  }
}

resource "kubectl_manifest" "node_labels" {
  count = length(var.node_names)
  
  yaml_body = <<-YAML
    apiVersion: v1
    kind: Node
    metadata:
      name: ${var.node_names[count.index]}
      labels:
        vnode.loft.sh/enabled: "true"
    YAML

  override_namespace = var.namespace
  force_new          = true
  server_side_apply  = true
  provider           = kubectl.gavinbunney
}
