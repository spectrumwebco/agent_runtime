"""
Data validation module for Llama 4 fine-tuning.

This module provides functionality to validate data collected from GitHub and Gitee
issue scrapers for fine-tuning Llama 4 models.
"""

import os
import json
import logging
import pandas as pd
import numpy as np
from typing import Dict, List, Optional, Union, Any
from datetime import datetime
import re
import jsonschema
from jsonschema import validate

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


class DataValidator:
    """
    Data validator class for Llama 4 fine-tuning.
    """

    def __init__(
        self,
        input_dir: str = "data",
        output_dir: str = "validated_data",
        schema_dir: str = "schemas",
    ):
        """
        Initialize the DataValidator.

        Args:
            input_dir: Directory containing data to validate
            output_dir: Directory to store validation results
            schema_dir: Directory containing JSON schemas for validation
        """
        self.input_dir = input_dir
        self.output_dir = output_dir
        self.schema_dir = schema_dir
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        os.makedirs(self.schema_dir, exist_ok=True)
        
        self.schemas = {
            "raw": self._get_raw_data_schema(),
            "chat": self._get_chat_format_schema(),
            "completion": self._get_completion_format_schema(),
        }
        
        logger.info("DataValidator initialized")

    def _get_raw_data_schema(self) -> Dict[str, Any]:
        """
        Get JSON schema for raw data validation.

        Returns:
            JSON schema for raw data validation
        """
        schema = {
            "type": "object",
            "required": ["input", "output", "metadata"],
            "properties": {
                "input": {
                    "type": "object",
                    "required": ["title", "description"],
                    "properties": {
                        "repository": {"type": "string"},
                        "topics": {
                            "type": "array",
                            "items": {"type": "string"}
                        },
                        "title": {"type": "string"},
                        "description": {"type": "string"},
                    },
                },
                "output": {
                    "type": "object",
                    "required": ["solution"],
                    "properties": {
                        "solution": {"type": "string"},
                    },
                },
                "metadata": {
                    "type": "object",
                    "required": ["id", "source"],
                    "properties": {
                        "id": {"type": "string"},
                        "source": {"type": "string"},
                        "repository": {"type": "string"},
                        "url": {"type": "string"},
                        "created_at": {"type": "string"},
                        "closed_at": {"type": "string"},
                        "labels": {
                            "type": "array",
                            "items": {"type": "string"}
                        },
                    },
                },
                "trajectory": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "required": ["step", "action", "content"],
                        "properties": {
                            "step": {"type": "integer"},
                            "action": {"type": "string"},
                            "content": {"type": "string"},
                            "timestamp": {"type": "string"},
                            "user": {"type": "string"},
                        },
                    },
                },
            },
        }
        
        return schema

    def _get_chat_format_schema(self) -> Dict[str, Any]:
        """
        Get JSON schema for chat format validation.

        Returns:
            JSON schema for chat format validation
        """
        schema = {
            "type": "object",
            "required": ["messages"],
            "properties": {
                "messages": {
                    "type": "array",
                    "minItems": 3,
                    "items": {
                        "type": "object",
                        "required": ["role", "content"],
                        "properties": {
                            "role": {
                                "type": "string",
                                "enum": ["system", "user", "assistant"]
                            },
                            "content": {"type": "string"},
                        },
                    },
                },
                "metadata": {
                    "type": "object",
                    "properties": {
                        "id": {"type": "string"},
                        "source": {"type": "string"},
                        "repository": {"type": "string"},
                        "url": {"type": "string"},
                        "created_at": {"type": "string"},
                        "closed_at": {"type": "string"},
                        "labels": {
                            "type": "array",
                            "items": {"type": "string"}
                        },
                    },
                },
                "trajectory": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "required": ["step", "action", "content"],
                        "properties": {
                            "step": {"type": "integer"},
                            "action": {"type": "string"},
                            "content": {"type": "string"},
                            "timestamp": {"type": "string"},
                            "user": {"type": "string"},
                        },
                    },
                },
            },
        }
        
        return schema

    def _get_completion_format_schema(self) -> Dict[str, Any]:
        """
        Get JSON schema for completion format validation.

        Returns:
            JSON schema for completion format validation
        """
        schema = {
            "type": "object",
            "required": ["prompt", "completion"],
            "properties": {
                "prompt": {"type": "string"},
                "completion": {"type": "string"},
                "metadata": {
                    "type": "object",
                    "properties": {
                        "id": {"type": "string"},
                        "source": {"type": "string"},
                        "repository": {"type": "string"},
                        "url": {"type": "string"},
                        "created_at": {"type": "string"},
                        "closed_at": {"type": "string"},
                        "labels": {
                            "type": "array",
                            "items": {"type": "string"}
                        },
                    },
                },
                "trajectory": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "required": ["step", "action", "content"],
                        "properties": {
                            "step": {"type": "integer"},
                            "action": {"type": "string"},
                            "content": {"type": "string"},
                            "timestamp": {"type": "string"},
                            "user": {"type": "string"},
                        },
                    },
                },
            },
        }
        
        return schema

    def save_schemas(self) -> Dict[str, str]:
        """
        Save JSON schemas to files.

        Returns:
            Dictionary containing paths to saved schema files
        """
        logger.info("Saving JSON schemas")
        
        os.makedirs(self.schema_dir, exist_ok=True)
        
        schema_paths = {}
        for schema_name, schema in self.schemas.items():
            schema_filename = f"{schema_name}_schema.json"
            schema_path = os.path.join(self.schema_dir, schema_filename)
            
            with open(schema_path, "w") as f:
                json.dump(schema, f, indent=2)
            
            schema_paths[schema_name] = schema_path
            logger.info(f"Saved {schema_name} schema to {schema_path}")
        
        return schema_paths

    def load_data(self, filename: str) -> List[Dict[str, Any]]:
        """
        Load data from a JSON file.

        Args:
            filename: Name of the file to load data from

        Returns:
            List of data examples
        """
        logger.info(f"Loading data from {filename}")
        
        input_path = os.path.join(self.input_dir, filename)
        with open(input_path, "r") as f:
            data = json.load(f)
        
        logger.info(f"Loaded {len(data)} examples from {filename}")
        return data

    def validate_data(
        self, data: List[Dict[str, Any]], schema_name: str = "raw"
    ) -> Dict[str, Any]:
        """
        Validate data against a JSON schema.

        Args:
            data: List of data examples to validate
            schema_name: Name of the schema to validate against

        Returns:
            Dictionary containing validation results
        """
        logger.info(f"Validating {len(data)} examples against {schema_name} schema")
        
        schema = self.schemas.get(schema_name)
        if not schema:
            raise ValueError(f"Schema {schema_name} not found")
        
        valid_examples = []
        invalid_examples = []
        validation_errors = []
        
        for i, example in enumerate(data):
            try:
                validate(instance=example, schema=schema)
                valid_examples.append(example)
            except jsonschema.exceptions.ValidationError as e:
                invalid_examples.append(example)
                validation_errors.append({
                    "index": i,
                    "error": str(e),
                })
        
        total_examples = len(data)
        valid_count = len(valid_examples)
        invalid_count = len(invalid_examples)
        valid_ratio = valid_count / total_examples if total_examples > 0 else 0
        
        validation_results = {
            "schema_name": schema_name,
            "total_examples": total_examples,
            "valid_examples": valid_count,
            "invalid_examples": invalid_count,
            "valid_ratio": valid_ratio,
            "validation_errors": validation_errors,
        }
        
        logger.info(f"Validation results: {valid_count}/{total_examples} examples valid ({valid_ratio:.2%})")
        return validation_results

    def save_validation_results(
        self, validation_results: Dict[str, Any], filename: str = "validation_results.json"
    ) -> str:
        """
        Save validation results to a file.

        Args:
            validation_results: Dictionary containing validation results
            filename: Name of the file to save validation results to

        Returns:
            Path to the saved file
        """
        logger.info(f"Saving validation results to {filename}")
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        output_path = os.path.join(self.output_dir, filename)
        with open(output_path, "w") as f:
            json.dump(validation_results, f, indent=2)
        
        logger.info(f"Validation results saved to {output_path}")
        return output_path

    def save_valid_examples(
        self, data: List[Dict[str, Any]], validation_results: Dict[str, Any], filename: str = "valid_examples.json"
    ) -> str:
        """
        Save valid examples to a file.

        Args:
            data: List of data examples
            validation_results: Dictionary containing validation results
            filename: Name of the file to save valid examples to

        Returns:
            Path to the saved file
        """
        logger.info(f"Saving valid examples to {filename}")
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        invalid_indices = [error["index"] for error in validation_results.get("validation_errors", [])]
        
        valid_examples = [example for i, example in enumerate(data) if i not in invalid_indices]
        
        output_path = os.path.join(self.output_dir, filename)
        with open(output_path, "w") as f:
            json.dump(valid_examples, f, indent=2)
        
        logger.info(f"Saved {len(valid_examples)} valid examples to {output_path}")
        return output_path

    def validate_file(
        self, filename: str, schema_name: str = "raw", save_valid: bool = True
    ) -> Dict[str, str]:
        """
        Validate data in a file against a JSON schema.

        Args:
            filename: Name of the file to validate
            schema_name: Name of the schema to validate against
            save_valid: Whether to save valid examples to a file

        Returns:
            Dictionary containing paths to saved files
        """
        logger.info(f"Validating file {filename} against {schema_name} schema")
        
        data = self.load_data(filename)
        
        validation_results = self.validate_data(data, schema_name=schema_name)
        
        results_filename = f"{os.path.splitext(filename)[0]}_validation_results.json"
        results_path = self.save_validation_results(validation_results, filename=results_filename)
        
        output_paths = {
            "results": results_path,
        }
        
        if save_valid:
            valid_filename = f"{os.path.splitext(filename)[0]}_valid.json"
            valid_path = self.save_valid_examples(data, validation_results, filename=valid_filename)
            output_paths["valid"] = valid_path
        
        logger.info(f"Validation completed: {output_paths}")
        return output_paths

    def check_data_quality(self, data: List[Dict[str, Any]]) -> Dict[str, Any]:
        """
        Check data quality metrics.

        Args:
            data: List of data examples

        Returns:
            Dictionary containing data quality metrics
        """
        logger.info(f"Checking data quality for {len(data)} examples")
        
        quality_metrics = {
            "total_examples": len(data),
            "empty_fields": {
                "input_title": 0,
                "input_description": 0,
                "output_solution": 0,
            },
            "length_metrics": {
                "input_title": {
                    "min": float("inf"),
                    "max": 0,
                    "mean": 0,
                    "median": 0,
                },
                "input_description": {
                    "min": float("inf"),
                    "max": 0,
                    "mean": 0,
                    "median": 0,
                },
                "output_solution": {
                    "min": float("inf"),
                    "max": 0,
                    "mean": 0,
                    "median": 0,
                },
            },
            "source_distribution": {},
            "repository_distribution": {},
            "topic_distribution": {},
            "label_distribution": {},
        }
        
        title_lengths = []
        description_lengths = []
        solution_lengths = []
        sources = {}
        repositories = {}
        topics = {}
        labels = {}
        
        for example in data:
            input_data = example.get("input", {})
            title = input_data.get("title", "")
            description = input_data.get("description", "")
            
            if not title:
                quality_metrics["empty_fields"]["input_title"] += 1
            else:
                title_length = len(title)
                title_lengths.append(title_length)
                quality_metrics["length_metrics"]["input_title"]["min"] = min(quality_metrics["length_metrics"]["input_title"]["min"], title_length)
                quality_metrics["length_metrics"]["input_title"]["max"] = max(quality_metrics["length_metrics"]["input_title"]["max"], title_length)
            
            if not description:
                quality_metrics["empty_fields"]["input_description"] += 1
            else:
                description_length = len(description)
                description_lengths.append(description_length)
                quality_metrics["length_metrics"]["input_description"]["min"] = min(quality_metrics["length_metrics"]["input_description"]["min"], description_length)
                quality_metrics["length_metrics"]["input_description"]["max"] = max(quality_metrics["length_metrics"]["input_description"]["max"], description_length)
            
            output_data = example.get("output", {})
            solution = output_data.get("solution", "")
            
            if not solution:
                quality_metrics["empty_fields"]["output_solution"] += 1
            else:
                solution_length = len(solution)
                solution_lengths.append(solution_length)
                quality_metrics["length_metrics"]["output_solution"]["min"] = min(quality_metrics["length_metrics"]["output_solution"]["min"], solution_length)
                quality_metrics["length_metrics"]["output_solution"]["max"] = max(quality_metrics["length_metrics"]["output_solution"]["max"], solution_length)
            
            metadata = example.get("metadata", {})
            source = metadata.get("source", "unknown")
            repository = metadata.get("repository", "unknown")
            example_topics = input_data.get("topics", [])
            example_labels = metadata.get("labels", [])
            
            sources[source] = sources.get(source, 0) + 1
            repositories[repository] = repositories.get(repository, 0) + 1
            
            for topic in example_topics:
                topics[topic] = topics.get(topic, 0) + 1
            
            for label in example_labels:
                labels[label] = labels.get(label, 0) + 1
        
        if title_lengths:
            quality_metrics["length_metrics"]["input_title"]["mean"] = np.mean(title_lengths)
            quality_metrics["length_metrics"]["input_title"]["median"] = np.median(title_lengths)
        
        if description_lengths:
            quality_metrics["length_metrics"]["input_description"]["mean"] = np.mean(description_lengths)
            quality_metrics["length_metrics"]["input_description"]["median"] = np.median(description_lengths)
        
        if solution_lengths:
            quality_metrics["length_metrics"]["output_solution"]["mean"] = np.mean(solution_lengths)
            quality_metrics["length_metrics"]["output_solution"]["median"] = np.median(solution_lengths)
        
        quality_metrics["source_distribution"] = sources
        quality_metrics["repository_distribution"] = repositories
        quality_metrics["topic_distribution"] = topics
        quality_metrics["label_distribution"] = labels
        
        logger.info("Data quality check completed")
        return quality_metrics

    def save_quality_metrics(
        self, quality_metrics: Dict[str, Any], filename: str = "quality_metrics.json"
    ) -> str:
        """
        Save data quality metrics to a file.

        Args:
            quality_metrics: Dictionary containing data quality metrics
            filename: Name of the file to save quality metrics to

        Returns:
            Path to the saved file
        """
        logger.info(f"Saving data quality metrics to {filename}")
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        output_path = os.path.join(self.output_dir, filename)
        with open(output_path, "w") as f:
            json.dump(quality_metrics, f, indent=2)
        
        logger.info(f"Data quality metrics saved to {output_path}")
        return output_path

    def check_file_quality(self, filename: str) -> Dict[str, str]:
        """
        Check data quality metrics for a file.

        Args:
            filename: Name of the file to check

        Returns:
            Dictionary containing paths to saved files
        """
        logger.info(f"Checking data quality for file {filename}")
        
        data = self.load_data(filename)
        
        quality_metrics = self.check_data_quality(data)
        
        metrics_filename = f"{os.path.splitext(filename)[0]}_quality_metrics.json"
        metrics_path = self.save_quality_metrics(quality_metrics, filename=metrics_filename)
        
        output_paths = {
            "metrics": metrics_path,
        }
        
        logger.info(f"Data quality check completed: {output_paths}")
        return output_paths


def main():
    """
    Main function for testing the DataValidator.
    """
    validator = DataValidator(
        input_dir="data",
        output_dir="validated_data",
        schema_dir="schemas",
    )
    
    schema_paths = validator.save_schemas()
    logger.info(f"Saved schemas: {schema_paths}")
    
    validation_paths = validator.validate_file(
        filename="github_gitee_issues.json",
        schema_name="raw",
        save_valid=True,
    )
    logger.info(f"Validation completed: {validation_paths}")
    
    quality_paths = validator.check_file_quality(
        filename="github_gitee_issues.json",
    )
    logger.info(f"Data quality check completed: {quality_paths}")


if __name__ == "__main__":
    main()
