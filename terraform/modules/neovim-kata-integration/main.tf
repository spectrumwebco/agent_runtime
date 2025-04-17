variable "namespace" {
  description = "Kubernetes namespace for Neovim Kata integration"
  type        = string
  default     = "neovim"
}

variable "kata_runtime_class" {
  description = "Kata Containers runtime class"
  type        = string
  default     = "kata"
}

variable "kata_version" {
  description = "Kata Containers version"
  type        = string
  default     = "2.4.0"
}

variable "enable_runtime_rs" {
  description = "Enable Rust-based runtime for Kata Containers"
  type        = bool
  default     = true
}

variable "enable_mem_agent" {
  description = "Enable Rust-based memory management agent for Kata Containers"
  type        = bool
  default     = true
}

variable "resource_limits" {
  description = "Resource limits for Kata Containers"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "500m"
    memory = "1Gi"
  }
}

resource "kubernetes_runtime_class" "kata" {
  metadata {
    name = var.kata_runtime_class
  }

  handler = "kata"
}

resource "kubernetes_config_map" "kata_config" {
  metadata {
    name      = "kata-config"
    namespace = var.namespace
  }

  data = {
    "configuration.toml" = templatefile("${path.module}/templates/kata-config.toml.tpl", {
      enable_runtime_rs = var.enable_runtime_rs
      enable_mem_agent  = var.enable_mem_agent
    })
  }
}

resource "kubernetes_daemonset" "kata_containers" {
  metadata {
    name      = "kata-containers"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"      = "kata-containers"
      "app.kubernetes.io/part-of"   = "agent-runtime"
      "app.kubernetes.io/component" = "runtime"
    }
  }

  spec {
    selector {
      match_labels = {
        name = "kata-containers"
      }
    }

    template {
      metadata {
        labels = {
          name = "kata-containers"
        }
      }

      spec {
        container {
          name  = "kata-containers"
          image = "katadocker/kata-deploy:${var.kata_version}"

          resources {
            limits = {
              cpu    = var.resource_limits.cpu
              memory = var.resource_limits.memory
            }
          }

          volume_mount {
            name       = "kata-config"
            mount_path = "/opt/kata/share/defaults/kata-containers"
          }

          volume_mount {
            name       = "kata-artifacts"
            mount_path = "/opt/kata"
          }

          volume_mount {
            name       = "kata-containers"
            mount_path = "/var/run/kata-containers"
          }

          volume_mount {
            name       = "crio-conf"
            mount_path = "/etc/crio/crio.conf.d"
          }

          security_context {
            privileged = true
          }
        }

        volume {
          name = "kata-config"
          config_map {
            name = kubernetes_config_map.kata_config.metadata[0].name
          }
        }

        volume {
          name = "kata-artifacts"
          host_path {
            path = "/opt/kata"
          }
        }

        volume {
          name = "kata-containers"
          host_path {
            path = "/var/run/kata-containers"
          }
        }

        volume {
          name = "crio-conf"
          host_path {
            path = "/etc/crio/crio.conf.d"
          }
        }
      }
    }
  }
}

output "kata_runtime_class" {
  description = "Kata Containers runtime class"
  value       = kubernetes_runtime_class.kata.metadata[0].name
}

output "kata_config_map" {
  description = "Kata Containers config map"
  value       = kubernetes_config_map.kata_config.metadata[0].name
}
