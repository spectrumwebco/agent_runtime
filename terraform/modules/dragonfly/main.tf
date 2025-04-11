# DragonflyDB Module for Agent Runtime
resource "helm_release" "dragonfly" {
  name       = "dragonfly"
  repository = "https://dragonflydb.github.io/dragonfly-operator" # Assuming official chart exists, adjust if needed
  chart      = "dragonfly"
  namespace  = var.namespace
  version    = "0.1.0" # Specify appropriate chart version

  set {
    name  = "replicas"
    value = var.replicas
  }

  set {
    name  = "persistence.enabled"
    value = "true"
  }

  set {
    name  = "persistence.size"
    value = "10Gi"
  }
  
  # Add other necessary DragonflyDB configurations here
  # e.g., resource limits, affinity rules, etc.
}

resource "kubernetes_secret" "dragonfly_auth" {
  metadata {
    name      = "dragonfly-auth"
    namespace = var.namespace
  }

  data = {
    # Store password securely, potentially fetched from Vault
    "requirepass" = var.dragonfly_password
  }
  type = "Opaque"

  depends_on = [helm_release.dragonfly]
}
