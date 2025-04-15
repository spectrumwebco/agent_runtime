resource "kubernetes_namespace" "vnode_system" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vnode"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_config_map" "vnode_config" {
  metadata {
    name      = "vnode-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vnode"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      runtime:
        image: ghcr.io/loft-sh/vnode-runtime:${var.vnode_version}
        imagePullPolicy: Always
      networking:
        hostPort: 8080
        containerPort: 8080
    EOT
  }
}

resource "kubernetes_daemon_set" "vnode_runtime" {
  metadata {
    name      = "vnode-runtime"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vnode"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector {
      match_labels = {
        "app.kubernetes.io/name" = "vnode-runtime"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "vnode-runtime"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "vnode"
          image = "ghcr.io/loft-sh/vnode-runtime:${var.vnode_version}"
          image_pull_policy = "Always"
          
          security_context {
            privileged = true
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/vnode"
          }
          
          volume_mount {
            name       = "runtime"
            mount_path = "/var/run/vnode"
          }
          
          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.vnode_config.metadata[0].name
          }
        }
        
        volume {
          name = "runtime"
          host_path {
            path = "/var/run/vnode"
            type = "DirectoryOrCreate"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "vnode_runtime" {
  metadata {
    name      = "vnode-runtime"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vnode"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name" = "vnode-runtime"
    }
    
    port {
      port        = 8080
      target_port = 8080
      name        = "http"
    }
    
    type = "ClusterIP"
  }
}
