
resource "kubernetes_namespace" "mcp" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = "agent-runtime-system"
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

locals {
  namespace = "agent-runtime-system"
  
  common_labels = {
    "app"                         = "agent-runtime-mcp-client"
    "agent-runtime/component"     = "mcp-client"
    "app.kubernetes.io/name"      = "agent-runtime"
    "app.kubernetes.io/part-of"   = "agent-runtime"
    "app.kubernetes.io/managed-by" = "terraform"
  }

resource "kubernetes_config_map" "mcp_config" {
  metadata {
    name      = "agent-runtime-config"
    namespace = local.namespace
    labels    = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      mcp:
        client:
          enabled: true
          port: 8082
          host_url: "http://agent-runtime-mcp-host:8080"
          server_url: "http://agent-runtime-sandbox:8081"
        container_runtimes:
          - lxc
          - podman
          - docker
          - kata
        handlers:
          - creation
          - deletion
          - maintenance
          - monitoring
          - insights
      integrations:
        librechat:
          enabled: true
          api_key_env: "LIBRECHAT_CODE_API_KEY"
    EOT
  }
}

resource "kubernetes_secret" "mcp_secrets" {
  metadata {
    name      = "agent-runtime-secrets"
    namespace = local.namespace
    labels    = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  type = "Opaque"
  
  string_data = {
    "librechat_code_api_key" = var.librechat_code_api_key
  }
}

resource "kubernetes_config_map" "mcp_client_config" {
  metadata {
    name      = "mcp-client"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "mcp-client" = var.mcp_client_config
  }
}

resource "kubernetes_deployment" "mcp_client" {
  metadata {
    name      = "agent-runtime-mcp-client"
    namespace = local.namespace
    labels    = {
      "app"                         = "agent-runtime-mcp-client"
      "agent-runtime/component"     = "mcp-client"
      "app.kubernetes.io/name"      = "agent-runtime"
      "app.kubernetes.io/part-of"   = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        "app" = "agent-runtime-mcp-client"
      }
    }

    template {
      metadata {
        labels = {
          "app"                         = "agent-runtime-mcp-client"
          "agent-runtime/component"     = "mcp-client"
          "app.kubernetes.io/name"      = "agent-runtime"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        containers {
          name  = "mcp-client"
          image = "${var.container_registry}/agent-runtime:latest"
          
          command = ["./agent-runtime"]
          args    = ["--mcp-client"]
          
          ports {
            container_port = 8082
            name          = "mcp-client"
          }
          
          env {
            name  = "MCP_HOST_URL"
            value = "http://agent-runtime-mcp-host:8080"
          }
          
          env {
            name  = "MCP_SERVER_URL"
            value = "http://agent-runtime-sandbox:8081"
          }
          
          env {
            name = "LIBRECHAT_CODE_API_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.mcp_secrets.metadata[0].name
                key  = "librechat_code_api_key"
              }
            }
          }
          
          resources {
            requests = {
              memory = "256Mi"
              cpu    = "250m"
            }
            limits = {
              memory = "512Mi"
              cpu    = "500m"
            }
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/agent-runtime/config"
          }
          
          volume_mount {
            name       = "shared-data"
            mount_path = "/var/lib/agent-runtime/data"
          }
        }
        
        volumes {
          name = "config"
          config_map {
            name = "agent-runtime-config"
          }
        }
        
        volumes {
          name = "shared-data"
          empty_dir {}
        }
      }
    }
  }
}

resource "kubernetes_service" "mcp_client" {
  metadata {
    name      = "agent-runtime-mcp-client"
    namespace = local.namespace
    labels    = {
      "app"                         = "agent-runtime-mcp-client"
      "agent-runtime/component"     = "mcp-client"
      "app.kubernetes.io/name"      = "agent-runtime"
      "app.kubernetes.io/part-of"   = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app" = "agent-runtime-mcp-client"
    }
    
    ports {
      port        = 8082
      target_port = "mcp-client"
      name        = "mcp-client"
    }
  }
}
