variable "namespace" {
  description = "Kubernetes namespace for Neovim deployment"
  type        = string
  default     = "neovim"
}

variable "replicas" {
  description = "Number of Neovim replicas"
  type        = number
  default     = 1
}

variable "image" {
  description = "Neovim container image"
  type        = string
  default     = "neovim/neovim:latest"
}

variable "supabase_url" {
  description = "Supabase URL for state persistence"
  type        = string
}

variable "supabase_key" {
  description = "Supabase key for state persistence"
  type        = string
  sensitive   = true
}

variable "storage_size" {
  description = "Size of persistent volume for Neovim data"
  type        = string
  default     = "1Gi"
}

variable "kata_runtime_class" {
  description = "Kata Containers runtime class"
  type        = string
  default     = "kata"
}

variable "enable_kata" {
  description = "Whether to enable Kata Containers for Neovim"
  type        = bool
  default     = true
}

variable "resource_limits" {
  description = "Resource limits for Neovim container"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "200m"
    memory = "512Mi"
  }
}

variable "resource_requests" {
  description = "Resource requests for Neovim container"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "100m"
    memory = "256Mi"
  }
}

locals {
  labels = {
    "app.kubernetes.io/name"      = "neovim"
    "app.kubernetes.io/part-of"   = "agent-runtime"
    "app.kubernetes.io/component" = "editor"
    "app.kubernetes.io/managed-by" = "terraform"
  }
}

resource "kubernetes_namespace" "neovim" {
  metadata {
    name = var.namespace
    labels = local.labels
  }
}

resource "kubernetes_config_map" "neovim_config" {
  metadata {
    name      = "neovim-config"
    namespace = kubernetes_namespace.neovim.metadata[0].name
    labels    = local.labels
  }

  data = {
    "supabase_url" = var.supabase_url
    "init.lua"     = file("${path.module}/templates/init.lua.tpl")
    "db_integration.lua" = file("${path.module}/templates/db_integration.lua.tpl")
  }
}

resource "kubernetes_secret" "neovim_secrets" {
  metadata {
    name      = "neovim-secrets"
    namespace = kubernetes_namespace.neovim.metadata[0].name
    labels    = local.labels
  }

  data = {
    "supabase_key" = var.supabase_key
  }

  type = "Opaque"
}

resource "kubernetes_persistent_volume_claim" "neovim_data" {
  metadata {
    name      = "neovim-data-pvc"
    namespace = kubernetes_namespace.neovim.metadata[0].name
    labels    = local.labels
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.storage_size
      }
    }
  }
}

resource "kubernetes_deployment" "neovim" {
  metadata {
    name      = "neovim"
    namespace = kubernetes_namespace.neovim.metadata[0].name
    labels    = local.labels
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "neovim"
      }
    }

    template {
      metadata {
        labels = merge(local.labels, {
          app = "neovim"
        })
      }

      spec {
        runtime_class_name = var.enable_kata ? var.kata_runtime_class : null

        container {
          name  = "neovim"
          image = var.image

          port {
            container_port = 8090
            name           = "http"
          }

          env {
            name  = "NEOVIM_API_BASE"
            value = "http://localhost:8090/neovim"
          }

          env {
            name = "SUPABASE_URL"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.neovim_config.metadata[0].name
                key  = "supabase_url"
              }
            }
          }

          env {
            name = "SUPABASE_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.neovim_secrets.metadata[0].name
                key  = "supabase_key"
              }
            }
          }

          volume_mount {
            name       = "neovim-config"
            mount_path = "/root/.config/nvim"
          }

          volume_mount {
            name       = "neovim-data"
            mount_path = "/root/.local/share/nvim"
          }

          resources {
            limits = {
              cpu    = var.resource_limits.cpu
              memory = var.resource_limits.memory
            }
            requests = {
              cpu    = var.resource_requests.cpu
              memory = var.resource_requests.memory
            }
          }
        }

        volume {
          name = "neovim-config"
          config_map {
            name = kubernetes_config_map.neovim_config.metadata[0].name
          }
        }

        volume {
          name = "neovim-data"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.neovim_data.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "neovim" {
  metadata {
    name      = "neovim"
    namespace = kubernetes_namespace.neovim.metadata[0].name
    labels    = local.labels
  }

  spec {
    selector = {
      app = "neovim"
    }

    port {
      port        = 8090
      target_port = 8090
      name        = "http"
    }

    type = "ClusterIP"
  }
}
