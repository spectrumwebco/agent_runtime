variable "namespace" {
  description = "Kubernetes namespace for PipeCD"
  type        = string
  default     = "pipecd"
}

variable "pipecd_version" {
  description = "Version of PipeCD to deploy"
  type        = string
  default     = "v0.51.1"
}

variable "control_plane_replicas" {
  description = "Number of replicas for the PipeCD control plane"
  type        = number
  default     = 1
}

variable "piped_replicas" {
  description = "Number of replicas for the PipeCD piped"
  type        = number
  default     = 1
}

variable "state_key" {
  description = "State key for PipeCD"
  type        = string
  default     = "pipecd-state-key"
}

variable "minio_endpoint" {
  description = "MinIO endpoint for PipeCD filestore"
  type        = string
  default     = "minio.minio.svc.cluster.local:9000"
}

variable "minio_bucket" {
  description = "MinIO bucket for PipeCD filestore"
  type        = string
  default     = "pipecd"
}

variable "minio_access_key" {
  description = "MinIO access key for PipeCD filestore"
  type        = string
  sensitive   = true
}

variable "minio_secret_key" {
  description = "MinIO secret key for PipeCD filestore"
  type        = string
  sensitive   = true
}

variable "ssh_key" {
  description = "SSH key for Git operations"
  type        = string
  sensitive   = true
  default     = ""
}

variable "project_id" {
  description = "Project ID for PipeCD"
  type        = string
  default     = "agent-runtime-project"
}

variable "piped_id" {
  description = "Piped ID for PipeCD"
  type        = string
  default     = "agent-runtime-piped"
}

variable "piped_key" {
  description = "Piped key for PipeCD"
  type        = string
  sensitive   = true
}

variable "kubeconfig" {
  description = "Kubeconfig for Kubernetes operations"
  type        = string
  sensitive   = true
  default     = ""
}

variable "repositories" {
  description = "Git repositories to be managed by PipeCD"
  type = list(object({
    repo_id = string
    remote  = string
    branch  = string
  }))
  default = [
    {
      repo_id = "agent-runtime"
      remote  = "https://github.com/spectrumwebco/agent_runtime.git"
      branch  = "main"
    }
  ]
}

variable "kubernetes_config" {
  description = "Kubernetes configuration for PipeCD"
  type = object({
    kubeconfig_path = string
  })
  default = {
    kubeconfig_path = "/etc/pipecd-secret/kubeconfig"
  }
}

variable "terraform_config" {
  description = "Terraform configuration for PipeCD"
  type = object({
    vars = map(string)
  })
  default = {
    vars = {
      project = "agent-runtime"
      region  = "us-central1"
    }
  }
}
