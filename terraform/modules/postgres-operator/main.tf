provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "postgres_operator" {
  metadata {
    name = var.postgres_operator_namespace
    labels = {
      "app.kubernetes.io/name" = "postgres-operator"
      "app.kubernetes.io/instance" = "postgres-operator"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_service_account" "postgres_operator" {
  metadata {
    name      = "postgres-operator"
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "postgres_operator" {
  metadata {
    name = "postgres-operator"
  }

  rule {
    api_groups = [""]
    resources  = ["configmaps", "endpoints", "events", "namespaces", "persistentvolumeclaims", "pods", "secrets", "services"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "replicasets", "statefulsets"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["batch"]
    resources  = ["cronjobs", "jobs"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["policy"]
    resources  = ["poddisruptionbudgets"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["postgres-operator.crunchydata.com"]
    resources  = ["postgresclusters"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["postgres-operator.crunchydata.com"]
    resources  = ["postgresclusters/status"]
    verbs      = ["get", "patch", "update"]
  }
}

resource "kubernetes_cluster_role_binding" "postgres_operator" {
  metadata {
    name = "postgres-operator"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.postgres_operator.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.postgres_operator.metadata[0].name
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
  }
}

resource "kubernetes_config_map" "postgres_operator_config" {
  metadata {
    name      = "postgres-operator-config"
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
  }

  data = {
    "postgresql.conf" = file("${path.module}/../../k8s/postgres-operator/config/postgres-config.yaml")
  }
}

resource "kubernetes_deployment" "postgres_operator" {
  metadata {
    name      = "postgres-operator"
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
    labels = {
      app = "postgres-operator"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "postgres-operator"
      }
    }

    template {
      metadata {
        labels = {
          app = "postgres-operator"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.postgres_operator.metadata[0].name

        container {
          name  = "postgres-operator"
          image = "registry.developers.crunchydata.com/crunchydata/postgres-operator:v5.3.1"
          
          env {
            name  = "CRUNCHY_DEBUG"
            value = "false"
          }
          
          env {
            name = "PGO_NAMESPACE"
            value_from {
              field_ref {
                field_path = "metadata.namespace"
              }
            }
          }
          
          env {
            name  = "PGO_TARGET_NAMESPACE"
            value = "*"
          }
          
          env {
            name  = "RELATED_IMAGE_PGBACKREST"
            value = "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.41-0"
          }
          
          env {
            name  = "RELATED_IMAGE_PGBOUNCER"
            value = "registry.developers.crunchydata.com/crunchydata/crunchy-pgbouncer:ubi8-1.17-0"
          }
          
          env {
            name  = "RELATED_IMAGE_PGEXPORTER"
            value = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres-exporter:ubi8-5.3.1-0"
          }
          
          env {
            name  = "RELATED_IMAGE_POSTGRES_13"
            value = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-13.10-0"
          }
          
          env {
            name  = "RELATED_IMAGE_POSTGRES_14"
            value = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-14.6-0"
          }
          
          env {
            name  = "RELATED_IMAGE_POSTGRES_15"
            value = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-15.1-0"
          }
          
          env {
            name  = "VAULT_ADDR"
            value = "http://vault.vault.svc.cluster.local:8200"
          }
          
          env {
            name = "VAULT_TOKEN"
            value_from {
              secret_key_ref {
                name = "vault-token"
                key  = "token"
              }
            }
          }
          
          resources {
            limits = {
              cpu    = "1"
              memory = "1Gi"
            }
            requests = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
          
          security_context {
            allow_privilege_escalation = false
            capabilities {
              drop = ["ALL"]
            }
            privileged = false
            read_only_root_filesystem = true
            run_as_non_root = true
          }
        }
      }
    }
  }
}

resource "kubernetes_custom_resource_definition" "postgres_cluster" {
  metadata {
    name = "postgresclusters.postgres-operator.crunchydata.com"
  }

  spec {
    group = "postgres-operator.crunchydata.com"
    names {
      kind     = "PostgresCluster"
      plural   = "postgresclusters"
      singular = "postgrescluster"
      short_names = ["pgc", "pgclusters"]
    }
    scope = "Namespaced"
    versions {
      name    = "v1beta1"
      served  = true
      storage = true
      schema {
        open_apiv3_schema = file("${path.module}/../../k8s/postgres-operator/manifests/crd.yaml")
      }
      subresources {
        status = {}
      }
      additional_printer_columns {
        name     = "Age"
        type     = "date"
        json_path = ".metadata.creationTimestamp"
      }
      additional_printer_columns {
        name     = "PG Version"
        type     = "string"
        json_path = ".spec.postgresVersion"
      }
      additional_printer_columns {
        name     = "Status"
        type     = "string"
        json_path = ".status.conditions[?(@.type==\"Ready\")].status"
      }
    }
  }
}

resource "helm_release" "vnode_runtime" {
  name       = "vnode-runtime"
  repository = "https://charts.loft.sh"
  chart      = "vnode-runtime"
  version    = "0.0.2"
  namespace  = "vnode-runtime"
  create_namespace = true

  set {
    name  = "global.imageRegistry"
    value = ""
  }

  set {
    name  = "vnodeRuntime.enabled"
    value = "true"
  }

  set {
    name  = "vnodeRuntime.image.repository"
    value = "loftsh/vnode-runtime"
  }

  set {
    name  = "vnodeRuntime.image.tag"
    value = "0.0.2"
  }

  set {
    name  = "vnodeRuntime.image.pullPolicy"
    value = "IfNotPresent"
  }

  set {
    name  = "vnodeRuntime.resources.limits.cpu"
    value = "1"
  }

  set {
    name  = "vnodeRuntime.resources.limits.memory"
    value = "1Gi"
  }

  set {
    name  = "vnodeRuntime.resources.requests.cpu"
    value = "500m"
  }

  set {
    name  = "vnodeRuntime.resources.requests.memory"
    value = "512Mi"
  }

  set {
    name  = "vnodeRuntime.integrations.postgres.enabled"
    value = "true"
  }

  set {
    name  = "vnodeRuntime.integrations.postgres.operatorNamespace"
    value = kubernetes_namespace.postgres_operator.metadata[0].name
  }
}

resource "kubernetes_manifest" "ml_postgres_cluster" {
  manifest = {
    apiVersion = "postgres-operator.crunchydata.com/v1beta1"
    kind       = "PostgresCluster"
    metadata = {
      name      = "ml-postgres-cluster"
      namespace = kubernetes_namespace.postgres_operator.metadata[0].name
    }
    spec = {
      image = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-15.1-0"
      postgresVersion = 15
      instances = [
        {
          name     = "instance1"
          replicas = 3
          dataVolumeClaimSpec = {
            accessModes = ["ReadWriteOnce"]
            resources = {
              requests = {
                storage = "10Gi"
              }
            }
          }
          affinity = {
            podAntiAffinity = {
              preferredDuringSchedulingIgnoredDuringExecution = [
                {
                  weight = 100
                  podAffinityTerm = {
                    topologyKey = "kubernetes.io/hostname"
                    labelSelector = {
                      matchLabels = {
                        "postgres-operator.crunchydata.com/cluster" = "ml-postgres-cluster"
                        "postgres-operator.crunchydata.com/instance-set" = "instance1"
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      ]
      backups = {
        pgbackrest = {
          image = "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.41-0"
          repos = [
            {
              name = "repo1"
              volume = {
                volumeClaimSpec = {
                  accessModes = ["ReadWriteOnce"]
                  resources = {
                    requests = {
                      storage = "20Gi"
                    }
                  }
                }
              }
            }
          ]
        }
      }
      patroni = {
        dynamicConfiguration = {
          postgresql = {
            parameters = {
              max_connections = "100"
              shared_buffers = "256MB"
              work_mem = "16MB"
              maintenance_work_mem = "64MB"
              effective_cache_size = "1GB"
              checkpoint_timeout = "15min"
              checkpoint_completion_target = "0.9"
              max_wal_size = "1GB"
              min_wal_size = "128MB"
              random_page_cost = "1.1"
              effective_io_concurrency = "200"
              log_min_duration_statement = "1000"
              log_checkpoints = "on"
              log_connections = "on"
              log_disconnections = "on"
              log_lock_waits = "on"
              log_temp_files = "0"
            }
          }
        }
      }
      users = [
        {
          name = "mlflow"
          databases = ["mlflow"]
          options = "CREATEDB"
        },
        {
          name = "jupyterlab"
          databases = ["jupyterlab"]
          options = "CREATEDB"
        },
        {
          name = "kubeflow"
          databases = ["kubeflow"]
          options = "CREATEDB"
        }
      ]
      monitoring = {
        pgmonitor = {
          exporter = {
            image = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres-exporter:ubi8-5.3.1-0"
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_custom_resource_definition.postgres_cluster
  ]
}
