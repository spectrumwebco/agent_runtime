output "module_name" {
  description = "The name of the generated Terraform module"
  value       = var.module_name
}

output "module_path" {
  description = "The path where the Terraform module was generated"
  value       = var.output_path
}
