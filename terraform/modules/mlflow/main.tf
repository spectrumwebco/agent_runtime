
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "mlflow" {
  metadata {
    name = var.mlflow_namespace
    labels = {
      "app.kubernetes.io/name" = "mlflow"
      "app.kubernetes.io/instance" = "mlflow"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_secret" "mlflow_db_secret" {
  metadata {
    name      = "mlflow-db-secret"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }

  data = {
    "username" = "mlflow"
    "password" = "mlflow-password"
    "database" = "mlflow"
    "host"     = "mlflow-db.${kubernetes_namespace.mlflow.metadata[0].name}.svc.cluster.local"
    "port"     = "5432"
  }

  type = "Opaque"
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "mlflow_config" {
  metadata {
    name      = "mlflow-config"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }

  data = {
    "MLFLOW_S3_ENDPOINT_URL" = "http://minio.minio.svc.cluster.local:9000"
    "MLFLOW_S3_IGNORE_TLS"   = "true"
    "MLFLOW_TRACKING_URI"    = var.mlflow_tracking_uri
    "DEFAULT_ARTIFACT_ROOT"  = "s3://mlflow/"
  }
}

resource "kubernetes_persistent_volume_claim" "mlflow_data" {
  metadata {
    name      = "mlflow-data"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.mlflow_storage_size
      }
    }
    storage_class_name = var.storage_class_name
  }
}

resource "kubernetes_deployment" "mlflow_db" {
  metadata {
    name      = "mlflow-db"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
    labels = {
      app = "mlflow-db"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "mlflow-db"
      }
    }

    template {
      metadata {
        labels = {
          app = "mlflow-db"
        }
      }

      spec {
        container {
          name  = "postgres"
          image = "postgres:13"

          port {
            container_port = 5432
          }

          env {
            name = "POSTGRES_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "username"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "password"
              }
            }
          }

          env {
            name = "POSTGRES_DB"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "database"
              }
            }
          }

          volume_mount {
            name       = "mlflow-db-data"
            mount_path = "/var/lib/postgresql/data"
          }

          resources {
            limits = {
              cpu    = "1000m"
              memory = "1Gi"
            }
            requests = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
        }

        volume {
          name = "mlflow-db-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.mlflow_data.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "mlflow_db" {
  metadata {
    name      = "mlflow-db"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }
  spec {
    selector = {
      app = kubernetes_deployment.mlflow_db.metadata[0].labels.app
    }
    port {
      port        = 5432
      target_port = 5432
    }
  }
}

resource "kubernetes_deployment" "mlflow_server" {
  metadata {
    name      = "mlflow-server"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
    labels = {
      app = "mlflow-server"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "mlflow-server"
      }
    }

    template {
      metadata {
        labels = {
          app = "mlflow-server"
        }
      }

      spec {
        container {
          name  = "mlflow"
          image = "ghcr.io/mlflow/mlflow:v${var.mlflow_version}"

          port {
            container_port = 5000
          }

          args = [
            "server",
            "--host=0.0.0.0",
            "--port=5000",
            "--backend-store-uri=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)",
            "--default-artifact-root=s3://mlflow/",
            "--artifacts-destination=s3://mlflow/"
          ]

          env {
            name = "POSTGRES_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "username"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "password"
              }
            }
          }

          env {
            name = "POSTGRES_DB"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "database"
              }
            }
          }

          env {
            name = "POSTGRES_HOST"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "host"
              }
            }
          }

          env {
            name = "POSTGRES_PORT"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mlflow_db_secret.metadata[0].name
                key  = "port"
              }
            }
          }

          env {
            name = "AWS_ACCESS_KEY_ID"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.minio_credentials.metadata[0].name
                key  = "accesskey"
              }
            }
          }

          env {
            name = "AWS_SECRET_ACCESS_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.minio_credentials.metadata[0].name
                key  = "secretkey"
              }
            }
          }

          env {
            name  = "MLFLOW_S3_ENDPOINT_URL"
            value = "http://minio.minio.svc.cluster.local:9000"
          }

          env {
            name  = "AWS_ENDPOINT_URL"
            value = "http://minio.minio.svc.cluster.local:9000"
          }

          env {
            name  = "MLFLOW_S3_IGNORE_TLS"
            value = "true"
          }

          resources {
            limits = {
              cpu    = "1000m"
              memory = "1Gi"
            }
            requests = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_service.mlflow_db
  ]
}

resource "kubernetes_service" "mlflow_server" {
  metadata {
    name      = "mlflow-server"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }
  spec {
    selector = {
      app = kubernetes_deployment.mlflow_server.metadata[0].labels.app
    }
    port {
      port        = 5000
      target_port = 5000
    }
  }
}

resource "kubernetes_ingress_v1" "mlflow_ingress" {
  metadata {
    name      = "mlflow-ingress"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
    annotations = {
      "kubernetes.io/ingress.class"                    = "nginx"
      "nginx.ingress.kubernetes.io/ssl-redirect"       = "false"
      "nginx.ingress.kubernetes.io/proxy-body-size"    = "0"
      "nginx.ingress.kubernetes.io/proxy-read-timeout" = "600"
      "nginx.ingress.kubernetes.io/proxy-send-timeout" = "600"
    }
  }
  spec {
    rule {
      host = "mlflow.example.com"
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.mlflow_server.metadata[0].name
              port {
                number = 5000
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_config_map" "llama4_experiment_config" {
  metadata {
    name      = "llama4-experiment-config"
    namespace = kubernetes_namespace.mlflow.metadata[0].name
  }

  data = {
    "llama4-maverick-experiment.json" = jsonencode({
      name = "llama4-maverick-fine-tuning"
      tags = {
        model_type = "llama4-maverick"
        task       = "fine-tuning"
        domain     = "software-engineering"
      }
      artifact_location = "s3://mlflow/llama4-maverick-fine-tuning"
    })
    "llama4-scout-experiment.json" = jsonencode({
      name = "llama4-scout-fine-tuning"
      tags = {
        model_type = "llama4-scout"
        task       = "fine-tuning"
        domain     = "software-engineering"
      }
      artifact_location = "s3://mlflow/llama4-scout-fine-tuning"
    })
  }
}
