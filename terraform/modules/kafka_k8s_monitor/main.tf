
resource "kubernetes_namespace" "monitoring" {
  metadata {
    name = "monitoring"
  }
}

resource "kubernetes_deployment" "kafka" {
  metadata {
    name      = "kafka"
    namespace = var.namespace
    labels = {
      app = "kafka"
    }
  }

  spec {
    replicas = var.kafka_replicas

    selector {
      match_labels = {
        app = "kafka"
      }
    }

    template {
      metadata {
        labels = {
          app = "kafka"
        }
      }

      spec {
        container {
          name  = "zookeeper"
          image = "wurstmeister/zookeeper:3.4.6"
          
          port {
            container_port = 2181
          }
          
          resources {
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
        
        container {
          name  = "kafka"
          image = "wurstmeister/kafka:2.13-2.7.0"
          
          port {
            container_port = 9092
          }
          
          env {
            name  = "KAFKA_ADVERTISED_HOST_NAME"
            value = "kafka.${var.namespace}.svc.cluster.local"
          }
          
          env {
            name  = "KAFKA_ADVERTISED_PORT"
            value = "9092"
          }
          
          env {
            name  = "KAFKA_ZOOKEEPER_CONNECT"
            value = "localhost:2181"
          }
          
          env {
            name  = "KAFKA_CREATE_TOPICS"
            value = "k8s-events:3:1,k8s-pods:3:1,k8s-deployments:3:1,k8s-services:3:1,k8s-configmaps:3:1,k8s-secrets:3:1,shared-state:3:1,event-stream:3:1"
          }
          
          resources {
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
      }
    }
  }
}

resource "kubernetes_service" "kafka" {
  metadata {
    name      = "kafka"
    namespace = var.namespace
  }

  spec {
    selector = {
      app = "kafka"
    }

    port {
      port        = 9092
      target_port = 9092
      name        = "kafka"
    }

    port {
      port        = 2181
      target_port = 2181
      name        = "zookeeper"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_deployment" "k8s_monitor" {
  metadata {
    name      = "k8s-monitor"
    namespace = var.namespace
    labels = {
      app = "k8s-monitor"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "k8s-monitor"
      }
    }

    template {
      metadata {
        labels = {
          app = "k8s-monitor"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.k8s_monitor.metadata[0].name
        
        container {
          name  = "k8s-monitor"
          image = "python:3.12-slim"
          
          command = ["/bin/bash", "-c"]
          args    = ["pip install kubernetes kafka-python && python -m django_app.manage start_k8s_monitor --daemon"]
          
          env {
            name  = "KAFKA_BOOTSTRAP_SERVERS"
            value = "kafka.${var.namespace}.svc.cluster.local:9092"
          }
          
          env {
            name  = "NAMESPACE"
            value = var.monitor_namespace
          }
          
          env {
            name  = "POLL_INTERVAL"
            value = var.poll_interval
          }
          
          env {
            name  = "RESOURCES"
            value = var.resources_to_monitor
          }
          
          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
          
          volume_mount {
            name       = "app-code"
            mount_path = "/app"
          }
        }
        
        volume {
          name = "app-code"
          config_map {
            name = kubernetes_config_map.k8s_monitor_code.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service_account" "k8s_monitor" {
  metadata {
    name      = "k8s-monitor"
    namespace = var.namespace
  }
}

resource "kubernetes_cluster_role" "k8s_monitor" {
  metadata {
    name = "k8s-monitor-role"
  }

  rule {
    api_groups = [""]
    resources  = ["pods", "services", "configmaps", "secrets", "namespaces", "nodes"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "statefulsets", "daemonsets"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = ["batch"]
    resources  = ["jobs", "cronjobs"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = ["networking.k8s.io"]
    resources  = ["ingresses"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_cluster_role_binding" "k8s_monitor" {
  metadata {
    name = "k8s-monitor-binding"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.k8s_monitor.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.k8s_monitor.metadata[0].name
    namespace = var.namespace
  }
}

resource "kubernetes_config_map" "k8s_monitor_code" {
  metadata {
    name      = "k8s-monitor-code"
    namespace = var.namespace
  }

  data = {
    "k8s_monitor.py" = file("${path.module}/../../backend/apps/python_agent/integrations/k8s_monitor.py")
  }
}
