
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "kubeflow" {
  metadata {
    name = "kubeflow"
    labels = {
      "app.kubernetes.io/name" = "kubeflow"
      "app.kubernetes.io/instance" = "kubeflow"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "helm_release" "kubeflow" {
  name       = "kubeflow"
  repository = "https://kubeflow.github.io/manifests"
  chart      = "kubeflow"
  namespace  = kubernetes_namespace.kubeflow.metadata[0].name
  version    = var.kubeflow_version
  timeout    = 1200

  set {
    name  = "pipeline.enabled"
    value = "true"
  }

  set {
    name  = "notebook.enabled"
    value = "true"
  }

  set {
    name  = "katib.enabled"
    value = "true"
  }

  set {
    name  = "training.enabled"
    value = "true"
  }

  depends_on = [
    kubernetes_namespace.kubeflow
  ]
}

resource "kubernetes_manifest" "pipeline_application" {
  manifest = {
    apiVersion = "app.k8s.io/v1beta1"
    kind       = "Application"
    metadata = {
      name      = "pipeline"
      namespace = kubernetes_namespace.kubeflow.metadata[0].name
      labels = {
        "app.kubernetes.io/name"      = "pipeline"
        "app.kubernetes.io/instance"  = "pipeline-v1.0.0"
        "app.kubernetes.io/version"   = "v1.0.0"
        "app.kubernetes.io/component" = "pipeline"
        "app.kubernetes.io/part-of"   = "kubeflow"
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      descriptor = {
        type        = "kubeflow-pipeline"
        version     = "v1.0.0"
        description = "Kubeflow Pipelines"
      }
      addOwnerRef = true
    }
  }

  depends_on = [
    helm_release.kubeflow
  ]
}

resource "kubernetes_manifest" "training_operator" {
  manifest = {
    apiVersion = "apps/v1"
    kind       = "Deployment"
    metadata = {
      name      = "training-operator"
      namespace = kubernetes_namespace.kubeflow.metadata[0].name
      labels = {
        "app.kubernetes.io/name"      = "training-operator"
        "app.kubernetes.io/instance"  = "training-operator-v1.0.0"
        "app.kubernetes.io/version"   = "v1.0.0"
        "app.kubernetes.io/component" = "training-operator"
        "app.kubernetes.io/part-of"   = "kubeflow"
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      replicas = 1
      selector = {
        matchLabels = {
          "app.kubernetes.io/name" = "training-operator"
        }
      }
      template = {
        metadata = {
          labels = {
            "app.kubernetes.io/name" = "training-operator"
          }
        }
        spec = {
          containers = [
            {
              name  = "training-operator"
              image = "kubeflow/training-operator:${var.training_operator_version}"
              command = [
                "/manager",
                "-enable-scheme=pytorch",
                "-enable-scheme=tensorflow",
                "-enable-scheme=mxnet",
                "-enable-scheme=xgboost"
              ]
              env = [
                {
                  name = "KUBEFLOW_NAMESPACE"
                  valueFrom = {
                    fieldRef = {
                      fieldPath = "metadata.namespace"
                    }
                  }
                }
              ]
              resources = {
                limits = {
                  cpu    = "500m"
                  memory = "512Mi"
                }
                requests = {
                  cpu    = "100m"
                  memory = "256Mi"
                }
              }
            }
          ]
          serviceAccountName = "training-operator"
        }
      }
    }
  }

  depends_on = [
    helm_release.kubeflow
  ]
}

resource "kubernetes_manifest" "katib_controller" {
  manifest = {
    apiVersion = "apps/v1"
    kind       = "Deployment"
    metadata = {
      name      = "katib-controller"
      namespace = kubernetes_namespace.kubeflow.metadata[0].name
      labels = {
        "app.kubernetes.io/name"      = "katib-controller"
        "app.kubernetes.io/instance"  = "katib-controller-v1.0.0"
        "app.kubernetes.io/version"   = "v1.0.0"
        "app.kubernetes.io/component" = "katib"
        "app.kubernetes.io/part-of"   = "kubeflow"
        "app.kubernetes.io/managed-by" = "terraform"
      }
    }
    spec = {
      replicas = 1
      selector = {
        matchLabels = {
          "app.kubernetes.io/name" = "katib-controller"
        }
      }
      template = {
        metadata = {
          labels = {
            "app.kubernetes.io/name" = "katib-controller"
          }
        }
        spec = {
          containers = [
            {
              name  = "katib-controller"
              image = "kubeflow/katib-controller:${var.katib_version}"
              command = [
                "./katib-controller"
              ]
              env = [
                {
                  name = "KATIB_CORE_NAMESPACE"
                  valueFrom = {
                    fieldRef = {
                      fieldPath = "metadata.namespace"
                    }
                  }
                }
              ]
              resources = {
                limits = {
                  cpu    = "500m"
                  memory = "512Mi"
                }
                requests = {
                  cpu    = "100m"
                  memory = "256Mi"
                }
              }
            }
          ]
          serviceAccountName = "katib-controller"
        }
      }
    }
  }

  depends_on = [
    helm_release.kubeflow
  ]
}

resource "kubernetes_persistent_volume_claim" "kubeflow_data" {
  metadata {
    name      = "kubeflow-data"
    namespace = kubernetes_namespace.kubeflow.metadata[0].name
    labels = {
      "app.kubernetes.io/name"      = "kubeflow-data"
      "app.kubernetes.io/instance"  = "kubeflow-data-v1.0.0"
      "app.kubernetes.io/version"   = "v1.0.0"
      "app.kubernetes.io/component" = "storage"
      "app.kubernetes.io/part-of"   = "kubeflow"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.kubeflow_data_storage_size
      }
    }
    storage_class_name = var.storage_class_name
  }

  depends_on = [
    kubernetes_namespace.kubeflow
  ]
}

resource "kubernetes_config_map" "llama4_training_config" {
  metadata {
    name      = "llama4-training-config"
    namespace = kubernetes_namespace.kubeflow.metadata[0].name
    labels = {
      "app.kubernetes.io/name"      = "llama4-training-config"
      "app.kubernetes.io/instance"  = "llama4-training-config-v1.0.0"
      "app.kubernetes.io/version"   = "v1.0.0"
      "app.kubernetes.io/component" = "config"
      "app.kubernetes.io/part-of"   = "kubeflow"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
  data = {
    "llama4-maverick-config.json" = jsonencode({
      model_id                    = "meta-llama/llama-4-maverick"
      output_dir                  = "/models/llama4-maverick"
      train_file                  = "/data/train.json"
      validation_file             = "/data/validation.json"
      test_file                   = "/data/test.json"
      max_seq_length              = 4096
      learning_rate               = 5e-5
      num_train_epochs            = 3
      per_device_train_batch_size = 8
      per_device_eval_batch_size  = 8
      gradient_accumulation_steps = 4
      warmup_steps                = 500
      weight_decay                = 0.01
      logging_steps               = 100
      evaluation_strategy         = "steps"
      eval_steps                  = 500
      save_steps                  = 1000
      save_total_limit            = 3
      fp16                        = true
      bf16                        = false
      load_best_model_at_end      = true
      metric_for_best_model       = "eval_loss"
      greater_is_better           = false
      seed                        = 42
      lora_r                      = 16
      lora_alpha                  = 32
      lora_dropout                = 0.05
      use_lora                    = true
      use_8bit_quantization       = false
      use_4bit_quantization       = false
    })
    "llama4-scout-config.json" = jsonencode({
      model_id                    = "meta-llama/llama-4-scout"
      output_dir                  = "/models/llama4-scout"
      train_file                  = "/data/train.json"
      validation_file             = "/data/validation.json"
      test_file                   = "/data/test.json"
      max_seq_length              = 4096
      learning_rate               = 5e-5
      num_train_epochs            = 3
      per_device_train_batch_size = 8
      per_device_eval_batch_size  = 8
      gradient_accumulation_steps = 4
      warmup_steps                = 500
      weight_decay                = 0.01
      logging_steps               = 100
      evaluation_strategy         = "steps"
      eval_steps                  = 500
      save_steps                  = 1000
      save_total_limit            = 3
      fp16                        = true
      bf16                        = false
      load_best_model_at_end      = true
      metric_for_best_model       = "eval_loss"
      greater_is_better           = false
      seed                        = 42
      lora_r                      = 16
      lora_alpha                  = 32
      lora_dropout                = 0.05
      use_lora                    = true
      use_8bit_quantization       = false
      use_4bit_quantization       = false
    })
    "data-config.json" = jsonencode({
      train_file                = "/data/github_issues_train.json"
      validation_file           = "/data/github_issues_validation.json"
      test_file                 = "/data/github_issues_test.json"
      max_seq_length            = 4096
      input_column              = "issue"
      output_column             = "solution"
      metadata_column           = "metadata"
      preprocessing_num_workers = null
      overwrite_cache           = false
      preprocessing_batch_size  = 1000
      streaming                 = false
      use_auth_token            = false
      ignore_pad_token_for_loss = true
      pad_to_max_length         = false
      max_train_samples         = null
      max_eval_samples          = null
      max_predict_samples       = null
    })
  }

  depends_on = [
    kubernetes_namespace.kubeflow
  ]
}
