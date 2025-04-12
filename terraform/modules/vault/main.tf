provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "vault" {
  metadata {
    name = var.vault_namespace
    labels = {
      "app.kubernetes.io/name" = "vault"
      "app.kubernetes.io/instance" = "vault"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_secret" "vault_token" {
  metadata {
    name = "vault-token"
    namespace = kubernetes_namespace.vault.metadata[0].name
  }

  data = {
    token = var.vault_token
  }

  type = "Opaque"
}

resource "kubernetes_service_account" "vault_auth" {
  metadata {
    name = "vault-auth"
    namespace = kubernetes_namespace.vault.metadata[0].name
  }
}

resource "kubernetes_cluster_role_binding" "vault_auth_delegator" {
  metadata {
    name = "vault-auth-delegator"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind = "ClusterRole"
    name = "system:auth-delegator"
  }

  subject {
    kind = "ServiceAccount"
    name = kubernetes_service_account.vault_auth.metadata[0].name
    namespace = kubernetes_namespace.vault.metadata[0].name
  }
}

resource "kubernetes_deployment" "vault" {
  metadata {
    name = "vault"
    namespace = kubernetes_namespace.vault.metadata[0].name
    labels = {
      app = "vault"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "vault"
      }
    }

    template {
      metadata {
        labels = {
          app = "vault"
        }
      }

      spec {
        container {
          name = "vault"
          image = "hashicorp/vault:${var.vault_version}"
          
          port {
            container_port = 8200
          }

          env {
            name = "VAULT_DEV_ROOT_TOKEN_ID"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.vault_token.metadata[0].name
                key = "token"
              }
            }
          }

          env {
            name = "VAULT_DEV_LISTEN_ADDRESS"
            value = "0.0.0.0:8200"
          }

          env {
            name = "VAULT_ADDR"
            value = "http://127.0.0.1:8200"
          }

          args = [
            "server",
            "-dev",
            "-dev-root-token-id=$(VAULT_DEV_ROOT_TOKEN_ID)",
            "-dev-listen-address=$(VAULT_DEV_LISTEN_ADDRESS)"
          ]

          resources {
            limits = {
              cpu = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu = "250m"
              memory = "256Mi"
            }
          }

          readiness_probe {
            http_get {
              path = "/v1/sys/health"
              port = 8200
            }
            initial_delay_seconds = 5
            period_seconds = 10
          }

          liveness_probe {
            http_get {
              path = "/v1/sys/health"
              port = 8200
            }
            initial_delay_seconds = 10
            period_seconds = 15
          }

          security_context {
            capabilities {
              add = ["IPC_LOCK"]
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "vault" {
  metadata {
    name = "vault"
    namespace = kubernetes_namespace.vault.metadata[0].name
    labels = {
      app = "vault"
    }
  }

  spec {
    selector = {
      app = "vault"
    }

    port {
      name = "vault"
      port = 8200
      target_port = 8200
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_config_map" "vault_config" {
  metadata {
    name = "vault-config"
    namespace = kubernetes_namespace.vault.metadata[0].name
  }

  data = {
    "vault-init.sh" = file("${path.module}/scripts/vault-init.sh")
  }
}

resource "kubernetes_job" "vault_init" {
  metadata {
    name = "vault-init"
    namespace = kubernetes_namespace.vault.metadata[0].name
  }

  spec {
    template {
      metadata {}

      spec {
        service_account_name = kubernetes_service_account.vault_auth.metadata[0].name
        
        container {
          name = "vault-init"
          image = "hashicorp/vault:${var.vault_version}"
          
          command = ["/bin/sh", "-c"]
          args = [
            "cp /vault-config/vault-init.sh /tmp/vault-init.sh && chmod +x /tmp/vault-init.sh && /tmp/vault-init.sh"
          ]

          env {
            name = "VAULT_ADDR"
            value = "http://vault.${kubernetes_namespace.vault.metadata[0].name}.svc.cluster.local:8200"
          }

          env {
            name = "VAULT_TOKEN"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.vault_token.metadata[0].name
                key = "token"
              }
            }
          }

          volume_mount {
            name = "vault-config"
            mount_path = "/vault-config"
          }
        }

        volume {
          name = "vault-config"
          config_map {
            name = kubernetes_config_map.vault_config.metadata[0].name
          }
        }

        restart_policy = "OnFailure"
      }
    }
  }

  depends_on = [
    kubernetes_deployment.vault,
    kubernetes_service.vault
  ]
}

resource "kubernetes_service_account" "vault_agent_injector" {
  metadata {
    name = "vault-agent-injector"
    namespace = kubernetes_namespace.vault.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "vault_agent_injector" {
  metadata {
    name = "vault-agent-injector"
  }

  rule {
    api_groups = ["admissionregistration.k8s.io"]
    resources = ["mutatingwebhookconfigurations"]
    verbs = ["get", "list", "watch", "patch"]
  }

  rule {
    api_groups = [""]
    resources = ["namespaces"]
    verbs = ["get", "list", "watch"]
  }

  rule {
    api_groups = [""]
    resources = ["pods"]
    verbs = ["get", "list", "watch", "patch"]
  }

  rule {
    api_groups = [""]
    resources = ["secrets"]
    verbs = ["get", "list", "watch", "create", "update", "delete"]
  }
}

resource "kubernetes_cluster_role_binding" "vault_agent_injector" {
  metadata {
    name = "vault-agent-injector"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind = "ClusterRole"
    name = kubernetes_cluster_role.vault_agent_injector.metadata[0].name
  }

  subject {
    kind = "ServiceAccount"
    name = kubernetes_service_account.vault_agent_injector.metadata[0].name
    namespace = kubernetes_namespace.vault.metadata[0].name
  }
}

resource "kubernetes_deployment" "vault_agent_injector" {
  metadata {
    name = "vault-agent-injector"
    namespace = kubernetes_namespace.vault.metadata[0].name
    labels = {
      app = "vault-agent-injector"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "vault-agent-injector"
      }
    }

    template {
      metadata {
        labels = {
          app = "vault-agent-injector"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.vault_agent_injector.metadata[0].name
        
        container {
          name = "vault-agent-injector"
          image = "hashicorp/vault-k8s:${var.vault_k8s_version}"
          
          args = [
            "agent-inject",
            "-tls-auto=vault-agent-injector",
            "-tls-auto-hosts=vault-agent-injector,vault-agent-injector.${kubernetes_namespace.vault.metadata[0].name},vault-agent-injector.${kubernetes_namespace.vault.metadata[0].name}.svc",
            "-log-level=info",
            "-log-format=standard"
          ]

          env {
            name = "AGENT_INJECT_LISTEN"
            value = ":8080"
          }

          env {
            name = "AGENT_INJECT_VAULT_ADDR"
            value = "http://vault.${kubernetes_namespace.vault.metadata[0].name}.svc.cluster.local:8200"
          }

          env {
            name = "AGENT_INJECT_VAULT_AUTH_PATH"
            value = "auth/kubernetes"
          }

          env {
            name = "AGENT_INJECT_LOG_LEVEL"
            value = "info"
          }

          env {
            name = "AGENT_INJECT_LOG_FORMAT"
            value = "standard"
          }

          env {
            name = "AGENT_INJECT_REVOKE_ON_SHUTDOWN"
            value = "false"
          }

          env {
            name = "AGENT_INJECT_CPU_REQUEST"
            value = "250m"
          }

          env {
            name = "AGENT_INJECT_CPU_LIMIT"
            value = "500m"
          }

          env {
            name = "AGENT_INJECT_MEM_REQUEST"
            value = "64Mi"
          }

          env {
            name = "AGENT_INJECT_MEM_LIMIT"
            value = "128Mi"
          }

          env {
            name = "AGENT_INJECT_DEFAULT_TEMPLATE"
            value = "map"
          }

          resources {
            limits = {
              cpu = "500m"
              memory = "128Mi"
            }
            requests = {
              cpu = "250m"
              memory = "64Mi"
            }
          }

          port {
            container_port = 8080
          }

          readiness_probe {
            http_get {
              path = "/health/ready"
              port = 8080
              scheme = "HTTPS"
            }
            initial_delay_seconds = 5
            period_seconds = 10
          }

          liveness_probe {
            http_get {
              path = "/health/ready"
              port = 8080
              scheme = "HTTPS"
            }
            initial_delay_seconds = 10
            period_seconds = 15
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "vault_agent_injector" {
  metadata {
    name = "vault-agent-injector"
    namespace = kubernetes_namespace.vault.metadata[0].name
    labels = {
      app = "vault-agent-injector"
    }
  }

  spec {
    selector = {
      app = "vault-agent-injector"
    }

    port {
      name = "https"
      port = 443
      target_port = 8080
    }

    type = "ClusterIP"
  }
}
