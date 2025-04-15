resource "kubernetes_stateful_set" "kafka" {
  metadata {
    name      = "kafka"
    namespace = var.namespace
    labels = {
      app     = "kafka"
      service = "messaging"
    }
  }

  spec {
    service_name = "kafka"
    replicas     = var.kafka_replicas

    selector {
      match_labels = {
        app = "kafka"
      }
    }

    template {
      metadata {
        labels = {
          app     = "kafka"
          service = "messaging"
        }
      }

      spec {
        container {
          name  = "kafka"
          image = "bitnami/kafka:${var.kafka_version}"

          port {
            container_port = 9092
            name           = "kafka"
          }

          env {
            name  = "KAFKA_CFG_ZOOKEEPER_CONNECT"
            value = "${kubernetes_service.zookeeper.metadata[0].name}:2181"
          }

          env {
            name  = "KAFKA_CFG_ADVERTISED_LISTENERS"
            value = "PLAINTEXT://${kubernetes_service.kafka.metadata[0].name}.${var.namespace}.svc.cluster.local:9092"
          }

          env {
            name  = "KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP"
            value = "PLAINTEXT:PLAINTEXT"
          }

          env {
            name  = "KAFKA_CFG_LISTENERS"
            value = "PLAINTEXT://:9092"
          }

          env {
            name  = "KAFKA_CFG_INTER_BROKER_LISTENER_NAME"
            value = "PLAINTEXT"
          }

          env {
            name  = "KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE"
            value = "true"
          }

          env {
            name  = "KAFKA_HEAP_OPTS"
            value = "-Xmx${var.kafka_heap_size} -Xms${var.kafka_heap_size}"
          }

          volume_mount {
            name       = "kafka-data"
            mount_path = "/bitnami/kafka"
          }

          resources {
            requests = {
              memory = var.kafka_memory_request
              cpu    = var.kafka_cpu_request
            }
            limits = {
              memory = var.kafka_memory_limit
              cpu    = var.kafka_cpu_limit
            }
          }

          readiness_probe {
            tcp_socket {
              port = 9092
            }
            initial_delay_seconds = 30
            period_seconds        = 10
            timeout_seconds       = 5
            success_threshold     = 1
            failure_threshold     = 3
          }

          liveness_probe {
            tcp_socket {
              port = 9092
            }
            initial_delay_seconds = 60
            period_seconds        = 20
            timeout_seconds       = 10
            success_threshold     = 1
            failure_threshold     = 6
          }
        }
      }
    }

    volume_claim_template {
      metadata {
        name = "kafka-data"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            storage = var.kafka_storage_size
          }
        }
      }
    }
  }
}

resource "kubernetes_deployment" "zookeeper" {
  metadata {
    name      = "zookeeper"
    namespace = var.namespace
    labels = {
      app     = "zookeeper"
      service = "messaging"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "zookeeper"
      }
    }

    template {
      metadata {
        labels = {
          app     = "zookeeper"
          service = "messaging"
        }
      }

      spec {
        container {
          name  = "zookeeper"
          image = "bitnami/zookeeper:${var.zookeeper_version}"

          port {
            container_port = 2181
            name           = "client"
          }

          port {
            container_port = 2888
            name           = "server"
          }

          port {
            container_port = 3888
            name           = "leader-election"
          }

          env {
            name  = "ALLOW_ANONYMOUS_LOGIN"
            value = "yes"
          }

          volume_mount {
            name       = "zookeeper-data"
            mount_path = "/bitnami/zookeeper"
          }

          resources {
            requests = {
              memory = var.zookeeper_memory_request
              cpu    = var.zookeeper_cpu_request
            }
            limits = {
              memory = var.zookeeper_memory_limit
              cpu    = var.zookeeper_cpu_limit
            }
          }

          readiness_probe {
            tcp_socket {
              port = 2181
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }

          liveness_probe {
            tcp_socket {
              port = 2181
            }
            initial_delay_seconds = 60
            period_seconds        = 20
          }
        }
      }
    }
  }

  volume {
    name = "zookeeper-data"
    persistent_volume_claim {
      claim_name = kubernetes_persistent_volume_claim.zookeeper_data.metadata[0].name
    }
  }
}

resource "kubernetes_persistent_volume_claim" "zookeeper_data" {
  metadata {
    name      = "zookeeper-data"
    namespace = var.namespace
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.zookeeper_storage_size
      }
    }
  }
}

resource "kubernetes_service" "kafka" {
  metadata {
    name      = "kafka"
    namespace = var.namespace
    labels = {
      app     = "kafka"
      service = "messaging"
    }
  }

  spec {
    selector = {
      app = "kafka"
    }

    port {
      port        = 9092
      target_port = 9092
      protocol    = "TCP"
      name        = "kafka"
    }
  }
}

resource "kubernetes_service" "zookeeper" {
  metadata {
    name      = "zookeeper"
    namespace = var.namespace
    labels = {
      app     = "zookeeper"
      service = "messaging"
    }
  }

  spec {
    selector = {
      app = "zookeeper"
    }

    port {
      port        = 2181
      target_port = 2181
      protocol    = "TCP"
      name        = "client"
    }

    port {
      port        = 2888
      target_port = 2888
      protocol    = "TCP"
      name        = "server"
    }

    port {
      port        = 3888
      target_port = 3888
      protocol    = "TCP"
      name        = "leader-election"
    }
  }
}
