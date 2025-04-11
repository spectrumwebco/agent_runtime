
resource "kubernetes_namespace" "rocketmq" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = var.namespace
    
    labels = merge({
      "app.kubernetes.io/name"       = "rocketmq"
      "app.kubernetes.io/part-of"    = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }, var.labels)
    
    annotations = var.annotations
  }
}

locals {
  namespace = var.create_namespace ? kubernetes_namespace.rocketmq[0].metadata[0].name : var.namespace
  
  common_labels = merge({
    "app.kubernetes.io/name"       = "rocketmq"
    "app.kubernetes.io/part-of"    = "agent-runtime"
    "app.kubernetes.io/managed-by" = "terraform"
  }, var.labels)
  
  kata_annotations = var.kata_container_integration ? {
    "io.kubernetes.cri.untrusted-workload" = "true"
    "io.kubernetes.cri.runtimeclass"       = "kata"
  } : {}
}

resource "kubernetes_config_map" "broker_config" {
  metadata {
    name      = "${var.name}-broker-config"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = merge({
    "broker.conf" = <<-EOT
      brokerClusterName = AgentRuntimeCluster
      brokerName = broker-${var.name}
      brokerId = 0
      deleteWhen = 04
      fileReservedTime = 48
      brokerRole = ASYNC_MASTER
      flushDiskType = ASYNC_FLUSH
      autoCreateTopicEnable = true
      autoCreateSubscriptionGroup = true
      messageDelayLevel = 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
      enablePropertyFilter = true
      traceTopicEnable = true
      transactionCheckInterval = 60000
      transactionCheckMax = 15
      transactionTimeOut = 6000
      aclEnable = ${var.acl_enabled}
      storePathRootDir = /data/store
      storePathCommitLog = /data/store/commitlog
      storePathConsumerQueue = /data/store/consumequeue
      storePathIndex = /data/store/index
      storeCheckpoint = /data/store/checkpoint
      abortFile = /data/store/abort
      maxMessageSize = 4194304
      flushCommitLogTimed = true
      flushCommitLogLeastPages = 4
      flushCommitLogThoroughInterval = 10000
      flushConsumeQueueThoroughInterval = 60000
      brokerIP1 = ${var.name}-broker-0.${var.name}-broker-headless.${local.namespace}.svc.cluster.local
      listenPort = 10911
      haListenPort = 10912
      EOT
  }, var.broker_config)
}

resource "kubernetes_config_map" "nameserver_config" {
  metadata {
    name      = "${var.name}-nameserver-config"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = merge({
    "nameserver.conf" = <<-EOT
      listenPort=9876
      serverWorkerThreads=8
      serverCallbackExecutorThreads=0
      serverSelectorThreads=3
      serverOnewaySemaphoreValue=256
      serverAsyncSemaphoreValue=64
      serverChannelMaxIdleTimeSeconds=120
      serverSocketSndBufSize=65535
      serverSocketRcvBufSize=65535
      serverPooledByteBufAllocatorEnable=true
      useEpollNativeSelector=false
      EOT
  }, var.name_server_config)
}

resource "kubernetes_secret" "rocketmq_acl" {
  count = var.acl_enabled ? 1 : 0
  
  metadata {
    name      = "${var.name}-acl"
    namespace = local.namespace
    labels    = local.common_labels
  }
  
  data = {
    "plain_acl.yml" = <<-EOT
      globalWhiteRemoteAddresses:
        - 10.0.0.0/8
        - 172.16.0.0/12
        - 192.168.0.0/16
      accounts:
        - accessKey: ${var.acl_access_key}
          secretKey: ${var.acl_secret_key}
          whiteRemoteAddress:
          admin: true
          defaultTopicPerm: DENY
          defaultGroupPerm: SUB
          topicPerms:
            - agent-events=DENY
            - k8s-lifecycle=PUB|SUB
            - kata-lifecycle=PUB|SUB
            - state-updates=PUB|SUB
          groupPerms:
            - agent-consumers=SUB
            - k8s-consumers=SUB
            - kata-consumers=SUB
            - state-consumers=SUB
      EOT
  }
}

resource "kubernetes_stateful_set" "nameserver" {
  metadata {
    name      = "${var.name}-nameserver"
    namespace = local.namespace
    labels    = local.common_labels
  }

  spec {
    service_name = "${var.name}-nameserver-headless"
    replicas     = var.high_availability ? var.name_server_replicas : 1
    
    selector {
      match_labels = {
        app     = "rocketmq"
        component = "nameserver"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          app       = "rocketmq"
          component = "nameserver"
        })
        
        annotations = merge(var.annotations, local.kata_annotations, {
          "prometheus.io/scrape" = var.prometheus_integration ? "true" : "false"
          "prometheus.io/port"   = "9876"
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
          name  = "nameserver"
          image = "${var.image_repository}:${var.image_tag}"
          args  = ["sh", "-c", "cd /opt/rocketmq-${var.image_tag}/bin && ./mqnamesrv -c /etc/rocketmq/nameserver.conf"]
          
          port {
            name           = "namesrv"
            container_port = 9876
          }
          
          resources {
            limits = {
              cpu    = var.resource_limits.name_server.cpu
              memory = var.resource_limits.name_server.memory
            }
            requests = {
              cpu    = var.resource_requests.name_server.cpu
              memory = var.resource_requests.name_server.memory
            }
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/rocketmq"
          }
          
          volume_mount {
            name       = "logs"
            mount_path = "/opt/rocketmq-${var.image_tag}/logs"
          }
          
          liveness_probe {
            tcp_socket {
              port = 9876
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            tcp_socket {
              port = 9876
            }
            initial_delay_seconds = 15
            period_seconds        = 5
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.nameserver_config.metadata[0].name
            items {
              key  = "nameserver.conf"
              path = "nameserver.conf"
            }
          }
        }
        
        volume {
          name = "logs"
          empty_dir {}
        }
      }
    }
  }
  
  depends_on = [kubernetes_namespace.rocketmq]
}

resource "kubernetes_service" "nameserver_headless" {
  metadata {
    name      = "${var.name}-nameserver-headless"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "nameserver"
    })
  }
  
  spec {
    selector = {
      app       = "rocketmq"
      component = "nameserver"
    }
    
    port {
      name        = "namesrv"
      port        = 9876
      target_port = 9876
    }
    
    cluster_ip = "None"
  }
  
  depends_on = [kubernetes_namespace.rocketmq]
}

