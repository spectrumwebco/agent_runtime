"""
Script to demonstrate Pydantic integration in the ML infrastructure.

This script shows how to use the Pydantic models defined in the ML infrastructure
to validate data and ensure type safety across the application.
"""

import os
import sys
import json
from typing import List, Dict, Any, Optional

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from src.ml_infrastructure.data.validation.models import (
    RawDataModel,
    ChatFormatModel,
    CompletionFormatModel,
    ValidationResult,
    QualityMetrics
)
from src.ml_infrastructure.api.models import (
    ModelList,
    ModelDetail,
    FineTuningJobCreate,
    FineTuningJobDetail,
    InferenceServiceCreate
)
from src.integrations.github.models import (
    GitHubRepository,
    GitHubIssue
)
from src.integrations.gitee.models import (
    GiteeRepository,
    GiteeIssue
)
from src.integrations.issue_collector.models import (
    CollectionConfig,
    CollectionResult,
    TrainingExample
)


def demonstrate_pydantic_validation():
    """Demonstrate Pydantic validation with example data."""
    print("Demonstrating Pydantic validation for ML infrastructure")
    
    raw_data = {
        "input": {
            "repository": "kubernetes/kubernetes",
            "topics": ["kubernetes", "gitops"],
            "title": "Issue with pod scheduling",
            "description": "Pods are not being scheduled correctly in my cluster."
        },
        "output": {
            "solution": "Check the node affinity rules and ensure resources are available."
        },
        "metadata": {
            "id": "12345",
            "source": "github",
            "repository": "kubernetes/kubernetes",
            "url": "https://github.com/kubernetes/kubernetes/issues/12345",
            "created_at": "2023-01-01T00:00:00Z",
            "closed_at": "2023-01-02T00:00:00Z",
            "labels": ["bug", "priority/medium"]
        },
        "trajectory": [
            {
                "step": 1,
                "action": "read_issue",
                "content": "Reading issue details",
                "timestamp": "2023-01-01T12:00:00Z"
            }
        ]
    }
    
    try:
        validated_raw_data = RawDataModel(**raw_data)
        print(f"✅ Raw data validation successful")
        print(f"  - Repository: {validated_raw_data.input.repository}")
        print(f"  - Issue title: {validated_raw_data.input.title}")
        print(f"  - Solution: {validated_raw_data.output.solution}")
    except Exception as e:
        print(f"❌ Raw data validation failed: {str(e)}")
    
    chat_data = {
        "messages": [
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": "I have an issue with pod scheduling in Kubernetes."},
            {"role": "assistant", "content": "Let me help you troubleshoot that."}
        ],
        "metadata": {
            "id": "12345",
            "source": "github",
            "repository": "kubernetes/kubernetes"
        }
    }
    
    try:
        validated_chat_data = ChatFormatModel(**chat_data)
        print(f"✅ Chat format validation successful")
        print(f"  - Number of messages: {len(validated_chat_data.messages)}")
        print(f"  - First message role: {validated_chat_data.messages[0].role}")
    except Exception as e:
        print(f"❌ Chat format validation failed: {str(e)}")
    
    fine_tuning_job = {
        "model_id": "llama4-maverick",
        "training_file": "training_data.jsonl",
        "validation_file": "validation_data.jsonl",
        "hyperparameters": {
            "epochs": 3,
            "batch_size": 4,
            "learning_rate": 1e-5
        },
        "suffix": "gitops-expert"
    }
    
    try:
        validated_fine_tuning_job = FineTuningJobCreate(**fine_tuning_job)
        print(f"✅ Fine-tuning job validation successful")
        print(f"  - Model ID: {validated_fine_tuning_job.model_id}")
        print(f"  - Epochs: {validated_fine_tuning_job.hyperparameters.epochs}")
    except Exception as e:
        print(f"❌ Fine-tuning job validation failed: {str(e)}")
    
    github_repo = {
        "id": 123456,
        "name": "kubernetes",
        "full_name": "kubernetes/kubernetes",
        "description": "Production-Grade Container Scheduling and Management",
        "html_url": "https://github.com/kubernetes/kubernetes",
        "stars": 98765,
        "forks": 32109,
        "topics": ["kubernetes", "containers", "orchestration"]
    }
    
    try:
        validated_github_repo = GitHubRepository(**github_repo)
        print(f"✅ GitHub repository validation successful")
        print(f"  - Repository: {validated_github_repo.full_name}")
        print(f"  - Stars: {validated_github_repo.stars}")
        print(f"  - Topics: {', '.join(validated_github_repo.topics)}")
    except Exception as e:
        print(f"❌ GitHub repository validation failed: {str(e)}")
    
    collection_config = {
        "topics": ["gitops", "terraform", "kubernetes"],
        "languages": ["python", "go"],
        "min_stars": 100,
        "max_repos_per_platform": 25,
        "max_issues_per_repo": 50,
        "include_pull_requests": False
    }
    
    try:
        validated_collection_config = CollectionConfig(**collection_config)
        print(f"✅ Collection config validation successful")
        print(f"  - Topics: {', '.join(validated_collection_config.topics)}")
        print(f"  - Languages: {', '.join(validated_collection_config.languages)}")
        print(f"  - Min stars: {validated_collection_config.min_stars}")
    except Exception as e:
        print(f"❌ Collection config validation failed: {str(e)}")
    
    print("\nPydantic integration ensures type safety and validation across the ML infrastructure")
    print("This enables seamless integration with the Autonomous AI Agent and Kubernetes deployment")


if __name__ == "__main__":
    demonstrate_pydantic_validation()
