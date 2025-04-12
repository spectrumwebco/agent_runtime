
resource "kubernetes_namespace" "flux_system" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = merge({
      "app.kubernetes.io/name"       = "flux-system"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
    annotations = var.annotations
  }
}

locals {
  namespace = var.create_namespace ? kubernetes_namespace.flux_system[0].metadata[0].name : var.namespace
}

resource "kubernetes_secret" "flux_system" {
  metadata {
    name      = "flux-system"
    namespace = local.namespace
    labels = merge({
      "app.kubernetes.io/name"       = "flux-system"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
  }

  data = {
    "identity"       = var.identity
    "identity.pub"   = var.identity_pub
    "known_hosts"    = var.known_hosts
  }

  type = "Opaque"
}

resource "kubernetes_manifest" "gotk_source_git" {
  manifest = {
    apiVersion = "source.toolkit.fluxcd.io/v1"
    kind       = "GitRepository"
    metadata = {
      name      = var.git_repository_name
      namespace = local.namespace
    }
    spec = {
      interval = var.sync_interval
      url      = var.git_repository_url
      ref = {
        branch = var.git_branch
      }
      secretRef = {
        name = kubernetes_secret.flux_system.metadata[0].name
      }
    }
  }
}

resource "kubernetes_manifest" "gotk_kustomization" {
  manifest = {
    apiVersion = "kustomize.toolkit.fluxcd.io/v1"
    kind       = "Kustomization"
    metadata = {
      name      = var.kustomization_name
      namespace = local.namespace
    }
    spec = {
      interval = var.sync_interval
      path     = var.kustomization_path
      prune    = var.prune
      sourceRef = {
        kind = "GitRepository"
        name = var.git_repository_name
      }
    }
  }

  depends_on = [
    kubernetes_manifest.gotk_source_git
  ]
}
