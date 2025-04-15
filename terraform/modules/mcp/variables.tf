variable "container_registry" {
  description = "Container registry for agent runtime images"
  type        = string
  default     = "ghcr.io/spectrumwebco"
}

variable "librechat_code_api_key" {
  description = "API key for LibreChat Code Interpreter"
  type        = string
  sensitive   = true
}

variable "replicas" {
  description = "Number of replicas for MCP client"
  type        = number
  default     = 1
}
