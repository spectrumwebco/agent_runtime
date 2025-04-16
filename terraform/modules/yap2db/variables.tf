variable "namespace" {
  description = "Kubernetes namespace to deploy Yap2DB"
  type        = string
  default     = "default"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = false
}

variable "replicas" {
  description = "Number of Yap2DB replicas"
  type        = number
  default     = 2
}

variable "container_registry" {
  description = "Container registry for Yap2DB images"
  type        = string
  default     = "spectrumwebco"
}

variable "yap2db_version" {
  description = "Yap2DB version"
  type        = string
  default     = "latest"
}

variable "supabase_host" {
  description = "Supabase host"
  type        = string
  default     = "supabase-db.default.svc.cluster.local"
}

variable "supabase_port" {
  description = "Supabase port"
  type        = string
  default     = "5432"
}

variable "dragonfly_host" {
  description = "DragonflyDB host"
  type        = string
  default     = "dragonfly-db.default.svc.cluster.local"
}

variable "dragonfly_port" {
  description = "DragonflyDB port"
  type        = string
  default     = "6379"
}

variable "rocketmq_host" {
  description = "RocketMQ host"
  type        = string
  default     = "rocketmq.default.svc.cluster.local"
}

variable "rocketmq_port" {
  description = "RocketMQ port"
  type        = string
  default     = "9876"
}

variable "ragflow_host" {
  description = "RAGflow host"
  type        = string
  default     = "ragflow.default.svc.cluster.local"
}

variable "ragflow_port" {
  description = "RAGflow port"
  type        = string
  default     = "8000"
}

variable "vcluster_enabled" {
  description = "Enable vCluster integration"
  type        = bool
  default     = true
}

variable "vnode_enabled" {
  description = "Enable vNode integration"
  type        = bool
  default     = true
}
