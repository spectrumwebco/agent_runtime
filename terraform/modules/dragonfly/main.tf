# DragonflyDB Module for Agent Runtime

resource "kubernetes_namespace" "dragonfly" {
  count = var.namespace != "default" ? 1 : 0

  metadata {
    name = var.namespace
    
    labels = merge({
      "app.kubernetes.io/name"       = "dragonfly"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
    
    annotations = var.annotations
  }
}

locals {
  namespace = var.namespace != "default" ? kubernetes_namespace.dragonfly[0].metadata[0].name : var.namespace
  
  common_labels = merge({
    "app.kubernetes.io/name"       = "dragonfly"
    "app.kubernetes.io/part-of"    = "agent-runtime"
    "app.kubernetes.io/managed-by" = "terraform"
  }, var.labels)
  
  kata_annotations = var.kata_container_integration ? {
    "io.kubernetes.cri.untrusted-workload" = "true"
    "io.kubernetes.cri.runtimeclass"       = "kata"
  } : {}
}

resource "kubernetes_config_map" "dragonfly_config" {
  metadata {
    name      = "dragonfly-config"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "dragonfly.conf" = <<-EOT
      bind 0.0.0.0
      port ${var.port}
      
      maxmemory ${var.maxmemory_percent}%
      maxmemory-policy ${var.maxmemory_policy}
      
      ${var.persistence_enabled ? "dir ${var.persistence_path}" : "# Persistence disabled"}
      ${var.persistence_enabled ? "dbfilename dump.rdb" : ""}
      ${var.persistence_enabled ? "save 900 1" : ""}
      ${var.persistence_enabled ? "save 300 10" : ""}
      ${var.persistence_enabled ? "save 60 10000" : ""}
      
      ${var.high_availability ? "replicaof dragonfly-0.dragonfly-headless.${local.namespace}.svc.cluster.local ${var.port}" : "# Replication disabled"}
      
      ${var.password_enabled ? "requirepass ${var.dragonfly_password}" : "# Password authentication disabled"}
      ${var.password_enabled ? "masterauth ${var.dragonfly_password}" : ""}
      
      ${var.event_stream_integration ? "notify-keyspace-events AKE" : "# Event notifications disabled"}
      
      ${var.cache_invalidation_enabled ? "lazyfree-lazy-eviction yes" : ""}
      ${var.cache_invalidation_enabled ? "lazyfree-lazy-expire yes" : ""}
      
      ${var.cache_ttl > 0 ? "expire-default-ttl ${var.cache_ttl}" : ""}
      
      ${join("\n      ", [for k, v in var.config_params : "${k} ${v}"])}
    EOT
  }
}

resource "kubernetes_stateful_set" "dragonfly" {
  metadata {
    name      = "dragonfly"
    namespace = local.namespace
    labels    = local.common_labels
  }

  spec {
    service_name = "dragonfly-headless"
    replicas     = var.high_availability ? var.replicas : 1
    
    selector {
      match_labels = {
        app = "dragonfly"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          app = "dragonfly"
        })
        
        annotations = merge(var.annotations, local.kata_annotations, {
          "prometheus.io/scrape" = var.prometheus_integration ? "true" : "false"
          "prometheus.io/port"   = "${var.port}"
        })
      }
      
      spec {
        dynamic "toleration" {
          for_each = var.kata_container_integration ? [1] : []
          content {
            key      = "kata"
            operator = "Exists"
            effect   = "NoSchedule"
          }
        }
        
        container {
          name  = "dragonfly"
          image = "${var.image_repository}:${var.image_tag}"
          args  = concat(["--config", "/etc/dragonfly/dragonfly.conf"], var.additional_args)
          
          port {
            name           = "redis"
            container_port = var.port
          }
          
          port {
            name           = "gossip"
            container_port = var.gossip_port
          }
          
          resources {
            limits = {
              cpu    = var.resource_limits.cpu
              memory = var.resource_limits.memory
            }
            requests = {
              cpu    = var.resource_requests.cpu
              memory = var.resource_requests.memory
            }
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/dragonfly"
          }
          
          volume_mount {
            name       = "data"
            mount_path = var.persistence_path
          }
          
          liveness_probe {
            exec {
              command = [
                "sh",
                "-c",
                "${var.password_enabled ? "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} ping" : "redis-cli -h localhost -p ${var.port} ping"}"
              ]
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            exec {
              command = [
                "sh",
                "-c",
                "${var.password_enabled ? "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} ping" : "redis-cli -h localhost -p ${var.port} ping"}"
              ]
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.dragonfly_config.metadata[0].name
            items {
              key  = "dragonfly.conf"
              path = "dragonfly.conf"
            }
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
            storage = var.storage_size
          }
        }
        
        storage_class_name = var.storage_class
      }
    }
    
    update_strategy {
      type = "RollingUpdate"
      
      rolling_update {
        partition = 0
      }
    }
  }
  
  depends_on = [kubernetes_namespace.dragonfly]
}

