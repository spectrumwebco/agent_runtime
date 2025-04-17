
module "otf" {
  source = "./otf"
  
  namespace = var.otf_namespace
  api_key   = var.otf_api_key
}

module "kubestack" {
  source = "./kubestack"
  
  namespace = var.kubestack_namespace
}

module "lynx" {
  source = "./lynx"
  
  namespace = var.lynx_namespace
}

module "terraform_operator" {
  source = "./terraform-operator"
  
  namespace = var.terraform_operator_namespace
}

module "generator_tf_module" {
  source = "./generator-tf-module"
  count  = var.enable_generator_tf_module ? 1 : 0
  
  module_name        = var.generator_module_name
  module_description = var.generator_module_description
  output_path        = var.generator_output_path
}

module "pipecd" {
  source = "./modules/pipecd"

  namespace             = var.pipecd_namespace
  pipecd_version        = var.pipecd_version
  control_plane_replicas = var.pipecd_control_plane_replicas
  piped_replicas        = var.pipecd_piped_replicas
  state_key             = var.pipecd_state_key
  minio_endpoint        = var.minio_endpoint
  minio_bucket          = var.pipecd_minio_bucket
  minio_access_key      = var.minio_access_key
  minio_secret_key      = var.minio_secret_key
  ssh_key               = var.ssh_key
  project_id            = var.pipecd_project_id
  piped_id              = var.pipecd_piped_id
  piped_key             = var.pipecd_piped_key
  kubeconfig            = var.kubeconfig
  repositories          = var.pipecd_repositories
  kubernetes_config     = var.pipecd_kubernetes_config
  terraform_config      = var.pipecd_terraform_config
}

module "kubefirst" {
  source = "./modules/kubefirst"

  namespace               = var.kubefirst_namespace
  kubefirst_version       = var.kubefirst_version
  gitea_version           = var.gitea_version
  vault_version           = var.vault_version
  replicas                = var.kubefirst_replicas
  git_provider            = var.kubefirst_git_provider
  git_username            = var.kubefirst_git_username
  git_password            = var.kubefirst_git_password
  gitea_postgres_password = var.kubefirst_gitea_postgres_password
  cloud_provider          = var.kubefirst_cloud_provider
  cluster_name            = var.kubefirst_cluster_name
  gitops_template_url     = var.kubefirst_gitops_template_url
  gitops_template_branch  = var.kubefirst_gitops_template_branch
  vault_token             = var.kubefirst_vault_token
}

module "neovim" {
  source = "./modules/neovim"

  namespace         = var.neovim_namespace
  replicas          = var.neovim_replicas
  image             = var.neovim_image
  supabase_url      = var.supabase_url
  supabase_key      = var.supabase_key
  storage_size      = var.neovim_storage_size
  enable_kata       = var.neovim_enable_kata
  kata_runtime_class = var.kata_runtime_class
  resource_limits   = var.neovim_resource_limits
  resource_requests = var.neovim_resource_requests
}
