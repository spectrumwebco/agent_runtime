
resource "kubernetes_namespace" "supabase" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    labels = merge({
      "app.kubernetes.io/name"       = "supabase"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
    annotations = var.annotations
  }
}

locals {
  namespace = var.create_namespace ? kubernetes_namespace.supabase[0].metadata[0].name : var.namespace
  
  common_labels = merge({
    "app.kubernetes.io/name"       = "supabase"
    "app.kubernetes.io/part-of"    = "agent-runtime"
    "app.kubernetes.io/managed-by" = "terraform"
  }, var.labels)
  
  kata_annotations = var.kata_container_integration ? {
    "io.kubernetes.cri.untrusted-workload" = "true"
    "io.kubernetes.cri.runtimeclass"       = "kata"
  } : {}
  
  postgres_password = var.postgres_password != "" ? var.postgres_password : random_password.postgres_password[0].result
  postgres_admin_password = var.postgres_admin_password != "" ? var.postgres_admin_password : random_password.postgres_admin_password[0].result
  postgres_replication_password = var.postgres_replication_password != "" ? var.postgres_replication_password : random_password.postgres_replication_password[0].result
  jwt_secret = var.jwt_secret != "" ? var.jwt_secret : random_password.jwt_secret[0].result
  anon_key = var.anon_key != "" ? var.anon_key : random_password.anon_key[0].result
  service_role_key = var.service_role_key != "" ? var.service_role_key : random_password.service_role_key[0].result
  
  instance_count = var.high_availability ? var.instance_count : 1
  postgres_replicas = var.high_availability ? var.postgres_replicas : 0
}

resource "random_password" "postgres_password" {
  count   = var.postgres_password == "" ? 1 : 0
  length  = 32
  special = false
}

resource "random_password" "postgres_admin_password" {
  count   = var.postgres_admin_password == "" ? 1 : 0
  length  = 32
  special = false
}

resource "random_password" "postgres_replication_password" {
  count   = var.postgres_replication_password == "" ? 1 : 0
  length  = 32
  special = false
}

resource "random_password" "jwt_secret" {
  count   = var.jwt_secret == "" ? 1 : 0
  length  = 64
  special = false
}

resource "random_password" "anon_key" {
  count   = var.anon_key == "" ? 1 : 0
  length  = 64
  special = false
}

resource "random_password" "service_role_key" {
  count   = var.service_role_key == "" ? 1 : 0
  length  = 64
  special = false
}

