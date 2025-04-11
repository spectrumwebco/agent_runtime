"""
Connector module for integrating GitHub and Gitee issue scrapers with ML pipeline.

This module provides functionality to connect the issue scrapers with the ML pipeline,
enabling data collection from GitHub and Gitee repositories focused on GitOps, Terraform,
and Kubernetes for fine-tuning Llama 4 models.
"""

import os
import json
import logging
import asyncio
from typing import Dict, List, Optional, Union, Any
from datetime import datetime

from integrations.github.scraper import GitHubIssueScraper
from integrations.gitee.scraper import GiteeIssueScraper
from integrations.issue_collector.collector import IssueCollector

import mlflow
from mlflow.tracking import MlflowClient

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


class ScraperConnector:
    """
    Connector class for integrating GitHub and Gitee issue scrapers with ML pipeline.
    """

    def __init__(
        self,
        github_token: Optional[str] = None,
        gitee_token: Optional[str] = None,
        mlflow_tracking_uri: Optional[str] = None,
        output_dir: str = "data",
    ):
        """
        Initialize the ScraperConnector.

        Args:
            github_token: GitHub API token for authentication
            gitee_token: Gitee API token for authentication
            mlflow_tracking_uri: MLflow tracking URI for dataset versioning
            output_dir: Directory to store collected data
        """
        self.github_token = github_token or os.environ.get("GITHUB_TOKEN")
        self.gitee_token = gitee_token or os.environ.get("GITEE_TOKEN")
        self.output_dir = output_dir
        
        self.github_scraper = GitHubIssueScraper(token=self.github_token)
        self.gitee_scraper = GiteeIssueScraper(token=self.gitee_token)
        self.issue_collector = IssueCollector()
        
        if mlflow_tracking_uri:
            mlflow.set_tracking_uri(mlflow_tracking_uri)
        self.mlflow_client = MlflowClient()
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        logger.info("ScraperConnector initialized")

    async def collect_github_issues(
        self,
        repositories: List[str],
        topics: List[str] = ["gitops", "terraform", "kubernetes"],
        state: str = "closed",
        limit: int = 100,
    ) -> List[Dict[str, Any]]:
        """
        Collect issues from GitHub repositories.

        Args:
            repositories: List of GitHub repositories to collect issues from
            topics: List of topics to filter repositories by
            state: State of issues to collect (open, closed, all)
            limit: Maximum number of issues to collect per repository

        Returns:
            List of collected issues
        """
        logger.info(f"Collecting GitHub issues from {len(repositories)} repositories")
        
        all_issues = []
        for repo in repositories:
            try:
                issues = await self.github_scraper.get_issues(
                    repo=repo, state=state, limit=limit
                )
                
                filtered_issues = []
                for issue in issues:
                    repo_topics = await self.github_scraper.get_repository_topics(repo)
                    if any(topic in repo_topics for topic in topics):
                        filtered_issues.append(issue)
                
                all_issues.extend(filtered_issues)
                logger.info(f"Collected {len(filtered_issues)} issues from {repo}")
            except Exception as e:
                logger.error(f"Error collecting issues from {repo}: {e}")
        
        return all_issues

    async def collect_gitee_issues(
        self,
        repositories: List[str],
        topics: List[str] = ["gitops", "terraform", "kubernetes"],
        state: str = "closed",
        limit: int = 100,
    ) -> List[Dict[str, Any]]:
        """
        Collect issues from Gitee repositories.

        Args:
            repositories: List of Gitee repositories to collect issues from
            topics: List of topics to filter repositories by
            state: State of issues to collect (open, closed, all)
            limit: Maximum number of issues to collect per repository

        Returns:
            List of collected issues
        """
        logger.info(f"Collecting Gitee issues from {len(repositories)} repositories")
        
        all_issues = []
        for repo in repositories:
            try:
                issues = await self.gitee_scraper.get_issues(
                    repo=repo, state=state, limit=limit
                )
                
                filtered_issues = []
                for issue in issues:
                    repo_topics = await self.gitee_scraper.get_repository_topics(repo)
                    if any(topic in repo_topics for topic in topics):
                        filtered_issues.append(issue)
                
                all_issues.extend(filtered_issues)
                logger.info(f"Collected {len(filtered_issues)} issues from {repo}")
            except Exception as e:
                logger.error(f"Error collecting issues from {repo}: {e}")
        
        return all_issues

    async def collect_issues(
        self,
        github_repositories: List[str] = [],
        gitee_repositories: List[str] = [],
        topics: List[str] = ["gitops", "terraform", "kubernetes"],
        state: str = "closed",
        limit: int = 100,
    ) -> Dict[str, List[Dict[str, Any]]]:
        """
        Collect issues from both GitHub and Gitee repositories.

        Args:
            github_repositories: List of GitHub repositories to collect issues from
            gitee_repositories: List of Gitee repositories to collect issues from
            topics: List of topics to filter repositories by
            state: State of issues to collect (open, closed, all)
            limit: Maximum number of issues to collect per repository

        Returns:
            Dictionary containing collected issues from GitHub and Gitee
        """
        logger.info("Collecting issues from GitHub and Gitee repositories")
        
        github_task = self.collect_github_issues(
            repositories=github_repositories, topics=topics, state=state, limit=limit
        )
        gitee_task = self.collect_gitee_issues(
            repositories=gitee_repositories, topics=topics, state=state, limit=limit
        )
        
        github_issues, gitee_issues = await asyncio.gather(github_task, gitee_task)
        
        return {
            "github": github_issues,
            "gitee": gitee_issues,
        }

    def format_issues_for_training(
        self, issues: Dict[str, List[Dict[str, Any]]]
    ) -> List[Dict[str, Any]]:
        """
        Format collected issues for training.

        Args:
            issues: Dictionary containing collected issues from GitHub and Gitee

        Returns:
            List of formatted issues for training
        """
        logger.info("Formatting issues for training")
        
        formatted_issues = []
        
        for issue in issues.get("github", []):
            try:
                formatted_issue = self._format_issue_for_training(issue, source="github")
                if formatted_issue:
                    formatted_issues.append(formatted_issue)
            except Exception as e:
                logger.error(f"Error formatting GitHub issue: {e}")
        
        for issue in issues.get("gitee", []):
            try:
                formatted_issue = self._format_issue_for_training(issue, source="gitee")
                if formatted_issue:
                    formatted_issues.append(formatted_issue)
            except Exception as e:
                logger.error(f"Error formatting Gitee issue: {e}")
        
        logger.info(f"Formatted {len(formatted_issues)} issues for training")
        return formatted_issues

    def _format_issue_for_training(
        self, issue: Dict[str, Any], source: str
    ) -> Optional[Dict[str, Any]]:
        """
        Format a single issue for training.

        Args:
            issue: Issue to format
            source: Source of the issue (github or gitee)

        Returns:
            Formatted issue for training or None if issue cannot be formatted
        """
        if not issue.get("solution") and not issue.get("body"):
            return None
        
        solution = issue.get("solution", "")
        if not solution and issue.get("body"):
            solution = issue.get("body", "")
        
        repository = issue.get("repository", {})
        repo_name = repository.get("full_name", "")
        repo_topics = repository.get("topics", [])
        
        issue_title = issue.get("title", "")
        issue_description = issue.get("body", "")
        issue_id = issue.get("id", "")
        issue_url = issue.get("html_url", "")
        issue_created_at = issue.get("created_at", "")
        issue_closed_at = issue.get("closed_at", "")
        issue_labels = [label.get("name", "") for label in issue.get("labels", [])]
        
        trajectory = []
        for comment in issue.get("comments", []):
            trajectory.append({
                "step": len(trajectory) + 1,
                "action": "comment",
                "content": comment.get("body", ""),
                "timestamp": comment.get("created_at", ""),
                "user": comment.get("user", {}).get("login", ""),
            })
        
        formatted_issue = {
            "input": {
                "repository": repo_name,
                "topics": repo_topics,
                "title": issue_title,
                "description": issue_description,
            },
            "output": {
                "solution": solution,
            },
            "metadata": {
                "id": issue_id,
                "source": source,
                "repository": repo_name,
                "url": issue_url,
                "created_at": issue_created_at,
                "closed_at": issue_closed_at,
                "labels": issue_labels,
            },
            "trajectory": trajectory,
        }
        
        return formatted_issue

    def save_training_data(
        self, training_data: List[Dict[str, Any]], filename: str = "training_data.json"
    ) -> str:
        """
        Save training data to a file.

        Args:
            training_data: Training data to save
            filename: Name of the file to save training data to

        Returns:
            Path to the saved file
        """
        logger.info(f"Saving {len(training_data)} training examples to {filename}")
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        output_path = os.path.join(self.output_dir, filename)
        with open(output_path, "w") as f:
            json.dump(training_data, f, indent=2)
        
        logger.info(f"Training data saved to {output_path}")
        return output_path

    def log_dataset_to_mlflow(
        self,
        dataset_path: str,
        experiment_name: str = "llama4-fine-tuning",
        dataset_name: str = "github-gitee-issues",
    ) -> str:
        """
        Log dataset to MLflow for versioning.

        Args:
            dataset_path: Path to the dataset file
            experiment_name: Name of the MLflow experiment
            dataset_name: Name of the dataset

        Returns:
            MLflow run ID
        """
        logger.info(f"Logging dataset to MLflow: {dataset_path}")
        
        experiment = self.mlflow_client.get_experiment_by_name(experiment_name)
        if experiment is None:
            experiment_id = self.mlflow_client.create_experiment(experiment_name)
        else:
            experiment_id = experiment.experiment_id
        
        with mlflow.start_run(experiment_id=experiment_id) as run:
            run_id = run.info.run_id
            
            mlflow.log_artifact(dataset_path, "datasets")
            
            with open(dataset_path, "r") as f:
                dataset = json.load(f)
                
                mlflow.log_param("dataset_name", dataset_name)
                mlflow.log_param("dataset_size", len(dataset))
                mlflow.log_param("dataset_path", dataset_path)
                mlflow.log_param("dataset_created_at", datetime.now().isoformat())
                
                github_count = sum(1 for item in dataset if item["metadata"]["source"] == "github")
                gitee_count = sum(1 for item in dataset if item["metadata"]["source"] == "gitee")
                
                mlflow.log_metric("github_issues_count", github_count)
                mlflow.log_metric("gitee_issues_count", gitee_count)
                mlflow.log_metric("total_issues_count", len(dataset))
        
        logger.info(f"Dataset logged to MLflow: {run_id}")
        return run_id


