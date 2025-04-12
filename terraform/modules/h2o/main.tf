
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "h2o" {
  metadata {
    name = var.h2o_namespace
    labels = {
      "app.kubernetes.io/name" = "h2o"
      "app.kubernetes.io/instance" = "h2o"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_persistent_volume_claim" "h2o_data" {
  metadata {
    name      = "h2o-data"
    namespace = kubernetes_namespace.h2o.metadata[0].name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.h2o_storage_size
      }
    }
    storage_class_name = var.storage_class_name
  }
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.h2o.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "h2o_config" {
  metadata {
    name      = "h2o-config"
    namespace = kubernetes_namespace.h2o.metadata[0].name
  }

  data = {
    "config.yaml" = <<-EOT
      h2o:
        version: ${var.h2o_version}
        resources:
          limits:
            cpu: 4
            memory: 16Gi
          requests:
            cpu: 2
            memory: 8Gi
        storage:
          size: ${var.h2o_storage_size}
          class: ${var.storage_class_name}
        s3:
          endpoint: http://minio.minio.svc.cluster.local:9000
          bucket: h2o
          region: us-east-1
          verify_ssl: false
        mlflow:
          tracking_uri: http://mlflow-server.mlflow.svc.cluster.local:5000
    EOT
  }
}

resource "kubernetes_deployment" "h2o" {
  metadata {
    name      = "h2o"
    namespace = kubernetes_namespace.h2o.metadata[0].name
    labels = {
      app = "h2o"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "h2o"
      }
    }

    template {
      metadata {
        labels = {
          app = "h2o"
        }
      }

      spec {
        container {
          name  = "h2o"
          image = "h2oai/h2o-automl:${var.h2o_version}"

          port {
            container_port = 54321
          }

          volume_mount {
            name       = "h2o-data"
            mount_path = "/data"
          }

          volume_mount {
            name       = "h2o-config"
            mount_path = "/etc/h2o"
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
            name  = "AWS_ENDPOINT_URL"
            value = "http://minio.minio.svc.cluster.local:9000"
          }

          env {
            name  = "MLFLOW_TRACKING_URI"
            value = "http://mlflow-server.mlflow.svc.cluster.local:5000"
          }

          resources {
            limits = {
              cpu    = "4"
              memory = "16Gi"
            }
            requests = {
              cpu    = "2"
              memory = "8Gi"
            }
          }
        }

        volume {
          name = "h2o-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.h2o_data.metadata[0].name
          }
        }

        volume {
          name = "h2o-config"
          config_map {
            name = kubernetes_config_map.h2o_config.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "h2o" {
  metadata {
    name      = "h2o"
    namespace = kubernetes_namespace.h2o.metadata[0].name
  }
  spec {
    selector = {
      app = kubernetes_deployment.h2o.metadata[0].labels.app
    }
    port {
      port        = 54321
      target_port = 54321
    }
  }
}

resource "kubernetes_config_map" "h2o_automl_config" {
  metadata {
    name      = "h2o-automl-config"
    namespace = kubernetes_namespace.h2o.metadata[0].name
  }

  data = {
    "automl_config.yaml" = <<-EOT
      automl:
        max_models: 20
        max_runtime_secs: 3600
        stopping_metric: "AUC"
        sort_metric: "AUC"
        seed: 42
        balance_classes: true
        class_sampling_factors: null
        max_after_balance_size: null
        keep_cross_validation_predictions: true
        keep_cross_validation_models: true
        keep_cross_validation_fold_assignment: true
        nfolds: 5
        fold_column: null
        ignored_columns: null
        exclude_algos: null
        include_algos: null
        project_name: "llama4-fine-tuning"
    EOT
  }
}
