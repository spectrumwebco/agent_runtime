resource "kubernetes_namespace" "vcluster_system" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_deployment" "vcluster_controller" {
  metadata {
    name      = "vcluster-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/component"  = "controller"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "vcluster"
        "app.kubernetes.io/component" = "controller"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "vcluster"
          "app.kubernetes.io/component" = "controller"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "controller"
          image = "ghcr.io/loft-sh/vcluster:${var.vcluster_version}"
          
          args = [
            "--name=vcluster",
            "--service-name=vcluster",
            "--kube-config=/etc/kubernetes/admin.conf"
          ]
          
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
            name           = "https"
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/kubernetes"
            read_only  = true
          }
          
          resources {
            requests = {
              memory = "256Mi"
              cpu    = "250m"
            }
            limits = {
              memory = "512Mi"
              cpu    = "500m"
            }
          }
        }
        
        volume {
          name = "config"
          secret {
            secret_name = kubernetes_secret.vcluster_config.metadata[0].name
          }
        }
        
        service_account_name = kubernetes_service_account.vcluster_controller.metadata[0].name
      }
    }
  }
}

resource "kubernetes_service" "vcluster" {
  metadata {
    name      = "vcluster"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/component"  = "api"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "vcluster"
      "app.kubernetes.io/component" = "controller"
    }
    
    port {
      port        = 443
      target_port = "https"
      name        = "https"
    }
    
    type = "ClusterIP"
  }
}

resource "kubernetes_service_account" "vcluster_controller" {
  metadata {
    name      = "vcluster-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_cluster_role" "vcluster_controller" {
  metadata {
    name = "vcluster-controller"
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  rule {
    api_groups = [""]
    resources  = ["pods", "services", "configmaps", "secrets", "persistentvolumeclaims"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
  
  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "statefulsets", "daemonsets"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
  
  rule {
    api_groups = ["networking.k8s.io"]
    resources  = ["ingresses"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_cluster_role_binding" "vcluster_controller" {
  metadata {
    name = "vcluster-controller"
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.vcluster_controller.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.vcluster_controller.metadata[0].name
    namespace = kubernetes_service_account.vcluster_controller.metadata[0].namespace
  }
}

resource "kubernetes_secret" "vcluster_config" {
  metadata {
    name      = "vcluster-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "admin.conf" = var.kube_config
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "vcluster_config" {
  metadata {
    name      = "vcluster-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      syncer:
        extraArgs: []
      networking:
        resolveDNS: true
      storage:
        persistence: true
    EOT
  }
}
