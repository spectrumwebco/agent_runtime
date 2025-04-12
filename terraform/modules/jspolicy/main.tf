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

resource "kubernetes_config_map" "jspolicy_config" {
  metadata {
    name      = "jspolicy-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      controller:
        logLevel: info
        webhookPort: 9443
      validator:
        logLevel: info
        port: 8443
      policies:
        defaultNamespace: ${var.namespace}
    EOT
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
          
          args = [
            "--config=/etc/jspolicy/config.yaml"
          ]
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/jspolicy"
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
          name = "config"
          config_map {
            name = kubernetes_config_map.jspolicy_config.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_deployment" "jspolicy_validator" {
  metadata {
    name      = "jspolicy-validator"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/component"  = "validator"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "jspolicy"
        "app.kubernetes.io/component" = "validator"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "jspolicy"
          "app.kubernetes.io/component" = "validator"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "validator"
          image = "ghcr.io/loft-sh/jspolicy-validator:${var.jspolicy_version}"
          
          args = [
            "--config=/etc/jspolicy/config.yaml"
          ]
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/jspolicy"
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
          name = "config"
          config_map {
            name = kubernetes_config_map.jspolicy_config.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "jspolicy_controller" {
  metadata {
    name      = "jspolicy-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/component"  = "controller"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "jspolicy"
      "app.kubernetes.io/component" = "controller"
    }
    
    port {
      port        = 9443
      target_port = 9443
      name        = "webhook"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_service" "jspolicy_validator" {
  metadata {
    name      = "jspolicy-validator"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "jspolicy"
      "app.kubernetes.io/component"  = "validator"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "jspolicy"
      "app.kubernetes.io/component" = "validator"
    }
    
    port {
      port        = 8443
      target_port = 8443
      name        = "api"
    }
    
    type = "ClusterIP"
  }
}
