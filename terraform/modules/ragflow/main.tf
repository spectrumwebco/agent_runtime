locals {
  namespace = var.create_namespace ? kubernetes_namespace.ragflow[0].metadata[0].name : var.namespace
  labels = merge({
    "app.kubernetes.io/name"       = "ragflow"
    "app.kubernetes.io/instance"   = var.name
    "app.kubernetes.io/component"  = "ragflow"
    "app.kubernetes.io/part-of"    = "agent-runtime"
    "app.kubernetes.io/managed-by" = "terraform"
  }, var.ragflow_labels)
  annotations = merge({
    "agent-runtime.io/component" = "ragflow"
  }, var.ragflow_annotations)
}

resource "kubernetes_namespace" "ragflow" {
  count = var.create_namespace ? 1 : 0
  
  metadata {
    name = var.namespace
    
    labels = local.labels
    annotations = local.annotations
  }
}

resource "helm_release" "ragflow_vcluster" {
  count = var.vcluster_enabled ? 1 : 0
  
  name       = var.vcluster_name
  namespace  = var.vcluster_namespace
  repository = "https://charts.loft.sh"
  chart      = "vcluster"
  version    = "0.15.0"
  
  set {
    name  = "syncer.extraArgs"
    value = "--tls-san=${var.vcluster_name}.${var.vcluster_namespace}"
  }
  
  set {
    name  = "sync.nodes.enabled"
    value = "true"
  }
  
  set {
    name  = "sync.nodes.nodeSelector"
    value = "ragflow=true"
  }
  
  set {
    name  = "isolation.enabled"
    value = "true"
  }
  
  set {
    name  = "service.type"
    value = "ClusterIP"
  }
  
  set {
    name  = "init.manifests"
    value = <<-EOT
      apiVersion: v1
      kind: Namespace
      metadata:
        name: ${local.namespace}
      ---
      apiVersion: v1
      kind: ServiceAccount
      metadata:
        name: ragflow
        namespace: ${local.namespace}
      ---
      apiVersion: rbac.authorization.k8s.io/v1
      kind: Role
      metadata:
        name: ragflow
        namespace: ${local.namespace}
      rules:
      - apiGroups: [""]
        resources: ["pods", "services", "configmaps", "secrets"]
        verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
      ---
      apiVersion: rbac.authorization.k8s.io/v1
      kind: RoleBinding
      metadata:
        name: ragflow
        namespace: ${local.namespace}
      subjects:
      - kind: ServiceAccount
        name: ragflow
        namespace: ${local.namespace}
      roleRef:
        kind: Role
        name: ragflow
        apiGroup: rbac.authorization.k8s.io
    EOT
  }
}

resource "kubernetes_config_map" "ragflow_config" {
  metadata {
    name      = "${var.name}-config"
    namespace = local.namespace
    labels    = local.labels
    annotations = local.annotations
  }
  
  data = {
    "config.yaml" = yamlencode({
      embedding_model = var.ragflow_config.embedding_model
      llm_model       = var.ragflow_config.llm_model
      vector_db       = var.ragflow_config.vector_db
      document_store  = var.ragflow_config.document_store
      api = {
        port           = var.ragflow_api_config.port
        max_tokens     = var.ragflow_api_config.max_tokens
        temperature    = var.ragflow_api_config.temperature
        top_p          = var.ragflow_api_config.top_p
        top_k          = var.ragflow_api_config.top_k
        chunk_size     = var.ragflow_api_config.chunk_size
        chunk_overlap  = var.ragflow_api_config.chunk_overlap
        max_documents  = var.ragflow_api_config.max_documents
      }
    })
  }
}

resource "kubernetes_secret" "ragflow_api_key" {
  metadata {
    name      = "${var.name}-api-key"
    namespace = local.namespace
    labels    = local.labels
    annotations = local.annotations
  }
  
  data = {
    "api-key" = var.ragflow_config.api_key
  }
  
  type = "Opaque"
}

