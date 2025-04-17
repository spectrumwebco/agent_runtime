variable "namespace" {
  description = "Kubernetes namespace for Lynx deployment"
  type        = string
  default     = "lynx"
}

resource "kubernetes_namespace" "lynx" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"      = "lynx"
      "app.kubernetes.io/part-of"   = "agent-runtime"
    }
  }
}

resource "kubernetes_deployment" "lynx" {
  metadata {
    name      = "lynx"
    namespace = kubernetes_namespace.lynx.metadata[0].name
    labels = {
      app                       = "lynx"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "lynx"
      }
    }

    template {
      metadata {
        labels = {
          app = "lynx"
        }
      }

      spec {
        container {
          name  = "lynx"
          image = "clivern/lynx:latest"

          port {
            container_port = 8080
          }

          resources {
            limits = {
              cpu    = "1"
              memory = "1Gi"
            }
            requests = {
              cpu    = "0.5"
              memory = "500Mi"
            }
          }

          env {
            name  = "LYNX_CONFIG"
            value = "/etc/lynx/config.yml"
          }

          volume_mount {
            name       = "config"
            mount_path = "/etc/lynx"
          }
        }

        volume {
          name = "config"
          config_map {
            name = "lynx-config"
          }
        }

        service_account_name = "lynx-sa"
      }
    }
  }
}

resource "kubernetes_service" "lynx" {
  metadata {
    name      = "lynx"
    namespace = kubernetes_namespace.lynx.metadata[0].name
  }

  spec {
    selector = {
      app = "lynx"
    }

    port {
      port        = 80
      target_port = 8080
    }
  }
}

resource "kubernetes_service_account" "lynx_sa" {
  metadata {
    name      = "lynx-sa"
    namespace = kubernetes_namespace.lynx.metadata[0].name
  }
}

resource "kubernetes_config_map" "lynx_config" {
  metadata {
    name      = "lynx-config"
    namespace = kubernetes_namespace.lynx.metadata[0].name
  }

  data = {
    "config.yml" = <<-EOT
      app:
        name: lynx
        port: 8080
    EOT
  }
}
