
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "kserve" {
  metadata {
    name = var.kserve_namespace
    labels = {
      "app.kubernetes.io/name" = "kserve"
      "app.kubernetes.io/instance" = "kserve"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.kserve.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}

resource "helm_release" "kserve" {
  name       = "kserve"
  repository = "https://kserve.github.io/kserve"
  chart      = "kserve"
  namespace  = kubernetes_namespace.kserve.metadata[0].name
  version    = var.kserve_version
  timeout    = 1200

  set {
    name  = "kserve.enabled"
    value = "true"
  }

  set {
    name  = "knative.enabled"
    value = "true"
  }

  set {
    name  = "models.enabled"
    value = "true"
  }

  depends_on = [
    kubernetes_namespace.kserve
  ]
}

resource "kubernetes_config_map" "kserve_config" {
  metadata {
    name      = "kserve-config"
    namespace = kubernetes_namespace.kserve.metadata[0].name
  }

  data = {
    "storageClassName" = var.storage_class_name
    "ingressGateway" = "knative-serving/knative-ingress-gateway"
    "ingressService" = "istio-ingressgateway.istio-system.svc.cluster.local"
    "localGateway" = "knative-serving/knative-local-gateway"
    "localGatewayService" = "knative-local-gateway.istio-system.svc.cluster.local"
  }

  depends_on = [
    helm_release.kserve
  ]
}

resource "kubernetes_config_map" "llama4_model_config" {
  metadata {
    name      = "llama4-model-config"
    namespace = kubernetes_namespace.kserve.metadata[0].name
  }

  data = {
    "llama4-maverick-model.yaml" = jsonencode({
      apiVersion = "serving.kserve.io/v1beta1"
      kind       = "InferenceService"
      metadata = {
        name      = "llama4-maverick"
        namespace = kubernetes_namespace.kserve.metadata[0].name
      }
      spec = {
        predictor = {
          model = {
            modelFormat = {
              name = "pytorch"
            }
            storageUri = "s3://models/llama4-maverick"
            resources = {
              limits = {
                cpu    = "4"
                memory = "16Gi"
                "nvidia.com/gpu" = "1"
              }
              requests = {
                cpu    = "2"
                memory = "8Gi"
              }
            }
            env = [
              {
                name  = "STORAGE_URI"
                value = "s3://models/llama4-maverick"
              },
              {
                name  = "MODEL_NAME"
                value = "llama4-maverick"
              },
              {
                name = "AWS_ACCESS_KEY_ID"
                valueFrom = {
                  secretKeyRef = {
                    name = "minio-credentials"
                    key  = "accesskey"
                  }
                }
              },
              {
                name = "AWS_SECRET_ACCESS_KEY"
                valueFrom = {
                  secretKeyRef = {
                    name = "minio-credentials"
                    key  = "secretkey"
                  }
                }
              },
              {
                name  = "AWS_ENDPOINT_URL"
                value = "http://minio.minio.svc.cluster.local:9000"
              },
              {
                name  = "S3_USE_HTTPS"
                value = "0"
              },
              {
                name  = "S3_VERIFY_SSL"
                value = "0"
              }
            ]
          }
        }
      }
    })
    "llama4-scout-model.yaml" = jsonencode({
      apiVersion = "serving.kserve.io/v1beta1"
      kind       = "InferenceService"
      metadata = {
        name      = "llama4-scout"
        namespace = kubernetes_namespace.kserve.metadata[0].name
      }
      spec = {
        predictor = {
          model = {
            modelFormat = {
              name = "pytorch"
            }
            storageUri = "s3://models/llama4-scout"
            resources = {
              limits = {
                cpu    = "4"
                memory = "16Gi"
                "nvidia.com/gpu" = "1"
              }
              requests = {
                cpu    = "2"
                memory = "8Gi"
              }
            }
            env = [
              {
                name  = "STORAGE_URI"
                value = "s3://models/llama4-scout"
              },
              {
                name  = "MODEL_NAME"
                value = "llama4-scout"
              },
              {
                name = "AWS_ACCESS_KEY_ID"
                valueFrom = {
                  secretKeyRef = {
                    name = "minio-credentials"
                    key  = "accesskey"
                  }
                }
              },
              {
                name = "AWS_SECRET_ACCESS_KEY"
                valueFrom = {
                  secretKeyRef = {
                    name = "minio-credentials"
                    key  = "secretkey"
                  }
                }
              },
              {
                name  = "AWS_ENDPOINT_URL"
                value = "http://minio.minio.svc.cluster.local:9000"
              },
              {
                name  = "S3_USE_HTTPS"
                value = "0"
              },
              {
                name  = "S3_VERIFY_SSL"
                value = "0"
              }
            ]
          }
        }
      }
    })
  }

  depends_on = [
    helm_release.kserve
  ]
}