resource "kubernetes_deployment" "qdrant" {
  count = var.ragflow_qdrant_config.enabled && var.ragflow_qdrant_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-qdrant"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "qdrant"
    })
    annotations = local.annotations
  }
  
  spec {
    replicas = var.enable_high_availability ? var.ragflow_qdrant_config.replicas : 1
    
    selector {
      match_labels = {
        app = "${var.name}-qdrant"
      }
    }
    
    template {
      metadata {
        labels = merge(local.labels, {
          app = "${var.name}-qdrant"
          "app.kubernetes.io/component" = "qdrant"
        })
        annotations = local.annotations
      }
      
      spec {
        container {
          name  = "qdrant"
          image = "qdrant/qdrant:latest"
          
          port {
            container_port = 6333
            name           = "http"
          }
          
          port {
            container_port = 6334
            name           = "grpc"
          }
          
          resources {
            limits = {
              cpu    = var.ragflow_qdrant_config.resources.limits.cpu
              memory = var.ragflow_qdrant_config.resources.limits.memory
            }
            requests = {
              cpu    = var.ragflow_qdrant_config.resources.requests.cpu
              memory = var.ragflow_qdrant_config.resources.requests.memory
            }
          }
          
          volume_mount {
            name       = "qdrant-data"
            mount_path = "/qdrant/storage"
          }
          
          liveness_probe {
            http_get {
              path = "/health"
              port = 6333
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            http_get {
              path = "/health"
              port = 6333
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
        
        volume {
          name = "qdrant-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.qdrant_data[0].metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_persistent_volume_claim" "qdrant_data" {
  count = var.ragflow_qdrant_config.enabled && var.ragflow_qdrant_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-qdrant-data"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "qdrant"
    })
    annotations = local.annotations
  }
  
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.ragflow_storage_size
      }
    }
    storage_class_name = var.ragflow_storage_class
  }
}

resource "kubernetes_service" "qdrant" {
  count = var.ragflow_qdrant_config.enabled && var.ragflow_qdrant_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-qdrant"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "qdrant"
    })
    annotations = local.annotations
  }
  
  spec {
    selector = {
      app = "${var.name}-qdrant"
    }
    
    port {
      port        = 6333
      target_port = 6333
      name        = "http"
    }
    
    port {
      port        = 6334
      target_port = 6334
      name        = "grpc"
    }
  }
}

resource "kubernetes_deployment" "minio" {
  count = var.ragflow_minio_config.enabled && var.ragflow_minio_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-minio"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "minio"
    })
    annotations = local.annotations
  }
  
  spec {
    replicas = var.enable_high_availability ? var.ragflow_minio_config.replicas : 1
    
    selector {
      match_labels = {
        app = "${var.name}-minio"
      }
    }
    
    template {
      metadata {
        labels = merge(local.labels, {
          app = "${var.name}-minio"
          "app.kubernetes.io/component" = "minio"
        })
        annotations = local.annotations
      }
      
      spec {
        container {
          name  = "minio"
          image = "minio/minio:latest"
          args  = ["server", "/data", "--console-address", ":9001"]
          
          port {
            container_port = 9000
            name           = "api"
          }
          
          port {
            container_port = 9001
            name           = "console"
          }
          
          env {
            name  = "MINIO_ROOT_USER"
            value = var.ragflow_minio_config.access_key != "" ? var.ragflow_minio_config.access_key : "minioadmin"
          }
          
          env {
            name  = "MINIO_ROOT_PASSWORD"
            value = var.ragflow_minio_config.secret_key != "" ? var.ragflow_minio_config.secret_key : "minioadmin"
          }
          
          resources {
            limits = {
              cpu    = var.ragflow_minio_config.resources.limits.cpu
              memory = var.ragflow_minio_config.resources.limits.memory
            }
            requests = {
              cpu    = var.ragflow_minio_config.resources.requests.cpu
              memory = var.ragflow_minio_config.resources.requests.memory
            }
          }
          
          volume_mount {
            name       = "minio-data"
            mount_path = "/data"
          }
          
          liveness_probe {
            http_get {
              path = "/minio/health/live"
              port = 9000
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            http_get {
              path = "/minio/health/ready"
              port = 9000
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
        
        volume {
          name = "minio-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.minio_data[0].metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_persistent_volume_claim" "minio_data" {
  count = var.ragflow_minio_config.enabled && var.ragflow_minio_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-minio-data"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "minio"
    })
    annotations = local.annotations
  }
  
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.ragflow_storage_size
      }
    }
    storage_class_name = var.ragflow_storage_class
  }
}

resource "kubernetes_service" "minio" {
  count = var.ragflow_minio_config.enabled && var.ragflow_minio_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-minio"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "minio"
    })
    annotations = local.annotations
  }
  
  spec {
    selector = {
      app = "${var.name}-minio"
    }
    
    port {
      port        = 9000
      target_port = 9000
      name        = "api"
    }
    
    port {
      port        = 9001
      target_port = 9001
      name        = "console"
    }
  }
}

