variable "name" {
  description = "Name of the jsPolicy deployment"
  type        = string
  default     = "jspolicy"
}

variable "namespace" {
  description = "Namespace to deploy jsPolicy"
  type        = string
  default     = "jspolicy-system"
}

variable "chart_version" {
  description = "Version of the jsPolicy Helm chart"
  type        = string
  default     = "0.5.0"
}

variable "image_tag" {
  description = "Tag of the jsPolicy image"
  type        = string
  default     = "latest"
}

variable "high_availability" {
  description = "Whether to enable high availability for jsPolicy"
  type        = bool
  default     = true
}

variable "resources" {
  description = "Resource limits and requests for jsPolicy"
  type = object({
    limits = object({
      cpu    = string
      memory = string
    })
    requests = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    limits = {
      cpu    = "500m"
      memory = "512Mi"
    }
    requests = {
      cpu    = "250m"
      memory = "256Mi"
    }
  }
}

variable "strict_mode" {
  description = "Whether to enforce policies strictly (fail on violation) or in advisory mode"
  type        = bool
  default     = false
}

variable "prometheus_enabled" {
  description = "Whether to enable Prometheus integration"
  type        = bool
  default     = true
}

variable "enable_reporting" {
  description = "Whether to enable policy reporting"
  type        = bool
  default     = true
}

variable "slack_webhook" {
  description = "Slack webhook URL for policy violation notifications"
  type        = string
  default     = ""
  sensitive   = true
}

variable "slack_channel" {
  description = "Slack channel for policy violation notifications"
  type        = string
  default     = "#policy-violations"
}

variable "enable_recovery_controller" {
  description = "Whether to enable the recovery controller"
  type        = bool
  default     = true
}

variable "recovery_controller_version" {
  description = "Version of the recovery controller image"
  type        = string
  default     = "latest"
}
