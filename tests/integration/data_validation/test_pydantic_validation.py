"""
Test module for Pydantic validation in the data validation module.
"""

import os
import json
import pytest
from typing import Dict, List, Any

from src.ml_infrastructure.data.validation.models import (
    RawDataModel,
    ChatFormatModel,
    CompletionFormatModel,
    ValidationResult,
    QualityMetrics
)
from src.ml_infrastructure.data.validation.validator import DataValidator


@pytest.fixture
def valid_raw_data():
    """Fixture for valid raw data."""
    return {
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


@pytest.fixture
def invalid_raw_data():
    """Fixture for invalid raw data."""
    return {
        "input": {
            "repository": "kubernetes/kubernetes",
            "topics": ["kubernetes", "gitops"],
            "description": "Pods are not being scheduled correctly in my cluster."
        },
        "output": {
            "solution": "Check the node affinity rules and ensure resources are available."
        },
        "metadata": {
            "id": "12345",
            "source": "github"
        }
    }


@pytest.fixture
def valid_chat_data():
    """Fixture for valid chat format data."""
    return {
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


@pytest.fixture
def valid_completion_data():
    """Fixture for valid completion format data."""
    return {
        "prompt": "Fix the issue with pod scheduling in Kubernetes",
        "completion": "Check the node affinity rules and ensure resources are available.",
        "metadata": {
            "id": "12345",
            "source": "github",
            "repository": "kubernetes/kubernetes"
        }
    }


def test_raw_data_model_validation(valid_raw_data, invalid_raw_data):
    """Test RawDataModel validation."""
    model = RawDataModel(**valid_raw_data)
    assert model.input.title == "Issue with pod scheduling"
    assert model.output.solution == "Check the node affinity rules and ensure resources are available."
    assert model.metadata.id == "12345"
    
    with pytest.raises(Exception):
        RawDataModel(**invalid_raw_data)


def test_chat_format_model_validation(valid_chat_data):
    """Test ChatFormatModel validation."""
    model = ChatFormatModel(**valid_chat_data)
    assert len(model.messages) == 3
    assert model.messages[0].role == "system"
    assert model.messages[1].role == "user"
    assert model.messages[2].role == "assistant"


def test_completion_format_model_validation(valid_completion_data):
    """Test CompletionFormatModel validation."""
    model = CompletionFormatModel(**valid_completion_data)
    assert model.prompt == "Fix the issue with pod scheduling in Kubernetes"
    assert model.completion == "Check the node affinity rules and ensure resources are available."


def test_data_validator_with_pydantic(valid_raw_data, invalid_raw_data, tmp_path):
    """Test DataValidator with Pydantic models."""
    input_dir = tmp_path / "input"
    output_dir = tmp_path / "output"
    schema_dir = tmp_path / "schemas"
    
    input_dir.mkdir()
    
    test_data = [valid_raw_data, valid_raw_data, invalid_raw_data]
    test_file = input_dir / "test_data.json"
    with open(test_file, "w") as f:
        json.dump(test_data, f)
    
    validator = DataValidator(
        input_dir=str(input_dir),
        output_dir=str(output_dir),
        schema_dir=str(schema_dir)
    )
    
    validation_results = validator.validate_data(test_data, schema_name="raw")
    
    assert isinstance(validation_results, ValidationResult)
    assert validation_results.total_examples == 3
    assert validation_results.valid_examples == 2
    assert validation_results.invalid_examples == 1
    assert validation_results.valid_ratio == 2/3
    
    quality_metrics = validator.check_data_quality(test_data)
    assert isinstance(quality_metrics, QualityMetrics)
    assert quality_metrics.total_examples == 3


def test_validator_file_operations(valid_raw_data, tmp_path):
    """Test DataValidator file operations with Pydantic models."""
    input_dir = tmp_path / "input"
    output_dir = tmp_path / "output"
    schema_dir = tmp_path / "schemas"
    
    input_dir.mkdir()
    
    test_data = [valid_raw_data, valid_raw_data]
    test_file = input_dir / "test_data.json"
    with open(test_file, "w") as f:
        json.dump(test_data, f)
    
    validator = DataValidator(
        input_dir=str(input_dir),
        output_dir=str(output_dir),
        schema_dir=str(schema_dir)
    )
    
    result_paths = validator.validate_file(
        filename="test_data.json",
        schema_name="raw",
        save_valid=True
    )
    
    assert "results" in result_paths
    assert "valid" in result_paths
    assert os.path.exists(result_paths["results"])
    assert os.path.exists(result_paths["valid"])
    
    quality_paths = validator.check_file_quality(filename="test_data.json")
    assert "metrics" in quality_paths
    assert os.path.exists(quality_paths["metrics"])
