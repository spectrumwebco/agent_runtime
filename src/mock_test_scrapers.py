"""
Mock test script for GitHub and Gitee scrapers.
"""

import os
import json
import asyncio
import logging
from typing import Dict, List, Any, Optional
from datetime import datetime

from integrations.github.integration import GitHubIntegration
from integrations.github.scraper import GitHubScraper
from integrations.gitee.integration import GiteeIntegration
from integrations.gitee.scraper import GiteeScraper
from integrations.issue_collector.collector import IssueCollector


async def mock_test_github_scraper(
    output_dir: str = "./data/test/github",
) -> Dict[str, Any]:
    """
    Mock test for GitHub scraper.

    Args:
        output_dir: Output directory for test results

    Returns:
        Test results
    """
    logging.info("Mock testing GitHub scraper")

    os.makedirs(output_dir, exist_ok=True)

    repositories = [
        {
            "id": 1001,
            "full_name": "kubernetes/kubernetes",
            "html_url": "https://github.com/kubernetes/kubernetes",
            "description": "Production-Grade Container Scheduling and Management",
            "stargazers_count": 100000,
            "topics": ["kubernetes", "k8s", "container", "orchestration"],
            "language": "Go",
        },
        {
            "id": 1002,
            "full_name": "hashicorp/terraform",
            "html_url": "https://github.com/hashicorp/terraform",
            "description": "Terraform enables you to safely and predictably create, change, and improve infrastructure.",
            "stargazers_count": 50000,
            "topics": ["terraform", "infrastructure", "iac", "hashicorp"],
            "language": "Go",
        },
        {
            "id": 1003,
            "full_name": "fluxcd/flux",
            "html_url": "https://github.com/fluxcd/flux",
            "description": "GitOps Kubernetes operator",
            "stargazers_count": 10000,
            "topics": ["gitops", "kubernetes", "devops", "cicd"],
            "language": "Go",
        },
    ]

    issues = [
        {
            "id": 101,
            "number": 1001,
            "title": "Fix pod scheduling issue in multi-zone clusters",
            "body": "When deploying pods across multiple zones, the scheduler is not respecting zone anti-affinity rules.",
            "html_url": "https://github.com/kubernetes/kubernetes/issues/1001",
            "state": "closed",
            "created_at": "2023-01-15T10:30:00Z",
            "closed_at": "2023-01-20T14:45:00Z",
            "labels": [
                {"name": "bug"},
                {"name": "priority/critical"},
                {"name": "area/scheduler"},
            ],
            "repository": repositories[0],
        },
        {
            "id": 102,
            "number": 2001,
            "title": "Terraform AWS provider fails to handle new instance types",
            "body": "The AWS provider is not recognizing the new instance types (m6g, c6g) when creating EC2 instances.",
            "html_url": "https://github.com/hashicorp/terraform/issues/2001",
            "state": "closed",
            "created_at": "2023-02-10T08:45:00Z",
            "closed_at": "2023-02-15T16:20:00Z",
            "labels": [
                {"name": "bug"},
                {"name": "provider/aws"},
                {"name": "enhancement"},
            ],
            "repository": repositories[1],
        },
        {
            "id": 103,
            "number": 3001,
            "title": "Improve GitOps reconciliation performance",
            "body": "The GitOps reconciliation process is slow when dealing with large repositories.",
            "html_url": "https://github.com/fluxcd/flux/issues/3001",
            "state": "closed",
            "created_at": "2023-03-05T09:15:00Z",
            "closed_at": "2023-03-10T11:30:00Z",
            "labels": [
                {"name": "enhancement"},
                {"name": "performance"},
                {"name": "gitops"},
            ],
            "repository": repositories[2],
        },
    ]

    with open(os.path.join(output_dir, "repositories.json"), "w") as f:
        json.dump(repositories, f, indent=2)

    issues_path = os.path.join(output_dir, "issues.json")
    with open(issues_path, "w") as f:
        json.dump(issues, f, indent=2)

    training_data = []

    for issue in issues:
        repo = issue["repository"]
        repo_name = repo["full_name"]
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

        if issue_number == 1001:
            output_text = "Fixed by updating the zone anti-affinity logic in scheduler/zone.go. The fix ensures that pods are properly distributed across zones according to the anti-affinity rules."
        elif issue_number == 2001:
            output_text = "Updated the instance type validation in the AWS provider to include all current instance types including the ARM-based ones (m6g, c6g)."
        elif issue_number == 3001:
            output_text = "Improved GitOps reconciliation performance by implementing caching for repository metadata and optimizing the diff algorithm."
        else:
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
                "action": "read_comment",
                "observation": "I've reproduced this issue in our test environment.",
                "response": "I'm considering this information in my analysis.",
            },
            {
                "action": "read_comment",
                "observation": (
                    "The issue is in the zone.go file. The scheduler is not correctly applying the zone anti-affinity rules when calculating pod placement."
                    if issue_number == 1001
                    else (
                        "I've implemented a fix by updating the instance type validation in the AWS provider. PR #2002 has the changes."
                        if issue_number == 2001
                        else "The performance issue is related to the diff algorithm used for reconciliation."
                    )
                ),
                "response": "I'm considering this information in my analysis.",
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

        training_data.append(
            {
                "input": input_text,
                "output": output_text,
                "metadata": metadata,
                "trajectory": trajectory,
            }
        )

    training_data_path = os.path.join(output_dir, "training_data.json")
    with open(training_data_path, "w") as f:
        json.dump(training_data, f, indent=2)

    return {
        "repositories": repositories,
        "issues": issues,
        "issues_path": issues_path,
        "training_data_path": training_data_path,
        "training_data": training_data,
    }


async def mock_test_gitee_scraper(
    output_dir: str = "./data/test/gitee",
) -> Dict[str, Any]:
    """
    Mock test for Gitee scraper.

    Args:
        output_dir: Output directory for test results

    Returns:
        Test results
    """
    logging.info("Mock testing Gitee scraper")

    os.makedirs(output_dir, exist_ok=True)

    repositories = [
        {
            "id": 2001,
            "path": "kubernetes",
            "namespace": {"path": "mindspore"},
            "html_url": "https://gitee.com/mindspore/kubernetes",
            "description": "Kubernetes implementation for Chinese cloud providers",
            "stargazers_count": 5000,
            "topics": ["kubernetes", "k8s", "container", "orchestration"],
            "language": "Go",
        },
        {
            "id": 2002,
            "path": "terraform",
            "namespace": {"path": "openeuler"},
            "html_url": "https://gitee.com/openeuler/terraform",
            "description": "Terraform implementation for Chinese cloud providers",
            "stargazers_count": 3000,
            "topics": ["terraform", "infrastructure", "iac", "cloud"],
            "language": "Go",
        },
        {
            "id": 2003,
            "path": "gitops",
            "namespace": {"path": "opengauss"},
            "html_url": "https://gitee.com/opengauss/gitops",
            "description": "GitOps tools for Chinese cloud providers",
            "stargazers_count": 2000,
            "topics": ["gitops", "kubernetes", "devops", "cicd"],
            "language": "Go",
        },
    ]

    issues = [
        {
            "id": 201,
            "number": 1001,
            "title": "修复多区域集群中的 Pod 调度问题",  # Fix pod scheduling issue in multi-zone clusters
            "body": "在跨多个区域部署 Pod 时，调度程序不遵守区域反亲和性规则。",  # When deploying pods across multiple zones, the scheduler is not respecting zone anti-affinity rules.
            "html_url": "https://gitee.com/mindspore/kubernetes/issues/1001",
            "state": "closed",
            "created_at": "2023-01-15T10:30:00Z",
            "closed_at": "2023-01-20T14:45:00Z",
            "labels": [
                {"name": "bug"},
                {"name": "priority/critical"},
                {"name": "area/scheduler"},
            ],
            "repository": repositories[0],
        },
        {
            "id": 202,
            "number": 2001,
            "title": "Terraform 华为云提供程序无法处理新实例类型",  # Terraform Huawei Cloud provider fails to handle new instance types
            "body": "华为云提供程序无法识别新的实例类型（c6s、c7）。",  # The Huawei Cloud provider is not recognizing the new instance types (c6s, c7).
            "html_url": "https://gitee.com/openeuler/terraform/issues/2001",
            "state": "closed",
            "created_at": "2023-02-10T08:45:00Z",
            "closed_at": "2023-02-15T16:20:00Z",
            "labels": [
                {"name": "bug"},
                {"name": "provider/huaweicloud"},
                {"name": "enhancement"},
            ],
            "repository": repositories[1],
        },
        {
            "id": 203,
            "number": 3001,
            "title": "提高 GitOps 协调性能",  # Improve GitOps reconciliation performance
            "body": "处理大型存储库时，GitOps 协调过程很慢。",  # The GitOps reconciliation process is slow when dealing with large repositories.
            "html_url": "https://gitee.com/opengauss/gitops/issues/3001",
            "state": "closed",
            "created_at": "2023-03-05T09:15:00Z",
            "closed_at": "2023-03-10T11:30:00Z",
            "labels": [
                {"name": "enhancement"},
                {"name": "performance"},
                {"name": "gitops"},
            ],
            "repository": repositories[2],
        },
    ]

    with open(os.path.join(output_dir, "repositories.json"), "w") as f:
        json.dump(repositories, f, indent=2)

    issues_path = os.path.join(output_dir, "issues.json")
    with open(issues_path, "w") as f:
        json.dump(issues, f, indent=2)

    training_data = []

    for issue in issues:
        repo = issue["repository"]
        repo_name = f"{repo['namespace']['path']}/{repo['path']}"
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

        if issue_number == 1001:
            output_text = "通过更新 scheduler/zone.go 中的区域反亲和性逻辑来修复。该修复确保根据反亲和性规则正确分配 Pod。"  # Fixed by updating the zone anti-affinity logic in scheduler/zone.go. The fix ensures that pods are properly distributed across zones according to the anti-affinity rules.
        elif issue_number == 2001:
            output_text = "更新了华为云提供程序中的实例类型验证，以包括所有当前实例类型，包括 c6s 和 c7。"  # Updated the instance type validation in the Huawei Cloud provider to include all current instance types including c6s and c7.
        elif issue_number == 3001:
            output_text = "通过为存储库元数据实现缓存并优化差异算法，提高了 GitOps 协调性能。"  # Improved GitOps reconciliation performance by implementing caching for repository metadata and optimizing the diff algorithm.
        else:
            output_text = "问题已成功解决。"  # Issue resolved successfully.

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
                "response": "我将分析这个问题以找到解决方案。",  # I'll analyze this issue to find a solution.
            },
            {
                "action": "analyze_issue",
                "observation": issue_body,
                "response": "根据问题描述，我需要了解问题并找到解决方案。",  # Based on the issue description, I need to understand the problem and find a solution.
            },
            {
                "action": "read_comment",
                "observation": "我已经在测试环境中重现了这个问题。",  # I've reproduced this issue in our test environment.
                "response": "我正在考虑这些信息进行分析。",  # I'm considering this information in my analysis.
            },
            {
                "action": "read_comment",
                "observation": (
                    "问题出在 zone.go 文件中。调度程序在计算 Pod 放置时没有正确应用区域反亲和性规则。"
                    if issue_number
                    == 1001  # The issue is in the zone.go file. The scheduler is not correctly applying the zone anti-affinity rules when calculating pod placement.
                    else (
                        "我已经通过更新华为云提供程序中的实例类型验证实现了修复。PR #2002 包含了这些更改。"
                        if issue_number
                        == 2001  # I've implemented a fix by updating the instance type validation in the Huawei Cloud provider. PR #2002 has the changes.
                        else "性能问题与协调使用的差异算法有关。"
                    )
                ),  # The performance issue is related to the diff algorithm used for reconciliation.
                "response": "我正在考虑这些信息进行分析。",  # I'm considering this information in my analysis.
            },
            {
                "action": "implement_solution",
                "observation": "测试解决方案...",  # Testing the solution...
                "response": output_text,
            },
            {
                "action": "verify_solution",
                "observation": "解决方案已实施并测试。",  # The solution has been implemented and tested.
                "response": "问题已成功解决。",  # The issue has been resolved successfully.
            },
        ]

        training_data.append(
            {
                "input": input_text,
                "output": output_text,
                "metadata": metadata,
                "trajectory": trajectory,
            }
        )

    training_data_path = os.path.join(output_dir, "training_data.json")
    with open(training_data_path, "w") as f:
        json.dump(training_data, f, indent=2)

    return {
        "repositories": repositories,
        "issues": issues,
        "issues_path": issues_path,
        "training_data_path": training_data_path,
        "training_data": training_data,
    }


