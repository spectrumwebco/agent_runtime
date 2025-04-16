resource "kubernetes_stateful_set" "doris_fe" {
  metadata {
    name      = "doris-fe"
    namespace = var.namespace
    labels = {
      app       = "doris"
      component = "fe"
      service   = "database"
    }
  }

  spec {
    service_name = "doris-fe"
    replicas     = var.fe_replicas

    selector {
      match_labels = {
        app       = "doris"
        component = "fe"
      }
    }

    template {
      metadata {
        labels = {
          app       = "doris"
          component = "fe"
          service   = "database"
        }
      }

      spec {
        container {
          name  = "doris-fe"
          image = "apache/doris:${var.doris_version}-fe"

          port {
            container_port = 8030
            name           = "http"
          }

          port {
            container_port = 9020
            name           = "rpc"
          }

          port {
            container_port = 9030
            name           = "query"
          }

          env {
            name  = "FE_SERVERS"
            value = "doris-fe-0.doris-fe.${var.namespace}.svc.cluster.local:9010"
          }

          env {
            name  = "DORIS_ADMIN_USER"
            value = "root"
          }

          env {
            name = "DORIS_ADMIN_PASSWORD"
            value_from {
              secret_key_ref {
                name = "doris-credentials"
                key  = "admin-password"
              }
            }
          }

          volume_mount {
            name       = "doris-fe-data"
            mount_path = "/opt/apache-doris/fe/doris-meta"
          }

          resources {
            requests = {
              memory = var.fe_memory_request
              cpu    = var.fe_cpu_request
            }
            limits = {
              memory = var.fe_memory_limit
              cpu    = var.fe_cpu_limit
            }
          }

          readiness_probe {
            http_get {
              path = "/api/bootstrap"
              port = 8030
            }
            initial_delay_seconds = 30
            period_seconds        = 15
          }

          liveness_probe {
            http_get {
              path = "/api/bootstrap"
              port = 8030
            }
            initial_delay_seconds = 60
            period_seconds        = 30
          }
        }
      }
    }

    volume_claim_template {
      metadata {
        name = "doris-fe-data"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            storage = var.fe_storage_size
          }
        }
      }
    }
  }
}

resource "kubernetes_stateful_set" "doris_be" {
  metadata {
    name      = "doris-be"
    namespace = var.namespace
    labels = {
      app       = "doris"
      component = "be"
      service   = "database"
    }
  }

  spec {
    service_name = "doris-be"
    replicas     = var.be_replicas

    selector {
      match_labels = {
        app       = "doris"
        component = "be"
      }
    }

    template {
      metadata {
        labels = {
          app       = "doris"
          component = "be"
          service   = "database"
        }
      }

      spec {
        container {
          name  = "doris-be"
          image = "apache/doris:${var.doris_version}-be"

          port {
            container_port = 9050
            name           = "be-port"
          }

          port {
            container_port = 8040
            name           = "webserver"
          }

          env {
            name  = "FE_SERVERS"
            value = "doris-fe-0.doris-fe.${var.namespace}.svc.cluster.local:9020"
          }

          volume_mount {
            name       = "doris-be-data"
            mount_path = "/opt/apache-doris/be/storage"
          }

          resources {
            requests = {
              memory = var.be_memory_request
              cpu    = var.be_cpu_request
            }
            limits = {
              memory = var.be_memory_limit
              cpu    = var.be_cpu_limit
            }
          }

          readiness_probe {
            http_get {
              path = "/api/health"
              port = 8040
            }
            initial_delay_seconds = 30
            period_seconds        = 15
          }

          liveness_probe {
            http_get {
              path = "/api/health"
              port = 8040
            }
            initial_delay_seconds = 60
            period_seconds        = 30
          }
        }
      }
    }

    volume_claim_template {
      metadata {
        name = "doris-be-data"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            storage = var.be_storage_size
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "doris_fe" {
  metadata {
    name      = "doris-fe"
    namespace = var.namespace
    labels = {
      app       = "doris"
      component = "fe"
      service   = "database"
    }
  }

  spec {
    selector = {
      app       = "doris"
      component = "fe"
    }

    port {
      port        = 8030
      target_port = 8030
      protocol    = "TCP"
      name        = "http"
    }

    port {
      port        = 9020
      target_port = 9020
      protocol    = "TCP"
      name        = "rpc"
    }

    port {
      port        = 9030
      target_port = 9030
      protocol    = "TCP"
      name        = "query"
    }
  }
}

resource "kubernetes_service" "doris_be" {
  metadata {
    name      = "doris-be"
    namespace = var.namespace
    labels = {
      app       = "doris"
      component = "be"
      service   = "database"
    }
  }

  spec {
    selector = {
      app       = "doris"
      component = "be"
    }

    port {
      port        = 9050
      target_port = 9050
      protocol    = "TCP"
      name        = "be-port"
    }

    port {
      port        = 8040
      target_port = 8040
      protocol    = "TCP"
      name        = "webserver"
    }
  }
}

resource "kubernetes_secret" "doris_credentials" {
  metadata {
    name      = "doris-credentials"
    namespace = var.namespace
  }

  data = {
    "admin-password" = var.admin_password
  }

  type = "Opaque"
}
