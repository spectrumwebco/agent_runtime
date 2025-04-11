"""
GitHub integration for D4E Agent.
"""

import os
import json
import asyncio
import logging
from typing import Dict, List, Any, Optional
import aiohttp


class GitHubIntegration:
    """GitHub API integration for D4E Agent."""

    def __init__(self, api_key: str, base_url: str = "https://api.github.com"):
        """
        Initialize the GitHub integration.

        Args:
            api_key: GitHub API key
            base_url: GitHub API base URL
        """
        self.api_key = api_key
        self.base_url = base_url
        self.headers = {
            "Authorization": f"token {api_key}",
            "Accept": "application/vnd.github.v3+json",
        }
        self.logger = logging.getLogger("GitHubIntegration")

    async def _make_request(
        self,
        method: str,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        data: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Make a request to the GitHub API.

        Args:
            method: HTTP method
            endpoint: API endpoint
            params: Query parameters
            data: Request data

        Returns:
            Response data
        """
        url = f"{self.base_url}/{endpoint}"

        async with aiohttp.ClientSession() as session:
            async with session.request(
                method=method,
                url=url,
                headers=self.headers,
                params=params,
                json=data,
            ) as response:
                response_data = await response.json()

                if response.status >= 400:
                    self.logger.error(
                        f"GitHub API error: {response.status} - {response_data}"
                    )
                    return {"error": response_data, "status_code": response.status}

                return response_data

    async def search_repositories(
        self,
        query: str,
        sort: str = "stars",
        order: str = "desc",
        per_page: int = 30,
        page: int = 1,
    ) -> Dict[str, Any]:
        """
        Search for repositories.

        Args:
            query: Search query
            sort: Sort field
            order: Sort order
            per_page: Results per page
            page: Page number

        Returns:
            Search results
        """
        params = {
            "q": query,
            "sort": sort,
            "order": order,
            "per_page": per_page,
            "page": page,
        }

        return await self._make_request("GET", "search/repositories", params=params)

    async def get_repository(self, owner: str, repo: str) -> Dict[str, Any]:
        """
        Get repository details.

        Args:
            owner: Repository owner
            repo: Repository name

        Returns:
            Repository details
        """
        return await self._make_request("GET", f"repos/{owner}/{repo}")

    async def get_issues(
        self,
        owner: str,
        repo: str,
        state: str = "all",
        sort: str = "created",
        direction: str = "desc",
        per_page: int = 30,
        page: int = 1,
    ) -> List[Dict[str, Any]]:
        """
        Get repository issues.

        Args:
            owner: Repository owner
            repo: Repository name
            state: Issue state (open, closed, all)
            sort: Sort field
            direction: Sort direction
            per_page: Results per page
            page: Page number

        Returns:
            Repository issues
        """
        params = {
            "state": state,
            "sort": sort,
            "direction": direction,
            "per_page": per_page,
            "page": page,
        }

        return await self._make_request(
            "GET", f"repos/{owner}/{repo}/issues", params=params
        )

    async def get_issue(
        self, owner: str, repo: str, issue_number: int
    ) -> Dict[str, Any]:
        """
        Get issue details.

        Args:
            owner: Repository owner
            repo: Repository name
            issue_number: Issue number

        Returns:
            Issue details
        """
        return await self._make_request(
            "GET", f"repos/{owner}/{repo}/issues/{issue_number}"
        )

    async def get_issue_comments(
        self,
        owner: str,
        repo: str,
        issue_number: int,
        per_page: int = 30,
        page: int = 1,
    ) -> List[Dict[str, Any]]:
        """
        Get issue comments.

        Args:
            owner: Repository owner
            repo: Repository name
            issue_number: Issue number
            per_page: Results per page
            page: Page number

        Returns:
            Issue comments
        """
        params = {
            "per_page": per_page,
            "page": page,
        }

        return await self._make_request(
            "GET", f"repos/{owner}/{repo}/issues/{issue_number}/comments", params=params
        )

    async def get_rate_limit(self) -> Dict[str, Any]:
        """
        Get rate limit status.

        Returns:
            Rate limit status
        """
        return await self._make_request("GET", "rate_limit")