resource "kubernetes_deployment" "redis" {
  count = var.ragflow_redis_config.enabled && var.ragflow_redis_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-redis"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "redis"
    })
    annotations = local.annotations
  }
  
  spec {
    replicas = var.enable_high_availability ? var.ragflow_redis_config.replicas : 1
    
    selector {
      match_labels = {
        app = "${var.name}-redis"
      }
    }
    
    template {
      metadata {
        labels = merge(local.labels, {
          app = "${var.name}-redis"
          "app.kubernetes.io/component" = "redis"
        })
        annotations = local.annotations
      }
      
      spec {
        container {
          name  = "redis"
          image = "redis:latest"
          
          port {
            container_port = 6379
            name           = "redis"
          }
          
          resources {
            limits = {
              cpu    = var.ragflow_redis_config.resources.limits.cpu
              memory = var.ragflow_redis_config.resources.limits.memory
            }
            requests = {
              cpu    = var.ragflow_redis_config.resources.requests.cpu
              memory = var.ragflow_redis_config.resources.requests.memory
            }
          }
          
          volume_mount {
            name       = "redis-data"
            mount_path = "/data"
          }
          
          liveness_probe {
            tcp_socket {
              port = 6379
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            tcp_socket {
              port = 6379
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
        
        volume {
          name = "redis-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.redis_data[0].metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_persistent_volume_claim" "redis_data" {
  count = var.ragflow_redis_config.enabled && var.ragflow_redis_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-redis-data"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "redis"
    })
    annotations = local.annotations
  }
  
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.ragflow_storage_size
      }
    }
    storage_class_name = var.ragflow_storage_class
  }
}

resource "kubernetes_service" "redis" {
  count = var.ragflow_redis_config.enabled && var.ragflow_redis_config.url == "" ? 1 : 0
  
  metadata {
    name      = "${var.name}-redis"
    namespace = local.namespace
    labels = merge(local.labels, {
      "app.kubernetes.io/component" = "redis"
    })
    annotations = local.annotations
  }
  
  spec {
    selector = {
      app = "${var.name}-redis"
    }
    
    port {
      port        = 6379
      target_port = 6379
      name        = "redis"
    }
  }
}

resource "kubernetes_deployment" "ragflow" {
  metadata {
    name      = var.name
    namespace = local.namespace
    labels    = local.labels
    annotations = local.annotations
  }
  
  spec {
    replicas = var.enable_high_availability ? var.ragflow_replicas : 1
    
    selector {
      match_labels = {
        app = var.name
      }
    }
    
    template {
      metadata {
        labels = merge(local.labels, {
          app = var.name
        })
        annotations = merge(local.annotations, {
          "prometheus.io/scrape" = var.enable_prometheus_integration ? "true" : "false"
          "prometheus.io/port"   = var.enable_prometheus_integration ? "8000" : ""
        })
      }
      
      spec {
        container {
          name  = "ragflow"
          image = "${var.ragflow_image}:${var.ragflow_image_tag}"
          
          port {
            container_port = var.ragflow_api_config.port
            name           = "http"
          }
          
          env {
            name  = "RAGFLOW_CONFIG_PATH"
            value = "/etc/ragflow/config.yaml"
          }
          
          env {
            name = "RAGFLOW_API_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.ragflow_api_key.metadata[0].name
                key  = "api-key"
              }
            }
          }
          
          env {
            name  = "QDRANT_URL"
            value = var.ragflow_qdrant_config.url != "" ? var.ragflow_qdrant_config.url : var.ragflow_qdrant_config.enabled ? "http://${kubernetes_service.qdrant[0].metadata[0].name}.${local.namespace}.svc.cluster.local:6333" : ""
          }
          
          env {
            name  = "QDRANT_API_KEY"
            value = var.ragflow_qdrant_config.api_key
          }
          
          env {
            name  = "QDRANT_COLLECTION"
            value = var.ragflow_qdrant_config.collection
          }
          
          env {
            name  = "MINIO_URL"
            value = var.ragflow_minio_config.url != "" ? var.ragflow_minio_config.url : var.ragflow_minio_config.enabled ? "http://${kubernetes_service.minio[0].metadata[0].name}.${local.namespace}.svc.cluster.local:9000" : ""
          }
          
          env {
            name  = "MINIO_ACCESS_KEY"
            value = var.ragflow_minio_config.access_key != "" ? var.ragflow_minio_config.access_key : "minioadmin"
          }
          
          env {
            name  = "MINIO_SECRET_KEY"
            value = var.ragflow_minio_config.secret_key != "" ? var.ragflow_minio_config.secret_key : "minioadmin"
          }
          
          env {
            name  = "MINIO_BUCKET"
            value = var.ragflow_minio_config.bucket
          }
          
          env {
            name  = "REDIS_URL"
            value = var.ragflow_redis_config.url != "" ? var.ragflow_redis_config.url : var.ragflow_redis_config.enabled ? "redis://${kubernetes_service.redis[0].metadata[0].name}.${local.namespace}.svc.cluster.local:6379" : ""
          }
          
          env {
            name  = "REDIS_PASSWORD"
            value = var.ragflow_redis_config.password
          }
          
          env {
            name  = "ENABLE_PROMETHEUS"
            value = var.enable_prometheus_integration ? "true" : "false"
          }
          
          env {
            name  = "ENABLE_JAEGER"
            value = var.enable_jaeger_integration ? "true" : "false"
          }
          
          env {
            name  = "ENABLE_OPENTELEMETRY"
            value = var.enable_opentelemetry_integration ? "true" : "false"
          }
          
          env {
            name  = "ENABLE_LOKI"
            value = var.enable_loki_integration ? "true" : "false"
          }
          
          env {
            name  = "ENABLE_VECTOR"
            value = var.enable_vector_integration ? "true" : "false"
          }
          
          resources {
            limits = {
              cpu    = var.ragflow_resources.limits.cpu
              memory = var.ragflow_resources.limits.memory
            }
            requests = {
              cpu    = var.ragflow_resources.requests.cpu
              memory = var.ragflow_resources.requests.memory
            }
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/ragflow"
          }
          
          liveness_probe {
            http_get {
              path = "/health"
              port = var.ragflow_api_config.port
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            http_get {
              path = "/health"
              port = var.ragflow_api_config.port
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.ragflow_config.metadata[0].name
          }
        }
        
        service_account_name = "ragflow"
      }
    }
  }
}

resource "kubernetes_service" "ragflow" {
  metadata {
    name      = var.name
    namespace = local.namespace
    labels    = local.labels
    annotations = local.annotations
  }
  
  spec {
    selector = {
      app = var.name
    }
    
    port {
      port        = var.ragflow_api_config.port
      target_port = var.ragflow_api_config.port
      name        = "http"
    }
  }
}

resource "kubernetes_ingress_v1" "ragflow" {
  count = var.enable_ingress ? 1 : 0
  
  metadata {
    name      = var.name
    namespace = local.namespace
    labels    = local.labels
    annotations = merge(local.annotations, {
      "kubernetes.io/ingress.class" = var.ingress_class
    })
  }
  
  spec {
    rule {
      host = var.ingress_domain
      
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          
          backend {
            service {
              name = kubernetes_service.ragflow.metadata[0].name
              port {
                number = var.ragflow_api_config.port
              }
            }
          }
        }
      }
    }
    
    dynamic "tls" {
      for_each = var.ingress_tls_enabled ? [1] : []
      
      content {
        hosts       = [var.ingress_domain]
        secret_name = var.ingress_tls_secret
      }
    }
  }
}

