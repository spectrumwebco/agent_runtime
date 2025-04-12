
resource "kubernetes_namespace" "argocd" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = merge({
      "app.kubernetes.io/name"       = "argocd"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
    annotations = var.annotations
  }
}

locals {
  namespace = var.create_namespace ? kubernetes_namespace.argocd[0].metadata[0].name : var.namespace
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = var.chart_version
  namespace  = local.namespace

  values = [
    var.values_yaml
  ]

  set {
    name  = "server.service.type"
    value = var.server_service_type
  }

  set {
    name  = "server.extraArgs"
    value = "{--insecure}"
  }

  set {
    name  = "controller.resources.limits.cpu"
    value = var.controller_cpu_limit
  }

  set {
    name  = "controller.resources.limits.memory"
    value = var.controller_memory_limit
  }

  set {
    name  = "controller.resources.requests.cpu"
    value = var.controller_cpu_request
  }

  set {
    name  = "controller.resources.requests.memory"
    value = var.controller_memory_request
  }

  set {
    name  = "server.resources.limits.cpu"
    value = var.server_cpu_limit
  }

  set {
    name  = "server.resources.limits.memory"
    value = var.server_memory_limit
  }

  set {
    name  = "server.resources.requests.cpu"
    value = var.server_cpu_request
  }

  set {
    name  = "server.resources.requests.memory"
    value = var.server_memory_request
  }

  set {
    name  = "repoServer.resources.limits.cpu"
    value = var.repo_server_cpu_limit
  }

  set {
    name  = "repoServer.resources.limits.memory"
    value = var.repo_server_memory_limit
  }

  set {
    name  = "repoServer.resources.requests.cpu"
    value = var.repo_server_cpu_request
  }

  set {
    name  = "repoServer.resources.requests.memory"
    value = var.repo_server_memory_request
  }
}

resource "kubernetes_manifest" "argocd_application" {
  count = length(var.applications)

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = var.applications[count.index].name
      namespace = local.namespace
    }
    spec = {
      project = var.applications[count.index].project
      source = {
        repoURL        = var.applications[count.index].repo_url
        targetRevision = var.applications[count.index].target_revision
        path           = var.applications[count.index].path
      }
      destination = {
        server    = var.applications[count.index].destination_server
        namespace = var.applications[count.index].destination_namespace
      }
      syncPolicy = {
        automated = {
          prune     = var.applications[count.index].prune
          selfHeal  = var.applications[count.index].self_heal
          allowEmpty = var.applications[count.index].allow_empty
        }
        syncOptions = var.applications[count.index].sync_options
      }
    }
  }

  depends_on = [
    helm_release.argocd
  ]
}
