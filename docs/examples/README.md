# ML Infrastructure Examples

This document provides examples for common workflows using the ML infrastructure.

## Table of Contents

- [Fine-Tuning Llama 4 Maverick](#fine-tuning-llama-4-maverick)
- [Fine-Tuning Llama 4 Scout](#fine-tuning-llama-4-scout)
- [Serving Models with KServe](#serving-models-with-kserve)
- [Evaluating Model Performance](#evaluating-model-performance)
- [Data Collection and Preprocessing](#data-collection-and-preprocessing)
- [Hyperparameter Tuning](#hyperparameter-tuning)
- [Model Registry Management](#model-registry-management)
- [Infrastructure Monitoring](#infrastructure-monitoring)

## Fine-Tuning Llama 4 Maverick

This example demonstrates how to fine-tune a Llama 4 Maverick model on GitHub issues data.

### Prerequisites

- Access to the ML infrastructure
- GitHub issues dataset
- Llama 4 Maverick base model

### Steps

1. **Prepare the dataset**

```python
from ml_infrastructure.data.preprocessing import DataPreprocessor
from ml_infrastructure.data.validation import DataValidator
from ml_infrastructure.data.versioning import DatasetVersioner

# Initialize components
preprocessor = DataPreprocessor()
validator = DataValidator()
versioner = DatasetVersioner()

# Load and preprocess data
raw_data = preprocessor.load_data("github-issues-raw")
processed_data = preprocessor.preprocess(raw_data)

# Validate data
validation_result = validator.validate(processed_data)
if validation_result.is_valid:
    # Version the dataset
    dataset_id = versioner.create_version(
        processed_data,
        name="github-issues-gitops",
        description="GitHub issues from GitOps repositories",
        version="v1.0.0"
    )
    print(f"Dataset created with ID: {dataset_id}")
else:
    print("Data validation failed:")
    for error in validation_result.errors:
        print(f"- {error}")
```

2. **Configure the training job**

```python
from ml_infrastructure.training.config import Llama4MaverickConfig
from ml_infrastructure.mlflow.config.client import MLFlowClient

# Initialize MLFlow client
mlflow_client = MLFlowClient()

# Create experiment
experiment_id = mlflow_client.create_experiment(
    name="llama4-maverick-gitops",
    artifact_location="s3://mlflow-artifacts/llama4-maverick-gitops"
)

# Configure training
config = Llama4MaverickConfig(
    model_name="meta-llama/Llama-4-Maverick",
    dataset_id="github-issues-gitops",
    learning_rate=2e-5,
    batch_size=8,
    epochs=3,
    use_lora=True,
    lora_r=16,
    lora_alpha=32,
    lora_dropout=0.05,
    max_seq_length=2048,
    gradient_accumulation_steps=4,
    warmup_steps=100,
    weight_decay=0.01,
    optimizer="adamw_torch",
    scheduler="cosine",
    fp16=True,
    evaluation_strategy="steps",
    eval_steps=500,
    save_steps=500,
    logging_steps=100
)

# Save configuration
config_path = config.save("configs/llama4-maverick-gitops.json")
print(f"Configuration saved to: {config_path}")
```

3. **Submit the training job**

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize API client
client = MLInfrastructureClient()

# Submit training job
job_id = client.training.start_job(
    model_type="maverick",
    dataset_id="github-issues-gitops",
    config_path=config_path,
    experiment_name="llama4-maverick-gitops",
    description="Fine-tuning Llama 4 Maverick on GitOps issues"
)

print(f"Training job submitted with ID: {job_id}")

# Monitor training progress
client.training.monitor_job(job_id)
```

4. **Register the model**

```python
# After training completes
model_id = client.models.register(
    model_type="maverick",
    training_job_id=job_id,
    description="Fine-tuned Llama 4 Maverick model for GitOps issues",
    stage="staging"
)

print(f"Model registered with ID: {model_id}")
```

## Fine-Tuning Llama 4 Scout

This example demonstrates how to fine-tune a Llama 4 Scout model on Terraform issues data.

### Prerequisites

- Access to the ML infrastructure
- Terraform issues dataset
- Llama 4 Scout base model

### Steps

1. **Prepare the dataset**

```python
from ml_infrastructure.data.preprocessing import DataPreprocessor
from ml_infrastructure.data.validation import DataValidator
from ml_infrastructure.data.versioning import DatasetVersioner

# Initialize components
preprocessor = DataPreprocessor()
validator = DataValidator()
versioner = DatasetVersioner()

# Load and preprocess data
raw_data = preprocessor.load_data("terraform-issues-raw")
processed_data = preprocessor.preprocess(raw_data)

# Validate data
validation_result = validator.validate(processed_data)
if validation_result.is_valid:
    # Version the dataset
    dataset_id = versioner.create_version(
        processed_data,
        name="terraform-issues",
        description="Terraform issues from GitHub",
        version="v1.0.0"
    )
    print(f"Dataset created with ID: {dataset_id}")
else:
    print("Data validation failed:")
    for error in validation_result.errors:
        print(f"- {error}")
```

2. **Configure the training job**

```python
from ml_infrastructure.training.config import Llama4ScoutConfig
from ml_infrastructure.mlflow.config.client import MLFlowClient

# Initialize MLFlow client
mlflow_client = MLFlowClient()

# Create experiment
experiment_id = mlflow_client.create_experiment(
    name="llama4-scout-terraform",
    artifact_location="s3://mlflow-artifacts/llama4-scout-terraform"
)

# Configure training
config = Llama4ScoutConfig(
    model_name="meta-llama/Llama-4-Scout",
    dataset_id="terraform-issues",
    learning_rate=1e-5,
    batch_size=4,
    epochs=5,
    use_lora=True,
    lora_r=8,
    lora_alpha=16,
    lora_dropout=0.1,
    max_seq_length=4096,
    gradient_accumulation_steps=8,
    warmup_steps=200,
    weight_decay=0.01,
    optimizer="adamw_torch",
    scheduler="cosine",
    fp16=True,
    evaluation_strategy="steps",
    eval_steps=250,
    save_steps=250,
    logging_steps=50
)

# Save configuration
config_path = config.save("configs/llama4-scout-terraform.json")
print(f"Configuration saved to: {config_path}")
```

3. **Submit the training job**

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize API client
client = MLInfrastructureClient()

# Submit training job
job_id = client.training.start_job(
    model_type="scout",
    dataset_id="terraform-issues",
    config_path=config_path,
    experiment_name="llama4-scout-terraform",
    description="Fine-tuning Llama 4 Scout on Terraform issues"
)

print(f"Training job submitted with ID: {job_id}")

# Monitor training progress
client.training.monitor_job(job_id)
```

4. **Register the model**

```python
# After training completes
model_id = client.models.register(
    model_type="scout",
    training_job_id=job_id,
    description="Fine-tuned Llama 4 Scout model for Terraform issues",
    stage="staging"
)

print(f"Model registered with ID: {model_id}")
```

## Serving Models with KServe

This example demonstrates how to serve a fine-tuned Llama 4 model using KServe.

### Prerequisites

- Access to the ML infrastructure
- Fine-tuned Llama 4 model in the model registry

### Steps

1. **Deploy the model**

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize API client
client = MLInfrastructureClient()

# Deploy model
deployment_id = client.serving.deploy_model(
    model_id="llama4-maverick-v1",
    deployment_name="llama4-maverick-gitops",
    replicas=2,
    resources={
        "requests": {
            "memory": "8Gi",
            "cpu": "2",
            "nvidia.com/gpu": "1"
        },
        "limits": {
            "memory": "16Gi",
            "cpu": "4",
            "nvidia.com/gpu": "1"
        }
    },
    autoscaling={
        "enabled": True,
        "min_replicas": 1,
        "max_replicas": 5,
        "target_concurrency": 10
    }
)

print(f"Model deployed with ID: {deployment_id}")
```

2. **Test the deployment**

```python
# Test the deployment
response = client.serving.predict(
    deployment_id=deployment_id,
    input={
        "repository": "kubernetes/kubernetes",
        "issue_title": "Pod fails to start with CrashLoopBackOff",
        "issue_description": "I'm trying to deploy a pod but it keeps failing with CrashLoopBackOff. The logs show that the container is exiting with code 1."
    },
    parameters={
        "temperature": 0.7,
        "max_tokens": 1024,
        "top_p": 0.9,
        "top_k": 50
    }
)

print("Model response:")
print(response.output)
```

3. **Configure canary deployment**

```python
# Deploy a canary version
canary_deployment_id = client.serving.deploy_canary(
    deployment_id=deployment_id,
    model_id="llama4-maverick-v2",
    traffic_percent=20
)

print(f"Canary deployment created with ID: {canary_deployment_id}")
```

4. **Monitor the deployment**

```python
# Get deployment metrics
metrics = client.serving.get_metrics(deployment_id)
print("Deployment metrics:")
print(f"- Request count: {metrics.request_count}")
print(f"- Average latency: {metrics.avg_latency} ms")
print(f"- P95 latency: {metrics.p95_latency} ms")
print(f"- Error rate: {metrics.error_rate}%")
```

## Evaluating Model Performance

This example demonstrates how to evaluate the performance of a fine-tuned Llama 4 model.

### Prerequisites

- Access to the ML infrastructure
- Fine-tuned Llama 4 model
- Evaluation dataset

### Steps

1. **Prepare the evaluation dataset**

```python
from ml_infrastructure.data.preprocessing import DataPreprocessor
from ml_infrastructure.data.validation import DataValidator

# Initialize components
preprocessor = DataPreprocessor()
validator = DataValidator()

# Load and preprocess evaluation data
raw_data = preprocessor.load_data("github-issues-eval")
eval_data = preprocessor.preprocess(raw_data)

# Validate data
validation_result = validator.validate(eval_data)
if validation_result.is_valid:
    print("Evaluation dataset is valid")
else:
    print("Data validation failed:")
    for error in validation_result.errors:
        print(f"- {error}")
```

2. **Run the evaluation**

```python
from ml_infrastructure.api.client import MLInfrastructureClient
from ml_infrastructure.training.evaluation import ModelEvaluationMetrics

# Initialize API client
client = MLInfrastructureClient()

# Run evaluation
evaluation_id = client.evaluation.start_evaluation(
    model_id="llama4-maverick-v1",
    dataset_id="github-issues-eval",
    metrics=[
        "accuracy",
        "f1",
        "precision",
        "recall",
        "trajectory_similarity"
    ]
)

print(f"Evaluation started with ID: {evaluation_id}")

# Wait for evaluation to complete
evaluation_result = client.evaluation.wait_for_evaluation(evaluation_id)

# Print evaluation results
print("Evaluation results:")
for metric, value in evaluation_result.metrics.items():
    print(f"- {metric}: {value}")
```

3. **Compare models**

```python
# Compare models
comparison = client.evaluation.compare_models(
    model_ids=["llama4-maverick-v1", "llama4-maverick-v2"],
    dataset_id="github-issues-eval",
    metrics=[
        "accuracy",
        "f1",
        "precision",
        "recall",
        "trajectory_similarity"
    ]
)

print("Model comparison:")
for metric, values in comparison.items():
    print(f"- {metric}:")
    for model_id, value in values.items():
        print(f"  - {model_id}: {value}")
```

4. **Generate evaluation report**

```python
# Generate evaluation report
report_path = client.evaluation.generate_report(
    evaluation_id=evaluation_id,
    output_format="html",
    include_examples=True
)

print(f"Evaluation report generated at: {report_path}")
```

## Data Collection and Preprocessing

This example demonstrates how to collect and preprocess data from GitHub and Gitee repositories.

### Prerequisites

- Access to the ML infrastructure
- GitHub and Gitee API credentials

### Steps

1. **Configure the scrapers**

```python
from ml_infrastructure.integration.scrapers import ScraperConnector

# Initialize scraper connector
connector = ScraperConnector(
    github_token=os.environ.get("GITHUB_TOKEN"),
    gitee_token=os.environ.get("GITEE_TOKEN")
)

# Configure repositories to scrape
github_repos = [
    "kubernetes/kubernetes",
    "terraform-providers/terraform-provider-aws",
    "fluxcd/flux2",
    "argoproj/argo-cd"
]

gitee_repos = [
    "openharmony/kernel_linux_5.10",
    "openeuler/iSulad"
]

# Configure filters
filters = {
    "state": "closed",
    "labels": ["bug", "enhancement", "feature"],
    "created_after": "2023-01-01",
    "created_before": "2025-01-01"
}
```

2. **Collect data**

```python
# Collect GitHub issues
github_issues = connector.collect_github_issues(
    repositories=github_repos,
    filters=filters,
    max_issues_per_repo=1000
)

print(f"Collected {len(github_issues)} GitHub issues")

# Collect Gitee issues
gitee_issues = connector.collect_gitee_issues(
    repositories=gitee_repos,
    filters=filters,
    max_issues_per_repo=1000
)

print(f"Collected {len(gitee_issues)} Gitee issues")

# Combine issues
all_issues = github_issues + gitee_issues
print(f"Total issues collected: {len(all_issues)}")
```

3. **Preprocess data**

```python
from ml_infrastructure.data.preprocessing import DataPreprocessor

# Initialize preprocessor
preprocessor = DataPreprocessor()

# Preprocess issues
processed_data = preprocessor.preprocess(all_issues)

print(f"Preprocessed {len(processed_data)} issues")
```

4. **Validate and version data**

```python
from ml_infrastructure.data.validation import DataValidator
from ml_infrastructure.data.versioning import DatasetVersioner

# Initialize components
validator = DataValidator()
versioner = DatasetVersioner()

# Validate data
validation_result = validator.validate(processed_data)
if validation_result.is_valid:
    # Split data into training and evaluation sets
    train_data, eval_data = preprocessor.split_data(
        processed_data,
        train_ratio=0.8,
        random_seed=42
    )
    
    # Version the datasets
    train_dataset_id = versioner.create_version(
        train_data,
        name="issues-train",
        description="Training dataset of GitHub and Gitee issues",
        version="v1.0.0"
    )
    
    eval_dataset_id = versioner.create_version(
        eval_data,
        name="issues-eval",
        description="Evaluation dataset of GitHub and Gitee issues",
        version="v1.0.0"
    )
    
    print(f"Training dataset created with ID: {train_dataset_id}")
    print(f"Evaluation dataset created with ID: {eval_dataset_id}")
else:
    print("Data validation failed:")
    for error in validation_result.errors:
        print(f"- {error}")
```

## Hyperparameter Tuning

This example demonstrates how to perform hyperparameter tuning for a Llama 4 model.

### Prerequisites

- Access to the ML infrastructure
- Training dataset

### Steps

1. **Configure the hyperparameter tuning job**

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize API client
client = MLInfrastructureClient()

# Define hyperparameter search space
search_space = {
    "learning_rate": {
        "type": "double",
        "min": 1e-6,
        "max": 1e-4,
        "scale": "log"
    },
    "batch_size": {
        "type": "int",
        "min": 1,
        "max": 16,
        "step": 1
    },
    "lora_r": {
        "type": "int",
        "min": 4,
        "max": 32,
        "step": 4
    },
    "lora_alpha": {
        "type": "int",
        "min": 8,
        "max": 64,
        "step": 8
    },
    "lora_dropout": {
        "type": "double",
        "min": 0.0,
        "max": 0.2,
        "step": 0.05
    }
}

# Configure tuning job
tuning_job_id = client.tuning.start_job(
    model_type="maverick",
    dataset_id="github-issues-gitops",
    search_space=search_space,
    objective_metric="eval_loss",
    objective_type="minimize",
    max_trials=20,
    parallel_trials=4,
    experiment_name="llama4-maverick-hparam-tuning",
    description="Hyperparameter tuning for Llama 4 Maverick"
)

print(f"Hyperparameter tuning job submitted with ID: {tuning_job_id}")
```

2. **Monitor the tuning job**

```python
# Monitor tuning progress
client.tuning.monitor_job(tuning_job_id)

# Get tuning results
tuning_results = client.tuning.get_results(tuning_job_id)

print("Hyperparameter tuning results:")
print(f"Best trial: {tuning_results.best_trial_id}")
print("Best parameters:")
for param, value in tuning_results.best_parameters.items():
    print(f"- {param}: {value}")
print(f"Best objective value: {tuning_results.best_objective_value}")
```

3. **Train with the best parameters**

```python
# Train with the best parameters
job_id = client.training.start_job(
    model_type="maverick",
    dataset_id="github-issues-gitops",
    parameters=tuning_results.best_parameters,
    experiment_name="llama4-maverick-best-params",
    description="Training Llama 4 Maverick with best parameters from hyperparameter tuning"
)

print(f"Training job submitted with ID: {job_id}")

# Monitor training progress
client.training.monitor_job(job_id)
```

## Model Registry Management

This example demonstrates how to manage models in the model registry.

### Prerequisites

- Access to the ML infrastructure
- Fine-tuned models

### Steps

1. **List models in the registry**

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize API client
client = MLInfrastructureClient()

# List models
models = client.models.list_models()

print(f"Found {len(models)} models in the registry:")
for model in models:
    print(f"- {model.model_id} ({model.model_type}): {model.description}")
```

2. **Get model details**

```python
# Get model details
model_id = "llama4-maverick-v1"
model = client.models.get_model(model_id)

print(f"Model: {model.model_id}")
print(f"Type: {model.model_type}")
print(f"Created at: {model.created_at}")
print(f"Status: {model.status}")
print(f"Description: {model.description}")
print(f"Stage: {model.stage}")
print(f"Version: {model.version}")
print("Metrics:")
for metric, value in model.metrics.items():
    print(f"- {metric}: {value}")
```

3. **Update model stage**

```python
# Update model stage
client.models.update_stage(
    model_id="llama4-maverick-v1",
    stage="production"
)

print(f"Updated model stage to 'production'")
```

4. **Compare models**

```python
# Compare models
comparison = client.models.compare(
    model_ids=["llama4-maverick-v1", "llama4-maverick-v2"]
)

print("Model comparison:")
for metric, values in comparison.metrics.items():
    print(f"- {metric}:")
    for model_id, value in values.items():
        print(f"  - {model_id}: {value}")
```

5. **Archive a model**

```python
# Archive a model
client.models.archive(model_id="llama4-maverick-v1")

print(f"Archived model 'llama4-maverick-v1'")
```

## Infrastructure Monitoring

This example demonstrates how to monitor the ML infrastructure.

### Prerequisites

- Access to the ML infrastructure

### Steps

1. **Get infrastructure status**

```python
from ml_infrastructure.api.client import MLInfrastructureClient

# Initialize API client
client = MLInfrastructureClient()

# Get infrastructure status
status = client.infrastructure.get_status()

print(f"Infrastructure status: {status.status}")
print("Component status:")
for component, info in status.components.items():
    print(f"- {component}: {info.status} (version: {info.version})")
print("Resource usage:")
for resource, usage in status.resources.items():
    print(f"- {resource}: {usage}")
```

2. **Get infrastructure metrics**

```python
# Get infrastructure metrics
metrics = client.infrastructure.get_metrics(
    component="all",
    metric="all",
    period="24h"
)

print("Infrastructure metrics:")
for metric_name, values in metrics.metrics.items():
    print(f"- {metric_name}:")
    for timestamp, value in values[-3:]:  # Show last 3 values
        print(f"  - {timestamp}: {value}")
```

3. **Get training job metrics**

```python
# Get training job metrics
job_id = "train-123456"
job_metrics = client.training.get_metrics(job_id)

print(f"Training job metrics for {job_id}:")
for metric_name, values in job_metrics.items():
    print(f"- {metric_name}:")
    for timestamp, value in values[-3:]:  # Show last 3 values
        print(f"  - {timestamp}: {value}")
```

4. **Get model serving metrics**

```python
# Get model serving metrics
deployment_id = "deploy-123456"
serving_metrics = client.serving.get_metrics(deployment_id)

print(f"Model serving metrics for {deployment_id}:")
print(f"- Request count: {serving_metrics.request_count}")
print(f"- Average latency: {serving_metrics.avg_latency} ms")
print(f"- P95 latency: {serving_metrics.p95_latency} ms")
print(f"- Error rate: {serving_metrics.error_rate}%")
print("- Request count over time:")
for timestamp, value in serving_metrics.request_count_over_time[-3:]:
    print(f"  - {timestamp}: {value}")
```
