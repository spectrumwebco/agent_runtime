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

module "dragonfly" {
  source = "./modules/dragonfly"
  
  namespace = var.namespace
  replicas  = var.dragonfly_replicas
  dragonfly_password = var.dragonfly_password # Pass password variable
  
  depends_on = [module.k8s]
}

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
  
  depends_on = [module.k8s]
}

module "rocketmq" {
  source = "./modules/rocketmq"
  
  namespace = var.namespace
  replicas  = var.rocketmq_replicas
  
  depends_on = [module.k8s]
}
