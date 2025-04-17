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

resource "null_resource" "generate_module" {
  provisioner "local-exec" {
    command = "npx -y generator-tf-module --name=${var.module_name} --description='${var.module_description}' --path=${var.output_path}"
  }
}
