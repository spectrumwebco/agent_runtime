resource "kubernetes_namespace" "yap2db" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "yap2db"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_deployment" "yap2db" {
  metadata {
    name      = "yap2db"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "yap2db"
      "app.kubernetes.io/component"  = "database-management"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "yap2db"
        "app.kubernetes.io/component" = "database-management"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "yap2db"
          "app.kubernetes.io/component" = "database-management"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "yap2db"
          image = "${var.container_registry}/yap2db:${var.yap2db_version}"
          
          port {
            container_port = 10824
            name           = "http"
          }
          
          env {
            name  = "SUPABASE_URL"
            value = "http://${var.supabase_host}:${var.supabase_port}"
          }
          
          env {
            name  = "DRAGONFLY_HOST"
            value = var.dragonfly_host
          }
          
          env {
            name  = "DRAGONFLY_PORT"
            value = var.dragonfly_port
          }
          
          env {
            name  = "ROCKETMQ_NAMESRV"
            value = "${var.rocketmq_host}:${var.rocketmq_port}"
          }
          
          env {
            name  = "RAGFLOW_ENDPOINT"
            value = "http://${var.ragflow_host}:${var.ragflow_port}"
          }
          
          env {
            name  = "JAVA_OPTS"
            value = "-Xms512m -Xmx1g"
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/app/data"
          }
          
          resources {
            requests = {
              memory = "1Gi"
              cpu    = "500m"
            }
            limits = {
              memory = "2Gi"
              cpu    = "1000m"
            }
          }
          
          readiness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            initial_delay_seconds = 30
            period_seconds        = 15
          }
          
          liveness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            initial_delay_seconds = 60
            period_seconds        = 30
          }
        }
        
        volume {
          name = "data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.yap2db_data.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "yap2db" {
  metadata {
    name      = "yap2db"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "yap2db"
      "app.kubernetes.io/component"  = "database-management"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "yap2db"
      "app.kubernetes.io/component" = "database-management"
    }
    
    port {
      port        = 10824
      target_port = "http"
      name        = "http"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_persistent_volume_claim" "yap2db_data" {
  metadata {
    name      = "yap2db-data"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "yap2db"
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

resource "kubernetes_config_map" "yap2db_config" {
  metadata {
    name      = "yap2db-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "yap2db"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "yap2db.conf" = <<-EOT
      supabase.url=${var.supabase_host}
      supabase.port=${var.supabase_port}
      
      dragonfly.host=${var.dragonfly_host}
      dragonfly.port=${var.dragonfly_port}
      
      rocketmq.namesrv=${var.rocketmq_host}:${var.rocketmq_port}
      
      ragflow.endpoint=http://${var.ragflow_host}:${var.ragflow_port}
      
      server.port=10824
      spring.profiles.active=release
    EOT
  }
}
