variable "namespace" {
  description = "Namespace for vNode resources"
  type        = string
  default     = "vnode-system"
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "vnode_version" {
  description = "Version of vNode runtime to deploy"
  type        = string
  default     = "0.0.2"
}
