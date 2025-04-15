resource "kubernetes_namespace" "dragonfly" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_deployment" "dragonfly" {
  metadata {
    name      = "dragonfly"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/component"  = "server"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "dragonfly"
        "app.kubernetes.io/component" = "server"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "dragonfly"
          "app.kubernetes.io/component" = "server"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "dragonfly"
          image = "${var.container_registry}/dragonfly:${var.dragonfly_version}"
          
          port {
            container_port = 6379
            name           = "redis"
          }
          
          port {
            container_port = 8000
            name           = "http"
          }
          
          env {
            name  = "DRAGONFLY_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.dragonfly_secrets.metadata[0].name
                key  = "password"
              }
            }
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/data"
          }
          
          resources {
            requests = {
              memory = "512Mi"
              cpu    = "250m"
            }
            limits = {
              memory = "1Gi"
              cpu    = "500m"
            }
          }
        }
        
        volume {
          name = "data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.dragonfly_data.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "dragonfly_redis" {
  metadata {
    name      = "dragonfly-redis"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/component"  = "redis"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "dragonfly"
      "app.kubernetes.io/component" = "server"
    }
    
    port {
      port        = 6379
      target_port = "redis"
      name        = "redis"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_service" "dragonfly_http" {
  metadata {
    name      = "dragonfly-http"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/component"  = "http"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "dragonfly"
      "app.kubernetes.io/component" = "server"
    }
    
    port {
      port        = 8000
      target_port = "http"
      name        = "http"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_persistent_volume_claim" "dragonfly_data" {
  metadata {
    name      = "dragonfly-data"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "10Gi"
      }
    }
  }
}

resource "kubernetes_secret" "dragonfly_secrets" {
  metadata {
    name      = "dragonfly-secrets"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "password" = var.password
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "dragonfly_config" {
  metadata {
    name      = "dragonfly-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "dragonfly.conf" = <<-EOT
      bind 0.0.0.0
      port 6379
      
      # Memory settings
      maxmemory ${var.max_memory}
      maxmemory-policy ${var.memory_policy}
      
      # Persistence settings
      appendonly yes
      appendfilename "appendonly.aof"
      appendfsync everysec
      
      # Event stream settings
      stream-node-max-bytes 4096
      stream-node-max-entries 100
      
      # Performance settings
      io-threads ${var.io_threads}
      io-threads-do-reads yes
    EOT
  }
}
