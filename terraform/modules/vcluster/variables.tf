variable "vcluster_name" {
  description = "Name of the vCluster deployment"
  type        = string
  default     = "agent-runtime-vcluster"
}

variable "vcluster_namespace" {
  description = "Namespace to deploy the vCluster"
  type        = string
  default     = "vcluster"
}

variable "kubernetes_version" {
  description = "Kubernetes version for the vCluster"
  type        = string
  default     = "1.28.0"
}

variable "distro" {
  description = "Kubernetes distribution to use (k3s, k0s, k8s)"
  type        = string
  default     = "k3s"
  
  validation {
    condition     = contains(["k3s", "k0s", "k8s"], var.distro)
    error_message = "Valid values for distro are: k3s, k0s, k8s."
  }
}

variable "persistent" {
  description = "Whether to use persistent storage for vCluster"
  type        = bool
  default     = true
}

variable "enable_vnode_integration" {
  description = "Whether to enable vNode runtime integration"
  type        = bool
  default     = true
}

variable "enable_high_availability" {
  description = "Whether to enable high availability for vCluster"
  type        = bool
  default     = true
}

variable "replicas" {
  description = "Number of vCluster replicas for high availability"
  type        = number
  default     = 3
}

variable "resource_limits" {
  description = "Resource limits for vCluster"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "2000m"
    memory = "4Gi"
  }
}

variable "resource_requests" {
  description = "Resource requests for vCluster"
  type = object({
    cpu    = string
    memory = string
  })
  default = {
    cpu    = "500m"
    memory = "1Gi"
  }
}

variable "storage_size" {
  description = "Size of persistent storage for vCluster"
  type        = string
  default     = "10Gi"
}

variable "sync_all_configmaps" {
  description = "Whether to sync all ConfigMaps from the host cluster"
  type        = bool
  default     = false
}

variable "sync_nodes" {
  description = "Whether to sync nodes from the host cluster"
  type        = bool
  default     = true
}

variable "sync_ingresses" {
  description = "Whether to sync ingresses from the host cluster"
  type        = bool
  default     = true
}

variable "isolation_enabled" {
  description = "Whether to enable network isolation for vCluster"
  type        = bool
  default     = true
}

variable "expose_via_service" {
  description = "Whether to expose vCluster via a service"
  type        = bool
  default     = true
}

variable "service_type" {
  description = "Type of service to expose vCluster (ClusterIP, NodePort, LoadBalancer)"
  type        = string
  default     = "ClusterIP"
  
  validation {
    condition     = contains(["ClusterIP", "NodePort", "LoadBalancer"], var.service_type)
    error_message = "Valid values for service_type are: ClusterIP, NodePort, LoadBalancer."
  }
}

variable "node_port" {
  description = "NodePort to use when service_type is NodePort"
  type        = number
  default     = 30080
}

variable "create_rbac" {
  description = "Whether to create RBAC resources for vCluster"
  type        = bool
  default     = true
}

variable "create_cluster_role" {
  description = "Whether to create a ClusterRole for vCluster"
  type        = bool
  default     = true
}

variable "create_service_account" {
  description = "Whether to create a ServiceAccount for vCluster"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Additional labels to add to vCluster resources"
  type        = map(string)
  default     = {}
}

variable "annotations" {
  description = "Additional annotations to add to vCluster resources"
  type        = map(string)
  default     = {}
}

variable "extra_args" {
  description = "Extra arguments to pass to vCluster"
  type        = list(string)
  default     = []
}