async def main():
    """
    Main function for testing the ScraperConnector.
    """
    connector = ScraperConnector(
        mlflow_tracking_uri="http://mlflow-server.mlflow.svc.cluster.local:5000",
        output_dir="data",
    )
    
    github_repositories = [
        "kubernetes/kubernetes",
        "hashicorp/terraform",
        "fluxcd/flux2",
        "argoproj/argo-cd",
        "jenkins-x/jx",
    ]
    
    gitee_repositories = [
        "openharmony/kernel_linux_5.10",
        "openeuler/infrastructure",
        "open-cluster-management/ocm",
    ]
    
    issues = await connector.collect_issues(
        github_repositories=github_repositories,
        gitee_repositories=gitee_repositories,
        topics=["gitops", "terraform", "kubernetes"],
        state="closed",
        limit=100,
    )
    
    training_data = connector.format_issues_for_training(issues)
    
    dataset_path = connector.save_training_data(
        training_data, filename="github_gitee_issues.json"
    )
    
    run_id = connector.log_dataset_to_mlflow(
        dataset_path=dataset_path,
        experiment_name="llama4-fine-tuning",
        dataset_name="github-gitee-issues",
    )
    
    logger.info(f"Dataset logged to MLflow: {run_id}")


if __name__ == "__main__":
    asyncio.run(main())
