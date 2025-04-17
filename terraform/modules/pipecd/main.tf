
resource "kubernetes_namespace" "pipecd" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name" = "pipecd"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }
}

resource "kubernetes_service_account" "pipecd" {
  metadata {
    name      = "pipecd"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "pipecd" {
  metadata {
    name = "pipecd-role"
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["*"]
  }
}

resource "kubernetes_cluster_role_binding" "pipecd" {
  metadata {
    name = "pipecd-role-binding"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.pipecd.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.pipecd.metadata[0].name
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }
}

resource "kubernetes_config_map" "pipecd_config" {
  metadata {
    name      = "pipecd-config"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  data = {
    "control-plane-config.yaml" = templatefile("${path.module}/templates/control-plane-config.yaml.tpl", {
      state_key = var.state_key
      minio_endpoint = var.minio_endpoint
      minio_bucket = var.minio_bucket
      repositories = var.repositories
    })
  }
}

resource "kubernetes_config_map" "pipecd_piped_config" {
  metadata {
    name      = "pipecd-piped-config"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  data = {
    "piped-config.yaml" = templatefile("${path.module}/templates/piped-config.yaml.tpl", {
      project_id = var.project_id
      piped_id = var.piped_id
      repositories = var.repositories
      kubernetes_config = var.kubernetes_config
      terraform_config = var.terraform_config
    })
  }
}

resource "kubernetes_secret" "pipecd_secret" {
  metadata {
    name      = "pipecd-secret"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  data = {
    "minio-access-key" = var.minio_access_key
    "minio-secret-key" = var.minio_secret_key
    "ssh-key"          = var.ssh_key
  }
}

resource "kubernetes_secret" "pipecd_piped_secret" {
  metadata {
    name      = "pipecd-piped-secret"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  data = {
    "piped-key"  = var.piped_key
    "ssh-key"    = var.ssh_key
    "kubeconfig" = var.kubeconfig
  }
}

resource "kubernetes_deployment" "pipecd_control_plane" {
  metadata {
    name      = "pipecd-control-plane"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  spec {
    replicas = var.control_plane_replicas

    selector {
      match_labels = {
        app = "pipecd-control-plane"
      }
    }

    template {
      metadata {
        labels = {
          app = "pipecd-control-plane"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.pipecd.metadata[0].name

        container {
          name  = "control-plane"
          image = "ghcr.io/pipe-cd/pipecd:${var.pipecd_version}"
          args  = ["server", "--config-file=/etc/pipecd-config/control-plane-config.yaml"]

          port {
            container_port = 9082
            name           = "http"
          }

          port {
            container_port = 9083
            name           = "grpc"
          }

          volume_mount {
            name       = "config-volume"
            mount_path = "/etc/pipecd-config"
          }

          volume_mount {
            name       = "secret-volume"
            mount_path = "/etc/pipecd-secret"
          }
        }

        volume {
          name = "config-volume"
          config_map {
            name = kubernetes_config_map.pipecd_config.metadata[0].name
          }
        }

        volume {
          name = "secret-volume"
          secret {
            secret_name = kubernetes_secret.pipecd_secret.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "pipecd_control_plane" {
  metadata {
    name      = "pipecd-control-plane"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  spec {
    selector = {
      app = "pipecd-control-plane"
    }

    port {
      port        = 80
      target_port = 9082
      name        = "http"
    }

    port {
      port        = 9083
      target_port = 9083
      name        = "grpc"
    }
  }
}

resource "kubernetes_deployment" "pipecd_piped" {
  metadata {
    name      = "pipecd-piped"
    namespace = kubernetes_namespace.pipecd.metadata[0].name
  }

  spec {
    replicas = var.piped_replicas

    selector {
      match_labels = {
        app = "pipecd-piped"
      }
    }

    template {
      metadata {
        labels = {
          app = "pipecd-piped"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.pipecd.metadata[0].name

        container {
          name  = "piped"
          image = "ghcr.io/pipe-cd/piped:${var.pipecd_version}"
          args  = ["piped", "--config-file=/etc/pipecd-config/piped-config.yaml"]

          volume_mount {
            name       = "config-volume"
            mount_path = "/etc/pipecd-config"
          }

          volume_mount {
            name       = "secret-volume"
            mount_path = "/etc/pipecd-secret"
          }
        }

        volume {
          name = "config-volume"
          config_map {
            name = kubernetes_config_map.pipecd_piped_config.metadata[0].name
          }
        }

        volume {
          name = "secret-volume"
          secret {
            secret_name = kubernetes_secret.pipecd_piped_secret.metadata[0].name
          }
        }
      }
    }
  }
}
