"""
Pydantic models for data validation in Llama 4 fine-tuning.

This module provides Pydantic models for validating data collected from GitHub and Gitee
issue scrapers for fine-tuning Llama 4 models.
"""

from typing import Dict, List, Optional, Any, Union
from datetime import datetime
from pydantic import BaseModel, Field, validator, HttpUrl


class TrajectoryStep(BaseModel):
    """Trajectory step model."""
    step: int = Field(..., description="Step number")
    action: str = Field(..., description="Action taken")
    content: str = Field(..., description="Content of the step")
    timestamp: Optional[str] = Field(None, description="Timestamp of the step")
    user: Optional[str] = Field(None, description="User who performed the action")


class InputData(BaseModel):
    """Input data model for raw data format."""
    repository: Optional[str] = Field(None, description="Repository name")
    topics: List[str] = Field(default_factory=list, description="Repository topics")
    title: str = Field(..., description="Issue title")
    description: str = Field(..., description="Issue description")


class OutputData(BaseModel):
    """Output data model for raw data format."""
    solution: str = Field(..., description="Solution to the issue")


class Metadata(BaseModel):
    """Metadata model for all data formats."""
    id: str = Field(..., description="Unique identifier")
    source: str = Field(..., description="Source of the data")
    repository: Optional[str] = Field(None, description="Repository name")
    url: Optional[str] = Field(None, description="URL to the issue")
    created_at: Optional[str] = Field(None, description="Creation date")
    closed_at: Optional[str] = Field(None, description="Closure date")
    labels: List[str] = Field(default_factory=list, description="Issue labels")


class RawDataModel(BaseModel):
    """Raw data model for validation."""
    input: InputData = Field(..., description="Input data")
    output: OutputData = Field(..., description="Output data")
    metadata: Metadata = Field(..., description="Metadata")
    trajectory: Optional[List[TrajectoryStep]] = Field(None, description="Solution trajectory")


class Message(BaseModel):
    """Message model for chat format."""
    role: str = Field(..., description="Message role")
    content: str = Field(..., description="Message content")

    @validator('role')
    def validate_role(cls, v):
        """Validate message role."""
        if v not in ["system", "user", "assistant"]:
            raise ValueError(f"Invalid role: {v}. Must be one of: system, user, assistant")
        return v


class ChatFormatModel(BaseModel):
    """Chat format model for validation."""
    messages: List[Message] = Field(..., description="Chat messages")
    metadata: Optional[Metadata] = Field(None, description="Metadata")
    trajectory: Optional[List[TrajectoryStep]] = Field(None, description="Solution trajectory")

    @validator('messages')
    def validate_messages(cls, v):
        """Validate chat messages."""
        if len(v) < 3:
            raise ValueError(f"Chat format must have at least 3 messages, got {len(v)}")
        return v


class CompletionFormatModel(BaseModel):
    """Completion format model for validation."""
    prompt: str = Field(..., description="Completion prompt")
    completion: str = Field(..., description="Completion response")
    metadata: Optional[Metadata] = Field(None, description="Metadata")
    trajectory: Optional[List[TrajectoryStep]] = Field(None, description="Solution trajectory")


class ValidationResult(BaseModel):
    """Validation result model."""
    schema_name: str = Field(..., description="Schema name used for validation")
    total_examples: int = Field(..., description="Total number of examples")
    valid_examples: int = Field(..., description="Number of valid examples")
    invalid_examples: int = Field(..., description="Number of invalid examples")
    valid_ratio: float = Field(..., description="Ratio of valid examples")
    validation_errors: List[Dict[str, Any]] = Field(default_factory=list, description="Validation errors")


class QualityMetrics(BaseModel):
    """Data quality metrics model."""
    total_examples: int = Field(..., description="Total number of examples")
    empty_fields: Dict[str, int] = Field(..., description="Count of empty fields")
    length_metrics: Dict[str, Dict[str, float]] = Field(..., description="Length metrics for fields")
    source_distribution: Dict[str, int] = Field(..., description="Distribution of sources")
    repository_distribution: Dict[str, int] = Field(..., description="Distribution of repositories")
    topic_distribution: Dict[str, int] = Field(..., description="Distribution of topics")
    label_distribution: Dict[str, int] = Field(..., description="Distribution of labels")
