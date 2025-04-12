"""
Gitee scraper for D4E Agent.
"""

import os
import json
import asyncio
import logging
from typing import Dict, List, Any, Optional, Set, Tuple
from datetime import datetime

from .integration import GiteeIntegration
from .models import GiteeRepository, GiteeIssue, TrainingExample


class GiteeScraper:
    """Gitee scraper for D4E Agent."""

    def __init__(
        self,
        gitee_integration: GiteeIntegration,
        output_dir: str = "./data/gitee",
    ):
        """
        Initialize the Gitee scraper.

        Args:
            gitee_integration: Gitee API integration
            output_dir: Directory to save scraped data
        """
        self.gitee_integration = gitee_integration
        self.output_dir = output_dir

        os.makedirs(output_dir, exist_ok=True)

        self.logger = logging.getLogger("GiteeScraper")

    async def scrape_repositories(
        self,
        topics: List[str],
        languages: Optional[List[str]] = None,
        min_stars: int = 100,
        max_repos: int = 25,
    ) -> List[Dict[str, Any]]:
        """
        Scrape repositories based on topics and languages.

        Args:
            topics: List of topics to search for
            languages: Optional list of languages to filter by
            min_stars: Minimum number of stars
            max_repos: Maximum number of repositories to return

        Returns:
            List of repositories
        """
        self.logger.info(f"Scraping repositories for topics: {topics}")

        repositories = []

        for topic in topics:
            query = f"{topic} stars:>={min_stars}"

            if languages:
                for language in languages:
                    language_query = f"{query} language:{language}"

                    self.logger.info(
                        f"Searching for repositories with query: {language_query}"
                    )

                    response = await self.gitee_integration.search_repositories(
                                                                               query=language_query,
                                                                               page=1,
                                                                               per_page=min(max_repos,
                                                                               100
                                                                           ),
                        order="desc",
                    )

                    if isinstance(response, list):
                        repositories.extend(response)

                        self.logger.info(
                            f"Found {len(response)} repositories for {topic} in {language}"
                        )
            else:
                self.logger.info(f"Searching for repositories with query: {query}")

                response = await self.gitee_integration.search_repositories(
                    query=query,
                    page=1,
                    per_page=min(max_repos, 100),
                    order="desc",
                )

                if isinstance(response, list):
                    repositories.extend(response)

                    self.logger.info(f"Found {len(response)} repositories for {topic}")

        unique_repos = {}
        for repo in repositories:
            repo_id = repo["id"]
            if repo_id not in unique_repos:
                unique_repos[repo_id] = repo

        repositories = list(unique_repos.values())

        repositories = [
            repo
            for repo in repositories
            if repo.get("stargazers_count", 0) >= min_stars
        ]

        repositories.sort(key=
    lambda x: x.get("stargazers_count", 0), reverse=True)

        repositories = repositories[:max_repos]

        self.logger.info(f"Scraped {len(repositories)} unique repositories")

        return repositories

    async def scrape_issues(
        self,
        repositories: List[Dict[str, Any]],
        state: str = "closed",
        max_issues_per_repo: int = 50,
        include_pull_requests: bool = False,
    ) -> List[Dict[str, Any]]:
        """
        Scrape issues from repositories.

        Args:
            repositories: List of repositories
            state: Issue state (open, closed, all)
            max_issues_per_repo: Maximum number of issues to scrape per repository
            include_pull_requests: Whether to include pull requests

        Returns:
            List of issues
        """
        self.logger.info(f"Scraping issues from {len(repositories)} repositories")

        all_issues = []

        for repo in repositories:
            owner = repo.get(
                            "owner",
                            {}).get("login") or repo.get("namespace",
                            {}).get(                "path"
                        )
            name = repo.get("name") or repo.get("path")

            if not owner or not name:
                self.logger.warning(
                    f"Skipping repository with missing owner or name: {repo}"
                )
                continue

            self.logger.info(f"Scraping issues from {owner}/{name}")

            try:
                issues = await self.gitee_integration.get_issues(
                    owner=owner,
                    repo=name,
                    state=state,
                    sort="created",
                    direction="desc",
                    page=1,
                    per_page=min(max_issues_per_repo, 100),
                )

                if not include_pull_requests:
                    issues = [
                        issue for issue in issues if not issue.get("pull_request")
                    ]

                for issue in issues:
                    issue["repository"] = repo

                all_issues.extend(issues)

                self.logger.info(f"Scraped {len(issues)} issues from {owner}/{name}")

                rate_limit = await self.gitee_integration.get_rate_limit()

                if "resources" in rate_limit and "core" in rate_limit["resources"]:
                    remaining = rate_limit["resources"]["core"]["remaining"]

                    if remaining < 10:
                        reset_time = rate_limit["resources"]["core"]["reset"]
                        reset_datetime = datetime.fromtimestamp(reset_time)

                        self.logger.warning(
                            f"Rate limit low: {remaining} requests remaining"
                        )
                        self.logger.warning(f"Rate limit resets at {reset_datetime}")

                        now = datetime.now()
                        sleep_time = (reset_datetime - now).total_seconds() + 10

                        if sleep_time > 0:
                            self.logger.warning(f"Sleeping for {sleep_time} seconds")
                            await asyncio.sleep(sleep_time)

            except Exception as e:
                self.logger.error(
                    f"Error scraping issues from {owner}/{name}: {str(e)}"
                )

        self.logger.info(f"Scraped {len(all_issues)} issues in total")

        return all_issues

    async def save_issues(
        self, issues: List[Dict[str, Any]], filename: str = "issues.json"
    ) -> str:
        """
        Save issues to a file.

        Args:
            issues: List of issues
            filename: Output filename

        Returns:
            Path to the saved file
        """
        output_path = os.path.join(self.output_dir, filename)

        with open(output_path, "w") as f:
            json.dump(issues, f, indent=2)

        self.logger.info(f"Saved {len(issues)} issues to {output_path}")

        return output_path

    async def format_for_training(
        self, issues: List[Dict[str, Any]], filename: str = "training_data.json"
    ) -> str:
        """
        Format issues for model training.

        Args:
            issues: List of issues
            filename: Output filename

        Returns:
            Path to the saved file
        """
        training_data = []

        for issue in issues:
            if "repository" not in issue:
                continue

            if not issue.get("body"):
                continue

            repo = issue["repository"]
            repo_name = (f"{repo.get('namespace', {}).get('path') or repo.get('owner', {}).get('login')}/"
                      f"{repo.get('path') or repo.get('name')}")
            repo_topics = repo.get("topics", [])

            issue_title = issue["title"]
            issue_body = issue["body"]
            issue_number = issue["number"]
            issue_url = issue["html_url"]
            issue_created_at = issue["created_at"]
            issue_closed_at = issue.get("closed_at")
            issue_labels = [label["name"] for label in issue.get("labels", [])]

            input_text = f"Repository: {repo_name}\n"

            if repo_topics:
                input_text += f"Topics: {', '.join(repo_topics)}\n"

            input_text += f"Issue Title: {issue_title}\n"
            input_text += f"Issue Description:\n{issue_body}\n"

            output_text = "Issue resolved successfully."

            metadata = {
                "issue_id": issue["id"],
                "issue_number": issue_number,
                "repository": repo_name,
                "url": issue_url,
                "created_at": issue_created_at,
                "closed_at": issue_closed_at,
                "labels": issue_labels,
            }

            trajectory = [
                {
                    "action": "read_issue",
                    "observation": f"Issue #{issue_number}: {issue_title}",
                    "response": "I'll analyze this issue to find a solution.",
                },
                {
                    "action": "analyze_issue",
                    "observation": issue_body,
                    "response": "Based on the issue description, I need to understand the problem and find a solution.",
                },
                {
                    "action": "implement_solution",
                    "observation": "Testing the solution...",
                    "response": output_text,
                },
                {
                    "action": "verify_solution",
                    "observation": "The solution has been implemented and tested.",
                    "response": "The issue has been resolved successfully.",
                },
            ]

            training_example = TrainingExample(
                input=input_text,
                output=output_text,
                metadata=metadata,
                trajectory=trajectory
            )
            
            training_data.append(training_example.dict())

        output_path = os.path.join(self.output_dir, filename)

        with open(output_path, "w") as f:
            json.dump(training_data, f, indent=2)

        self.logger.info(
            f"Saved {len(training_data)} training examples to {output_path}"
        )

        return output_path

    async def scrape_and_save(
        self,
        topics: List[str] = ["gitops", "terraform", "kubernetes", "k8s"],
        languages: Optional[List[str]] = None,
        min_stars: int = 100,
        max_repos: int = 25,
        max_issues_per_repo: int = 50,
        include_pull_requests: bool = False,
    ) -> Tuple[str, str]:
        """
        Scrape repositories and issues, and save them to files.

        Args:
            topics: List of topics to search for
            languages: Optional list of languages to filter by
            min_stars: Minimum number of stars
            max_repos: Maximum number of repositories to scrape
            max_issues_per_repo: Maximum number of issues to scrape per repository
            include_pull_requests: Whether to include pull requests

        Returns:
            Tuple of (issues_path, training_data_path)
        """
        repositories = await self.scrape_repositories(
            topics=topics,
            languages=languages,
            min_stars=min_stars,
            max_repos=max_repos,
        )

        issues = await self.scrape_issues(
            repositories=repositories,
            state="closed",
            max_issues_per_repo=max_issues_per_repo,
            include_pull_requests=include_pull_requests,
        )

        issues_path = await self.save_issues(issues)

        training_data_path = await self.format_for_training(issues)

        return issues_path, training_data_path
