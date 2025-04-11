variable "namespace" {
  description = "Kubernetes namespace for MCP components"
  type        = string
  default     = "agent-runtime-system"
}

variable "librechat_code_api_key" {
  description = "LibreChat Code Interpreter API key"
  type        = string
  sensitive   = true
}
