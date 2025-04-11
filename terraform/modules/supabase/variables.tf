variable "namespace" {
  description = "Kubernetes namespace for Supabase deployment"
  type        = string
  default     = "supabase"
}

variable "create_namespace" {
  description = "Whether to create the Supabase namespace"
  type        = bool
  default     = true
}

variable "name" {
  description = "Name for Supabase resources"
  type        = string
  default     = "supabase"
}

variable "high_availability" {
  description = "Whether to enable high availability for Supabase"
  type        = bool
  default     = true
}

variable "instance_count" {
  description = "Number of Supabase instances to deploy"
  type        = number
  default     = 4
}

variable "postgres_version" {
  description = "PostgreSQL version for Supabase"
  type        = string
  default     = "15.2"
}

variable "postgres_replicas" {
  description = "Number of PostgreSQL read replicas"
  type        = number
  default     = 3
}

variable "postgres_storage_size" {
  description = "Storage size for PostgreSQL"
  type        = string
  default     = "50Gi"
}

variable "postgres_storage_class" {
  description = "Storage class for PostgreSQL"
  type        = string
  default     = null
}

variable "postgres_resources" {
  description = "Resource limits and requests for PostgreSQL"
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
      cpu    = "2000m"
      memory = "4Gi"
    }
    requests = {
      cpu    = "1000m"
      memory = "2Gi"
    }
  }
}

variable "postgres_backup_schedule" {
  description = "Schedule for PostgreSQL backups (cron format)"
  type        = string
  default     = "0 2 * * *"
}

variable "postgres_backup_retention" {
  description = "Number of PostgreSQL backups to retain"
  type        = number
  default     = 7
}

variable "postgres_rollback_enabled" {
  description = "Whether to enable rollback capability for PostgreSQL"
  type        = bool
  default     = true
}

variable "postgres_previous_state_retention" {
  description = "Number of previous states to retain for rollback"
  type        = number
  default     = 5
}

variable "postgres_max_connections" {
  description = "Maximum number of connections for PostgreSQL"
  type        = number
  default     = 200
}

variable "postgres_shared_buffers" {
  description = "Shared buffers for PostgreSQL"
  type        = string
  default     = "1GB"
}

variable "postgres_wal_level" {
  description = "WAL level for PostgreSQL"
  type        = string
  default     = "logical"
}

variable "postgres_max_wal_senders" {
  description = "Maximum number of WAL senders for PostgreSQL"
  type        = number
  default     = 10
}

variable "postgres_max_replication_slots" {
  description = "Maximum number of replication slots for PostgreSQL"
  type        = number
  default     = 10
}

variable "postgres_password" {
  description = "Password for PostgreSQL"
  type        = string
  sensitive   = true
  default     = ""
}

variable "postgres_admin_password" {
  description = "Admin password for PostgreSQL"
  type        = string
  sensitive   = true
  default     = ""
}

variable "postgres_replication_password" {
  description = "Replication password for PostgreSQL"
  type        = string
  sensitive   = true
  default     = ""
}

variable "postgres_databases" {
  description = "List of PostgreSQL databases to create"
  type = list(object({
    name  = string
    owner = string
  }))
  default = [
    {
      name  = "agent_state"
      owner = "agent"
    },
    {
      name  = "task_state"
      owner = "agent"
    },
    {
      name  = "tool_state"
      owner = "agent"
    },
    {
      name  = "mcp_state"
      owner = "agent"
    },
    {
      name  = "prompts_state"
      owner = "agent"
    },
    {
      name  = "modules_state"
      owner = "agent"
    }
  ]
}

variable "postgres_users" {
  description = "List of PostgreSQL users to create"
  type = list(object({
    name     = string
    password = string
    superuser = bool
    create_db = bool
    create_role = bool
  }))
  default = [
    {
      name     = "agent"
      password = ""
      superuser = false
      create_db = false
      create_role = false
    },
    {
      name     = "readonly"
      password = ""
      superuser = false
      create_db = false
      create_role = false
    }
  ]
}

variable "postgres_extensions" {
  description = "List of PostgreSQL extensions to enable"
  type        = list(string)
  default     = ["pg_stat_statements", "pgcrypto", "uuid-ossp", "pg_trgm"]
}

variable "realtime_enabled" {
  description = "Whether to enable Supabase Realtime"
  type        = bool
  default     = true
}

variable "realtime_replicas" {
  description = "Number of Realtime replicas"
  type        = number
  default     = 2
}

variable "realtime_resources" {
  description = "Resource limits and requests for Realtime"
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
      cpu    = "1000m"
      memory = "1Gi"
    }
    requests = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "auth_enabled" {
  description = "Whether to enable Supabase Auth"
  type        = bool
  default     = true
}

variable "auth_replicas" {
  description = "Number of Auth replicas"
  type        = number
  default     = 2
}

variable "auth_resources" {
  description = "Resource limits and requests for Auth"
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
      cpu    = "1000m"
      memory = "1Gi"
    }
    requests = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "storage_enabled" {
  description = "Whether to enable Supabase Storage"
  type        = bool
  default     = true
}

variable "storage_replicas" {
  description = "Number of Storage replicas"
  type        = number
  default     = 2
}

variable "storage_resources" {
  description = "Resource limits and requests for Storage"
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
      cpu    = "1000m"
      memory = "1Gi"
    }
    requests = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "rest_enabled" {
  description = "Whether to enable Supabase REST API"
  type        = bool
  default     = true
}

variable "rest_replicas" {
  description = "Number of REST API replicas"
  type        = number
  default     = 2
}

variable "rest_resources" {
  description = "Resource limits and requests for REST API"
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
      cpu    = "1000m"
      memory = "1Gi"
    }
    requests = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "kong_enabled" {
  description = "Whether to enable Kong API Gateway"
  type        = bool
  default     = true
}

variable "kong_replicas" {
  description = "Number of Kong replicas"
  type        = number
  default     = 2
}

variable "kong_resources" {
  description = "Resource limits and requests for Kong"
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
      cpu    = "1000m"
      memory = "1Gi"
    }
    requests = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "ingress_enabled" {
  description = "Whether to enable ingress for Supabase"
  type        = bool
  default     = false
}

variable "ingress_domain" {
  description = "Domain for Supabase ingress"
  type        = string
  default     = ""
}

variable "ingress_class" {
  description = "Ingress class for Supabase"
  type        = string
  default     = "nginx"
}

variable "ingress_tls_enabled" {
  description = "Whether to enable TLS for Supabase ingress"
  type        = bool
  default     = false
}

variable "ingress_tls_secret" {
  description = "TLS secret for Supabase ingress"
  type        = string
  default     = ""
}

variable "prometheus_integration" {
  description = "Whether to enable Prometheus integration"
  type        = bool
  default     = true
}

variable "kata_container_integration" {
  description = "Whether to enable Kata Containers integration"
  type        = bool
  default     = true
}

variable "labels" {
  description = "Additional labels for Supabase resources"
  type        = map(string)
  default     = {}
}

variable "annotations" {
  description = "Additional annotations for Supabase resources"
  type        = map(string)
  default     = {}
}

variable "jwt_secret" {
  description = "JWT secret for Supabase"
  type        = string
  sensitive   = true
  default     = ""
}

variable "anon_key" {
  description = "Anonymous key for Supabase"
  type        = string
  sensitive   = true
  default     = ""
}

variable "service_role_key" {
  description = "Service role key for Supabase"
  type        = string
  sensitive   = true
  default     = ""
}
