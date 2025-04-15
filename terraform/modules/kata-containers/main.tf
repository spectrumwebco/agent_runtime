resource "kubernetes_namespace" "kata_containers" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "kata-containers"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_config_map" "kata_runtime_config" {
  metadata {
    name      = "kata-runtime-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "kata-containers"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      runtime:
        kata_install_mode: "host"
        rust_log: "info"
      components:
        runtime_rs: true
        mem_agent: true
    EOT
  }
}

resource "kubernetes_daemon_set" "kata_containers_runtime" {
  metadata {
    name      = "kata-containers-runtime"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "kata-containers"
      "app.kubernetes.io/component"  = "runtime"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "kata-containers"
        "app.kubernetes.io/component" = "runtime"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "kata-containers"
          "app.kubernetes.io/component" = "runtime"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "kata-runtime"
          image = "ghcr.io/kata-containers/kata-containers:latest"
          
          security_context {
            privileged = true
          }
          
          env {
            name  = "KATA_INSTALL_MODE"
            value = "host"
          }
          
          env {
            name  = "RUST_LOG"
            value = "info"
          }
          
          volume_mount {
            name       = "host"
            mount_path = "/host"
          }
          
          volume_mount {
            name       = "runtime-data"
            mount_path = "/var/run/kata-containers"
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
        }
        
        volume {
          name = "host"
          host_path {
            path = "/"
            type = "Directory"
          }
        }
        
        volume {
          name = "runtime-data"
          host_path {
            path = "/var/run/kata-containers"
            type = "DirectoryOrCreate"
          }
        }
        
        node_selector = var.node_selector
      }
    }
  }
}

resource "kubernetes_deployment" "kata_runtime_components" {
  metadata {
    name      = "kata-runtime-components"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "kata-containers"
      "app.kubernetes.io/component"  = "components"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "kata-containers"
        "app.kubernetes.io/component" = "components"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "kata-containers"
          "app.kubernetes.io/component" = "components"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        container {
          name  = "runtime-rs"
          image = "ghcr.io/kata-containers/runtime-rs:latest"
          
          env {
            name  = "RUST_LOG"
            value = "info"
          }
          
          volume_mount {
            name       = "runtime-data"
            mount_path = "/var/run/kata-containers"
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
        }
        
        container {
          name  = "mem-agent"
          image = "ghcr.io/kata-containers/mem-agent:latest"
          
          env {
            name  = "RUST_LOG"
            value = "info"
          }
          
          volume_mount {
            name       = "runtime-data"
            mount_path = "/var/run/kata-containers"
          }
          
          resources {
            requests = {
              memory = "128Mi"
              cpu    = "100m"
            }
            limits = {
              memory = "256Mi"
              cpu    = "250m"
            }
          }
        }
        
        volume {
          name = "runtime-data"
          host_path {
            path = "/var/run/kata-containers"
            type = "DirectoryOrCreate"
          }
        }
      }
    }
  }
}

resource "kubernetes_runtime_class" "kata_containers" {
  metadata {
    name = "kata-containers"
  }

  handler = "kata"
}

resource "kubernetes_deployment" "agent_runtime_sandbox" {
  metadata {
    name      = "agent-runtime-sandbox"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/component"  = "sandbox"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "agent-runtime"
        "app.kubernetes.io/component" = "sandbox"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "agent-runtime"
          "app.kubernetes.io/component" = "sandbox"
          "app.kubernetes.io/part-of"   = "agent-runtime"
        }
      }

      spec {
        runtime_class_name = "kata-containers"
        
        container {
          name  = "sandbox"
          image = "ghcr.io/spectrumwebco/agent-runtime-sandbox:latest"
          
          env {
            name  = "LIBRECHAT_CODE_API_KEY"
            value_from {
              secret_key_ref {
                name = "agent-runtime-secrets"
                key  = "librechat_code_api_key"
              }
            }
          }
          
          port {
            container_port = 8081
            name           = "http"
          }
          
          volume_mount {
            name       = "workspace"
            mount_path = "/workspace"
          }
          
          resources {
            requests = {
              memory = "1Gi"
              cpu    = "500m"
            }
            limits = {
              memory = "4Gi"
              cpu    = "2"
            }
          }
        }
        
        volume {
          name = "workspace"
          persistent_volume_claim {
            claim_name = "agent-runtime-workspace"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "agent_runtime_sandbox" {
  metadata {
    name      = "agent-runtime-sandbox"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/component"  = "sandbox"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "agent-runtime"
      "app.kubernetes.io/component" = "sandbox"
    }
    
    port {
      port        = 8081
      target_port = "http"
      name        = "http"
    }
  }
}

resource "kubernetes_persistent_volume_claim" "agent_runtime_workspace" {
  metadata {
    name      = "agent-runtime-workspace"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = "10Gi"
      }
    }
  }
}

resource "kubernetes_secret" "agent_runtime_secrets" {
  metadata {
    name      = "agent-runtime-secrets"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  type = "Opaque"
  
  data = {
    "librechat_code_api_key" = var.librechat_code_api_key
    "rdp_password"           = var.rdp_password
    "ssh_key"                = var.ssh_key
  }
}