resource "kubernetes_secret" "postgres_credentials" {
  metadata {
    name      = "${var.name}-postgres-credentials"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "POSTGRES_PASSWORD"         = local.postgres_password
    "POSTGRES_ADMIN_PASSWORD"   = local.postgres_admin_password
    "POSTGRES_REPLICATION_PASSWORD" = local.postgres_replication_password
    "JWT_SECRET"                = local.jwt_secret
    "ANON_KEY"                  = local.anon_key
    "SERVICE_ROLE_KEY"          = local.service_role_key
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "postgres_config" {
  metadata {
    name      = "${var.name}-postgres-config"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "postgresql.conf" = <<-EOT
      listen_addresses = '*'
      port = 5432
      
      shared_buffers = ${var.postgres_shared_buffers}
      work_mem = '64MB'
      
      wal_level = ${var.postgres_wal_level}
      max_wal_senders = ${var.postgres_max_wal_senders}
      max_replication_slots = ${var.postgres_max_replication_slots}
      
      max_connections = ${var.postgres_max_connections}
      
      hot_standby = on
      hot_standby_feedback = on
    EOT
    
    "pg_hba.conf" = <<-EOT
      
      local   all             all                                     trust
      
      host    all             all             127.0.0.1/32            md5
      
      host    all             all             ::1/128                 md5
      
      host    replication     replicator      all                     md5
      
      host    all             all             all                     md5
    EOT
    
    "setup-primary.sh" = <<-EOT
      set -e
      
      psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
        CREATE USER replicator WITH REPLICATION PASSWORD '${local.postgres_replication_password}';
      EOSQL
      
      psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
        CREATE USER agent WITH PASSWORD 'agent_password';
        CREATE USER readonly WITH PASSWORD 'readonly_password';
      EOSQL
      
      for db in agent_state task_state tool_state mcp_state prompts_state modules_state; do
        psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
          CREATE DATABASE $db OWNER agent;
        EOSQL
        
        psql -v ON_ERROR_STOP=1 --username postgres -d $db <<-EOSQL
          CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
          CREATE EXTENSION IF NOT EXISTS "pgcrypto";
          CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        EOSQL
      done
      
      psql -v ON_ERROR_STOP=1 --username postgres -d task_state <<-EOSQL
        CREATE SCHEMA IF NOT EXISTS task;
        CREATE TABLE IF NOT EXISTS task.state (
          id SERIAL PRIMARY KEY,
          task_id UUID NOT NULL,
          agent_id UUID NOT NULL,
          state JSONB NOT NULL,
          status VARCHAR(50) NOT NULL,
          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
          updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );
        
        CREATE TABLE IF NOT EXISTS task.state_history (
          id SERIAL PRIMARY KEY,
          state_id INTEGER NOT NULL REFERENCES task.state(id),
          state JSONB NOT NULL,
          status VARCHAR(50) NOT NULL,
          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );
        
        -- Create rollback trigger
        CREATE OR REPLACE FUNCTION task.state_history_trigger()
        RETURNS TRIGGER AS \$\$
        BEGIN
          INSERT INTO task.state_history (state_id, state, status, created_at)
          VALUES (OLD.id, OLD.state, OLD.status, NOW());
          RETURN NEW;
        END;
        \$\$ LANGUAGE plpgsql;
        
        DROP TRIGGER IF EXISTS state_history_trigger ON task.state;
        CREATE TRIGGER state_history_trigger
        BEFORE UPDATE ON task.state
        FOR EACH ROW
        EXECUTE FUNCTION task.state_history_trigger();
      EOSQL
    EOT
    
    "setup-replica.sh" = <<-EOT
      set -e
      
      pg_basebackup -h ${var.name}-postgres-primary -U replicator -p 5432 -D /var/lib/postgresql/data -Fp -Xs -P -R
      
      cat > /var/lib/postgresql/data/recovery.conf <<-EOF
      standby_mode = 'on'
      primary_conninfo = 'host=${var.name}-postgres-primary port=5432 user=replicator password=${local.postgres_replication_password}'
      EOF
    EOT
  }
}

resource "kubernetes_stateful_set" "postgres_primary" {
  metadata {
    name      = "${var.name}-postgres-primary"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      "app.kubernetes.io/component" = "postgres-primary"
    })
  }

  spec {
    service_name = "${var.name}-postgres-headless"
    replicas     = 1
    
    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "supabase"
        "app.kubernetes.io/component" = "postgres-primary"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          "app.kubernetes.io/component" = "postgres-primary"
        })
        annotations = merge(var.annotations, local.kata_annotations)
      }
      
      spec {
        container {
          name  = "postgres"
          image = "postgres:${var.postgres_version}"
          
          env {
            name  = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_credentials.metadata[0].name
                key  = "POSTGRES_PASSWORD"
              }
            }
          }
          
          port {
            container_port = 5432
            name           = "postgres"
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/var/lib/postgresql/data"
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/docker-entrypoint-initdb.d/setup-primary.sh"
            sub_path   = "setup-primary.sh"
          }
          
          liveness_probe {
            exec {
              command = ["pg_isready", "-U", "postgres"]
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.postgres_config.metadata[0].name
            default_mode = "0755"
          }
        }
      }
    }
    
    volume_claim_template {
      metadata {
        name = "data"
      }
      
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            storage = var.postgres_storage_size
          }
        }
      }
    }
  }
}

