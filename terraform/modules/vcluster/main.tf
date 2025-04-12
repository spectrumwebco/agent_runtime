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

resource "kubernetes_config_map" "vcluster_config" {
  metadata {
    name      = "vcluster-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      vcluster:
        image: rancher/k3s:v1.21.4-k3s1
        imagePullPolicy: Always
      networking:
        hostPort: 8443
        containerPort: 8443
    EOT
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
        service_account_name = kubernetes_service_account.vcluster_controller.metadata[0].name
        
        container {
          name            = "controller"
          image           = "rancher/k3s:v1.21.4-k3s1"
          image_pull_policy = "Always"
          
          args = [
            "--config=/etc/vcluster/config.yaml"
          ]
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/vcluster"
          }
          
          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.vcluster_config.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service_account" "vcluster_controller" {
  metadata {
    name      = "vcluster-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }
}

resource "kubernetes_cluster_role" "vcluster_controller" {
  metadata {
    name = "vcluster-controller"
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  rule {
    api_groups = [""]
    resources  = ["namespaces", "pods", "services", "configmaps", "secrets"]
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
    namespace = var.namespace
  }
}

resource "kubernetes_service" "vcluster_webhook" {
  metadata {
    name      = "vcluster-webhook"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "vcluster"
      "app.kubernetes.io/component"  = "webhook"
      "app.kubernetes.io/part-of"    = "agent-runtime"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "vcluster"
      "app.kubernetes.io/component" = "controller"
    }
    
    port {
      port        = 443
      target_port = 8443
      name        = "webhook"
    }
    
    type = "ClusterIP"
  }
}
