variable "namespace" {
  description = "Kubernetes namespace for Kubefirst"
  type        = string
  default     = "kubefirst"
}

variable "kubefirst_version" {
  description = "Version of Kubefirst to deploy"
  type        = string
  default     = "v2.7.0"
}

variable "gitea_version" {
  description = "Version of Gitea to deploy"
  type        = string
  default     = "1.20"
}

variable "vault_version" {
  description = "Version of Vault to deploy"
  type        = string
  default     = "1.14"
}

variable "replicas" {
  description = "Number of replicas for the Kubefirst deployment"
  type        = number
  default     = 1
}

variable "git_provider" {
  description = "Git provider for Kubefirst"
  type        = string
  default     = "gitea"
}

variable "git_username" {
  description = "Git username for Kubefirst"
  type        = string
  default     = "ove"
}

variable "git_password" {
  description = "Git password for Kubefirst"
  type        = string
  sensitive   = true
}

variable "gitea_postgres_password" {
  description = "Password for Gitea Postgres database"
  type        = string
  sensitive   = true
  default     = "Password123"
}

variable "cloud_provider" {
  description = "Cloud provider for Kubefirst"
  type        = string
  default     = "k3d"
}

variable "cluster_name" {
  description = "Kubernetes cluster name for Kubefirst"
  type        = string
  default     = "agent-runtime-cluster"
}

variable "gitops_template_url" {
  description = "GitOps template URL for Kubefirst"
  type        = string
  default     = "https://github.com/kubefirst/gitops-template.git"
}

variable "gitops_template_branch" {
  description = "GitOps template branch for Kubefirst"
  type        = string
  default     = "main"
}

variable "vault_token" {
  description = "Vault root token for Kubefirst"
  type        = string
  sensitive   = true
  default     = "kubefirst-root-token"
}
