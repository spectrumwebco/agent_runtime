"""
Pydantic models for GitHub data structures.
"""

from typing import List, Optional
from datetime import datetime
from pydantic import BaseModel, Field, HttpUrl


class GitHubRepository(BaseModel):
    """GitHub repository model."""
    id: int = Field(..., description="Repository ID")
    name: str = Field(..., description="Repository name")
    full_name: str = Field(..., description="Full repository name (owner/repo)")
    owner: str = Field(..., description="Repository owner")
    url: HttpUrl = Field(..., description="Repository URL")
    html_url: HttpUrl = Field(..., description="Repository HTML URL")
    description: Optional[str] = Field(None, description="Repository description")
    stargazers_count: int = Field(..., description="Number of stars")
    topics: List[str] = Field(default_factory=list, description="Repository topics")
    language: Optional[str] = Field(None, description="Primary repository language")
    created_at: datetime = Field(..., description="Repository creation date")
    updated_at: datetime = Field(..., description="Repository last update date")


class GitHubLabel(BaseModel):
    """GitHub issue label model."""
    id: int = Field(..., description="Label ID")
    name: str = Field(..., description="Label name")
    description: Optional[str] = Field(None, description="Label description")
    color: str = Field(..., description="Label color")


class GitHubIssue(BaseModel):
    """GitHub issue model."""
    id: int = Field(..., description="Issue ID")
    number: int = Field(..., description="Issue number")
    title: str = Field(..., description="Issue title")
    body: Optional[str] = Field(None, description="Issue body")
    state: str = Field(..., description="Issue state")
    html_url: HttpUrl = Field(..., description="Issue HTML URL")
    created_at: datetime = Field(..., description="Issue creation date")
    updated_at: datetime = Field(..., description="Issue last update date")
    closed_at: Optional[datetime] = Field(None, description="Issue closure date")
    labels: List[GitHubLabel] = Field(default_factory=list, description="Issue labels")
    repository: GitHubRepository = Field(..., description="Parent repository")


class TrainingExample(BaseModel):
    """Training example model."""
    input: str = Field(..., description="Input text for training")
    output: str = Field(..., description="Expected output text")
    metadata: dict = Field(..., description="Additional metadata")
    trajectory: List[dict] = Field(..., description="Solution trajectory")