async def mock_test_issue_collector(
    output_dir: str = "./data/test/collector",
) -> Dict[str, Any]:
    """
    Mock test for issue collector.

    Args:
        output_dir: Output directory for test results

    Returns:
        Test results
    """
    logging.info("Mock testing issue collector")

    os.makedirs(output_dir, exist_ok=True)

    github_results = await mock_test_github_scraper(os.path.join(output_dir, "github"))
    gitee_results = await mock_test_gitee_scraper(os.path.join(output_dir, "gitee"))

    combined_training_data = (
        github_results["training_data"] + gitee_results["training_data"]
    )

    combined_training_data_path = os.path.join(
        output_dir, "combined_training_data.json"
    )
    with open(combined_training_data_path, "w") as f:
        json.dump(combined_training_data, f, indent=2)

    return {
        "github": github_results,
        "gitee": gitee_results,
        "combined_training_data": combined_training_data,
        "combined_training_data_path": combined_training_data_path,
    }


async def main():
    """
    Main function to run the mock tests.
    """
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )

    output_dir = "./data/test"
    os.makedirs(output_dir, exist_ok=True)

    try:
        github_results = await mock_test_github_scraper(
            os.path.join(output_dir, "github")
        )

        gitee_results = await mock_test_gitee_scraper(os.path.join(output_dir, "gitee"))

        collector_results = await mock_test_issue_collector(
            os.path.join(output_dir, "collector")
        )

        test_results = {
            "github": {
                "repositories_count": len(github_results["repositories"]),
                "issues_count": len(github_results["issues"]),
                "training_data_count": len(github_results["training_data"]),
                "topics": ["kubernetes", "terraform", "gitops"],
                "languages": ["Go"],
            },
            "gitee": {
                "repositories_count": len(gitee_results["repositories"]),
                "issues_count": len(gitee_results["issues"]),
                "training_data_count": len(gitee_results["training_data"]),
                "topics": ["kubernetes", "terraform", "gitops"],
                "languages": ["Go"],
            },
            "combined_training_data_count": len(
                collector_results["combined_training_data"]
            ),
            "data_format": {
                "includes_issue_description": True,
                "includes_solution": True,
                "includes_repository_metadata": True,
                "includes_trajectory": True,
            },
        }

        with open(os.path.join(output_dir, "test_results.json"), "w") as f:
            json.dump(test_results, f, indent=2)

        logging.info("All mock tests completed successfully")

        example_data_path = os.path.join(output_dir, "example_training_data.json")
        with open(example_data_path, "w") as f:
            json.dump(collector_results["combined_training_data"][:2], f, indent=2)

        logging.info(
            f"Saved example training data to {os.path.abspath(example_data_path)}"
        )

        return {
            "github_results": github_results,
            "gitee_results": gitee_results,
            "collector_results": collector_results,
            "test_results": test_results,
            "example_data_path": example_data_path,
        }

    except Exception as e:
        logging.error(f"Error running mock tests: {str(e)}")
        raise


if __name__ == "__main__":
    asyncio.run(main())
