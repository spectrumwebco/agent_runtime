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

variable "mcp_host_url" {
  description = "URL for the MCP host"
  type        = string
  default     = "http://mcp-host:8080"
}

variable "mcp_server_url" {
  description = "URL for the MCP server"
  type        = string
  default     = "http://mcp-server:8080"
}

variable "librechat_code_api_key" {
  description = "API key for LibreChat Code Interpreter"
  type        = string
  sensitive   = true
}

variable "argocd_chart_version" {
  description = "Version of the ArgoCD Helm chart"
  type        = string
  default     = "5.51.4"
}

variable "argocd_values_yaml" {
  description = "Values YAML for the ArgoCD Helm chart"
  type        = string
  default     = ""
}

variable "flux_git_repository_url" {
  description = "URL of the Git repository for Flux"
  type        = string
}

variable "flux_git_branch" {
  description = "Branch of the Git repository for Flux"
  type        = string
  default     = "main"
}

variable "flux_sync_interval" {
  description = "Interval for Flux synchronization"
  type        = string
  default     = "1m0s"
}

variable "vnode_version" {
  description = "Version of vNode runtime to deploy"
  type        = string
  default     = "0.0.2"
}

variable "jspolicy_version" {
  description = "Version of jsPolicy to deploy"
  type        = string
  default     = "0.3.0-beta.5"
}

variable "vcluster_version" {
  description = "Version of vCluster to deploy"
  type        = string
  default     = "0.15.0"
}

variable "argocd_chart_version" {
  description = "Version of ArgoCD Helm chart to deploy"
  type        = string
  default     = "5.16.14"
}

variable "argocd_values_yaml" {
  description = "Values YAML for ArgoCD Helm chart"
  type        = string
  default     = ""
}

variable "flux_git_repository_url" {
  description = "Git repository URL for Flux"
  type        = string
  default     = "https://github.com/spectrumwebco/agent_runtime"
}

variable "flux_git_branch" {
  description = "Git branch for Flux"
  type        = string
  default     = "main"
}

variable "flux_sync_interval" {
  description = "Sync interval for Flux"
  type        = string
  default     = "1m"
}

variable "dragonfly_password" {
  description = "Password for DragonflyDB"
  type        = string
  sensitive   = true
  default     = ""
}

variable "mcp_host_url" {
  description = "URL for MCP host"
  type        = string
  default     = "http://mcp-host.mcp-system.svc.cluster.local:8080"
}

variable "mcp_server_url" {
  description = "URL for MCP server"
  type        = string
  default     = "http://mcp-server.mcp-system.svc.cluster.local:8080"
}

variable "librechat_code_api_key" {
  description = "API key for LibreChat Code Interpreter"
  type        = string
  sensitive   = true
  default     = ""
}
