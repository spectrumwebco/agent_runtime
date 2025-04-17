
resource "kubernetes_namespace" "kubefirst" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name" = "kubefirst"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }
}

resource "kubernetes_service_account" "kubefirst" {
  metadata {
    name      = "kubefirst"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "kubefirst" {
  metadata {
    name = "kubefirst-role"
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["*"]
  }
}

resource "kubernetes_cluster_role_binding" "kubefirst" {
  metadata {
    name = "kubefirst-role-binding"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.kubefirst.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.kubefirst.metadata[0].name
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }
}

resource "kubernetes_config_map" "kubefirst_config" {
  metadata {
    name      = "kubefirst-config"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  data = {
    "config.yaml" = templatefile("${path.module}/templates/config.yaml.tpl", {
      git_provider = var.git_provider
      git_username = var.git_username
      cloud_provider = var.cloud_provider
      cluster_name = var.cluster_name
      gitops_template_url = var.gitops_template_url
      gitops_template_branch = var.gitops_template_branch
    })
  }
}

resource "kubernetes_secret" "kubefirst_secret" {
  metadata {
    name      = "kubefirst-secret"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  data = {
    "gitea-password" = var.git_password
    "vault-token"    = var.vault_token
  }
}

resource "kubernetes_deployment" "kubefirst" {
  metadata {
    name      = "kubefirst"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "kubefirst"
      }
    }

    template {
      metadata {
        labels = {
          app = "kubefirst"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.kubefirst.metadata[0].name

        container {
          name  = "kubefirst"
          image = "kubefirst/kubefirst:${var.kubefirst_version}"
          args  = ["server", "--config-file=/etc/kubefirst/config.yaml"]

          port {
            container_port = 8080
            name           = "http"
          }

          volume_mount {
            name       = "config-volume"
            mount_path = "/etc/kubefirst"
          }

          volume_mount {
            name       = "secret-volume"
            mount_path = "/etc/kubefirst-secret"
          }
        }

        volume {
          name = "config-volume"
          config_map {
            name = kubernetes_config_map.kubefirst_config.metadata[0].name
          }
        }

        volume {
          name = "secret-volume"
          secret {
            secret_name = kubernetes_secret.kubefirst_secret.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "kubefirst" {
  metadata {
    name      = "kubefirst"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    selector = {
      app = "kubefirst"
    }

    port {
      port        = 80
      target_port = 8080
      name        = "http"
    }
  }
}

resource "kubernetes_deployment" "gitea" {
  metadata {
    name      = "gitea"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "gitea"
      }
    }

    template {
      metadata {
        labels = {
          app = "gitea"
        }
      }

      spec {
        container {
          name  = "gitea"
          image = "gitea/gitea:${var.gitea_version}"

          port {
            container_port = 3000
            name           = "http"
          }

          port {
            container_port = 22
            name           = "ssh"
          }

          env {
            name  = "GITEA__database__DB_TYPE"
            value = "postgres"
          }

          env {
            name  = "GITEA__database__HOST"
            value = "gitea-postgres:5432"
          }

          env {
            name  = "GITEA__database__NAME"
            value = "gitea"
          }

          env {
            name  = "GITEA__database__USER"
            value = "gitea"
          }

          env {
            name = "GITEA__database__PASSWD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.gitea_postgres_secret.metadata[0].name
                key  = "password"
              }
            }
          }

          volume_mount {
            name       = "gitea-data"
            mount_path = "/data"
          }
        }

        volume {
          name = "gitea-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.gitea_data_pvc.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "gitea" {
  metadata {
    name      = "gitea"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    selector = {
      app = "gitea"
    }

    port {
      port        = 80
      target_port = 3000
      name        = "http"
    }

    port {
      port        = 22
      target_port = 22
      name        = "ssh"
    }
  }
}

resource "kubernetes_persistent_volume_claim" "gitea_data_pvc" {
  metadata {
    name      = "gitea-data-pvc"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "10Gi"
      }
    }
  }
}

resource "kubernetes_deployment" "gitea_postgres" {
  metadata {
    name      = "gitea-postgres"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "gitea-postgres"
      }
    }

    template {
      metadata {
        labels = {
          app = "gitea-postgres"
        }
      }

      spec {
        container {
          name  = "postgres"
          image = "postgres:15"

          port {
            container_port = 5432
          }

          env {
            name  = "POSTGRES_USER"
            value = "gitea"
          }

          env {
            name  = "POSTGRES_DB"
            value = "gitea"
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.gitea_postgres_secret.metadata[0].name
                key  = "password"
              }
            }
          }

          volume_mount {
            name       = "postgres-data"
            mount_path = "/var/lib/postgresql/data"
          }
        }

        volume {
          name = "postgres-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.gitea_postgres_pvc.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "gitea_postgres" {
  metadata {
    name      = "gitea-postgres"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    selector = {
      app = "gitea-postgres"
    }

    port {
      port        = 5432
      target_port = 5432
    }
  }
}

resource "kubernetes_persistent_volume_claim" "gitea_postgres_pvc" {
  metadata {
    name      = "gitea-postgres-pvc"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "5Gi"
      }
    }
  }
}

resource "kubernetes_secret" "gitea_postgres_secret" {
  metadata {
    name      = "gitea-postgres-secret"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  data = {
    password = var.gitea_postgres_password
  }
}

resource "kubernetes_deployment" "vault" {
  metadata {
    name      = "vault"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "vault"
      }
    }

    template {
      metadata {
        labels = {
          app = "vault"
        }
      }

      spec {
        container {
          name  = "vault"
          image = "hashicorp/vault:${var.vault_version}"

          port {
            container_port = 8200
            name           = "http"
          }

          port {
            container_port = 8201
            name           = "internal"
          }

          env {
            name = "VAULT_DEV_ROOT_TOKEN_ID"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.vault_secret.metadata[0].name
                key  = "root-token"
              }
            }
          }

          env {
            name  = "VAULT_DEV_LISTEN_ADDRESS"
            value = "0.0.0.0:8200"
          }

          volume_mount {
            name       = "vault-data"
            mount_path = "/vault/data"
          }
        }

        volume {
          name = "vault-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.vault_data_pvc.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "vault" {
  metadata {
    name      = "vault"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    selector = {
      app = "vault"
    }

    port {
      port        = 8200
      target_port = 8200
      name        = "http"
    }

    port {
      port        = 8201
      target_port = 8201
      name        = "internal"
    }
  }
}

resource "kubernetes_persistent_volume_claim" "vault_data_pvc" {
  metadata {
    name      = "vault-data-pvc"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "5Gi"
      }
    }
  }
}

resource "kubernetes_secret" "vault_secret" {
  metadata {
    name      = "vault-secret"
    namespace = kubernetes_namespace.kubefirst.metadata[0].name
  }

  data = {
    "root-token" = var.vault_token
  }
}
