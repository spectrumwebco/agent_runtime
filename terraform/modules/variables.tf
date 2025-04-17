
variable "otf_namespace" {
  description = "Kubernetes namespace for OTF deployment"
  type        = string
  default     = "otf"
}

variable "otf_api_key" {
  description = "API key for OTF"
  type        = string
  sensitive   = true
  default     = ""
}

variable "kubestack_namespace" {
  description = "Kubernetes namespace for Kubestack deployment"
  type        = string
  default     = "kubestack"
}

variable "lynx_namespace" {
  description = "Kubernetes namespace for Lynx deployment"
  type        = string
  default     = "lynx"
}

variable "terraform_operator_namespace" {
  description = "Kubernetes namespace for Terraform Operator deployment"
  type        = string
  default     = "terraform-operator"
}


variable "enable_generator_tf_module" {
  description = "Whether to enable the Terraform module generator"
  type        = bool
  default     = false
}

variable "generator_module_name" {
  description = "Name of the Terraform module to generate"
  type        = string
  default     = "example"
}

variable "generator_module_description" {
  description = "Description of the Terraform module"
  type        = string
  default     = "Example Terraform module"
}

variable "generator_output_path" {
  description = "Path to output the generated module"
  type        = string
  default     = "/tmp/example-module"
}
