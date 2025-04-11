# Kata Containers Module for Agent Runtime
resource "kubernetes_manifest" "kata_runtime_class" {
  manifest = yamldecode(file("${path.module}/../../k8s/kata-containers/config.yaml"))
}

resource "kubernetes_manifest" "kata_sandbox" {
  manifest = yamldecode(file("${path.module}/../../k8s/kata-containers/sandbox.yaml"))
  
  depends_on = [
    kubernetes_manifest.kata_runtime_class
  ]
}

resource "kubernetes_config_map" "kata_config" {
  metadata {
    name      = "kata-containers-config"
    namespace = var.namespace
  }
  
  data = {
    "configuration.toml" = file("${path.module}/../../configs/kata/configuration.toml")
  }
}

resource "kubernetes_secret" "kata_certs" {
  metadata {
    name      = "kata-containers-certs"
    namespace = var.namespace
  }
  
  data = {
    "ca.crt"     = var.kata_ca_cert
    "server.crt" = var.kata_server_cert
    "server.key" = var.kata_server_key
  }
  
  type = "Opaque"
}
