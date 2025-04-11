"""
Gitee integration for D4E Agent.
"""

import os
import json
import asyncio
import logging
from typing import Dict, List, Any, Optional
from datetime import datetime
import aiohttp


class GiteeIntegration:
    """Gitee API integration for D4E Agent."""

    def __init__(self, api_key: str, base_url: str = "https://gitee.com/api/v5"):
        """
        Initialize the Gitee integration.

        Args:
            api_key: Gitee API key
            base_url: Gitee API base URL
        """
        self.api_key = api_key
        self.base_url = base_url
        self.headers = {
            "Content-Type": "application/json;charset=UTF-8",
        }
        self.logger = logging.getLogger("GiteeIntegration")

    async def _make_request(
        self, method: str, endpoint: str, params: Optional[Dict[str, Any]] = None, data: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """
        Make a request to the Gitee API.

        Args:
            method: HTTP method
            endpoint: API endpoint
            params: Query parameters
            data: Request data

        Returns:
            Response data
        """
        url = f"{self.base_url}/{endpoint}"
        
        if params is None:
            params = {}
        
        params["access_token"] = self.api_key
        
        async with aiohttp.ClientSession() as session:
            async with session.request(
                method=method,
                url=url,
                headers=self.headers,
                params=params,
                json=data,
            ) as response:
                try:
                    response_data = await response.json()
                except Exception as e:
                    self.logger.error(f"Error parsing response: {str(e)}")
                    response_data = {"error": str(e)}
                
                if response.status >= 400:
                    self.logger.error(f"Gitee API error: {response.status} - {response_data}")
                    return {"error": response_data, "status_code": response.status}
                
                return response_data

    async def search_repositories(
        self,
        query: str,
        page: int = 1,
        per_page: int = 30,
        order: str = "desc",
    ) -> Dict[str, Any]:
        """
        Search for repositories.

        Args:
            query: Search query
            page: Page number
            per_page: Results per page
            order: Sort order

        Returns:
            Search results
        """
        params = {
            "q": query,
            "page": page,
            "per_page": per_page,
            "order": order,
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
        page: int = 1,
        per_page: int = 30,
    ) -> List[Dict[str, Any]]:
        """
        Get repository issues.

        Args:
            owner: Repository owner
            repo: Repository name
            state: Issue state (open, closed, all)
            sort: Sort field
            direction: Sort direction
            page: Page number
            per_page: Results per page

        Returns:
            Repository issues
        """
        params = {
            "state": state,
            "sort": sort,
            "direction": direction,
            "page": page,
            "per_page": per_page,
        }
        
        return await self._make_request("GET", f"repos/{owner}/{repo}/issues", params=params)

    async def get_issue(self, owner: str, repo: str, issue_number: int) -> Dict[str, Any]:
        """
        Get issue details.

        Args:
            owner: Repository owner
            repo: Repository name
            issue_number: Issue number

        Returns:
            Issue details
        """
        return await self._make_request("GET", f"repos/{owner}/{repo}/issues/{issue_number}")

    async def get_issue_comments(
        self, owner: str, repo: str, issue_number: int, page: int = 1, per_page: int = 30
    ) -> List[Dict[str, Any]]:
        """
        Get issue comments.

        Args:
            owner: Repository owner
            repo: Repository name
            issue_number: Issue number
            page: Page number
            per_page: Results per page

        Returns:
            Issue comments
        """
        params = {
            "page": page,
            "per_page": per_page,
        }
        
        return await self._make_request(
            "GET", f"repos/{owner}/{repo}/issues/{issue_number}/comments", params=params
        )

    async def get_rate_limit(self) -> Dict[str, Any]:
        """
        Get rate limit status.

        Note: Gitee does not have a direct endpoint for rate limit status.
        This is a mock implementation.

        Returns:
            Rate limit status
        """
        return {
            "resources": {
                "core": {
                    "limit": 5000,
                    "remaining": 4500,
                    "reset": int(datetime.now().timestamp()) + 3600,
                }
            }
        }
