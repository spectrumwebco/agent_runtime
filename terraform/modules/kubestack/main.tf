variable "namespace" {
  description = "Kubernetes namespace for Kubestack deployment"
  type        = string
  default     = "kubestack"
}

resource "kubernetes_namespace" "kubestack" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"      = "kubestack"
      "app.kubernetes.io/part-of"   = "agent-runtime"
    }
  }
}

resource "kubernetes_deployment" "kubestack" {
  metadata {
    name      = "kubestack"
    namespace = kubernetes_namespace.kubestack.metadata[0].name
    labels = {
      app                       = "kubestack"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "kubestack"
      }
    }

    template {
      metadata {
        labels = {
          app = "kubestack"
        }
      }

      spec {
        container {
          name  = "kubestack"
          image = "kubestack/kubestack:latest"

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
        }

        service_account_name = "kubestack-sa"
      }
    }
  }
}

resource "kubernetes_service" "kubestack" {
  metadata {
    name      = "kubestack"
    namespace = kubernetes_namespace.kubestack.metadata[0].name
  }

  spec {
    selector = {
      app = "kubestack"
    }

    port {
      port        = 80
      target_port = 8080
    }
  }
}

resource "kubernetes_service_account" "kubestack_sa" {
  metadata {
    name      = "kubestack-sa"
    namespace = kubernetes_namespace.kubestack.metadata[0].name
  }
}
