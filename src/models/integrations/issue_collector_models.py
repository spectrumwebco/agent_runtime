"""
Pydantic models for issue collector.
"""

from typing import Dict, List, Optional, Any
from datetime import datetime
from pydantic import BaseModel, Field, HttpUrl, validator


class CollectionConfig(BaseModel):
    """Configuration for issue collection."""
    topics: List[str] = Field(
        default_factory=lambda: ["gitops", "terraform", "kubernetes", "k8s"],
        description="List of topics to search for"
    )
    languages: Optional[List[str]] = Field(
        None, 
        description="Optional list of languages to filter by"
    )
    min_stars: int = Field(
        100, 
        description="Minimum number of stars",
        ge=0
    )
    max_repos_per_platform: int = Field(
        25, 
        description="Maximum number of repositories to scrape per platform",
        gt=0
    )
    max_issues_per_repo: int = Field(
        50, 
        description="Maximum number of issues to scrape per repository",
        gt=0
    )
    include_pull_requests: bool = Field(
        False, 
        description="Whether to include pull requests"
    )


class CollectionResult(BaseModel):
    """Result of issue collection."""
    github_issues_path: str = Field(..., description="Path to GitHub issues")
    github_training_data_path: str = Field(..., description="Path to GitHub training data")
    gitee_issues_path: str = Field(..., description="Path to Gitee issues")
    gitee_training_data_path: str = Field(..., description="Path to Gitee training data")
    combined_training_data_path: str = Field(..., description="Path to combined training data")
    
    @validator('*')
    def check_file_paths(cls, v):
        """Validate that file paths exist."""
        if not v.endswith('.json'):
            raise ValueError(f"File path must end with .json: {v}")
        return v


class TrainingExample(BaseModel):
    """Training example model."""
    input: str = Field(..., description="Input text for training")
    output: str = Field(..., description="Expected output text")
    metadata: Dict[str, Any] = Field(..., description="Additional metadata")
    trajectory: List[Dict[str, str]] = Field(..., description="Solution trajectory")
    
    @validator('trajectory')
    def validate_trajectory(cls, v):
        """Validate trajectory structure."""
        required_keys = {"action", "observation", "response"}
        for step in v:
            if not all(key in step for key in required_keys):
                missing = required_keys - set(step.keys())
                raise ValueError(f"Trajectory step missing required keys: {missing}")
        return v
