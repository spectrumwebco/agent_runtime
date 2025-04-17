variable "module_name" {
  description = "Name of the Terraform module to generate"
  type        = string
}

variable "module_description" {
  description = "Description of the Terraform module"
  type        = string
}

variable "output_path" {
  description = "Path to output the generated module"
  type        = string
}
