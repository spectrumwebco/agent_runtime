
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "seldon" {
  metadata {
    name = var.seldon_namespace
    labels = {
      "app.kubernetes.io/name" = "seldon"
      "app.kubernetes.io/instance" = "seldon"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "helm_release" "seldon_core" {
  name       = "seldon-core"
  repository = "https://storage.googleapis.com/seldon-charts"
  chart      = "seldon-core-operator"
  namespace  = kubernetes_namespace.seldon.metadata[0].name
  version    = var.seldon_version
  timeout    = 600

  set {
    name  = "usageMetrics.enabled"
    value = "false"
  }

  set {
    name  = "istio.enabled"
    value = "true"
  }

  set {
    name  = "certManager.enabled"
    value = "false"
  }

  depends_on = [
    kubernetes_namespace.seldon
  ]
}

resource "kubernetes_config_map" "seldon_model_config" {
  metadata {
    name      = "seldon-model-config"
    namespace = kubernetes_namespace.seldon.metadata[0].name
  }

  data = {
    "llama4-maverick-seldon.yaml" = jsonencode({
      apiVersion = "machinelearning.seldon.io/v1"
      kind       = "SeldonDeployment"
      metadata = {
        name      = "llama4-maverick"
        namespace = kubernetes_namespace.seldon.metadata[0].name
      }
      spec = {
        name = "llama4-maverick"
        predictors = [
          {
            name        = "default"
            replicas    = 1
            annotations = {
              "seldon.io/no-engine" = "true"
            }
            graph = {
              name        = "llama4-maverick-container"
              type        = "MODEL"
              implementation = "TRITON_SERVER"
              modelUri    = "s3://models/llama4-maverick"
              parameters = [
                {
                  name  = "model_name"
                  value = "llama4-maverick"
                  type  = "STRING"
                },
                {
                  name  = "signature_name"
                  value = "serving_default"
                  type  = "STRING"
                }
              ]
              env = [
                {
                  name  = "AWS_ENDPOINT_URL"
                  value = "http://minio.minio.svc.cluster.local:9000"
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
                  name  = "S3_USE_HTTPS"
                  value = "0"
                },
                {
                  name  = "S3_VERIFY_SSL"
                  value = "0"
                }
              ]
            }
            componentSpecs = [
              {
                spec = {
                  containers = [
                    {
                      name  = "llama4-maverick-container"
                      image = "nvcr.io/nvidia/tritonserver:22.12-py3"
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
                    }
                  ]
                }
              }
            ]
          }
        ]
      }
    })
    "llama4-scout-seldon.yaml" = jsonencode({
      apiVersion = "machinelearning.seldon.io/v1"
      kind       = "SeldonDeployment"
      metadata = {
        name      = "llama4-scout"
        namespace = kubernetes_namespace.seldon.metadata[0].name
      }
      spec = {
        name = "llama4-scout"
        predictors = [
          {
            name        = "default"
            replicas    = 1
            annotations = {
              "seldon.io/no-engine" = "true"
            }
            graph = {
              name        = "llama4-scout-container"
              type        = "MODEL"
              implementation = "TRITON_SERVER"
              modelUri    = "s3://models/llama4-scout"
              parameters = [
                {
                  name  = "model_name"
                  value = "llama4-scout"
                  type  = "STRING"
                },
                {
                  name  = "signature_name"
                  value = "serving_default"
                  type  = "STRING"
                }
              ]
              env = [
                {
                  name  = "AWS_ENDPOINT_URL"
                  value = "http://minio.minio.svc.cluster.local:9000"
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
                  name  = "S3_USE_HTTPS"
                  value = "0"
                },
                {
                  name  = "S3_VERIFY_SSL"
                  value = "0"
                }
              ]
            }
            componentSpecs = [
              {
                spec = {
                  containers = [
                    {
                      name  = "llama4-scout-container"
                      image = "nvcr.io/nvidia/tritonserver:22.12-py3"
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
                    }
                  ]
                }
              }
            ]
          }
        ]
      }
    })
  }

  depends_on = [
    helm_release.seldon_core
  ]
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.seldon.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}
