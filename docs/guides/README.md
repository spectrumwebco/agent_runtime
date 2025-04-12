# ML Infrastructure Setup Guides

This document provides setup and usage guides for all components of the ML infrastructure.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Kubernetes Cluster Setup](#kubernetes-cluster-setup)
- [Terraform Installation](#terraform-installation)
- [KubeFlow Setup](#kubeflow-setup)
- [MLFlow Setup](#mlflow-setup)
- [KServe Setup](#kserve-setup)
- [MinIO Setup](#minio-setup)
- [Feast Setup](#feast-setup)
- [Seldon Core Setup](#seldon-core-setup)
- [JupyterLab Setup](#jupyterlab-setup)
- [h2o.ai Setup](#h2oai-setup)
- [Data Pipeline Setup](#data-pipeline-setup)

## Prerequisites

Before setting up the ML infrastructure, ensure you have the following prerequisites:

- Kubernetes cluster (v1.21+)
- Terraform (v1.0.0+)
- kubectl (v1.21+)
- Helm (v3.0.0+)
- Docker (v20.10.0+)
- Python (v3.8+)

## Kubernetes Cluster Setup

The ML infrastructure requires a Kubernetes cluster with sufficient resources for training and serving models.

### Resource Requirements

- At least 4 nodes
- Each node with at least 8 CPU cores and 32GB RAM
- GPU nodes for training (NVIDIA T4 or better)
- At least 500GB of storage

### Setup Steps

1. Create a Kubernetes cluster using your preferred provider (e.g., GKE, EKS, AKS).
2. Configure kubectl to connect to your cluster:

```bash
kubectl config use-context <your-cluster-context>
```

3. Verify the cluster is running:

```bash
kubectl get nodes
```

## Terraform Installation

Terraform is used to manage the infrastructure as code.

### Installation Steps

1. Download and install Terraform from the [official website](https://www.terraform.io/downloads.html).
2. Verify the installation:

```bash
terraform --version
```

3. Initialize Terraform in the project directory:

```bash
cd terraform
terraform init
```

## KubeFlow Setup

KubeFlow is used to orchestrate ML workflows, including training, hyperparameter tuning, and pipelines.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.kubeflow
```

2. Verify the installation:

```bash
kubectl get pods -n kubeflow
```

### Configuration

1. Configure KubeFlow pipelines:

```bash
kubectl apply -f src/ml_infrastructure/kubeflow/manifests/pipelines/pipeline.yaml
```

2. Configure KubeFlow training operator:

```bash
kubectl apply -f src/ml_infrastructure/kubeflow/manifests/training/training-operator.yaml
```

3. Configure KubeFlow Katib for hyperparameter tuning:

```bash
kubectl apply -f src/ml_infrastructure/kubeflow/manifests/katib/hyperparameter-tuning.yaml
```

## MLFlow Setup

MLFlow is used for experiment tracking and model registry.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.mlflow
```

2. Verify the installation:

```bash
kubectl get pods -n ml-infrastructure
```

### Configuration

1. Configure MLFlow server:

```bash
kubectl apply -f src/ml_infrastructure/mlflow/config/server/mlflow-server.yaml
```

2. Configure MLFlow database:

```bash
kubectl apply -f src/ml_infrastructure/mlflow/config/database/mlflow-db.yaml
```

3. Configure MLFlow experiment:

```bash
kubectl apply -f src/ml_infrastructure/mlflow/config/experiment/llama4-experiment-config.yaml
```

## KServe Setup

KServe is used for model serving with advanced capabilities like canary deployments and autoscaling.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.kserve
```

2. Verify the installation:

```bash
kubectl get pods -n kserve
```

### Configuration

1. Configure KServe namespace:

```bash
kubectl apply -f src/ml_infrastructure/kserve/manifests/namespace/kserve-namespace.yaml
```

2. Configure KServe models:

```bash
kubectl apply -f src/ml_infrastructure/kserve/manifests/models/llama4-model.yaml
```

3. Configure KServe scaling:

```bash
kubectl apply -f src/ml_infrastructure/kserve/manifests/scaling/scaling-config.yaml
```

4. Configure KServe canary deployments:

```bash
kubectl apply -f src/ml_infrastructure/kserve/manifests/canary/canary-deployment.yaml
```

## MinIO Setup

MinIO is used for object storage for artifacts, datasets, and model files.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.minio
```

2. Verify the installation:

```bash
kubectl get pods -n ml-infrastructure
```

### Configuration

1. Create buckets:

```bash
./src/ml_infrastructure/minio/scripts/create-buckets.sh
```

2. Configure lifecycle management:

```bash
./src/ml_infrastructure/minio/scripts/configure-lifecycle.sh
```

3. Set up backup and replication:

```bash
./src/ml_infrastructure/minio/scripts/setup-backup.sh
```

## Feast Setup

Feast is used for feature store management.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.feast
```

2. Verify the installation:

```bash
kubectl get pods -n ml-infrastructure
```

### Configuration

1. Configure feature store:

```bash
kubectl apply -f src/ml_infrastructure/feast/config/feature-store.yaml
```

2. Configure features:

```bash
kubectl apply -f src/ml_infrastructure/feast/config/features/issue_features.py
```

## Seldon Core Setup

Seldon Core is used for advanced model serving and inference pipelines.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.seldon
```

2. Verify the installation:

```bash
kubectl get pods -n seldon-system
```

### Configuration

1. Configure Seldon Core deployment:

```bash
kubectl apply -f src/ml_infrastructure/seldon/manifests/seldon-deployment.yaml
```

## JupyterLab Setup

JupyterLab is used for interactive development and experimentation.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.jupyterhub
```

2. Verify the installation:

```bash
kubectl get pods -n ml-infrastructure
```

### Configuration

1. Configure JupyterLab deployment:

```bash
kubectl apply -f src/ml_infrastructure/jupyter/config/jupyter-deployment.yaml
```

## h2o.ai Setup

h2o.ai is used for automated machine learning capabilities.

### Installation Steps

1. Apply the Terraform configuration:

```bash
cd terraform
terraform apply -target=module.h2o
```

2. Verify the installation:

```bash
kubectl get pods -n ml-infrastructure
```

### Configuration

1. Configure h2o.ai deployment:

```bash
kubectl apply -f src/ml_infrastructure/h2o/config/h2o-deployment.yaml
```

## Data Pipeline Setup

The data pipeline is used to collect, process, and prepare training data from GitHub and Gitee repositories.

### Installation Steps

1. Configure the data pipeline:

```bash
kubectl apply -f src/ml_infrastructure/integration/scrapers/connector.py
kubectl apply -f src/ml_infrastructure/data/preprocessing/preprocessor.py
kubectl apply -f src/ml_infrastructure/data/validation/validator.py
kubectl apply -f src/ml_infrastructure/data/versioning/versioner.py
```

2. Verify the installation:

```bash
kubectl get pods -n ml-infrastructure
```