resource "kubernetes_service" "dragonfly_headless" {
  metadata {
    name      = "dragonfly-headless"
    namespace = local.namespace
    labels    = local.common_labels
  }
  
  spec {
    selector = {
      app = "dragonfly"
    }
    
    port {
      name        = "redis"
      port        = var.port
      target_port = var.port
    }
    
    port {
      name        = "gossip"
      port        = var.gossip_port
      target_port = var.gossip_port
    }
    
    cluster_ip = "None"
  }
  
  depends_on = [kubernetes_namespace.dragonfly]
}

resource "kubernetes_service" "dragonfly" {
  metadata {
    name      = "dragonfly"
    namespace = local.namespace
    labels    = local.common_labels
  }
  
  spec {
    selector = {
      app = "dragonfly"
    }
    
    port {
      name        = "redis"
      port        = var.port
      target_port = var.port
    }
    
    type = "ClusterIP"
  }
  
  depends_on = [kubernetes_namespace.dragonfly]
}

resource "kubernetes_cron_job_v1" "dragonfly_snapshot" {
  count = var.snapshot_enabled ? 1 : 0
  
  metadata {
    name      = "dragonfly-snapshot"
    namespace = local.namespace
    labels    = local.common_labels
  }
  
  spec {
    schedule                      = var.snapshot_schedule
    concurrency_policy            = "Forbid"
    successful_jobs_history_limit = var.snapshot_retention
    failed_jobs_history_limit     = 3
    
    job_template {
      metadata {
        labels = local.common_labels
      }
      
      spec {
        template {
          metadata {
            labels = local.common_labels
          }
          
          spec {
            container {
              name  = "snapshot"
              image = "redis:alpine"
              
              command = [
                "sh",
                "-c",
                "${var.password_enabled ? 
                  "redis-cli -h dragonfly.${local.namespace}.svc.cluster.local -p ${var.port} -a ${var.dragonfly_password} SAVE" : 
                  "redis-cli -h dragonfly.${local.namespace}.svc.cluster.local -p ${var.port} SAVE"}"
              ]
            }
            
            restart_policy = "OnFailure"
          }
        }
      }
    }
  }
  
  depends_on = [
    kubernetes_stateful_set.dragonfly,
    kubernetes_service.dragonfly
  ]
}

resource "kubernetes_config_map" "dragonfly_rollback" {
  count = var.rollback_enabled ? 1 : 0
  
  metadata {
    name      = "dragonfly-rollback"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "rollback.sh" = <<-EOT
      set -e
      
      
      SNAPSHOTS=$(ls -1t ${var.persistence_path}/backup/ | grep -E '^dragonfly-snapshot-[0-9]+\.rdb$')
      
      if [ -z "$SNAPSHOTS" ]; then
        echo "No snapshots available for rollback"
        exit 1
      fi
      
      SNAPSHOT_INDEX=${1:-1}
      SNAPSHOT=$(echo "$SNAPSHOTS" | sed -n "${SNAPSHOT_INDEX}p")
      
      if [ -z "$SNAPSHOT" ]; then
        echo "Snapshot index $SNAPSHOT_INDEX not found"
        exit 1
      fi
      
      echo "Rolling back to snapshot: $SNAPSHOT"
      
      ${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} SHUTDOWN SAVE" : 
        "redis-cli -h localhost -p ${var.port} SHUTDOWN SAVE"}
      
      while pgrep -f "dragonfly" > /dev/null; do
        echo "Waiting for DragonflyDB to stop..."
        sleep 1
      done
      
      cp ${var.persistence_path}/dump.rdb ${var.persistence_path}/dump.rdb.bak
      
      cp ${var.persistence_path}/backup/$SNAPSHOT ${var.persistence_path}/dump.rdb
      
      dragonfly --config /etc/dragonfly/dragonfly.conf &
      
      echo "Rollback completed successfully"
    EOT
    
    "create-snapshot.sh" = <<-EOT
      set -e
      
      
      mkdir -p ${var.persistence_path}/backup
      
      TIMESTAMP=$(date +%Y%m%d%H%M%S)
      
      ${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} SAVE" : 
        "redis-cli -h localhost -p ${var.port} SAVE"}
      
      cp ${var.persistence_path}/dump.rdb ${var.persistence_path}/backup/dragonfly-snapshot-$TIMESTAMP.rdb
      
      ls -1t ${var.persistence_path}/backup/dragonfly-snapshot-*.rdb | tail -n +${var.snapshot_retention + 1} | xargs -r rm
      
      echo "Snapshot created: dragonfly-snapshot-$TIMESTAMP.rdb"
    EOT
    
    "list-snapshots.sh" = <<-EOT
      
      
      if [ ! -d "${var.persistence_path}/backup" ]; then
        echo "No snapshots available"
        exit 0
      fi
      
      SNAPSHOTS=$(ls -1t ${var.persistence_path}/backup/ | grep -E '^dragonfly-snapshot-[0-9]+\.rdb$')
      
      if [ -z "$SNAPSHOTS" ]; then
        echo "No snapshots available"
        exit 0
      fi
      
      echo "Available snapshots:"
      echo "$SNAPSHOTS" | nl
    EOT
  }
}

