"""
Pydantic models for Gitee integration.

This module provides Pydantic models for Gitee API responses and data structures
used in the Gitee issue scraper.
"""

from typing import Dict, List, Optional, Any
from datetime import datetime
from pydantic import BaseModel, Field, HttpUrl


class GiteeUser(BaseModel):
    """Gitee user model."""
    id: int = Field(..., description="User ID")
    login: str = Field(..., description="Username")
    name: Optional[str] = Field(None, description="User's real name")
    html_url: HttpUrl = Field(..., description="User profile URL")
    type: str = Field(..., description="User type")
    site_admin: bool = Field(False, description="Whether the user is a site admin")


class GiteeLabel(BaseModel):
    """Gitee label model."""
    id: int = Field(..., description="Label ID")
    name: str = Field(..., description="Label name")
    color: str = Field(..., description="Label color")
    description: Optional[str] = Field(None, description="Label description")


class GiteeRepository(BaseModel):
    """Gitee repository model."""
    id: int = Field(..., description="Repository ID")
    name: str = Field(..., description="Repository name")
    full_name: str = Field(..., description="Full repository name")
    description: Optional[str] = Field(None, description="Repository description")
    html_url: HttpUrl = Field(..., description="Repository URL")
    stars: int = Field(0, description="Number of stars", alias="stargazers_count")
    forks: int = Field(0, description="Number of forks", alias="forks_count")
    topics: List[str] = Field(default_factory=list, description="Repository topics")
    language: Optional[str] = Field(None, description="Primary language")
    created_at: Optional[datetime] = Field(None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(None, description="Last update timestamp")
    owner: Optional[GiteeUser] = Field(None, description="Repository owner")


class GiteeIssue(BaseModel):
    """Gitee issue model."""
    id: int = Field(..., description="Issue ID")
    number: int = Field(..., description="Issue number")
    title: str = Field(..., description="Issue title")
    body: Optional[str] = Field(None, description="Issue body")
    html_url: HttpUrl = Field(..., description="Issue URL")
    state: str = Field(..., description="Issue state")
    created_at: datetime = Field(..., description="Creation timestamp")
    updated_at: datetime = Field(..., description="Last update timestamp")
    closed_at: Optional[datetime] = Field(None, description="Closure timestamp")
    user: GiteeUser = Field(..., description="Issue creator")
    labels: List[GiteeLabel] = Field(default_factory=list, description="Issue labels")
    comments: int = Field(0, description="Number of comments")
    pull_request: Optional[Dict[str, Any]] = Field(None, description="Pull request data if this is a PR")
    repository: Optional[GiteeRepository] = Field(None, description="Repository this issue belongs to")


class GiteeComment(BaseModel):
    """Gitee comment model."""
    id: int = Field(..., description="Comment ID")
    body: str = Field(..., description="Comment body")
    html_url: HttpUrl = Field(..., description="Comment URL")
    created_at: datetime = Field(..., description="Creation timestamp")
    updated_at: datetime = Field(..., description="Last update timestamp")
    user: GiteeUser = Field(..., description="Comment author")


class GiteeSearchResult(BaseModel):
    """Gitee search result model."""
    total_count: int = Field(..., description="Total number of results")
    incomplete_results: bool = Field(False, description="Whether results are incomplete")
    items: List[GiteeIssue] = Field(..., description="Search results")


class GiteeScraperConfig(BaseModel):
    """Gitee scraper configuration model."""
    topics: List[str] = Field(
        default_factory=lambda: ["gitops", "terraform", "kubernetes", "k8s"],
        description="List of topics to search for"
    )
    languages: Optional[List[str]] = Field(None, description="Optional list of languages to filter by")
    min_stars: int = Field(100, description="Minimum number of stars", ge=0)
    max_repos: int = Field(25, description="Maximum number of repositories to scrape", gt=0)
    max_issues_per_repo: int = Field(50, description="Maximum number of issues to scrape per repository", gt=0)
    include_pull_requests: bool = Field(False, description="Whether to include pull requests")
    token: Optional[str] = Field(None, description="Gitee API token")
