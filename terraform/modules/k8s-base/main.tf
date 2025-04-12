resource "kubernetes_namespace" "agent_runtime_system" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = "agent-runtime-system"
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_service_account" "agent_runtime_controller" {
  metadata {
    name      = "agent-runtime-controller"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_cluster_role" "agent_runtime_controller" {
  metadata {
    name = "agent-runtime-controller"
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  rule {
    api_groups = [""]
    resources  = ["pods", "services", "configmaps", "secrets"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "statefulsets", "daemonsets"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["networking.k8s.io"]
    resources  = ["ingresses"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_cluster_role_binding" "agent_runtime_controller" {
  metadata {
    name = "agent-runtime-controller"
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.agent_runtime_controller.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.agent_runtime_controller.metadata[0].name
    namespace = kubernetes_service_account.agent_runtime_controller.metadata[0].namespace
  }
}

resource "kubernetes_config_map" "agent_runtime_config" {
  metadata {
    name      = "agent-runtime-config"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/name"       = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }

  data = {
    "config.yaml" = <<-EOT
      controller:
        workers: 2
        resyncPeriod: 30s
      logging:
        level: info
    EOT
  }
}
