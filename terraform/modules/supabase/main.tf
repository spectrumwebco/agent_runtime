resource "kubernetes_namespace" "supabase" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_deployment" "supabase" {
  metadata {
    name      = "supabase"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/component"  = "server"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "supabase"
        "app.kubernetes.io/component" = "server"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "supabase"
          "app.kubernetes.io/component" = "server"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "supabase"
          image = "${var.container_registry}/supabase:${var.supabase_version}"
          
          port {
            container_port = 8000
            name           = "http"
          }
          
          port {
            container_port = 5432
            name           = "postgres"
          }
          
          env {
            name  = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_secrets.metadata[0].name
                key  = "postgres_password"
              }
            }
          }
          
          env {
            name  = "JWT_SECRET"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_secrets.metadata[0].name
                key  = "jwt_secret"
              }
            }
          }
          
          env {
            name  = "ANON_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_secrets.metadata[0].name
                key  = "anon_key"
              }
            }
          }
          
          env {
            name  = "SERVICE_ROLE_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_secrets.metadata[0].name
                key  = "service_role_key"
              }
            }
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/var/lib/postgresql/data"
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
        }
        
        volume {
          name = "data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.supabase_data.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "supabase_http" {
  metadata {
    name      = "supabase-http"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/component"  = "http"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "supabase"
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

resource "kubernetes_service" "supabase_postgres" {
  metadata {
    name      = "supabase-postgres"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/component"  = "postgres"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "supabase"
      "app.kubernetes.io/component" = "server"
    }
    
    port {
      port        = 5432
      target_port = "postgres"
      name        = "postgres"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_persistent_volume_claim" "supabase_data" {
  metadata {
    name      = "supabase-data"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "20Gi"
      }
    }
  }
}

resource "kubernetes_secret" "supabase_secrets" {
  metadata {
    name      = "supabase-secrets"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "postgres_password" = var.postgres_password
    "jwt_secret"        = var.jwt_secret
    "anon_key"          = var.anon_key
    "service_role_key"  = var.service_role_key
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "supabase_config" {
  metadata {
    name      = "supabase-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.json" = <<-EOT
      {
        "project_ref": "agent-runtime",
        "db_host": "localhost",
        "db_port": 5432,
        "db_name": "postgres",
        "db_user": "postgres",
        "api_external_url": "http://supabase-http.${var.namespace}.svc.cluster.local:8000",
        "studio_port": 8000,
        "kong_port": 8000,
        "kong_url": "http://supabase-http.${var.namespace}.svc.cluster.local:8000",
        "auth_site_url": "http://supabase-http.${var.namespace}.svc.cluster.local:8000",
        "storage_backend": "file",
        "storage_file_backend_path": "/var/lib/storage"
      }
    EOT
  }
}
