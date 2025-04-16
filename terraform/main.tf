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
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.14.0"
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

provider "kubectl" {
  config_path = var.kubeconfig_path
  alias       = "gavinbunney"
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

module "doris" {
  source = "./modules/doris"
  
  namespace = var.namespace
  doris_version = var.doris_version
  fe_replicas = var.doris_fe_replicas
  be_replicas = var.doris_be_replicas
  admin_password = var.doris_admin_password
  
  depends_on = [module.k8s]
}

module "kafka" {
  source = "./modules/kafka"
  
  namespace = var.namespace
  kafka_version = var.kafka_version
  zookeeper_version = var.zookeeper_version
  kafka_replicas = var.kafka_replicas
  
  depends_on = [module.k8s]
}

module "postgres_operator" {
  source = "./modules/postgres"
  
  namespace = var.namespace
  cluster_name = var.postgres_cluster_name
  replicas = var.postgres_replicas
  
  depends_on = [module.k8s]
}

module "kafka_k8s_monitor" {
  source = "./modules/kafka_k8s_monitor"
  
  namespace = var.namespace
  kafka_replicas = var.kafka_replicas
  monitor_namespace = var.monitor_namespace
  poll_interval = var.poll_interval
  resources_to_monitor = var.resources_to_monitor
  
  depends_on = [module.k8s, module.kafka]
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
  
  namespace = "${var.namespace}-mcp"
  create_namespace = true
  
  mcp_host_url = var.mcp_host_url
  mcp_server_url = var.mcp_server_url
  librechat_code_api_key = var.librechat_code_api_key
  
  kata_container_integration = true
  
  depends_on = [module.k8s, module.kata]
}

module "argocd" {
  source = "./modules/argocd"
  
  namespace = "argocd"
  create_namespace = true
  
  chart_version = var.argocd_chart_version
  values_yaml = var.argocd_values_yaml
  
  depends_on = [module.k8s]
}

module "flux_system" {
  source = "./modules/flux-system"
  
  namespace = "flux-system"
  create_namespace = true
  
  git_repository_url = var.flux_git_repository_url
  git_branch = var.flux_git_branch
  sync_interval = var.flux_sync_interval
  
  depends_on = [module.k8s]
}

module "vnode" {
  source = "./modules/vnode"
  
  namespace = "vnode-system"
  create_namespace = true
  
  vnode_version = var.vnode_version
  
  depends_on = [module.k8s]
}

module "jspolicy" {
  source = "./modules/jspolicy"
  
  namespace = "jspolicy-system"
  create_namespace = true
  
  jspolicy_version = var.jspolicy_version
  
  depends_on = [module.k8s]
}

module "vcluster" {
  source = "./modules/vcluster"
  
  namespace = "vcluster-system"
  create_namespace = true
  
  vcluster_version = var.vcluster_version
  
  depends_on = [module.k8s]
}

module "supabase" {
  source = "./modules/supabase"
  
  namespace = "supabase-system"
  create_namespace = true
  
  depends_on = [module.k8s]
}

module "dragonfly" {
  source = "./modules/dragonfly"
  
  namespace = "dragonfly-system"
  create_namespace = true
  
  password = var.dragonfly_password
  
  depends_on = [module.k8s]
}

module "k8s_base" {
  source = "./modules/k8s-base"
  
  namespace = "agent-runtime-system"
  create_namespace = true
}