resource "kubernetes_manifest" "service_monitor" {
  count = var.prometheus_integration ? 1 : 0
  
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"
    
    metadata = {
      name      = "dragonfly-service-monitor"
      namespace = local.namespace
      labels    = local.common_labels
    }
    
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/name" = "dragonfly"
        }
      }
      
      endpoints = [
        {
          port     = "redis"
          interval = "30s"
          path     = "/metrics"
        }
      ]
    }
  }
  
  depends_on = [kubernetes_service.dragonfly]
}

resource "kubernetes_config_map" "event_stream_integration" {
  count = var.event_stream_integration ? 1 : 0
  
  metadata {
    name      = "dragonfly-event-stream"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "event-stream-integration.sh" = <<-EOT
      set -e
      
      
      ${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} CONFIG SET notify-keyspace-events AKE" : 
        "redis-cli -h localhost -p ${var.port} CONFIG SET notify-keyspace-events AKE"}
      
      ${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} CONFIG SET expire-default-ttl ${var.cache_ttl}" : 
        "redis-cli -h localhost -p ${var.port} CONFIG SET expire-default-ttl ${var.cache_ttl}"}
      
      ${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} CONFIG SET context-cache-ttl ${var.context_cache_ttl}" : 
        "redis-cli -h localhost -p ${var.port} CONFIG SET context-cache-ttl ${var.context_cache_ttl}"}
      
      echo "Event Stream integration configured successfully"
    EOT
    
    "cache-invalidation.sh" = <<-EOT
      set -e
      
      
      if [ -z "$1" ]; then
        echo "Usage: $0 <pattern>"
        exit 1
      fi
      
      PATTERN=$1
      
      KEYS=$(${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} KEYS \"$PATTERN\"" : 
        "redis-cli -h localhost -p ${var.port} KEYS \"$PATTERN\""})
      
      if [ -z "$KEYS" ]; then
        echo "No keys found matching pattern: $PATTERN"
        exit 0
      fi
      
      echo "$KEYS" | xargs -I{} ${var.password_enabled ? 
        "redis-cli -h localhost -p ${var.port} -a ${var.dragonfly_password} DEL {}" : 
        "redis-cli -h localhost -p ${var.port} DEL {}"}
      
      echo "Cache invalidation completed for pattern: $PATTERN"
    EOT
  }
}

resource "kubernetes_config_map" "context_rebuild_registry" {
  count = var.event_stream_integration && var.cache_invalidation_enabled ? 1 : 0
  
  metadata {
    name      = "dragonfly-context-rebuild"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "context-rebuild-registry.json" = <<-EOT
      {
        "rebuild_functions": [
          {
            "pattern": "agent:context:*",
            "rebuild_function": "agent.RebuildAgentContext",
            "ttl": ${var.context_cache_ttl}
          },
          {
            "pattern": "tool:state:*",
            "rebuild_function": "tools.RebuildToolState",
            "ttl": ${var.cache_ttl}
          },
          {
            "pattern": "execution:state:*",
            "rebuild_function": "execution.RebuildExecutionState",
            "ttl": ${var.cache_ttl}
          }
        ]
      }
    EOT
  }
}
