
resource "kubernetes_deployment" "django_backend" {
  metadata {
    name      = "django-backend"
    namespace = var.namespace
    labels = {
      app       = "django-backend"
      component = "backend"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "django-backend"
      }
    }

    template {
      metadata {
        labels = {
          app       = "django-backend"
          component = "backend"
        }
      }

      spec {
        container {
          name              = "django"
          image             = var.image
          image_pull_policy = "Always"

          port {
            container_port = 8000
            name           = "http"
          }

          env {
            name = "DATABASE_URL"
            value_from {
              secret_key_ref {
                name = "django-db-credentials"
                key  = "database_url"
              }
            }
          }

          env {
            name = "SECRET_KEY"
            value_from {
              secret_key_ref {
                name = "django-secrets"
                key  = "secret_key"
              }
            }
          }

          env {
            name  = "ALLOWED_HOSTS"
            value = "*"
          }

          env {
            name  = "DEBUG"
            value = "False"
          }

          env {
            name = "REDIS_URL"
            value_from {
              secret_key_ref {
                name = "django-redis-credentials"
                key  = "redis_url"
              }
            }
          }

          env {
            name = "GITHUB_CLIENT_ID"
            value_from {
              secret_key_ref {
                name = "oauth-credentials"
                key  = "github_client_id"
              }
            }
          }

          env {
            name = "GITHUB_CLIENT_SECRET"
            value_from {
              secret_key_ref {
                name = "oauth-credentials"
                key  = "github_client_secret"
              }
            }
          }

          env {
            name = "GITEE_CLIENT_ID"
            value_from {
              secret_key_ref {
                name = "oauth-credentials"
                key  = "gitee_client_id"
              }
            }
          }

          env {
            name = "GITEE_CLIENT_SECRET"
            value_from {
              secret_key_ref {
                name = "oauth-credentials"
                key  = "gitee_client_secret"
              }
            }
          }

          env {
            name = "POLAR_API_KEY"
            value_from {
              secret_key_ref {
                name = "billing-credentials"
                key  = "polar_api_key"
              }
            }
          }

          env {
            name  = "DEVIN_API_URL"
            value = "http://185.196.220.224:8000"
          }

          env {
            name = "DEVIN_API_KEY"
            value_from {
              secret_key_ref {
                name = "api-credentials"
                key  = "devin_api_key"
              }
            }
          }

          resources {
            limits = {
              cpu    = "1"
              memory = "2Gi"
            }
            requests = {
              cpu    = "500m"
              memory = "1Gi"
            }
          }

          liveness_probe {
            http_get {
              path = "/api/health/"
              port = "http"
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }

          readiness_probe {
            http_get {
              path = "/api/health/"
              port = "http"
            }
            initial_delay_seconds = 15
            period_seconds        = 5
          }

          volume_mount {
            name       = "django-config"
            mount_path = "/app/config"
          }

          volume_mount {
            name       = "workspaces"
            mount_path = "/app/workspaces"
          }

          volume_mount {
            name       = "vault-token"
            mount_path = "/vault/token"
            read_only  = true
          }
        }

        volume {
          name = "django-config"
          config_map {
            name = "django-config"
          }
        }

        volume {
          name = "workspaces"
          persistent_volume_claim {
            claim_name = "workspaces-pvc"
          }
        }

        volume {
          name = "vault-token"
          secret {
            secret_name = "vault-token"
          }
        }

        security_context {
          run_as_user  = 1000
          run_as_group = 1000
          fs_group     = 1000
        }
      }
    }
  }
}

resource "kubernetes_service" "django_backend" {
  metadata {
    name      = "django-backend"
    namespace = var.namespace
    labels = {
      app       = "django-backend"
      component = "backend"
    }
  }

  spec {
    selector = {
      app = "django-backend"
    }

    port {
      port        = 8000
      target_port = "http"
      name        = "http"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_config_map" "django_config" {
  metadata {
    name      = "django-config"
    namespace = var.namespace
  }

  data = {
    "settings.py" = file("${path.module}/files/settings.py")
    "urls.py"     = file("${path.module}/files/urls.py")
  }
}

resource "kubernetes_persistent_volume_claim" "workspaces_pvc" {
  metadata {
    name      = "workspaces-pvc"
    namespace = var.namespace
  }

  spec {
    access_modes = ["ReadWriteMany"]
    resources {
      requests = {
        storage = "10Gi"
      }
    }
    storage_class_name = "standard"
  }
}

resource "kubernetes_ingress_v1" "django_ingress" {
  metadata {
    name      = "django-ingress"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/ingress.class"                 = "nginx"
      "nginx.ingress.kubernetes.io/ssl-redirect"    = "true"
      "nginx.ingress.kubernetes.io/proxy-body-size" = "50m"
    }
  }

  spec {
    rule {
      host = var.ingress_host

      http {
        path {
          path      = "/"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.django_backend.metadata[0].name
              port {
                name = "http"
              }
            }
          }
        }
      }
    }

    tls {
      hosts       = [var.ingress_host]
      secret_name = "django-tls-secret"
    }
  }
}
