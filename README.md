# Agent Runtime with ML Infrastructure

A comprehensive Go/Python implementation of SWE-Agent and SWE-ReX frameworks with integrated ML infrastructure for fine-tuning and evaluating Llama 4 models.

## Overview

This repository provides a high-performance runtime environment for autonomous software engineering agents with integrated ML infrastructure. The system enables:

- Creation, configuration, and execution of AI-powered agents for software engineering tasks
- Fine-tuning Llama 4 models (Maverick and Scout variants) using data from GitHub and Gitee repositories
- Deployment in containerized environments (Kubernetes with Kata Containers)
- Robust infrastructure management, state tracking, and event distribution

The infrastructure includes:

- Web scrapers for collecting solved issues from GitOps, Terraform, and Kubernetes repositories
- Data pipeline for preprocessing and validating training data
- ML infrastructure components (KubeFlow, MLflow, KServe)
- API client for interacting with the ML infrastructure
- Terraform modules for infrastructure deployment

## Agent Runtime Core Features

- **Agent System**: Manages the agent's execution loop, state tracking, and command execution
- **Tool Registry**: Provides a structured interface for defining and executing tools
- **FFI System**: Enables execution of Python and C++ code from Go
- **Event Stream System**: Distributes events between components with DragonflyDB caching
- **MCP (Model Control Plane)**: Provides interfaces for models and tools to communicate

## ML Infrastructure API Client

The ML Infrastructure API Client provides a Python interface for interacting with the ML infrastructure for Llama 4 fine-tuning.

### Installation

```bash
# Clone the repository
git clone https://github.com/spectrumwebco/agent_runtime.git
cd agent_runtime

# Install dependencies
pip install -r requirements.txt

# Install the package in development mode
pip install -e .
```

### Environment Setup

Create a `.env` file based on the provided `.env.example`:

```bash
cp .env.example .env
```

Edit the `.env` file to set your API credentials:

```
# ML Infrastructure API Configuration
ML_API_BASE_URL=http://your-api-server:8000
ML_API_USERNAME=your_username
ML_API_PASSWORD=your_password

# Other configuration...
```

### Usage

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize client using environment variables
client = MLInfrastructureClient()

# Or provide credentials directly
client = MLInfrastructureClient(
    base_url="http://your-api-server:8000",
    username="your_username",
    password="your_password"
)

# Get API status
status = client.get_api_status()
print(f"API Status: {status}")

# Create a dataset
dataset = client.create_dataset(
    name="gitops-terraform-k8s-issues",
    description="Solved issues from GitOps, Terraform, and Kubernetes repositories",
    source="github-gitee",
    version="1.0.0",
    metadata={
        "topics": ["gitops", "terraform", "kubernetes"],
        "issue_type": "solved"
    }
)

# Submit a training job
job = client.submit_training_job(
    name="llama4-maverick-fine-tuning",
    model_type="llama4-maverick",
    dataset_id=dataset["id"],
    config={
        "training_type": "fine-tuning",
        "epochs": 3,
        "batch_size": 8,
        "learning_rate": 2e-5
    }
)

# Monitor training progress
job_status = client.get_training_job(job["id"])
print(f"Job status: {job_status}")
```

For more detailed examples, see the example scripts in `src/ml_infrastructure/api/examples/`.

## Components

- **Agent System**: Core agent implementation with execution loop and state management
- **Tool Registry**: Registry for managing and executing tools
- **Web Scrapers**: Collect solved issues from GitHub and Gitee repositories
- **Data Pipeline**: Preprocess and validate training data
- **k8s/**: All Kubernetes manifests (KubeFlow, MLflow, KServe, etc.)
- **terraform/**: Infrastructure as Code configurations
- **src/models/**: Pydantic models for type safety and validation

## Directory Structure

- **cmd/**: Contains CLI entry points
- **pkg/**: Contains core packages intended for external use
- **internal/**: Contains packages for internal use
- **k8s/**: All Kubernetes manifests and configurations
  - kubeflow/: KubeFlow manifests for orchestrating ML workflows
  - mlflow/: MLflow configurations for experiment tracking
  - kserve/: KServe manifests for model serving
  - minio/: MinIO configurations for artifact storage
  - monitoring/: Prometheus, Grafana, and Loki for monitoring
  - argocd/: ArgoCD configurations for GitOps
  
- **terraform/**: Infrastructure as Code
  - modules/: Terraform modules for each component
  - main.tf: Main Terraform configuration

- **src/models/**: Pydantic models for type safety
  - api/: API models
  - data_validation/: Data validation models
  - integrations/: Integration models (GitHub, Gitee)
  - ml_infrastructure/: ML infrastructure models

## License

[MIT License](LICENSE)