resource "kubernetes_service" "nameserver" {
  metadata {
    name      = "${var.name}-nameserver"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "nameserver"
    })
  }
  
  spec {
    selector = {
      app       = "rocketmq"
      component = "nameserver"
    }
    
    port {
      name        = "namesrv"
      port        = 9876
      target_port = 9876
    }
    
    type = "ClusterIP"
  }
  
  depends_on = [kubernetes_namespace.rocketmq]
}

resource "kubernetes_stateful_set" "broker" {
  metadata {
    name      = "${var.name}-broker"
    namespace = local.namespace
    labels    = local.common_labels
  }

  spec {
    service_name = "${var.name}-broker-headless"
    replicas     = var.high_availability ? var.broker_replicas : 1
    
    selector {
      match_labels = {
        app       = "rocketmq"
        component = "broker"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          app       = "rocketmq"
          component = "broker"
        })
        
        annotations = merge(var.annotations, local.kata_annotations, {
          "prometheus.io/scrape" = var.prometheus_integration ? "true" : "false"
          "prometheus.io/port"   = "10911"
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
          name  = "broker"
          image = "${var.image_repository}:${var.image_tag}"
          args  = ["sh", "-c", "cd /opt/rocketmq-${var.image_tag}/bin && ./mqbroker -c /etc/rocketmq/broker.conf -n ${var.name}-nameserver.${local.namespace}.svc.cluster.local:9876"]
          
          port {
            name           = "broker"
            container_port = 10911
          }
          
          port {
            name           = "haport"
            container_port = 10912
          }
          
          resources {
            limits = {
              cpu    = var.resource_limits.broker.cpu
              memory = var.resource_limits.broker.memory
            }
            requests = {
              cpu    = var.resource_requests.broker.cpu
              memory = var.resource_requests.broker.memory
            }
          }
          
          volume_mount {
            name       = "config"
            mount_path = "/etc/rocketmq"
          }
          
          volume_mount {
            name       = "store"
            mount_path = "/data/store"
          }
          
          volume_mount {
            name       = "logs"
            mount_path = "/opt/rocketmq-${var.image_tag}/logs"
          }
          
          dynamic "volume_mount" {
            for_each = var.acl_enabled ? [1] : []
            content {
              name       = "acl"
              mount_path = "/opt/rocketmq-${var.image_tag}/conf/plain_acl.yml"
              sub_path   = "plain_acl.yml"
            }
          }
          
          liveness_probe {
            tcp_socket {
              port = 10911
            }
            initial_delay_seconds = 60
            period_seconds        = 15
          }
          
          readiness_probe {
            tcp_socket {
              port = 10911
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.broker_config.metadata[0].name
            items {
              key  = "broker.conf"
              path = "broker.conf"
            }
          }
        }
        
        volume {
          name = "logs"
          empty_dir {}
        }
        
        dynamic "volume" {
          for_each = var.acl_enabled ? [1] : []
          content {
            name = "acl"
            secret {
              secret_name = kubernetes_secret.rocketmq_acl[0].metadata[0].name
            }
          }
        }
      }
    }
    
    volume_claim_template {
      metadata {
        name = "store"
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
  }
  
  depends_on = [
    kubernetes_namespace.rocketmq,
    kubernetes_stateful_set.nameserver,
    kubernetes_service.nameserver
  ]
}

resource "kubernetes_service" "broker_headless" {
  metadata {
    name      = "${var.name}-broker-headless"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "broker"
    })
  }
  
  spec {
    selector = {
      app       = "rocketmq"
      component = "broker"
    }
    
    port {
      name        = "broker"
      port        = 10911
      target_port = 10911
    }
    
    port {
      name        = "haport"
      port        = 10912
      target_port = 10912
    }
    
    cluster_ip = "None"
  }
  
  depends_on = [kubernetes_namespace.rocketmq]
}

