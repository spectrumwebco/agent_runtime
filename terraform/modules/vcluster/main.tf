resource "helm_release" "vcluster" {
  name       = var.vcluster_name
  namespace  = var.vcluster_namespace
  repository = "https://charts.loft.sh"
  chart      = "vcluster"
  version    = "0.15.0"

  set {
    name  = "syncer.extraArgs"
    value = "{--tls-san=${var.vcluster_name}.${var.vcluster_namespace}}"
  }

  set {
    name  = "vcluster.image"
    value = "rancher/k3s:v${var.kubernetes_version}-k3s1"
  }

  set {
    name  = "persistent"
    value = var.persistent
  }

  set {
    name  = "storage.persistence.enabled"
    value = var.persistent
  }

  set {
    name  = "storage.persistence.size"
    value = "10Gi"
  }

  set {
    name  = "distro"
    value = var.distro
  }

  set {
    name  = "sync.nodes.enabled"
    value = "true"
  }

  set {
    name  = "sync.ingresses.enabled"
    value = "true"
  }

  set {
    name  = "isolation.enabled"
    value = "true"
  }

  set {
    name  = "replicas"
    value = "3"
  }

  set {
    name  = "vcluster.resources.requests.cpu"
    value = "500m"
  }

  set {
    name  = "vcluster.resources.requests.memory"
    value = "1Gi"
  }

  set {
    name  = "vcluster.resources.limits.cpu"
    value = "2000m"
  }

  set {
    name  = "vcluster.resources.limits.memory"
    value = "4Gi"
  }
}

resource "kubernetes_deployment" "vnode_integration" {
  metadata {
    name      = "vnode-integration"
    namespace = var.vcluster_namespace
    labels = {
      app = "vnode-integration"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "vnode-integration"
      }
    }

    template {
      metadata {
        labels = {
          app = "vnode-integration"
        }
      }

      spec {
        container {
          image = "ghcr.io/loft-sh/vnode-runtime:0.0.1-alpha.1"
          name  = "vnode-runtime"

          env {
            name  = "VCLUSTER_NAME"
            value = var.vcluster_name
          }

          env {
            name  = "VCLUSTER_NAMESPACE"
            value = var.vcluster_namespace
          }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
        }
      }
    }
  }

  depends_on = [helm_release.vcluster]
}

resource "kubernetes_service_account" "vcluster" {
  metadata {
    name      = "${var.vcluster_name}-sa"
    namespace = var.vcluster_namespace
  }
}

resource "kubernetes_cluster_role" "vcluster" {
  metadata {
    name = "${var.vcluster_name}-role"
  }

  rule {
    api_groups = [""]
    resources  = ["nodes", "namespaces", "pods", "services", "configmaps", "secrets", "serviceaccounts", "persistentvolumes", "persistentvolumeclaims"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "statefulsets", "daemonsets", "replicasets"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["networking.k8s.io"]
    resources  = ["ingresses", "networkpolicies"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["rbac.authorization.k8s.io"]
    resources  = ["roles", "rolebindings", "clusterroles", "clusterrolebindings"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["storage.k8s.io"]
    resources  = ["storageclasses"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["batch"]
    resources  = ["jobs", "cronjobs"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["policy"]
    resources  = ["podsecuritypolicies", "poddisruptionbudgets"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }

  rule {
    api_groups = ["autoscaling"]
    resources  = ["horizontalpodautoscalers"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_cluster_role_binding" "vcluster" {
  metadata {
    name = "${var.vcluster_name}-rolebinding"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.vcluster.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.vcluster.metadata[0].name
    namespace = kubernetes_service_account.vcluster.metadata[0].namespace
  }
}

data "kubernetes_secret" "vcluster_kubeconfig" {
  metadata {
    name      = "vc-${var.vcluster_name}"
    namespace = var.vcluster_namespace
  }

  depends_on = [helm_release.vcluster]
}

output "kubeconfig" {
  value     = data.kubernetes_secret.vcluster_kubeconfig.data.config
  sensitive = true
}
