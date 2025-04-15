provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "minio" {
  metadata {
    name = var.minio_namespace
    labels = {
      "app.kubernetes.io/name" = "minio"
      "app.kubernetes.io/instance" = "minio"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.minio.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}

resource "kubernetes_persistent_volume_claim" "minio_data" {
  metadata {
    name      = "minio-data"
    namespace = kubernetes_namespace.minio.metadata[0].name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.minio_storage_size
      }
    }
    storage_class_name = var.storage_class_name
  }
}

resource "helm_release" "minio" {
  name       = "minio"
  repository = "https://charts.min.io/"
  chart      = "minio"
  namespace  = kubernetes_namespace.minio.metadata[0].name
  version    = var.minio_version
  timeout    = 600

  set {
    name  = "mode"
    value = "standalone"
  }

  set {
    name  = "persistence.enabled"
    value = "true"
  }

  set {
    name  = "persistence.existingClaim"
    value = kubernetes_persistent_volume_claim.minio_data.metadata[0].name
  }

  set {
    name  = "accessKey"
    value = var.minio_access_key
  }

  set {
    name  = "secretKey"
    value = var.minio_secret_key
  }

  set {
    name  = "resources.requests.memory"
    value = "1Gi"
  }

  set {
    name  = "resources.requests.cpu"
    value = "250m"
  }

  set {
    name  = "resources.limits.memory"
    value = "2Gi"
  }

  set {
    name  = "resources.limits.cpu"
    value = "500m"
  }

  depends_on = [
    kubernetes_namespace.minio,
    kubernetes_persistent_volume_claim.minio_data
  ]
}

resource "kubernetes_config_map" "minio_bucket_config" {
  metadata {
    name      = "minio-bucket-config"
    namespace = kubernetes_namespace.minio.metadata[0].name
  }

  data = {
    "create-buckets.sh" = <<-EOT
      
      wget -q https://dl.min.io/client/mc/release/linux-amd64/mc -O /usr/local/bin/mc
      chmod +x /usr/local/bin/mc
      
      mc config host add minio http://minio.${kubernetes_namespace.minio.metadata[0].name}.svc.cluster.local:9000 ${var.minio_access_key} ${var.minio_secret_key} --api s3v4
      
      mc mb --ignore-existing minio/mlflow
      mc mb --ignore-existing minio/models
      mc mb --ignore-existing minio/datasets
      mc mb --ignore-existing minio/feast
      
      mc policy set download minio/mlflow
      mc policy set download minio/models
      mc policy set download minio/datasets
      mc policy set download minio/feast
      
      echo "MinIO buckets created and configured successfully"
    EOT
  }

  depends_on = [
    helm_release.minio
  ]
}

resource "kubernetes_job" "minio_init" {
  metadata {
    name      = "minio-init"
    namespace = kubernetes_namespace.minio.metadata[0].name
  }

  spec {
    template {
      metadata {
        labels = {
          app = "minio-init"
        }
      }

      spec {
        container {
          name    = "minio-init"
          image   = "alpine:3.15"
          command = ["/bin/sh", "-c", "/scripts/create-buckets.sh"]

          volume_mount {
            name       = "scripts"
            mount_path = "/scripts"
          }

          env {
            name = "MINIO_ACCESS_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.minio_credentials.metadata[0].name
                key  = "accesskey"
              }
            }
          }

          env {
            name = "MINIO_SECRET_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.minio_credentials.metadata[0].name
                key  = "secretkey"
              }
            }
          }
        }

        volume {
          name = "scripts"
          config_map {
            name = kubernetes_config_map.minio_bucket_config.metadata[0].name
            default_mode = "0755"
          }
        }

        restart_policy = "OnFailure"
      }
    }

    backoff_limit = 3
  }

  depends_on = [
    helm_release.minio,
    kubernetes_config_map.minio_bucket_config
  ]
}