resource "kubernetes_service" "broker" {
  metadata {
    name      = "${var.name}-broker"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "broker"
    })
  }
  
  spec {
    selector = {
      app       = "rocketmq"
      component = "broker"
    }
    
    port {
      name        = "broker"
      port        = 10911
      target_port = 10911
    }
    
    port {
      name        = "haport"
      port        = 10912
      target_port = 10912
    }
    
    type = "ClusterIP"
  }
  
  depends_on = [kubernetes_namespace.rocketmq]
}

resource "kubernetes_deployment" "dashboard" {
  count = var.dashboard_enabled ? 1 : 0
  
  metadata {
    name      = "${var.name}-dashboard"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "dashboard"
    })
  }
  
  spec {
    replicas = var.high_availability ? 2 : 1
    
    selector {
      match_labels = {
        app       = "rocketmq"
        component = "dashboard"
      }
    }
    
    template {
      metadata {
        labels = merge(local.common_labels, {
          app       = "rocketmq"
          component = "dashboard"
        })
      }
      
      spec {
        container {
          name  = "dashboard"
          image = "apacherocketmq/rocketmq-dashboard:latest"
          
          port {
            name           = "http"
            container_port = 8080
          }
          
          env {
            name  = "JAVA_OPTS"
            value = "-Drocketmq.namesrv.addr=${var.name}-nameserver.${local.namespace}.svc.cluster.local:9876 -Dserver.port=8080"
          }
          
          resources {
            limits = {
              cpu    = var.resource_limits.dashboard.cpu
              memory = var.resource_limits.dashboard.memory
            }
            requests = {
              cpu    = var.resource_requests.dashboard.cpu
              memory = var.resource_requests.dashboard.memory
            }
          }
          
          liveness_probe {
            http_get {
              path = "/"
              port = 8080
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            http_get {
              path = "/"
              port = 8080
            }
            initial_delay_seconds = 15
            period_seconds        = 5
          }
        }
      }
    }
  }
  
  depends_on = [
    kubernetes_namespace.rocketmq,
    kubernetes_stateful_set.nameserver,
    kubernetes_service.nameserver
  ]
}

resource "kubernetes_service" "dashboard" {
  count = var.dashboard_enabled ? 1 : 0
  
  metadata {
    name      = "${var.name}-dashboard"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "dashboard"
    })
  }
  
  spec {
    selector = {
      app       = "rocketmq"
      component = "dashboard"
    }
    
    port {
      name        = "http"
      port        = 8080
      target_port = 8080
    }
    
    type = "ClusterIP"
  }
  
  depends_on = [kubernetes_namespace.rocketmq]
}

