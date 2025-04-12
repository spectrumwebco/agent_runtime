# Fine-Tune

A comprehensive infrastructure for fine-tuning and evaluating Llama 4 Maverick and Llama 4 Scout models on GitOps, Terraform, and Kubernetes repositories.

## Overview

This repository provides a complete ML infrastructure for fine-tuning Llama 4 models using data from GitHub and Gitee repositories. The infrastructure includes:

- Web scrapers for collecting solved issues from GitOps, Terraform, and Kubernetes repositories
- Data pipeline for preprocessing and validating training data
- ML infrastructure components (KubeFlow, MLflow, KServe)
- API client for interacting with the ML infrastructure
- Terraform modules for infrastructure deployment

## ML Infrastructure API Client

The ML Infrastructure API Client provides a Python interface for interacting with the ML infrastructure for Llama 4 fine-tuning.

### Installation

```bash
# Clone the repository
git clone https://github.com/spectrumwebco/Fine-Tune.git
cd Fine-Tune

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

- **Web Scrapers**: Collect solved issues from GitHub and Gitee repositories
- **Data Pipeline**: Preprocess and validate training data
- **KubeFlow**: Orchestrate ML workflows
- **MLflow**: Track experiments and manage models
- **KServe**: Serve models on Kubernetes
- **Terraform**: Deploy infrastructure components

## License

[MIT License](LICENSE)
