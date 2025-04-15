"""
Dataset versioning module for Llama 4 fine-tuning.

This module provides functionality to version datasets collected from GitHub and Gitee
issue scrapers for fine-tuning Llama 4 models using MLflow.
"""

import os
import json
import logging
import pandas as pd
import numpy as np
from typing import Dict, List, Optional, Union, Any
from datetime import datetime
import mlflow
from mlflow.tracking import MlflowClient
import hashlib

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


class DatasetVersioner:
    """
    Dataset versioner class for Llama 4 fine-tuning.
    """

    def __init__(
        self,
        input_dir: str = "data",
        mlflow_tracking_uri: Optional[str] = None,
        experiment_name: str = "llama4-fine-tuning",
    ):
        """
        Initialize the DatasetVersioner.

        Args:
            input_dir: Directory containing datasets to version
            mlflow_tracking_uri: MLflow tracking URI
            experiment_name: Name of the MLflow experiment
        """
        self.input_dir = input_dir
        self.experiment_name = experiment_name
        
        if mlflow_tracking_uri:
            mlflow.set_tracking_uri(mlflow_tracking_uri)
        self.mlflow_client = MlflowClient()
        
        experiment = self.mlflow_client.get_experiment_by_name(experiment_name)
        if experiment is None:
            self.experiment_id = self.mlflow_client.create_experiment(experiment_name)
        else:
            self.experiment_id = experiment.experiment_id
        
        logger.info(f"DatasetVersioner initialized with experiment ID: {self.experiment_id}")

    def load_dataset(self, filename: str) -> List[Dict[str, Any]]:
        """
        Load dataset from a JSON file.

        Args:
            filename: Name of the file to load dataset from

        Returns:
            List of dataset examples
        """
        logger.info(f"Loading dataset from {filename}")
        
        input_path = os.path.join(self.input_dir, filename)
        with open(input_path, "r") as f:
            dataset = json.load(f)
        
        logger.info(f"Loaded {len(dataset)} examples from {filename}")
        return dataset

    def compute_dataset_hash(self, dataset: List[Dict[str, Any]]) -> str:
        """
        Compute a hash for the dataset.

        Args:
            dataset: List of dataset examples

        Returns:
            Hash of the dataset
        """
        dataset_json = json.dumps(dataset, sort_keys=True)
        
        dataset_hash = hashlib.sha256(dataset_json.encode()).hexdigest()
        
        return dataset_hash

    def compute_dataset_stats(self, dataset: List[Dict[str, Any]]) -> Dict[str, Any]:
        """
        Compute statistics for the dataset.

        Args:
            dataset: List of dataset examples

        Returns:
            Dictionary containing dataset statistics
        """
        logger.info(f"Computing statistics for dataset with {len(dataset)} examples")
        
        stats = {
            "count": len(dataset),
            "sources": {},
            "repositories": {},
            "topics": {},
            "labels": {},
        }
        
        for example in dataset:
            metadata = example.get("metadata", {})
            source = metadata.get("source", "unknown")
            repository = metadata.get("repository", "unknown")
            
            input_data = example.get("input", {})
            topics = input_data.get("topics", [])
            
            labels = metadata.get("labels", [])
            
            stats["sources"][source] = stats["sources"].get(source, 0) + 1
            stats["repositories"][repository] = stats["repositories"].get(repository, 0) + 1
            
            for topic in topics:
                stats["topics"][topic] = stats["topics"].get(topic, 0) + 1
            
            for label in labels:
                stats["labels"][label] = stats["labels"].get(label, 0) + 1
        
        logger.info("Dataset statistics computed")
        return stats

    def log_dataset_to_mlflow(
        self,
        dataset_path: str,
        dataset_name: str,
        dataset_version: str = "1.0.0",
        tags: Optional[Dict[str, str]] = None,
    ) -> str:
        """
        Log dataset to MLflow for versioning.

        Args:
            dataset_path: Path to the dataset file
            dataset_name: Name of the dataset
            dataset_version: Version of the dataset
            tags: Tags to add to the MLflow run

        Returns:
            MLflow run ID
        """
        logger.info(f"Logging dataset to MLflow: {dataset_path}")
        
        dataset = self.load_dataset(os.path.basename(dataset_path))
        
        dataset_hash = self.compute_dataset_hash(dataset)
        
        dataset_stats = self.compute_dataset_stats(dataset)
        
        with mlflow.start_run(experiment_id=self.experiment_id) as run:
            run_id = run.info.run_id
            
            mlflow.log_artifact(dataset_path, "datasets")
            
            mlflow.log_param("dataset_name", dataset_name)
            mlflow.log_param("dataset_version", dataset_version)
            mlflow.log_param("dataset_hash", dataset_hash)
            mlflow.log_param("dataset_size", len(dataset))
            mlflow.log_param("dataset_path", dataset_path)
            mlflow.log_param("dataset_created_at", datetime.now().isoformat())
            
            for source, count in dataset_stats["sources"].items():
                mlflow.log_metric(f"source_{source}_count", count)
            
            mlflow.log_metric("total_examples_count", len(dataset))
            
            stats_path = os.path.join(os.path.dirname(dataset_path), f"{os.path.basename(dataset_path)}_stats.json")
            with open(stats_path, "w") as f:
                json.dump(dataset_stats, f, indent=2)
            
            mlflow.log_artifact(stats_path, "statistics")
            
            if tags:
                for key, value in tags.items():
                    mlflow.set_tag(key, value)
        
        logger.info(f"Dataset logged to MLflow: {run_id}")
        return run_id

    def get_dataset_versions(self, dataset_name: str) -> List[Dict[str, Any]]:
        """
        Get all versions of a dataset from MLflow.

        Args:
            dataset_name: Name of the dataset

        Returns:
            List of dataset versions
        """
        logger.info(f"Getting versions of dataset: {dataset_name}")
        
        runs = self.mlflow_client.search_runs(
            experiment_ids=[self.experiment_id],
            filter_string=f"params.dataset_name = '{dataset_name}'",
        )
        
        versions = []
        for run in runs:
            run_id = run.info.run_id
            run_data = self.mlflow_client.get_run(run_id)
            
            version = {
                "run_id": run_id,
                "dataset_name": run_data.data.params.get("dataset_name"),
                "dataset_version": run_data.data.params.get("dataset_version"),
                "dataset_hash": run_data.data.params.get("dataset_hash"),
                "dataset_size": int(run_data.data.params.get("dataset_size", 0)),
                "dataset_path": run_data.data.params.get("dataset_path"),
                "dataset_created_at": run_data.data.params.get("dataset_created_at"),
                "metrics": run_data.data.metrics,
                "tags": run_data.data.tags,
            }
            
            versions.append(version)
        
        logger.info(f"Found {len(versions)} versions of dataset: {dataset_name}")
        return versions

    def get_latest_dataset_version(self, dataset_name: str) -> Optional[Dict[str, Any]]:
        """
        Get the latest version of a dataset from MLflow.

        Args:
            dataset_name: Name of the dataset

        Returns:
            Latest version of the dataset or None if no versions found
        """
        logger.info(f"Getting latest version of dataset: {dataset_name}")
        
        versions = self.get_dataset_versions(dataset_name)
        
        versions.sort(key=lambda x: x["dataset_created_at"], reverse=True)
        
        if versions:
            logger.info(f"Latest version of dataset {dataset_name}: {versions[0]['dataset_version']}")
            return versions[0]
        
        logger.warning(f"No versions found for dataset: {dataset_name}")
        return None

    def download_dataset(self, run_id: str, output_dir: str = "downloaded_datasets") -> str:
        """
        Download a dataset from MLflow.

        Args:
            run_id: MLflow run ID
            output_dir: Directory to save the downloaded dataset

        Returns:
            Path to the downloaded dataset
        """
        logger.info(f"Downloading dataset from MLflow run: {run_id}")
        
        os.makedirs(output_dir, exist_ok=True)
        
        run = self.mlflow_client.get_run(run_id)
        dataset_name = run.data.params.get("dataset_name")
        dataset_version = run.data.params.get("dataset_version")
        
        artifacts = self.mlflow_client.list_artifacts(run_id, "datasets")
        
        for artifact in artifacts:
            artifact_path = os.path.join("datasets", artifact.path)
            output_path = os.path.join(output_dir, f"{dataset_name}_v{dataset_version}_{os.path.basename(artifact.path)}")
            
            self.mlflow_client.download_artifacts(run_id, artifact_path, output_dir)
            os.rename(os.path.join(output_dir, artifact_path), output_path)
            
            logger.info(f"Downloaded dataset to: {output_path}")
            return output_path
        
        logger.warning(f"No dataset found in MLflow run: {run_id}")
        return ""

    def compare_dataset_versions(
        self, dataset_name: str, version1: str, version2: str
    ) -> Dict[str, Any]:
        """
        Compare two versions of a dataset.

        Args:
            dataset_name: Name of the dataset
            version1: First version to compare
            version2: Second version to compare

        Returns:
            Dictionary containing comparison results
        """
        logger.info(f"Comparing dataset versions: {version1} and {version2}")
        
        versions = self.get_dataset_versions(dataset_name)
        
        v1 = next((v for v in versions if v["dataset_version"] == version1), None)
        v2 = next((v for v in versions if v["dataset_version"] == version2), None)
        
        if not v1 or not v2:
            logger.warning(f"One or both versions not found: {version1}, {version2}")
            return {"error": "One or both versions not found"}
        
        comparison = {
            "dataset_name": dataset_name,
            "version1": version1,
            "version2": version2,
            "size_diff": v2["dataset_size"] - v1["dataset_size"],
            "size_diff_percent": (v2["dataset_size"] - v1["dataset_size"]) / v1["dataset_size"] * 100 if v1["dataset_size"] > 0 else 0,
            "hash_diff": v1["dataset_hash"] != v2["dataset_hash"],
            "metrics_diff": {},
        }
        
        for metric, value in v1["metrics"].items():
            if metric in v2["metrics"]:
                comparison["metrics_diff"][metric] = v2["metrics"][metric] - value
        
        logger.info(f"Comparison results: {comparison}")
        return comparison

    def create_dataset_version(
        self,
        dataset_path: str,
        dataset_name: str,
        increment_type: str = "patch",
    ) -> str:
        """
        Create a new version of a dataset.

        Args:
            dataset_path: Path to the dataset file
            dataset_name: Name of the dataset
            increment_type: Type of version increment (major, minor, patch)

        Returns:
            MLflow run ID
        """
        logger.info(f"Creating new version of dataset: {dataset_name}")
        
        latest_version = self.get_latest_dataset_version(dataset_name)
        
        if latest_version:
            current_version = latest_version["dataset_version"]
            major, minor, patch = map(int, current_version.split("."))
            
            if increment_type == "major":
                new_version = f"{major + 1}.0.0"
            elif increment_type == "minor":
                new_version = f"{major}.{minor + 1}.0"
            else:  # patch
                new_version = f"{major}.{minor}.{patch + 1}"
        else:
            new_version = "1.0.0"
        
        logger.info(f"New version: {new_version}")
        
        run_id = self.log_dataset_to_mlflow(
            dataset_path=dataset_path,
            dataset_name=dataset_name,
            dataset_version=new_version,
            tags={"increment_type": increment_type},
        )
        
        return run_id


def main():
    """
    Main function for testing the DatasetVersioner.
    """
    versioner = DatasetVersioner(
        input_dir="data",
        mlflow_tracking_uri="http://mlflow-server.mlflow.svc.cluster.local:5000",
        experiment_name="llama4-fine-tuning",
    )
    
    run_id = versioner.log_dataset_to_mlflow(
        dataset_path="data/github_gitee_issues.json",
        dataset_name="github-gitee-issues",
        dataset_version="1.0.0",
        tags={"source": "github-gitee", "purpose": "fine-tuning"},
    )
    
    logger.info(f"Dataset logged to MLflow: {run_id}")
    
    versions = versioner.get_dataset_versions("github-gitee-issues")
    logger.info(f"Dataset versions: {versions}")
    
    latest_version = versioner.get_latest_dataset_version("github-gitee-issues")
    logger.info(f"Latest dataset version: {latest_version}")


if __name__ == "__main__":
    main()
