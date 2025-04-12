# ML Infrastructure Documentation

This documentation provides comprehensive information about the ML infrastructure for fine-tuning and evaluating Llama 4 models on GitOps, Terraform, and Kubernetes issue data.

## Table of Contents

- [Architecture](./architecture/README.md): Architecture diagrams and component descriptions
- [Setup Guides](./guides/README.md): Installation and configuration guides for all components
- [API Documentation](./api/README.md): API endpoints and parameters
- [Examples](./examples/README.md): Example workflows and use cases

## Overview

The ML infrastructure is designed to support fine-tuning and evaluation of Llama 4 Maverick and Scout models on datasets created from GitHub and Gitee repositories. The infrastructure follows a modular architecture with components deployed on Kubernetes, using a combination of Kubernetes manifests and Terraform configurations.

### Key Components

- **KubeFlow**: Orchestrates ML workflows, including training, hyperparameter tuning, and pipelines
- **MLFlow**: Tracks experiments and manages model registry
- **KServe**: Serves models with advanced capabilities like canary deployments and autoscaling
- **MinIO**: Provides object storage for artifacts, datasets, and model files
- **Feast**: Manages feature stores for ML features
- **Seldon Core**: Enables advanced model serving and inference pipelines
- **JupyterLab**: Supports interactive development and experimentation
- **h2o.ai**: Provides automated machine learning capabilities

### Data Pipeline

The infrastructure ingests training data from GitHub and Gitee issue scrapers, which collect solved issues from GitOps, Terraform, and Kubernetes repositories. The data is processed, validated, and versioned before being used for model training.

### Model Training

The infrastructure supports fine-tuning Llama 4 models using a variety of techniques, including:

- Parameter-efficient fine-tuning (PEFT) with LoRA
- Instruction fine-tuning
- Trajectory-based fine-tuning

### Model Evaluation

Models are evaluated using a combination of metrics, including:

- Standard metrics (accuracy, F1 score, etc.)
- Trajectory similarity
- SWE Agent benchmarking

### Model Serving

Trained models are served using KServe, which provides:

- Autoscaling
- Canary deployments
- A/B testing
- Monitoring and logging

## Getting Started

To get started with the ML infrastructure, see the [Setup Guides](./guides/README.md).
