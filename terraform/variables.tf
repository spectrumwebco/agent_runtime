# Variables for Agent Runtime Terraform Configuration
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "vault_addr" {
  description = "Vault server address"
  type        = string
  default     = "https://vault.example.com:8200"
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "agent-runtime"
}

variable "namespace" {
  description = "Kubernetes namespace for Agent Runtime"
  type        = string
  default     = "agent-runtime-system"
}

variable "vcluster_enabled" {
  description = "Enable vCluster deployment"
  type        = bool
  default     = true
}

variable "vcluster_version" {
  description = "vCluster version to deploy"
  type        = string
  default     = "1.27"
}

variable "jspolicy_enabled" {
  description = "Enable jsPolicy deployment"
  type        = bool
  default     = true
}

variable "kata_node_selector" {
  description = "Node selector for Kata Containers"
  type        = map(string)
  default     = {
    "kata-containers" = "true"
  }
}

variable "dragonfly_replicas" {
  description = "Number of DragonflyDB replicas"
  type        = number
  default     = 3
}

variable "rocketmq_replicas" {
  description = "Number of RocketMQ replicas"
  type        = number
  default     = 3
}

variable "dragonfly_password" {
  description = "Password for DragonflyDB authentication"
  type        = string
  sensitive   = true
  default     = "changeme" 
}

variable "librechat_code_api_key" {
  description = "LibreChat Code Interpreter API key for MCP integration"
  type        = string
  sensitive   = true
}
