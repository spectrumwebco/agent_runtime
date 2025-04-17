variable "namespace" {
  description = "Kubernetes namespace for terraform-operator deployment"
  type        = string
  default     = "terraform-operator"
}

resource "kubernetes_namespace" "terraform_operator" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/name"      = "terraform-operator"
      "app.kubernetes.io/part-of"   = "agent-runtime"
    }
  }
}

resource "kubernetes_deployment" "terraform_operator" {
  metadata {
    name      = "terraform-operator"
    namespace = kubernetes_namespace.terraform_operator.metadata[0].name
    labels = {
      app                       = "terraform-operator"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "terraform-operator"
      }
    }

    template {
      metadata {
        labels = {
          app = "terraform-operator"
        }
      }

      spec {
        container {
          name  = "terraform-operator"
          image = "galleybytes/terraform-operator:latest"

          resources {
            limits = {
              cpu    = "1"
              memory = "1Gi"
            }
            requests = {
              cpu    = "0.5"
              memory = "500Mi"
            }
          }
        }

        service_account_name = "terraform-operator-sa"
      }
    }
  }
}

resource "kubernetes_service_account" "terraform_operator_sa" {
  metadata {
    name      = "terraform-operator-sa"
    namespace = kubernetes_namespace.terraform_operator.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "terraform_operator_role" {
  metadata {
    name = "terraform-operator-role"
  }

  rule {
    api_groups = [""]
    resources  = ["configmaps", "secrets"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["deployments"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["batch"]
    resources  = ["jobs"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_cluster_role_binding" "terraform_operator_rolebinding" {
  metadata {
    name = "terraform-operator-rolebinding"
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.terraform_operator_sa.metadata[0].name
    namespace = kubernetes_namespace.terraform_operator.metadata[0].name
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.terraform_operator_role.metadata[0].name
  }
}
