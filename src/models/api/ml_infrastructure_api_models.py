"""
Pydantic models for ML Infrastructure API.

This module provides Pydantic models for the ML Infrastructure API client,
ensuring type safety and validation for API requests and responses.
"""

from typing import Dict, List, Optional, Any, Union
from datetime import datetime
from pydantic import BaseModel, Field, HttpUrl


class ModelList(BaseModel):
    """Model list response."""
    models: List[str] = Field(..., description="List of available models")


class ModelDetail(BaseModel):
    """Model detail response."""
    id: str = Field(..., description="Model ID")
    name: str = Field(..., description="Model name")
    version: str = Field(..., description="Model version")
    description: Optional[str] = Field(None, description="Model description")
    parameters: Dict[str, Any] = Field(default_factory=dict, description="Model parameters")
    created_at: datetime = Field(..., description="Creation timestamp")
    updated_at: datetime = Field(..., description="Last update timestamp")


class HyperParameters(BaseModel):
    """Hyperparameters for fine-tuning."""
    epochs: int = Field(3, description="Number of epochs", ge=1, le=10)
    batch_size: int = Field(4, description="Batch size", ge=1, le=64)
    learning_rate: float = Field(1e-5, description="Learning rate", gt=0, le=1)
    weight_decay: Optional[float] = Field(None, description="Weight decay")
    warmup_steps: Optional[int] = Field(None, description="Warmup steps")
    max_grad_norm: Optional[float] = Field(None, description="Maximum gradient norm")


class FineTuningJobCreate(BaseModel):
    """Fine-tuning job creation request."""
    model_id: str = Field(..., description="Base model ID")
    training_file: str = Field(..., description="Training file path or ID")
    validation_file: Optional[str] = Field(None, description="Validation file path or ID")
    hyperparameters: HyperParameters = Field(default_factory=HyperParameters, description="Hyperparameters")
    suffix: Optional[str] = Field(None, description="Model name suffix")
    compute_config: Optional[Dict[str, Any]] = Field(None, description="Compute configuration")


class FineTuningJobDetail(BaseModel):
    """Fine-tuning job detail response."""
    id: str = Field(..., description="Job ID")
    model_id: str = Field(..., description="Base model ID")
    status: str = Field(..., description="Job status")
    created_at: datetime = Field(..., description="Creation timestamp")
    updated_at: datetime = Field(..., description="Last update timestamp")
    fine_tuned_model: Optional[str] = Field(None, description="Fine-tuned model ID")
    training_file: str = Field(..., description="Training file path or ID")
    validation_file: Optional[str] = Field(None, description="Validation file path or ID")
    hyperparameters: HyperParameters = Field(..., description="Hyperparameters")
    metrics: Optional[Dict[str, Any]] = Field(None, description="Training metrics")
    error: Optional[str] = Field(None, description="Error message if failed")


class InferenceServiceCreate(BaseModel):
    """Inference service creation request."""
    name: str = Field(..., description="Service name")
    model_id: str = Field(..., description="Model ID")
    replicas: int = Field(1, description="Number of replicas", ge=1, le=10)
    resources: Dict[str, Any] = Field(default_factory=dict, description="Resource requirements")
    scaling_config: Optional[Dict[str, Any]] = Field(None, description="Scaling configuration")
    timeout: Optional[int] = Field(None, description="Timeout in seconds")


class InferenceServiceDetail(BaseModel):
    """Inference service detail response."""
    id: str = Field(..., description="Service ID")
    name: str = Field(..., description="Service name")
    model_id: str = Field(..., description="Model ID")
    status: str = Field(..., description="Service status")
    url: Optional[HttpUrl] = Field(None, description="Service URL")
    created_at: datetime = Field(..., description="Creation timestamp")
    updated_at: datetime = Field(..., description="Last update timestamp")
    replicas: int = Field(..., description="Number of replicas")
    resources: Dict[str, Any] = Field(..., description="Resource requirements")
    scaling_config: Optional[Dict[str, Any]] = Field(None, description="Scaling configuration")
    error: Optional[str] = Field(None, description="Error message if failed")
