
provider "kubernetes" {
  config_path = var.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = var.kubeconfig_path
  }
}

resource "kubernetes_namespace" "jupyterhub" {
  metadata {
    name = var.jupyterhub_namespace
    labels = {
      "app.kubernetes.io/name" = "jupyterhub"
      "app.kubernetes.io/instance" = "jupyterhub"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_persistent_volume_claim" "jupyterhub_data" {
  metadata {
    name      = "jupyterhub-data"
    namespace = kubernetes_namespace.jupyterhub.metadata[0].name
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.jupyterhub_storage_size
      }
    }
    storage_class_name = var.storage_class_name
  }
}

resource "kubernetes_secret" "minio_credentials" {
  metadata {
    name      = "minio-credentials"
    namespace = kubernetes_namespace.jupyterhub.metadata[0].name
  }

  data = {
    "accesskey" = var.minio_access_key
    "secretkey" = var.minio_secret_key
  }

  type = "Opaque"
}

resource "kubernetes_config_map" "jupyterhub_config" {
  metadata {
    name      = "jupyterhub-config"
    namespace = kubernetes_namespace.jupyterhub.metadata[0].name
  }

  data = {
    "jupyterhub_config.py" = <<-EOT
      c.JupyterHub.authenticator_class = 'jupyterhub.auth.DummyAuthenticator'
      c.DummyAuthenticator.password = 'password'
      c.Authenticator.admin_users = {'admin'}
      
      c.JupyterHub.spawner_class = 'kubespawner.KubeSpawner'
      c.KubeSpawner.image = 'jupyter/datascience-notebook:${var.jupyterhub_version}'
      c.KubeSpawner.cpu_limit = 2
      c.KubeSpawner.mem_limit = '4G'
      
      c.KubeSpawner.service_account = 'default'
      c.KubeSpawner.start_timeout = 600
      c.KubeSpawner.http_timeout = 600
      
      c.KubeSpawner.environment = {
          'AWS_ACCESS_KEY_ID': 'minioadmin',
          'AWS_SECRET_ACCESS_KEY': 'minioadmin',
          'AWS_ENDPOINT_URL': 'http://minio.minio.svc.cluster.local:9000',
          'MLFLOW_TRACKING_URI': 'http://mlflow-server.mlflow.svc.cluster.local:5000',
      }
      
      c.KubeSpawner.extra_container_config = {
          'env': [
              {
                  'name': 'AWS_ACCESS_KEY_ID',
                  'valueFrom': {
                      'secretKeyRef': {
                          'name': 'minio-credentials',
                          'key': 'accesskey',
                      }
                  }
              },
              {
                  'name': 'AWS_SECRET_ACCESS_KEY',
                  'valueFrom': {
                      'secretKeyRef': {
                          'name': 'minio-credentials',
                          'key': 'secretkey',
                      }
                  }
              }
          ]
      }
      
      c.JupyterHub.hub_ip = '0.0.0.0'
      c.JupyterHub.hub_port = 8081
      c.JupyterHub.bind_url = 'http://0.0.0.0:8000'
      
      c.JupyterHub.allow_named_servers = True
      c.JupyterHub.named_server_limit_per_user = 5
    EOT
  }
}