resource "kubernetes_stateful_set" "postgres_replica" {
  metadata {
    name      = "${var.name}-postgres-replica"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      "app.kubernetes.io/component" = "postgres-replica"
    })
  }

  spec {
    service_name = "${var.name}-postgres-replica-headless"
    replicas     = local.postgres_replicas
    
    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "supabase"
        "app.kubernetes.io/component" = "postgres-replica"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          "app.kubernetes.io/component" = "postgres-replica"
        })
        annotations = merge(var.annotations, local.kata_annotations)
      }
      
      spec {
        init_container {
          name  = "init-replica"
          image = "postgres:${var.postgres_version}"
          
          command = ["/scripts/setup-replica.sh"]
          
          env {
            name  = "PGPASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_credentials.metadata[0].name
                key  = "POSTGRES_REPLICATION_PASSWORD"
              }
            }
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/var/lib/postgresql/data"
          }
          
          volume_mount {
            name       = "scripts"
            mount_path = "/scripts"
          }
        }
        
        container {
          name  = "postgres"
          image = "postgres:${var.postgres_version}"
          
          port {
            container_port = 5432
            name           = "postgres"
          }
          
          volume_mount {
            name       = "data"
            mount_path = "/var/lib/postgresql/data"
          }
          
          liveness_probe {
            exec {
              command = ["pg_isready", "-U", "postgres"]
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
        }
        
        volume {
          name = "scripts"
          config_map {
            name = kubernetes_config_map.postgres_config.metadata[0].name
            default_mode = "0755"
          }
        }
      }
    }
    
    volume_claim_template {
      metadata {
        name = "data"
      }
      
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            storage = var.postgres_storage_size
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "postgres_primary" {
  metadata {
    name      = "${var.name}-postgres-primary"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      "app.kubernetes.io/component" = "postgres-primary"
    })
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "supabase"
      "app.kubernetes.io/component" = "postgres-primary"
    }
    
    port {
      port        = 5432
      target_port = 5432
      name        = "postgres"
    }
  }
}

resource "kubernetes_service" "postgres_replica" {
  metadata {
    name      = "${var.name}-postgres-replica"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      "app.kubernetes.io/component" = "postgres-replica"
    })
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "supabase"
      "app.kubernetes.io/component" = "postgres-replica"
    }
    
    port {
      port        = 5432
      target_port = 5432
      name        = "postgres"
    }
  }
}

resource "kubernetes_deployment" "supabase" {
  count = local.instance_count
  
  metadata {
    name      = "${var.name}-${count.index}"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      "app.kubernetes.io/component" = "supabase"
      "app.kubernetes.io/instance"  = "${count.index}"
    })
  }

  spec {
    replicas = 1
    
    selector {
      match_labels = {
        "app.kubernetes.io/name"      = "supabase"
        "app.kubernetes.io/component" = "supabase"
        "app.kubernetes.io/instance"  = "${count.index}"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          "app.kubernetes.io/component" = "supabase"
          "app.kubernetes.io/instance"  = "${count.index}"
        })
        annotations = merge(var.annotations, local.kata_annotations)
      }
      
      spec {
        container {
          name  = "auth"
          image = "supabase/auth:latest"
          
          env {
            name  = "POSTGRES_HOST"
            value = count.index == 0 ? kubernetes_service.postgres_primary.metadata[0].name : kubernetes_service.postgres_replica.metadata[0].name
          }
          
          env {
            name  = "POSTGRES_PORT"
            value = "5432"
          }
          
          env {
            name  = "POSTGRES_DB"
            value = "postgres"
          }
          
          env {
            name  = "POSTGRES_USER"
            value = "postgres"
          }
          
          env {
            name  = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_credentials.metadata[0].name
                key  = "POSTGRES_PASSWORD"
              }
            }
          }
          
          env {
            name  = "JWT_SECRET"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_credentials.metadata[0].name
                key  = "JWT_SECRET"
              }
            }
          }
          
          port {
            container_port = 9999
            name           = "auth"
          }
        }
        
        container {
          name  = "rest"
          image = "postgrest/postgrest:latest"
          
          env {
            name  = "PGRST_DB_URI"
            value = "postgres://postgres:${local.postgres_password}@${count.index == 0 ? kubernetes_service.postgres_primary.metadata[0].name : kubernetes_service.postgres_replica.metadata[0].name}:5432/postgres"
          }
          
          env {
            name  = "PGRST_JWT_SECRET"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_credentials.metadata[0].name
                key  = "JWT_SECRET"
              }
            }
          }
          
          port {
            container_port = 3000
            name           = "rest"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "supabase" {
  metadata {
    name      = var.name
    namespace = local.namespace
    labels    = local.common_labels
  }

  spec {
    selector = {
      "app.kubernetes.io/name"      = "supabase"
      "app.kubernetes.io/component" = "supabase"
    }
    
    port {
      port        = 80
      target_port = 3000
      name        = "rest"
    }
    
    port {
      port        = 9999
      target_port = 9999
      name        = "auth"
    }
  }
}
