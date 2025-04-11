# Fine-Tune

A web scraper for collecting solved issues from GitHub and Gitee repositories to fine-tune Llama 4 Maverick and Scout models.

## Overview

This project provides tools to scrape solved issues from GitHub and Gitee repositories, focusing on GitOps, Terraform, and Kubernetes domains. The collected data is formatted for fine-tuning Llama 4 Maverick and Scout models, capturing issue trajectories similar to SWE Agent benchmarking.

## Features

- Scrape repositories from GitHub and Gitee based on topics and languages
- Collect closed issues with solutions
- Format issues for model training with trajectories
- Support for GitOps, Terraform, and Kubernetes repositories
- Mock testing capabilities for development without API keys

## Installation

```bash
# Clone the repository
git clone https://github.com/spectrumwebco/Fine-Tune.git
cd Fine-Tune

# Install dependencies
pip install -e .
```

## Usage

### Configuration

Set up your API keys as environment variables:

```bash
export GITHUB_API_KEY="your_github_api_key"
export GITEE_API_KEY="your_gitee_api_key"
```

### Collecting Issues

```python
import asyncio
from src.integrations.issue_collector import IssueCollector

async def main():
    collector = IssueCollector(
        github_api_key="your_github_api_key",
        gitee_api_key="your_gitee_api_key",
        output_dir="./data/collected_issues"
    )
    
    results = await collector.collect_and_save(
        topics=["gitops", "terraform", "kubernetes", "k8s"],
        languages=["python", "go"],
        min_stars=100,
        max_repos_per_platform=25,
        max_issues_per_repo=50
    )
    
    print(f"Collected data saved to: {results['combined_training_data_path']}")

if __name__ == "__main__":
    asyncio.run(main())
```

### Mock Testing

For development without API keys, use the mock testing functionality:

```python
import asyncio
from src.mock_test_scrapers import mock_test_issue_collector

async def main():
    results = await mock_test_issue_collector()
    print(f"Mock data saved to: {results['combined_training_data_path']}")

if __name__ == "__main__":
    asyncio.run(main())
```

## Data Format

The collected data is formatted for fine-tuning Llama 4 models with the following structure:

```json
{
  "input": "Repository: kubernetes/kubernetes\nTopics: kubernetes, k8s, container, orchestration\nIssue Title: Fix pod scheduling issue in multi-zone clusters\nIssue Description:\nWhen deploying pods across multiple zones, the scheduler is not respecting zone anti-affinity rules.\n",
  "output": "Fixed by updating the zone anti-affinity logic in scheduler/zone.go. The fix ensures that pods are properly distributed across zones according to the anti-affinity rules.",
  "metadata": {
    "issue_id": 101,
    "issue_number": 1001,
    "repository": "kubernetes/kubernetes",
    "url": "https://github.com/kubernetes/kubernetes/issues/1001",
    "created_at": "2023-01-15T10:30:00Z",
    "closed_at": "2023-01-20T14:45:00Z",
    "labels": ["bug", "priority/critical", "area/scheduler"]
  },
  "trajectory": [
    {
      "action": "read_issue",
      "observation": "Issue #1001: Fix pod scheduling issue in multi-zone clusters",
      "response": "I'll analyze this issue to find a solution."
    },
    // Additional trajectory steps...
  ]
}
```

## Integration with ML Infrastructure

The collected data can be used with:

- KubeFlow for model fine-tuning
- MLFlow for experiment tracking
- KServe for model serving
- JupyterLab notebooks for experimentation

## License

MIT