resource "helm_release" "jupyterhub" {
  name       = "jupyterhub"
  repository = "https://jupyterhub.github.io/helm-chart/"
  chart      = "jupyterhub"
  namespace  = kubernetes_namespace.jupyterhub.metadata[0].name
  version    = var.jupyterhub_version
  timeout    = 1200

  values = [
    <<-EOT
    hub:
      config:
        JupyterHub:
          authenticator_class: dummy
        DummyAuthenticator:
          password: password
      extraConfig:
        myConfig: |
          c.KubeSpawner.extra_container_config = {
              'env': [
                  {
                      'name': 'AWS_ACCESS_KEY_ID',
                      'valueFrom': {
                          'secretKeyRef': {
                              'name': 'minio-credentials',
                              'key': 'accesskey',
                          }
                      }
                  },
                  {
                      'name': 'AWS_SECRET_ACCESS_KEY',
                      'valueFrom': {
                          'secretKeyRef': {
                              'name': 'minio-credentials',
                              'key': 'secretkey',
                          }
                      }
                  },
                  {
                      'name': 'AWS_ENDPOINT_URL',
                      'value': 'http://minio.minio.svc.cluster.local:9000',
                  },
                  {
                      'name': 'MLFLOW_TRACKING_URI',
                      'value': 'http://mlflow-server.mlflow.svc.cluster.local:5000',
                  }
              ]
          }
      db:
        pvc:
          storageClassName: ${var.storage_class_name}
          storage: 10Gi
    
    singleuser:
      image:
        name: jupyter/datascience-notebook
        tag: ${var.jupyterhub_version}
      cpu:
        limit: 2
        guarantee: 0.5
      memory:
        limit: 4G
        guarantee: 1G
      storage:
        capacity: ${var.jupyterhub_storage_size}
        dynamic:
          storageClass: ${var.storage_class_name}
      extraEnv:
        MLFLOW_TRACKING_URI: http://mlflow-server.mlflow.svc.cluster.local:5000
        FEAST_FEATURE_STORE_CONFIG_PATH: /etc/feast/feature_store.yaml
    
    proxy:
      service:
        type: ClusterIP
    EOT
  ]

  depends_on = [
    kubernetes_namespace.jupyterhub,
    kubernetes_persistent_volume_claim.jupyterhub_data,
    kubernetes_config_map.jupyterhub_config
  ]
}

