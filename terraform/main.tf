# Main Terraform Configuration for Agent Runtime
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 3.20.0"
    }
  }
  required_version = ">= 1.0.0"
  
  backend "s3" {
    bucket = "agent-runtime-terraform-state"
    key    = "agent-runtime/terraform.tfstate"
    region = "us-west-2"
    encrypt = true
  }
}

provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

provider "vault" {
  address = var.vault_addr
}

module "k8s" {
  source = "./modules/k8s"
  
  cluster_name = var.cluster_name
  namespace    = var.namespace
  
  vcluster_enabled = var.vcluster_enabled
  vcluster_version = var.vcluster_version
  
  jspolicy_enabled = var.jspolicy_enabled
}

module "vcluster" {
  source = "./modules/vcluster"
  
  count = var.vcluster_enabled ? 1 : 0
  
  name      = "${var.cluster_name}-vcluster"
  namespace = var.namespace
  version   = var.vcluster_version
  
  depends_on = [module.k8s]
}

module "jspolicy" {
  source = "./modules/jspolicy"
  
  count = var.jspolicy_enabled ? 1 : 0
  
  name      = "${var.cluster_name}-jspolicy"
  namespace = var.namespace
  
  depends_on = [module.k8s]
}

module "vnode" {
  source = "./modules/vnode"
  
  name      = "${var.cluster_name}-vnode"
  namespace = var.namespace
  vnode_image = "ghcr.io/loft-sh/vnode-runtime"
  vnode_version = "0.0.1-alpha.1"
  
  depends_on = [module.k8s, module.vcluster]
}

module "kata" {
  source = "./modules/kata"
  
  namespace = var.namespace
  node_selector = var.kata_node_selector
  
  depends_on = [module.k8s]
}

module "dragonfly" {
  source = "./modules/dragonfly"
  
  namespace = var.namespace
  replicas  = var.dragonfly_replicas
  dragonfly_password = var.dragonfly_password
  
  depends_on = [module.k8s]
}

module "rocketmq" {
  source = "./modules/rocketmq"
  
  namespace = var.namespace
  replicas  = var.rocketmq_replicas
  
  depends_on = [module.k8s]
}

module "monitoring" {
  source = "./modules/monitoring"
  
  namespace = var.namespace
  
  enable_prometheus = true
  enable_grafana = true
  enable_thanos = true
  enable_loki = true
  enable_jaeger = true
  enable_vector = true
  enable_opentelemetry = true
  enable_kube_state_metrics = true
  enable_cadvisor = true
  enable_kubernetes_dashboard = true
  
  depends_on = [module.k8s]
}

module "ragflow" {
  source = "./modules/ragflow"
  
  name = "${var.cluster_name}-ragflow"
  namespace = "ragflow"
  create_namespace = true
  vcluster_enabled = true
  
  enable_kata_container_integration = true
  
  depends_on = [module.k8s, module.vcluster, module.kata]
}

module "mcp" {
  source = "./modules/mcp"
  
  namespace = var.namespace
  librechat_code_api_key = var.librechat_code_api_key
  
  depends_on = [module.k8s, module.kata]
}
