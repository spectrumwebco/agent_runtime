"""
Issue collector for D4E Agent.
"""

import os
import json
import asyncio
from typing import Any, Dict, List, Optional, Set, Tuple
from datetime import datetime
import logging

from ..github.integration import GitHubIntegration
from ..github.scraper import GitHubScraper
from ..gitee.integration import GiteeIntegration
from ..gitee.scraper import GiteeScraper


class IssueCollector:
    """Issue collector for D4E Agent."""

    def __init__(
        self,
        github_api_key: str,
        gitee_api_key: str,
        output_dir: str = "./data/collected_issues",
        log_level: int = logging.INFO,
    ):
        """
        Initialize the issue collector.

        Args:
            github_api_key: GitHub API key
            gitee_api_key: Gitee API key
            output_dir: Directory to save collected data
            log_level: Logging level
        """
        self.github_api_key = github_api_key
        self.gitee_api_key = gitee_api_key
        self.output_dir = output_dir

        os.makedirs(output_dir, exist_ok=True)

        logging.basicConfig(
            level=log_level,
            format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        )
        self.logger = logging.getLogger("IssueCollector")

        self.github_integration = GitHubIntegration(github_api_key)
        self.gitee_integration = GiteeIntegration(gitee_api_key)

        self.github_scraper = GitHubScraper(
            self.github_integration,
            output_dir=os.path.join(output_dir, "github"),
        )
        self.gitee_scraper = GiteeScraper(
            self.gitee_integration,
            output_dir=os.path.join(output_dir, "gitee"),
        )

    async def collect_issues(
        self,
        topics: List[str],
        languages: Optional[List[str]] = None,
        min_stars: int = 100,
        max_repos_per_platform: int = 25,
        max_issues_per_repo: int = 50,
        include_pull_requests: bool = False,
    ) -> Dict[str, Any]:
        """
        Collect issues from GitHub and Gitee.

        Args:
            topics: List of topics to search for
            languages: Optional list of languages to filter by
            min_stars: Minimum number of stars
            max_repos_per_platform: Maximum number of repositories to scrape per platform
            max_issues_per_repo: Maximum number of issues to scrape per repository
            include_pull_requests: Whether to include pull requests

        Returns:
            Collection results
        """
        self.logger.info(f"Collecting issues for topics: {topics}")

        self.logger.info("Collecting issues from GitHub")
        github_issues_path, github_training_data_path = (
            await self.github_scraper.scrape_and_save(
                topics=topics,
                languages=languages,
                min_stars=min_stars,
                max_repos=max_repos_per_platform,
                max_issues_per_repo=max_issues_per_repo,
                include_pull_requests=include_pull_requests,
            )
        )

        self.logger.info("Collecting issues from Gitee")
        gitee_issues_path, gitee_training_data_path = (
            await self.gitee_scraper.scrape_and_save(
                topics=topics,
                languages=languages,
                min_stars=min_stars,
                max_repos=max_repos_per_platform,
                max_issues_per_repo=max_issues_per_repo,
                include_pull_requests=include_pull_requests,
            )
        )

        combined_training_data = await self.combine_training_data(
            github_training_data_path,
            gitee_training_data_path,
        )

        return {
            "github_issues_path": github_issues_path,
            "github_training_data_path": github_training_data_path,
            "gitee_issues_path": gitee_issues_path,
            "gitee_training_data_path": gitee_training_data_path,
            "combined_training_data_path": combined_training_data,
        }

    async def combine_training_data(
        self, github_training_data_path: str, gitee_training_data_path: str
    ) -> str:
        """
        Combine training data from GitHub and Gitee.

        Args:
            github_training_data_path: Path to GitHub training data
            gitee_training_data_path: Path to Gitee training data

        Returns:
            Path to combined training data
        """
        self.logger.info("Combining training data from GitHub and Gitee")

        with open(github_training_data_path, "r") as f:
            github_training_data = json.load(f)

        with open(gitee_training_data_path, "r") as f:
            gitee_training_data = json.load(f)

        combined_training_data = github_training_data + gitee_training_data

        output_path = os.path.join(self.output_dir, "combined_training_data.json")
        with open(output_path, "w") as f:
            json.dump(combined_training_data, f, indent=2)

        self.logger.info(
            f"Saved {len(combined_training_data)} combined training examples to {output_path}"
        )
        return output_path

    async def collect_and_save(
        self,
        topics: List[str] = ["gitops", "terraform", "kubernetes", "k8s"],
        languages: Optional[List[str]] = None,
        min_stars: int = 100,
        max_repos_per_platform: int = 25,
        max_issues_per_repo: int = 50,
        include_pull_requests: bool = False,
    ) -> Dict[str, Any]:
        """
        Collect and save issues from GitHub and Gitee.

        Args:
            topics: List of topics to search for
            languages: Optional list of languages to filter by
            min_stars: Minimum number of stars
            max_repos_per_platform: Maximum number of repositories to scrape per platform
            max_issues_per_repo: Maximum number of issues to scrape per repository
            include_pull_requests: Whether to include pull requests

        Returns:
            Collection results
        """
        return await self.collect_issues(
            topics=topics,
            languages=languages,
            min_stars=min_stars,
            max_repos_per_platform=max_repos_per_platform,
            max_issues_per_repo=max_issues_per_repo,
            include_pull_requests=include_pull_requests,
        )