resource "kubernetes_config_map" "jupyter_notebooks" {
  metadata {
    name      = "jupyter-notebooks"
    namespace = kubernetes_namespace.jupyterhub.metadata[0].name
  }

  data = {
    "llama4_fine_tuning_example.ipynb" = <<-EOT
      {
        "cells": [
          {
            "cell_type": "markdown",
            "metadata": {},
            "source": [
              "# Llama 4 Fine-Tuning Example Notebook\\n",
              "\\n",
              "This notebook demonstrates how to fine-tune Llama 4 models using the ML infrastructure."
            ]
          },
          {
            "cell_type": "code",
            "execution_count": null,
            "metadata": {},
            "source": [
              "import os\\n",
              "import sys\\n",
              "import mlflow\\n",
              "import pandas as pd\\n",
              "import numpy as np\\n",
              "import matplotlib.pyplot as plt\\n",
              "\\n",
              "# Set MLFlow tracking URI\\n",
              "mlflow.set_tracking_uri(os.environ.get('MLFLOW_TRACKING_URI', 'http://mlflow-server.mlflow.svc.cluster.local:5000'))\\n",
              "print(f\"MLFlow tracking URI: {mlflow.get_tracking_uri()}\")"
            ]
          },
          {
            "cell_type": "code",
            "execution_count": null,
            "metadata": {},
            "source": [
              "# Connect to MinIO for data access\\n",
              "import boto3\\n",
              "from botocore.client import Config\\n",
              "\\n",
              "s3_endpoint_url = os.environ.get('AWS_ENDPOINT_URL', 'http://minio.minio.svc.cluster.local:9000')\\n",
              "s3_access_key = os.environ.get('AWS_ACCESS_KEY_ID', 'minioadmin')\\n",
              "s3_secret_key = os.environ.get('AWS_SECRET_ACCESS_KEY', 'minioadmin')\\n",
              "\\n",
              "s3_client = boto3.client(\\n",
              "    's3',\\n",
              "    endpoint_url=s3_endpoint_url,\\n",
              "    aws_access_key_id=s3_access_key,\\n",
              "    aws_secret_access_key=s3_secret_key,\\n",
              "    config=Config(signature_version='s3v4'),\\n",
              "    region_name='us-east-1'\\n",
              ")\\n",
              "\\n",
              "# List buckets\\n",
              "response = s3_client.list_buckets()\\n",
              "print(\"Available buckets:\")\\n",
              "for bucket in response['Buckets']:\\n",
              "    print(f\"- {bucket['Name']}\")"
            ]
          },
          {
            "cell_type": "markdown",
            "metadata": {},
            "source": [
              "## Load GitHub Issue Dataset\\n",
              "\\n",
              "Load the GitHub issue dataset from MinIO for fine-tuning."
            ]
          },
          {
            "cell_type": "code",
            "execution_count": null,
            "metadata": {},
            "source": [
              "# Load dataset from MinIO\\n",
              "try:\\n",
              "    s3_client.download_file('datasets', 'github_issues.parquet', 'github_issues.parquet')\\n",
              "    issues_df = pd.read_parquet('github_issues.parquet')\\n",
              "    print(f\"Loaded {len(issues_df)} GitHub issues\")\\n",
              "    issues_df.head()\\n",
              "except Exception as e:\\n",
              "    print(f\"Error loading dataset: {e}\")"
            ]
          },
          {
            "cell_type": "markdown",
            "metadata": {},
            "source": [
              "## Prepare Data for Fine-Tuning\\n",
              "\\n",
              "Prepare the GitHub issue dataset for fine-tuning Llama 4 models."
            ]
          },
          {
            "cell_type": "code",
            "execution_count": null,
            "metadata": {},
            "source": [
              "# Prepare data for fine-tuning\\n",
              "def prepare_data_for_fine_tuning(df):\\n",
              "    # Create input-output pairs for fine-tuning\\n",
              "    data = []\\n",
              "    for _, row in df.iterrows():\\n",
              "        input_text = f\"Issue: {row['title']}\\n\\nDescription: {row['body']}\"\\n",
              "        output_text = row['solution']\\n",
              "        data.append({\\n",
              "            'input': input_text,\\n",
              "            'output': output_text,\\n",
              "            'metadata': {\\n",
              "                'repository': row['repository'],\\n",
              "                'issue_id': row['issue_id'],\\n",
              "                'language': row['language']\\n",
              "            }\\n",
              "        })\\n",
              "    return data\\n",
              "\\n",
              "try:\\n",
              "    fine_tuning_data = prepare_data_for_fine_tuning(issues_df)\\n",
              "    print(f\"Prepared {len(fine_tuning_data)} examples for fine-tuning\")\\n",
              "    \\n",
              "    # Split into train, validation, and test sets\\n",
              "    from sklearn.model_selection import train_test_split\\n",
              "    train_data, test_data = train_test_split(fine_tuning_data, test_size=0.2, random_state=42)\\n",
              "    train_data, val_data = train_test_split(train_data, test_size=0.25, random_state=42)\\n",
              "    \\n",
              "    print(f\"Train: {len(train_data)}, Validation: {len(val_data)}, Test: {len(test_data)}\")\\n",
              "except Exception as e:\\n",
              "    print(f\"Error preparing data: {e}\")"
            ]
          },
          {
            "cell_type": "markdown",
            "metadata": {},
            "source": [
              "## Start MLFlow Experiment\\n",
              "\\n",
              "Create or get an MLFlow experiment for tracking the fine-tuning process."
            ]
          },
          {
            "cell_type": "code",
            "execution_count": null,
            "metadata": {},
            "source": [
              "# Set up MLFlow experiment\\n",
              "experiment_name = \"llama4-maverick-fine-tuning\"\\n",
              "mlflow.set_experiment(experiment_name)\\n",
              "\\n",
              "# Start a new run\\n",
              "with mlflow.start_run(run_name=\"github-issues-fine-tuning\") as run:\\n",
              "    # Log parameters\\n",
              "    mlflow.log_param(\"model_type\", \"llama4-maverick\")\\n",
              "    mlflow.log_param(\"dataset\", \"github_issues\")\\n",
              "    mlflow.log_param(\"train_examples\", len(train_data))\\n",
              "    mlflow.log_param(\"val_examples\", len(val_data))\\n",
              "    mlflow.log_param(\"test_examples\", len(test_data))\\n",
              "    \\n",
              "    # Log dataset statistics\\n",
              "    mlflow.log_param(\"avg_input_length\", np.mean([len(example['input']) for example in train_data]))\\n",
              "    mlflow.log_param(\"avg_output_length\", np.mean([len(example['output']) for example in train_data]))\\n",
              "    \\n",
              "    print(f\"MLFlow run ID: {run.info.run_id}\")"
            ]
          },
          {
            "cell_type": "markdown",
            "metadata": {},
            "source": [
              "## Submit Fine-Tuning Job to KubeFlow\\n",
              "\\n",
              "Submit a fine-tuning job to KubeFlow for distributed training."
            ]
          },
          {
            "cell_type": "code",
            "execution_count": null,
            "metadata": {},
            "source": [
              "# Submit fine-tuning job to KubeFlow\\n",
              "import yaml\\n",
              "import json\\n",
              "from kubernetes import client, config\\n",
              "\\n",
              "# Load Kubernetes configuration\\n",
              "try:\\n",
              "    config.load_incluster_config()\\n",
              "except:\\n",
              "    config.load_kube_config()\\n",
              "\\n",
              "# Create Kubernetes API client\\n",
              "api_client = client.ApiClient()\\n",
              "custom_api = client.CustomObjectsApi(api_client)\\n",
              "\\n",
              "# Define PyTorchJob for fine-tuning\\n",
              "pytorch_job = {\\n",
              "    \"apiVersion\": \"kubeflow.org/v1\",\\n",
              "    \"kind\": \"PyTorchJob\",\\n",
              "    \"metadata\": {\\n",
              "        \"name\": f\"llama4-fine-tuning-{run.info.run_id[:8]}\",\\n",
              "        \"namespace\": \"kubeflow\"\\n",
              "    },\\n",
              "    \"spec\": {\\n",
              "        \"pytorchReplicaSpecs\": {\\n",
              "            \"Master\": {\\n",
              "                \"replicas\": 1,\\n",
              "                \"restartPolicy\": \"OnFailure\",\\n",
              "                \"template\": {\\n",
              "                    \"spec\": {\\n",
              "                        \"containers\": [\\n",
              "                            {\\n",
              "                                \"name\": \"pytorch\",\\n",
              "                                \"image\": \"pytorch/pytorch:2.0.0-cuda11.7-cudnn8-runtime\",\\n",
              "                                \"command\": [\\n",
              "                                    \"python\",\\n",
              "                                    \"-m\",\\n",
              "                                    \"src.ml_infrastructure.training.scripts.train_llama4\",\\n",
              "                                    \"--model_type=llama4-maverick\",\\n",
              "                                    f\"--mlflow_run_id={run.info.run_id}\",\\n",
              "                                    \"--train_file=s3://datasets/train.json\",\\n",
              "                                    \"--validation_file=s3://datasets/validation.json\",\\n",
              "                                    \"--test_file=s3://datasets/test.json\"\\n",
              "                                ],\\n",
              "                                \"env\": [\\n",
              "                                    {\\n",
              "                                        \"name\": \"AWS_ACCESS_KEY_ID\",\\n",
              "                                        \"valueFrom\": {\\n",
              "                                            \"secretKeyRef\": {\\n",
              "                                                \"name\": \"minio-credentials\",\\n",
              "                                                \"key\": \"accesskey\"\\n",
              "                                            }\\n",
              "                                        }\\n",
              "                                    },\\n",
              "                                    {\\n",
              "                                        \"name\": \"AWS_SECRET_ACCESS_KEY\",\\n",
              "                                        \"valueFrom\": {\\n",
              "                                            \"secretKeyRef\": {\\n",
              "                                                \"name\": \"minio-credentials\",\\n",
              "                                                \"key\": \"secretkey\"\\n",
              "                                            }\\n",
              "                                        }\\n",
              "                                    },\\n",
              "                                    {\\n",
              "                                        \"name\": \"AWS_ENDPOINT_URL\",\\n",
              "                                        \"value\": \"http://minio.minio.svc.cluster.local:9000\"\\n",
              "                                    },\\n",
              "                                    {\\n",
              "                                        \"name\": \"MLFLOW_TRACKING_URI\",\\n",
              "                                        \"value\": \"http://mlflow-server.mlflow.svc.cluster.local:5000\"\\n",
              "                                    }\\n",
              "                                ],\\n",
              "                                \"resources\": {\\n",
              "                                    \"limits\": {\\n",
              "                                        \"nvidia.com/gpu\": 1,\\n",
              "                                        \"memory\": \"16Gi\",\\n",
              "                                        \"cpu\": 4\\n",
              "                                    },\\n",
              "                                    \"requests\": {\\n",
              "                                        \"memory\": \"8Gi\",\\n",
              "                                        \"cpu\": 2\\n",
              "                                    }\\n",
              "                                }\\n",
              "                            }\\n",
              "                        ]\\n",
              "                    }\\n",
              "                }\\n",
              "            },\\n",
              "            \"Worker\": {\\n",
              "                \"replicas\": 2,\\n",
              "                \"restartPolicy\": \"OnFailure\",\\n",
              "                \"template\": {\\n",
              "                    \"spec\": {\\n",
              "                        \"containers\": [\\n",
              "                            {\\n",
              "                                \"name\": \"pytorch\",\\n",
              "                                \"image\": \"pytorch/pytorch:2.0.0-cuda11.7-cudnn8-runtime\",\\n",
              "                                \"command\": [\\n",
              "                                    \"python\",\\n",
              "                                    \"-m\",\\n",
              "                                    \"src.ml_infrastructure.training.scripts.train_llama4\",\\n",
              "                                    \"--model_type=llama4-maverick\",\\n",
              "                                    f\"--mlflow_run_id={run.info.run_id}\",\\n",
              "                                    \"--train_file=s3://datasets/train.json\",\\n",
              "                                    \"--validation_file=s3://datasets/validation.json\",\\n",
              "                                    \"--test_file=s3://datasets/test.json\",\\n",
              "                                    \"--distributed=True\"\\n",
              "                                ],\\n",
              "                                \"env\": [\\n",
              "                                    {\\n",
              "                                        \"name\": \"AWS_ACCESS_KEY_ID\",\\n",
              "                                        \"valueFrom\": {\\n",
              "                                            \"secretKeyRef\": {\\n",
              "                                                \"name\": \"minio-credentials\",\\n",
              "                                                \"key\": \"accesskey\"\\n",
              "                                            }\\n",
              "                                        }\\n",
              "                                    },\\n",
              "                                    {\\n",
              "                                        \"name\": \"AWS_SECRET_ACCESS_KEY\",\\n",
              "                                        \"valueFrom\": {\\n",
              "                                            \"secretKeyRef\": {\\n",
              "                                                \"name\": \"minio-credentials\",\\n",
              "                                                \"key\": \"secretkey\"\\n",
              "                                            }\\n",
              "                                        }\\n",
              "                                    },\\n",
              "                                    {\\n",
              "                                        \"name\": \"AWS_ENDPOINT_URL\",\\n",
              "                                        \"value\": \"http://minio.minio.svc.cluster.local:9000\"\\n",
              "                                    },\\n",
              "                                    {\\n",
              "                                        \"name\": \"MLFLOW_TRACKING_URI\",\\n",
              "                                        \"value\": \"http://mlflow-server.mlflow.svc.cluster.local:5000\"\\n",
              "                                    }\\n",
              "                                ],\\n",
              "                                \"resources\": {\\n",
              "                                    \"limits\": {\\n",
              "                                        \"nvidia.com/gpu\": 1,\\n",
              "                                        \"memory\": \"16Gi\",\\n",
              "                                        \"cpu\": 4\\n",
              "                                    },\\n",
              "                                    \"requests\": {\\n",
              "                                        \"memory\": \"8Gi\",\\n",
              "                                        \"cpu\": 2\\n",
              "                                    }\\n",
              "                                }\\n",
              "                            }\\n",
              "                        ]\\n",
              "                    }\\n",
              "                }\\n",
              "            }\\n",
              "        }\\n",
              "    }\\n",
              "}\\n",
              "\\n",
              "# Save PyTorchJob YAML for reference\\n",
              "with open('pytorch_job.yaml', 'w') as f:\\n",
              "    yaml.dump(pytorch_job, f)\\n",
              "\\n",
              "# Submit PyTorchJob to KubeFlow\\n",
              "try:\\n",
              "    response = custom_api.create_namespaced_custom_object(\\n",
              "        group=\"kubeflow.org\",\\n",
              "        version=\"v1\",\\n",
              "        namespace=\"kubeflow\",\\n",
              "        plural=\"pytorchjobs\",\\n",
              "        body=pytorch_job\\n",
              "    )\\n",
              "    print(f\"PyTorchJob created: {response['metadata']['name']}\")"
            ]
          }
        ],
        "metadata": {
          "kernelspec": {
            "display_name": "Python 3",
            "language": "python",
            "name": "python3"
          },
          "language_info": {
            "codemirror_mode": {
              "name": "ipython",
              "version": 3
            },
            "file_extension": ".py",
            "mimetype": "text/x-python",
            "name": "python",
            "nbconvert_exporter": "python",
            "pygments_lexer": "ipython3",
            "version": "3.8.10"
          }
        },
        "nbformat": 4,
        "nbformat_minor": 4
      }
    EOT
  }
}