resource "kubernetes_manifest" "ragflow_service_monitor" {
  count = var.enable_prometheus_integration ? 1 : 0
  
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"
    metadata = {
      name      = var.name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      selector = {
        matchLabels = {
          app = var.name
        }
      }
      endpoints = [
        {
          port     = "http"
          interval = "15s"
          path     = "/metrics"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "ragflow_kata_integration" {
  count = var.enable_kata_container_integration ? 1 : 0
  
  manifest = {
    apiVersion = "apps/v1"
    kind       = "Deployment"
    metadata = {
      name      = "${var.name}-kata-integration"
      namespace = local.namespace
      labels = merge(local.labels, {
        "app.kubernetes.io/component" = "kata-integration"
      })
    }
    spec = {
      replicas = 1
      selector = {
        matchLabels = {
          app = "${var.name}-kata-integration"
        }
      }
      template = {
        metadata = {
          labels = merge(local.labels, {
            app = "${var.name}-kata-integration"
            "app.kubernetes.io/component" = "kata-integration"
          })
        }
        spec = {
          runtimeClassName = "kata-containers"
          containers = [
            {
              name  = "kata-integration"
              image = "ubuntu:22.04"
              command = [
                "/bin/bash",
                "-c",
                "apt-get update && apt-get install -y curl python3 python3-pip && pip3 install ragflow-client && sleep infinity"
              ]
              env = [
                {
                  name  = "RAGFLOW_API_URL"
                  value = "http://${kubernetes_service.ragflow.metadata[0].name}.${local.namespace}.svc.cluster.local:${var.ragflow_api_config.port}"
                },
                {
                  name = "RAGFLOW_API_KEY"
                  valueFrom = {
                    secretKeyRef = {
                      name = kubernetes_secret.ragflow_api_key.metadata[0].name
                      key  = "api-key"
                    }
                  }
                }
              ]
              resources = {
                limits = {
                  cpu    = "500m"
                  memory = "1Gi"
                }
                requests = {
                  cpu    = "250m"
                  memory = "512Mi"
                }
              }
            }
          ]
        }
      }
    }
  }
}

resource "null_resource" "create_k8s_directory" {
  provisioner "local-exec" {
    command = "mkdir -p ${path.module}/../../k8s/ragflow"
  }
}

resource "local_file" "ragflow_kubernetes_yaml" {
  depends_on = [null_resource.create_k8s_directory]
  
  filename = "${path.module}/../../k8s/ragflow/ragflow.yaml"
  content  = file("${path.module}/../../k8s/ragflow/ragflow.yaml")
}
