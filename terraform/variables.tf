
variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "postgres_operator_namespace" {
  description = "Namespace for the PostgreSQL Operator"
  type        = string
  default     = "postgres-operator"
}

variable "postgres_version" {
  description = "PostgreSQL version to deploy"
  type        = number
  default     = 15
}

variable "postgres_replicas" {
  description = "Number of PostgreSQL replicas"
  type        = number
  default     = 3
}

variable "postgres_storage_size" {
  description = "Storage size for PostgreSQL data"
  type        = string
  default     = "10Gi"
}

variable "backup_storage_size" {
  description = "Storage size for PostgreSQL backups"
  type        = string
  default     = "20Gi"
}

variable "vault_integration_enabled" {
  description = "Enable Vault integration for PostgreSQL"
  type        = bool
  default     = true
}

variable "vault_address" {
  description = "Vault server address"
  type        = string
  default     = "http://vault.vault.svc.cluster.local:8200"
}

variable "vnode_runtime_enabled" {
  description = "Enable vNode runtime integration"
  type        = bool
  default     = true
}

variable "vnode_runtime_version" {
  description = "vNode runtime version"
  type        = string
  default     = "0.0.2"
}

variable "vnode_runtime_namespace" {
  description = "Namespace for vNode runtime deployment"
  type        = string
  default     = "vnode-runtime"
}

variable "vnode_runtime_replica_count" {
  description = "Number of vNode runtime replicas"
  type        = number
  default     = 1
}

variable "vnode_runtime_resources_limits_cpu" {
  description = "CPU limits for vNode runtime"
  type        = string
  default     = "500m"
}

variable "vnode_runtime_resources_limits_memory" {
  description = "Memory limits for vNode runtime"
  type        = string
  default     = "512Mi"
}

variable "vnode_runtime_resources_requests_cpu" {
  description = "CPU requests for vNode runtime"
  type        = string
  default     = "250m"
}

variable "vnode_runtime_resources_requests_memory" {
  description = "Memory requests for vNode runtime"
  type        = string
  default     = "256Mi"
}

variable "vault_namespace" {
  description = "Namespace for Vault deployment"
  type        = string
  default     = "vault"
}

variable "vault_version" {
  description = "Version of Vault to deploy"
  type        = string
  default     = "1.13.0"
}

variable "vault_k8s_version" {
  description = "Version of Vault K8s to deploy"
  type        = string
  default     = "1.1.0"
}

variable "vault_token" {
  description = "Root token for Vault"
  type        = string
  default     = "vault-token-secret-key"
  sensitive   = true
}

variable "vault_resources_limits_cpu" {
  description = "CPU limits for Vault"
  type        = string
  default     = "500m"
}

variable "vault_resources_limits_memory" {
  description = "Memory limits for Vault"
  type        = string
  default     = "512Mi"
}

variable "vault_resources_requests_cpu" {
  description = "CPU requests for Vault"
  type        = string
  default     = "250m"
}

variable "vault_resources_requests_memory" {
  description = "Memory requests for Vault"
  type        = string
  default     = "256Mi"
}

variable "kubeflow_version" {
  description = "Version of KubeFlow to deploy"
  type        = string
  default     = "1.7.0"
}

variable "training_operator_version" {
  description = "Version of KubeFlow Training Operator to deploy"
  type        = string
  default     = "v1.7.0"
}

variable "katib_version" {
  description = "Version of KubeFlow Katib to deploy"
  type        = string
  default     = "v0.15.0"
}

variable "kubeflow_data_storage_size" {
  description = "Size of persistent volume claim for KubeFlow data"
  type        = string
  default     = "100Gi"
}

variable "mlflow_version" {
  description = "Version of MLFlow to deploy"
  type        = string
  default     = "2.3.0"
}

variable "mlflow_namespace" {
  description = "Namespace for MLFlow deployment"
  type        = string
  default     = "mlflow"
}

variable "mlflow_storage_size" {
  description = "Size of persistent volume claim for MLFlow"
  type        = string
  default     = "50Gi"
}

variable "mlflow_tracking_uri" {
  description = "MLFlow tracking URI"
  type        = string
  default     = "http://mlflow-server.mlflow.svc.cluster.local:5000"
}

variable "kserve_version" {
  description = "Version of KServe to deploy"
  type        = string
  default     = "0.10.0"
}

variable "kserve_namespace" {
  description = "Namespace for KServe deployment"
  type        = string
  default     = "kserve"
}

variable "minio_version" {
  description = "Version of MinIO to deploy"
  type        = string
  default     = "RELEASE.2023-05-04T21-44-30Z"
}

variable "minio_namespace" {
  description = "Namespace for MinIO deployment"
  type        = string
  default     = "minio"
}

variable "minio_storage_size" {
  description = "Size of persistent volume claim for MinIO"
  type        = string
  default     = "100Gi"
}

variable "minio_access_key" {
  description = "MinIO access key"
  type        = string
  default     = "minioadmin"
  sensitive   = true
}

variable "minio_secret_key" {
  description = "MinIO secret key"
  type        = string
  default     = "minioadmin"
  sensitive   = true
}

variable "feast_version" {
  description = "Version of Feast to deploy"
  type        = string
  default     = "0.30.0"
}

variable "feast_namespace" {
  description = "Namespace for Feast deployment"
  type        = string
  default     = "feast"
}

variable "jupyterhub_version" {
  description = "Version of JupyterHub to deploy"
  type        = string
  default     = "2.0.0"
}

variable "jupyterhub_namespace" {
  description = "Namespace for JupyterHub deployment"
  type        = string
  default     = "jupyter"
}

variable "jupyterhub_storage_size" {
  description = "Size of persistent volume claim for JupyterHub"
  type        = string
  default     = "50Gi"
}

variable "seldon_version" {
  description = "Version of Seldon Core to deploy"
  type        = string
  default     = "1.15.0"
}

variable "seldon_namespace" {
  description = "Namespace for Seldon Core deployment"
  type        = string
  default     = "seldon"
}

variable "h2o_version" {
  description = "Version of h2o.ai to deploy"
  type        = string
  default     = "3.38.0.1"
}

variable "h2o_namespace" {
  description = "Namespace for h2o.ai deployment"
  type        = string
  default     = "h2o"
}

variable "h2o_storage_size" {
  description = "Size of persistent volume claim for h2o.ai"
  type        = string
  default     = "50Gi"
}

variable "storage_class_name" {
  description = "Storage class name for persistent volume claims"
  type        = string
  default     = "standard"
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
