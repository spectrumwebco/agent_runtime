resource "kubernetes_namespace" "jspolicy_system" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_deployment" "jspolicy_controller" {
  metadata {
    name      = "jspolicy-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/component"  = "controller"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "jspolicy"
        "app.kubernetes.io/component" = "controller"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "jspolicy"
          "app.kubernetes.io/component" = "controller"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "controller"
          image = "ghcr.io/loft-sh/jspolicy-controller:${var.jspolicy_version}"
          
          env {
            name  = "POD_NAMESPACE"
            value_from {
              field_ref {
                field_path = "metadata.namespace"
              }
            }
          }
          
          port {
            container_port = 8443
            name           = "webhook"
          }
          
          volume_mount {
            name       = "cert"
            mount_path = "/tmp/k8s-webhook-server/serving-certs"
            read_only  = true
          }
          
          resources {
            requests = {
              memory = "128Mi"
              cpu    = "100m"
            }
            limits = {
              memory = "256Mi"
              cpu    = "200m"
            }
          }
        }
        
        volume {
          name = "cert"
          secret {
            secret_name = kubernetes_secret.jspolicy_webhook_cert.metadata[0].name
          }
        }
        
        service_account_name = kubernetes_service_account.jspolicy_controller.metadata[0].name
      }
    }
  }
}

resource "kubernetes_service" "jspolicy_webhook" {
  metadata {
    name      = "jspolicy-webhook"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/component"  = "webhook"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "jspolicy"
      "app.kubernetes.io/component" = "controller"
    }
    
    port {
      port        = 443
      target_port = "webhook"
      name        = "webhook"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_service_account" "jspolicy_controller" {
  metadata {
    name      = "jspolicy-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_cluster_role" "jspolicy_controller" {
  metadata {
    name = "jspolicy-controller"
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  rule {
    api_groups = [""]
    resources  = ["configmaps"]
    verbs      = ["get", "list", "watch"]
  }
  
  rule {
    api_groups = ["jspolicy.com"]
    resources  = ["policies", "policies/status"]
    verbs      = ["get", "list", "watch", "update", "patch"]
  }
  
  rule {
    api_groups = ["admissionregistration.k8s.io"]
    resources  = ["validatingwebhookconfigurations"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_cluster_role_binding" "jspolicy_controller" {
  metadata {
    name = "jspolicy-controller"
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.jspolicy_controller.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.jspolicy_controller.metadata[0].name
    namespace = kubernetes_service_account.jspolicy_controller.metadata[0].namespace
  }
}

resource "kubernetes_secret" "jspolicy_webhook_cert" {
  metadata {
    name      = "jspolicy-webhook-cert"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  type = "kubernetes.io/tls"
  
  data = {
    "tls.crt" = var.webhook_cert
    "tls.key" = var.webhook_key
  }
}

resource "kubernetes_config_map" "jspolicy_config" {
  metadata {
    name      = "jspolicy-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      controller:
        logLevel: info
        metricsAddr: ":8080"
        webhookPort: 8443
      policies:
        defaultTimeout: 5s
    EOT
  }
}
