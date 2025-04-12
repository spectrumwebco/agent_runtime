"""
Integration test for Gitee scraper.
"""

import os
import pytest
import asyncio
from src.integrations.gitee.integration import GiteeIntegration
from src.integrations.gitee.scraper import GiteeScraper


@pytest.mark.asyncio
async def test_gitee_scraper_initialization():
    """Test Gitee scraper initialization."""
    api_key = "test_api_key"
    
    integration = GiteeIntegration(api_key=api_key)
    
    scraper = GiteeScraper(integration=integration)
    
    assert scraper.integration == integration
    assert hasattr(scraper, "logger")


@pytest.mark.asyncio
async def test_gitee_scraper_search_repositories():
    """Test Gitee scraper search repositories method."""
    gitee_token = os.environ.get("GITEE_TOKEN")
    if not gitee_token:
        pytest.skip("No Gitee token available")
    
    integration = GiteeIntegration(api_key=gitee_token)
    
    scraper = GiteeScraper(integration=integration)
    
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
async def test_gitee_scraper_get_issues():
    """Test Gitee scraper get issues method."""
    gitee_token = os.environ.get("GITEE_TOKEN")
    if not gitee_token:
        pytest.skip("No Gitee token available")
    
    integration = GiteeIntegration(api_key=gitee_token)
    
    scraper = GiteeScraper(integration=integration)
    
    issues = await scraper.get_issues(
        owner="gitee",
        repo="gitee",
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
    asyncio.run(test_gitee_scraper_initialization())
