
resource "kubernetes_namespace" "postgres_operator" {
  metadata {
    name = "pgo"
  }
}

resource "kubernetes_service_account" "postgres_operator" {
  metadata {
    name      = "pgo-deployer-sa"
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "postgres_operator" {
  metadata {
    name = "pgo-deployer-cr"
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["*"]
  }
}

resource "kubernetes_cluster_role_binding" "postgres_operator" {
  metadata {
    name = "pgo-deployer-crb"
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

resource "kubernetes_config_map" "postgres_operator" {
  metadata {
    name      = "pgo-deployer-cm"
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
  }

  data = {
    "values.yaml" = <<-EOT
      archive_mode: "true"
      archive_timeout: "60"
      backrest_aws_s3_key: ""
      backrest_aws_s3_secret: ""
      backrest_aws_s3_bucket: ""
      backrest_aws_s3_endpoint: ""
      backrest_aws_s3_region: ""
      backrest_aws_s3_uri_style: ""
      backrest_aws_s3_verify_tls: "true"
      backrest_gcs_key: ""
      backrest_gcs_bucket: ""
      backrest_gcs_endpoint: ""
      backrest_gcs_keytype: "service"
      backrest_port: "2022"
      badger: "false"
      ccp_image_prefix: "registry.developers.crunchydata.com/crunchydata"
      ccp_image_pull_secret: ""
      ccp_image_pull_secret_manifest: ""
      ccp_image_tag: "centos8-13.6-0"
      create_rbac: "true"
      db_name: ""
      db_password_age_days: "0"
      db_password_length: "24"
      db_port: "5432"
      db_replicas: "0"
      db_user: "postgres"
      default_instance_memory: "128Mi"
      default_pgbackrest_memory: "48Mi"
      default_pgbouncer_memory: "24Mi"
      default_exporter_memory: "24Mi"
      delete_metrics_namespace: "false"
      delete_operator_namespace: "false"
      delete_watched_namespaces: "false"
      disable_auto_failover: "false"
      disable_fsgroup: "false"
      reconcile_rbac: "true"
      exporterport: "9187"
      metrics: "false"
      metrics_namespace: "pgo"
      metrics_image_prefix: "registry.developers.crunchydata.com/crunchydata"
      metrics_image_tag: "centos8-5.0.2-0"
      namespace: "pgo"
      namespace_mode: "dynamic"
      pgbadgerport: "10000"
      pgo_admin_password: "admin"
      pgo_admin_perms: "*"
      pgo_admin_role_name: "pgoadmin"
      pgo_admin_username: "admin"
      pgo_apiserver_port: "8443"
      pgo_apiserver_url: "https://postgres-operator"
      pgo_client_cert_secret: "pgo.tls"
      pgo_client_container_install: "false"
      pgo_client_install: "true"
      pgo_client_version: "v5.1.0"
      pgo_cluster_admin: "false"
      pgo_disable_eventing: "false"
      pgo_disable_tls: "false"
      pgo_image_prefix: "registry.developers.crunchydata.com/crunchydata"
      pgo_image_tag: "centos8-5.1.0"
      pgo_installation_name: "devtest"
      pgo_noauth_routes: ""
      pgo_operator_namespace: "pgo"
      pgo_tls_ca_store: ""
      pgo_tls_no_verify: "false"
      pod_anti_affinity: "preferred"
      pod_anti_affinity_pgbackrest: ""
      pod_anti_affinity_pgbouncer: ""
      scheduler_timeout: "3600"
      service_type: "ClusterIP"
      sync_replication: "false"
      backrest_storage: "default"
      backup_storage: "default"
      primary_storage: "default"
      replica_storage: "default"
      wal_storage: ""
      storage1_name: "default"
      storage1_access_mode: "ReadWriteOnce"
      storage1_size: "1G"
      storage1_type: "dynamic"
      storage2_name: "hostpathstorage"
      storage2_access_mode: "ReadWriteMany"
      storage2_size: "1G"
      storage2_type: "create"
      storage3_name: "nfsstorage"
      storage3_access_mode: "ReadWriteMany"
      storage3_size: "1G"
      storage3_type: "create"
      storage3_supplemental_groups: "65534"
      storage4_name: "nfsstoragered"
      storage4_access_mode: "ReadWriteMany"
      storage4_size: "1G"
      storage4_match_labels: "crunchyzone=red"
      storage4_type: "create"
      storage4_supplemental_groups: "65534"
      storage5_name: "storageos"
      storage5_access_mode: "ReadWriteOnce"
      storage5_size: "5Gi"
      storage5_type: "dynamic"
      storage5_class: "fast"
      storage6_name: "primarysite"
      storage6_access_mode: "ReadWriteOnce"
      storage6_size: "4G"
      storage6_type: "dynamic"
      storage6_class: "primarysite"
      storage7_name: "alternatesite"
      storage7_access_mode: "ReadWriteOnce"
      storage7_size: "4G"
      storage7_type: "dynamic"
      storage7_class: "alternatesite"
      storage8_name: "gce"
      storage8_access_mode: "ReadWriteOnce"
      storage8_size: "300M"
      storage8_type: "dynamic"
      storage8_class: "standard"
      storage9_name: "rook"
      storage9_access_mode: "ReadWriteOnce"
      storage9_size: "1Gi"
      storage9_type: "dynamic"
      storage9_class: "rook-ceph-block"
    EOT
  }
}

resource "kubernetes_job" "postgres_operator_deploy" {
  metadata {
    name      = "pgo-deploy"
    namespace = kubernetes_namespace.postgres_operator.metadata[0].name
  }

  spec {
    backoff_limit = 0
    template {
      metadata {
        name = "pgo-deploy"
      }
      spec {
        service_account_name = kubernetes_service_account.postgres_operator.metadata[0].name
        restart_policy       = "Never"
        container {
          name              = "pgo-deploy"
          image             = "registry.developers.crunchydata.com/crunchydata/pgo-deployer:v5.1.0"
          image_pull_policy = "IfNotPresent"
          
          env {
            name  = "DEPLOY_ACTION"
            value = "install"
          }
          
          env {
            name  = "WATCH_NAMESPACE"
            value = "pgo,default"
          }
          
          env {
            name  = "PGO_NAMESPACE"
            value = "pgo"
          }
          
          env {
            name  = "PGO_TARGET_NAMESPACE"
            value = "default"
          }
          
          env {
            name  = "DISABLE_TELEMETRY"
            value = "true"
          }
          
          env {
            name  = "TLS_CA_TRUST"
            value = ""
          }
          
          env {
            name  = "TLS_REPLICATION_CA_TRUST"
            value = ""
          }
          
          volume_mount {
            name       = "deployer-conf"
            mount_path = "/conf"
          }
        }
        
        volume {
          name = "deployer-conf"
          config_map {
            name = kubernetes_config_map.postgres_operator.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "postgres_cluster" {
  manifest = {
    apiVersion = "postgres-operator.crunchydata.com/v1beta1"
    kind       = "PostgresCluster"
    metadata = {
      name      = var.cluster_name
      namespace = var.namespace
    }
    spec = {
      image            = "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-14.5-0"
      postgresVersion  = 14
      instances = [
        {
          name     = "instance1"
          replicas = var.replicas
          dataVolumeClaimSpec = {
            accessModes = ["ReadWriteOnce"]
            resources = {
              requests = {
                storage = var.storage_size
              }
            }
          }
          affinity = {
            podAntiAffinity = {
              preferredDuringSchedulingIgnoredDuringExecution = [
                {
                  weight = 1
                  podAffinityTerm = {
                    topologyKey = "kubernetes.io/hostname"
                    labelSelector = {
                      matchLabels = {
                        "postgres-operator.crunchydata.com/cluster"     = var.cluster_name
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
          image = "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.38-0"
          repos = [
            {
              name = "repo1"
              volume = {
                volumeClaimSpec = {
                  accessModes = ["ReadWriteOnce"]
                  resources = {
                    requests = {
                      storage = "5Gi"
                    }
                  }
                }
              }
            }
          ]
        }
      }
      users = [
        {
          name      = "agent_user"
          databases = ["agent_runtime", "agent_db", "trajectory_db", "ml_db"]
          options   = "SUPERUSER CREATEDB"
        },
        {
          name      = "app_user"
          databases = ["agent_runtime"]
        }
      ]
      patroni = {
        dynamicConfiguration = {
          postgresql = {
            parameters = {
              shared_buffers  = "256MB"
              max_connections = "200"
              log_statement   = "all"
            }
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_job.postgres_operator_deploy
  ]
}
