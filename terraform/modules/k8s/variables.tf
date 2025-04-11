# Variables for Kubernetes Module
variable "namespace" {
  description = "Kubernetes namespace for Agent Runtime"
  type        = string
  default     = "agent-runtime"
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "agent-runtime-cluster"
}

variable "kubernetes_version" {
  description = "Kubernetes version to use"
  type        = string
  default     = "1.28.0"
}

variable "cloud_provider" {
  description = "Cloud provider to use (aws, azure, ovh, fly)"
  type        = string
  default     = "aws"
  
  validation {
    condition     = contains(["aws", "azure", "ovh", "fly"], var.cloud_provider)
    error_message = "Valid values for cloud_provider are: aws, azure, ovh, fly."
  }
}

variable "region" {
  description = "Region to deploy the cluster in"
  type        = string
  default     = "us-east-1"
}

variable "node_count" {
  description = "Number of nodes in the cluster"
  type        = number
  default     = 3
}

variable "node_size" {
  description = "Size of the nodes (small, medium, large, xlarge)"
  type        = string
  default     = "medium"
  
  validation {
    condition     = contains(["small", "medium", "large", "xlarge"], var.node_size)
    error_message = "Valid values for node_size are: small, medium, large, xlarge."
  }
}

variable "vcluster_enabled" {
  description = "Enable vCluster deployment"
  type        = bool
  default     = true
}

variable "vcluster_version" {
  description = "vCluster version to deploy"
  type        = string
  default     = "0.15.0"
}

variable "vnode_enabled" {
  description = "Whether to enable vNode runtime"
  type        = bool
  default     = true
}

variable "jspolicy_enabled" {
  description = "Enable jsPolicy deployment"
  type        = bool
  default     = true
}

variable "monitoring_enabled" {
  description = "Whether to enable monitoring stack"
  type        = bool
  default     = true
}

variable "rocketmq_enabled" {
  description = "Whether to enable RocketMQ"
  type        = bool
  default     = true
}

variable "dragonfly_enabled" {
  description = "Whether to enable DragonflyDB"
  type        = bool
  default     = true
}

variable "supabase_enabled" {
  description = "Whether to enable Supabase"
  type        = bool
  default     = true
}

variable "kata_enabled" {
  description = "Whether to enable Kata Containers"
  type        = bool
  default     = true
}

variable "ragflow_enabled" {
  description = "Whether to enable RAGflow"
  type        = bool
  default     = true
}

variable "aws_vpc_cidr" {
  description = "CIDR block for AWS VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "admin_group_object_ids" {
  description = "Azure AD group object IDs for cluster admin access"
  type        = list(string)
  default     = []
}

variable "ovh_service_name" {
  description = "OVH Public Cloud service name"
  type        = string
  default     = ""
}

variable "monthly_billed" {
  description = "Whether to use monthly billing for OVH nodes"
  type        = bool
  default     = true
}

variable "k3s_token" {
  description = "Token for k3s cluster on Fly.io"
  type        = string
  default     = ""
  sensitive   = true
}
