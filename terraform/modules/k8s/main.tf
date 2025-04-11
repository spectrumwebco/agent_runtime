# Kubernetes Module for Agent Runtime
resource "kubernetes_namespace" "agent_runtime" {
  metadata {
    name = var.namespace
    
    labels = {
      "app.kubernetes.io/name" = "agent-runtime"
      "app.kubernetes.io/part-of" = "agent-runtime"
    }
  }
}

resource "helm_release" "vcluster" {
  count = var.vcluster_enabled ? 1 : 0
  
  name       = "vcluster"
  repository = "https://charts.loft.sh"
  chart      = "vcluster"
  version    = var.vcluster_version
  namespace  = kubernetes_namespace.agent_runtime.metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/vcluster/values.yaml")
  ]
}

resource "kubernetes_manifest" "jspolicy" {
  count = var.jspolicy_enabled ? 1 : 0
  
  manifest = yamldecode(file("${path.module}/../../k8s/jspolicy/policies.yaml"))
  
  depends_on = [
    helm_release.vcluster
  ]
}
