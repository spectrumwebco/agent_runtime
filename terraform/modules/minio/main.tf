/**
 * # MinIO Terraform Module
 *
 * This module deploys MinIO for artifact storage in the ML infrastructure.
 */

resource "kubernetes_namespace" "ml_infrastructure" {
  count = var.create_namespace ? 1 : 0
  
  metadata {
    name = var.namespace
    
    labels = {
      "app.kubernetes.io/part-of" = "ml-infrastructure"
      "app.kubernetes.io/component" = "storage"
    }
  }
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
  }
  
  data = {
    accessKey = var.minio_access_key
    secretKey = var.minio_secret_key
  }
  
  depends_on = [kubernetes_namespace.ml_infrastructure]
}

resource "kubernetes_config_map" "minio_config" {
  metadata {
    name      = "minio-config"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
  }
  
  data = {
    MINIO_BROWSER_REDIRECT_URL = var.minio_console_url
    MINIO_PROMETHEUS_AUTH_TYPE = "public"
    MINIO_PROMETHEUS_URL       = var.prometheus_url
    MINIO_PROMETHEUS_JOB_ID    = "minio"
    MINIO_REGION               = var.region
    MINIO_DOMAIN               = var.minio_domain
    MINIO_STORAGE_CLASS_STANDARD = "EC:2"
    MINIO_STORAGE_CLASS_RRS    = "EC:1"
  }
  
  depends_on = [kubernetes_namespace.ml_infrastructure]
}

resource "kubernetes_persistent_volume_claim" "minio_pvc" {
  metadata {
    name      = "minio-pvc"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
  }
  
  spec {
    access_modes = ["ReadWriteOnce"]
    
    resources {
      requests = {
        storage = var.storage_size
      }
    }
    
    storage_class_name = var.storage_class
  }
  
  depends_on = [kubernetes_namespace.ml_infrastructure]
}

resource "kubernetes_deployment" "minio" {
  metadata {
    name      = "minio"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
  }
  
  spec {
    selector {
      match_labels = {
        app = "minio"
      }
    }
    
    strategy {
      type = "Recreate"
    }
    
    template {
      metadata {
        labels = {
          app       = "minio"
          component = "artifact-storage"
        }
      }
      
      spec {
        container {
          name  = "minio"
          image = var.minio_image
          
          args = [
            "server",
            "/data",
            "--console-address",
            ":9001"
          ]
          
          env {
            name = "MINIO_ROOT_USER"
            
            value_from {
              secret_key_ref {
                name = kubernetes_secret.minio_credentials.metadata[0].name
                key  = "accessKey"
              }
            }
          }
          
          env {
            name = "MINIO_ROOT_PASSWORD"
            
            value_from {
              secret_key_ref {
                name = kubernetes_secret.minio_credentials.metadata[0].name
                key  = "secretKey"
              }
            }
          }
          
          dynamic "env" {
            for_each = kubernetes_config_map.minio_config.data
            
            content {
              name  = env.key
              value = env.value
            }
          }
          
          port {
            container_port = 9000
            name           = "api"
          }
          
          port {
            container_port = 9001
            name           = "console"
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/data"
          }
          
          resources {
            requests = {
              memory = var.memory_request
              cpu    = var.cpu_request
            }
            
            limits = {
              memory = var.memory_limit
              cpu    = var.cpu_limit
            }
          }
          
          liveness_probe {
            http_get {
              path = "/minio/health/live"
              port = "api"
            }
            
            initial_delay_seconds = 120
            period_seconds        = 20
          }
          
          readiness_probe {
            http_get {
              path = "/minio/health/ready"
              port = "api"
            }
            
            initial_delay_seconds = 120
            period_seconds        = 20
          }
        }
        
        volume {
          name = "data"
          
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.minio_pvc.metadata[0].name
          }
        }
      }
    }
  }
  
  depends_on = [
    kubernetes_namespace.ml_infrastructure,
    kubernetes_secret.minio_credentials,
    kubernetes_config_map.minio_config,
    kubernetes_persistent_volume_claim.minio_pvc
  ]
}

resource "kubernetes_service" "minio" {
  metadata {
    name      = "minio"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
  }
  
  spec {
    port {
      port        = 9000
      target_port = 9000
      protocol    = "TCP"
      name        = "api"
    }
    
    port {
      port        = 9001
      target_port = 9001
      protocol    = "TCP"
      name        = "console"
    }
    
    selector = {
      app = "minio"
    }
    
    type = "ClusterIP"
  }
  
  depends_on = [kubernetes_namespace.ml_infrastructure]
}

resource "kubernetes_service" "minio_headless" {
  metadata {
    name      = "minio-headless"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
  }
  
  spec {
    port {
      port = 9000
      name = "api"
    }
    
    port {
      port = 9001
      name = "console"
    }
    
    selector = {
      app = "minio"
    }
    
    cluster_ip = "None"
  }
  
  depends_on = [kubernetes_namespace.ml_infrastructure]
}

