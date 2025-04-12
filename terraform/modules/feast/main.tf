
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "feast" {
  metadata {
    name = var.feast_namespace
    labels = {
      "app.kubernetes.io/name" = "feast"
      "app.kubernetes.io/instance" = "feast"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_config_map" "feast_config" {
  metadata {
    name      = "feast-config"
    namespace = kubernetes_namespace.feast.metadata[0].name
  }

  data = {
    "feature_store.yaml" = <<-EOT
      project: llama4-fine-tuning
      registry: s3://feast/registry.db
      provider: local
      online_store:
        type: redis
        connection_string: redis.${kubernetes_namespace.feast.metadata[0].name}.svc.cluster.local:6379
      offline_store:
        type: file
      entity_key_serialization_version: 2
    EOT
  }
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.feast.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}

resource "kubernetes_deployment" "feast_redis" {
  metadata {
    name      = "redis"
    namespace = kubernetes_namespace.feast.metadata[0].name
    labels = {
      app = "redis"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "redis"
      }
    }

    template {
      metadata {
        labels = {
          app = "redis"
        }
      }

      spec {
        container {
          name  = "redis"
          image = "redis:6.2-alpine"

          port {
            container_port = 6379
          }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "feast_redis" {
  metadata {
    name      = "redis"
    namespace = kubernetes_namespace.feast.metadata[0].name
  }
  spec {
    selector = {
      app = kubernetes_deployment.feast_redis.metadata[0].labels.app
    }
    port {
      port        = 6379
      target_port = 6379
    }
  }
}

resource "kubernetes_deployment" "feast_server" {
  metadata {
    name      = "feast-server"
    namespace = kubernetes_namespace.feast.metadata[0].name
    labels = {
      app = "feast-server"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "feast-server"
      }
    }

    template {
      metadata {
        labels = {
          app = "feast-server"
        }
      }

      spec {
        container {
          name  = "feast-server"
          image = "feastdev/feature-server:${var.feast_version}"

          port {
            container_port = 6566
          }

          volume_mount {
            name       = "feast-config"
            mount_path = "/etc/feast"
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
            name  = "S3_ENDPOINT_URL"
            value = "http://minio.minio.svc.cluster.local:9000"
          }

          env {
            name  = "FEAST_S3_ENDPOINT_URL"
            value = "http://minio.minio.svc.cluster.local:9000"
          }

          env {
            name  = "FEAST_FEATURE_STORE_CONFIG_PATH"
            value = "/etc/feast/feature_store.yaml"
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
          name = "feast-config"
          config_map {
            name = kubernetes_config_map.feast_config.metadata[0].name
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_service.feast_redis
  ]
}

resource "kubernetes_service" "feast_server" {
  metadata {
    name      = "feast-server"
    namespace = kubernetes_namespace.feast.metadata[0].name
  }
  spec {
    selector = {
      app = kubernetes_deployment.feast_server.metadata[0].labels.app
    }
    port {
      port        = 6566
      target_port = 6566
    }
  }
}

resource "kubernetes_config_map" "feast_features" {
  metadata {
    name      = "feast-features"
    namespace = kubernetes_namespace.feast.metadata[0].name
  }

  data = {
    "issue_features.py" = <<-EOT
      from datetime import timedelta
      from feast import Entity, Feature, FeatureView, FileSource, ValueType
      from feast.types import Float32, Int64, String

      issue = Entity(
          name="issue",
          value_type=ValueType.STRING,
          description="GitHub issue identifier",
      )

      issue_source = FileSource(
          path="s3://datasets/github_issues.parquet",
          event_timestamp_column="timestamp",
      )

      issue_features = FeatureView(
          name="issue_features",
          entities=["issue"],
          ttl=timedelta(days=365),
          features=[
              Feature(name="title_length", dtype=Int64),
              Feature(name="description_length", dtype=Int64),
              Feature(name="num_comments", dtype=Int64),
              Feature(name="num_labels", dtype=Int64),
              Feature(name="has_code", dtype=Int64),
              Feature(name="repository", dtype=String),
              Feature(name="language", dtype=String),
              Feature(name="time_to_resolution", dtype=Float32),
              Feature(name="complexity_score", dtype=Float32),
          ],
          online=True,
          input=issue_source,
          tags={"team": "ml", "owner": "data-science"},
      )
    EOT
  }
}
