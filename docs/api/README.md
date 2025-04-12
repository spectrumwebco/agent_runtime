# ML Infrastructure API Documentation

This document provides comprehensive documentation for the ML Infrastructure API, which enables interaction with the fine-tuning and evaluation infrastructure for Llama 4 models.

## Table of Contents

- [Authentication](#authentication)
- [Base URL](#base-url)
- [API Endpoints](#api-endpoints)
  - [Training](#training)
  - [Inference](#inference)
  - [Model Management](#model-management)
  - [Data Management](#data-management)
  - [Infrastructure Management](#infrastructure-management)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Versioning](#versioning)

## Authentication

The ML Infrastructure API uses API key authentication. You must include your API key in the request headers.

```bash
curl -X GET https://api.example.com/v1/models \
  -H "Content-Type: application/json" \
  -H "x-api-key: your-api-key"
```

## Base URL

The base URL for the ML Infrastructure API is:

```
https://api.example.com/v1
```

## API Endpoints

### Training

#### Start Training Job

```
POST /training/jobs
```

Start a new training job for fine-tuning a Llama 4 model.

**Request Body:**

```json
{
  "model_type": "maverick",
  "dataset_id": "github-issues-v1",
  "config": {
    "learning_rate": 2e-5,
    "batch_size": 8,
    "epochs": 3,
    "use_lora": true,
    "lora_r": 16,
    "lora_alpha": 32,
    "lora_dropout": 0.05
  },
  "experiment_name": "llama4-maverick-gitops-issues",
  "description": "Fine-tuning Llama 4 Maverick on GitOps issues"
}
```

**Response:**

```json
{
  "job_id": "train-123456",
  "status": "queued",
  "created_at": "2025-04-11T18:30:00Z",
  "model_type": "maverick",
  "dataset_id": "github-issues-v1",
  "experiment_name": "llama4-maverick-gitops-issues"
}
```

#### Get Training Job Status

```
GET /training/jobs/{job_id}
```

Get the status of a training job.

**Response:**

```json
{
  "job_id": "train-123456",
  "status": "running",
  "created_at": "2025-04-11T18:30:00Z",
  "started_at": "2025-04-11T18:31:00Z",
  "model_type": "maverick",
  "dataset_id": "github-issues-v1",
  "experiment_name": "llama4-maverick-gitops-issues",
  "progress": {
    "current_epoch": 2,
    "total_epochs": 3,
    "current_step": 1500,
    "total_steps": 2000,
    "loss": 1.234,
    "learning_rate": 1.5e-5
  }
}
```

#### List Training Jobs

```
GET /training/jobs
```

List all training jobs.

**Query Parameters:**

- `status` (optional): Filter by status (queued, running, completed, failed)
- `model_type` (optional): Filter by model type (maverick, scout)
- `limit` (optional): Limit the number of results (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**

```json
{
  "jobs": [
    {
      "job_id": "train-123456",
      "status": "running",
      "created_at": "2025-04-11T18:30:00Z",
      "model_type": "maverick",
      "dataset_id": "github-issues-v1",
      "experiment_name": "llama4-maverick-gitops-issues"
    },
    {
      "job_id": "train-123457",
      "status": "queued",
      "created_at": "2025-04-11T18:35:00Z",
      "model_type": "scout",
      "dataset_id": "terraform-issues-v1",
      "experiment_name": "llama4-scout-terraform-issues"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

#### Cancel Training Job

```
POST /training/jobs/{job_id}/cancel
```

Cancel a running training job.

**Response:**

```json
{
  "job_id": "train-123456",
  "status": "cancelling",
  "message": "Job cancellation in progress"
}
```

### Inference

#### Run Inference

```
POST /inference
```

Run inference using a fine-tuned Llama 4 model.

**Request Body:**

```json
{
  "model_id": "llama4-maverick-v1",
  "input": {
    "repository": "kubernetes/kubernetes",
    "issue_title": "Pod fails to start with CrashLoopBackOff",
    "issue_description": "I'm trying to deploy a pod but it keeps failing with CrashLoopBackOff. The logs show that the container is exiting with code 1."
  },
  "parameters": {
    "temperature": 0.7,
    "max_tokens": 1024,
    "top_p": 0.9,
    "top_k": 50
  }
}
```

**Response:**

```json
{
  "model_id": "llama4-maverick-v1",
  "output": "The issue is likely due to an error in your container configuration. Here are some steps to troubleshoot:\n\n1. Check the container logs using `kubectl logs <pod-name>`\n2. Verify that the command and arguments in your pod spec are correct\n3. Ensure that the container image exists and is accessible\n4. Check resource constraints (CPU/memory)\n\nBased on the limited information, I would recommend checking the container logs first to see the exact error message.",
  "metadata": {
    "tokens_generated": 512,
    "generation_time": 2.5,
    "model_version": "v1.0.0"
  }
}
```

#### Get Model Information

```
GET /models/{model_id}
```

Get information about a specific model.

**Response:**

```json
{
  "model_id": "llama4-maverick-v1",
  "model_type": "maverick",
  "created_at": "2025-04-10T12:00:00Z",
  "status": "active",
  "description": "Fine-tuned Llama 4 Maverick model for GitOps issues",
  "metrics": {
    "accuracy": 0.85,
    "f1": 0.87,
    "trajectory_similarity": 0.78
  },
  "training_job_id": "train-123456",
  "version": "v1.0.0"
}
```

#### List Models

```
GET /models
```

List all available models.

**Query Parameters:**

- `model_type` (optional): Filter by model type (maverick, scout)
- `status` (optional): Filter by status (active, inactive)
- `limit` (optional): Limit the number of results (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**

```json
{
  "models": [
    {
      "model_id": "llama4-maverick-v1",
      "model_type": "maverick",
      "created_at": "2025-04-10T12:00:00Z",
      "status": "active",
      "description": "Fine-tuned Llama 4 Maverick model for GitOps issues"
    },
    {
      "model_id": "llama4-scout-v1",
      "model_type": "scout",
      "created_at": "2025-04-09T10:00:00Z",
      "status": "active",
      "description": "Fine-tuned Llama 4 Scout model for Terraform issues"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

### Model Management

#### Register Model

```
POST /models/register
```

Register a model in the model registry.

**Request Body:**

```json
{
  "model_id": "llama4-maverick-v1",
  "model_type": "maverick",
  "description": "Fine-tuned Llama 4 Maverick model for GitOps issues",
  "training_job_id": "train-123456",
  "metrics": {
    "accuracy": 0.85,
    "f1": 0.87,
    "trajectory_similarity": 0.78
  },
  "stage": "staging"
}
```

**Response:**

```json
{
  "model_id": "llama4-maverick-v1",
  "model_type": "maverick",
  "created_at": "2025-04-11T19:00:00Z",
  "status": "active",
  "description": "Fine-tuned Llama 4 Maverick model for GitOps issues",
  "stage": "staging",
  "version": "v1.0.0"
}
```

#### Update Model Stage

```
POST /models/{model_id}/stage
```

Update the stage of a model in the model registry.

**Request Body:**

```json
{
  "stage": "production"
}
```

**Response:**

```json
{
  "model_id": "llama4-maverick-v1",
  "model_type": "maverick",
  "status": "active",
  "description": "Fine-tuned Llama 4 Maverick model for GitOps issues",
  "stage": "production",
  "version": "v1.0.0",
  "updated_at": "2025-04-11T19:30:00Z"
}
```

#### Delete Model

```
DELETE /models/{model_id}
```

Delete a model from the model registry.

**Response:**

```json
{
  "model_id": "llama4-maverick-v1",
  "status": "deleted",
  "message": "Model deleted successfully"
}
```

### Data Management

#### Upload Dataset

```
POST /datasets/upload
```

Upload a dataset for training.

**Request Body:**

```json
{
  "name": "github-issues-v1",
  "description": "GitHub issues dataset for GitOps repositories",
  "source": "github",
  "format": "json",
  "version": "v1.0.0"
}
```

**Response:**

```json
{
  "dataset_id": "github-issues-v1",
  "name": "github-issues-v1",
  "description": "GitHub issues dataset for GitOps repositories",
  "source": "github",
  "format": "json",
  "version": "v1.0.0",
  "created_at": "2025-04-11T20:00:00Z",
  "status": "pending",
  "upload_url": "https://api.example.com/v1/datasets/upload/123456"
}
```

#### Get Dataset Information

```
GET /datasets/{dataset_id}
```

Get information about a specific dataset.

**Response:**

```json
{
  "dataset_id": "github-issues-v1",
  "name": "github-issues-v1",
  "description": "GitHub issues dataset for GitOps repositories",
  "source": "github",
  "format": "json",
  "version": "v1.0.0",
  "created_at": "2025-04-11T20:00:00Z",
  "status": "ready",
  "size": 1024000,
  "num_examples": 5000,
  "metadata": {
    "repositories": ["kubernetes/kubernetes", "terraform-providers/terraform-provider-aws"],
    "issue_types": ["bug", "feature", "question"],
    "date_range": ["2023-01-01", "2025-01-01"]
  }
}
```

#### List Datasets

```
GET /datasets
```

List all available datasets.

**Query Parameters:**

- `source` (optional): Filter by source (github, gitee)
- `status` (optional): Filter by status (pending, ready, failed)
- `limit` (optional): Limit the number of results (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**

```json
{
  "datasets": [
    {
      "dataset_id": "github-issues-v1",
      "name": "github-issues-v1",
      "description": "GitHub issues dataset for GitOps repositories",
      "source": "github",
      "format": "json",
      "version": "v1.0.0",
      "created_at": "2025-04-11T20:00:00Z",
      "status": "ready"
    },
    {
      "dataset_id": "terraform-issues-v1",
      "name": "terraform-issues-v1",
      "description": "Terraform issues dataset",
      "source": "github",
      "format": "json",
      "version": "v1.0.0",
      "created_at": "2025-04-10T15:00:00Z",
      "status": "ready"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

#### Delete Dataset

```
DELETE /datasets/{dataset_id}
```

Delete a dataset.

**Response:**

```json
{
  "dataset_id": "github-issues-v1",
  "status": "deleted",
  "message": "Dataset deleted successfully"
}
```

### Infrastructure Management

#### Get Infrastructure Status

```
GET /infrastructure/status
```

Get the status of the ML infrastructure.

**Response:**

```json
{
  "status": "healthy",
  "components": {
    "kubeflow": {
      "status": "healthy",
      "version": "v1.7.0"
    },
    "mlflow": {
      "status": "healthy",
      "version": "v2.3.0"
    },
    "kserve": {
      "status": "healthy",
      "version": "v0.10.0"
    },
    "minio": {
      "status": "healthy",
      "version": "RELEASE.2023-06-29T05-12-28Z"
    }
  },
  "resources": {
    "cpu_usage": "45%",
    "memory_usage": "60%",
    "gpu_usage": "30%",
    "storage_usage": "55%"
  }
}
```

#### Get Infrastructure Metrics

```
GET /infrastructure/metrics
```

Get metrics for the ML infrastructure.

**Query Parameters:**

- `component` (optional): Filter by component (kubeflow, mlflow, kserve, minio)
- `metric` (optional): Filter by metric (cpu, memory, gpu, storage)
- `period` (optional): Time period (1h, 6h, 24h, 7d, 30d)

**Response:**

```json
{
  "metrics": {
    "cpu_usage": [
      {"timestamp": "2025-04-11T18:00:00Z", "value": 45},
      {"timestamp": "2025-04-11T19:00:00Z", "value": 50},
      {"timestamp": "2025-04-11T20:00:00Z", "value": 40}
    ],
    "memory_usage": [
      {"timestamp": "2025-04-11T18:00:00Z", "value": 60},
      {"timestamp": "2025-04-11T19:00:00Z", "value": 65},
      {"timestamp": "2025-04-11T20:00:00Z", "value": 55}
    ],
    "gpu_usage": [
      {"timestamp": "2025-04-11T18:00:00Z", "value": 30},
      {"timestamp": "2025-04-11T19:00:00Z", "value": 35},
      {"timestamp": "2025-04-11T20:00:00Z", "value": 25}
    ],
    "storage_usage": [
      {"timestamp": "2025-04-11T18:00:00Z", "value": 55},
      {"timestamp": "2025-04-11T19:00:00Z", "value": 56},
      {"timestamp": "2025-04-11T20:00:00Z", "value": 57}
    ]
  },
  "period": "3h",
  "interval": "1h"
}
```

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of a request.

- 200 OK: The request was successful
- 400 Bad Request: The request was invalid
- 401 Unauthorized: Authentication failed
- 403 Forbidden: The authenticated user does not have permission to access the requested resource
- 404 Not Found: The requested resource was not found
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: An error occurred on the server

Error responses include a JSON object with an error message:

```json
{
  "error": {
    "code": "invalid_request",
    "message": "Invalid request: missing required field 'model_type'",
    "status": 400
  }
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse. Rate limits are applied on a per-API key basis.

- Training endpoints: 10 requests per minute
- Inference endpoints: 100 requests per minute
- Model management endpoints: 60 requests per minute
- Data management endpoints: 60 requests per minute
- Infrastructure management endpoints: 60 requests per minute

Rate limit information is included in the response headers:

- `X-RateLimit-Limit`: The maximum number of requests allowed in a time window
- `X-RateLimit-Remaining`: The number of requests remaining in the current time window
- `X-RateLimit-Reset`: The time at which the current rate limit window resets, in UTC epoch seconds

## Versioning

The API is versioned using the URL path. The current version is `v1`.

```
https://api.example.com/v1/models
```

Future versions will be available at:

```
https://api.example.com/v2/models
```

API changes within a version are backward compatible. Breaking changes will only be introduced in a new API version.