resource "kubernetes_ingress_v1" "minio" {
  count = var.create_ingress ? 1 : 0
  
  metadata {
    name      = "minio"
    namespace = var.namespace
    
    labels = {
      app       = "minio"
      component = "artifact-storage"
    }
    
    annotations = {
      "kubernetes.io/ingress.class"                 = "nginx"
      "nginx.ingress.kubernetes.io/ssl-redirect"    = "true"
      "nginx.ingress.kubernetes.io/proxy-body-size" = "0"
      "nginx.ingress.kubernetes.io/proxy-buffering" = "off"
      "nginx.ingress.kubernetes.io/proxy-read-timeout" = "600"
      "nginx.ingress.kubernetes.io/proxy-send-timeout" = "600"
    }
  }
  
  spec {
    rule {
      host = var.minio_domain
      
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          
          backend {
            service {
              name = kubernetes_service.minio.metadata[0].name
              
              port {
                name = "api"
              }
            }
          }
        }
      }
    }
    
    rule {
      host = var.minio_console_domain
      
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          
          backend {
            service {
              name = kubernetes_service.minio.metadata[0].name
              
              port {
                name = "console"
              }
            }
          }
        }
      }
    }
    
    dynamic "tls" {
      for_each = var.tls_secret_name != "" ? [1] : []
      
      content {
        hosts = [
          var.minio_domain,
          var.minio_console_domain
        ]
        
        secret_name = var.tls_secret_name
      }
    }
  }
  
  depends_on = [
    kubernetes_namespace.ml_infrastructure,
    kubernetes_service.minio
  ]
}

resource "null_resource" "create_buckets" {
  count = var.create_buckets ? 1 : 0
  
  provisioner "local-exec" {
    command = <<-EOT
      export MINIO_ENDPOINT=${var.minio_endpoint}
      export MINIO_ACCESS_KEY=${var.minio_access_key}
      export MINIO_SECRET_KEY=${var.minio_secret_key}
      
      if ! command -v mc &> /dev/null; then
        curl -O https://dl.min.io/client/mc/release/linux-amd64/mc
        chmod +x mc
        sudo mv mc /usr/local/bin/
      fi
      
      mc alias set minio "${MINIO_ENDPOINT}" "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}"
      
      for bucket in ${join(" ", var.buckets)}; do
        mc mb --ignore-existing "minio/$bucket"
      done
      
      mc policy set download "minio/mlflow-artifacts"
      mc policy set download "minio/model-registry"
      mc policy set download "minio/training-data"
      mc policy set download "minio/model-serving"
      
      mc version enable "minio/model-registry"
      mc version enable "minio/checkpoints"
    EOT
  }
  
  depends_on = [
    kubernetes_deployment.minio,
    kubernetes_service.minio
  ]
}

resource "null_resource" "configure_lifecycle" {
  count = var.configure_lifecycle ? 1 : 0
  
  provisioner "local-exec" {
    command = <<-EOT
      export MINIO_ENDPOINT=${var.minio_endpoint}
      export MINIO_ACCESS_KEY=${var.minio_access_key}
      export MINIO_SECRET_KEY=${var.minio_secret_key}
      
      if ! command -v mc &> /dev/null; then
        curl -O https://dl.min.io/client/mc/release/linux-amd64/mc
        chmod +x mc
        sudo mv mc /usr/local/bin/
      fi
      
      mc alias set minio "${MINIO_ENDPOINT}" "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}"
      
      cat > /tmp/logs-lifecycle.json << EOF
      {
        "Rules": [
          {
            "ID": "expire-old-logs",
            "Status": "Enabled",
            "Filter": {
              "Prefix": ""
            },
            "Expiration": {
              "Days": 30
            }
          }
        ]
      }
      EOF
      
      mc ilm import minio/logs < /tmp/logs-lifecycle.json
      rm /tmp/logs-lifecycle.json
      
      cat > /tmp/checkpoints-lifecycle.json << EOF
      {
        "Rules": [
          {
            "ID": "expire-old-checkpoints",
            "Status": "Enabled",
            "Filter": {
              "Prefix": ""
            },
            "NoncurrentVersionExpiration": {
              "NoncurrentDays": 90
            },
            "Expiration": {
              "Days": 365
            }
          }
        ]
      }
      EOF
      
      mc ilm import minio/checkpoints < /tmp/checkpoints-lifecycle.json
      rm /tmp/checkpoints-lifecycle.json
      
      cat > /tmp/model-registry-lifecycle.json << EOF
      {
        "Rules": [
          {
            "ID": "expire-old-model-versions",
            "Status": "Enabled",
            "Filter": {
              "Prefix": ""
            },
            "NoncurrentVersionExpiration": {
              "NoncurrentDays": 180
            }
          }
        ]
      }
      EOF
      
      mc ilm import minio/model-registry < /tmp/model-registry-lifecycle.json
      rm /tmp/model-registry-lifecycle.json
    EOT
  }
  
  depends_on = [
    null_resource.create_buckets
  ]
}
