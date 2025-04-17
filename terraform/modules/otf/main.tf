variable "namespace" {
  description = "Kubernetes namespace for OTF deployment"
  type        = string
  default     = "otf"
}

variable "api_key" {
  description = "API key for OTF"
  type        = string
  sensitive   = true
}

resource "kubernetes_namespace" "otf" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"      = "otf"
      "app.kubernetes.io/part-of"   = "agent-runtime"
    }
  }
}

resource "kubernetes_secret" "otf_secrets" {
  metadata {
    name      = "otf-secrets"
    namespace = kubernetes_namespace.otf.metadata[0].name
  }

  data = {
    "api-key" = var.api_key
  }
}

resource "kubernetes_deployment" "otf" {
  metadata {
    name      = "otf"
    namespace = kubernetes_namespace.otf.metadata[0].name
    labels = {
      app                       = "otf"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "otf"
      }
    }

    template {
      metadata {
        labels = {
          app = "otf"
        }
      }

      spec {
        container {
          name  = "otf"
          image = "otfninja/otf:latest"

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
            name = "OTF_API_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.otf_secrets.metadata[0].name
                key  = "api-key"
              }
            }
          }
        }

        service_account_name = "otf-sa"
      }
    }
  }
}

resource "kubernetes_service" "otf" {
  metadata {
    name      = "otf"
    namespace = kubernetes_namespace.otf.metadata[0].name
  }

  spec {
    selector = {
      app = "otf"
    }

    port {
      port        = 80
      target_port = 8080
    }
  }
}

resource "kubernetes_service_account" "otf_sa" {
  metadata {
    name      = "otf-sa"
    namespace = kubernetes_namespace.otf.metadata[0].name
  }
}
