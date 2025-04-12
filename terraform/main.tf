
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.20.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.9.0"
    }
  }
  required_version = ">= 1.0.0"
}

provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "ml_infrastructure" {
  metadata {
    name = "ml-infrastructure"
    labels = {
      "app.kubernetes.io/name"       = "ml-infrastructure"
      "app.kubernetes.io/instance"   = "ml-infrastructure-v1.0.0"
      "app.kubernetes.io/version"    = "v1.0.0"
      "app.kubernetes.io/component"  = "infrastructure"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

module "postgres_operator" {
  source = "./modules/postgres-operator"

  kubeconfig_path             = var.kubeconfig_path
  postgres_operator_namespace = var.postgres_operator_namespace
  postgres_version            = var.postgres_version
  postgres_replicas           = var.postgres_replicas
  postgres_storage_size       = var.postgres_storage_size
  backup_storage_size         = var.backup_storage_size
  storage_class_name          = var.storage_class_name
  vault_integration_enabled   = var.vault_integration_enabled
  vault_address               = var.vault_address
  vnode_runtime_enabled       = var.vnode_runtime_enabled
  vnode_runtime_version       = var.vnode_runtime_version
  depends_on                  = [kubernetes_namespace.ml_infrastructure]
}

module "kubeflow" {
  source = "./modules/kubeflow"

  kubeconfig_path             = var.kubeconfig_path
  kubeflow_version            = var.kubeflow_version
  training_operator_version   = var.training_operator_version
  katib_version               = var.katib_version
  kubeflow_data_storage_size  = var.kubeflow_data_storage_size
  storage_class_name          = var.storage_class_name
  depends_on                  = [module.postgres_operator]
}

module "mlflow" {
  source = "./modules/mlflow"

  kubeconfig_path          = var.kubeconfig_path
  mlflow_version           = var.mlflow_version
  mlflow_namespace         = var.mlflow_namespace
  mlflow_storage_size      = var.mlflow_storage_size
  storage_class_name       = var.storage_class_name
  minio_access_key         = var.minio_access_key
  minio_secret_key         = var.minio_secret_key
  mlflow_tracking_uri      = var.mlflow_tracking_uri
  depends_on               = [module.postgres_operator]
}

module "kserve" {
  source = "./modules/kserve"

  kubeconfig_path          = var.kubeconfig_path
  kserve_version           = var.kserve_version
  kserve_namespace         = var.kserve_namespace
  storage_class_name       = var.storage_class_name
  minio_access_key         = var.minio_access_key
  minio_secret_key         = var.minio_secret_key
  depends_on               = [kubernetes_namespace.ml_infrastructure]
}

module "minio" {
  source = "./modules/minio"

  kubeconfig_path          = var.kubeconfig_path
  minio_version            = var.minio_version
  minio_namespace          = var.minio_namespace
  minio_storage_size       = var.minio_storage_size
  storage_class_name       = var.storage_class_name
  minio_access_key         = var.minio_access_key
  minio_secret_key         = var.minio_secret_key
  depends_on               = [kubernetes_namespace.ml_infrastructure]
}

module "feast" {
  source = "./modules/feast"

  kubeconfig_path          = var.kubeconfig_path
  feast_version            = var.feast_version
  feast_namespace          = var.feast_namespace
  storage_class_name       = var.storage_class_name
  minio_access_key         = var.minio_access_key
  minio_secret_key         = var.minio_secret_key
  depends_on               = [module.minio]
}

module "jupyterhub" {
  source = "./modules/jupyterhub"

  kubeconfig_path          = var.kubeconfig_path
  jupyterhub_version       = var.jupyterhub_version
  jupyterhub_namespace     = var.jupyterhub_namespace
  jupyterhub_storage_size  = var.jupyterhub_storage_size
  storage_class_name       = var.storage_class_name
  depends_on               = [kubernetes_namespace.ml_infrastructure]
}

module "seldon" {
  source = "./modules/seldon"

  kubeconfig_path          = var.kubeconfig_path
  seldon_version           = var.seldon_version
  seldon_namespace         = var.seldon_namespace
  storage_class_name       = var.storage_class_name
  depends_on               = [kubernetes_namespace.ml_infrastructure]
}

module "h2o" {
  source = "./modules/h2o"

  kubeconfig_path          = var.kubeconfig_path
  h2o_version              = var.h2o_version
  h2o_namespace            = var.h2o_namespace
  h2o_storage_size         = var.h2o_storage_size
  storage_class_name       = var.storage_class_name
  depends_on               = [kubernetes_namespace.ml_infrastructure]
}

module "vault" {
  source = "./modules/vault"

  kubeconfig_path                 = var.kubeconfig_path
  vault_namespace                 = var.vault_namespace
  vault_version                   = var.vault_version
  vault_k8s_version               = var.vault_k8s_version
  vault_token                     = var.vault_token
  vault_resources_limits_cpu      = var.vault_resources_limits_cpu
  vault_resources_limits_memory   = var.vault_resources_limits_memory
  vault_resources_requests_cpu    = var.vault_resources_requests_cpu
  vault_resources_requests_memory = var.vault_resources_requests_memory
}