resource "kubernetes_ingress_v1" "dashboard" {
  count = var.dashboard_enabled && var.ingress_enabled ? 1 : 0
  
  metadata {
    name      = "${var.name}-dashboard"
    namespace = local.namespace
    labels    = merge(local.common_labels, {
      component = "dashboard"
    })
    
    annotations = {
      "kubernetes.io/ingress.class" = var.ingress_class
    }
  }
  
  spec {
    rule {
      host = "rocketmq-dashboard.${var.ingress_domain}"
      
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          
          backend {
            service {
              name = kubernetes_service.dashboard[0].metadata[0].name
              
              port {
                number = 8080
              }
            }
          }
        }
      }
    }
    
    dynamic "tls" {
      for_each = var.ingress_tls_enabled ? [1] : []
      content {
        hosts       = ["rocketmq-dashboard.${var.ingress_domain}"]
        secret_name = var.ingress_tls_secret
      }
    }
  }
  
  depends_on = [kubernetes_service.dashboard]
}

resource "kubernetes_config_map" "topic_creation" {
  metadata {
    name      = "${var.name}-topic-creation"
    namespace = local.namespace
    labels    = local.common_labels
  }

  data = {
    "create-topics.sh" = <<-EOT
      set -e
      
      NAMESRV_ADDR=${var.name}-nameserver.${local.namespace}.svc.cluster.local:9876
      
      echo "Waiting for name server to be ready..."
      until nc -z ${var.name}-nameserver.${local.namespace}.svc.cluster.local 9876; do
        sleep 5
      done
      
      echo "Waiting for broker to be ready..."
      until nc -z ${var.name}-broker.${local.namespace}.svc.cluster.local 10911; do
        sleep 5
      done
      
      cd /opt/rocketmq-${var.image_tag}/bin
      
      %{for topic in var.topic_configs}
      echo "Creating topic ${topic.name}..."
      ./mqadmin updateTopic -n $NAMESRV_ADDR -t ${topic.name} -c AgentRuntimeCluster -r ${topic.read_queue_nums} -w ${topic.write_queue_nums} -p ${topic.perm}
      %{endfor}
      
      echo "All topics created successfully!"
    EOT
  }
}

resource "kubernetes_job" "topic_creation" {
  metadata {
    name      = "${var.name}-topic-creation"
    namespace = local.namespace
    labels    = local.common_labels
  }

  spec {
    template {
      metadata {
        labels = local.common_labels
      }
      
      spec {
        container {
          name  = "topic-creation"
          image = "${var.image_repository}:${var.image_tag}"
          command = ["/bin/sh", "-c", "chmod +x /scripts/create-topics.sh && /scripts/create-topics.sh"]
          
          volume_mount {
            name       = "scripts"
            mount_path = "/scripts"
          }
          
          dynamic "volume_mount" {
            for_each = var.acl_enabled ? [1] : []
            content {
              name       = "acl"
              mount_path = "/opt/rocketmq-${var.image_tag}/conf/plain_acl.yml"
              sub_path   = "plain_acl.yml"
            }
          }
        }
        
        volume {
          name = "scripts"
          config_map {
            name = kubernetes_config_map.topic_creation.metadata[0].name
            default_mode = "0755"
          }
        }
        
        dynamic "volume" {
          for_each = var.acl_enabled ? [1] : []
          content {
            name = "acl"
            secret {
              secret_name = kubernetes_secret.rocketmq_acl[0].metadata[0].name
            }
          }
        }
        
        restart_policy = "OnFailure"
      }
    }
    
    backoff_limit = 5
  }
  
  depends_on = [
    kubernetes_stateful_set.broker,
    kubernetes_service.broker,
    kubernetes_stateful_set.nameserver,
    kubernetes_service.nameserver
  ]
}

resource "kubernetes_manifest" "service_monitor" {
  count = var.prometheus_integration ? 1 : 0
  
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"
    
    metadata = {
      name      = "${var.name}-service-monitor"
      namespace = local.namespace
      labels    = local.common_labels
    }
    
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/name" = "rocketmq"
        }
      }
      
      endpoints = [
        {
          port     = "namesrv"
          interval = "30s"
          path     = "/metrics"
        },
        {
          port     = "broker"
          interval = "30s"
          path     = "/metrics"
        }
      ]
    }
  }
  
  depends_on = [
    kubernetes_service.nameserver,
    kubernetes_service.broker
  ]
}
