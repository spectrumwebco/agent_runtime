terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.20.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.9.0"
    }
  }
}

resource "kubernetes_namespace" "vnode_runtime" {
  metadata {
    name = var.namespace
    labels = {
      app = "vnode-runtime"
    }
  }
}

resource "helm_release" "vnode_runtime" {
  name       = "vnode-runtime"
  repository = "https://charts.loft.sh"
  chart      = "vnode-runtime"
  version    = var.vnode_runtime_version
  namespace  = kubernetes_namespace.vnode_runtime.metadata[0].name

  set {
    name  = "replicaCount"
    value = var.replica_count
  }

  set {
    name  = "resources.limits.cpu"
    value = var.resources_limits_cpu
  }

  set {
    name  = "resources.limits.memory"
    value = var.resources_limits_memory
  }

  set {
    name  = "resources.requests.cpu"
    value = var.resources_requests_cpu
  }

  set {
    name  = "resources.requests.memory"
    value = var.resources_requests_memory
  }
}

resource "kubernetes_config_map" "vnode_runtime_integration" {
  metadata {
    name      = "vnode-runtime-integration"
    namespace = kubernetes_namespace.vnode_runtime.metadata[0].name
  }

  data = {
    "kubeflow-integration.yaml" = <<-EOT
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: vnode-kubeflow-integration
        namespace: kubeflow
      data:
        vnode-runtime-enabled: "true"
        vnode-runtime-endpoint: "http://vnode-runtime.${kubernetes_namespace.vnode_runtime.metadata[0].name}.svc.cluster.local:8080"
    EOT

    "mlflow-integration.yaml" = <<-EOT
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: vnode-mlflow-integration
        namespace: mlflow
      data:
        vnode-runtime-enabled: "true"
        vnode-runtime-endpoint: "http://vnode-runtime.${kubernetes_namespace.vnode_runtime.metadata[0].name}.svc.cluster.local:8080"
    EOT

    "kserve-integration.yaml" = <<-EOT
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: vnode-kserve-integration
        namespace: kserve
      data:
        vnode-runtime-enabled: "true"
        vnode-runtime-endpoint: "http://vnode-runtime.${kubernetes_namespace.vnode_runtime.metadata[0].name}.svc.cluster.local:8080"
    EOT
  }
}

resource "kubernetes_deployment" "kubeflow_vnode_integration" {
  metadata {
    name      = "kubeflow-vnode-integration"
    namespace = "kubeflow"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "kubeflow-vnode-integration"
      }
    }

    template {
      metadata {
        labels = {
          app = "kubeflow-vnode-integration"
        }
      }

      spec {
        container {
          name  = "vnode-integration"
          image = "loftsh/vnode-runtime-agent:${var.vnode_runtime_version}"

          env {
            name = "VNODE_RUNTIME_ENDPOINT"
            value_from {
              config_map_key_ref {
                name = "vnode-kubeflow-integration"
                key  = "vnode-runtime-endpoint"
              }
            }
          }

          env {
            name = "VNODE_RUNTIME_ENABLED"
            value_from {
              config_map_key_ref {
                name = "vnode-kubeflow-integration"
                key  = "vnode-runtime-enabled"
              }
            }
          }
        }
      }
    }
  }

  depends_on = [kubernetes_config_map.vnode_runtime_integration]
}

resource "kubernetes_deployment" "mlflow_vnode_integration" {
  metadata {
    name      = "mlflow-vnode-integration"
    namespace = "mlflow"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "mlflow-vnode-integration"
      }
    }

    template {
      metadata {
        labels = {
          app = "mlflow-vnode-integration"
        }
      }

      spec {
        container {
          name  = "vnode-integration"
          image = "loftsh/vnode-runtime-agent:${var.vnode_runtime_version}"

          env {
            name = "VNODE_RUNTIME_ENDPOINT"
            value_from {
              config_map_key_ref {
                name = "vnode-mlflow-integration"
                key  = "vnode-runtime-endpoint"
              }
            }
          }

          env {
            name = "VNODE_RUNTIME_ENABLED"
            value_from {
              config_map_key_ref {
                name = "vnode-mlflow-integration"
                key  = "vnode-runtime-enabled"
              }
            }
          }
        }
      }
    }
  }

  depends_on = [kubernetes_config_map.vnode_runtime_integration]
}

resource "kubernetes_deployment" "kserve_vnode_integration" {
  metadata {
    name      = "kserve-vnode-integration"
    namespace = "kserve"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "kserve-vnode-integration"
      }
    }

    template {
      metadata {
        labels = {
          app = "kserve-vnode-integration"
        }
      }

      spec {
        container {
          name  = "vnode-integration"
          image = "loftsh/vnode-runtime-agent:${var.vnode_runtime_version}"

          env {
            name = "VNODE_RUNTIME_ENDPOINT"
            value_from {
              config_map_key_ref {
                name = "vnode-kserve-integration"
                key  = "vnode-runtime-endpoint"
              }
            }
          }

          env {
            name = "VNODE_RUNTIME_ENABLED"
            value_from {
              config_map_key_ref {
                name = "vnode-kserve-integration"
                key  = "vnode-runtime-enabled"
              }
            }
          }
        }
      }
    }
  }

  depends_on = [kubernetes_config_map.vnode_runtime_integration]
}
