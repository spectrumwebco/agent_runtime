"""
Integration test for GitHub scraper.
"""

import os
import pytest
import asyncio
from src.integrations.github.integration import GitHubIntegration
from src.integrations.github.scraper import GitHubScraper


@pytest.mark.asyncio
async def test_github_scraper_initialization():
    """Test GitHub scraper initialization."""
    api_key = "test_api_key"
    
    integration = GitHubIntegration(api_key=api_key)
    
    scraper = GitHubScraper(integration=integration)
    
    assert scraper.integration == integration
    assert hasattr(scraper, "logger")


@pytest.mark.asyncio
async def test_github_scraper_search_repositories():
    """Test GitHub scraper search repositories method."""
    github_token = os.environ.get("GITHUB_TOKEN")
    if not github_token:
        pytest.skip("No GitHub token available")
    
    integration = GitHubIntegration(api_key=github_token)
    
    scraper = GitHubScraper(integration=integration)
    
    repositories = await scraper.search_repositories(
        query="kubernetes", 
        limit=5
    )
    
    assert isinstance(repositories, list)
    assert len(repositories) <= 5
    
    if repositories:
        assert "name" in repositories[0]
        assert "owner" in repositories[0]
        assert "url" in repositories[0]


@pytest.mark.asyncio
async def test_github_scraper_get_issues():
    """Test GitHub scraper get issues method."""
    github_token = os.environ.get("GITHUB_TOKEN")
    if not github_token:
        pytest.skip("No GitHub token available")
    
    integration = GitHubIntegration(api_key=github_token)
    
    scraper = GitHubScraper(integration=integration)
    
    issues = await scraper.get_issues(
        owner="kubernetes",
        repo="kubernetes",
        state="closed",
        limit=5
    )
    
    assert isinstance(issues, list)
    assert len(issues) <= 5
    
    if issues:
        assert "number" in issues[0]
        assert "title" in issues[0]
        assert "state" in issues[0]
        assert issues[0]["state"] == "closed"


if __name__ == "__main__":
    asyncio.run(test_github_scraper_initialization())
