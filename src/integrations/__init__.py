"""
Integrations module for D4E Agent.
"""

from .github import GitHubIntegration, GitHubScraper
from .gitee import GiteeIntegration, GiteeScraper
from .issue_collector import IssueCollector

__all__ = [
    "GitHubIntegration",
    "GitHubScraper",
    "GiteeIntegration",
    "GiteeScraper",
    "IssueCollector",
]
