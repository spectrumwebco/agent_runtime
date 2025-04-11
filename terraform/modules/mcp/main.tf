terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

resource "kubernetes_manifest" "mcp_client" {
  manifest = yamldecode(file("${path.module}/../../k8s/mcp/client.yaml"))
}

resource "kubernetes_manifest" "mcp_host" {
  manifest = yamldecode(file("${path.module}/../../k8s/mcp/host.yaml"))
}

resource "kubernetes_secret" "agent_runtime_secrets" {
  metadata {
    name      = "agent-runtime-secrets"
    namespace = var.namespace
  }

  data = {
    librechat_code_api_key = var.librechat_code_api_key
  }

  type = "Opaque"
}
